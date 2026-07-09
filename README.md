![bacot banner](./bacot.png)

# Bacot
### *"Mulut lo bacot — tapi bacot filter kata-kata lo."*

**Bacot** adalah profanity filter untuk Bahasa Indonesia yang handle affix (imbuhan), leet speak, dan unstack character. Zero dependency. Sekali `go get` langsung jalan, tanpa drama.

---

## 🧐 Kenapa Ada Library Kayak Gini?

Karena internet Indonesia itu **liar**.

Ada user komen "4njiiing b4ngs4t" di kolom komentar. Ada yang chat "m3m3k" di grup keluarga. Ada yang bales "k0nt0l" di tweet lo.

Lo sebagai developer punya dua pilihan:
1. **Pasrah** — aplikasi lo jadi tong sampah digital 🗑️
2. **Pasang Bacot** — biar filter otomatis, lo tidur nyenyak

Kami sarankan opsi 2. Tapi terserah lo, ini hidup lo.

---

## 🔧 Installation

```go
import bacot "github.com/yumiodd/bacot/src"
```

Itu aja. Tidak ada `go get` 20 library, tidak ada `npm install` pusing-pusing, tidak ada `requirements.txt` penuh misteri.

*Kami bukan Python.*

---

## 🚀 Quick Start

```go
b := bacot.New()

// Lo mau ngasih lihat ke user yang sotoy?
res := b.Text("4njiiing").Scan()
res.IsProfane() // true — kena. Bablas? Nggak.
res.Censor()    // "*******" — aman disajikan ke publik

// Kalo mau sekalian dikumpulin semua dosanya
res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
```

> **⚠️ PENTING:** `bacot.New()` itu generate dictionary plus affix variants. Berdasarkan penelitian ilmiah (dan juga pengalaman pribadi), **jangan `New()` tiap scan**. Buat satu instance, pakai berulang. Ini bukan sikat gigi, gak perlu ganti tiap hari.

---

## 📋 Fitur (Biar Keliatan Profesional)

| Fitur | Kerjaannya |
|-------|-----------|
| **Affix-aware** | `mebabi`? `penganjing`? `kebacotan`? Semua kena. Karena Bacot ngerti Bahasa Indonesia pake imbuhan. |
| **Leet speak** | `4njing` → `anjing`. `8abi` → `babi`. `5ialan` → `sialan`. *Kami bukan polisi bahasa, kami lebih dari itu.* |
| **Unstack char** | `anjiiiiiiiiing` → `anjing`. Lo kira pake karakter berulang lolos? *Nice try, bocah.* |
| **Dua mode scan** | `Scan()` = per kata. `RecursiveScan()` = substring, sampe ke akar-akarnya. |
| **Custom dictionary** | Lo punya kata kesayangan? `AddWord()` aja. |
| **Censor otomatis** | Mapping leet speak → `*`. User nulis `4njing`, sensor balikin `******`. Puas? |
| **Zero dependencies** | Cuma stdlib Go. **Zero. Kosong. Nothing. Nihil.** Kalau lo nemu dependency selain stdlib, *itu halusinasi lo.* |
| **Collection mode** | Mau dikumpulin semua match? `Collect(true)`. Biar sekalian dihajar. |

---

## 📖 Cara Pake (Biar Gak Malu sama Teknisi)

### 🟢 Basic Detection
Deteksi dasar, tinggal tulis, scan, beres.

```go
b := bacot.New()
res := b.Text("anjing").Scan()
res.IsProfane() // true (surprise!)
res.Censor()    // "******"
res.Extract()   // ["anjing"]
```

### 🟢 Collect All Mode — *"Kumpulin Semua Biar Puas"*
Mau tau total kata kasar sekaligus? Pake collect mode.

```go
res := b.Text("babi dan anjing").Collect(true).Scan()
res.Count()     // 2
res.Extract()   // ["babi", "anjing"]
res.First()     // "babi"
res.Last()      // "anjing"
```

Default `Collect(false)` = berhenti di match pertama. **Seperti mantan: tau satu aja udah cukup, sisanya sakit hati.**

### 🟢 Custom Words — *"Rizky Febian Baru Aja Rilis Single Baru"*
Tiap hari ada aja istilah baru. Santai, tinggal ditambah.

```go
b.AddWord(false, "umiakucihanya")   // exact match
b.AddWord(true, "cepet")            // + affix variants (encepet, dicepet, dst)
```

**Ada kata kasar baru lagi viral?** Lo tinggal tambahin. *Go-Live sambil update dictionary, resiko ditanggung pengguna.*

### 🟢 Affix Mode
Kata "mebabi" itu sebenernya "babi" pake imbuhan `me-`. Di Indonesia, kenyataan ini pahit tapi harus diterima.

