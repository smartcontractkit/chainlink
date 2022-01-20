package monitoring

type Logger interface {
	With(args ...interface{}) Logger

	Tracew(format string, values ...interface{})
	Debugw(format string, values ...interface{})
	Infow(format string, values ...interface{})
	Warnw(format string, values ...interface{})
	Errorw(format string, values ...interface{})
	Criticalw(format string, values ...interface{})
	Panicw(format string, values ...interface{})
	Fatalw(format string, values ...interface{})
}
