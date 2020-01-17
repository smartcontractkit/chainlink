package utils

// Copied from bleeding-edge testify
// https://github.com/stretchr/testify/blob/8c465a0c/assert/assertions.go
//
// Remove this and redirect its callers to testify's assert/require, once
// chainlink depends on a version greater than v1.4.0

import (
	"fmt"
	"runtime/debug"

	"github.com/stretchr/testify/assert"
)

type tHelper interface {
	Helper()
}

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (bool, interface{}, string) {

	didPanic := false
	var message interface{}
	var stack string
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
				stack = string(debug.Stack())
			}
		}()

		// call the target function
		f()

	}()

	return didPanic, message, stack

}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   assert.PanicsWithError(t, "crazy error", func(){ GoCrazy() })
func PanicsWithError(t assert.TestingT, errString string, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	funcDidPanic, panicValue, panickedStack := didPanic(f)
	if !funcDidPanic {
		return assert.Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}
	panicErr, ok := panicValue.(error)
	if !ok || panicErr.Error() != errString {
		return assert.Fail(t, fmt.Sprintf("func %#v should panic with error message:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s", f, errString, panicValue, panickedStack), msgAndArgs...)
	}

	return true
}
