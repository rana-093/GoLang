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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func averageWaitingTime(customers [][]int) float64 {
	numRows := len(customers)
	if numRows == 1 {
		return float64(customers[0][1])
	}
	currentTime, waiting := customers[0][0]+customers[0][1], customers[0][1]
	for i := 1; i < numRows; i++ {
		if customers[i][0] >= currentTime+customers[i][1] {
			currentTime = customers[i][0] + customers[i][1]
		} else {
			currentTime = currentTime + customers[i][1]
		}
		waiting += max(0, currentTime-customers[i][0])
	}
	return float64(waiting) / float64(numRows)
}

func main() {
	// nums, target := []int{2, 7, 11, 15}, 9
	// result := twoSum(nums, target)
	// fmt.Printf("result is: %d\n", result)

	customers := [][]int{{5, 2}, {5, 4}, {10, 3}, {20, 1}}
	ans := averageWaitingTime(customers)
	fmt.Println("Result is %.2f ", ans)
}
