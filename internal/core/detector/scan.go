package detector

import (
	"fmt"
	"slices"
	"strings"
)

func (d *Detector) Scan(s string) bool {

	s = strings.ToLower(s)
	lenS := len(s)

	if d.dict.IsStopWord(s) {
		return false
	}

	if d.withExactWord {
		if !(lenS <= d.maxChar) || !(slices.Contains(d.wordCount, lenS)) {
			return false
		}
	}

	// level 1 scanning (exact length)
	if slices.Contains(d.wordCount, lenS) {
		if d.dict.Contains(s) {
			return true
		}
	}

	// level 1.5 scanning for indonesian langguange
	// with pre-scanning (stemming)
	// {...}

	// leve 2 scanning (sub-string)
	// using left-right cursor
	var (
		startIdx = 0
		endIdx   = 0
		compound = map[string][]int{}
	)
	if !(d.withExactWord) {
		for l := 0; l <= lenS; l++ {

			sub := s[l:]
			for _, r := range d.wordCount {
				if r > len(sub) {
					break
				}

				word := sub[:r]
				if d.dict.Contains(word) {
					startIdx = l
					endIdx = l + r

					if d.WithCompound {
						compound[word] = []int{l, l + r}
						continue
					}

					fmt.Println("found in", s, "word at", startIdx, ":", endIdx)
					fmt.Println("word: ", word)
					return true
				}
			}
		}
	}

	if d.WithCompound && len(compound) > 0 {
		for k, v := range compound {
			fmt.Println("found in", s, "word at", v[0], ":", v[1])
			fmt.Println("word: ", k)
		}

		return true
	}

	return false
}
