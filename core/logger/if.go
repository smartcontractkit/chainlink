package logger

type Logger interface {
	Loggable

	Swap(Logger) Logger
	With(...interface{}) Logger

	Sync() error
}

type Loggable interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugw(string, ...interface{})

	Error(...interface{})
	ErrorIf(error, ...string)
	ErrorIfCalling(func() error, ...string)
	Errorf(string, ...interface{})
	Errorw(string, ...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalw(string, ...interface{})

	Info(...interface{})
	Infof(string, ...interface{})
	Infow(string, ...interface{})

	Panic(...interface{})
	PanicIf(error)
	Panicf(string, ...interface{})

	Trace(...interface{})
	Tracef(string, ...interface{})
	Tracew(string, ...interface{})

	Warn(...interface{})
	WarnIf(error)
	Warnf(string, ...interface{})
	Warnw(string, ...interface{})
}

// FIXME: this is ugly, make it a function
// NewErrorw(string, ...interface{})
