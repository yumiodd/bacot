package bacot

import (
	"maps"
	"slices"
)

type Bacot struct {
	// scan behavior
	withStemming  bool
	withQuickScan bool
	withTrimSpace bool
	withExactWord bool // if false then it become very sensitive
	withCompound  bool
	withLeetSpeak bool

	// modal scanning
	minChar     int
	maxChar     int
	commonFound int
	wordCount   []int

	Dict *Dictionary
}

func New() *Bacot {

	new := &Bacot{
		withStemming:  true,
		withQuickScan: true,
		withTrimSpace: false,
		withExactWord: true,
		withCompound:  false,
		Dict:          NewDictionary(),
	}

	var wordCount = map[int]int{}
	new.minChar = 99999
	for _, word := range slices.Collect(maps.Keys(new.Dict.badWords)) {
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

	return new

}
