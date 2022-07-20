package logger

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ ocrtypes.Logger = &ocrWrapper{}

type ocrWrapper struct {
	internal  Logger
	trace     bool
	saveError func(string)
}

func NewOCRWrapper(l Logger, trace bool, saveError func(string)) ocrtypes.Logger {
	return &ocrWrapper{
		internal:  l.Helper(2),
		trace:     trace,
		saveError: saveError,
	}
}

// TODO(sam): Zap does not support trace level logging yet, so this hack is
// necessary to silence excessive logging
func (ol *ocrWrapper) Trace(msg string, fields ocrtypes.LogFields) {
	if ol.trace {
		ol.internal.Debugw(msg, toKeysAndValues(fields)...)
	}
}

func (ol *ocrWrapper) Debug(msg string, fields ocrtypes.LogFields) {
	ol.internal.Debugw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Info(msg string, fields ocrtypes.LogFields) {
	ol.internal.Infow(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Warn(msg string, fields ocrtypes.LogFields) {
	ol.internal.Warnw(msg, toKeysAndValues(fields)...)
}

// Note that the structured fields may contain dynamic data (timestamps etc.)
// So when saving the error, we only save the top level string, details
// are included in the log.
func (ol *ocrWrapper) Error(msg string, fields ocrtypes.LogFields) {
	ol.saveError(msg)
	ol.internal.Errorw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Critical(msg string, fields ocrtypes.LogFields) {
	ol.internal.Criticalw(msg, toKeysAndValues(fields)...)
}

func toKeysAndValues(fields ocrtypes.LogFields) []interface{} {
	out := []interface{}{}
	for key, val := range fields {
		out = append(out, key, val)
	}
	return out
}
