package bacot

import (
	"fmt"
	"slices"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

func Test(t *testing.T) {
	b := bacot.New()

	fmt.Println(b.Text("babi123").Collect(true).Scan().Extract())
}

func TestNew(t *testing.T) {
	b := bacot.New()
	if b == nil {
		t.Fatal("New() returned nil")
	}
	if b.Dict == nil {
		t.Fatal("Dict is nil")
	}
	if b.Dict.Min() == 0 || b.Dict.Max() == 0 {
		t.Error("Dict Min/Max should be non-zero after initialization")
	}
	if len(b.Dict.GetWordsLen()) == 0 {
		t.Error("GetWordsLen() should not be empty")
	}
}

func TestTextDefaults(t *testing.T) {
	b := bacot.New()
	ms := b.Text("Hello World")
	if ms == nil {
		t.Fatal("Text() returned nil")
	}

	got := ms.GetText()
	// "hello" → UnstackChar removes consecutive 'l' → "helo"
	if got != "helo world" {
		t.Errorf("GetText() = %q, want %q", got, "helo world")
	}
}

func TestTextProcessesLeetSpeak(t *testing.T) {
	b := bacot.New()
	ms := b.Text("4njing").WithLeetSpeak()
	got := ms.GetText()
	if got != "anjing" {
		t.Errorf("GetText() after leet = %q, want %q", got, "anjing")
	}
}

func TestTextUnstacksChars(t *testing.T) {
	b := bacot.New()
	ms := b.Text("baaabiii")
	got := ms.GetText()
	// "baaabiii" → removes consecutive 'a' and 'i' → "babi"
	if got != "babi" {
		t.Errorf("GetText() after unstack = %q, want %q", got, "babi")
	}
}

func TestTextLowercase(t *testing.T) {
	b := bacot.New()
	ms := b.Text("ANJING")
	got := ms.GetText()
	if got != "anjing" {
		t.Errorf("GetText() = %q, want %q", got, "anjing")
	}
}

func TestAddWordNoAffix(t *testing.T) {
	b := bacot.New()
	b.AddWord("xyzword")
	if !b.Dict.Contains("xyzword") {
		t.Error("Dict should contain 'xyzword' after AddWord(false, ...)")
	}
}

func TestAddWordWithAffix(t *testing.T) {
	b := bacot.New()
	b.AddWord("babi")
	added := b.Dict.Contains("babi")
	if !added {
		t.Error("Dict should contain base word 'babi' after AddWord(true, 'babi')")
	}
}

func TestAddWordChain(t *testing.T) {
	b := bacot.New()
	b.AddWord("foo").AddWord("bar")
	if !b.Dict.Contains("foo") || !b.Dict.Contains("bar") {
		t.Error("AddWord chaining should add all words")
	}
}

func TestAddWordScansCorrectly(t *testing.T) {
	b := bacot.New()
	b.AddWord("customword")
	res := b.Text("customword").Scan()
	if !res.IsProfane() {
		t.Error("Scan() should detect 'customword' after AddWord")
	}
}

func TestDictionaryContains(t *testing.T) {
	b := bacot.New()
	if !b.Dict.Contains("anjing") {
		t.Error("Dict should contain 'anjing'")
	}
	if !b.Dict.Contains("babi") {
		t.Error("Dict should contain 'babi'")
	}
	if b.Dict.Contains("nonexistent") {
		t.Error("Dict should not contain 'nonexistent'")
	}
}

func TestDictionaryStopWords(t *testing.T) {
	b := bacot.New()
	if !b.Dict.IsStopWord("dan") {
		t.Error("'dan' should be a stop word")
	}
	if !b.Dict.IsStopWord("di") {
		t.Error("'di' should be a stop word")
	}
	if b.Dict.IsStopWord("anjing") {
		t.Error("'anjing' should not be a stop word")
	}
}

func TestDictionaryAddWords(t *testing.T) {
	b := bacot.New()
	b.Dict.AddWords("testword1", "testword2")
	if !b.Dict.Contains("testword1") {
		t.Error("Dict should contain 'testword1' after AddWords")
	}
	if !b.Dict.Contains("testword2") {
		t.Error("Dict should contain 'testword2' after AddWords")
	}
}

func TestDictionaryDelWords(t *testing.T) {
	b := bacot.New()
	b.Dict.DelWords("anjing")
	if b.Dict.Contains("anjing") {
		t.Error("Dict should not contain 'anjing' after DelWords")
	}
}

func TestDictionaryGetDict(t *testing.T) {
	b := bacot.New()
	dict := b.Dict.GetDict()
	if dict == nil {
		t.Error("GetDict() should not return nil")
	}
	if _, ok := dict["anjing"]; !ok {
		t.Error("GetDict() should contain 'anjing'")
	}
}

func TestDictionaryMetrics(t *testing.T) {
	b := bacot.New()
	if b.Dict.Min() <= 0 {
		t.Error("Min() should be positive")
	}
	if b.Dict.Max() < b.Dict.Min() {
		t.Error("Max() should be >= Min()")
	}
	if b.Dict.Majorty() <= 0 {
		t.Error("Majorty() should be positive")
	}
	if len(b.Dict.GetWordsLen()) == 0 {
		t.Error("GetWordsLen() should not be empty")
	}
	if !b.Dict.IsContainLen(b.Dict.Max()) {
		t.Error("IsContainLen(Max()) should be true")
	}
	if b.Dict.IsContainLen(999) {
		t.Error("IsContainLen(999) should be false")
	}
}

func TestDictionaryWordsLenSorted(t *testing.T) {
	b := bacot.New()
	lengths := b.Dict.GetWordsLen()
	if !slices.IsSorted(lengths) {
		t.Error("GetWordsLen() should return sorted slice")
	}
}

func TestScanDetectsProfanity(t *testing.T) {
	b := bacot.New()
	res := b.Text("anjing").Scan()
	if !res.IsProfane() {
		t.Error("Scan() should detect 'anjing'")
	}
}

func TestScanCleanText(t *testing.T) {
	b := bacot.New()
	res := b.Text("hari ini cuaca cerah").Scan()
	if res.IsProfane() {
		t.Error("Scan() should not detect profanity in clean text")
	}
}

func TestScanCollectSingleDefault(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Scan()
	if !res.IsProfane() {
		t.Fatal("Scan() should detect profanity")
	}
	if res.Count() != 1 {
		t.Errorf("With default Collect=false, Count() = %d, want 1", res.Count())
	}
}

func TestScanCollectAll(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	if !res.IsProfane() {
		t.Fatal("Scan() should detect profanity")
	}
	if res.Count() < 2 {
		t.Errorf("With Collect=true, Count() = %d, want at least 2", res.Count())
	}
}

func TestScanDetectsMultipleProfanities(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing asu").Collect(true).Scan()
	if res.Count() < 3 {
		t.Errorf("Should detect all 3 profanities, got %d", res.Count())
	}
}

func TestScanRespectsAffixDisabled(t *testing.T) {
	b := bacot.New()
	ms := b.Text("mebabi").Affix(false)
	got := ms.GetText()
	if got != "mebabi" {
		t.Fatalf("GetText() = %q", got)
	}
	res := ms.Scan()
	if res.IsProfane() {
		t.Error("With Affix=false, 'mebabi' should not be detected")
	}
}

func TestScanAffixEnabled(t *testing.T) {
	b := bacot.New()
	res := b.Text("mebabi").Scan()
	if !res.IsProfane() {
		t.Error("With Affix=true, 'mebabi' should be detected (contains 'babi')")
	}
}

func TestScanAffixWithVariousPrefixes(t *testing.T) {
	b := bacot.New()

	detected := []struct {
		input string
		word  string
		note  string
	}{
		{"mebabi", "babi", "non-vowel-start, prefix me-"},
		{"pebabi", "babi", "non-vowel-start, prefix pe-"},
		{"terasu", "asu", "vowel-start, 'r' from ter-"},
		{"berasu", "asu", "vowel-start, 'r' from ber-"},
		{"ngasu", "asu", "vowel-start, 'g' from ng-"},
		{"mengasu", "asu", "vowel-start, 'g' from meng-"},
		{"pengasu", "asu", "vowel-start, 'g' from peng-"},
	}
	for _, tc := range detected {
		t.Run(tc.input+"_detected", func(t *testing.T) {
			res := b.Text(tc.input).Scan()
			if !res.IsProfane() {
				t.Errorf("'%s' should be detected (%s)", tc.input, tc.note)
			}
		})
	}

	notDetected := []struct {
		input string
		word  string
		note  string
	}{
		{"diasu", "asu", "vowel-start, 'i' is not g/r"},
		{"teasu", "asu", "vowel-start, 'e' is not g/r"},
		{"beasu", "asu", "vowel-start, 'e' is not g/r"},
		{"menyasu", "asu", "vowel-start, 'y' from meny-, 'y' is not g/r"},
		{"penyasu", "asu", "vowel-start, 'y' from peny-, 'y' is not g/r"},
	}
	for _, tc := range notDetected {
		t.Run(tc.input+"_not_detected", func(t *testing.T) {
			res := b.Text(tc.input).Scan()
			if res.IsProfane() {
				t.Errorf("'%s' should NOT be detected (%s)", tc.input, tc.note)
			}
		})
	}
}

func TestScanLeetSpeakDetection(t *testing.T) {
	b := bacot.New()
	cases := []struct {
		input string
		word  string
	}{
		{"4njing", "anjing"},
		{"@njing", "anjing"},
		{"8abi", "babi"},
		{"4su", "asu"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			res := b.Text(tc.input).WithLeetSpeak().Scan()
			if !res.IsProfane() {
				t.Errorf("Leet input '%s' should be detected as '%s'", tc.input, tc.word)
			}
		})
	}
}

