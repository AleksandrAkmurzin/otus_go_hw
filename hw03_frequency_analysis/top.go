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
	wordsFrequency := countWordFrequency(text)

	freq2words := make(map[int][]string, maxLen)
	for word, freq := range wordsFrequency {
		freq2words[freq] = append(freq2words[freq], word)
	}

	frequencies := make([]int, 0, len(freq2words))
	for freq := range freq2words {
		frequencies = append(frequencies, freq)
	}
	sort.Ints(frequencies)

	topFreqWords := make([]string, 0, maxLen)
	for i := len(frequencies) - 1; i >= 0; i-- {
		currentFreqWords := freq2words[frequencies[i]]
		sort.Strings(currentFreqWords)
		topFreqWords = append(topFreqWords, currentFreqWords...)
		if len(topFreqWords) >= maxLen {
			break
		}
	}

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
	for i := 0; i < len(words); i++ {
		word := words[i]
		if taskWithAsteriskIsCompleted && (word == "-" || word == "") {
			continue
		}
		wordsFreq[word]++
	}

	return wordsFreq
}
