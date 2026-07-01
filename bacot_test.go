package bacot

import (
	"fmt"
	"testing"
	"time"
)

func TestBacotSingleWord(t *testing.T) {

	startTime := time.Now()
	b := New()
	b.withExactWord = false
	r := b.Scan("babi")
	endTime := time.Now()
	f := endTime.Sub(startTime)
	fmt.Println("time:", f.String())

	if len(r.Words) == 0 {
		t.Fatal("expect found word but result none")
	}

	if r.Words[0].Word != "babi" {
		t.Fatal("expect found 'babi' but result ", r.Words[0].Word)
		return
	}
}
