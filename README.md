
![bacot.png](./bacot.png)
-

## Usage

Kesusahan buat daftar kata-kata yang dibatasi dalam bahasa Indonesia? Pake ini.

```
package main

import (
	"fmt"

	bacot "github.com/yumiodd/bacot/src"
)

func main() {

	// Inisial main struct
	b := bacot.New()

	// Buat modal scan
	text1 := b.Text("hallo")
	text2 := b.Text("anjing")
	text3 := b.Text("menganjing")
	text4 := b.Text("ngampret")
	text5 := b.Text("memukul")

	fmt.Println(text1.Scan().IsProfine()) // output: false
	fmt.Println(text2.Scan().IsProfine()) // output: true
	fmt.Println(text3.Scan().IsProfine()) // output: true
	fmt.Println(text4.Scan().IsProfine()) // output: true
	fmt.Println(text5.Scan().IsProfine()) // output: false

	// Tambah kata
	b = b.AddWord(true, "pukul")
	fmt.Println(text5.Scan().IsProfine()) // output: true

}
```
'memukul' tidak ada dalam dictionary tapi digenerate dari kata 'pukul'. Daftar kata yang disediakan kebanyakan adalah kata dasar, semua variasi
kata akan digenerate saat `bacot.New()`. Kamu bisa ambil daftar mapnya di `Dict.GetDict() `

```
func main() {

	b := bacot.New()
	fmt.Println(b.Dict.GetDict())   // output: malas, coba sendiri aja.
}
```

## Pre-scan Options
Beberapa opsi yang bisa di pakai untuk mengubah prilaku scan:
```
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {

    // initial
    b := bacot.New()

    text := b.Text("dasar kamu babi")

    // Menghilangkan whitespace
    text = text.WithSanitizeWhiteSpace(true)

    // Menormalisasikan semua leetSpeak karakter
    text = text.WithNormalizeLeetSpeak(true)

    // Mengumpulkan semua kata yang terdapat dalam kamus
    text = text.Collect()

    // Menggunakan Stemming, digunakan sebagai default
    text = text.WithStemming()

}

```
Semua opsi diatas bisa digunakan untuk semua jenis scan, `Scan() ` dan `RecursiveScan()`. 

Untuk whitespace kami memiliki daftar default, namun bisa kamu custom di saat pembuatan New(), tapi masih dalam tahap pengembangan.

## Hasil Penyaringan
hasil penyaringan, dapat di proses dalam berbagai rupa, di dalam module sudah di sediakan dan akan di tambah seiring kebutuhan dan dalam pengambangan.
```
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {


    res := bacot.New().Text("dasar babi anjing").Scan()

    // Memerika apakah mengandung kata yang dibatasi
    fmt.Println(res.IsToxic())                  // output: true

    // Mengambil kata pertama dan terakhir jumpa
    fmt.Println(res.First())                    // output: "babi"
    fmt.Println(res.Last())                     // output: "anjig"

    // Menyensor kata yang dibatasi dalam kalimat
    fmt.Println(res.CensorText())               // output: "dasar **** ******"

    // Menghitung kata yang dibatasi
    fmt.Println(res.CountFoundWord())           // output: 2

    // Mengekstrak kata yang dibatasi yang ditemukan
    words := res.Extract()                      // output: []string{"babi", "anjing"}

    // Mengabil data mentah
    wordIndexes := res.GetFountWords()          // output: []*WordIndex{}

    // Membuat word generator, mirip fungsi yield di python
    gen := res.Generator()

    fmt.Println(gen.Yield().Word())              // output : "babi"   
    fmt.Println(gen.Yield().Word())              // output : "anjing"   
}

``` 
## Pengembangan
Seperti module pada umumnya, module ini open source saya harap para module dev yang tertarik dapat berkonstribusi dan membuat module ini lebih baik, bahkan untuk yang tidak memahami coding masih bisa melakukan konstribusi dengan melakukan update pada kamus perkataan kasar, saya harap module ini berguna dan bisa dipakai oleh orang-orang yang memiliki kesulitan yang serupa.