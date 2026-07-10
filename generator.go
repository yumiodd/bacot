package bacot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Struktur data untuk mem-parsing JSON
type DictionaryData struct {
	Dictionary []Entry `json:"dictionary"`
}

type Entry struct {
	Word string `json:"word"`
}

func main() {
	// 1. Membaca file JSON sumber
	filePath := "dictionary__JSON.json"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Gagal membaca file %s: %v", filePath, err)
	}

	// 2. Melakukan Unmarshal data JSON
	var data DictionaryData
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		log.Fatalf("Gagal melakukan unmarshal JSON: %v", err)
	}

	// 3. Mengumpulkan kata-kata unik menggunakan map
	wordMap := make(map[string]struct{})
	for _, entry := range data.Dictionary {
		words := strings.Split(entry.Word, " ")
		for _, w := range words {
			cleanWord := strings.TrimSpace(w)
			// Memastikan tidak menyimpan string kosong atau tanda strip tunggal
			if cleanWord != "" && cleanWord != "-" {
				wordMap[cleanWord] = struct{}{}
			}
		}
	}

	// 4. Membuat file .go baru untuk menampung hasil
	outFileName := "generated_dictionary.go"
	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("Gagal membuat file %s: %v", outFileName, err)
	}
	defer outFile.Close() // Pastikan file ditutup setelah selesai

	// 5. Menulis struktur kode Go ke dalam file tersebut
	// Kamu bisa mengganti "package main" menjadi package modulmu, misalnya "package bacot"
	fmt.Fprintln(outFile, "package main")
	fmt.Fprintln(outFile)
	fmt.Fprintln(outFile, "// KamusKata adalah map yang di-generate otomatis dari file JSON")
	fmt.Fprintln(outFile, "var KamusKata = map[string]struct{}{")

	// Iterasi map unik dan tulis ke dalam file sebagai key dari map[string]struct{}
	for word := range wordMap {
		// Menggunakan %q agar string aman dibungkus tanda kutip (escaped)
		fmt.Fprintf(outFile, "\t%q: {},\n", word)
	}

	// Menutup deklarasi map
	fmt.Fprintln(outFile, "}")

	// Menampilkan status berhasil di terminal
	fmt.Printf("Sukses! Berhasil men-generate file '%s' dengan %d kata unik.\n", outFileName, len(wordMap))
}
