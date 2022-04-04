package hw03frequencyanalysis

type WordCount struct {
	word string
	n    int
}

func (w *WordCount) Count() {
	w.n++
}

func (w *WordCount) CheckWord(word string) bool {
	return word == w.word
}

func MakeWordCount(word string) WordCount {
	return WordCount{word, 1}
}
