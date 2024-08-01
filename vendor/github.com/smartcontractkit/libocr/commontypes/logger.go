package commontypes

// Loggers logs things using a structured-logging approach.
// All its functions should be thread-safe.
// It is acceptable to pass a nil LogFields to all of its functions.
type Logger interface {
	Trace(msg string, fields LogFields)
	Debug(msg string, fields LogFields)
	Info(msg string, fields LogFields)
	Warn(msg string, fields LogFields)
	Error(msg string, fields LogFields)
	Critical(msg string, fields LogFields)
}

type LogFields map[string]interface{}
