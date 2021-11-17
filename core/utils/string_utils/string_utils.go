package stringutils

import "strconv"

func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func FromInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func ToInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(n), nil
}

func FromInt32(n int32) string {
	return FromInt64(int64(n))
}
