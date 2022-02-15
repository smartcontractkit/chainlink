package logger

// SugaredLogger extends the base Logger interface with syntactic sugar, similar to zap.SugaredLogger.
type SugaredLogger interface {
	Logger
	//TODO document
	AssumptionViolation(args ...interface{})
	AssumptionViolationf(format string, vals ...interface{})
	AssumptionViolationw(msg string, keyvals ...interface{})
}

func Sugared(l Logger) SugaredLogger {
	return &sugared{
		Logger: l,
		h:      l.Helper(1),
	}
}

type sugared struct {
	Logger
	h Logger // helper with stack trace skip level
}

func (s *sugared) AssumptionViolation(args ...interface{}) {
	s.h.Error(append([]interface{}{"AssumptionViolation:"}, args...))
}

func (s *sugared) AssumptionViolationf(format string, vals ...interface{}) {
	s.h.Errorf("AssumptionViolation: "+format, vals...)
}

func (s *sugared) AssumptionViolationw(msg string, keyvals ...interface{}) {
	s.h.Errorw("AssumptionViolation: "+msg, keyvals...)
}
