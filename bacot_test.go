package bacot

import (
	"fmt"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

var b = bacot.New()

func TestBacotLeetSpeak(t *testing.T) {

	fmt.Println(
		b.Text("4nj!n6").
			SanitizeDuplicateChar(true).
			WithLeetSpeak(true).
			RecursiveScan().
			CensoredText(),
	)
}
