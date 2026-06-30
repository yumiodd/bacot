package detector

import (
	"bacot/dictionary"
	"fmt"
	"maps"
	"slices"
)

type Detector struct {
	withStemming  bool
	withQuickScan bool
	withTrimSpace bool

	// config
	withExactWord bool // if false then it become very sensitive
	WithCompound  bool

	// modal scanning
	minChar     int
	maxChar     int
	commonFound int
	wordCount   []int

	dict *dictionary.Dictionary
}

func New(dict *dictionary.Dictionary) *Detector {
	new := &Detector{
		withStemming: true,
		dict:         dict}

	words := dict.GetDict()

	var wordCount = map[int]int{}
	new.minChar = 99999
	for _, word := range slices.Collect(maps.Keys(words)) {
		lw := len(word)
		if lw > new.maxChar {
			new.maxChar = lw
		}
		if lw < new.minChar {
			new.minChar = lw
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
	new.wordCount = slices.Collect(maps.Keys(wordCount))
	slices.Sort(new.wordCount)
	new.commonFound = common

	fmt.Println("Detector Config")
	fmt.Println("max", new.maxChar, wordCount[new.maxChar])
	fmt.Println("min", new.minChar, wordCount[new.minChar])
	fmt.Println("common", new.commonFound, wordCount[new.commonFound])
	fmt.Println("word count", new.wordCount)
	return new
}