func TestScanUnstackCharDetection(t *testing.T) {
	b := bacot.New()
	cases := []string{
		"anjiiing",
		"baaabi",
		"assuu",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			res := b.Text(input).Scan()
			if !res.IsProfane() {
				t.Errorf("Unstacked input '%s' should be detected", input)
			}
		})
	}
}

func TestScanTrimSpace(t *testing.T) {
	b := bacot.New()
	res := b.Text("a n j i n g").ClearSpace().Scan()
	if !res.IsProfane() {
		t.Error("TrimSpace input 'a n j i n g' should be detected as 'anjing'")
	}
}

func TestScanTrimSpaceAndLeet(t *testing.T) {
	b := bacot.New()
	res := b.Text("4 n j i n g").WithLeetSpeak().ClearSpace().Scan()
	if !res.IsProfane() {
		t.Error("TrimSpace+Leet input '4 n j i n g' should be detected as 'anjing'")
	}
}

func TestScanResultIsProfane(t *testing.T) {
	b := bacot.New()
	profane := b.Text("anjing").Scan()
	if !profane.IsProfane() {
		t.Error("IsProfane() should return true when profanity found")
	}
	clean := b.Text("halo").Scan()
	if clean.IsProfane() {
		t.Error("IsProfane() should return false for clean text")
	}
}

