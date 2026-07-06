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

	// This type of scan is easy if you use sanitize space
	temp := ms.withSanitizeSpace
	ms.WithSanitizeSpace(true)

	res := &ScanResult{
		text:        ms.text,
		praScanText: ms.generateText(),
	}

	ms.withSanitizeSpace = temp

	var (
		s        = res.praScanText
		finished = false
	)
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
					Start: l,
					End:   l + r - 1})

				if ms.collect {
					break
				}

				finished = true
			}
		}
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

// func (b *Bacot) Scan(s string) *Result {

// 	if b.withSanitizeSpace {

// 		var (
// 			sanitizeString strings.Builder
// 			spaces         = make(map[int]rune, 100)
// 		)
// 		for i, c := range s {
// 			if c == ' ' {
// 				spaces[i] = c
// 				continue
// 			}
// 			sanitizeString.WriteRune(c)
// 		}

// 		res := b.scanning(sanitizeString.String())
// 		res.spaceIndex = spaces
// 		return res
// 	}

// 	if b.withExactWord {
// 		// TODO: create a list that is not just space
// 		var (
// 			words = strings.Split(s, " ")
// 			idx   = 0
// 			res   = &Result{}
// 		)
// 		for _, w := range words {
// 			if r := b.scanning(w); len(r.Words) > 0 {
// 				r.Words[0].Start = idx
// 				r.Words[0].End += idx
// 				res.Words = append(res.Words, r.Words[0])

// 				// if not using quick scan
// 				if !(b.withCompound) {
// 					break
// 				}
// 			}
// 			idx += len(w) + 1
// 		}
// 		return res
// 	}

// 	res := b.scanning(s)

// 	return res
// }

// func (b *Bacot) scanning(s string) *Result {

// 	result := &Result{Text: s}

// 	s = strings.ToLower(s)
// 	lenS := len(s)

// 	if b.withExactWord && !(slices.Contains(b.wordCount, lenS)) {
// 		return result
// 	}

// 	// level 1 scanning (exact length)
// 	if slices.Contains(b.wordCount, lenS) {
// 		if b.Dict.Contains(s) {
// 			return &Result{
// 				Words: []*WordIndex{{s, 0, lenS - 1}},
// 				Error: nil,
// 			}
// 		}
// 	}

// 	// level 1.5 scanning for indonesian langguange
// 	// with pre-scanning (stemming)
// 	// {...}

// 	// level 2 scanning (sub-string)
// 	// using left-right cursor
// 	if !(b.withExactWord) {
// 		for l := 0; l <= lenS; l++ {

// 			sub := s[l:]
// 			for _, r := range b.wordCount {

// 				if r > len(sub) {
// 					break
// 				}

// 				word := sub[:r]

// 				// if using trim space,
// 				// even if you use a word whose value s is already free of space characters, it is handled by func Scan() above
// 				if b.Dict.Contains(word) {

// 					result.Words = append(result.Words, &WordIndex{word, l, l + r - 1})
// 					if b.withCompound {
// 						break
// 					}

// 					return result
// 				}
// 			}
// 		}
// 	}

// 	return result
// }
