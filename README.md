![bacot banner](./bacot.png)

# Bacot

*"Mulut lo bacot — tapi bacot filter kata-kata lo."*

Profanity filter untuk Bahasa Indonesia dengan affix detection, leet speak obfuscation, dan unstack character. Zero dependency.

Bahasa Indonesia aglutinatif — satu kata dasar bisa berubah dengan imbuhan (`babi` → `mebabi`, `pebabi`, `dibabi`, `kebabi`). Bacot handle ini otomatis tanpa perlu daftar semua variasi imbuhan. *Capek kali, masa "babi" doang, "mebabi" juga, "pebabi" lagi.*

## Installation

```go
import bacot "github.com/yumiodd/bacot/src"
```

Itu aja, gak ada dependency lain. *Kami bukan npm.*

## Quick Start

```go
b := bacot.New()

res := b.Text("4njiiing").Scan()
res.IsProfane() // true
res.Censor()    // "*******"

res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
```

> **Satu** instance untuk banyak scan. Jangan `New()` tiap kali — ini bukan sikat gigi.

## Features

| Feature | Description |
|---------|-------------|
| **Affix-aware** | Deteksi kata berimbuhan (`mebabi`, `penganjing`) |
| **Leet speak** | `4njing` → `anjing`, `8abi` → `babi` |
| **Unstack char** | `anjiiing` → `anjing` |
| **Two scan modes** | `Scan()` per-word, `RecursiveScan()` sampe ke dalem |
| **Custom words** | `AddWord(affix bool, words...)` |
| **Censor** | Mapping otomatis leet → `*` |
| **Zero dependencies** | Hanya stdlib Go |
| **Collection mode** | Kumpulin semua match sekaligus |

## Usage

### Basic Detection

```go
res := b.Text("anjing").Scan()
res.IsProfane() // true
res.Censor()    // "******"
res.Extract()   // ["anjing"]
```

### Collect All — *"Sekalian Dihajar"*

```go
res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
```

Default `Collect(false)` — cukup tau satu doang, kayak mantan.

### Custom Words — *"Update Kamus, Update Perasaan"*

```go
b.AddWord(false, "custombad")
b.AddWord(true, "pukul") // + affix variants
```

Istilah baru bermunculan tiap hari? Tambahin aja. *Udah biasa.*

### Affix Mode

```go
b.Text("mebabi").Scan().IsProfane()      // true (ya jelas)
b.Text("mebabi").Affix(false).Scan().IsProfane() // false
```

Di Indonesia, "mebabi" tetaplah "babi". *Jangan munafik.*

### Recursive Scan

Nyari sampe ke substring — buat yang suka nyelipin karakter.

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Leet Speak & Unstack

```go
b.Text("4njing").Scan().IsProfane()     // true
b.Text("anjiiing").Scan().IsProfane()   // true
b.Text("4njiiing").Scan().IsProfane()   // true
```

Aktif default, bisa dimatiin kalo emang mau.

### Trim Space — *"Pake Spasi Biar Lolos?"*

```go
b.Text("a n j i n g").TrimSpace().Scan().IsProfane() // true
```

*Nice try.*

### Dictionary

```go
b.Dict.Contains("anjing")    // true
b.Dict.IsStopWord("dan")     // true
b.Dict.AddWords("custom1", "custom2")
b.Dict.DelWords("anjing")
b.Dict.GetDict()             // map[string]struct{}
```

### Generator / Iterate Results

```go
gen := res.Generator()
if w := gen.Yield(); w != nil {
    fmt.Println(w.Word)
}
```

## How It Works

```
Input → ToLower → LeetSpeak → UnstackChar → [TrimSpace] → Scan / RecursiveScan
```

### Sliding Window Affix Detection

Kata berimbuhan? Bacot geser-geser kayak lagi swipe Tinder:

```
mebabi → ebabi → babi → DITEMUKAN ✅
```

