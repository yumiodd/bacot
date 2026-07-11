package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

// data test:
// https://github.com/damzaky/kumpulan-kata-bahasa-indonesia-KBBI/blob/master/list_1.0.0.txt
func TestKosaKata(t *testing.T) {

	// Buka file data
	file, err := os.Open("./list_kosa_kata_100K.txt")
	if err != nil {
		t.Fatalf("cant open file: %s", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Buat file output
	output, err := os.Create("kosa_kata_test_output.csv")
	if err != nil {
		t.Fatalf("cant create file: %s", err.Error())
	}
	defer output.Close()
	header := "input|badword|got\n"
	n, err := output.WriteString(header) // ==> buat header
	if err != nil {
		t.Fatalf("cant write string into the output file: %s", err.Error())
	}
	if len(header) != n {
		t.Fatal("write string into the output not complete")
	}

	// Buat New Bacot
	b := bacot.New()

	var count = 0
	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			t.Fatalf("error when scan: %s", err.Error())
		}

		// Get data per baris dan scan dengan bacot
		input := scanner.Text()
		res := b.Text(input).Scan()

		if !res.IsProfane() {
			continue
		}

		// tulis output
		var sb strings.Builder
		sb.WriteString(input)
		sb.WriteRune('|')

		if res.IsProfane() {
			sb.WriteString("1|")
		} else {
			sb.WriteString("0|")
		}

		for _, s := range res.Extract() {
			sb.WriteString(s)
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')

		n, err := output.WriteString(sb.String())
		if err != nil {
			t.Fatalf("cant write string into the output file: %s", err.Error())
		}
		if len(sb.String()) != n {
			t.Fatal("write string into the output not complete")
		}

		count++
	}

	fmt.Println("file sudah siap")
	fmt.Println("kata terditeksi abbusive: ", count)
}

func TestFalsePositive(t *testing.T) {

	b := bacot.New()

	txt := b.Text("ancok")
	fmt.Println(txt.Scan().IsProfane())

}
