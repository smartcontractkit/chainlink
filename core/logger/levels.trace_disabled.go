//go:build !trace

package logger

func (l *zapLogger) Trace(args ...interface{}) {}

func (l *zapLogger) Tracef(format string, values ...interface{}) {}

func (l *zapLogger) Tracew(msg string, keysAndValues ...interface{}) {}
