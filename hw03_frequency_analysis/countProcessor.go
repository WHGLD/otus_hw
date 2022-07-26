package hw03frequencyanalysis

type CountProcessor struct {
	list []WordCount
	n    int
}

func (p *CountProcessor) CountWord(word string) {
	if index, ok := p.getListIndex(word); ok == false {
		wordCount := MakeWordCount(word)
		p.list = append(p.list, wordCount)
	} else {
		p.list[index].Count()
	}
}

func (p *CountProcessor) getListIndex(word string) (int, bool) {
	for i := 0; i < len(p.list); i++ {
		if p.list[i].CheckWord(word) {
			return i, true
		}
	}

	return -1, false
}

func (p *CountProcessor) Sort() {
	p.list = RunSort(p.list)
}

func (p *CountProcessor) GetTop10Words() []string {
	var onlyWords []string
	for index, wordCount := range p.list {
		if index < 10 {
			onlyWords = append(onlyWords, wordCount.word)
		} else {
			break
		}
	}

	return onlyWords
}