func TestScanResultFirst(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	if res.First() != "babi" {
		t.Errorf("First() = %q, want %q", res.First(), "babi")
	}
}

func TestScanResultFirstEmpty(t *testing.T) {
	b := bacot.New()
	res := b.Text("halo").Scan()
	if res.First() != "" {
		t.Errorf("First() for clean text should be empty, got %q", res.First())
	}
}

func TestScanResultLast(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	if res.Last() != "anjing" {
		t.Errorf("Last() = %q, want %q", res.Last(), "anjing")
	}
}

func TestScanResultLastEmpty(t *testing.T) {
	b := bacot.New()
	res := b.Text("halo").Scan()
	if res.Last() != "" {
		t.Errorf("Last() for clean text should be empty, got %q", res.Last())
	}
}

func TestScanResultCount(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing asu kontol").Collect(true).Scan()
	if res.Count() != 4 {
		t.Errorf("Count() = %d, want 4", res.Count())
	}
}

func TestScanResultExtract(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	extracted := res.Extract()
	if len(extracted) != 2 {
		t.Fatalf("Extract() length = %d, want 2", len(extracted))
	}
	if extracted[0] != "babi" || extracted[1] != "anjing" {
		t.Errorf("Extract() = %v, want [babi anjing]", extracted)
	}
}

func TestScanResultGetText(t *testing.T) {
	b := bacot.New()
	original := "Hello World"
	res := b.Text(original).Scan()
	if res.GetText() != original {
		t.Errorf("GetText() = %q, want %q", res.GetText(), original)
	}
}

func TestScanResultCensorBasic(t *testing.T) {
	b := bacot.New()
	res := b.Text("anjing").Scan()
	censored := res.Censor()
	if censored != "******" {
		t.Errorf("Censor() = %q, want %q", censored, "******")
	}
}

