package detector

import (
	"bacot/dictionary"
)

type Detector struct {
	withStemming  bool
	withQuickScan bool

	dict *dictionary.Dictionary
}

func New(dict *dictionary.Dictionary) *Detector {
	new := &Detector{
		withStemming: true,
		dict:         dict}
	return new
}
