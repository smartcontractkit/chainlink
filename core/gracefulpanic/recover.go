package gracefulpanic

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"runtime/debug"
)

func WrapRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Default.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))
		}
	}()
	fn()
}
