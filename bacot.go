package bacot

import (
	"bacot/dictionary"
	"bacot/internal/core/detector"
)

type Bacot struct {
	dictionary *dictionary.Dictionary
	detector   *detector.Detector
}

func New() *Bacot {
	bacot := &Bacot{dictionary: dictionary.New()}
	bacot.detector = detector.New(bacot.dictionary)
	return bacot

}

func (b *Bacot) Tokenizer(text string) []string {
	return b.detector.Tokenizer(text)
}

func (b *Bacot) IsToxic(s string) bool {
	return b.detector.Scan(s)
}

func (b *Bacot) WithCompound(compound bool) {
	b.detector.WithCompound = compound
}
