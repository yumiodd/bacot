package bacot

import (
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

var bt = bacot.New()

func repeatWord(w string, n int) string {
	var out string
	for i := 0; i < n; i++ {
		out += w + " "
	}
	return out
}

func repeatJoined(s string, n int) string {
	var out string
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

// --------------------------------------------------
// RINGAN — single word, basic sentences, dictionary
// --------------------------------------------------

func BenchmarkContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Dict.Contains("anjing")
	}
}

func BenchmarkAddWord(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.AddWord("newbadword")
	}
}

func BenchmarkTextPreprocessing(b *testing.B) {
	input := "4njiiing d4n b4by d4r4nc4ng"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input)
	}
}

func BenchmarkScanSingleWord(b *testing.B) {
	input := "anjing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}

func BenchmarkScan(b *testing.B) {
	input := "aku suka 4njiiing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkScanCollect(b *testing.B) {
	input := "aku suka babi dan anjing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkRecursiveScan(b *testing.B) {
	input := "aku suka 4njiiing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).RecursiveScan()
	}
}

func BenchmarkCensor(b *testing.B) {
	res := bt.Text("aku suka 4njiiing").Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

func BenchmarkCensorSentence(b *testing.B) {
	res := bt.Text("babi dan anjing asu kontol").Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

// --------------------------------------------------
// SEDANG — paragraph, stop words, affix, leet, unstack
// --------------------------------------------------

var longParagraph = "aku suka 4njiiing dan b4by b4ngs4t itu sangat menjengkelkan karena mereka selalu b3g0 dan b0d0h " +
	"setiap hari mereka selalu ngomong kasar k4nt0l m3m3k dan p3c3k tanpa henti " +
	"ini adalah sebuah kalimat yang sangat panjang untuk menguji performa bacot dalam menangani input besar " +
	"dengan berbagai macam kata kasar b4bi anjing asu k0nt0l dan varian l33t speak lainnya " +
	"tolong jangan Tiru perilaku mereka yang suka mengumpat dan berkata kasar setiap saat"

func BenchmarkSedangScanParagraph(b *testing.B) {
	input := longParagraph
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkSedangRecursiveScanParagraph(b *testing.B) {
	input := longParagraph
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).RecursiveScan()
	}
}

func BenchmarkSedangCensorParagraph(b *testing.B) {
	res := bt.Text(longParagraph).Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

func BenchmarkSedangScanCleanLong(b *testing.B) {
	input := "hari ini cuaca sangat cerah dan indah sekali saya suka berjalan jalan di taman " +
		"bersama dengan teman teman sambil menikmati udara segar dan sinar matahari " +
		"semoga hari ini menjadi hari yang menyenangkan bagi kita semua tanpa ada masalah " +
		"kita bisa belajar bersama dan saling membantu satu sama lain"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}

func BenchmarkSedangScanManyStopWords(b *testing.B) {
	input := "dan di dengan untuk dan di dengan untuk dan di dengan untuk dan di dengan untuk " +
		"dan di dengan untuk dan di dengan untuk dan di dengan untuk dan di dengan untuk " +
		"dan di dengan untuk dan di dengan untuk dan di dengan untuk dan di dengan untuk " +
		"dan di dengan untuk dan di dengan untuk dan di dengan untuk dan di dengan untuk"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkSedangScanAffixLong(b *testing.B) {
	input := "mebabi pebabi tebabi dibabi beasu mebabi pebabi tebabi dibabi beasu " +
		"mebabi pebabi tebabi dibabi beasu mebabi pebabi tebabi dibabi beasu " +
		"mebabi pebabi tebabi dibabi beasu mebabi pebabi tebabi dibabi beasu"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}

func BenchmarkSedangScanAllLeet(b *testing.B) {
	input := "4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h 4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h " +
		"4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h 4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h " +
		"4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h 4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkSedangScanHeavyUnstack(b *testing.B) {
	input := "aaannnjjjiiinnnggg bbbaaabbbiiii aaassssuuu kkkkkoooonnnnttttooolll"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkSedangScanMixedAffixLeet(b *testing.B) {
	input := "m3b4bi p38abi m3n4njing p3b4ngs4t b34su m3k0nt0l"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}

func BenchmarkSedangRecursiveScanNoMatch(b *testing.B) {
	input := "superkalifragilistikekspialidocious"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).RecursiveScan()
	}
}

func BenchmarkSedangCensorLeet(b *testing.B) {
	res := bt.Text("4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h").Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

// --------------------------------------------------
// BERAT — ~10KB text with various mixes
// --------------------------------------------------

var tenKbParagraph = func() string {
	base := "4njing 8abi b4ngs4t k0nt0l b3g0 b0d0h m3m3k p3c3k " +
		"babi anjing asu kontol bangsat bego bodoh memek pecun " +
		"mebabi pebabi tebabi dibabi beasu menjing penganjing "
	return repeatWord(base, 100)
}()

func BenchmarkBeratScan10Kb(b *testing.B) {
	input := tenKbParagraph
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkBeratRecursiveScan10Kb(b *testing.B) {
	input := tenKbParagraph
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).RecursiveScan()
	}
}

func BenchmarkBeratCensor10Kb(b *testing.B) {
	res := bt.Text(tenKbParagraph).Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

var tenKbLeetText = repeatWord("4njing 8abi 4su k0nt0l b4ngs4t b3g0 b0d0h m3m3k p3c3k d4n b4by", 200)

func BenchmarkBeratScan10KbAllLeet(b *testing.B) {
	input := tenKbLeetText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkBeratTextPreprocessing10Kb(b *testing.B) {
	input := tenKbLeetText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input)
	}
}

var tenKbContinuous = repeatJoined("babi", 2500)

func BenchmarkBeratScan10KbContinuous(b *testing.B) {
	input := tenKbContinuous
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).ClearSpace().Scan()
	}
}

