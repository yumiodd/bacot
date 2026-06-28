package bacot

import (
	"bacot/dictionary"
	"bacot/internal/core/detector"
)

type Bacot struct {
	dictionary *dictionary.Dictionary
	detector   *detector.Detector
}

func NewDefault() *Bacot {
	return NewWithCustomDictionary(dictionary.NewDefault())
}

// to change to another language, make sure all data from the dict uses that language.
func NewWithCustomDictionary(dict *dictionary.Dictionary) *Bacot {
	new := &Bacot{dictionary: dictionary.NewDefault()}
	new.detector = detector.New(new.dictionary)

	return new
}

var (
	stemmingDisable  = 0
	stemOnlyBadwords = 1 // default
	stemFull         = 2
)

func (b *Bacot) StemmingLevel(level int) *Bacot {
	switch level {
	case stemmingDisable:
		b.detector.WithStemming(false)
	case stemOnlyBadwords:
		b.dictionary.SetScanLevel(dictionary.ScanBadWords)
	case stemFull:
		b.dictionary.SetScanLevel(dictionary.ScanBoth)
	default:
		b.dictionary.SetScanLevel(dictionary.ScanBadWords)
	}
	return b
}

func (b *Bacot) AddBadWords(words ...string) *Bacot {
	b.dictionary.AddBadWords(words...)
	return b
}

func (b *Bacot) Stem(word string) string {
	return b.detector.Stem(word)
}
