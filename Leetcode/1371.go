package Leetcode

func isVowel(ch rune) bool {
	return ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u'
}

func vowelPosition(ch rune) int {
	switch ch {
	case 'a':
		return 0
	case 'e':
		return 1
	case 'i':
		return 2
	case 'o':
		return 3
	case 'u':
		return 4
	}
	return -1
}

func findTheLongestSubstring(s string) int {
	freq := make(map[int]int)
	mask, maxLen := 0, 0
	freq[0] = -1
	for idx, val := range s {
		if isVowel(val) {
			pos := vowelPosition(val)
			mask ^= 1 << pos
		}

		if _, exists := freq[mask]; exists {
			maxLen = max(maxLen, idx-freq[mask])
		}

		if _, exists := freq[mask]; !exists {
			freq[mask] = idx
		}
	}
	return maxLen
}
