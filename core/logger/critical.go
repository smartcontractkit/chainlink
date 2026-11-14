package logger

func (l *zapLogger) Critical(args ...any) {
	// DPanic is used for the appropriate numerical level (between error and panic), but we never actually panic.
	l.sugaredHelper(1).DPanic(args...)
}

func (l *zapLogger) Criticalf(format string, values ...any) {
	l.sugaredHelper(1).DPanicf(format, values...)
}

func (l *zapLogger) Criticalw(msg string, keysAndValues ...any) {
	l.sugaredHelper(1).DPanicw(msg, keysAndValues...)
}
