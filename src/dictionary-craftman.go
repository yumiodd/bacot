package bacot

import (
	"slices"
	"strings"
)

// Prefix Indonesia berdasarkan kaidah morfologi:
//   - meng- → nasal fusion: meng- + kata mulai vokal/k/g/h
//   - mem- → nasal fusion: mem- + kata mulai b/f/p/v
//   - men- → nasal fusion: men- + kata mulai c/d/j/s/t
//
// Prefix di Scan() digunakan untuk stripping imbuhan
// dengan prioritas: 4 char → 3 char → 2 char.
var (
	vocals = []rune{'a', 'e', 'i', 'o', 'u'}

	mengPrefix = []rune{'a', 'e', 'i', 'o', 'u', 'k', 'g', 'h'}
	memPrefix  = []rune{'b', 'f', 'p', 'v'}
	menPrefix  = []rune{'c', 'd', 'j', 's', 't'}

	prefixes2 = []string{"me", "ng", "ny", "pe", "di", "te", "be"}
	prefixes3 = []string{"ber", "ter", "per", "mem", "pem", "men", "pen"}
	prefixes4 = []string{"meng", "peng", "peny", "meny"}

	// Suffix Indonesia yang umum untuk validasi affix:
	// sisa setelah stem bukan suffix → cek isOneSyllable.
	suffixes = []string{
		"kan",
		"an",
		"i", "in",
		"nya",
		"ku",
		"mu",
	}
)
// isOneSyllable mengecek apakah string adalah 1 suku kata
// dalam kaidah Bahasa Indonesia.
//
// Digunakan oleh affix detection untuk mencegah false positive:
// jika sisa setelah stripping prefix/stem adalah tepat 1 suku
// kata, kemungkinan besar kata sudah berubah makna.
// Contoh: "babiru" → stem "babi" + sisa "ru" (1 suku kata) → skip
//
// Aturan suku kata:
//   - Maksimal 1 vokal (diftong dihitung 1 vokal)
//   - Maksimal 3 konsonan sebelum vokal
//   - Maksimal 3 konsonan setelah vokal
//   - 'ng' dan 'ny' dihitung 1 konsonan
func isOneSyllable(s string) bool {
	if len(s) > 5 || len(s) < 1 {
		return false
	}

	var (
		preVocal   = 0
		vocalFound = 0
		pascaVoal  = 0
	)
	for i := 0; i < len(s); i++ {
		c := rune(s[i])

		if slices.Contains(vocals, c) {

			vocalFound += 1
			if vocalFound > 1 {
				return false
			}

			// Diftong {ai, au, ei, oi} dihitung 1 vokal
			if i < len(s)-2 || (len(s) == 2 && i == 0) {
				diftongPair := s[i+1]
				if (c == 'a' && (diftongPair == 'i' || diftongPair == 'u')) ||
					(c == 'e' && diftongPair == 'i') ||
					(c == 'o' && diftongPair == 'i') {
					i += 1
				}
			}
			continue
		}

		// 'ng' dan 'ny' dihitung 1 konsonan
		if c == 'n' {
			if i < len(s)-1 && (s[i+1] == 'g' || s[i+1] == 'y') {
				i += 1
			}
		}

		if vocalFound > 0 {
			pascaVoal += 1
		} else {
			preVocal += 1
		}

		if preVocal > 3 || pascaVoal > 3 {
			return false
		}

	}

	if vocalFound == 0 {
		return false
	}

	return true
}

func getFirstSyllable(s string) string {

	var i = 5
	if len(s) <= i {
		i = len(s)
	}

	for ; i >= 3; i-- {
		if isOneSyllable(s[:i]) {
			return s[:i]
		}
	}

	return ""
}