```go
// Affix ON = default
b.Text("mebabi").Scan().IsProfane()      // true

// Affix OFF = cuma exact match
b.Text("mebabi").Affix(false).Scan().IsProfane() // false
```

Lo bisa matiin affix kalo mau cari exact match doang. Tapi inget, **di Indonesia, "mebabi" tetep "babi"**. Jangan munafik.

### 🟢 Recursive Scan — *"Nakal Kamu, Nyelip terus"*
Punya user kreatif yang nyisipin karakter biar lolos filter? Kita punya *RecursiveScan*.

```go
b.Text("xbabi").Scan().IsProfane()           // false (katanya sih lolos)
b.Text("xbabi").RecursiveScan().IsProfane()  // true (maaf, nyata)
```

**Cara kerja:** Bacot bakal nyari setiap kemungkinan substring. Mau lo taro apapun di depan/belakang kata kasar, tetap ketahuan. *Udah, nyerah aja.*

### 🟢 Leet Speak & Unstack — *"4njiiing b4ngs4t, Bacot Tetap Tahu"*
Ini andalan Bacot. User ngetik apa pun, tetap terfilter.

```go
b.Text("4njing").Scan().IsProfane()     // true (4 → a)
b.Text("anjiiing").Scan().IsProfane()   // true (akar kata anjing)
b.Text("4njiiing").Scan().IsProfane()   // true (kombinasi maut, tetap kena)
```

**Kami udah liat semua trik kotor anak internet Indonesia.** Pernah main chatojol? Kami berasal dari dunia sana. *Kami tau kelakuan lo.*

### 🟢 Trim Space — *"Pake Spasi? Maaf, 2012 Memaanggil, Trik Lo Udah Kedaluwarsa"*
Coba aja kalo berani.

```go
b.Text("a n j i n g").TrimSpace().Scan().IsProfane() // true
```

`a n j i n g`? Dipisah pake spasi? **Masih kena.** Mau dijadiin `a n j i n g` juga tetap kena. Pake titik? Kena. Pake garis miring? Kena. **Pokoknya kena.** 😎

### 🟢 Dictionary — *"The Blacklist"*
Ini kamusnya, lo bisa ngatur sesuka hati.

```go
b.Dict.Contains("anjing")    // true
b.Dict.IsStopWord("dan")     // true (kata "dan" bukan kasar ya, jangan lebay)
b.Dict.AddWords("tai", "keparat")
b.Dict.DelWords("anjing")    // kasian juga si an*jing, maaf ya bang*  
b.Dict.GetDict()             // map[string]struct{}
```

### 🟢 Generator / Iterator
Buat lo yang suka pamer *design pattern* di code review.

```go
res := b.Text("babi anjing").Collect(true).Scan()
gen := res.Generator()
for w := gen.Yield(); w != nil; w = gen.Yield() {
    fmt.Println(w.Word)  // "babi", trus "anjing"
}
```

---

## ⚙️ Cara Kerjanya (Biar Value Lo Naik Di Mata Atasan)

### Pipeline Sederhana
```
Input String → ToLower → LeetSpeak → UnstackChar → [TrimSpace] → Scan / RecursiveScan
```

Mikir pake pipeline ribet? **Enggak.** Kami bukan Netflix. Kami cuma filter kata kasar. *Sederhana, cepat, beres.*

### Sliding Window Affix Detection — *"Ngakalin Imbuhan"*
Bahasa Indonesia itu unik. Kata dasar dikasih imbuhan jadi beda arti. Tapi bukan berarti lolos filter.

```
Input: "mebabi"
Mekanisme sliding window:
├── "mebabi" → cek dictionary? gak ada
├── "mebab" → gak ada
├── "meba" → gak ada
├── "meb" → gak ada
├── "me" → gak ada
├── "ebabi" → coba-coba
│   └── "babi" → 🚨 DITEMUKAN! 🚨
```

**"BABI" NYANGKUT!** Gampang kan?

Bacot gak perlu stemming rumit, gak perlu ML, gak perlu neural network. **Cuma geser-geser doang.** Mirip lo geser Tinder — bedanya kalo ini dapetnya `MATCH ❤️`, bukan ghosting. 

### RecursiveScan — *"No Place to Hide"*

```
Input: "xbabi"
├── xbabi → cek dictionary? Siapa yang punya kata "xbabi"? 
├── babi  → 🚨 KETANGKAP LO! 🚨
├── bab
├── ba
├── b
├── abi
├── ab
├── a
├── bi
├── b
└── i
```

