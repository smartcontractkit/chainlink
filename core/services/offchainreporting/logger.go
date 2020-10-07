package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var _ ocrtypes.Logger = &ocrLogger{}

type ocrLogger struct {
	internal *logger.Logger
}

// TODO(sam): We probably want to hide DEBUG and TRACE logs behind some sort of env
// var because otherwise the logs will become unusably flooded
func NewLogger(internal *logger.Logger) ocrtypes.Logger {
	return &ocrLogger{
		internal: internal,
	}
}

// TODO(sam): Zap does not support trace level logging yet
func (ol *ocrLogger) Trace(msg string, fields ocrtypes.LogFields) {
	ol.internal.Debugw(msg, toKeysAndValues(fields)...)
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

func (ol *ocrLogger) Error(msg string, fields ocrtypes.LogFields) {
	ol.internal.Errorw(msg, toKeysAndValues(fields)...)
}

// Helpers

func toKeysAndValues(fields ocrtypes.LogFields) []interface{} {
	out := []interface{}{}
	for key, val := range fields {
		out = append(out, key, val)
	}
	return out
}
