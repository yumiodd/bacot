![test](bacot.png)
# Bacot

Profanity filter for Indonesian language. Lightweight, fast, affix-aware, zero dependency.

## Installation

```go
go get -u github.com/yumiodd/bacot/src
```

## Usage

```go
package main

import bacot "github.com/yumiodd/bacot/src"

func main() {
	b := bacot.New()

	// Check: is it profane?
	b.Text("mebabi").Scan().IsProfane()     // true
	b.Text("kelas").Scan().IsProfane()      // false

	// Censor: make it presentable
	b.Text("babi dan anjing").Collect(true).Scan().Censor()
	// "**** dan ******"

	// Extract: who's being toxic?
	res := b.Text("asu babi bacot").Collect(true).Scan()
	res.Extract() // ["asu", "babi", "bacot"]
	res.Count()   // 3
}
```
**Note: `New()` is heavy** because it precomputes all affix variants. Call it once and keep it as a singleton -- don't create a new instance per request.


## Quick overview

| Problem | Example | Bacot says |
|---------|---------|------------|
| Affix | `mebabi` | Nice try, still `babi` underneath |
| Leet speak | `4njing` | Numbers won't save you |
| Repeated chars | `anjiiiiing` | De-stacked, still `anjing` |
| False positive | `babiru` | `ru` isn't a suffix, skip (OK) |
| Valid suffix | `memakan` | stem preffix `me`+`makan`, clean, skip (OK) |
| Speed | 10Kb censor | **172 microseconds** |

## Features

### Affix-aware -- the main weapon

Bacot understands Indonesian affixes natively. No manual regex patterns needed.

```go
b.Text("mebabi").Scan().IsProfane()        // true
b.Text("penganjing").Scan().IsProfane()    // true
b.Text("dimakani").Scan().IsProfane()      // false -- di-makan-i
b.Text("mebabi").Affix(false).Scan()...    // false -- exact match only
```

Supported prefixes: `me-`, `pe-`, `di-`, `te-`, `be-`, `ber-`, `ter-`, `per-`, `meng-`, `peng-`, `men-`, `pen-`, `meny-`, `peny-`, `mem-`, `pem-`, `ng-`, `ny-`.

Plus **nasal fusion**: `memukul` = `mem-` + `pukul`. Dictionary only needs `pukul`.

### Leet speak

```go
b.Text("4njing").WithLeetSpeak().Scan().IsProfane() // true
```

### Unstack char

```go
b.Text("anjiiiiing").Scan().IsProfane() // true
```

### Recursive scan

Detects profanities embedded inside substrings. `xbabi` still catches `babi`.

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Custom dictionary

```go
b.AddWord("setan")              // auto-generates all affix variants
b.Dict.DelWords("anjing")       // remove from dictionary
b.AddFalsePositive("kelas")     // prevent false match
```

## Configuration

Customize the preprocessing pipeline via `ModalScanConfig`:

```go
b.Config(&bacot.ModalScanConfig{
	Affix:   false,    // exact match only
	Collect: true,     // collect all matches
	Order: []bacot.SanitizeOrder{
		bacot.WithLeetSpeak,
		bacot.UnstackChar,
	},
})
```

Default pipeline: `Emoji -> ReplaceWhiteSpace -> SanitizeReadSign -> ReplaceWhiteSpace -> UnstackChar -> Affix(true)`

## How it works

1. **Precomputed affix variants** -- at `New()`, all root words plus their `me-`, `ber-`, `meng-` variants are generated into a map. Scan = O(1) map lookup. No runtime stemming overhead.

2. **Length histogram pre-filter** -- if the dictionary only has words of length 3-8, an 11-character token is skipped immediately without a lookup.

3. **False positive filter** -- after prefix stripping, the remainder is checked: known suffixes (`-kan`, `-an`, `-i`, `-nya`) are considered valid. Only non-suffix remainders <=1 syllable are skipped (`babiru` -> `ru`).

## Benchmark

| Operation | Time |
|-----------|------|
| Check word | ~6 ns |
| Scan sentence | ~531 ns |
| Censor 10Kb | ~172 us |

## License

[MIT](https://choosealicense.com/licenses/mit/)
