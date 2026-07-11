![bacot banner](./bacot.png)

# Bacot
### *"Mulut lo bacot — tapi bacot filter kata-kata lo."*

**Bacot** adalah profanity filter untuk Bahasa Indonesia yang handle affix (imbuhan), leet speak, dan unstack character. Zero dependency. Sekali `go get` langsung jalan.

Bahasa Indonesia aglutinatif — satu kata dasar bisa berubah dengan berbagai imbuhan (`babi` → `mebabi`, `pebabi`, `dibabi`). Bacot handle semua varian ini otomatis tanpa perlu daftarin satu-satu.

---

## 🔧 Installation

```go
import bacot "github.com/yumiodd/bacot/src"
```

Itu aja. Nggak ada dependency lain. *Kami bukan npm.*

---

## 🚀 Quick Start

```go
b := bacot.New()

res := b.Text("4njiiing").WithLeetSpeak().Scan()
res.IsProfane() // true
res.Censor()    // "*******"

res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
```

> **⚠️:** Jangan `New()` tiap scan. Buat satu instance, pakai berulang. Ini bukan sikat gigi.

---

## 📋 Fitur

| Fitur | Kerjaannya |
|-------|-----------|
| **Affix-aware** | `mebabi`? `penganjing`? Semua kena. Bacot ngerti Bahasa Indonesia pake imbuhan. |
| **Leet speak** | `4njing` → `anjing`, `8abi` → `babi`. *Kami bukan polisi bahasa, kami lebih dari itu.* |
| **Unstack char** | `anjiiiiiiiing` → `anjing`. Pake karakter berulang? *Nice try.* |
| **Dua mode scan** | `Scan()` per kata. `RecursiveScan()` substring, sampe ke akar. |
| **Custom dictionary** | `AddWord()`, `DelWord()` — lo punya kata kesayangan? Tambahin aja. |
| **Censor otomatis** | Mapping leet → `*`. User nulis `4njing`, sensor balikin `******`. |
| **Zero dependencies** | Cuma stdlib Go. **Zero.** Kalau lo nemu dependency lain, *itu halusinasi lo.* |
| **Collection mode** | Mau dikumpulin semua match? `Collect(true)`. Biar sekalian dihajar. |

---

## 📖 Cara Pake

### Basic Detection

```go
b := bacot.New()
res := b.Text("anjing").Scan()
res.IsProfane() // true
res.Censor()    // "******"
res.Extract()   // ["anjing"]
```

### Collect Mode — *"Kumpulin Semua Biar Puas"*

```go
res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
res.First()     // "babi"
res.Last()      // "anjing"
```

Default `Collect(false)` = berhenti di match pertama. **Kayak mantan: tau satu aja udah cukup.**

### Custom Words — *"Update Kamus, Update Perasaan"*

```go
b.AddWord("umiakucihanya")
b.AddWord("pukul") // otomatis detect mepukul, dipukul, pemukul
```

### Affix Mode

```go
b.Text("mebabi").Scan().IsProfane()      // true (default ON)
b.Text("mebabi").Affix(false).Scan()...  // false
```

Di Indonesia, "mebabi" tetep "babi". *Jangan munafik.*

### Recursive Scan — *"Nakal Kamu"*

```go
b.Text("xbabi").Scan().IsProfane()           // false
b.Text("xbabi").RecursiveScan().IsProfane()  // true
```

Mau lo taruh apapun di depan kata kasar, tetap ketahuan. *Udah, nyerah aja.*

