package utils

import (
	"fmt"
)

// BatchSplit splits an slices into an slices of slicess with a maximum length
func BatchSplit[T any](list []T, max int) (out [][]T, err error) {
	if max == 0 {
		return out, fmt.Errorf("max batch length cannot be 0")
	}

	// batch list into no more than max each
	for len(list) > max {
		// assign to list: remaining after taking slice from beginning
		// append to out: max length slice from beginning of list
		list, out = list[max:], append(out, list[:max])
	}
	out = append(out, list) // append remaining to list (slice len < max)
	return out, nil
}
