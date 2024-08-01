package limits

import "sort"

func max(x int, xs ...int) int {
	sort.Ints(xs)
	if len(xs) == 0 || xs[len(xs)-1] < x {
		return x
	} else {
		return xs[len(xs)-1]
	}
}
