package Leetcode

import "strconv"

var memo map[string][]int

func solve(str string) []int {
	if res, exists := memo[str]; exists {
		return res
	}
	var result []int
	for i := 0; i < len(str); i++ {
		if str[i] == '+' || str[i] == '-' || str[i] == '*' {
			left := solve(str[:i])
			right := solve(str[i+1:])
			for _, l := range left {
				for _, r := range right {
					switch str[i] {
					case '+':
						result = append(result, l+r)
					case '-':
						result = append(result, l-r)
					case '*':
						result = append(result, l*r)
					}
				}
			}
		}
	}
	if len(result) == 0 {
		num, _ := strconv.Atoi(str)
		result = append(result, num)
	}
	return result
}

func diffWaysToCompute(expression string) []int {
	memo = make(map[string][]int)
	return solve(expression)
}
