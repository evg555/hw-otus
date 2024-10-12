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

	keys := make([]string, 0, len(mapped))

	for k := range mapped {
		keys = append(keys, k)
	}

	sorted := sortByFreq(keys, mapped)

	if len(sorted) < top {
		top = len(sorted)
	}

	return sorted[:top]
}

func sortByFreq(keys []string, m map[string]int) []string {
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}
