//go:build !trace

package logger

func Trace(l Logger, args ...interface{}) {}

func Tracef(l Logger, format string, values ...interface{}) {}

func Tracew(l Logger, msg string, keysAndValues ...interface{}) {}
