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

res := b.Text("4njiiing").WithLeetSpeak().Scan()
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
| **Custom words** | `AddWord(words...)` |
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
b.AddWord("custombad")           // + varian craftMan
b.AddWord("pukul")               // + imbuhan otomatis
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

Leet speak **tidak aktif** secara default. Aktifkan manual:

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true
```

Unstack character otomatis (hapus huruf berulang):

```go
b.Text("anjiiing").Scan().IsProfane()   // true (unstack: anjiiing → anjing)
```

Kombinasi leet + unstack:

```go
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true
```

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

### Config — *"Sekali Atur, Scan Enteng"*

Buat standarisasi pipeline tanpa chaining manual tiap scan pake `ModalScanConfig`.
**Urutan `Order` penting** — hasil akhir tergantung urutan preprocessing.

```go
b := bacot.New()

b.Config(&bacot.ModalScanConfig{
    Affix:   true,
    Collect: true,
    Order: []bacot.SanitizeOrder{
        bacot.WithLeetSpeak,           // 1. leet dulu
        bacot.UnstackChar,             // 2. baru unstack
    },
})

res := b.Text("4njiiing").Scan()
res.IsProfane() // true (leet + unstack otomatis)
```

Urutan di atas: leet dulu (`4` → `a`), baru unstack (hapus huruf berulang).
Config ini berlaku untuk semua panggilan `Text()` setelahnya sampai dipanggil ulang.

Pipeline yang bisa diatur:

| Order | Fungsi |
|-------|--------|
| `bacot.SanitizeNewLine` | Hapus newline |
| `bacot.TrimSpace` | Hapus spasi antar huruf |
| `bacot.WithLeetSpeak` | Konversi leet (`4` → `a`, dll) |
| `bacot.UnstackChar` | Hapus huruf berulang |

Chaining manual tetep jalan kapan aja, dan **gak** terpengaruh config:

```go
b.Text("4 n j i n g").WithLeetSpeak().TrimSpace().Scan().IsProfane() // true
```

## How It Works

```
Input → ToLower → SanitizeNewLine → UnstackChar → [TrimSpace] → Scan / RecursiveScan
                → [WithLeetSpeak] (manual)
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
| 🟢 | `Contains` | lookup | **5.7 ns** |
| 🟢 | `ScanSingleWord` | `"anjing"` | **221 ns** |
| 🟢 | `Scan` | `"aku suka 4njiiing"` | **531 ns** |
| 🟢 | `RecursiveScan` | `"aku suka 4njiiing"` | **433 ns** |
| 🟢 | `Censor` | `"anjing"` | **0.89 ns** |
| 🟡 | `AddWord` | 1 kata baru | **607 ns** |
| 🟡 | `ScanParagraph` | 5 kalimat leet+affix | **10.6 µs** |
| 🟡 | `ScanMixedAffixLeet` | 6 kata campur | **1.3 µs** |
| 🟡 | `ScanAffixLong` | 20 kata imbuhan | **1.8 µs** |
| 🔴 | `Scan10Kb` | 1.000 kata | **344 µs** |
| 🔴 | `RecursiveScan10Kb` | 1.000 kata | **346 µs** |
| ⚫ | `Scan100Kb` | 100KB lorem | **20 ms** |
| ⚫ | `RecursiveScanLongSingleWord` | 1 kata 10K char | **744 µs** |
| ⚫ | `ScanAllProfanity` | 10KB, 4.000 match | **457 µs** |

<details>
<summary>📊 Raw benchmark data</summary>

```
goos: linux | goarch: amd64 | cpu: 13th Gen Intel(R) Core(TM) i5-13400

--- RINGAN ---
Contains-16                     217,744,636    5.711 ns/op       0 B/op    0 allocs/op
BenchmarkAddWord-16               1,960,904   607.800 ns/op     280 B/op    7 allocs/op
BenchmarkTextPreprocessing-16     3,364,468   348.700 ns/op     160 B/op    7 allocs/op
BenchmarkScanSingleWord-16        5,298,542   221.700 ns/op     168 B/op    6 allocs/op
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

**Singkatnya:** Bacot bisa scan ~2 juta kalimat per detik. Bahkan di skenario terberat — 10KB teks, 4.000 kata kasar, 1 kata sepanjang 10K karakter — semuanya masih sub-millisecond. Intel i5 aja nggak kerasa, apalagi server production.

## FAQ

### Kodenya sederhana, emang work?

Work. 87 test, 35 benchmark, zero dep. Kadang solusi terbaik ya yang sederhana.

### Dictionary cuma 190 kata?

Tambahin aja. PR diterima. *Makin kotor makin kita senang.*

## Credits

- Dictionary: [indonesian-badwords](https://github.com/drizki/indonesian-badwords) *makasih dah ngumpuli*
- Stop words: [go-sastrawi](https://github.com/truthfulses/go-sastrawi) *pinjem, makasih*

Kontribusi kata kasar baru? Buka aja PR. 😄

## License

[MIT](https://choosealicense.com/licenses/mit/)
