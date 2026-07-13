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
	PrefixNasalFusion bool
}

type Dictionary struct {
	words          DictWords
	stops          DictWords
	falsePositives DictWords

	min       int
	max       int
	majority  int
	wordCount []int

	histogramFrequentWordLen map[int]int
}

func (d *Dictionary) Min() int {
	return d.min
}
func (d *Dictionary) Max() int {
	return d.max
}
func (d *Dictionary) Majority() int {
	return d.majority
}
func (d *Dictionary) IsContainLen(n int) bool {
	return slices.Contains(d.wordCount, n)
}
func (d *Dictionary) GetWordsLen() []int {
	return d.wordCount
}

func NewDictionary() *Dictionary {

	var words []string
	for _, w := range badwords {
		words = append(words, craftMan(w)...)
	}
	new := &Dictionary{
		words:                    NewDictWord(words...),
		stops:                    NewDictWord(defaultStopWords...),
		falsePositives:           NewDictWord(falsePositives...),
		histogramFrequentWordLen: map[int]int{},
	}

	new.counting()

	return new
}

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

		_, ok := d.histogramFrequentWordLen[lw]
		if ok {
			d.histogramFrequentWordLen[lw] += 1
		} else {
			d.histogramFrequentWordLen[lw] = 1
		}
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWordLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}
	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWordLen))
	slices.Sort(d.wordCount)
	d.majority = common

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
		_, ok := d.histogramFrequentWordLen[lw]
		if ok {
			d.histogramFrequentWordLen[lw] += 1
		} else {
			d.histogramFrequentWordLen[lw] = 1
		}

		// add the word
		d.words[word] = struct{}{}
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWordLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}

	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWordLen))
	slices.Sort(d.wordCount)
	d.majority = common
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

		d.histogramFrequentWordLen[lenW] -= 1
	}

	common := 0
	maxCount := 0
	for k, v := range d.histogramFrequentWordLen {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}

	d.wordCount = slices.Collect(maps.Keys(d.histogramFrequentWordLen))
	slices.Sort(d.wordCount)

	if len(d.wordCount) == 0 {
		d.min = 0
		d.max = 0
		d.majority = 0
	} else {
		d.min = d.wordCount[0]
		d.max = d.wordCount[len(d.wordCount)-1]
	}
	d.majority = common
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

func (d *Dictionary) AddFalsePositive(words ...string) {
	for _, w := range words {
		d.falsePositives[strings.ToLower(w)] = struct{}{}
	}
}
