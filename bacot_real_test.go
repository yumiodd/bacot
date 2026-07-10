package bacot

import (
	"fmt"
	"strings"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

// TestRealCleanInput — masukan pengguna nyata yang seharusnya BERSIH
func TestRealCleanInput(t *testing.T) {
	b := bacot.New()

	inputs := []struct {
		input string
		note  string
	}{
		// Percakapan sehari-hari
		{"selamat pagi semua", "sapaan"},
		{"lagi apa sekarang", "chat santai"},
		{"udah makan siang", "chat casual"},
		{"di mana rumahmu", "pertanyaan, stop word di-"},
		{"kapan kita ketemu", "pertanyaan"},
		{"terima kasih banyak", "ucapan terima kasih"},
		{"saya suka makan nasi goreng", "opini makanan"},

		// Artikel & berita
		{"pemerintah mengadakan rapat kabinet", "berita politik"},
		{"berdasarkan hasil penelitian yang dilakukan", "artikel ilmiah"},
		{"para mahasiswa mengerjakan tugas akhir", "berita kampus"},
		{"menteri pendidikan memberikan sambutan", "berita formal"},
		{"perusahaan akan mengumumkan hasil", "berita ekonomi"},
		{"buku ini sangat bermanfaat", "review buku"},

		// Kata umum yang mirip dengan bentuk berimbuhan
		{"pergi ke pasar", "per- + gi, gi bukan kata kasar"},
		{"terus belajar", "stop word"},
		{"terbakar habis", "ter- + bakar, bakar bukan kata kasar"},
		{"berdiri tegak", "ber- + diri, diri bukan kata kasar"},
		{"tolong ambilkan buku", "permintaan"},
		{"sedang apa sekarang", "pertanyaan"},

		// Kata normal yang tidak kasar
		{"buku", "bukan kata kasar"},
		{"kelas", "bukan kata kasar"},
		{"susu", "bukan kata kasar"},
		{"baja", "bukan kata kasar"},
		{"baku", "bukan kata kasar"},
		{"terasi", "bukan kata kasar (ter- + asi, asi bukan kata kasar)"},

		// Diftong
		{"sungai", "akhiran diftong, bukan kata kasar"},
		{"pantai", "akhiran diftong, bukan kata kasar"},
		{"sampai", "akhiran diftong, bukan kata kasar"},
		{"pakai", "akhiran diftong, bukan kata kasar"},
	}

	var failed bool
	for _, tc := range inputs {
		t.Run(tc.note, func(t *testing.T) {
			res := b.Text(tc.input).Scan()
			if res.IsProfane() {
				failed = true
				t.Errorf("CLEAN EXPECTED but detected as profane: %q (%s) — words: %v", tc.input, tc.note, res.Extract())
			}
		})
	}

	if failed {
		t.Log("╔══════════════════════════════════════════════╗")
		t.Log("║  BERSIH — ada false positive!               ║")
		t.Log("╚══════════════════════════════════════════════╝")
	}
}

// TestRealProfanityInput — masukan pengguna nyata yang seharusnya TERDETEKSI
func TestRealProfanityInput(t *testing.T) {
	b := bacot.New()

	inputs := []struct {
		input   string
		note    string
		options func(*bacot.ModalScan) *bacot.ModalScan
	}{
		// Direct profanity
		{"dasar babi", "komentar kasar", nil},
		{"anjing lu", "hinaan", nil},
		{"goblok amat", "makian", nil},
		{"dasar memek", "hinaan", nil},
		{"dasar asu", "hinaan", nil},
		{"kontol", "kata kasar", nil},
		{"bangsat", "kata kasar", nil},
		{"bacot", "kata kasar", nil},
		{"ngentot", "kata kasar", nil},
		{"sial", "kata kasar", nil},
		{"bodoh", "kata kasar", nil},
		{"tolol", "kata kasar", nil},

		// Multi profanity
		{"kontol bangsat", "gabungan 2 kata kasar", nil},
		{"babi anjing", "gabungan 2 kata kasar", nil},
		{"dasar babi lu babi", "komentar kasar berulang", nil},

		// Affixed (real Indonesian usage)
		{"terbabi", "ter- + babi, dipakai di sosial media", nil},
		{"membabi buta", "mem- + babi, frasa 'membabi buta'", nil},
		{"ngasu", "ng- + asu", nil},
		{"mengasu", "meng- + asu", nil},
		{"pengasu", "peng- + asu", nil},

		// Obfuscation — leet speak
		{"4njing", "leet '4' → 'a'", func(ms *bacot.ModalScan) *bacot.ModalScan { return ms.WithLeetSpeak() }},
		{"8abi", "leet '8' → 'b'", func(ms *bacot.ModalScan) *bacot.ModalScan { return ms.WithLeetSpeak() }},
		{"4su", "leet '4' → 'a'", func(ms *bacot.ModalScan) *bacot.ModalScan { return ms.WithLeetSpeak() }},

		// Obfuscation — trim space
		{"b a b i", "spasi di antara huruf", func(ms *bacot.ModalScan) *bacot.ModalScan { return ms.TrimSpace() }},
		{"a s u", "spasi di antara huruf", func(ms *bacot.ModalScan) *bacot.ModalScan { return ms.TrimSpace() }},

		// Obfuscation — stacked chars
		{"kontolll", "huruf 'l' di-stack", nil},
		{"anjiiing", "huruf 'i' di-stack", nil},
		{"baaabi", "huruf 'a' di-stack", nil},

		// Obfuscation — number suffix (babi ditemukan via stem)
		{"babi123", "angka di belakang, 'babi' tetap terdeteksi", nil},

		// Leet + trim + unstack combined
		{"4 n j i n g", "leet + spasi", func(ms *bacot.ModalScan) *bacot.ModalScan {
			return ms.WithLeetSpeak().TrimSpace()
		}},
		{"4njiiing", "leet + stack", func(ms *bacot.ModalScan) *bacot.ModalScan {
			return ms.WithLeetSpeak()
		}},
	}

	var failed bool
	for _, tc := range inputs {
		name := strings.ReplaceAll(tc.note, " ", "_")
		t.Run(name, func(t *testing.T) {
			ms := b.Text(tc.input)
			if tc.options != nil {
				ms = tc.options(ms)
			}
			res := ms.Scan()
			if !res.IsProfane() {
				failed = true
				t.Errorf("PROFANE EXPECTED but not detected: %q (%s)", tc.input, tc.note)
			}
		})
	}

	if failed {
		t.Log("╔══════════════════════════════════════════════╗")
		t.Log("║  PROFAN — ada yang tidak terdeteksi!        ║")
		t.Log("╚══════════════════════════════════════════════╝")
	}
}

// TestMixedCleanThenProfane — simulasi percakapan utuh
func TestMixedCleanThenProfane(t *testing.T) {
	b := bacot.New()

	scenarios := []struct {
		name    string
		cleans  []string
		profane string
	}{
		{
			"chat santai lalu makian",
			[]string{"halo apa kabar", "lagi apa", "oke siap"},
			"dasar babi",
		},
		{
			"obrolan lucu lalu hinaan",
			[]string{"hahaha", "iya betul", "lucu banget"},
			"anjing lu",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			var cleanFailed bool
			for _, c := range sc.cleans {
				res := b.Text(c).Scan()
				if res.IsProfane() {
					cleanFailed = true
					t.Errorf("CLEAN in conversation but detected: %q → %v", c, res.Extract())
				}
			}
			if cleanFailed {
				t.Log("⚠️  Clean texts in conversation had false positives")
			}

			res := b.Text(sc.profane).Scan()
			if !res.IsProfane() {
				t.Errorf("PROFANE in conversation but not detected: %q", sc.profane)
			}
		})
	}
}

