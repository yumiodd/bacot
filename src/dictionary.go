package bacot

import (
	"maps"
	"slices"
)

type DictWords = map[string]struct{}

func NewDictWord(words ...string) DictWords {
	new := map[string]struct{}{}
	for _, w := range words {
		new[w] = struct{}{}
	}
	return new
}

type Dictionary struct {
	badWords DictWords
	stops    DictWords

	min       int
	max       int
	majorty   int
	wordCount []int
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
func (d *Dictionary) IsContainLen(n int) bool {
	return slices.Contains(d.wordCount, n)
}
func (d *Dictionary) GetWordsLen() []int {
	return d.wordCount
}

func NewDictionary() *Dictionary {

	new := &Dictionary{
		badWords: NewDictWord(badwords...),
		stops:    NewDictWord(defaultStopWords...),
	}

	wordCount := map[int]int{}
	new.min = 99999
	for _, word := range slices.Collect(maps.Keys(new.badWords)) {
		lw := len(word)
		if lw > new.max {
			new.max = lw
		}
		if lw < new.min {
			new.min = lw
		}

		_, ok := wordCount[lw]
		if ok {
			wordCount[lw] += 1
		} else {
			wordCount[lw] = 1
		}
	}

	common := 0
	maxCount := 0
	for k, v := range wordCount {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}
	new.majorty = common

	new.wordCount = slices.Collect(maps.Keys(wordCount))
	slices.Sort(new.wordCount)

	return new
}

func (d *Dictionary) setUp() *Dictionary {

	var wordCount = map[int]int{}
	d.min = 99999
	for _, word := range slices.Collect(maps.Keys(d.badWords)) {
		lw := len(word)
		if lw > d.max {
			d.max = lw
		}
		if lw < d.min {
			d.min = lw
		}

		_, ok := wordCount[lw]
		if ok {
			wordCount[lw] += 1
		} else {
			wordCount[lw] = 1
		}
	}

	common := 0
	maxCount := 0
	for k, v := range wordCount {
		if v > maxCount {
			common = k
			maxCount = v
		}
	}
	d.wordCount = slices.Collect(maps.Keys(wordCount))
	slices.Sort(d.wordCount)
	d.majorty = common

	return d
}

func (d *Dictionary) AddWords(words ...string) {
	for _, w := range words {
		d.badWords[w] = struct{}{}
		d.setUp()
	}
}

func (d *Dictionary) DelWords(words ...string) {
	for _, w := range words {
		delete(d.badWords, w)
		d.setUp()
	}
}

// Remove both stop and badwords
func (d *Dictionary) Clear() *Dictionary {
	d.badWords = DictWords{}
	d.stops = DictWords{}
	d.setUp()
	return nil
}

func (d *Dictionary) Contains(word string) bool {
	if _, ok := d.badWords[word]; ok {
		return true
	}

	return false
}

func (d *Dictionary) GetDict() DictWords {
	return d.badWords
}

func (d *Dictionary) IsStopWord(s string) bool {
	if _, ok := d.stops[s]; ok {
		return true
	}

	return false
}
