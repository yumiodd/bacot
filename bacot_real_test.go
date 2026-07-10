package bacot

import (
	"slices"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

func TestU(t *testing.T) {
	bacot.New().Text("menganjing").Scan()
}

// SocialMediaComments — simulasi komentar dari media sosial
func TestSocialMediaComments(t *testing.T) {
	b := bacot.New()

	t.Run("dasar babi lu", func(t *testing.T) {
		res := b.Text("dasar babi lu").Collect(true).Scan()
		if !res.IsProfane() {
			t.Fatal("should detect 'babi'")
		}
		if !slices.Contains(res.Extract(), "babi") {
			t.Errorf("expected 'babi' in %v", res.Extract())
		}
	})

	t.Run("anjing lu babi", func(t *testing.T) {
		res := b.Text("anjing lu babi").Collect(true).Scan()
		if !res.IsProfane() {
			t.Fatal("should detect profanity")
		}
		if !slices.Contains(res.Extract(), "anjing") {
			t.Errorf("expected 'anjing' in %v", res.Extract())
		}
		if !slices.Contains(res.Extract(), "babi") {
			t.Errorf("expected 'babi' in %v", res.Extract())
		}
	})

	t.Run("goblok amat sih", func(t *testing.T) {
		res := b.Text("goblok amat sih").Scan()
		if !res.IsProfane() {
			t.Error("should detect 'goblok'")
		}
	})

	t.Run("kontol bangsat", func(t *testing.T) {
		res := b.Text("kontol bangsat").Collect(true).Scan()
		if !res.IsProfane() {
			t.Fatal("should detect profanity")
		}
		if !slices.Contains(res.Extract(), "kontol") {
			t.Errorf("expected 'kontol' in %v", res.Extract())
		}
		if !slices.Contains(res.Extract(), "bangsat") {
			t.Errorf("expected 'bangsat' in %v", res.Extract())
		}
	})

	t.Run("dasar memek", func(t *testing.T) {
		res := b.Text("dasar memek").Scan()
		if !res.IsProfane() {
			t.Error("should detect 'memek'")
		}
	})
}

// ChatMessages — simulasi percakapan sehari-hari
func TestChatMessages(t *testing.T) {
	b := bacot.New()

	clean := []struct {
		input string
		note  string
	}{
		{"lu di mana", "stop word 'di', clean"},
		{"udah makan belum", "percakapan normal"},
		{"dari tadi ngapain", "pertanyaan biasa"},
		{"oke siap", "konfirmasi"},
	}
	for _, tc := range clean {
		t.Run(tc.input, func(t *testing.T) {
			res := b.Text(tc.input).Scan()
			if res.IsProfane() {
				t.Errorf("%q should be clean (%s): got %v", tc.input, tc.note, res.Extract())
			}
		})
	}

	profane := []struct {
		input string
		words []string
	}{
		{"dari tadi ngentot aja", []string{"ngentot"}},
		{"dasar bacot", []string{"bacot"}},
	}
	for _, tc := range profane {
		t.Run(tc.input, func(t *testing.T) {
			res := b.Text(tc.input).Collect(true).Scan()
			if !res.IsProfane() {
				t.Errorf("should detect profanity in %q", tc.input)
				return
			}
			for _, w := range tc.words {
				if !slices.Contains(res.Extract(), w) {
					t.Errorf("expected %q in %v", w, res.Extract())
				}
			}
		})
	}
}

// ArticleText — teks artikel/berita yang seharusnya bersih
func TestArticleText(t *testing.T) {
	b := bacot.New()

	cases := []struct {
		input string
		note  string
	}{
		{"Berdasarkan hasil penelitian yang dilakukan", "formasi ber- + 'dasar' bukan kata kasar"},
		{"Pemerintah Indonesia mengadakan rapat kabinet", "teks berita normal"},
		{"Para mahasiswa mengerjakan tugas akhir", "kalimat kampus"},
		{"Menteri Pendidikan memberikan sambutan", "berita formal"},
		{"Perusahaan akan mengumumkan hasil", "berita ekonomi"},
	}
	for _, tc := range cases {
		t.Run(tc.note, func(t *testing.T) {
			res := b.Text(tc.input).Scan()
			if res.IsProfane() {
				t.Errorf("clean article text should not be profane: %q → %v", tc.input, res.Extract())
			}
		})
	}
}

// ObfuscatedInput — berbagai bentuk pengelabuan (leet, spasi, stacked)
func TestObfuscatedInput(t *testing.T) {
	b := bacot.New()

	t.Run("leet 4njing", func(t *testing.T) {
		res := b.Text("4njing").WithLeetSpeak().Scan()
		if !res.IsProfane() {
			t.Error("'4njing' + leet should detect 'anjing'")
		}
	})

	t.Run("leet @njing", func(t *testing.T) {
		res := b.Text("@njing").WithLeetSpeak().Scan()
		if !res.IsProfane() {
			t.Error("'@njing' + leet should detect 'anjing'")
		}
	})

	t.Run("trim space a n j i n g", func(t *testing.T) {
		res := b.Text("a n j i n g").TrimSpace().Scan()
		if !res.IsProfane() {
			t.Error("'a n j i n g' + trim should detect 'anjing'")
		}
	})

	t.Run("unstack kontolll", func(t *testing.T) {
		res := b.Text("kontolll").Scan()
		if !res.IsProfane() {
			t.Error("'kontolll' + unstack should detect 'kontol'")
		}
	})

	t.Run("unstack anjiiing", func(t *testing.T) {
		res := b.Text("anjiiing").Scan()
		if !res.IsProfane() {
			t.Error("'anjiiing' + unstack should detect 'anjing'")
		}
	})

	t.Run("leet+trim 4njing spaced", func(t *testing.T) {
		res := b.Text("4 n j i n g").WithLeetSpeak().TrimSpace().Scan()
		if !res.IsProfane() {
			t.Error("'4 n j i n g' + leet + trim should detect 'anjing'")
		}
	})

	t.Run("number suffix babi123", func(t *testing.T) {
		res := b.Text("babi123").Scan()
		if !res.IsProfane() {
			t.Error("'babi123' should detect 'babi' (prefix match, rest '123' not one-syllable)")
		}
	})

	t.Run("exclamation not leet by default", func(t *testing.T) {
		res := b.Text("anjing!").Scan()
		if res.IsProfane() {
			t.Error("'anjing!' should NOT be detected — '!' not converted by default, prevChar=' ' blocks vowel match")
		}
	})

	t.Run("leet 8abi WithLeetSpeak", func(t *testing.T) {
		res := b.Text("8abi").WithLeetSpeak().Scan()
		if !res.IsProfane() {
			t.Error("'8abi' + leet should detect 'babi'")
		}
	})
}

// AffixedWords — kata berimbuhan yang terdeteksi atau ditolak
func TestAffixedWordsReal(t *testing.T) {
	b := bacot.New()

	detected := []struct {
		input string
		word  string
		note  string
	}{
		{"mebabi", "babi", "me- + babi, konsonan"},
		{"pebabi", "babi", "pe- + babi, konsonan"},
		{"terbabi", "babi", "ter- + babi, konsonan"},
		{"berbabi", "babi", "ber- + babi, konsonan"},
		{"ngasu", "asu", "ng- + asu, vokal + 'g'"},
		{"mengasu", "asu", "meng- + asu, vokal + 'g'"},
		{"pengasu", "asu", "peng- + asu, vokal + 'g'"},
		{"menyasu", "asu", "meny- + asu, vokal + 'y'"},
		{"penyasu", "asu", "peny- + asu, vokal + 'y'"},
		{"terasu", "asu", "ter- + asu, vokal + 'r'"},
		{"berasu", "asu", "ber- + asu, vokal + 'r'"},
		{"perasu", "asu", "per- + asu, vokal + 'r'"},
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
		{"diasu", "asu", "di- + vokal, tanpa g/r"},
		{"teasu", "asu", "te- + vokal"},
		{"beasu", "asu", "be- + vokal"},
		{"peasu", "asu", "pe- + vokal"},
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

// FalsePositivePrevention — kata normal yang mirip dengan bentuk berimbuhan
func TestFalsePositivePrevention(t *testing.T) {
	b := bacot.New()

	cases := []string{
		"pergi",      // per- + gi, "gi" bukan kata kasar
		"terus",      // stop word
		"terbakar",   // ter- + bakar, "bakar" bukan kata kasar
		"berdiri",    // ber- + diri, "diri" bukan kata kasar
		"baku",       // bukan kata kasar
		"baja",       // bukan kata kasar
		"susu",       // bukan kata kasar
		"kelas",      // bukan kata kasar
		"perasuai",   // per- + asu + "ai", rest "ai" adalah diftong 1 suku kata → ditolak
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			res := b.Text(input).Scan()
			if res.IsProfane() {
				t.Errorf("'%s' should NOT be detected (false positive)", input)
			}
		})
	}
}

// MixedScenarios — kalimat campuran dalam 1 percakapan
func TestMixedScenarios(t *testing.T) {
	b := bacot.New()

	t.Run("clean then profane", func(t *testing.T) {
		clean := b.Text("halo apa kabar").Scan()
		if clean.IsProfane() {
			t.Error("clean opening should not be profane")
		}

		profane := b.Text("dasar babi").Scan()
		if !profane.IsProfane() {
			t.Error("follow-up should detect profanity")
		}
	})

	t.Run("multiple bacot instance", func(t *testing.T) {
		b1 := bacot.New()
		b2 := bacot.New()

		r1 := b1.Text("halo").Scan()
		r2 := b2.Text("anjing").Scan()

		if r1.IsProfane() {
			t.Error("b1 clean should not be profane")
		}
		if !r2.IsProfane() {
			t.Error("b2 should detect 'anjing'")
		}
	})

	t.Run("affix false positive", func(t *testing.T) {
		res := b.Text("berdasarkan hasil").Scan()
		if res.IsProfane() {
			t.Error("'berdasarkan' should not be detected ('dasar' not in dict)")
		}
	})

	t.Run("prefix non-vowel always detected", func(t *testing.T) {
		res := b.Text("mebabi").Scan()
		if !res.IsProfane() {
			t.Error("'mebabi' should be detected ('babi' starts with consonant)")
		}
	})
}
