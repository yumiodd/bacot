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
		withLeetSpeak:         true,
		sanitizeDuplicateChar: true,
		affix:                 true,

		dict: b.Dict,
		text: strings.ToLower(s),
	}

	return b.modalScan
}
