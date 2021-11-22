package recovery

import (
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func ReportPanics(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(logger.SentryFlushDeadline)

			panic(err)
		}
	}()
	fn()
}

func WrapRecover(lggr logger.Logger, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(logger.SentryFlushDeadline)

			lggr.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))
		}
	}()
	fn()
}

func WrapRecoverHandle(lggr logger.Logger, fn func(), onPanic func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(logger.SentryFlushDeadline)

			lggr.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))

			if onPanic != nil {
				onPanic(err)
			}
		}
	}()
	fn()
}
