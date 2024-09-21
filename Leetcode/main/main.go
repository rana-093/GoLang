package main

import (
	"fmt"
)

func main() {
	ans := []int{1, 2, 3, 4}
	for _, val := range ans {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}
