package stringutils

import "strconv"

// ToInt64 parses s as a base 10 int64.
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// FromInt64 formats n as a base 10 string.
func FromInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}

// ToInt32 parses s as a base 10 int32.
func ToInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(n), nil
}

// FromInt32 formats n as a base 10 string.
func FromInt32(n int32) string {
	return FromInt64(int64(n))
}
