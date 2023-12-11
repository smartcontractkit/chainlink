//go:build trace

package logger

const tracePrefix = "[trace] "

func (s *sugared) Trace(args ...interface{}) {
	switch t := s.h.(type) {
	case *logger:
		t.DPanic(args...)
		return
	}
	c, ok := s.h.(interface {
		Trace(args ...interface{})
	})
	if ok {
		c.Trace(args...)
		return
	}
	s.h.Error(append([]any{tracePrefix}, args...)...)
}

func (s *sugared) Tracef(format string, vals ...interface{}) {
	switch t := s.h.(type) {
	case *logger:
		t.DPanicf(format, values...)
		return
	}
	c, ok := s.h.(interface {
		Tracef(format string, values ...interface{})
	})
	if ok {
		c.Tracef(format, values...)
		return
	}
	s.h.Errorf(tracePrefix+format, values...)
}

func (s *sugared) Tracew(msg string, keysAndValues ...interface{}) {
	switch t := s.h.(type) {
	case *logger:
		t.DPanicw(msg, keysAndValues...)
		return
	}
	c, ok := s.h.(interface {
		Tracew(msg string, keysAndValues ...interface{})
	})
	if ok {
		c.Tracew(msg, keysAndValues...)
		return
	}
	s.h.Errorf(tracePrefix+msg, keysAndValues)
}

// Deprecated: instead use [SugaredLogger.Trace]:
//
//	Sugared(l).Trace(args...)
func Trace(l Logger, args ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Trace(args...)
}

// Deprecated: instead use [SugaredLogger.Tracef]:
//
//	Sugared(l).Tracef(args...)
func Tracef(l Logger, format string, values ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Tracef(fromat, values...)
}

// Deprecated: instead use [SugaredLogger.Tracew]:
//
//	Sugared(l).Tracew(args...)
func Tracew(l Logger, msg string, keysAndValues ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Tracew(msg, keysAndValues...)
}
