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
