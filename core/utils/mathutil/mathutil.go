package mathutil

import "golang.org/x/exp/constraints"

func Max[V constraints.Ordered](first V, vals ...V) V {
	max := first
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[V constraints.Ordered](first V, vals ...V) V {
	min := first
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}
