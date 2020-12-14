package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var _ ocrtypes.Logger = &ocrLogger{}

type ocrLogger struct {
	internal  *logger.Logger
	trace     bool
	saveError func(string)
}

func NewLogger(internal *logger.Logger, trace bool, saveError func(string)) ocrtypes.Logger {
	return &ocrLogger{
		internal:  internal,
		trace:     trace,
		saveError: saveError,
	}
}

// TODO(sam): Zap does not support trace level logging yet, so this hack is
// necessary to silence excessive logging
func (ol *ocrLogger) Trace(msg string, fields ocrtypes.LogFields) {
	if ol.trace {
		ol.internal.Debugw(msg, toKeysAndValues(fields)...)
	}
}

func (ol *ocrLogger) Debug(msg string, fields ocrtypes.LogFields) {
	ol.internal.Debugw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrLogger) Info(msg string, fields ocrtypes.LogFields) {
	ol.internal.Infow(msg, toKeysAndValues(fields)...)
}

func (ol *ocrLogger) Warn(msg string, fields ocrtypes.LogFields) {
	ol.internal.Warnw(msg, toKeysAndValues(fields)...)
}

// Note that the structured fields may contain dynamic data (timestamps etc.)
// So when saving the error, we only save the top level string, details
// are included in the log.
func (ol *ocrLogger) Error(msg string, fields ocrtypes.LogFields) {
	ol.saveError(msg)
	ol.internal.Errorw(msg, toKeysAndValues(fields)...)
}

func toKeysAndValues(fields ocrtypes.LogFields) []interface{} {
	out := []interface{}{}
	for key, val := range fields {
		out = append(out, key, val)
	}
	return out
}
