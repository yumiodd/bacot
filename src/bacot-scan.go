package bacot

import (
	"slices"
	"strings"
)

type ModalScan struct {
	collect bool // def: false
	affix   bool // def: true

	// temp, upper layer set by (bacot)
	input string
	text  string
	dict  *Dictionary
}

func (ms *ModalScan) Collect(v bool) *ModalScan {
	ms.collect = v
	return ms
}

func (ms *ModalScan) Affix(v bool) *ModalScan {
	ms.affix = v
	return ms
}

// per word scanning
func (ms *ModalScan) Scan() *ScanResult {

	res := &ScanResult{
		text:        ms.input,
		praScanText: ms.text,
	}

	var (
		words = strings.SplitSeq(res.praScanText, " ")
		idx   = 0
	)

	// example input : "word in sentence" or "word-join"
	for w := range words {
		var lenW = len(w)

		if !(ms.dict.IsContainLen(lenW)) || !(ms.affix && !ms.dict.IsStopWord(w)) {
			continue
		}

		var found = ms.dict.Contains(w)

		// stemming
		if !found && ms.affix && (lenW > 3) && slices.Contains(prefixes, w[:2]) {

			// recursive scan / sliding window
			for l := 0; l <= len(w); l++ {
				sub := w[l:]

				for _, r := range ms.dict.GetWordsLen() {
					if r > len(sub) {
						break
					}

					word := sub[:r]
					if ms.dict.Contains(word) {
						found = true
						break
					}
				}
			}
		}

		if found {
			res.words = append(res.words,
				&WordIndex{
					Word:  w,
					Start: idx,
					End:   len(w) + idx - 1,
				})

			if !(ms.collect) {
				break
			}
		}

		idx += len(w) + 1
	}

	return res
}

// Reursive will scan in sub-string fashion
func (ms *ModalScan) RecursiveScan() *ScanResult {

	res := &ScanResult{
		text:        ms.input,
		praScanText: ms.text,
	}

	var (
		finished = false
		idx      = 0
	)
	for s := range strings.SplitSeq(res.praScanText, " ") {
		if finished {
			break
		}

		if s == "" {
			idx++
			continue
		}

		for l := 0; l <= len(s); l++ {
			if finished {
				break
			}

			sub := s[l:]
			for _, r := range ms.dict.GetWordsLen() {

				if r > len(sub) {
					break
				}

				word := sub[:r]
				if ms.dict.Contains(word) {
					res.words = append(res.words, &WordIndex{
						Word:  word,
						Start: l + idx,
						End:   idx + l + r - 1,
					})

					if ms.collect {
						break
					}
					finished = true
				}
			}
		}
		idx += len(s) + 1
	}

	return res
}

func (ms *ModalScan) UnstackChar() *ModalScan {
	var (
		sb   strings.Builder
		prev rune
	)

	for _, c := range ms.text {
		if c == prev {
			continue
		} else {
			sb.WriteRune(c)
		}
		prev = c
	}

	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) TrimSpace() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if c == ' ' {
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) WithLeetSpeak() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if r, ok := simpleLeetSpeaks[c]; ok {
			sb.WriteRune(r)
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) GetText() string {
	return ms.text
}
