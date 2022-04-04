package hw03frequencyanalysis

import (
	"strings"
)

func Top10(input string) []string {
	var processor CountProcessor

	wordsList := checkForNewLineCases(strings.Split(input, " "))

	for i := 0; i < len(wordsList); i++ {
		currentWord := strings.TrimSpace(wordsList[i])

		if currentWord == " " || currentWord == "" {
			continue
		}

		processor.CountWord(currentWord)
	}

	processor.Sort()

	return processor.GetTop10Words()
}

func checkForNewLineCases(wordsList []string) []string {
	for i := 0; i < len(wordsList); i++ {
		checkForNewSpaces := strings.Fields(wordsList[i])
		if len(checkForNewSpaces) > 1 {
			wordsList = remove(wordsList, i)
			wordsList = append(wordsList, checkForNewSpaces...)
		}
	}

	return wordsList
}

func remove(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
