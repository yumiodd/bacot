package bacot

import (
	"strings"
)

type SanitizeOrder int

const (
	SanitizeNewLine SanitizeOrder = iota
	ClearSpace
	WithLeetSpeak
	UnstackChar
	// new
	TrimSpce
	ReplaceWhiteSpace
	SanitazeReadSign
)

type ModalScanConfig struct {
	Affix   bool
	Collect bool

	Order []SanitizeOrder
}

type Bacot struct {
	modalScan             *ModalScan
	customModalScanConfig *ModalScanConfig

	// modal scanning
	Dict *Dictionary
}

func New() *Bacot {
	return &Bacot{Dict: NewDictionary()}
}

// dengan Defoul setting
func (b *Bacot) Text(s string) *ModalScan {
	b.modalScan = &ModalScan{
		affix: true,
		dict:  b.Dict,
		input: s,
		text:  strings.ToLower(s),
	}

	if b.customModalScanConfig != nil {

		c := b.customModalScanConfig

		if c.Affix {
			b.modalScan.affix = c.Affix
		}
		if c.Collect {
			b.modalScan.collect = c.Collect
		}

		for _, f := range c.Order {
			switch f {
			case SanitizeNewLine:
				b.modalScan.SanitizeNewLine()
			case ClearSpace:
				b.modalScan.ClearSpace()
			case WithLeetSpeak:
				b.modalScan.WithLeetSpeak()
			case UnstackChar:
				b.modalScan.UnstackChar()
			case TrimSpce:
				b.modalScan.TrimSpace()
			case SanitazeReadSign:
				b.modalScan.SanitazeReadSign()
			case ReplaceWhiteSpace:
				b.modalScan.ReplaceWhiteSpace()
			}
		}

		return b.modalScan
	}

	// Default settings
	b.modalScan.
		SanitizeEmoji().
		ReplaceWhiteSpace().
		SanitazeReadSign().
		ReplaceWhiteSpace().
		UnstackChar().
		Affix(true)

	return b.modalScan
}

// Mentah, setting sendiri
func (b *Bacot) Raw(s string) *ModalScan {
	return &ModalScan{
		affix: false,
		dict:  b.Dict,
		input: s,
		text:  strings.ToLower(s)}
}

// Config
func (b *Bacot) Config(config *ModalScanConfig) *Bacot {
	b.customModalScanConfig = config
	return b
}

func (b *Bacot) AddWord(words ...string) *Bacot {
	b.Dict.AddWords(words...)
	return b
}

func (b *Bacot) DelWord(words ...string) *Bacot {
	b.Dict.DelWords(words...)
	return b
}
