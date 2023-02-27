package logger

import "go.uber.org/zap/zapcore"

// encodeLevel is a zapcore.EncodeLevel that encodes 'crit' in place of dpanic for our custom Critical* level.
func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == zapcore.DPanicLevel {
		enc.AppendString("crit")
	} else {
		zapcore.LowercaseLevelEncoder(l, enc)
	}
}

func (l *zapLogger) Critical(args ...any) {
	// DPanic is used for the appropriate numerical level (between error and panic), but we never actually panic.
	l.sugaredHelper(1).DPanic(args...)
}

func (l *zapLogger) Criticalf(format string, values ...any) {
	l.sugaredHelper(1).DPanicf(format, values...)
}

func (l *zapLogger) Criticalw(msg string, keysAndValues ...any) {
	l.sugaredHelper(1).DPanicw(msg, keysAndValues...)
}
