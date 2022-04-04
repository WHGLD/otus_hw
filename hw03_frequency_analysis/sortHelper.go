package hw03frequencyanalysis

import (
	"sort"
)

func RunSort(wordCountsList []WordCount) []WordCount {
	// разбиваем на группы с одинаковой частотностью
	m := map[int][]WordCount{}
	for _, wordCount := range wordCountsList {
		m[wordCount.n] = append(m[wordCount.n], wordCount)
	}

	// сортируем группы по частотности
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// сортируем по лексике слова внутри групп
	var sortedWordsCount []WordCount
	for _, k := range keys {
		groupedWordsCount := m[k]
		sort.Slice(groupedWordsCount, func(i, j int) bool {
			return groupedWordsCount[i].word < groupedWordsCount[j].word
		})
		// собираем итоговый результат
		sortedWordsCount = append(sortedWordsCount, groupedWordsCount...)
	}

	return sortedWordsCount
}
