package bacot

import (
	"slices"
	"strings"
)

var (
	vocals = []rune{'a', 'e', 'i', 'o', 'u'}

	mengPrefix = []rune{'a', 'e', 'i', 'o', 'u', 'k', 'g', 'h'}
	memPrefix  = []rune{'b', 'f', 'p', 'v'}
	menPrefix  = []rune{'c', 'd', 'j', 's', 't'}

	prefixes2 = []string{"me", "ng", "ny", "pe", "di", "te", "be"}
	prefixes3 = []string{"ber", "ter", "per", "mem", "pem", "men", "pen"}
	prefixes4 = []string{"meng", "peng", "peny", "meny"}

	suffixes = []string{
		"kan",
		"an",
		"i", "in",
		"nya",
		"ku",
		"mu",
	}
)

func isOneSyllable(s string) bool {
	if len(s) > 5 || len(s) < 1 {
		return false
	}

	var (
		preVowel   = 0
		vowelFound = 0
		postVowel = 0
	)
	for i := 0; i < len(s); i++ {
		c := rune(s[i])

		if slices.Contains(vocals, c) {

			vowelFound += 1
			if vowelFound > 1 {
				return false
			}

			// kalau ketemu vokal dan itu diftong, maka
			// akan dihitung 1
			// vocal diftong {ai, au, ei, oi}
			if i < len(s)-2 || (len(s) == 2 && i == 0) {
				diphthong := s[i+1]
				if (c == 'a' && (diphthong == 'i' || diphthong == 'u')) ||
					(c == 'e' && diphthong == 'i') ||
					(c == 'o' && diphthong == 'i') {
					// geser cursornya 1 langkah
					i += 1
				}
			}
			continue
		}

		// untuk kasus dimana suku kata terdairi dari -ng atau -ny
		// karena mereka dihitung 1 bentuk konsonan
		if c == 'n' {
			if i < len(s)-1 && (s[i+1] == 'g' || s[i+1] == 'y') {
				i += 1
			}
		}

		if vowelFound > 0 {
			postVowel += 1
		} else {
			preVowel += 1
		}

		if preVowel > 3 || postVowel > 3 {
			return false
		}
	}

	if vowelFound == 0 {
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
func dynamicVowelAlteration(s string) []string {

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

func consonantForEmphasis(s string) []string {

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

func avoidVowelToAbbrev(s string) []string {

	var ret []string
	return ret
}

func craftMan(s string) []string {

	var ret []string

	// kata asli
	ret = append(ret, s)

	// imbuhan yang melebur, prefix-
	ret = append(ret, nasalFusionWord(s)...)

	return ret
}
