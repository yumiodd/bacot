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

	suffixes = []string{"kan", "an", "i", "if", "al", "is", "ni", "ik", "wan", "wati", "man", "is", "or", "er", "tas", "isme", "tas", "nya", "ku", "mu"}
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

func ky_consonanAddition(s string) []string {
	s = strings.ToLower(s)

	if len(s) >= 3 && slices.Contains(vocals, rune(s[len(s)-1])) {
		return []string{s + "k", s, s + "y"}
	}
	return []string{}
}

func craftMan(s string) []string {

	var ret []string
	ret = append(ret, nasalFusionWord(s)...)
	ret = append(ret, dyamicVocalAlteration(s)...)
	ret = append(ret, ky_consonanAddition(s)...)

	return ret

}
