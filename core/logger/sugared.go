package logger

// SugaredLogger extends the base Logger interface with syntactic sugar, similar to zap.SugaredLogger.
type SugaredLogger interface {
	Logger
	// AssumptionViolation variants log at error level with the message prefix "AssumptionViolation: ".
	AssumptionViolation(args ...interface{})
	AssumptionViolationf(format string, vals ...interface{})
	AssumptionViolationw(msg string, keyvals ...interface{})
	// ErrorIf logs the error if present.
	ErrorIf(err error, msg string)
	// ErrorIfFn calls fn() and logs any returned error along with msg.
	// Unlike ErrorIf, this can be deffered inline, since the function call is delayed.
	ErrorIfFn(fn func() error, msg string)
}

// Sugared returns a new SugaredLogger wrapping the given Logger.
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

// AssumptionViolation wraps Error logs with assumption violation tag.
func (s *sugared) AssumptionViolation(args ...interface{}) {
	s.h.Error(append([]interface{}{"AssumptionViolation:"}, args...))
}

// AssumptionViolationf wraps Errorf logs with assumption violation tag.
func (s *sugared) AssumptionViolationf(format string, vals ...interface{}) {
	s.h.Errorf("AssumptionViolation: "+format, vals...)
}

// AssumptionViolationw wraps Errorw logs with assumption violation tag.
func (s *sugared) AssumptionViolationw(msg string, keyvals ...interface{}) {
	s.h.Errorw("AssumptionViolation: "+msg, keyvals...)
}

func (s *sugared) ErrorIf(err error, msg string) {
	if err != nil {
		s.h.Errorw(msg, "err", err)
	}
}

func (s *sugared) ErrorIfFn(fn func() error, msg string) {
	if err := fn(); err != nil {
		s.h.Errorw(msg, "err", err)
	}
}
