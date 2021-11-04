package shutdown

import (
	"os"
	"runtime/debug"

	"github.com/smartcontractkit/chainlink/core/logger"
)

var HardPanic bool

func init() {
	if os.Getenv("ENABLE_HARD_PANIC") == "true" {
		HardPanic = true
	}
}

func WrapRecover(lggr logger.Logger, fn func()) {
	if !HardPanic {
		defer func() {
			if err := recover(); err != nil {
				lggr.Errorw("goroutine panicked", "panic", err, "stacktrace", string(debug.Stack()))
			}
		}()
	}
	fn()
}
