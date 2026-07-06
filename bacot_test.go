package bacot

import (
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
