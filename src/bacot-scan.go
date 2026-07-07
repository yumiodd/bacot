package bacot

import (
	"strings"
)

type ModalScan struct {
	// pre scan, user config
	withReplaceSpecialChar bool // true
	withSanitizeSpace      bool // false

	// in scan, user config
	collect bool // false
	stem    bool // false

	// temp, upper layer set (bacot)
	text string
	dict *Dictionary
}

func (ms *ModalScan) WithReplaceSpecialChar(v bool) *ModalScan {
	ms.withReplaceSpecialChar = v
	return ms
}

func (ms *ModalScan) WithSanitizeSpace(v bool) *ModalScan {
	ms.withSanitizeSpace = v
	return ms
}

func (ms *ModalScan) Collect(v bool) *ModalScan {
	ms.collect = v
	return ms
}

func (ms *ModalScan) Stem(v bool) *ModalScan {
	ms.stem = v
	return ms
}

// Scan is per-word scanning
func (ms *ModalScan) Scan() *ScanResult {

	res := &ScanResult{
		text:        ms.text,
		praScanText: ms.generateText(),
	}

	var (
		words = strings.SplitSeq(res.praScanText, " ")
		idx   = 0
	)

	// example input : "word in sentence" or "word-join"
	for w := range words {
		var lenW = len(w)

		if !(ms.dict.IsContainLen(lenW)) {
			return res
		}

		if ms.dict.Contains(w) {
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

// Reursive willl scann in sub-string fashion
func (ms *ModalScan) RecursiveScan() *ScanResult {

	res := &ScanResult{
		text:        ms.text,
		praScanText: ms.text,
	}

	if ms.withSanitizeSpace {
		res.praScanText = ms.generateText()
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

func (ms *ModalScan) generateText() string {

	if ms.text == "" {
		return ""
	}

	var sb strings.Builder
	for _, c := range ms.text {
		if ms.withSanitizeSpace {
			if _, ok := whiteSpace[c]; ok {
				continue
			}
		}
		// if ms.withReplaceSpecialChar {
		// 	if sc, ok := specialChar[c]; ok {
		// 		ms.res.specialCharIndex[i] = c
		// 		c = sc
		// 	}
		// }

		sb.WriteRune(c)
	}
	return sb.String()
}