Setiap posisi `l` diambil substring `w[l:]`, dicocokkan dengan dictionary. Pas ketemu, selesai.

### RecursiveScan

```
xbabi → babi → DITEMUKAN ✅
```

Simple. Efektif. *Gak perlu pake ML segala.*

## Benchmark

Intel i5-13400, Linux, Go 1.26.4.

| Skala | Operasi | Input | Waktu |
|-------|---------|-------|-------|
| 🟢 | `Contains` | lookup | **5 ns** |
| 🟢 | `Scan` | `"aku suka 4njiiing"` | **483 ns** |
| 🟢 | `RecursiveScan` | `"aku suka 4njiiing"` | **520 ns** |
| 🟡 | `ScanParagraph` | 5 kalimat leet+affix | **10 µs** |
| 🟡 | `ScanMixedAffixLeet` | 6 kata campur | **1 µs** |
| 🔴 | `Scan10Kb` | 1.000 kata | **411 µs** |
| 🔴 | `RecursiveScan10Kb` | 1.000 kata | **444 µs** |
| ⚫ | `Scan100Kb` | 100KB lorem | **17 ms** |
| ⚫ | `RecursiveScanLongSingleWord` | 1 kata 10K char | **758 µs** |
| ⚫ | `ScanAllProfanity` | 10KB, 4.000 match | **523 µs** |

<details>
<summary>📊 Raw benchmark data</summary>

```
goos: linux | goarch: amd64 | cpu: 13th Gen Intel(R) Core(TM) i5-13400

--- RINGAN ---
Contains-16                  43,944,686    5.5 ns/op       0 B/op    0 allocs/op
ScanSingleWord-16             1,000,000  216.0 ns/op     168 B/op    6 allocs/op
Scan-16                         435,322  492.7 ns/op     232 B/op    9 allocs/op
RecursiveScan-16                431,160  534.3 ns/op     232 B/op    9 allocs/op
Censor-16                     1,234,800  189.0 ns/op      56 B/op    3 allocs/op

--- SEDANG ---
ScanParagraph-16                 25,159   9,065 ns/op   3,160 B/op   32 allocs/op
ScanMixedAffixLeet-16           248,612   1,028 ns/op     392 B/op   12 allocs/op
ScanAffixLong-16                 86,047   2,577 ns/op   1,160 B/op   16 allocs/op

--- BERAT ---
Scan10Kb-16                         570  404,140 ns/op  244,697 B/op  2,249 allocs/op
RecursiveScan10Kb-16                464  443,880 ns/op  244,697 B/op  2,249 allocs/op
Scan10KbAllLeet-16                  768  289,176 ns/op  225,498 B/op  1,649 allocs/op
Scan100Kb-16                         13 17,101,498 ns/op  8,255 B/op     66 allocs/op
RecursiveScanLongSingleWord-16      316  757,932 ns/op   93,320 B/op     36 allocs/op
ScanAllProfanity-16                 451  522,762 ns/op  460,507 B/op  4,055 allocs/op
```
</details>

**Singkatnya:** Bacot bisa scan ~2 juta kalimat per detik. Bahkan di skenario terberat — 10KB teks, 4.000 kata kasar, 1 kata sepanjang 10K karakter — semuanya masih sub-millisecond. Intel i5 aja nggak kerasa, apalagi server production.

## FAQ

### Kodenya sederhana, emang work?

Work. 74 test, 33 benchmark, zero dep. Kadang solusi terbaik ya yang sederhana.

### Dictionary cuma 190 kata?

Tambahin aja. PR diterima. *Makin kotor makin kita senang.*

## Credits

- Dictionary: [indonesian-badwords](https://github.com/drizki/indonesian-badwords) *makasih dah ngumpuli*
- Stop words: [go-sastrawi](https://github.com/truthfulses/go-sastrawi) *pinjem, makasih*

Kontribusi kata kasar baru? Buka aja PR. 😄

## License

[MIT](https://choosealicense.com/licenses/mit/)
