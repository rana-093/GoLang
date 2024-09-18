package Leetcode

import (
	"sort"
	"strconv"
)

func largestNumber(nums []int) string {
	sort.Slice(nums, func(i, j int) bool {
		strxy := strconv.Itoa(nums[i]) + strconv.Itoa(nums[j])
		stryx := strconv.Itoa(nums[j]) + strconv.Itoa(nums[i])
		return strxy > stryx
	})
	ans := ""
	for _, val := range nums {
		ans += strconv.Itoa(val)
	}
	if ans[0] == '0' {
		return "0"
	}
	return ans
}
