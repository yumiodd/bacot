package bacot

import (
	"slices"
	"strings"
)

func (b *Bacot) Scan(s string) *Result {

	res := b.scanning(s)

	return res
}

func (b *Bacot) scanning(s string) *Result {

	result := &Result{}

	s = strings.ToLower(s)
	lenS := len(s)

	if b.withExactWord && !(slices.Contains(b.wordCount, lenS)) {
		return result
	}

	// level 1 scanning (exact length)
	if slices.Contains(b.wordCount, lenS) {
		if b.Dict.Contains(s) {
			return &Result{
				Words: []*WordIndext{{s, 0, lenS - 1}},
				Error: nil,
			}
		}
	}

	// level 1.5 scanning for indonesian langguange
	// with pre-scanning (stemming)
	// {...}

	// level 2 scanning (sub-string)
	// using left-right cursor
	if !(b.withExactWord) {
		for l := 0; l <= lenS; l++ {

			sub := s[l:]
			for _, r := range b.wordCount {

				if r > len(sub) {
					break
				}

				word := sub[:r]
				if b.Dict.Contains(word) {

					result.Words = append(result.Words, &WordIndext{word, l, l + r})
					if b.withCompound {
						continue
					}

					return result
				}
			}
		}
	}

	return result
}
