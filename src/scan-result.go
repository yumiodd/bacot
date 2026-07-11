package bacot

import (
	"math"
	"strings"
)

type WordIndex struct {
	Word  string
	Start int
	End   int
}

type WordIndexGenerator struct {
	buff []*WordIndex
	cur  int
}

func (wig *WordIndexGenerator) Yield() *WordIndex {
	if len(wig.buff) == 0 {
		return nil
	}

	wig.cur++
	if wig.cur >= len(wig.buff) {
		return nil
	}

	return wig.buff[wig.cur]
}

// ScanResult adalah hasil dari proses scanning.
// Ia menyimpan dua versi teks:
//   - text: teks asli input (preserve huruf besar/kecil)
//   - praScanText: teks setelah preprocessing (lowercased + sanitized)
//
// WordIndex menyimpan posisi start/end tiap kata kotor yang ditemukan,
// memungkinkan operasi seperti Censor() tanpa perlu scanning ulang.
type ScanResult struct {
	specialCharIndex map[int]rune

	text string

	// praScanText adalah teks setelah preprocessing (lowercased, sanitasi).
	// Posisi di WordIndex mengacu ke index rune di praScanText.
	praScanText string
	words       []*WordIndex
}

func (sr *ScanResult) GetText() string {
	return sr.text
}

// Censor() mengganti character kata kotor dengan '*'.
// Arsitektur: menggunakan precomputed WordIndex dari Scan().
//
// Strategi:
//   1. Replace langsung di []rune praScanText menggunakan
//      WordIndex (O(1) per kata, tanpa perlu scanning ulang).
//   2. Jika teks asli punya huruf besar/kecil, petakan kembali
//      dari teks asli ke teks yang sudah dicensor.
//   3. Jika teks sudah lowercase semua, tidak perlu mapping.
func (sr *ScanResult) Censor() string {

	if sr.Count() == 0 {
		return sr.text
	}

	c := []rune(sr.praScanText)
	for _, w := range sr.words {
		for i := w.Start; i <= w.End; i++ {
			c[i] = '*'
		}
	}

	if strings.ToLower(sr.text) == sr.praScanText {
		return string(c)
	}

	var (
		sb   strings.Builder
		diff = 0
		r    = []rune(c)
	)
	for i, s := range sr.text {

		if s == ' ' {
			sb.WriteRune(s)
			diff++
			continue
		}

		if s != r[i-diff] {
			if math.Abs(float64(s-r[i-diff])) == 32 {
				sb.WriteRune(s)
				continue
			}
			sb.WriteRune(r[i-diff])
		}
	}

	return sb.String()
}

func (sr *ScanResult) IsProfane() bool {
	return len(sr.words) > 0
}

func (sr *ScanResult) First() string {
	if len(sr.words) == 0 {
		return ""
	}

	return sr.words[0].Word
}

func (sr *ScanResult) Last() string {
	if len(sr.words) == 0 {
		return ""
	}

	return sr.words[len(sr.words)-1].Word
}

func (sr *ScanResult) Generator() *WordIndexGenerator {
	return &WordIndexGenerator{
		buff: sr.words,
		cur:  -1,
	}
}

func (sr *ScanResult) Count() int {
	return len(sr.words)
}

func (sr *ScanResult) Extract() []string {
	var words []string
	for _, w := range sr.words {
		words = append(words, w.Word)
	}
	return words
}
