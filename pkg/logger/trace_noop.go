//go:build !trace

package logger

// Tracew is a no-op.
func Tracew(l Logger, msg string, keysAndValues ...interface{}) {}

// Tracef is a no-op.
func Tracef(l Logger, format string, values ...interface{}) {}
