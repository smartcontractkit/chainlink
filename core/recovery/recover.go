package recovery

import (
	"runtime/debug"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func WrapRecover(lggr logger.Logger, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			lggr.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))
		}
	}()
	fn()
}

func WrapRecoverHandle(lggr logger.Logger, fn func(), onPanic func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			lggr.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))

			if onPanic != nil {
				onPanic(err)
			}
		}
	}()
	fn()
}
