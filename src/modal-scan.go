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

// Scan() menggunakan strategi token-based (split by spasi),
// bukan substring search. Ini mengurangi false positive:
// "kelas" tidak akan terdeteksi sebagai "elas" karena token exact match.
//
// Optimasi: length histogram pre-filter — jika tidak ada kata
// di dictionary dengan panjang yang sama, token langsung skip.
// Ini menghindari map lookup yang tidak perlu.
func (ms *ModalScan) Scan() *ScanResult {

	res := &ScanResult{
		text:        ms.input,
		praScanText: ms.text,
	}

	var (
		words = strings.SplitSeq(res.praScanText, " ")
		idx   = 0
	)

	for w := range words {
		var lenW = len(w)

		// Length histogram: skip token kalau di dictionary
		// tidak ada satupun kata dengan panjang yang sama.
		if !(ms.dict.IsContainLen(lenW)) {
			idx += len(w) + 1
			continue
		}
		// Stop words dan false positive langsung dilewati.
		if ms.dict.IsStopWord(w) || ms.dict.IsFalsePositive(w) {
			idx += len(w) + 1
			continue
		}

		// Exact match di dictionary (performa O(1) karena map).
		var found = ms.dict.Contains(w)

		// Affix detection dengan sliding window.
		// Jika exact match gagal dan affix aktif, kita coba
		// stripping prefix Indonesia (me-, pe-, ber-, ter-, di-, dll)
		// lalu cek apakah stem-nya ada di dictionary.
		// Contoh: "mebabi" → strip "me" → "babi" → ditemukan ✅
		if !found && ms.affix && (lenW > 3) {

			var (
				wTemp    string
				prevChar rune
			)
			// Coba stripping prefix Indonesia dengan prioritas:
			// 4 char (meng-, peng-) → 3 char (ber-, ter-, per-) → 2 char (me-, di-, pe-)
			// prevChar digunakan untuk validasi nasal fusion dan peluluhan.
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

			// Sliding window: iterate semua panjang kata yang ada
			// di dictionary untuk mencari stem yang cocok.
			for _, r := range ms.dict.GetWordsLen() {

				if found {
					break
				}
				if r > len(wTemp) {
					break
				}

				word := wTemp[:r]
				if ms.dict.Contains(word) {

					// Validasi nasal fusion:
					// Jika stem diawali vokal, prefix harus berakhiran
					// 'g' (meng-, peng-) atau 'r' (ber-, per-, ter-).
					// Contoh: "mengikat" → prevChar 'g' ✅
					//          "mengeja" → prevChar 'n' ❌
					if slices.Contains(vocals, rune(word[0])) && !(prevChar == 'g' || prevChar == 'r' || prevChar == 'n') {
						continue
					}

					// Koreksi over-stripping:
					// "memakan" → strip "mem-" → "akan"
					// Tapi kata asli adalah "makan". Cek juga w[2:] = "emakan".
					if prevChar == 'm' && word[0] == 'm' && ms.dict.Contains(w[2:]) {
						found = true
						break
					}

					// Cegah false positive: jika sisa setelah stem
					// adalah 1 suku kata, kemungkinan kata sudah
					// berubah makna. Contoh: "babiru" → "babi" + "ru"
					// "ru" adalah 1 suku kata → bukan "babi" ❌
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

// RecursiveScan() mencari kata kotor di dalam substring token.
// Berguna untuk kasus seperti "xbabi" di mana kata kotor
// menempel dengan karakter lain. RecursiveScan akan sliding
// window dari setiap posisi l di dalam token dan mencocokkan
// dengan dictionary.
//
// Kompleksitas: O(token_length × dict_word_lengths).
// Untuk token normal (<20 char) sangat cepat.
// Untuk token sangat panjang (10K char) perlu diwaspadai.
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

		// Sliding window: iterasi dari setiap posisi l di token.
		// Setiap posisi diambil substring s[l:] dan dicocokkan
		// dengan dictionary menggunakan semua panjang kata yang ada.
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

func (ms *ModalScan) SanitizeReadSign() *ModalScan {

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
