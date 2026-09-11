package utils

import "errors"

// BatchSplit splits an slices into an slices of slicess with a maximum length
func BatchSplit[T any](list []T, maxLen int) (out [][]T, err error) {
	if maxLen == 0 {
		return out, errors.New("max batch length cannot be 0")
	}

	// batch list into no more than max each
	for len(list) > maxLen {
		// assign to list: remaining after taking slice from beginning
		// append to out: max length slice from beginning of list
		list, out = list[maxLen:], append(out, list[:maxLen])
	}
	out = append(out, list) // append remaining to list (slice len < max)
	return out, nil
}
