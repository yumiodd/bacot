package bacot

import (
	"fmt"
	"slices"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

func TestTrimSpace(t *testing.T) {

	b := bacot.New()

	text := b.Raw("     halo hhhhhhh hhhh h h h h h h     h").TrimSpace().GetText()

	fmt.Println(text)

	if text != " halo hhhhhhh hhhh h h h h h h h" {
		t.Fatal("output not match")
	}

}

func TestReplaceWhiteSpace(t *testing.T) {

	b := bacot.New()

	txt := b.Raw("\tmakasi banget loh \rya \v").ReplaceWhiteSpace().GetText()
	fmt.Println(txt)

	if txt != " makasi banget loh  ya  " {
		t.Fatal("output not match")
	}
}
func TestSanitizeReadSign(t *testing.T) {

	b := bacot.New()

	txt := b.Raw("jallan-jalan. ke kota|padang, berk^ota k3 sur@baya").SanitazeReadSign().GetText()
	fmt.Println(txt)

	if txt != "jallan jalan  ke kota padang  berkota k3 sur@baya" {
		t.Fatal("output not match")
	}
}
func TestS(t *testing.T) {

	b := bacot.New().Dict.GetDict()

	fmt.Println("string[]{")
	var s []string
	for k := range b {
		s = append(s, k)
	}
	slices.Sort(s)
	for _, k := range s {
		fmt.Printf("\"%s\",\n", k)
	}
	fmt.Println("}")

}
