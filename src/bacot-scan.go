package bacot

import (
	"slices"
	"strings"
)

type ModalScan struct {
	// pre scan, user config
	withLeetSpeak     bool // true
	withReplaceSpace  bool // true
	withSanitizeSpace bool // false

	// in scan, user config
	collect               bool // false
	affix                 bool // true
	sanitizeDuplicateChar bool // true

	// temp, upper layer set (bacot)
	text string
	dict *Dictionary
}

func (ms *ModalScan) WithSanitizeSpace(v bool) *ModalScan {
	ms.withReplaceSpace = v
	return ms
}

func (ms *ModalScan) Collect(v bool) *ModalScan {
	ms.collect = v
	return ms
}

func (ms *ModalScan) WithAffix(v bool) *ModalScan {
	ms.affix = v
	return ms
}

func (ms *ModalScan) SanitizeDuplicateChar(v bool) *ModalScan {
	ms.sanitizeDuplicateChar = v
	return ms
}

func (ms *ModalScan) WithLeetSpeak(v bool) *ModalScan {
	ms.withLeetSpeak = v
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

// Reursive willl scann in sub-string fashion
func (ms *ModalScan) RecursiveScan() *ScanResult {

	res := &ScanResult{
		text:        ms.text,
		praScanText: ms.generateText(),
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

	if ms.withLeetSpeak ||
		ms.sanitizeDuplicateChar ||
		ms.withReplaceSpace ||
		ms.withSanitizeSpace {
		var sb strings.Builder

		for _, c := range ms.text {
			if ms.withReplaceSpace {
				if _, ok := whiteSpaces[c]; ok {
					sb.WriteRune(' ')
					continue
				}
			}

			if ms.withLeetSpeak {
				if replace, ok := simpleLeetSpeaks[c]; ok {
					sb.WriteRune(replace)
					continue
				}
			}

			sb.WriteRune(c)
		}

		if ms.sanitizeDuplicateChar {
			s := sanitizeDuplicateChar(sb.String())
			if s != sb.String() {
				sb.Reset()
				sb.WriteString(s)
			}
		}

		if ms.withSanitizeSpace {
			s := sanitizeSpace(sb.String())
			sb.Reset()
			sb.WriteString(s)
		}

		return sb.String()

	}
	return ms.text
}

func sanitizeDuplicateChar(s string) string {

	var (
		sb   strings.Builder
		prev rune
	)

	for _, c := range s {
		if c == prev {
			continue
		} else {
			sb.WriteRune(c)
		}
		prev = c
	}
	return sb.String()
}

func sanitizeSpace(s string) string {

	var sb strings.Builder
	for _, c := range s {
		if c == ' ' {
			sb.WriteRune(c)
			continue
		}
		sb.WriteRune(c)
	}

	return sb.String()
}
