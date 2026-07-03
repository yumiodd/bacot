package bacot

import "strings"

type WordIndext struct {
	Word  string
	Start int
	End   int
}

type Result struct {
	Words []*WordIndext
	Error error
	Text  string
}

func (r *Result) Censor() string {
	if len(r.Words) == 0 {
		return ""
	}

	st := []rune(r.Text)
	for _, w := range r.Words {
		for i := w.Start; i <= w.End; i++ {
			if st[i] == ' ' {
				continue
			}

			st[i] = '*'
		}
	}

	return string(st)
}

func (r *Result) Extract() []string {
	if len(r.Words) == 0 {
		return []string{}
	}

	var words []string
	for _, w := range r.Words {
		words = append(words, w.Word)
	}
	return words
}

func (r *Result) Result() *Result {
	return r
}

func (r *Result) IsToxic() bool {
	return len(r.Words) > 0
}

func (r *Result) First() string {
	if len(r.Words) > 0 {
		return r.Words[0].Word
	}
	return ""
}

func (r *Result) Last() string {
	if len(r.Words) > 0 {
		return r.Words[len(r.Words)-1].Word
	}
	return ""
}

func (r *Result) SpotLight() string {
	if len(r.Words) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, v := range r.Text {
		if v == ' ' {
			sb.WriteRune(' ')
			continue
		}
		sb.WriteRune('_')
	}

	res := []rune(sb.String())
	for _, w := range r.Words {
		for i := w.Start; i <= w.End; i++ {
			res[i] = rune(r.Text[i])
		}
	}

	return string(res)
}