**Mau ditambahin karakter apapun di sela-sela, Bacot tau.** Saking detailnya, mungkin Bacot tau lo lagi makan apa sambil ngetik kata kasar. *We're that good (and annoying).*

---

## 🏎️ Performance

Bacot di benchmark pake Intel i5-13400, Linux, Go 1.26.4:

### Ringan — Chat biasah

| Operasi | Input | Waktu |
|---------|-------|-------|
| `Contains` | dictionary lookup | **5 ns** — lebih cepet dari CPU lo mikir |
| `Scan` | `"aku suka 4njiiing"` | **483 ns** — 2 juta kalimat/detik |
| `RecursiveScan` | `"aku suka 4njiiing"` | **520 ns** — 1,9 juta kalimat/detik |
| `Censor` | hasil scan | **190 ns** — 5,3 juta/detik |

> 💡 Kalo Bacot scan setahun non-stop = 63 triliun chat. CPU lo cuma *"yoi silahkan"*.

### Sedang — Komentar 200-500 char

| Operasi | Input | Waktu |
|---------|-------|-------|
| `ScanParagraph` | 5 kalimat, campur leet+affix | **10 µs** — 50× lebih cepet dari kedip mata |
| `ScanMixedAffixLeet` | 6 kata, affix+leet campur | **1 µs** — 400× lebih cepet dari kedip mata |
| `ScanAffixLong` | 30 kata berimbuhan | **2,9 µs** — 100× lebih cepet dari kedip |
| `CensorParagraph` | hasil scan paragraph | **4,2 µs** — 25 paragraph kesensor dalam 1 kedipan |

### Berat — 10KB s/d 100KB

| Operasi | Input | Waktu |
|---------|-------|-------|
| `Scan10Kb` | 1.000 kata, campur leet+affix | **411 µs** — 500× lebih cepet dari baca judul |
| `Scan10KbAllLeet` | 1.000 kata, 100% leet | **302 µs** — full leet tetap sub-ms |
| `Scan10KbContinuous` | 10KB tanpa spasi | **155 µs** — 6× lebih cepet dari kedip |
| `Scan100Kb` | 100KB lorem ipsum bersih | **21 ms** — 1,2 frame di 60fps |
| `RecursiveScanLongSingleWord` | 1 kata 10K char, match di ujung | **761 µs** — worst case masih sub-ms |
| `ScanAllProfanity` | 10KB, 100% kata kasar | **533 µs** — 4.000 match dalam 0,5 ms |

> **Bottom line:** Bahkan skenario paling ekstrem tetap sub-millisecond (kecuali 100KB yang cuma 21 ms). CPU lo gak bakal ngerasa Bacot ada. *Kecuali kalo lo pake Pentium 4 — itu mah urusan lo.*

<details>
<summary>📊 Raw benchmark data</summary>

```
goos: linux | goarch: amd64 | cpu: 13th Gen Intel(R) Core(TM) i5-13400
```

**RINGAN**

| Benchmark | Iterations | ns/op | B/op | allocs |
|-----------|-----------|-------|------|--------|
| Contains | 43,944,686 | 5.5 | 0 | 0 |
| AddWord | 5,450 | 44,262 | 71,424 | 34 |
| TextPreprocessing | 513,392 | 435.1 | 160 | 7 |
| ScanSingleWord | 1,000,000 | 216.0 | 168 | 6 |
| Scan | 435,322 | 492.7 | 232 | 9 |
| ScanCollect | 300,020 | 679.4 | 312 | 12 |
| RecursiveScan | 431,160 | 534.3 | 232 | 9 |
| Censor | 1,234,800 | 189.0 | 56 | 3 |
| CensorSentence | 1,675,422 | 141.7 | 32 | 1 |

**SEDANG**

| Benchmark | Iterations | ns/op | B/op | allocs |
|-----------|-----------|-------|------|--------|
| ScanParagraph | 25,159 | 9,065 | 3,160 | 32 |
| RecursiveScanParagraph | 22,011 | 12,281 | 3,160 | 32 |
| CensorParagraph | 55,858 | 4,461 | 2,808 | 8 |
| ScanCleanLong | 38,013 | 6,217 | 2,144 | 16 |
| ScanManyStopWords | 49,251 | 4,913 | 2,144 | 16 |
| ScanAffixLong | 86,047 | 2,577 | 1,160 | 16 |
| ScanAllLeet | 41,974 | 5,310 | 3,480 | 63 |
| ScanMixedAffixLeet | 248,612 | 1,028 | 392 | 12 |
| RecursiveScanNoMatch | 92,424 | 2,490 | 352 | 10 |
| CensorLeet | 522,712 | 472.4 | 296 | 5 |

**BERAT**

