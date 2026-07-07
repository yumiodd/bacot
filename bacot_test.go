package bacot

import (
	"fmt"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

func TestBacotDetectInOneWord(t *testing.T) {

	b := bacot.New()
	res := b.Text("asu").Scan()

	if !res.IsToxic() {
		t.Fatal("expect toxic")
	}

	if f := res.CountFoundWord(); f != 1 {
		t.Fatal("expect found 1 got:", f)
	}

	if f := res.First(); f != "asu" {
		t.Fatal("expect 'asu' got: ", f)
	}
}

func TestBacotDetectOneWordInMultiWord(t *testing.T) {

	b := bacot.New()

	res := b.Text("dasar kamu asu ya").Scan()

	if !res.IsToxic() {
		t.Fatal("expect toxic")
	}

	if f := res.CountFoundWord(); f != 1 {
		t.Fatal("expect found 1 got:", f)
	}

	if f := res.First(); f != "asu" {
		t.Fatal("expect 'asu' got: ", f)
	}
}

func TestBacotNotDetectInOneWord(t *testing.T) {

	b := bacot.New()

	res := b.Text("hallo").Scan()

	if res.IsToxic() {
		t.Fatal("not toxic")
	}

	if f := res.CountFoundWord(); f != 0 {
		t.Fatal("expect found 0 got:", f)
	}

	if f := res.First(); f != "" {
		t.Fatal("expect empty string got: ", f)
	}
}

func TestBacotNotDetectInOneMultiWords(t *testing.T) {

	b := bacot.New()

	res := b.Text("hallo semua nama saya yumi").Scan()

	if res.IsToxic() {
		t.Fatal("not toxic")
	}

	if f := res.CountFoundWord(); f != 0 {
		t.Fatal("expect found 0 got:", f)
	}

	if f := res.First(); f != "" {
		t.Fatal("expect empty string got: ", f)
	}
}

func TestBacotDetectCollection(t *testing.T) {

	b := bacot.New()

	res := b.Text("anjing asu babi puki kimak bajingan").Collect(true).Scan()

	if f := res.CountFoundWord(); f != 6 {
		t.Fatal("expected found 6 got:", f)
	}

	var (
		test = []string{"anjing", "asu", "babi", "puki", "kimak", "bajingan"}
		gen  = res.WordGenerator()
	)

	for _, e := range test {
		g := gen.Yield()
		if e != g.Word {
			t.Fatalf("expect word %s got: %s", e, g.Word)
		}
	}
}

func TestBacotDetectOneWordWithSanitizeSpace(t *testing.T) {

	b := bacot.New()

	res := b.Text("a su").WithSanitizeSpace(true).Scan()

	if f := res.CountFoundWord(); f != 1 {
		t.Fatal("expected found 1 got:", f)
	}

	if f := res.First(); f != "asu" {
		t.Fatal("expected found `asu` got:", f)
	}
}

func TestBacotCensoringResult(t *testing.T) {

	b := bacot.New()

	res := b.Text("asu").Scan()

	if f := res.CountFoundWord(); f != 1 {
		t.Fatal("expected found 1 got:", f)
	}

	e := "***"
	if s := res.CensoredText(); s != e {
		t.Fatalf("expected found \"%s\" got: %s", e, s)
	}
}

func TestBacotCensoringResultMultiWord(t *testing.T) {

	b := bacot.New()

	res := b.Text("asu anjing babi").Collect(true).Scan()

	if f := res.CountFoundWord(); f != 3 {
		t.Fatal("expected found 3 got:", f)
	}

	e := "*** ****** ****"
	if s := res.CensoredText(); s != e {
		t.Fatalf("expected found \"%s\" got: %s", e, s)
	}
}

func TestBacotDetectWithSanitizeWhiteSpace(t *testing.T) {

	b := bacot.New()

	res := b.Text(" a\"su").WithSanitizeSpace(true).Collect(true).Scan()

	n := 1
	if f := res.CountFoundWord(); f != n {
		t.Fatalf("expected found %d got: %d", n, f)
	}

	e := "asu"
	if f := res.First(); f != e {
		t.Fatalf("expected found `%s` got: %s", e, f)
	}
}

func TestBacotDetectUsingRecursiveScan(t *testing.T) {

	b := bacot.New()

	res := b.Text(" pasu").RecursiveScan()

	expInt := 1
	if f := res.CountFoundWord(); f != expInt {
		t.Fatalf("expected found %d got: %d", expInt, f)
	}

	expStr := "asu"
	if f := res.First(); f != expStr {
		t.Fatalf("expected found \"%s\" got: %s", expStr, f)
	}

	expStr = " p***"
	if s := res.CensoredText(); s != expStr {
		t.Fatalf("expected found \"%s\" got: %s", expStr, s)
	}
}

func TestBacotDetectUsingRecursiveScanWithCollect(t *testing.T) {

	b := bacot.New()

	res := b.Text(" pasu anjingm babingsa").Collect(true).RecursiveScan()

	expInt := 3
	if f := res.CountFoundWord(); f != expInt {
		t.Fatalf("expected found %d got: %d", expInt, f)
	}

	gen := res.WordGenerator()
	expStrArr := []string{"asu", "anjing", "babi"}
	for _, e := range expStrArr {
		if f := gen.Yield().Word; f != e {
			t.Fatalf("expected found \"%s\" got: %s", e, f)
		}
	}

	expStr := " p*** ******m ****ngsa"
	if s := res.CensoredText(); s != expStr {
		t.Fatalf("expected found \"%s\" got: %s", expStr, s)
	}
}

func TestBacotDetectRecursiveScanWithSanitizeWhiteSpaceFalse(t *testing.T) {

	b := bacot.New()

	res := b.Text("  an jingm babingsa").Collect(true).RecursiveScan()

	e := []string{"babi"}
	if f := res.CountFoundWord(); f != len(e) {
		fmt.Println(res.Extract())
		t.Fatalf("expected found %d got: %d", len(e), f)
	}
}

func TestBacotDetectRecursiveScanWithSanitizeWhiteSpaceTrue(t *testing.T) {

	b := bacot.New()

	res := b.Text("  an jingm babingsa").WithSanitizeSpace(true).Collect(true).RecursiveScan()

	e := []string{"anjing", "babi"}
	if f := res.CountFoundWord(); f != len(e) {
		t.Fatalf("expected found %d got: %d", len(e), f)
	}

	eCensore := "  ** ****m ****ngsa"
	if s := res.CensoredText(); s != eCensore {
		t.Fatalf("expected %s got: %s", eCensore, s)
	}
}
