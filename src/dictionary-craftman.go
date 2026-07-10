package bacot

import (
	"slices"
	"strings"
)

var (
	vocals = []rune{'a', 'e', 'i', 'o', 'u'}

	mengPrefix = []rune{'a', 'e', 'i', 'o', 'u', 'k', 'g', 'h'}
	// mePrefix   = []rune{'l', 'm', 'n', 'r', 'w', 'y'}
	memPrefix = []rune{'b', 'f', 'p', 'v'}
	menPrefix = []rune{'c', 'd', 'j', 's', 't'}

	prefixes2 = []string{"me", "mg", "ng", "pg", "mm", "pm", "mn", "ny", "pe", "di", "te", "be", "tr", "br"}
	prefixes3 = []string{"ber", "ter", "per", "mem", "pem"}

	suffixes = []string{
		"kan",
		"an",
		"i",
		"if",
		"al",
		"is",
		"ni",
		"ik",
		"wan",
		"wati",
		"man",
		"or",
		"er",
		"tas",
		"isme",
		"tas",
		"nya",
		"ku",
		"mu",
		"kahn",
		"khan",
	}
)

func isOneSyllable(s string) bool {
	if len(s) > 5 || len(s) < 1 {
		return false
	}

	var (
		preVocal   = 0
		vocalFound = 0
		pascaVoal  = 0
	)
	for i := 0; i <= len(s); i++ {
		c := rune(s[i])

		if slices.Contains(vocals, c) {

			vocalFound += 1
			if vocalFound > 1 {
				return false
			}

			// kalau ketemu vokal dan itu diftong, maka
			// akan dihitung 1
			// vocal diftong {ai, au, ei, oi}
			if i < len(s) {
				diftongPair := s[i+1]
				if (c == 'a' && (diftongPair == 'i' || diftongPair == 'u')) ||
					(c == 'e' && diftongPair == 'i') ||
					(c == 'o' && diftongPair == 'i') {
					// geser cursornya 1 langkah
					i += 1
				}
			}
			continue
		}

		// untuk kasus dimana suku kata terdairi dari -ng atau -ny
		// karena mereka dihitung 1 bentuk konsonan
		if c == 'n' {
			if i < len(s) && (s[i+1] == 'g' || s[i+1] == 'y') {
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

		for _, pfx := range []string{"meng", "mng", "mg", "ng", "peng", "png", "pg"} {
			ret = append(ret, pfx+s[1:])
		}

	} else if slices.Contains(memPrefix, sr[0]) && sr[0] == 'p' {
		if slices.Contains(vocals, sr[1]) {
			for _, pfx := range []string{"mm", "mem"} {
				ret = append(ret, pfx+s[1:])
			}
		}

		for _, pfx := range []string{"pem", "pm"} {
			ret = append(ret, pfx+s[1:])
		}

	} else if slices.Contains(menPrefix, sr[0]) {
		if (sr[0] == 't') && (slices.Contains(vocals, sr[1])) {

			for _, pfx := range []string{"men", "mn", "pen", "pen"} {
				ret = append(ret, pfx+s[1:])
			}

		} else if (sr[0] == 's') && (slices.Contains(vocals, sr[1])) {

			for _, pfx := range []string{"meny", "mny", "my", "ny", "peny", "pny"} {
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
		return []string{s + "k", s}
	}
	return []string{}
}

func suffix(s string) []string {

	s = strings.ToLower(s)
	lastChar := len(s) - 1
	var ret []string

	ret = append(ret,
		s+"kan", s+"kn",
		s+"an",
		s+"in",
		s+"ku",
		s+"mu",
		s+"nya", s+"ny",
		s+"lah", s+"lh",
		s+"kah", s+"kh",
		s+"tah", s+"th",
		s+"pun", s+"pn",
		// tolong tambah jika masih ada
	)

	// -i
	if s[lastChar] != 'i' {
		ret = append(ret, s+"i")
	}
	if s[lastChar] == 'l' {
		ret = append(ret, s+"ah")
	}

	return ret
}

func avoidVocalToAbv(s string) []string {

	var ret []string
	return ret
}

func craftMan(s string) []string {

	var ret []string
	// kata asli
	ret = append(ret, s)

	// vocal o dan i yang berubah
	alteredVocal := dyamicVocalAlteration(s)
	ret = append(ret, alteredVocal...)

	// imbuhan yang melebur, prefix-
	pref := nasalFusionWord(s)
	ret = append(ret, pref...)

	// suffix
	suff := suffix(s)
	ret = append(ret, suff...)

	// penekanan akhir kata dengan vokal
	emphasis := consonanForEmphasis(s)
	ret = append(ret, emphasis...)

	// prefix + suffix
	for _, p := range pref {
		ret = append(ret, suffix(p)...)
	}

	// preffix + alteredVocal
	for _, a := range alteredVocal {
		alteredVocalWithPreffix := nasalFusionWord(a)
		ret = append(ret, alteredVocalWithPreffix...)

		// preffix + altered vocal + suffix
		for _, ap := range alteredVocalWithPreffix {
			ret = append(ret, suffix(ap)...)
		}
	}

	return ret
}
