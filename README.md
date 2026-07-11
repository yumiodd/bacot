![bacot.png](./bacot.png)
# Bacot

Filter kata kasar untuk Bahasa Indonesia.

```
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true
b.Text("babi dan anjing").Collect(true).Scan().Censor() // "**** dan ******"
```

Bahasa Indonesia itu aglutinatif. `babi` bisa jadi `mebabi`, `pebabi`, `dibabi`, `kebabi`, `memperbabikan`. Coba lo handle pake `strings.Contains` sendirian — *good luck, you'll need it.*

Bacot lahir karena capek liat developer nulis regex panjang 20 baris cuma buat deteksi "4njiiing". *Kami sudah melalui masa gelap itu (bohong). Kami tidak akan kembali (bohong lagi).*

## Tidak klise pake regex aja

| Masalah | Contoh | Regex | Bacot |
|---------|--------|-------|-------|
| Imbuhan | `mebabi`, `penganjing` | lo harus bikin pattern buat tiap imbuhan | otomatis |
| Leet speak | `4njing`, `8abi`, `k0nt0l` | `[4a][nN]...` — capek? | built-in |
| Unstack | `anjiiiiing` | `anj+ing` — gampang sih, tapi itu baru 1 kata | otomatis |
| False positive | `babiru` bukan `babi` | lo bikin pengecualian? ya nulis lagi | udah di-handle |
| 10Kb censor | 4000 kata kasar | 5-20ms, rerata CPU merintih | 172µs, CPU nggak sadar |

## Instalasi

```go
import bacot "github.com/yumiodd/bacot/src"
```

Nggak ada `go get` lagi ini itu. ga ada `requirements.txt` instalasi misterius. Zero dep (*untuk sekarang*)

*kami bukan python*

## Quick Start

```go
b := bacot.New()

// Cek — biar tau
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true

// Sensor — biar aman
b.Text("babi dan anjing").Collect(true).Scan().Censor() // "**** dan ******"

// Hitung — biar puas
res := b.Text("babi dan anjing asu kontol").Collect(true).Scan()
res.Count()   // 4
res.Extract() // ["babi", "anjing", "asu", "kontol"]
res.First()   // "babi"
res.Last()    // "kontol"
```

> `New()` itu bikin dictionary + generate ribuan varian imbuhan. Jangan dipanggil tiap scan — reuse instance. *Ini bukan `new StringBuilder()` tiap loop.*

## Fitur

### Affix-aware

"mebabi" itu "babi" pake imbuhan `me-`. `penganjing` itu `peng-` + `anjing`. Bacot tau ini karena dia pake **sliding window**: stripping prefix Indonesia satu per satu sampe ketemu kata dasarnya.

```go
b.Text("mebabi").Scan().IsProfane()      // true
b.Text("mebabi").Affix(false).Scan()...  // false — kalo mau exact match aja
```

Yang didukung: `me-`, `pe-`, `di-`, `te-`, `be-`, `ber-`, `ter-`, `per-`, `meng-`, `peng-`, `men-`, `pen-`, `meny-`, `peny-`, `mem-`, `pem-`, `ng-`, `ny-`.

Bonus: **nasal fusion**. `memukul` itu `mem-` + `pukul` (p-nya luluh). Di kamus Bacot cuma ada `pukul`, tapi `memukul` tetep kedeteksi. *Seperti mantan — apapun kamu kedoknya, tetap ketahuan.*

### Leet Speak

