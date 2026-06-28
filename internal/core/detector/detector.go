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

func (d *Detector) WithStemming(v bool) *Detector {
	d.withStemming = v
	return d
}

func (d *Detector) WithQuickScan(v bool) *Detector {
	d.withQuickScan = v
	return d
}

func (d *Detector) Scan(fulltext string) *ScanResult {

	tokens := d.Tokenizer(fulltext)

	result := &ScanResult{}

	// scan fulltext
	// {...}

	// scan per words
	for _, t := range tokens {

		word := d.Naturalize(t)
		if d.withStemming {
			word = d.Stem(t)
		}

		if d.dict.IsBadword(word) {
			result.Badwords = append(result.Badwords, word)
			if d.withQuickScan {
				return result
			}
			continue
		}

	}

	return result
}

func (d *Detector) Naturalize(word string) string {
	return ""
}
