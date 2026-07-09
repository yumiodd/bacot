
![bacot banner](./bacot.png)


# Bacot

Module go sederhana penyaring kata, dikhususkan untuk bahasa Indonesia.

Bahasa Indonesia itu termasuk bahasa aglutinatif, bergantung banget sama affix (imbuhan) yang mana makna satu kata aja bisa berubah-ubah tergantung imbuhannya yang nempel. Capek bikin variasi satu kata dengan imbuhannya? pake module ini, cukup daftarkan kata dasar saja, sisanya beres.

## Installation



```go
import bacot "github.com/yumiodd/bacot/src"
```
    
## Usage/Examples

```go

    // quick use
    b := bacot.New()

    text := b.Text("bacot")

    scanResult := text.Scan() 

    scanResult.IsProfane()      // output: true
    scanResult.Censor()         //output: *****
    scanResult.Extract()        //output: bacot

    // using chain method
    bacot.New().Text("bacot").Scan().IsProfane()    // output: true

    // add new word
    b.Text("memukul").Scan().IsProfane()    // output: false

    b.AddWord("pukul")
    b.Text("memukul").Scan().IsProfane()    // output: true
```
Karena `bacot.New()` akan melakukan generate pada kamus jadi saya tidak menyarankan menggunakannya berulang setiap kali memeriksa text, itu akan mempengaruhi kecepatan proses.

`Bacot.Text("text")` akan mengembalikan struct `ModalScan{}`  yang mana adalah struct utama yang melakuakn scanning pada inputan text, setiap `ModalScan` itu independen jadi configurasi tidak saling mempengaruhi ke yang lain.
```go
    b := b.New()

    // text1 dan text2 adalah 2 ModalScan yang berbeda
    text1 := b.Text("bacot anj ing")
    text2 := b.Text("ba8!")
```
`ModalScan{}` memiliki beberpaa konifgurasi yang bisa terapkan sebelum melakukan `Scan()`.
```go
    b := bacot.New()

    // Chaining Method builder, jadi bisa diset sekaligus secara runtut
    text := b.Text("bacot").WithLeetSpeak().Collect().Affix().TrimSpace().UnstackChar()
    result := text.Scan()

```
> **opsi pre-scan** `Options() // (default value)`
- `TrimSpace() // false`, seperti namanya space input text akan dihilangkan.
- `WithLeetSpeak() // true`, menggati semua karakter leet speak ke bentuk alfabet biasa, cth: @ -> a, 8 -> b, 6 -> g. 
- `UnstackChar() // true`, menghapus karakter duplikat pada input text, mempercepat proses pencarian dan lebih akurat.
- `Affix(bool) // true`, pada dasarnya daftar kata pada kamus sebagian besar hanya terdiri dari kata dasar, kemudian semua kata yang berpotensi luluh karena penggabungan imbuhan men- dan peng- digenerate ketika `bacot.New()`. Dengan opsi ini kata-kata bentuk kompleks yang memiliki potensi terdaftar dalam kamus seperti "berpembacotan" tidak akan lolos. *note: hanya bekerja untuk `Scan()`.
- `Collect(true) // false`, defaultnya `Scan()`  akan selesai ketika mendapatkan satu kata atau tidak sama sekali, dengan opsi ini semua kata yang ditemukan akan dikumpulkan hingga akhir dari text input.
- method `GetText()`, mengembalikan bentuk akhir string dari text yang akan discan, cth: `Text("ba co t").TrimSpace().GetText() // output: "bacot"`

> Terdapat dua jenis scan,
- `Scan()`, mencocokan perkata.
- `RecursiveScan()`, mencocokan untuk setiap substring dari kata, tidak termasuk spasi.

###### **Scan Return**
semua jenis scan akan mengambalikan struct `ScanResult{}`:
``` go
    
    res := bacot.New().Text("jangan banyak bacot kamu ya, sialan").Scan()

    res.Extract()   // output: ["bacot", "sialan"]
    res.Censor()    // output: "jangan banyak ***** kamu ya, ******"
    res.GetText()   // output: "jangan banyak bacot kamu ya, sialan"
    res.IsProfane() // output: true
    res.First()     // output: "bacot"
    res.Last()      // output: "sialan"
    res.Count()     // output: 2
```

## Contributing

Segala konstribusi sangat wellcome ya, terutama untuk daftar kamus.

## License

[MIT](https://choosealicense.com/licenses/mit/)

