//go:build trace

package logger

const tracePrefix = "[trace] "

func Trace(l Logger, args ...interface{}) {
	switch t := l.(type) {
	case *logger:
		t.DPanic(args...)
		return
	}
	c, ok := l.(interface {
		Trace(args ...interface{})
	})
	if ok {
		c.Trace(args...)
		return
	}
	l.Error(append([]any{tracePrefix}, args...)...)
}

func Tracef(l Logger, format string, values ...interface{}) {
	switch t := l.(type) {
	case *logger:
		t.DPanicf(format, values...)
		return
	}
	c, ok := l.(interface {
		Tracef(format string, values ...interface{})
	})
	if ok {
		c.Tracef(format, values...)
		return
	}
	l.Errorf(tracePrefix+format, values...)
}

func Tracew(l Logger, msg string, keysAndValues ...interface{}) {
	switch t := l.(type) {
	case *logger:
		t.DPanicw(msg, keysAndValues...)
		return
	}
	c, ok := l.(interface {
		Tracew(msg string, keysAndValues ...interface{})
	})
	if ok {
		c.Tracew(msg, keysAndValues...)
		return
	}
	l.Errorf(tracePrefix+msg, keysAndValues)
}
