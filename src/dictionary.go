package bacot

type DictWords = map[string]struct{}

func NewDictWord(words ...string) DictWords {
	new := map[string]struct{}{}
	for _, w := range words {
		new[w] = struct{}{}
	}
	return new
}

type Dictionary struct {
	badWords DictWords
	stops    DictWords
}

func NewDictionary() *Dictionary {

	newDict := &Dictionary{
		badWords: NewDictWord(badwords...),
		stops:    NewDictWord(defaultStopWords...),
	}

	return newDict
}

func (d *Dictionary) AddWords(words ...string) {
	for _, w := range words {
		d.badWords[w] = struct{}{}
	}
}

func (d *Dictionary) DelWords(words ...string) {
	for _, w := range words {
		delete(d.badWords, w)
	}
}

// Remove both stop and badwords
func (d *Dictionary) Clear() *Dictionary {
	d.badWords = DictWords{}
	d.stops = DictWords{}
	return nil
}

func (d *Dictionary) Contains(word string) bool {
	if _, ok := d.badWords[word]; ok {
		return true
	}

	return false
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

func (d *Dictionary) GetDict() DictWords {
	return d.badWords
}

func (d *Dictionary) IsStopWord(s string) bool {
	if _, ok := d.stops[s]; ok {
		return true
	}

	return false
}
