package slicelib

// GroupBy groups a slice based on a specific item property. The returned groups slice is deterministic.
func GroupBy[T any, K comparable](items []T, prop func(T) K) ([]K, map[K][]T) {
	groups := make([]K, 0)
	grouped := make(map[K][]T)
	for _, item := range items {
		k := prop(item)
		if _, exists := grouped[k]; !exists {
			groups = append(groups, k)
		}
		grouped[k] = append(grouped[k], item)
	}
	return groups, grouped
}

// CountUnique counts the unique items of the provided slice.
func CountUnique[T comparable](items []T) int {
	m := make(map[T]struct{})
	for _, item := range items {
		m[item] = struct{}{}
	}
	return len(m)
}

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](slices [][]T) []T {
	res := make([]T, 0)
	for _, s := range slices {
		res = append(res, s...)
	}
	return res
}
