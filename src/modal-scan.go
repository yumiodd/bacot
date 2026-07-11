package bacot

import (
	"slices"
	"strings"
)

type ModalScan struct {
	collect bool // def: false
	affix   bool // def: true

	// temp, upper layer set by (bacot)
	input string
	text  string
	dict  *Dictionary
}

func (ms *ModalScan) Collect(v bool) *ModalScan {
	ms.collect = v
	return ms
}

func (ms *ModalScan) Affix(v bool) *ModalScan {
	ms.affix = v
	return ms
}

// per word scanning
func (ms *ModalScan) Scan() *ScanResult {

	res := &ScanResult{
		text:        ms.input,
		praScanText: ms.text,
	}

	var (
		words = strings.SplitSeq(res.praScanText, " ")
		idx   = 0
	)

	// example input : "word in sentence" or "word-join"
	for w := range words {
		var lenW = len(w)

		// jika panjang w tidak termasuk dalam histogram frequent di dict, skip
		if !(ms.dict.IsContainLen(lenW)) {
			idx += len(w) + 1
			continue
		}
		// jika ternyata terdaftar dalam stopWord atau falsePositive
		if ms.dict.IsStopWord(w) || ms.dict.IsFalsePositive(w) {
			idx += len(w) + 1
			continue
		}

		// match perkata di kamus
		var found = ms.dict.Contains(w)

		// affic - prefix first
		if !found && ms.affix && (lenW > 3) {

			var (
				wTemp    string
				prevChar rune
			)
			if slices.Contains(prefixes4, w[:4]) {
				wTemp = w[4:]
				prevChar = rune(w[3])
			} else if slices.Contains(prefixes3, w[:3]) {
				wTemp = w[3:]
				prevChar = rune(w[2])
			} else if slices.Contains(prefixes2, w[:2]) {
				wTemp = w[2:]
				prevChar = rune(w[1])
			} else {
				wTemp = w[0:]
				prevChar = ' '
			}

			if len(wTemp) == 0 {
				break
			}

			for _, r := range ms.dict.GetWordsLen() {

				if found {
					break
				}
				if r > len(wTemp) {
					break
				}

				word := wTemp[:r]
				if ms.dict.Contains(word) {

					// kalau ketemu dan awal kata berawal huruf vokal,
					// maka sebelum huruf ini harus 'g' keranena pasti dari imbuhan {peng- , meng-}
					// atau peluluhan huruf 'r' dari imbuhan {per-, ber-, ter-}
					if slices.Contains(vocals, rune(word[0])) && !(prevChar == 'g' || prevChar == 'r' || prevChar == 'n') {
						continue
					}

					// membenarkan jika ada salah pengambilan
					// memakan -> akan, kerna memangkas [3:] mem-,
					// karena bisa jadi kata asli adalah makan, me- makan
					if prevChar == 'e' && word[0] == 'm' && ms.dict.Contains(w[1:]) {
						found = true
						break
					}

					// jika sisa potongan string adalah 1 suku kata sangat besar
					// kemungkinan kalau kata ini kata yang berbeda atau katanya sudah berubah makna
					// walaupun kata aslinya terdaftar penambahan 1 suku kata akan cukup untuk memperkeruh
					// intensi nya.
					rest := wTemp[r:]
					if rest != "" && !slices.Contains(suffixes, rest) && isOneSyllable(rest) {
						continue
					}

					found = true
					break
				}
			}
		}

		if found {
			res.words = append(res.words,
				&WordIndex{
					Word:  w,
					Start: idx,
					End:   len(w) + idx - 1,
				})

			if !(ms.collect) {
				break
			}
		}
		idx += len(w) + 1
	}
	return res
}

// Recursive will scan in sub-string fashion
func (ms *ModalScan) RecursiveScan() *ScanResult {

	res := &ScanResult{
		text:        ms.input,
		praScanText: ms.text,
	}

	var (
		finished = false
		idx      = 0
	)
	for s := range strings.SplitSeq(res.praScanText, " ") {
		if finished {
			break
		}

		if s == "" {
			idx++
			continue
		}

		for l := 0; l <= len(s); l++ {
			if finished {
				break
			}

			sub := s[l:]
			for _, r := range ms.dict.GetWordsLen() {

				if r > len(sub) {
					break
				}

				word := sub[:r]
				if ms.dict.Contains(word) {
					res.words = append(res.words, &WordIndex{
						Word:  word,
						Start: l + idx,
						End:   idx + l + r - 1,
					})

					if ms.collect {
						break
					}
					finished = true
				}
			}
		}
		idx += len(s) + 1
	}

	return res
}

func (ms *ModalScan) UnstackChar() *ModalScan {
	var (
		sb   strings.Builder
		prev rune
	)

	for _, c := range ms.text {
		if c == prev {
			continue
		} else {
			sb.WriteRune(c)
		}
		prev = c
	}

	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) ClearSpace() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if c == ' ' {
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()
	return ms
}

// Leet mengandung angka
// jika case nya := "babi123" -> "babiire" -> "babire"
// dengan affix "babire" tidak dianggap kata kotor karena "babi" + sukukata, beda makna
// maka jika lolos di leetspeak, scan kedua kali yang tanpa sanitasi leetspeak
//
// kesimpulan:
// jika input mengandung elemen leetspeak dan lolos scan
// scan ulang dengan mematikan leet
func (ms *ModalScan) WithLeetSpeak() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if r, ok := simpleLeetSpeaks[c]; ok {
			sb.WriteRune(r)
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) GetText() string {
	return ms.text
}

func (ms *ModalScan) SanitizeNewLine() *ModalScan {
	var sb strings.Builder
	for _, c := range ms.text {
		if c == '\n' {
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) TrimSpace() *ModalScan {

	var (
		sb   strings.Builder
		prev rune
	)

	for i, c := range ms.text {
		if i == 0 && c == ' ' {
			continue
		}
		if prev == ' ' && c == ' ' {
			continue
		}
		sb.WriteRune(c)
		prev = c
	}
	ms.text = sb.String()

	return ms
}

func (ms *ModalScan) ReplaceWhiteSpace() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if _, ok := whiteSpace[c]; ok && c != ' ' {
			sb.WriteRune(' ')
			continue
		}
		sb.WriteRune(c)
	}
	ms.text = sb.String()

	return ms
}

func (ms *ModalScan) SanitazeReadSign() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {
		if _, ok := replaceToSpace[c]; ok {
			sb.WriteRune(' ')
		} else if _, ok := replaceToClear[c]; ok {
			continue
		} else {
			sb.WriteRune(c)
		}
	}

	ms.text = sb.String()
	return ms
}

func (ms *ModalScan) SanitizeEmoji() *ModalScan {

	var sb strings.Builder
	for _, c := range ms.text {

		c := rune(c)

		if ((c >= '\U0001F600') && (c <= '\U0001F64F')) || // Emoticons
			((c >= '\U0001F300') && (c <= '\U0001F5FF')) || // Misc Symbols & Pictographs
			((c >= '\U0001F680') && (c <= '\U0001F6FF')) || // Transport & Map Symbols
			((c >= '\U0001F900') && (c <= '\U0001F9FF')) || // Supplemental Symbols
			((c >= '\U0001FA70') && (c <= '\U0001FAFF')) || // Symbols & Pictographs Ext-A
			((c >= '\u2600') && (c <= '\u26FF')) || // Misc Symbols (Blok Lama)
			((c >= '\u2700') && (c <= '\u27BF')) { // Dingbats (Blok Lama)
			sb.WriteRune(' ')
		} else {
			sb.WriteRune(c)
		}

	}

	ms.text = sb.String()
	return ms
}
