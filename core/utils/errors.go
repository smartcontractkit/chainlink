package utils

import (
	"fmt"
	"strings"
)

type multiErrorList []error

// MultiErrorList returns an error which formats underlying errors as a list, or nil if err is nil.
func MultiErrorList(err error) (int, error) {
	if err == nil {
		return 0, nil
	}
	errs := Flatten(err)
	return len(errs), multiErrorList(errs)
}

func (m multiErrorList) Error() string {
	l := len(m)
	if l == 1 {
		return m[0].Error()
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d errors:", l)
	for _, e := range m {
		fmt.Fprintf(&sb, "\n\t- %v", e)
	}
	return sb.String()
}

func (m multiErrorList) Unwrap() []error {
	return m
}

// Flatten calls `Unwrap() []error` on each error and subsequent returned error that implement the method, returning a fully flattend sequence.
//
//nolint:errorlint // error type checks will fail on wrapped errors. Disabled since we are not doing checks on error types.
func Flatten(errs ...error) (flat []error) {
	for _, err := range errs {
		if me, ok := err.(interface{ Unwrap() []error }); ok {
			flat = append(flat, Flatten(me.Unwrap()...)...)
			continue
		}
		flat = append(flat, err)
	}
	return
}
