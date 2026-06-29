package bacot

import (
	"bacot/dictionary"
	"bacot/internal/core/detector"
)

type Bacot struct {
	dictionary *dictionary.Dictionary
	detector   *detector.Detector
}

func New() *Bacot {
	bacot := &Bacot{dictionary: dictionary.New()}
	bacot.detector = detector.New(bacot.dictionary)
	return bacot

}
