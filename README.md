[English Version](./readme.en.md) 

# Bacot

Filter kata kasar / profanity khusus Bahasa Indonesia. Ringan, cepat, zero dependency, dan **affix-aware**.

## Instalasi

```go
import bacot "github.com/yumiodd/bacot/src"
```

## Penggunaan

```go
b := bacot.New()

// Deteksi
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true

// Sensor
b.Text("babi dan anjing").Collect(true).Scan().Censor() // "**** dan ******"

// Ekstrak
res := b.Text("babi dan anjing asu kontol").Collect(true).Scan()
res.Count()   // 4
res.Extract() // ["babi", "anjing", "asu", "kontol"]
res.First()   // "babi"
```

## Kenapa Bacot?

| Masalah | Contoh | Library lain | Bacot |
|---------|--------|-------------|-------|
| Imbuhan | `mebabi`, `penganjing` | bikin pattern tiap imbuhan | **otomatis** |
| Leet speak | `4njing`, `k0nt0l` | regex tak berujung | built-in |
| Karakter berulang | `anjiiiiing` | per-kata regex | otomatis |
| False positive | `babiru` bukan `babi` | manual exception | sudah di-handle |
| Kecepatan | Sensor 10Kb | 5-20ms | **172µs** |

## Fitur

### Affix-aware — keunggulan utama

Bacot paham imbuhan Indonesia secara native. `mebabi` → `me-` + `babi`. `penganjing` → `peng-` + `anjing`. Tanpa pattern manual.

Imbuhan yang didukung: `me-`, `pe-`, `di-`, `te-`, `be-`, `ber-`, `ter-`, `per-`, `meng-`, `peng-`, `men-`, `pen-`, `meny-`, `peny-`, `mem-`, `pem-`, `ng-`, `ny-`.

Plus **nasal fusion**: `memukul` = `mem-` + `pukul` (p luluh). Kamus cukup berisi `pukul`.

```go
b.Text("mebabi").Scan().IsProfane()      // true
b.Text("mebabi").Affix(false).Scan()...  // false — exact match saja
```

### Leet speak

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true
```

### Unstack char

```go
b.Text("anjiiiiing").Scan().IsProfane() // true
```

### Recursive scan

Mendeteksi substring di dalam token. `"xbabi"` → scan tiap posisi → ketemu `"babi"`.

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Kamus kustom

```go
b.AddWord("setan")     // otomatis generate semua varian imbuhan
b.Dict.DelWords("anjing")
b.AddFalsePositive("kelas") // cegah false match pada "kelas"
```

## Benchmark

| Operasi | Waktu |
|---------|-------|
| Cek kata | ~6 ns |
| Scan kalimat | ~531 ns |
| Sensor 10Kb | ~172 µs |

## Cara kerja

1. **Varian imbuhan di-generate di awal** — semua bentuk berimbuhan dibuat saat `New()`. Scan = pencarian exact di map. Tanpa stemming di runtime.
2. **Histogram panjang kata** — pre-filter: skip token yang panjangnya tidak ada di kamus. O(log n).
3. **Filter false positive** — jika sisa stem ≤1 suku kata, kemungkinan bukan imbuhan beneran. `"babiru"` → stem `"babi"` → sisa `"ru"` → skip.

## Lisensi

[MIT](https://choosealicense.com/licenses/mit/)
