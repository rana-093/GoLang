package Leetcode

import "strings"

func countConsistentStrings(allowed string, words []string) int {
	ans := 0
	for i := 0; i < len(words); i++ {
		for _, ch := range words[i] {
			res := strings.ContainsRune(allowed, ch)
			if !res {
				ans--
				break
			}
		}
		ans++
	}
	return ans
}
