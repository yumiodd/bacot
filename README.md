![test](./bacot.png)

- [English Version](./readme.en.md)

# Bacot
Library Go untuk deteksi kata kotor Bahasa Indonesia. Ringan, cepet, paham imbuhan, zero dependency.


## Instalasi

```go
go get -u github.com/yumiodd/bacot/src
```

## Penggunaan

```go
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {
	b := bacot.New()

	// Cek: kasar apa nggak?
	b.Text("mebabi").Scan().IsProfane()     // true
	b.Text("kelas").Scan().IsProfane()      // false

	// Sensor: biar keliatan sopan
	b.Text("babi dan anjing").Collect(true).Scan().Censor()
	// "**** dan ******"

	// Ekstrak: siapa aja yang kena?
	res := b.Text("asu babi bacot").Collect(true).Scan()
	res.Extract() // ["asu", "babi", "bacot"]
	res.Count()   // 3
}
```
**Catatan: `New()` itu berat** karena precompute semua varian imbuhan. Panggil sekali aja, simpan sebagai singleton, jangan bikin baru tiap request.

## Kenalan

| Masalah | Contoh | Bacot |
|---------|--------|-------|
| Imbuhan | `mebabi` | Sok pake imbuhan? ketauan `babi`-nya |
| Leet speak | `4njing` | Ganti angka percuma, tetep kecium |
| Stack karakter | `anjiiiiing` | Udah di-unstack, tetep `anjing` |
| False positive | `babiru` | `ru` bukan suffix, skip (OK) |
| Suffix valid | `memakan` | stem `makan` + suffix `-an`, valid (OK) |
| Speed | Sensor 10Kb | **172 microsecond** |

## Fitur

### Affix-aware -- senjata utama

Bacot paham imbuhan Indonesia native. Nggak perlu bikin pattern manual.

```go
b.Text("mebabi").Scan().IsProfane()        // true
b.Text("penganjing").Scan().IsProfane()    // true
b.Text("dimakani").Scan().IsProfane()      // true -- di-makan-i
b.Text("mebabi").Affix(false).Scan()...    // false -- exact match aja
```

Prefix yang didukung: `me-`, `pe-`, `di-`, `te-`, `be-`, `ber-`, `ter-`, `per-`, `meng-`, `peng-`, `men-`, `pen-`, `meny-`, `peny-`, `mem-`, `pem-`, `ng-`, `ny-`.

Plus **nasal fusion**: `memukul` = `mem-` + `pukul`. Kamus isi `pukul`, Bacot ngerti sendiri.

### Leet speak -- pura-pura pake angka

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true -- nice try
```

### Unstack char

```go
b.Text("anjiiiiing").Scan().IsProfane() // true
```

### Recursive scan

Mendeteksi kata kotor yang nempel di substring. `xbabi` tetap ketemu `babi`.

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Kamus kustom

```go
b.AddWord("setan")              // varian imbuhan digenerate otomatis
b.Dict.DelWords("anjing")       // hapus dari kamus
b.AddFalsePositive("kelas")     // cegah false match
```

## Konfigurasi

Pipeline bisa diatur sendiri lewat `ModalScanConfig`:

```go
b.Config(&bacot.ModalScanConfig{
	Affix:   false,    // exact match aja
	Collect: true,     // kumpulin semua
	Order: []bacot.SanitizeOrder{
		bacot.WithLeetSpeak,
		bacot.UnstackChar,
	},
})
```

Default pipeline: `Emoji -> ReplaceWhiteSpace -> SanitizeReadSign -> ReplaceWhiteSpace -> UnstackChar -> Affix(true)`

## Cara kerja

1. **Varian imbuhan precomputed** -- pas `New()`, semua kata dasar + varian `me-`, `ber-`, `meng-`, dll udah masuk map. Scan = O(1) lookup. Nggak ada stemming di runtime.

2. **Length histogram pre-filter** -- kalau dictionary cuma punya kata 3-8 karakter, token 11 karakter langsung di-skip tanpa lookup.

3. **False positive filter** -- sisa stripping dicek: kalau suffix dikenal (`-kan`, `-an`, `-i`, `-nya`) dianggap valid. Cuma sisa non-suffix <=1 suku kata yang di-skip (`babiru` -> `ru`).

## Benchmark

| Operasi | Waktu |
|---------|-------|
| Cek kata | ~6 ns |
| Scan kalimat | ~531 ns |
| Sensor 10Kb | ~172 us |

## Lisensi

[MIT](https://choosealicense.com/licenses/mit/)
