package bacot

import (
	"testing"
)

func TestBacot(t *testing.T) {

	text := "maling kundang"

	bacot := New()

	tokens := bacot.Tokenizer(text)

	if tokens[0] != "maling" {
		t.Error("failet")
	}
	if tokens[1] != "kundang" {
		t.Error("failet")
	}

}
