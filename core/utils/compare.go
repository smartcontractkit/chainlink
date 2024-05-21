package utils

func IsZero[C comparable](val C) bool {
	var zero C
	return zero == val
}