// TestRealReport — laporan akhir ringkas
func TestRealReport(t *testing.T) {
	b := bacot.New()

	cleanSet := []string{
		"selamat pagi",
		"terima kasih",
		"lagi apa",
		"pergi ke pasar",
		"terus belajar",
		"buku",
		"kelas",
		"susu",
		"terasi",
		"sungai",
	}

	profaneSet := []string{
		"babi",
		"anjing",
		"asu",
		"kontol",
		"memek",
		"goblok",
		"bangsat",
		"bacot",
		"ngentot",
		"sial",
		"tolol",
		"bodoh",
	}

	var (
		falsePositives []string
		falseNegatives []string
	)

	for _, input := range cleanSet {
		if b.Text(input).Scan().IsProfane() {
			falsePositives = append(falsePositives, input)
		}
	}

	for _, input := range profaneSet {
		if !b.Text(input).Scan().IsProfane() {
			falseNegatives = append(falseNegatives, input)
		}
	}

	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║         BACOT — LAPORAN REAL TEST           ║")
	fmt.Println("╠══════════════════════════════════════════════╣")
	fmt.Printf("║  Clean inputs tested     : %2d                  ║\n", len(cleanSet))
	fmt.Printf("║  Profane inputs tested   : %2d                  ║\n", len(profaneSet))
	fmt.Printf("║  False positives         : %2d                  ║\n", len(falsePositives))
	fmt.Printf("║  False negatives         : %2d                  ║\n", len(falseNegatives))
	fmt.Println("╠══════════════════════════════════════════════╣")

	if len(falsePositives) > 0 {
		fmt.Println("║  FALSE POSITIVES:                           ║")
		for _, fp := range falsePositives {
			fmt.Printf("║    • %-36s  ║\n", fp)
		}
	}

	if len(falseNegatives) > 0 {
		fmt.Println("║  FALSE NEGATIVES:                           ║")
		for _, fn := range falseNegatives {
			fmt.Printf("║    • %-36s  ║\n", fn)
		}
	}

	fmt.Println("╚══════════════════════════════════════════════╝")

	if len(falsePositives) > 0 || len(falseNegatives) > 0 {
		t.Error("Real test report menunjukkan masalah")
	}
}
