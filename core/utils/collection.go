package utils

import (
	"errors"
	"maps"
	"slices"
)

// BatchSplit splits a slice into a slice of slices with a maximum length.
// Returns an error if max is less than or equal to zero.
func BatchSplit[T any](list []T, max int) (out [][]T, err error) {
	if max <= 0 {
		return out, errors.New("max batch length must be greater than 0")
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

// Flatten takes a slice of slices and returns a single concatenated slice.
func Flatten[T any](lists [][]T) []T {
	var total int
	for _, l := range lists {
		total += len(l)
	}
	result := make([]T, 0, total)
	for _, l := range lists {
		result = append(result, l...)
	}
	return result
}

// UniqueValues returns the unique values from a map, in sorted order if the
// values are comparable. The caller receives a new slice.
func UniqueValues[K comparable, V comparable](m map[K]V) []V {
	seen := make(map[V]struct{}, len(m))
	result := make([]V, 0, len(m))
	for _, v := range m {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// MergeMaps merges multiple maps into one. Later maps take precedence for
// duplicate keys.
func MergeMaps[K comparable, V any](ms ...map[K]V) map[K]V {
	out := make(map[K]V)
	for _, m := range ms {
		maps.Copy(out, m)
	}
	return out
}

// FilterSlice returns a new slice containing only the elements for which the
// predicate returns true.
func FilterSlice[T any](s []T, pred func(T) bool) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if pred(v) {
			result = append(result, v)
		}
	}
	return slices.Clip(result)
}

// MapSlice applies a function to each element of a slice and returns the results.
func MapSlice[T any, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// ContainsDuplicate returns true if the slice contains duplicate elements.
func ContainsDuplicate[T comparable](s []T) bool {
	seen := make(map[T]struct{}, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}
	return false
}
