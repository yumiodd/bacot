package bacot

import (
	"slices"
	"strings"
)

func (b *Bacot) Scan(s string) *Result {

	if b.withExactWord {
		// TODO: create a list that is not just space
		words := strings.Split(s, " ")
		res := &Result{Text: s}
		for i, w := range words {
			if r := b.scanning(w); len(r.Words) > 0 {
				r.Words[0].Start += i
				r.Words[0].End += i
				res.Words = append(res.Words, r.Words[0])

				// if not using quick scan
				if !(b.withCompound) {
					break
				}
			}
		}
		return res
	}

	res := b.scanning(s)

	return res
}

func (b *Bacot) scanning(s string) *Result {

	result := &Result{Text: s}

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

				// if using trim space,
				// even if you use a word whose value s is already free of space characters, it is handled by func Scan() above
				if !(b.withExactWord) {
					word = strings.ReplaceAll(word, " ", "")
				}

				if b.Dict.Contains(word) {

					result.Words = append(result.Words, &WordIndext{word, l, l + r})
					if b.withCompound {
						break
					}

					return result
				}
			}
		}
	}

	return result
}
