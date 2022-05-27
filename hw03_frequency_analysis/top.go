package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	taskWithAsteriskIsCompleted = true
	rBase                       = regexp.MustCompile(`[\s]+`)
	rWithAsterisk               = regexp.MustCompile(`[^а-я\-]+`)
)

const maxLen = 10

func Top10(text string) []string {
	wordsFrequencyMap := countWordFrequency(text)

	wordFrequencies := make([]wordFrequency, 0, len(wordsFrequencyMap))
	for word, freq := range wordsFrequencyMap {
		wordFrequencies = append(wordFrequencies, wordFrequency{word, freq})
	}

	topFreqWords := sortWords(wordFrequencies)
	if len(topFreqWords) > maxLen {
		topFreqWords = topFreqWords[0:maxLen]
	}

	return topFreqWords
}

func countWordFrequency(text string) map[string]int {
	if len(text) == 0 {
		return nil
	}

	r := rBase
	if taskWithAsteriskIsCompleted {
		text = strings.ToLower(text)
		r = rWithAsterisk
	}
	words := r.Split(text, -1)

	wordsFreq := make(map[string]int, len(words))
	for _, word := range words {
		if taskWithAsteriskIsCompleted && (word == "-" || word == "") {
			continue
		}
		wordsFreq[word]++
	}

	return wordsFreq
}

type wordFrequency struct {
	word string
	freq int
}

func sortWords(frequencies []wordFrequency) []string {
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].freq == frequencies[j].freq {
			return frequencies[i].word < frequencies[j].word
		}

		return frequencies[i].freq > frequencies[j].freq
	})

	words := make([]string, 0, len(frequencies))
	for _, wordFreq := range frequencies {
		words = append(words, wordFreq.word)
	}

	return words
}