func nasalFusionWord(s string) []string {
	s = strings.ToLower(s)
	sr := []rune(s)

	var ret []string

	if (slices.Contains(mengPrefix, sr[0])) && sr[0] == 'k' && slices.Contains(vocals, sr[1]) {

		// ended -g
		for _, pfx := range []string{"meng", "peng", "ng"} {
			ret = append(ret, pfx+s[1:])
		}

	} else if slices.Contains(memPrefix, sr[0]) && sr[0] == 'p' {
		if slices.Contains(vocals, sr[1]) {
			for _, pfx := range []string{"mem", "m"} {
				ret = append(ret, pfx+s[1:])
			}
		}

		for _, pfx := range []string{"pem"} {
			ret = append(ret, pfx+s[1:])
		}

	} else if slices.Contains(menPrefix, sr[0]) {
		if (sr[0] == 't') && (slices.Contains(vocals, sr[1])) {

			for _, pfx := range []string{"men", "pen"} {
				ret = append(ret, pfx+s[1:])
			}

		} else if (sr[0] == 's') && (slices.Contains(vocals, sr[1])) {

			for _, pfx := range []string{"meny", "peny", "ny"} {
				ret = append(ret, pfx+s[1:])
			}
		}
	}

	return ret
}
func dyamicVocalAlteration(s string) []string {

	s = strings.ToLower(s)
	sr := []rune(s)
	var ret []string

	var idx = map[rune][]int{
		'u': {},
		'e': {},
	}

	for i, r := range sr {
		if r == 'o' {
			idx['u'] = append(idx['u'], i)
		}
		if r == 'i' {
			idx['e'] = append(idx['e'], i)
		}
	}

	var (
		a    = slices.Clone(sr)
		flag = 0
	)
	for k, v := range idx {
		if len(v) == 0 {
			continue
		}

		alt := slices.Clone(sr)
		for _, i := range v {
			alt[i] = k
			a[i] = k
		}

		flag += 1
		ret = append(ret, string(alt))
	}
	if flag == 2 {
		ret = append(ret, string(a))
	}

	return ret
}

func consonanForEmphasis(s string) []string {

	// contoh: tai -> taik, asu - > asuk, puki -> pukik
	// note: mempertimbangan menambah penekanan 'g'
	s = strings.ToLower(s)

	if len(s) >= 3 && slices.Contains(vocals, rune(s[len(s)-1])) {
		return []string{s + "k"}
	}
	return []string{}
}

// func suffix(s string) []string {

// 	s = strings.ToLower(s)
// 	lastChar := len(s) - 1
// 	var ret []string

// 	ret = append(ret,
// 		s+"kan", s+"kn",
// 		s+"an",
// 		s+"in",
// 		s+"ku",
// 		s+"mu",
// 		s+"nya", s+"ny",
// 		s+"lah", s+"lh",
// 		s+"kah", s+"kh",
// 		s+"tah", s+"th",
// 		s+"pun", s+"pn",
// 		// tolong tambah jika masih ada
// 	)

// 	// -i
// 	if s[lastChar] != 'i' {
// 		ret = append(ret, s+"i")
// 	}
// 	if s[lastChar] == 'l' {
// 		ret = append(ret, s+"ah")
// 	}

// 	return ret
// }

func avoidVocalToAbv(s string) []string {

	var ret []string
	return ret
}

// craftMan menghasilkan varian imbuhan dari kata dasar.
// Ini adalah optimization: generate semua varian di init,
// bukan stemming di runtime. Tiap Scan() cukup exact match
// di dictionary tanpa perlu proses imbuhan ulang.
//
// Contoh: "babi" → "babi", "mebabi", "pebabi", "tebabi",
// "dibabi", "kebabi", "ngbabi", "bebabi" (dll via nasal fusion)
func craftMan(s string) []string {

	var ret []string

	ret = append(ret, s)

	// Nasal fusion: prefix yang melebur dengan kata dasar
	// Contoh: "pukul" → "memukul", "pemukul", "mukul"
	ret = append(ret, nasalFusionWord(s)...)

	return ret
}
