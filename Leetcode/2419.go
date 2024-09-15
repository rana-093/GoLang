package Leetcode

import (
	"math"
)

func longestSubarray(nums []int) int {
	highest := math.MinInt32
	for _, val := range nums {
		highest = max(highest, val)
	}
	length, maxLen := 0, 1
	for i := 0; i < len(nums); i++ {
		if nums[i] == highest {
			length++
		} else {
			length = 0
		}
		maxLen = max(maxLen, length)
	}
	return maxLen
}
