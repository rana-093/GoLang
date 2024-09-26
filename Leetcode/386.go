package Leetcode

func solve1(n, limit int) []int {
	if n > limit {
		return []int{}
	}
	var cur []int
	if n != 0 {
		cur = append(cur, n)
	}
	for i := 0; i < 10; i++ {
		temp := n*10 + i
		if temp == 0 || temp > limit {
			continue
		}
		xyz := solve1(temp, limit)
		cur = append(cur, xyz...)
	}
	return cur
}

func lexicalOrder(n int) []int {
	return solve1(0, n)
}
