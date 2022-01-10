//go:build trace

package logger

import "fmt"

const tracePrefix = "[TRACE] "

func (l *zapLogger) Trace(args ...interface{}) {
	args[0] = fmt.Sprint(tracePrefix, args[0])
	l.sugaredHelper(1).Debug(args...)
}

func (l *zapLogger) Tracef(format string, values ...interface{}) {
	l.sugaredHelper(1).Debugf(fmt.Sprint(tracePrefix, format), values...)
}

func (l *zapLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.sugaredHelper(1).Debugw(fmt.Sprint(tracePrefix, msg), keysAndValues...)
}
