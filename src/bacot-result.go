package bacot

import "strings"

type WordIndex struct {
	Word  string
	Start int
	End   int
}

type WordIndexGenerator struct {
	buff []*WordIndex
	cur  int
}

func (wig *WordIndexGenerator) Yield() *WordIndex {
	if len(wig.buff) == 0 {
		return nil
	}

	// fallback for out of range index
	if wig.cur > len(wig.buff)-1 {
		return &WordIndex{}
	}

	wig.cur++
	return wig.buff[wig.cur]
}

type ScanResult struct {
	specialCharIndex map[int]rune

	// text is the original string of input and will not change
	// text ini akan di proses dan nilainya di assign ke praScanbaText
	text string

	// praScanText is the text resulting from the scan process, and this will be used using words
	praScanText string
	words       []*WordIndex
}

func (sr *ScanResult) GetText() string {
	return sr.text
}

func (sr *ScanResult) GetDetectWords() []string {
	var words []string
	for _, w := range sr.words {
		words = append(words, w.Word)
	}
	return words
}

func (sr *ScanResult) CensoredText() string {

	c := []rune(sr.praScanText)
	for _, w := range sr.words {
		for i := w.Start; i <= w.End; i++ {
			c[i] = '*'
		}
	}

	if sr.text == sr.praScanText {
		return string(c)
	}

	var (
		sb   strings.Builder
		diff = 0
		r    = []rune(c)
	)
	for i, s := range sr.text {
		if _, ok := whiteSpaces[s]; ok {
			sb.WriteRune(s)
			diff++
			continue
		}
		if s != r[i-diff] {

			if v, ok := simpleLeetSpeaks[s]; ok && v == c[i-diff] {
				sb.WriteRune(c[i-diff])
				continue
			} else {
				sb.WriteRune('*')
				diff++
				continue
			}
		}
		sb.WriteRune(c[i-diff])
	}

	return sb.String()
}

func (sr *ScanResult) IsToxic() bool {
	return len(sr.words) > 0
}

func (sr *ScanResult) First() string {
	if len(sr.words) == 0 {
		return ""
	}

	return sr.words[0].Word
}

func (sr *ScanResult) Last() string {
	if len(sr.words) == 0 {
		return ""
	}

	return sr.words[len(sr.words)-1].Word
}

func (sr *ScanResult) WordGenerator() *WordIndexGenerator {
	return &WordIndexGenerator{
		buff: sr.words,
		cur:  -1,
	}
}

func (sr *ScanResult) GetFoundWord() []*WordIndex {
	return sr.words
}
func (sr *ScanResult) CountFoundWord() int {
	return len(sr.words)
}

func (sr *ScanResult) Extract() []string {
	var words []string
	for _, w := range sr.words {
		words = append(words, w.Word)
	}
	return words
}
