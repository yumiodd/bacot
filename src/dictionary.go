package bacot

import (
	"maps"
	"slices"
	"strings"
)

type DictWords = map[string]struct{}

func NewDictWord(words ...string) DictWords {
	new := make(map[string]struct{}, 5000)
	for _, w := range words {
		w = strings.ToLower(w)
		new[w] = struct{}{}
	}
	return new
}

type CraftManConfig struct {
	PrefixNassalFusion bool
}

type Dictionary struct {
	words          DictWords
	stops          DictWords
	falsePositives DictWords

	min       int
	max       int
	majorty   int
	wordCount []int

	histogramFrequentWorldLen map[int]int
}

func (d *Dictionary) Min() int {
	return d.min
}
func (d *Dictionary) Max() int {
	return d.max
}
func (d *Dictionary) Majorty() int {
	return d.majorty
}
// IsContainLen adalah pre-filter optimization.
// Mengecek apakah ada kata di dictionary dengan panjang n.
// Dipanggil Scan() sebelum map lookup untuk skip token
// yang pasti tidak ada di dictionary.
func (d *Dictionary) IsContainLen(n int) bool {
	return slices.Contains(d.wordCount, n)
}
func (d *Dictionary) GetWordsLen() []int {
	return d.wordCount
}

// NewDictionary menginisialisasi dictionary dengan:
//   1. Setiap kata dasar dijalankan melalui craftMan() untuk
//      menghasilkan semua varian imbuhan (me-, pe-, di-, dll).
//   2. Stop words dan false positive dictionary di-load.
//   3. counting() membangun histogram panjang kata untuk
//      optimasi IsContainLen.
func NewDictionary() *Dictionary {

	var words []string
	for _, w := range badwords {
		words = append(words, craftMan(w)...)
	}
	new := &Dictionary{
		words:                     NewDictWord(words...),
		stops:                     NewDictWord(defaultStopWords...),
		falsePositives:            NewDictWord(falsePositives...),
		histogramFrequentWorldLen: map[int]int{},
	}

	new.counting()

	return new
}

// counting() membangun histogram panjang kata dari dictionary.
// Data ini digunakan untuk:
//   1. IsContainLen() — pre-filter token berdasarkan panjang
//   2. GetWordsLen() — daftar panjang yang tersedia untuk sliding window
//   3. Majorty() — panjang kata yang paling umum (debug/analisis)
//
// Optimization: histogram ini mencegah map lookup yang tidak perlu
// di Scan() dengan skip token jika panjangnya tidak ada di dictionary.
func (d *Dictionary) counting() *Dictionary {

	d.min = 99999
	for _, word := range slices.Collect(maps.Keys(d.words)) {
		lw := len(word)
		if lw > d.max {
			d.max = lw
		}
		if lw < d.min {
			d.min = lw
		}

		_, ok := d.histogramFrequentWorldLen[lw]
		if ok {
			d.histogramFrequentWorldLen[lw] += 1
		} else {
			d.histogramFrequentWorldLen[lw] = 1
		}
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWorldLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}
	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWorldLen))
	slices.Sort(d.wordCount)
	d.majorty = common

	return d
}

func (d *Dictionary) AddWords(words ...string) {
	if len(words) == 0 {
		return
	}

	var add []string
	for _, w := range words {
		add = append(add, craftMan(w)...)
	}

	for _, word := range slices.Collect(maps.Keys(NewDictWord(add...))) {

		// Min Max
		lw := len(word)
		if lw > d.max {
			d.max = lw
		}
		if lw < d.min {
			d.min = lw
		}

		// Add frequent
		_, ok := d.histogramFrequentWorldLen[lw]
		if ok {
			d.histogramFrequentWorldLen[lw] += 1
		} else {
			d.histogramFrequentWorldLen[lw] = 1
		}

		// add the word
		d.words[word] = struct{}{}
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWorldLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}

	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWorldLen))
	slices.Sort(d.wordCount)
	d.majorty = common
}

func (d *Dictionary) DelWords(words ...string) {
	if len(words) == 0 {
		return
	}

	for _, word := range words {
		lenW := len(word)

		if _, ok := d.words[word]; ok {
			delete(d.words, word)
		} else {
			continue
		}

		d.histogramFrequentWorldLen[lenW] -= 1
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWorldLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}

	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWorldLen))
	slices.Sort(d.wordCount)

	if len(d.wordCount) == 0 {
		d.min = 0
		d.max = 0
		d.majorty = 0
	} else {
		d.min = d.wordCount[0]
		d.max = d.wordCount[len(d.wordCount)-1]
	}
	d.majorty = common
}

func (d *Dictionary) Contains(word string) bool {
	if _, ok := d.words[word]; ok {
		return true
	}

	return false
}

func (d *Dictionary) GetDict() DictWords {

	var words = DictWords{}
	maps.Copy(words, d.words)

	return words
}

func (d *Dictionary) IsStopWord(s string) bool {
	if _, ok := d.stops[s]; ok {
		return true
	}

	return false
}

func (d *Dictionary) IsFalsePositive(s string) bool {
	if _, ok := d.falsePositives[s]; ok {
		return true
	}

	return false
}
