package logger

import (
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ ocrtypes.Logger = &ocrWrapper{}

type ocrWrapper struct {
	l         Logger
	saveError func(string)
}

// NewOCRWrapper returns a new [ocrtypes.Logger] backed by the given Logger.
func NewOCRWrapper(l Logger, saveError func(string)) ocrtypes.Logger {
	return &ocrWrapper{
		// Skip an extra level since we are passed along to another wrapper.
		l:         Helper(l, 2),
		saveError: saveError,
	}
}
func (ol *ocrWrapper) Trace(msg string, fields ocrtypes.LogFields) {
	Tracew(ol.l, msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Debug(msg string, fields ocrtypes.LogFields) {
	ol.l.Debugw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Info(msg string, fields ocrtypes.LogFields) {
	ol.l.Infow(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Warn(msg string, fields ocrtypes.LogFields) {
	ol.l.Warnw(msg, toKeysAndValues(fields)...)
}

// Note that the structured fields may contain dynamic data (timestamps etc.)
// So when saving the error, we only save the top level string, details
// are included in the log.
func (ol *ocrWrapper) Error(msg string, fields ocrtypes.LogFields) {
	ol.saveError(msg)
	ol.l.Errorw(msg, toKeysAndValues(fields)...)
}

func (ol *ocrWrapper) Critical(msg string, fields ocrtypes.LogFields) {
	Criticalw(ol.l, msg, toKeysAndValues(fields)...)
}

func toKeysAndValues(fields ocrtypes.LogFields) []interface{} {
	out := []interface{}{}
	for key, val := range fields {
		out = append(out, key, val)
	}
	return out
}
