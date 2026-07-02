package bacot

import (
	bacot "bacot/src"
	"fmt"
	"testing"
	"time"
)

func SpeedTest(f func() *bacot.Result) *bacot.Result {
	startTime := time.Now()
	r := f()
	endTime := time.Now()
	final := endTime.Sub(startTime)
	fmt.Println("time:", final.String())
	return r
}

// func TestBacotSingleWord(t *testing.T) {

// 	r := SpeedTest(func() *Result {
// 		b := New()
// 		b.withExactWord = false
// 		return b.Scan("asdfasdfaBabiABdb")
// 	})

// 	if len(r.Words) == 0 {
// 		t.Fatal("expect found word but result none")
// 	}

// 	if r.Words[0].Word != "babi" {
// 		t.Fatal("expect found 'babi' but result ", r.Words[0].Word)
// 		return
// 	}
// }

// func TestBacotTrimSpace(t *testing.T) {

// 	r := SpeedTest(func() *Result {
// 		b := New()
// 		b.withExactWord = false
// 		b.withTrimSpace = true
// 		b.withCompound = true
// 		return b.Scan("asdfasdfaBa biAdab an jing")
// 	})

// 	if len(r.Words) == 0 {
// 		t.Fatal("expect found word but result none")
// 	}

// 	if r.Words[0].Word != "babi" {
// 		t.Fatal("expect found 'babi' but result ", r.Words[0].Word)
// 		return
// 	}

// 	fmt.Println(r.Censor())

// }

func TestSpotLight(t *testing.T) {

	r := SpeedTest(func() *bacot.Result {
		b := bacot.New().WithTrimSpace().WithCompound()
		return b.Scan("asdfbb ab Iadsfasdf")
	})

	fmt.Println("dari r:", r.SpotLight())
}
