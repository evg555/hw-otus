package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	top := 10
	words := strings.Fields(str)

	if len(words) == 0 {
		return []string{}
	}

	mapped := make(map[string]int, len(words))

	for _, word := range words {
		mapped[word]++
	}

	sorted := sortByFreq(mapped)

	if len(sorted) < top {
		top = len(sorted)
	}

	return sorted[:top]
}

func sortByFreq(m map[string]int) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}
