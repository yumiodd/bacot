# Bacot

Profanity filter for Indonesian language. Lightweight, fast, zero dependency, and **affix-aware**.

## Installation

```go
import bacot "github.com/yumiodd/bacot/src"
```

## Usage

```go
b := bacot.New()

// Detection
b.Text("4njiiing").WithLeetSpeak().Scan().IsProfane() // true

// Censor
b.Text("babi dan anjing").Collect(true).Scan().Censor() // "**** dan ******"

// Extract
res := b.Text("babi dan anjing asu kontol").Collect(true).Scan()
res.Count()   // 4
res.Extract() // ["babi", "anjing", "asu", "kontol"]
res.First()   // "babi"
```

## Why Bacot?

| Problem | Example | Other libs | Bacot |
|---------|---------|-----------|-------|
| Affix (prefix/suffix) | `mebabi`, `penganjing` | manual pattern per affix | **automatic** |
| Leet speak | `4njing`, `k0nt0l` | endless regex | built-in |
| Repeated chars | `anjiiiiing` | per-word regex | automatic |
| False positive | `babiru` ≠ `babi` | manual exceptions | handled |
| Speed | 10Kb censor | 5-20ms | **172µs** |

## Features

### Affix-aware — the main advantage

Bacot understands Indonesian affixes natively. `mebabi` → `me-` + `babi`. `penganjing` → `peng-` + `anjing`. No manual patterns needed.

Supported affixes: `me-`, `pe-`, `di-`, `te-`, `be-`, `ber-`, `ter-`, `per-`, `meng-`, `peng-`, `men-`, `pen-`, `meny-`, `peny-`, `mem-`, `pem-`, `ng-`, `ny-`.

Plus **nasal fusion**: `memukul` = `mem-` + `pukul` (p dropped). Dictionary only needs `pukul`.

```go
b.Text("mebabi").Scan().IsProfane()      // true
b.Text("mebabi").Affix(false).Scan()...  // false — exact match only
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

Detects substrings inside tokens. `"xbabi"` → scan every position → finds `"babi"`.

```go
b.Text("xbabi").RecursiveScan().IsProfane() // true
```

### Custom dictionary

```go
b.AddWord("setan")     // auto-generates all affix variants
b.Dict.DelWords("anjing")
b.AddFalsePositive("kelas") // prevent false match on "kelas"
```

## Benchmark

| Operation | Time |
|-----------|------|
| Check word | ~6 ns |
| Scan sentence | ~531 ns |
| Censor 10Kb | ~172 µs |

## How it works

1. **Precomputed affix variants** — all affixed forms generated at `New()`. Scan = exact map lookup. No runtime stemming.
2. **Length histogram** — pre-filter: skip tokens with no matching length in dictionary. O(log n).
3. **False positive filter** — if stem remainder ≤1 syllable, it's likely not a real affix. `"babiru"` → stem `"babi"` → remainder `"ru"` → skip.

## License

[MIT](https://choosealicense.com/licenses/mit/)
