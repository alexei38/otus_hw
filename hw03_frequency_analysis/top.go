package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reWord = regexp.MustCompile("([А-Яа-яA-Za-z0-9]+([А-Яа-яA-Za-z0-9-][А-Яа-яA-Za-z0-9]+)?)+")

type Freq struct {
	Word  string
	Count int
}

type FreqList []Freq

func (f FreqList) Top(top int) []string {
	var result []string
	if len(f) > 0 {
		for i := 0; i < top && len(f) > i; i++ {
			result = append(result, f[i].Word)
		}
	}
	return result
}

func (f FreqList) Sort() {
	if len(f) > 0 {
		sort.Slice(f, func(i, j int) bool {
			if f[i].Count == f[j].Count {
				return f[i].Word < f[j].Word
			}
			return f[i].Count > f[j].Count
		})
	}
}

func createFreqStruct(freqs map[string]int) FreqList {
	freqList := FreqList{}
	for word, cnt := range freqs {
		freqList = append(freqList, Freq{word, cnt})
	}
	return freqList
}

func countWords(s string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(s) {
		word = strings.ToLower(word)
		match := reWord.FindString(word)
		if match != "" {
			counts[match]++
		}
	}
	return counts
}

func Top(text string, top int) []string {
	counts := countWords(text)
	freqStruct := createFreqStruct(counts)
	freqStruct.Sort()
	return freqStruct.Top(top)
}

func Top10(text string) []string {
	return Top(text, 10)
}
