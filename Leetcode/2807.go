package Leetcode

func gcd(x, y int) int {
	if y == 0 {
		return x
	}
	return gcd(y, x%y)
}

func insertGreatestCommonDivisors(head *ListNode) *ListNode {
	var curr, savedHead = head, head
	for curr != nil {
		curVal := curr.Val
		nextVal := -1
		if curr.Next != nil {
			nextVal = curr.Next.Val
		}
		if nextVal != -1 {
			value := gcd(curVal, nextVal)
			nextNode := curr.Next
			temp := &ListNode{Val: value, Next: nil}
			curr.Next = temp
			temp.Next = nextNode
			curr = temp.Next
		} else {
			curr = curr.Next
		}
	}

	return savedHead
}
