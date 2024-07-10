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

// LC - 1525
func numSplits(s string) int {
	m1, m2, N, sz := make([]bool, 26), make([]bool, 26), len(s), len(s)
	left, right, ans := make([]int, N), make([]int, N+1), 0
	left[0], right[N] = 1, 0
	m1[rune(s[0])-'a'] = true

	for i, j := 0, sz-1; i < sz; i, j = i+1, j-1 {
		if i > 0 {
			left[i] = left[i-1]
		}
		right[j] = right[j+1]

		if !m1[rune(s[i])-'a'] {
			m1[rune(s[i])-'a'] = true
			left[i]++
		}
		if !m2[rune(s[j])-'a'] {
			m2[rune(s[j])-'a'] = true
			right[j]++
		}
	}

	for i, _ := range s {
		uniqueLeft, uniqueRight := left[i], right[i+1]
		if uniqueLeft == uniqueRight {
			ans++
		}
	}

	return ans
}

func main() {
	// nums, target := []int{2, 7, 11, 15}, 9
	// result := twoSum(nums, target)
	// fmt.Printf("result is: %d\n", result)

	customers := [][]int{{5, 2}, {5, 4}, {10, 3}, {20, 1}}
	ans := averageWaitingTime(customers)
	fmt.Println("Result is %.2f ", ans)
}
