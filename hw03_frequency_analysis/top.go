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

	sort.Slice(keys, func(i, j int) bool {
		if mapped[keys[i]] == mapped[keys[j]] {
			return keys[i] < keys[j]
		}
		return mapped[keys[i]] > mapped[keys[j]]
	})

	if len(keys) < top {
		top = len(keys)
	}

	return keys[:top]
}