func TestScanResultCensorInSentence(t *testing.T) {
	b := bacot.New()
	res := b.Text("dia adalah anjing").Collect(true).Scan()
	censored := res.Censor()
	want := "dia adalah ******"
	if censored != want {
		t.Errorf("Censor() = %q, want %q", censored, want)
	}
}

func TestScanResultCensorMultiple(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi dan anjing").Collect(true).Scan()
	censored := res.Censor()
	want := "**** dan ******"
	if censored != want {
		t.Errorf("Censor() = %q, want %q", censored, want)
	}
}

func TestScanResultCensorLeetSpeak(t *testing.T) {
	b := bacot.New()
	res := b.Text("4njing").WithLeetSpeak().Scan()
	censored := res.Censor()
	// "4njing" → leet → "anjing" → censored → all asterisks mapped back
	if censored != "******" {
		t.Errorf("Censor() for leet input = %q, want %q", censored, "******")
	}
}

func TestScanResultCensorPartialProfanity(t *testing.T) {
	b := bacot.New()
	res := b.Text("aku suka babi").Collect(true).Scan()
	censored := res.Censor()
	want := "aku suka ****"
	if censored != want {
		t.Errorf("Censor() = %q, want %q", censored, want)
	}
}

func TestRecursiveScan(t *testing.T) {
	b := bacot.New()
	res := b.Text("anjing").RecursiveScan()
	if !res.IsProfane() {
		t.Error("RecursiveScan() should detect 'anjing'")
	}
}

func TestRecursiveScanCleanText(t *testing.T) {
	b := bacot.New()
	res := b.Text("halo dunia").RecursiveScan()
	if res.IsProfane() {
		t.Error("RecursiveScan() should not detect clean text")
	}
}

func TestRecursiveScanCollect(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).RecursiveScan()
	if res.Count() < 2 {
		t.Errorf("With Collect=true, RecursiveScan Count() = %d, want at least 2", res.Count())
	}
}

func TestRecursiveScanNoCollect(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").RecursiveScan()
	if res.Count() != 1 {
		t.Errorf("With default Collect=false, RecursiveScan Count() = %d, want 1", res.Count())
	}
}

func TestRecursiveScanSubstring(t *testing.T) {
	b := bacot.New()
	res := b.Text("xxbabi").RecursiveScan()
	if !res.IsProfane() {
		t.Error("RecursiveScan() should detect 'babi' substring in 'xxbabi'")
	}
}

func TestRecursiveScanMultipleSubstrings(t *testing.T) {
	b := bacot.New()
	res := b.Text("babianjing").Collect(true).RecursiveScan()
	if res.Count() < 2 {
		t.Errorf("RecursiveScan should detect multiple substrings in 'babianjing', got %d", res.Count())
	}
}

func TestRecursiveScanLeetSpeak(t *testing.T) {
	b := bacot.New()
	cases := []string{
		"4njing",
		"8abi",
		"4su",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			res := b.Text(input).WithLeetSpeak().RecursiveScan()
			if !res.IsProfane() {
				t.Errorf("RecursiveScan should detect leet input '%s'", input)
			}
		})
	}
}

func TestRecursiveScanUnstackChar(t *testing.T) {
	b := bacot.New()
	cases := []string{
		"anjiiing",
		"baaabi",
		"assuu",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			res := b.Text(input).RecursiveScan()
			if !res.IsProfane() {
				t.Errorf("RecursiveScan should detect unstacked input '%s'", input)
			}
		})
	}
}

func TestRecursiveScanLeetAndUnstack(t *testing.T) {
	b := bacot.New()
	res := b.Text("4njiiing").WithLeetSpeak().RecursiveScan()
	if !res.IsProfane() {
		t.Error("RecursiveScan should detect combined leet+unstack input '4njiiing'")
	}
}

func TestRecursiveScanWordIndexPositions(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi").RecursiveScan()
	gen := res.Generator()
	wi := gen.Yield()
	if wi == nil {
		t.Fatal("Generator should yield a word")
	}
	if wi.Word != "babi" {
		t.Errorf("Word = %q, want %q", wi.Word, "babi")
	}
	if wi.Start != 0 {
		t.Errorf("Start = %d, want 0", wi.Start)
	}
	if wi.End != 3 {
		t.Errorf("End = %d, want 3", wi.End)
	}
}

