package Leetcode

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

func FindMinDifference(timePoints []string) int {
	var arr []int

	for _, val := range timePoints {
		parts := strings.Split(val, ":")
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		totalMinutes := 60*hour + minute
		arr = append(arr, totalMinutes)
		totalMinutesIn24HrFormat := 60*(hour+24) + minute
		arr = append(arr, totalMinutesIn24HrFormat)
	}

	sort.Ints(arr)
	ans := math.MaxInt

	for i := 0; i+1 < len(arr); i++ {
		diff := arr[i+1] - arr[i]
		ans = min(ans, diff)
	}

	return ans
}