| Benchmark | Iterations | ns/op | B/op | allocs |
|-----------|-----------|-------|------|--------|
| Scan10Kb | 570 | 404,140 | 244,697 | 2,249 |
| RecursiveScan10Kb | 464 | 443,880 | 244,697 | 2,249 |
| Censor10Kb | 1,476 | 157,493 | 128,505 | 18 |
| Scan10KbAllLeet | 768 | 289,176 | 225,498 | 1,649 |
| TextPreprocessing10Kb | 1,473 | 158,865 | 125,984 | 35 |
| Scan10KbContinuous | 1,473 | 153,782 | 139,865 | 50 |
| RecursiveScan10KbContinuous | 1,533 | 146,431 | 139,905 | 52 |

**SUPER**

| Benchmark | Iterations | ns/op | B/op | allocs |
|-----------|-----------|-------|------|--------|
| Scan100Kb | 13 | 17,101,498 | 8,255,096 | 66 |
| RecursiveScan100Kb | 10 | 21,919,385 | 8,255,088 | 66 |

**EXTREM**

| Benchmark | Iterations | ns/op | B/op | allocs |
|-----------|-----------|-------|------|--------|
| ScanAllProfanity | 451 | 522,762 | 460,507 | 4,055 |
| RecursiveScanAllProfanity | 350 | 649,306 | 460,507 | 4,055 |
| CensorAllProfanity | 1,903 | 127,041 | 122,880 | 2 |
| RecursiveScanLongSingleWord | 316 | 757,932 | 93,320 | 36 |
| ScanLongSingleWord | 2,223 | 115,330 | 93,280 | 34 |
</details>

---

## 📛 Kenapa Namanya "Bacot"?

Kata "bacot" artinya **ngomong ngawur**. Ironis? Ya. Tapi pas — library ini lahir buat nge-filter orang-orang yang lagi bacot di aplikasi lo.

Nama-nama lain yang sempat dipertimbangkan:
- ~~`go-kasar`~~ — *terlalu vulgar, samaan kayak bokap kalo marah*
- ~~`filter-bad-words-indonesia`~~ — *25 karakter? Capek ngetiknya*
- ~~`janganKasarYa`~~ — *kayak ibu-ibu PKK*
- ~~`siPenyaring`~~ — *kedengaran kayak brand detergent*
- ~~`omongKotor`~~ — *ngelindur*
- ~~`kataKasarFilterGo`~~ — *gak kreatif*
- ~~`goblok`~~ — *udah dipake orang (dan bagus juga sih)*

Bacot aja. **Singkat, jelas, sedikit kurang ajar.** Sesuai fungsinya.

---

## ❓ Tanya Jawab Receh (FAQ)

### "Ini library Go pertama saya, gimana cara pakenya?"

`New()` → `Text()` → `Scan()`. Lo bisa Go sekarang. ***Pun intended.***

### "Saya mau pake, tapi takut ada kata kasar yang lolos."

Namanya juga buatan manusia. Kalau nemu, bikin **Issue** atau **Pull Request**. Kami bukan dukun — kami butuh kontribusi lo. *Lo yang mantengin internet setiap hari pasti tau istilah baru.*

### "Kok dictionary-nya cuma 190 kata dasar?"

Karena kami orang baik dan gak mau kesusun kamus sendirian. **PR diterima.** Tambahin kata-kata favorit lo. **Makin kotor makin bagus.** Ini profanity filter, bukan TPA.

### "Saya malu kalo kontribusi soalnya kata-katanya jorok banget."

Bro/sis... Ini **library profanity filter.** Kami tinggal di got. Kami udah biasa baca kata-kata kotor. **Makanya kami buat ini.** Lo kontribusi kata kotor, kami senang. Malu itu buat orang yang nge-filter pake regex doang.

### "Bisa dipake production?"

**Udah.** 74 unit test. 33 benchmark. Zero dep. 5,5 nanosecond lookup. Nggak bakal numpukin memory server lo. *Malah CPU i5 jadi nganggur karena loadingnya terlalu cepet.*

---

## 🙏 Credits

- **Dictionary:** [indonesian-badwords](https://github.com/drizki/indonesian-badwords) — *Riset ulang itu buang waktu. Makasih bang udah dikumpulin.*
- **Stop Words:** [go-sastrawi](https://github.com/truthfulses/go-sastrawi) — *pinjem, makasi.*

Ada kata kasar yang kurang dalam dictionary? **Buka PR.** Tulis kata-kata paling kotor yang lo tau. Kami janji gak akan tersinggung.

*Kami malah senang.* 😄

---

## ⚖️ License

[MIT](https://choosealicense.com/licenses/mit/)
---