func TestRecursiveScanWordIndexPrefix(t *testing.T) {
	b := bacot.New()
	res := b.Text("xbabi").RecursiveScan()
	gen := res.Generator()
	wi := gen.Yield()
	if wi == nil {
		t.Fatal("Generator should yield a word")
	}
	if wi.Word != "babi" {
		t.Errorf("Word = %q, want %q", wi.Word, "babi")
	}
	if wi.Start != 1 {
		t.Errorf("Start = %d, want 1", wi.Start)
	}
	if wi.End != 4 {
		t.Errorf("End = %d, want 4", wi.End)
	}
}

func TestRecursiveScanWordIndexMultiWord(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).RecursiveScan()
	if res.Count() != 2 {
		t.Fatalf("Count = %d, want 2", res.Count())
	}
	x := res.Extract()
	if x[0] != "babi" {
		t.Errorf("First word = %q, want %q", x[0], "babi")
	}
	if x[1] != "anjing" {
		t.Errorf("Second word = %q, want %q", x[1], "anjing")
	}

	gen := res.Generator()
	wi := gen.Yield()
	if wi.Start != 0 || wi.End != 3 {
		t.Errorf("babi position: Start=%d End=%d, want 0 3", wi.Start, wi.End)
	}
	wi = gen.Yield()
	if wi.Start != 5 || wi.End != 10 {
		t.Errorf("anjing position: Start=%d End=%d, want 5 10", wi.Start, wi.End)
	}
}

func TestScanResultGenerator(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	gen := res.Generator()
	if gen == nil {
		t.Fatal("Generator() should not return nil")
	}
	if w := gen.Yield(); w == nil {
		t.Error("Generator should yield first word")
	} else if w.Word != "babi" {
		t.Errorf("First word = %q, want %q", w.Word, "babi")
	}
	if w := gen.Yield(); w == nil {
		t.Error("Generator should yield second word")
	} else if w.Word != "anjing" {
		t.Errorf("Second word = %q, want %q", w.Word, "anjing")
	}
	if w := gen.Yield(); w != nil {
		t.Error("Generator should be exhausted after 2 words")
	}
}

func TestScanResultGeneratorEmpty(t *testing.T) {
	b := bacot.New()
	res := b.Text("halo").Scan()
	gen := res.Generator()
	wi := gen.Yield()
	if wi != nil {
		t.Error("Generator should yield nil for empty result")
	}
}

func TestScanResultMethodsOnEmpty(t *testing.T) {
	b := bacot.New()
	res := b.Text("halo").Scan()
	if res.IsProfane() {
		t.Error("IsProfane() should be false")
	}
	if res.First() != "" {
		t.Errorf("First() = %q, want empty", res.First())
	}
	if res.Last() != "" {
		t.Errorf("Last() = %q, want empty", res.Last())
	}
	if res.Count() != 0 {
		t.Errorf("Count() = %d, want 0", res.Count())
	}
	if len(res.Extract()) != 0 {
		t.Errorf("Extract() should be empty, got %v", res.Extract())
	}
	if res.Censor() != "halo" {
		t.Errorf("Censor() for clean text = %q, want %q", res.Censor(), "halo")
	}
}

func TestEdgeCaseEmptyString(t *testing.T) {
	b := bacot.New()
	res := b.Text("").Scan()
	if res.IsProfane() {
		t.Error("Empty string should not be profane")
	}
	if res.Censor() != "" {
		t.Errorf("Censor() for empty = %q", res.Censor())
	}
}

func TestEdgeCaseOnlyStopWords(t *testing.T) {
	b := bacot.New()
	res := b.Text("dan atau di").Scan()
	if res.IsProfane() {
		t.Error("Text with only stop words should not be profane")
	}
}

func TestEdgeCaseMixedCase(t *testing.T) {
	b := bacot.New()
	res := b.Text("AnJiNg").Scan()
	if !res.IsProfane() {
		t.Error("Mixed case 'AnJiNg' should be detected (lowercased)")
	}
}

func TestEdgeCaseSpecialCharacters(t *testing.T) {
	b := bacot.New()
	// leet is off by default, '!' is not converted
	res := b.Text("anjing!").Scan()
	if res.IsProfane() {
		t.Error("'anjing!' should NOT be detected — '!' not converted by default, prevChar=' ' blocks vowel match")
	}
}

func TestEdgeCaseProfanityWithNumbers(t *testing.T) {
	b := bacot.New()
	// "babi123" → stem finds "babi" (len 4) within "babi123"
	// rest "123" has no vowels → isOneSyllable returns false → not filtered
	res := b.Text("babi123").Scan()
	if !res.IsProfane() {
		t.Error("'babi123' should be detected — 'babi' found via stem, rest '123' not one-syllable")
	}
}

