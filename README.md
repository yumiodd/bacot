![bacot banner](./bacot.png)

# Bacot

*"Mulut lo bacot — tapi bacot filter kata-kata lo."*

Profanity filter untuk Bahasa Indonesia dengan affix detection, leet speak, unstack character. Zero dependency.

Bahasa Indonesia aglutinatif — satu kata dasar bisa berubah dengan berbagai imbuhan (`babi` → `mebabi`, `pebabi`, `dibabi`, `kebabi`). Bacot handle semua varian ini otomatis tanpa perlu daftarin satu-satu.

```go
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true
b.Text("babi dan anjing").Collect(true).Scan().Censor() // "**** dan ******"
```

## Fitur

| Fitur | Contoh | Keterangan |
|-------|--------|------------|
| **Affix-aware** | `mebabi`, `penganjing` | Deteksi kata berimbuhan otomatis (default ON) |
| **Leet speak** | `4njing` → `anjing`, `8abi` → `babi` | 25+ karakter mapping, aktif manual |
| **Unstack char** | `anjiiing` → `anjing` | Hapus huruf berulang (default ON) |
| **Recursive scan** | `xbabi` → ketemu | Nyari sampe ke dalam substring |
| **Two scan modes** | `Scan()` / `RecursiveScan()` | Per word atau sliding window per token |
| **Collection mode** | `Collect(true)` | Kumpulin semua match atau berhenti di pertama |
| **Custom dictionary** | `AddWord()`, `DelWord()` | Tambah atau hapus kata kapan aja |
| **Censor** | `Censor()` → `******` | Replace character berdasarkan posisi |
| **False positive filter** | "babiru" bukan "babi" | Deteksi 1 suku kata sisa setelah stem |
| **Zero dependencies** | — | Hanya stdlib Go |

## Quick Start

```go
import bacot "github.com/yumiodd/bacot/src"

b := bacot.New()

// Deteksi + sensor
res := b.Text("4njiiing").WithLeetSpeak().Scan()
res.IsProfane() // true
res.Censor()    // "*******"

// Kumpulin semua match
res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
res.First()     // "babi"
res.Last()      // "anjing"
```

> **Satu** instance untuk banyak scan. Jangan `New()` tiap kali.

## Usage

### Detection & Censor

```go
res := b.Text("anjing").Scan()
res.IsProfane() // true
res.Censor()    // "******"
res.Extract()   // ["anjing"]
```

### Affix / Imbuhan

```go
b.Text("mebabi").Scan().IsProfane()      // true (affix default ON)
b.Text("mebabi").Affix(false).Scan()...  // false
```

### Leet Speak

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true
```

### Unstack + Leet

```go
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true
```

### Recursive Scan

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Trim Space — *"Pake Spasi Biar Lolos?"*

```go
b.Text("a n j i n g").TrimSpace().Scan().IsProfane() // true
```

### Config — Sekali Atur

```go
b.Config(&bacot.ModalScanConfig{
    Affix:   true,
    Collect: true,
    Order: []bacot.SanitizeOrder{
        bacot.WithLeetSpeak,
        bacot.UnstackChar,
    },
})
b.Text("4njiiing").Scan().IsProfane() // true, otomatis
```

### Dictionary

```go
b.Dict.AddWords("custom")
b.Dict.DelWords("anjing")
b.Dict.Contains("anjing")   // true
b.Dict.IsStopWord("dan")    // true
```

## Bagaimana Cara Kerjanya

```
Input → ToLower → [Sanitize] → Scan → ScanResult ─→ Censor()
                                                   ├→ IsProfane()
                                                   ├→ Extract()
                                                   ├→ Count()
                                                   └→ Generator()
```

**Precomputed WordIndex** — `Scan()` menghitung posisi start/end tiap kata kotor sekali. Operasi setelahnya (`Censor`, `Extract`, `Count`) tinggal pakai data itu tanpa scan ulang.

**Length histogram** — token yang panjangnya ga ada di dictionary langsung di-skip sebelum map lookup. Mengurangi operasi yang tidak perlu.

**Affix detection** — stripping prefix Indonesia (meng-, mem-, men-, ber-, ter-, per-, me-, pe-, di-, dll) dengan validasi nasal fusion. Contoh: `memukul` (p luluh jadi m), `menggali` (g tetap). Ditambah filter false positive: kalau sisa setelah stem cuma 1 suku kata, kemungkinan kata sudah beda makna.

## Benchmark

Intel i5-13400, Linux, Go 1.26.4.

| Skala | Operasi | Waktu |
|-------|---------|-------|
| 🟢 | `Contains` lookup | 5.7 ns |
| 🟢 | `Scan` kalimat pendek | 531 ns |
| 🟢 | `Censor` kata | 0.89 ns |
| 🟡 | `Scan` paragraf (5 kalimat) | 10.6 µs |
| 🔴 | `Scan` 10Kb (1000 kata) | 344 µs |
| 🔴 | `Censor` 10Kb (4000 match) | 172 µs |
| ⚫ | `Scan` 100Kb | 20 ms |

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

## FAQ

### Kodenya sederhana, emang work?

85 test, 34 benchmark, zero dep. Kadang solusi terbaik ya yang sederhana.

### Dictionary cuma 190 kata?

200+ kata dasar. Tapi tiap kata otomatis di-generate varian imbuhannya jadi ribuan (mebabi, pebabi, dibabi, menjing, penganjing, dll). Kalau kurang tinggal `AddWord()`. PR juga diterima.

### Thread safe?

Tidak — no mutex, no sync. Buat instance baru per goroutine, atau bungkus pake `sync.Mutex`.

### Kok ada method typo?

(`TrimSpce`, `SanitazeReadSign`, `Majorty`) — iya. Backward compatibility. Mungkin next major version dibenerin.

## Credits

- Dictionary: [indonesian-badwords](https://github.com/drizki/indonesian-badwords)
- Stop words: [go-sastrawi](https://github.com/truthfulses/go-sastrawi)

Kontribusi kata kasar baru? Buka aja PR. 😄

## License

[MIT](https://choosealicense.com/licenses/mit/)
