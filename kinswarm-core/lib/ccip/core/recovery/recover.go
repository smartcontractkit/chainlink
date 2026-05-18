package recovery

import (
	"github.com/getsentry/sentry-go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

func ReportPanics(fn func()) {
	HandleFn(fn, func(err any) {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(corelogger.SentryFlushDeadline)

		panic(err)
	})
}

func WrapRecover(lggr logger.Logger, fn func()) {
	WrapRecoverHandle(lggr, fn, nil)
}

func WrapRecoverHandle(lggr logger.Logger, fn func(), onPanic func(recovered any)) {
	HandleFn(fn, func(recovered any) {
		logger.Sugared(lggr).Criticalw("Recovered goroutine panic", "panic", recovered)

		if onPanic != nil {
			onPanic(recovered)
		}
	})
}

func HandleFn(fn func(), onPanic func(recovered any)) {
	defer func() {
		if recovered := recover(); recovered != nil {
			onPanic(recovered)
		}
	}()
	fn()
}
