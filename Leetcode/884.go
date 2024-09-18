package Leetcode

import "strings"

func uncommonFromSentences(s1 string, s2 string) []string {
	var ans []string
	wordPresence := make(map[string]int)
	s1Parts := strings.Split(s1, " ")
	s2Parts := strings.Split(s2, " ")

	for _, s1Part := range s1Parts {
		wordPresence[s1Part]++
	}

	for _, s2Part := range s2Parts {
		wordPresence[s2Part]++
	}

	for key, val := range wordPresence {
		if val == 1 {
			ans = append(ans, key)
		}
	}

	return ans
}
