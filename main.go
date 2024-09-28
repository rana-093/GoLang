package main

import (
	"GoLang/Backend"
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
	fmt.Println("==================================================")

	err, _ := Backend.ParseXSLXFromObjectHistoryReport("csvparsing/object_history_report_2024_09_01_00_00_00_2024_10_01_00_00_00_1727255721.xlsx")
	if err != nil {
		return
	}
	Backend.ParseXSLXFromDailyUsageReport("csvparsing/driver_daily_distance_report_2024_09_01_00_00_00_2024_09_29_00_00_00_1727541640.xlsx")
}
