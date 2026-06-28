package dictionary

type DictWords = map[string]struct{}

func NewDictWord(words ...string) DictWords {
	new := make(map[string]struct{}, 10000)
	for _, w := range words {
		new[w] = struct{}{}
	}
	return new
}

var (
	ScanBadWords = 1
	ScanWords    = 2
	ScanBoth     = 3
)

type Dictionary struct {
	scanLevel int

	words    DictWords
	badWords DictWords

	stops DictWords
}

func New(words DictWords, stops ...string) *Dictionary {

	newStops := DictWords{}
	if len(stops) > 0 {
		newStops = NewDictWord(stops...)
	}

	newDict := &Dictionary{
		words: words,
		stops: newStops,
	}

	return newDict
}

func NewDefault() *Dictionary {
	return &Dictionary{
		words:    NewDictWord(defaultWords...),
		badWords: NewDictWord(badwords...),
		stops:    NewDictWord(defaultStopWords...),
	}
}

func (d *Dictionary) SetScanLevel(level int) {
	d.scanLevel = level
}

func (d *Dictionary) AddWords(words ...string) {
	for _, w := range words {
		d.words[w] = struct{}{}
	}
}

func (d *Dictionary) DelWords(words ...string) {
	for _, w := range words {
		delete(d.words, w)
	}
}

func (d *Dictionary) Clear() *Dictionary {
	d.words = DictWords{}
	d.stops = DictWords{}
	return nil
}

func (d *Dictionary) Contains(word string) bool {

	found := false
	if (d.scanLevel == ScanBadWords) || d.scanLevel == ScanBoth {
		_, ok := d.badWords[word]
		found = found || ok
	}
	if (d.scanLevel == ScanBadWords) || d.scanLevel == ScanBoth {
		_, ok := d.words[word]
		found = found || ok
	}

	return found
}

func (d *Dictionary) AddStopWords(words ...string) {
	for _, w := range words {
		d.stops[w] = struct{}{}
	}
}

func (d *Dictionary) DelStopWords(words ...string) {
	for _, w := range words {
		delete(d.stops, w)
	}
}

func (d *Dictionary) IsBadword(word string) bool {
	_, ok := d.badWords[word]
	return ok
}

func (d *Dictionary) AddBadWords(words ...string) {
	for _, w := range words {
		d.stops[w] = struct{}{}
		d.words[w] = struct{}{}
	}
}
