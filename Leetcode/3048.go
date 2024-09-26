package Leetcode

import (
	"strconv"
)

func longestCommonPrefix(arr1 []int, arr2 []int) int {
	presence := make(map[string]int)
	ans := 0

	for _, val := range arr1 {
		tempStr := strconv.Itoa(val)
		for i := 0; i < len(tempStr); i++ {
			curString := tempStr[:i+1]
			presence[curString]++
		}
	}

	for _, val := range arr2 {
		tempStr := strconv.Itoa(val)
		for i := 0; i < len(tempStr); i++ {
			curString := tempStr[:i+1]
			if _, exist := presence[curString]; exist {
				ans = max(ans, len(curString))
			}
		}
	}

	return ans
}
