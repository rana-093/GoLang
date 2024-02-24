package main

import (
	"fmt"
	"sort"
)

func twoSum(num []int, target int) []int {
	sort.Ints(num)
	for i, j := 0, len(num)-1; j > i; {
		curSum := num[i] + num[j]
		if curSum > target {
			j = j - 1
		} else if curSum < target {
			i = i + 1
		} else {
			return []int{i, j}
		}
	}
	return []int{-1, -1}
}

func main() {
	nums, target := []int{2, 7, 11, 15}, 9
	result := twoSum(nums, target)
	fmt.Printf("result is: %d\n", result)
}
