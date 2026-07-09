package bacot

import (
	"strings"
)

type Bacot struct {
	modalScan             *ModalScan
	customModalScanConfig *ModalScan

	// modal scanning
	Dict *Dictionary
}

func New() *Bacot {
	return &Bacot{Dict: NewDictionary()}
}

func (b *Bacot) Text(s string) *ModalScan {
	b.modalScan = &ModalScan{
		affix: true,
		dict:  b.Dict,
		input: s,
		text:  strings.ToLower(s),
	}

	// Default settings
	b.modalScan.
		WithLeetSpeak().
		UnstackChar().
		Affix(true)

	return b.modalScan
}

func (b *Bacot) AddWord(affix bool, words ...string) *Bacot {

	if affix {
		var s []string
		for _, w := range words {
			s = append(s, craftMan(w)...)
		}

		b.Dict.AddWords(s...)

	} else {
		b.Dict.AddWords(words...)
	}

	b.Dict.setUp()
	return b
}
