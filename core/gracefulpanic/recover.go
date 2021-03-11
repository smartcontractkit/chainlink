package gracefulpanic

import "github.com/smartcontractkit/chainlink/core/logger"

func WrapRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Default.Errorw("goroutine panicked", "panic", err)
		}
	}()
	fn()
}
