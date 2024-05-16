package mathutil

import (
	"fmt"
	"math"
	"slices"

	"golang.org/x/exp/constraints"
)

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

func Avg[V constraints.Integer](arr ...V) (V, error) {
	total := V(0)

	for _, v := range arr {
		prev := total
		total += v

		// check addition overflow (positive + negative)
		if (total < prev && !math.Signbit(float64(v))) ||
			(total > prev && math.Signbit(float64(v))) {
			return 0, fmt.Errorf("overflow: addition %T", V(0))
		}
	}

	// length overflow
	// assumes array len is always less than MaxInt
	if uint64(V(len(arr))) != uint64(len(arr)) {
		return 0, fmt.Errorf("overflow: array len %d in type %T", len(arr), V(0))
	}

	return total / V(len(arr)), nil
}

// Median mirrors implementation with generics: https://github.com/montanaflynn/stats/blob/249b5aaa10484bb7e8f3b866b0925aaebdac8170/median.go#L6
func Median[V constraints.Integer](arr ...V) (V, error) {
	slices.Sort(arr)

	l := len(arr)
	if l == 0 {
		return 0, fmt.Errorf("empty input")
	}
	if l%2 == 0 {
		return Avg(arr[l/2-1 : l/2+1]...)
	}
	return arr[l/2], nil
}
