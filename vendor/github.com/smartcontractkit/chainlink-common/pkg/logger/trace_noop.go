//go:build !trace

package logger

func (s *sugared) Trace(args ...interface{}) {}

func (s *sugared) Tracef(format string, vals ...interface{}) {}

func (s *sugared) Tracew(msg string, keysAndValues ...interface{}) {}

// Deprecated: instead use [SugaredLogger.Trace]:
//
//	Sugared(l).Trace(args...)
func Trace(l Logger, args ...interface{}) {}

// Deprecated: instead use [SugaredLogger.Tracef]:
//
//	Sugared(l).Tracef(args...)
func Tracef(l Logger, format string, values ...interface{}) {}

// Deprecated: instead use [SugaredLogger.Tracew]:
//
//	Sugared(l).Tracew(args...)
func Tracew(l Logger, msg string, keysAndValues ...interface{}) {}
