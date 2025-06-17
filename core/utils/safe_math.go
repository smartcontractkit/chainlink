package utils

// safely cast int64 to uint64, treating negative values as 0
func NonNegativeInt64ToUint64(i int64) uint64 {
	if i < 0 {
		i = 0
	}
	return uint64(i) //nolint:gosec // already checked for negative
}
