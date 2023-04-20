package utils

func IsZero[C comparable](val C) bool {
	var zero C
	return zero == val
}
func Equal[C comparable](val, other C) bool {
	return val == other
}
