package Leetcode

func doesExist(length int, s string) (bool, int) {
	freq1, freq2 := make(map[int]int), make(map[int]int)
	const mod1, mod2 = 1e9 + 7, 2147483629
	base1, hash1, power1, hash2, power2, base2 := 31, 0, 1, 0, 1, 131

	for i := 0; i < length-1; i++ {
		power1 = (power1 * base1) % mod1
		power2 = (power2 * base2) % mod2
	}

	for i := 0; i < length; i++ {
		hash1 = (hash1*base1 + int(s[i])) % mod1
		hash2 = (hash2*base2 + int(s[i])) % mod2
	}

	freq1[hash1%mod1]++
	freq2[hash2%mod2]++

	for i := length; i < len(s); i++ {
		hash1 = (hash1 - (int(s[i-length])*power1 + mod1)) % mod1
		hash1 += mod1
		hash1 = (hash1*base1 + (int(s[i]))) % mod1
		hash1 %= mod1
		freq1[hash1]++

		hash2 = (hash2 - (int(s[i-length])*power2 + mod2)) % mod2
		hash2 += mod2
		hash2 = (hash2*base2 + (int(s[i]))) % mod2
		hash2 %= mod2
		freq2[hash2]++

		if freq1[hash1] > 1 && freq2[hash2] > 1 {
			return true, i - length + 1
		}
	}
	return false, -1
}

// LongestDupSubstring banana, -> ana, ana ||  abed -> ""
func longestDupSubstring(s string) string {
	lo, hi := 0, len(s)-1
	st, ans := -1, -1
	for lo <= hi {
		mid := (lo + hi + 1) / 2
		if exists, start := doesExist(mid, s); exists {
			st = start
			ans = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	if ans == -1 {
		return ""
	}
	return s[st : st+ans]
}
