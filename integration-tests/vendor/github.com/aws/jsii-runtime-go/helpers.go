package jsii

import "time"

// Bool obtains a pointer to the provided bool.
func Bool(v bool) *bool { return &v }

// Bools obtains a pointer to a slice of pointers to all the provided booleans.
func Bools(v ...bool) *[]*bool {
	slice := make([]*bool, len(v))
	for i := 0; i < len(v); i++ {
		slice[i] = Bool(v[i])
	}
	return &slice
}

// Number obtains a pointer to the provided float64.
func Number(v float64) *float64 { return &v }

// Numbers obtains a pointer to a slice of pointers to all the provided numbers.
func Numbers(v ...float64) *[]*float64 {
	slice := make([]*float64, len(v))
	for i := 0; i < len(v); i++ {
		slice[i] = Number(v[i])
	}
	return &slice
}

// String obtains a pointer to the provided string.
func String(v string) *string { return &v }

// Strings obtains a pointer to a slice of pointers to all the provided strings.
func Strings(v ...string) *[]*string {
	slice := make([]*string, len(v))
	for i := 0; i < len(v); i++ {
		slice[i] = String(v[i])
	}
	return &slice
}

// Time obtains a pointer to the provided time.Time.
func Time(v time.Time) *time.Time { return &v }

// Times obtains a pointer to a slice of pointers to all the provided time.Time.
func Times(v ...time.Time) *[]*time.Time {
	slice := make([]*time.Time, len(v))
	for i := 0; i < len(v); i++ {
		slice[i] = Time(v[i])
	}
	return &slice
}
