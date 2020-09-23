package types

type Logger interface {
	Trace(msg string, fields LogFields)
	Debug(msg string, fields LogFields)
	Info(msg string, fields LogFields)
	Warn(msg string, fields LogFields)
	Error(msg string, fields LogFields)
}

type LogFields map[string]interface{}
