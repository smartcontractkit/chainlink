package logger

import (
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
)

func TestOCRWrapper(t *testing.T) {
	t.Parallel()

	testFields := commontypes.LogFields{
		"field": 123,
	}

	args := []interface{}{"msg"}
	args = append(args, toKeysAndValues(testFields)...)

	ml := NewMockLogger(t)
	ml.On("Helper", 2).Return(ml).Once()
	ml.On("Debugw", args...).Twice() // due to Trace
	ml.On("Infow", args...).Once()
	ml.On("Warnw", args...).Once()
	ml.On("Criticalw", args...).Once()
	ml.On("Errorw", args...).Once()

	var savedError string
	saveError := func(err string) {
		savedError = err
	}

	w := NewOCRWrapper(ml, true, saveError)
	w.Trace("msg", testFields)
	w.Debug("msg", testFields)
	w.Info("msg", testFields)
	w.Warn("msg", testFields)
	w.Critical("msg", testFields)
	w.Error("msg", testFields)
	assert.Equal(t, "msg", savedError)
}

func TestOCRWrapper_NoTrace(t *testing.T) {
	t.Parallel()

	ml := NewMockLogger(t)
	ml.On("Helper", 2).Return(ml).Once()

	w := NewOCRWrapper(ml, false, nil)
	w.Trace("msg", commontypes.LogFields{})
}
