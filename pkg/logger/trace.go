//go:build trace

package logger

const tracePrefix = "[TRACE] "

// Tracew emits trace level logs, which are debug level with a '[trace]' prefix.
func Tracew(l Logger, msg string, keysAndValues ...interface{}) {
	l = Helper(l, 1)
	t, ok := l.(interface {
		Tracew(string, ...interface{})
	})
	if ok {
		t.Tracew(msg, keysAndValues...)
		return
	}
	l.Debugw(tracePrefix+msg, keysAndValues...)
}

// Tracef emits trace level logs, which are debug level with a '[trace]' prefix.
func Tracef(l Logger, format string, values ...interface{}) {
	l = Helper(l, 1)
	t, ok := l.(interface {
		Tracef(string, ...interface{})
	})
	if ok {
		t.Tracef(format, values...)
		return
	}
	l.Debugf(tracePrefix+format, values...)
}
