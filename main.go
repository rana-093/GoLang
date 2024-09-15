package main

import (
	"GoLang/RandomConcepts"
	"fmt"
	"sort"
)

func twoSum(num []int, target int) []int {
	sort.Ints(num)
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		curSum := num[i] + num[j]
		fmt.Printf("cursum : %d, i = %d, j = %d", curSum, i, target)
		if curSum > target {
			j = j - 1
		} else if curSum < target {
			i = i + 1
		} else {
			return []int{i, j}
		}
		fmt.Printf("cursum : %d, i = %d, j = %d", curSum, i, j)
	}
	return []int{-1, -1}
}

func main() {
	//checkGoRoutinesAndChannels()
	calc()
	err := RandomConcepts.TestInterface()
	fmt.Println(err)
}
