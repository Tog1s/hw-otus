package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// var regex regexp.Regexp = *regexp.MustCompile("[а-яА-Я-]+")

type WordCount struct {
	Word  string
	Count int
}

func Top10(inputText string) []string {
	words := make(map[string]int)
	text := string(inputText)
	// matches := regex.FindAllString(text, -1)
	matches := strings.Fields(text)

	for _, match := range matches {
		words[match]++
	}

	var wordCounts []WordCount
	for k, v := range words {
		wordCounts = append(wordCounts, WordCount{k, v})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	})

	var result = make([]string, 0, len(wordCounts))
	for i, wc := range wordCounts {
		result = append(result, wc.Word)
		if i == 9 {
			break
		}
	}
	return result
}
