package Leetcode

func isConsistent(mp map[rune]bool, str string) bool {
	for _, ch := range str {
		if !mp[ch] {
			return false
		}
	}
	return true
}

func countConsistentStrings(allowed string, words []string) int {
	presence := make(map[rune]bool)
	for _, ch := range allowed {
		presence[ch] = true
	}
	ans := 0
	for i := 0; i < len(words); i++ {
		if isConsistent(presence, words[i]) {
			ans++
		}
	}
	return ans
}
