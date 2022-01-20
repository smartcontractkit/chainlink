package shutdown

import (
	"os"
	ossignal "os/signal"
	"syscall"
)

// CancelOnShutdown waits for SIGINT/SIGTERM signals and calls cancelFunc
func CancelOnShutdown(cancelFunc func()) {
	ch := make(chan os.Signal, 1)
	ossignal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	cancelFunc()
}
