package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestNullLogger(t *testing.T) {
	t.Parallel()

	t.Run("names", func(t *testing.T) {
		t.Parallel()

		l := logger.NullLogger
		assert.Equal(t, l, l.Named("foo"))
		assert.Equal(t, l, l.With("foo"))
		assert.Equal(t, l, l.Helper(123))
	})

	t.Run("no-op", func(t *testing.T) {
		t.Parallel()

		l := logger.NullLogger
		l.SetLogLevel(zapcore.DebugLevel)
		l.Trace()
		l.Debug()
		l.Info()
		l.Warn()
		l.Error()
		l.Critical()
		l.Panic()
		l.Fatal()
		l.Tracef("msg")
		l.Debugf("msg")
		l.Infof("msg")
		l.Warnf("msg")
		l.Errorf("msg")
		l.Criticalf("msg")
		l.Panicf("msg")
		l.Fatalf("msg")
		l.Tracew("msg")
		l.Debugw("msg")
		l.Infow("msg")
		l.Warnw("msg")
		l.Errorw("msg")
		l.Criticalw("msg")
		l.Panicw("msg")
		l.Fatalw("msg")
		l.Recover(nil)
		assert.Nil(t, l.Sync())
	})
}
