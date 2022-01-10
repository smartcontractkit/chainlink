package logger

func (l *zapLogger) Critical(args ...interface{}) {
	// DPanic is used for the appropriate numerical level (between error and panic), but we never actually panic.
	l.sugaredHelper(1).DPanic(args...)
}

func (l *zapLogger) Criticalf(format string, values ...interface{}) {
	l.sugaredHelper(1).DPanicf(format, values...)
}

func (l *zapLogger) CriticalW(msg string, keysAndValues ...interface{}) {
	l.sugaredHelper(1).DPanicw(msg, keysAndValues...)
}
