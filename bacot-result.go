package bacot

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
		for i := w.Start; i < w.End; i++ {
			if st[i] == ' ' {
				continue
			}

			st[i] = '*'
		}
	}

	return string(st)
}

func (r *Result) GetWords() []string {
	if len(r.Words) == 0 {
		return []string{}
	}

	var words []string
	for _, w := range r.Words {
		words = append(words, w.Word)
	}
	return words
}
