package Leetcode

func xorQueries(arr []int, queries [][]int) []int {
	left := make([]int, len(arr))
	left[0] = arr[0]
	for i := 1; i < len(arr); i++ {
		left[i] = left[i-1] ^ arr[i]
	}
	ans := make([]int, len(queries))
	for idx, query := range queries {
		l, r := query[0], query[1]
		leftXor, rightXor := 0, left[r]
		if l-1 < 0 {
			leftXor = 0
		} else {
			leftXor = left[l-1]
		}
		ans[idx] = leftXor ^ rightXor
	}
	return ans
}