### Leet Speak & Unstack

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane()   // true
b.Text("anjiiing").Scan().IsProfane()                  // true (unstack default ON)
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane()  // true (kombinasi)
```

### Trim Space — *"Pake Spasi? Trik Lo Udah Kedaluwarsa"*

```go
b.Text("a n j i n g").TrimSpace().Scan().IsProfane() // true
```

**Pokoknya kena.**

### Config — *"Sekali Atur, Scan Enteng"*

```go
b.Config(&bacot.ModalScanConfig{
    Affix:   true,
    Collect: true,
    Order: []bacot.SanitizeOrder{
        bacot.WithLeetSpeak,
        bacot.UnstackChar,
    },
})
b.Text("4njiiing").Scan().IsProfane() // otomatis
```

### Dictionary

```go
b.Dict.Contains("anjing")    // true
b.Dict.IsStopWord("dan")     // true
b.Dict.AddWords("tai", "keparat")
b.Dict.DelWords("anjing")
b.Dict.GetDict()             // map[string]struct{}
```

### Generator

```go
gen := res.Generator()
for w := gen.Yield(); w != nil; w = gen.Yield() {
    fmt.Println(w.Word)
}
```

---

## ⚙️ Cara Kerjanya

### Pipeline

```
Input → ToLower → [Sanitize] → Scan → ScanResult ─→ Censor()
                                                    ├→ IsProfane()
                                                    ├→ Extract()
                                                    ├→ Count()
                                                    └→ Generator()
```

**Precomputed WordIndex** — `Scan()` menghitung posisi start/end tiap kata kotor sekali. Setelah itu `Censor()`, `Extract()`, `Count()` tinggal pakai data itu tanpa scan ulang. Ini yang bikin Censor cepat — replace in-place berdasarkan index yang udah ada.

**Length histogram** — token yang panjangnya ga ada di dictionary langsung di-skip sebelum map lookup. Hemat operasi ga perlu.

### Sliding Window Affix Detection

Kata berimbuhan? Bacot geser-geser kayak swipe Tinder:

```
Input: "mebabi"
├── "mebabi" → cek dictionary? gak ada
├── "ebabi" → coba-coba
│   └── "babi" → 🚨 DITEMUKAN!
```

Ditambah validasi nasal fusion: `memukul` (mem- + pukul, p luluh), `menggali` (meng- + gali, g tetap). Dan filter false positive: sisa 1 suku kata setelah stem = kemungkinan kata udah beda makna.

### RecursiveScan

```
Input: "xbabi"
├── xbabi → gak ada
├── babi  → 🚨 KETANGKAP!
```

*Gak perlu ML, gak perlu neural network. Cuma geser-geser doang.*

---

## 🏎️ Benchmark

Intel i5-13400, Linux, Go 1.26.4.

| Skala | Operasi | Waktu |
|-------|---------|-------|
| 🟢 Ringan | `Contains` lookup | 5.7 ns |
| 🟢 Ringan | `Scan` kalimat | 531 ns |
| 🟢 Ringan | `Censor` kata | 0.89 ns |
| 🟡 Sedang | `Scan` paragraf (5 kalimat) | 10.6 µs |
| 🔴 Berat | `Scan` 10Kb (1000 kata) | 344 µs |
| 🔴 Berat | `Censor` 10Kb (4000 match) | 172 µs |
| ⚫ Ekstrim | `Scan` 100Kb | 20 ms |

**Singkatnya:** 2 juta kalimat per detik. Bahkan skenario terberat tetep sub-millisecond. Intel i5 aja nggak kerasa, apalagi server production.

<details>
<summary>📊 Detail benchmark</summary>

```
goos: linux | goarch: amd64 | cpu: 13th Gen Intel(R) Core(TM) i5-13400

--- RINGAN ---
Contains-16                     217,744,636    5.711 ns/op       0 B/op    0 allocs/op
AddWord-16                        1,960,904   607.800 ns/op     280 B/op    7 allocs/op
TextPreprocessing-16              3,364,468   348.700 ns/op     160 B/op    7 allocs/op
ScanSingleWord-16                 5,298,542   221.700 ns/op     168 B/op    6 allocs/op
Scan-16                           2,273,468   531.400 ns/op     192 B/op    7 allocs/op
RecursiveScan-16                  2,765,752   433.400 ns/op     192 B/op    7 allocs/op
Censor-16                      1,000,000,000   0.8916 ns/op       0 B/op    0 allocs/op
CensorSentence-16                 6,700,214   181.500 ns/op      32 B/op    1 allocs/op

