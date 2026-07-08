
![bacot.png](./bacot.png)
-
# **Bacot**
Is a simple word filter module specifically for Indonesian language content. Having trouble making a list of words that are change its form because of KPST rules? Use this.
## Basic Use

Basically this is just a word filter module, but with a little algorithm for Indonesian.

```
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {

    // Membuat struct baru
    b := bacot.New()

    // Membuat modal scan
    text1 := b.Text("dasar kamu Anjing")
    text2 := b.Text("babi anjing")

    // melakukan scan, menghasilkan struct result
    //
    // kita menyediakan 2 opsi:
    // Scan(), melakuakan pemeriksaan per kata yang dipisah dengan spasi.
    // RecursiveScan(), proses ini akan memeriksa substring dari input.

    res1 := text1.Scan()
    res2 := text2.RecursiveScan()

    // Mendapatkan informasi hasil
    fmt.Println(res1.IsToxic())         // output: true , profined->"Anjing"
    fmt.Println(res2.CountFoundWord())  // output: 2, profined-> "babi", "anjing"
}

```
## Option pre-scan
Beberapa opsi yang bisa kamu gunakan untuk mengubah prilaku scan:
```
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {

    // initial
    b := bacot.New()

    text := b.Text("dasar kamu babi")

    // Menghilangkan whitespace
    text = text.WithSanitizeWhiteSpace(true)

    // Menormalisasikan semua leetSpeak karakter
    text = text.WithNormalizeLeetSpeak(true)

    // Mengumpulkan semua kata yang terdapat dalam kamus
    text = text.Collect()

    // Menggunakan Stemming, digunakan sebagai default
    text = text.WithStemming()

}

```
Semua opsi diatas bisa digunakan untuk semua jenis scan, `Scan() ` dan `RecursiveScan()`. 

Untuk whitespace kami memiliki daftar default, namun bisa kamu custom di saat pembuatan New(), tapi masih dalam tahap pengembangan.

## Hasil Penyaringan
hasil penyaringan, dapat di proses dalam berbagai rupa, di dalam module sudah di sediakan dan akan di tambah seiring kebutuhan dan dalam pengambangan.
```
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {


    res := bacot.New().Text("dasar babi anjing").Scan()

    // Memerika apakah mengandung kata yang dibatasi
    fmt.Println(res.IsToxic())                  // output: true

    // Mengambil kata pertama dan terakhir jumpa
    fmt.Println(res.First())                    // output: "babi"
    fmt.Println(res.Last())                     // output: "anjig"

    // Menyensor kata yang dibatasi dalam kalimat
    fmt.Println(res.CensorText())               // output: "dasar **** ******"

    // Menghitung kata yang dibatasi
    fmt.Println(res.CountFoundWord())           // output: 2

    // Mengekstrak kata yang dibatasi yang ditemukan
    words := res.Extract()                      // output: []string{"babi", "anjing"}

    // Mengabil data mentah
    wordIndexes := res.GetFountWords()          // output: []*WordIndex{}

    // Membuat word generator, mirip fungsi yield di python
    gen := res.Generator()

    fmt.Println(gen.Yield().Word())              // output : "babi"   
    fmt.Println(gen.Yield().Word())              // output : "anjing"   
}

``` 
## Pengembangan
Seperti module pada umumnya, module ini open source saya harap para module dev yang tertarik dapat berkonstribusi dan membuat module ini lebih baik, bahkan untuk yang tidak memahami coding masih bisa melakukan konstribusi dengan melakukan update pada kamus perkataan kasar, saya harap module ini berguna dan bisa dipakai oleh orang-orang yang memiliki kesulitan yang serupa.