package Leetcode

type ListNode struct {
	Val  int
	Next *ListNode
}

func spiralMatrix(m int, n int, head *ListNode) [][]int {
	matrix := make([][]int, m)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			matrix[i][j] = -1
		}
	}

	left, right, top, bottom := 0, n-1, 0, m-1

	for head != nil {
		for i := left; i <= right && head != nil; i++ {
			matrix[top][i] = head.Val
			head = head.Next
		}

		for i := top + 1; i <= bottom && head != nil; i++ {
			matrix[i][right] = head.Val
			head = head.Next
		}

		for i := right - 1; i >= left && head != nil; i-- {
			matrix[bottom][i] = head.Val
			head = head.Next
		}

		for i := bottom - 1; i > top && head != nil; i-- {
			matrix[i][left] = head.Val
			head = head.Next
		}

		left++
		right--
		top++
		bottom--
	}

	return matrix
}
