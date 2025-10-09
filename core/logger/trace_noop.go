//go:build !trace

package logger

func (l *zapLogger) Trace(args ...any) {}

func (l *zapLogger) Tracef(format string, values ...any) {}

func (l *zapLogger) Tracew(msg string, keysAndValues ...any) {}