func BenchmarkBeratRecursiveScan10KbContinuous(b *testing.B) {
	input := tenKbContinuous
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).ClearSpace().RecursiveScan()
	}
}

// --------------------------------------------------
// SUPER BERAT — 100KB, worst-case single word
// --------------------------------------------------

var hundredKbLorem = repeatJoined(
	"lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore "+
		"et dolore magna aliqua ut enim ad minim veniam quis nostrud exercitation ullamco laboris nisi "+
		"ut aliquip ex ea commodo consequat duis aute irure dolor in reprehenderit in voluptate velit "+
		"esse cillum dolore eu fugiat nulla pariatur excepteur sint occaecat cupidatat non proident "+
		"sunt in culpa qui officia deserunt mollit anim id est laborum ",
	2000,
)

func BenchmarkSuperScan100Kb(b *testing.B) {
	input := hundredKbLorem
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}

func BenchmarkSuperRecursiveScan100Kb(b *testing.B) {
	input := hundredKbLorem
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).RecursiveScan()
	}
}

// --------------------------------------------------
// EXTREM — worst case: all profanity, 10K char single word
// --------------------------------------------------

var allProfanityText = repeatWord("anjing babi asu kontol bangsat bego bodoh memek", 500)

func BenchmarkExtremScanAllProfanity(b *testing.B) {
	input := allProfanityText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).Scan()
	}
}

func BenchmarkExtremRecursiveScanAllProfanity(b *testing.B) {
	input := allProfanityText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Collect(true).RecursiveScan()
	}
}

func BenchmarkExtremCensorAllProfanity(b *testing.B) {
	res := bt.Text(allProfanityText).Collect(true).Scan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.Censor()
	}
}

var tenKCharWord = func() string {
	return repeatJoined("abcdefghij", 995) + "babi"
}()

func BenchmarkExtremRecursiveScanLongSingleWord(b *testing.B) {
	input := tenKCharWord
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).RecursiveScan()
	}
}

func BenchmarkExtremScanLongSingleWord(b *testing.B) {
	input := tenKCharWord
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Text(input).Scan()
	}
}
