package bacot

import "bacot/dictionary"

func (b *Bacot) WithDefaultDict() *Bacot {
	b.dictionary = dictionary.New()
	return b
}

func (b *Bacot) ClearDict() *Bacot {
	b.dictionary.Clear()
	return b
}

func (b *Bacot) DictContains(word string) bool {
	return b.dictionary.Contains(word)
}
