package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
}

func Top10(inputText string) []string {
	// if taskWithAsteriskIsCompleted = true
	// text := strings.ToLower(inputText)
	// text = strings.ReplaceAll(text, " - ", "")
	// text = strings.ReplaceAll(text, ".", "")

	text := inputText
	allWords := strings.Fields(text)

	words := make(map[string]int)
	for _, match := range allWords {
		words[match]++
	}

	wordCounts := make([]WordCount, 0, len(words))
	for k, v := range words {
		wordCounts = append(wordCounts, WordCount{k, v})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			return wordCounts[i].Word < wordCounts[j].Word
		}
		return wordCounts[i].Count > wordCounts[j].Count
	})

	result := make([]string, 0, len(wordCounts))
	for i, wc := range wordCounts {
		result = append(result, wc.Word)
		if i == 9 {
			break
		}
	}
	return result
}
