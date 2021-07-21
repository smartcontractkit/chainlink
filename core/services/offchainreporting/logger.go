package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"go.uber.org/zap"
)

var _ ocrcommontypes.Logger = &ocrLogger{}

type ocrLogger struct {
	internal  *logger.Logger
	trace     bool
	saveError func(string)
}

func NewLogger(l *logger.Logger, trace bool, saveError func(string)) ocrcommontypes.Logger {
	internal := logger.CreateLogger(l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar())
	return &ocrLogger{
		internal:  internal,
		trace:     trace,
		saveError: saveError,
	}
}

// TODO(sam): Zap does not support trace level logging yet, so this hack is
// necessary to silence excessive logging
func (ol *ocrLogger) Trace(msg string, fields ocrcommontypes.LogFields) {
	if ol.trace {
		ol.internal.Debugw(msg, toKeysAndValues(fields)...)
	}
}

func (ol *ocrLogger) Debug(msg string, fields ocrcommontypes.LogFields) {
	ol.internal.Debugw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrLogger) Info(msg string, fields ocrcommontypes.LogFields) {
	ol.internal.Infow(msg, toKeysAndValues(fields)...)
}

func (ol *ocrLogger) Warn(msg string, fields ocrcommontypes.LogFields) {
	ol.internal.Warnw(msg, toKeysAndValues(fields)...)
}

// Note that the structured fields may contain dynamic data (timestamps etc.)
// So when saving the error, we only save the top level string, details
// are included in the log.
func (ol *ocrLogger) Error(msg string, fields ocrcommontypes.LogFields) {
	ol.saveError(msg)
	ol.internal.Errorw(msg, toKeysAndValues(fields)...)
}

func toKeysAndValues(fields ocrcommontypes.LogFields) []interface{} {
	out := []interface{}{}
	for key, val := range fields {
		out = append(out, key, val)
	}
	return out
}