--- SEDANG ---
ScanParagraph-16                    114,300  10,681 ns/op   2,744 B/op   23 allocs/op
RecursiveScanParagraph-16           129,304   9,309 ns/op   2,744 B/op   23 allocs/op
CensorParagraph-16                  224,820   5,680 ns/op   3,704 B/op   10 allocs/op
ScanCleanLong-16                    165,366   6,979 ns/op   2,144 B/op   16 allocs/op
ScanManyStopWords-16                171,465   6,520 ns/op   2,144 B/op   16 allocs/op
ScanAffixLong-16                    642,274   1,829 ns/op   1,160 B/op   16 allocs/op
ScanAllLeet-16                      193,651   5,823 ns/op   1,120 B/op   14 allocs/op
ScanHeavyUnstack-16               1,269,277   942.6 ns/op     600 B/op   17 allocs/op
ScanMixedAffixLeet-16               879,368   1,312 ns/op     352 B/op   10 allocs/op
RecursiveScanNoMatch-16             478,560   2,376 ns/op     352 B/op   10 allocs/op
CensorLeet-16                  1,000,000,000   0.8819 ns/op       0 B/op    0 allocs/op

--- BERAT ---
Scan10Kb-16                           3,390 344,589 ns/op 200,666 B/op  1,448 allocs/op
RecursiveScan10Kb-16                  3,636 346,540 ns/op 203,865 B/op  1,548 allocs/op
Censor10Kb-16                         6,565 207,774 ns/op 144,889 B/op     19 allocs/op
Scan10KbAllLeet-16                    3,864 347,040 ns/op 126,049 B/op     36 allocs/op
TextPreprocessing10Kb-16            10,000 109,373 ns/op 125,985 B/op     35 allocs/op
Scan10KbContinuous-16                10,084 116,301 ns/op 139,865 B/op     50 allocs/op
RecursiveScan10KbContinuous-16        8,287 131,111 ns/op 139,905 B/op     52 allocs/op
Scan100Kb-16                             57 20,040,838 ns/op 8,255,097 B/op 66 allocs/op
RecursiveScan100Kb-16                    62 18,785,788 ns/op 8,255,106 B/op 66 allocs/op
ScanAllProfanity-16                   2,328 457,570 ns/op 460,508 B/op  4,055 allocs/op
RecursiveScanAllProfanity-16          2,036 573,975 ns/op 460,507 B/op  4,055 allocs/op
CensorAllProfanity-16                 7,722 147,106 ns/op 122,880 B/op      2 allocs/op
RecursiveScanLongSingleWord-16        1,544 744,142 ns/op  93,321 B/op     36 allocs/op
ScanLongSingleWord-16                14,616  81,681 ns/op  93,280 B/op     34 allocs/op
```
</details>

---

## ❓ FAQ

### Kodenya sederhana, emang work?

85 test, 34 benchmark, zero dep. Kadang solusi terbaik ya yang sederhana.

### Dictionary cuma 190 kata?

200+ kata dasar. Tapi tiap kata otomatis di-generate varian imbuhan jadi ribuan (mebabi, pebabi, dibabi, menjing, penganjing, dll). Kalau kurang tinggal `AddWord()`. PR diterima. *Makin kotor makin kita senang.*

### Thread safe?

Buat instance baru per goroutine atau bungkus `sync.Mutex`.

### Kok ada method typo?

(`TrimSpce`, `SanitazeReadSign`, `Majorty`) — iya. Backward compatibility. Mungkin next major version dibenerin. Atau lo benerin sendiri.

### Bisa dipake production?

85 test, 34 benchmark, zero dep. Nggak bakal numpukin memory server lo.

---

## 🙏 Credits

- **Dictionary:** [indonesian-badwords](https://github.com/drizki/indonesian-badwords) — *makasih bang udah dikumpulin.*
- **Stop Words:** [go-sastrawi](https://github.com/truthfulses/go-sastrawi) — *pinjem, makasi.*

Kontribusi kata kasar baru? Buka aja PR. 😄

---

## ⚖️ License

[MIT](https://choosealicense.com/licenses/mit/)