func TestChainingFullFlow(t *testing.T) {
	b := bacot.New()
	res := b.Text("aku suka 4njing").
		WithLeetSpeak().
		Collect(true).
		Scan()
	if !res.IsProfane() {
		t.Fatal("Should detect profanity")
	}
	censored := res.Censor()
	if censored == "aku suka 4njing" {
		t.Error("Censor() should replace profanity with asterisks")
	}
	if res.Count() >= 1 {
		t.Logf("Detected: %v, censored: %q", res.Extract(), censored)
	}
}

func TestGetTextReturnsProcessedText(t *testing.T) {
	b := bacot.New()
	ms := b.Text("4njing")
	got := ms.GetText()
	// leet is off by default, "4njing" stays as-is
	if got != "4njing" {
		t.Errorf("GetText() = %q, want %q", got, "4njing")
	}
}

func TestModalScanCollectChaining(t *testing.T) {
	b := bacot.New()
	ms := b.Text("test").Collect(true)
	if ms == nil {
		t.Fatal("Collect() should return ModalScan for chaining")
	}
}

func TestModalScanAffixChaining(t *testing.T) {
	b := bacot.New()
	ms := b.Text("test").Affix(false)
	if ms == nil {
		t.Fatal("Affix() should return ModalScan for chaining")
	}
}

func TestModalScanUnstackCharChaining(t *testing.T) {
	b := bacot.New()
	ms := b.Text("test").UnstackChar()
	if ms == nil {
		t.Fatal("UnstackChar() should return ModalScan for chaining")
	}
}

func TestModalScanTrimSpaceChaining(t *testing.T) {
	b := bacot.New()
	ms := b.Text("test").ClearSpace()
	if ms == nil {
		t.Fatal("TrimSpace() should return ModalScan for chaining")
	}
}

func TestScanResultWordIndexPositions(t *testing.T) {
	b := bacot.New()
	res := b.Text("babi anjing").Collect(true).Scan()
	gen := res.Generator()
	wi := gen.Yield()
	if wi == nil {
		t.Fatal("Generator should yield first word")
	}
	if wi.Word != "babi" {
		t.Errorf("First word = %q, want %q", wi.Word, "babi")
	}
	if wi.Start != 0 {
		t.Errorf("First word Start = %d, want 0", wi.Start)
	}
}

func TestAddWordThenScanNoAffix(t *testing.T) {
	b := bacot.New()
	b.AddWord("custombad")
	res := b.Text("custombad").Scan()
	if !res.IsProfane() {
		t.Error("Scan should detect custom word added without affix")
	}
}

func TestDictionaryContainsLowercase(t *testing.T) {
	b := bacot.New()
	if !b.Dict.Contains("anjing") {
		t.Error("Dict.Contains('anjing') should be true")
	}
}

func TestDictionaryIsStopWordLowercase(t *testing.T) {
	b := bacot.New()
	if !b.Dict.IsStopWord("dan") {
		t.Error("Dict.IsStopWord('dan') should be true")
	}
}

func TestMultipleScansSameInstance(t *testing.T) {
	b := bacot.New()
	r1 := b.Text("babi").Scan()
	r2 := b.Text("halo").Scan()
	if !r1.IsProfane() {
		t.Error("First scan should detect profanity")
	}
	if r2.IsProfane() {
		t.Error("Second scan should not detect profanity in clean text")
	}
}

func TestAddWordNonAffixPreservesDictionaryLengths(t *testing.T) {
	b := bacot.New()
	before := len(b.Dict.GetWordsLen())
	b.AddWord("abcdefghij")
	after := len(b.Dict.GetWordsLen())
	if after < before {
		t.Error("Adding a word should not reduce word length count")
	}
}

func TestScanChainedConfigurations(t *testing.T) {
	b := bacot.New()
	res := b.Text("4 n j i n g").
		WithLeetSpeak().
		ClearSpace().
		UnstackChar().
		Collect(true).
		Scan()
	if !res.IsProfane() {
		t.Error("Chained configuration should detect profanity from leet+trim+unstack")
	}
}

func TestU(t *testing.T) {

	res := bacot.New().Text("bacotin").Collect(true).Scan()
	fmt.Println(res.Censor())
	fmt.Println(res.Extract())
}
