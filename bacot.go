package bacot

import (
	"maps"
	"slices"
)

type Bacot struct {

	// Perform steming on words that contain affixes,
	// specifically for words from Indonesian
	// feature on under development
	// default true
	withStemming bool

	// The sentence will be broken down into tokens and each word will be scanned
	// if you use this then withTrimSpace is automatically turned off
	// if it is turned off then the word search will be very sensitive, sub-strings of the word will be checked too
	// default true
	//
	// # some people prefer to turn it off but there are also those who are bothered by it,
	// # well I also get distracted sometimes so I turn it on as default if you want to turn it off you can do it manually
	withExactWord bool

	// Option to continue searching until the end of the given text even if the word has already been found,
	// Using this will automatically turn off WithQuickScan
	// default false
	withCompound bool

	// Option replaces each character with an alphabet, for example: @ -> a, 3 -> e
	// default true
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

func (b *Bacot) WithStemming(v bool) *Bacot {
	b.withStemming = v
	return b
}

func (b *Bacot) WithExactWord() *Bacot {
	b.withExactWord = true
	return b
}

func (b *Bacot) WithTrimSpace() *Bacot {
	b.withExactWord = false
	return b
}

func (b *Bacot) WithQuickScan() *Bacot {
	b.withCompound = true
	return b
}
func (b *Bacot) WithCompound() *Bacot {
	b.withCompound = false
	return b
}