`4 → a`, `8 → b`, `3 → e`, `@ → a`, `$ → s`, `! → i`, `0 → o`, `7 → t`, `5 → s`, `# → h`, dan masih banyak lagi.

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true
// 4 → a | njing → anjing | match ✅
```

Leet speak nggak aktif default karena butuh extra pass. Tapi kalo lo tau user base lo banyak pake leet, aktifkan aja.

### Unstack Char

```go
b.Text("anjiiiiing").Scan().IsProfane() // true
// anji{i}ng → anjing | match ✅
```

User kira pake huruf berulang lolos? *Apakah ini pertama kalinya lo berhadapan dengan internet Indonesia?*

### Recursive Scan

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

Nyari substring match di dalam token. Nge-scan dari tiap posisi karakter di kata. Jadi kalau user lo kreatif nambahin prefiks acak — tetap kena.

Konsekuensi: O(token_length × dict_variants). Buat kata normal (<20 char) — nggak kerasa. Buat teks 10Kb tanpa spasi — ya agak mikir dikit. *Worst case complexity: istilah keren buat "bakal lebih lambat dikit".*

### Custom Dictionary

```go
b.AddWord("pukul")     // otomatis detect memukul, dipukul, pemukul, mukul
b.AddWord("setan")     // + nyetan, kesetanan, dll
b.Dict.AddWords("custom")
b.Dict.DelWords("anjing")
```

Istilah baru muncul tiap hari. Lo tinggal `AddWord()`, kami urus varian imbuhannya.

### Config — biar nggak repot chaining manual tiap scan

```go
b.Config(&bacot.ModalScanConfig{
    Affix:   true,
    Collect: true,
    Order: []bacot.SanitizeOrder{
        bacot.WithLeetSpeak,
        bacot.UnstackChar,
    },
})
b.Text("4njiiing").Scan().IsProfane() // true — leet + unstack otomatis
```

## Cara Kerja

Pipeline dari input sampe hasil:

```
Input → ToLower → [Sanitize] → Scan → ScanResult ─→ Censor()
                                                     ├→ IsProfane()
                                                     ├→ Extract()
                                                     └→ Count()
```

Arsitektur:

1. **Precomputed WordIndex.** `Scan()` itu bagian paling berat — dia nge-tokenize, stripping affix, ngecek dictionary. Tapi hasilnya adalah `[]WordIndex` yang nyimpen posisi start/end tiap kata kotor. `Censor()` tinggal replace `*` di posisi itu — O(n) per kata, O(1) per replace. Nggak perlu scan ulang. *Bayangin kayak index di database — alokasi sekali, query berkali-kali.*

2. **Length histogram.** Sebelum map lookup, Bacot ngecek: "apa ada kata di dictionary dengan panjang segini?". Kalau nggak ada, token langsung di-skip. Ini pre-filter O(log n) yang ngurangin map lookup nggak perlu. *Micro-optimization yang hasilnya keliatan di 10Kb.*

3. **Affix variant generation (`craftMan`).** Di `NewDictionary()`, setiap kata dasar di-generate semua varian imbuhan yang mungkin. Inilah kenapa `New()` agak berat (beberapa ratus µs) tapi scan jadi cepat — cukup exact match di map, nggak perlu stemming di runtime. *Tradeoff klasik: bayar sekali di init, gratis selamanya.*

4. **False positive filter.** Kalau setelah stripping prefix ditemukan match, tapi sisa stringnya cuma 1 suku kata — kemungkinan itu bukan kata imbuhan beneran. Contoh: `"babiru"` → strip `me-` gagal → strip `di-` gagal → ambil stem `"babi"` (ditemukan!) → sisa `"ru"` = 1 suku kata → skip. *Heuristic sederhana yang manjur buat ngejaga akurasi.*

## Benchmark

Intel i5-13400, Linux, Go 1.26.4.

| Level | Operasi | Waktu |
|-------|---------|-------|
| Ringan | `Contains` lookup aja | 5.7 ns |
| Ringan | `Scan` kalimat pendek | 531 ns |
| Ringan | `Censor` kata tunggal | 0.89 ns |
| Sedang | `Scan` paragraf 5 kalimat | 10.6 µs |
| Berat | `Scan` 10Kb (1000 kata) | 344 µs |
| Berat | `Censor` 10Kb (4000 match) | 172 µs |
| Ekstrim | `Scan` 100Kb lorem | 20 ms |

*"2 juta kalimat per detik."* — 

<details>
<summary>📊 Raw benchmark — buat yang demen angka</summary>

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

**Dictionary cuma 200 kata?**  
200+ kata dasar. Tiap kata dasar generate varian imbuhan — jadinya ribuan. Kurang? `AddWord()`. Atau buka PR. *Makin kotor makin kita senang.*,
gunakan `AddFalsePosive()` untuk menambal kesalahan *hehe*

**Thread safe?**  
Enggak. `Bacot` pake pointer sharing. Lo mau pake di banyak goroutine? Buat instance baru tiap goroutine atau bungkus `sync.Mutex`. *nekad, resiko tanggung sendiri*

**Cocok buat production?**  
Well coba aja dulu.

## License

[MIT](https://choosealicense.com/licenses/mit/)
