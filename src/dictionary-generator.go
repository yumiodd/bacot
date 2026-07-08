package bacot

import (
	"slices"
	"strings"
)

// refrence: https://id.wikipedia.org/wiki/Prefiks

var (
	vocals = []rune{'a', 'e', 'i', 'o', 'u'}

	mengPrefix = []rune{'a', 'e', 'i', 'o', 'u', 'k', 'g', 'h'}
	// mePrefix   = []rune{'l', 'm', 'n', 'r', 'w', 'y'}
	memPrefix = []rune{'b', 'f', 'p', 'v'}
	menPrefix = []rune{'c', 'd', 'j', 's', 't'}

	prefixes = []string{"me", "mg", "ng", "pg", "mm", "pm", "mn", "my", "ny", "pe", "di", "te", "be"}
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
	for _, c := range s {
		if slices.Contains(vocals, c) {
			vocalFound += 1
			if vocalFound > 1 {
				return false
			}
			continue
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

func affix(s string) []string {

	// prefix
	// meng-, di-, ber-, ter-, peng-, per- se-
	// note:
	// has alomorf:
	// 	meng-,	{me-, mem-, men-, meny-, menge-}
	// 	ber-, 	{be-, bel-}
	// 	ter-, 	{te-}
	// 	per-, 	{pe-, pel-}
	// 	peng-, 	{pe-, pem-, pen-, peny-, penge-}

	// 1. meng-, peng-
	// meng-{k,t,s,p}+v...= meng-v...

	return []string{}
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

// below are the affic generator features, for now please comment first =======

// func meng_Prefix(s string) []string {
// 	s = strings.ToLower(s)
// 	sr := []rune(s)

// 	var ret []string

// 	if isOneSyllable(s) {
// 		ret = append(ret, "menge"+s)
// 		// meng-
// 	} else if slices.Contains(mengPrefix, sr[0]) {
// 		if sr[0] == 'k' && slices.Contains(vocals, sr[1]) {
// 			ret = append(ret, "meng"+s[1:])
// 		} else {
// 			ret = append(ret, "meng"+s)
// 		}

// 		// me-
// 	} else if slices.Contains(mePrefix, sr[0]) {
// 		if (sr[0] == 'n') && (sr[1] == 'g' || sr[1] == 'y') {
// 			ret = append(ret, "me"+s)
// 		} else {
// 			ret = append(ret, "men"+s)
// 		}

// 		// mem-
// 	} else if slices.Contains(memPrefix, sr[0]) {
// 		if sr[0] == 'p' && slices.Contains(vocals, sr[1]) {
// 			ret = append(ret, "mem"+s[1:])
// 		} else {
// 			ret = append(ret, "mem"+s)
// 		}

// 		// men-
// 	} else if slices.Contains(menPrefix, sr[0]) {
// 		if (sr[0] == 't') && (slices.Contains(vocals, sr[1])) {
// 			ret = append(ret, "men"+s[1:])
// 		} else if (sr[0] == 's') && (slices.Contains(vocals, sr[1])) {
// 			ret = append(ret, "meny"+s[1:])
// 		} else {

// 			ret = append(ret, "men"+s)
// 		}
// 	}

// 	return ret
// }
// func peng_Prefix(s string) []string {
// 	s = strings.ToLower(s)
// 	sr := []rune(s)

// 	var ret []string

// 	if isOneSyllable(s) {
// 		ret = append(ret, "penge"+s)
// 		// meng-
// 	} else if slices.Contains(mengPrefix, sr[0]) {
// 		if sr[0] == 'k' && slices.Contains(vocals, sr[1]) {
// 			ret = append(ret, "peng"+s[1:])
// 		} else {
// 			ret = append(ret, "peng"+s)
// 		}

// 		// me-
// 	} else if slices.Contains(mePrefix, sr[0]) {
// 		if (sr[0] == 'n') && (sr[1] == 'g' || sr[1] == 'y') {
// 			ret = append(ret, "pe"+s)
// 		} else {
// 			ret = append(ret, "pen"+s)
// 		}

// 		// mem-
// 	} else if slices.Contains(memPrefix, sr[0]) {
// 		if sr[0] == 'p' {
// 			ret = append(ret, "pem"+s[1:])
// 		} else {
// 			ret = append(ret, "pem"+s)
// 		}

// 		// men-
// 	} else if slices.Contains(menPrefix, sr[0]) {
// 		if (sr[0] == 't') && (slices.Contains(vocals, sr[1])) {
// 			ret = append(ret, "pen"+s[1:])
// 		} else if (sr[0] == 's') && (slices.Contains(vocals, sr[1])) {
// 			ret = append(ret, "peny"+s[1:])
// 		} else {

// 			ret = append(ret, "pen"+s)
// 		}
// 	}

// 	return ret
// }

// func ber_Prefix(s string) []string {

// 	s = strings.ToLower(s)
// 	sr := []rune(s)

// 	var ret []string
// 	if (sr[0] == 'r') || (strings.Contains(getFirstSyllable(s), "er")) {
// 		ret = append(ret, "be"+s, "kebe"+s, "sebe"+s)

// 	} else {
// 		ret = append(ret, "ber"+s, "keber"+s, "seber"+s)
// 	}

// 	return ret
// }
// func ter_Prefix(s string) []string {

// 	s = strings.ToLower(s)
// 	sr := []rune(s)

// 	var ret []string
// 	if (sr[0] == 'r') || (strings.Contains(getFirstSyllable(s), "er")) {
// 		ret = append(ret, "te"+s, "kete"+s)

// 	} else {
// 		ret = append(ret, "ter"+s, "keter"+s)
// 	}

// 	return ret
// }

// func per_Prefix(s string) []string {

// 	s = strings.ToLower(s)
// 	sr := []rune(s)

// 	var ret []string
// 	if (sr[0] == 'r') || (strings.Contains(getFirstSyllable(s), "er")) {
// 		ret = append(ret, "pe"+s, "mempe"+s, "dipe"+s, "terpe"+s)

// 	} else {
// 		ret = append(ret, "per"+s, "memper"+s, "diper"+s, "terper"+s, "seper"+s)
// 	}

// 	ret = append(ret, "pe"+s)
// 	return ret
// }

// func ke_Prefix(s string) []string {
// 	s = strings.ToLower(s)
// 	return []string{"ke" + s}
// }

// func se_Prefix(s string) []string {
// 	s = strings.ToLower(s)
// 	return []string{"se" + s}
// }

// func di_Prefix(s string) []string {
// 	s = strings.ToLower(s)
// 	return []string{"di" + s}
// }

// func an_Suffix(s string) []string { // peng-an, per-an, ke-an, se-an
// 	s = strings.ToLower(s)

// 	if s[:2] == "an" {
// 		return []string{}
// 	}

// 	return []string{s + "an"}
// }

// func kan_Suffix(s string) []string { // me-kan, ter-kan
// 	s = strings.ToLower(s)
// 	return []string{s + "kan"}
// }

// func i_Suffix(s string) []string { // me-i, di-i, ter-i
// 	s = strings.ToLower(s)
// 	return []string{s + "i"}
// }

// func in_Suffix(s string) []string { // me-i, di-i, ter-i // informal
// 	s = strings.ToLower(s)
// 	return []string{s + "in"}
// }

// func ku_Suffix(s string) []string {
// 	s = strings.ToLower(s)
// 	return []string{s + "ku"}
// }
