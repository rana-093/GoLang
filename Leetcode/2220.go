package Leetcode

func minBitFlips(start int, goal int) int {
	xor, ans := start^goal, 0
	for i := 0; i < 32; i++ {
		if xor>>uint(i)&1 == 1 {
			ans++
		}
	}
	return ans
}
