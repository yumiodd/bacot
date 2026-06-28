package detector

import (
	"html"
	"strings"
)

func (d *Detector) Tokenizer(text string) []string {

	text = strings.ToLower(text)
	text = html.UnescapeString(text)
	text = rxURL.ReplaceAllString(text, "")
	text = rxEmail.ReplaceAllString(text, "")
	text = rxTwitter.ReplaceAllString(text, "")
	text = rxEscapeStr.ReplaceAllString(text, "")
	text = rxSymbol.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	return strings.Fields(text)
}
