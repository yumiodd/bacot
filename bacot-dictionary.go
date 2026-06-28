package bacot

import "bacot/dictionary"

func (b *Bacot) WithDefaultDict() *Bacot {
	b.dictionary = dictionary.NewDefault()
	return b
}

func (b *Bacot) WithCustomDictionary(dict *dictionary.Dictionary) *Bacot {
	b.dictionary = dict
	return b
}

func (b *Bacot) ClearDict() *Bacot {
	b.dictionary.Clear()
	return b
}

func (b *Bacot) DictContains(word string) bool {
	return b.dictionary.Contains(word)
}
