package shutdown

import (
	"os"
	ossignal "os/signal"
	"syscall"
)

// HandleShutdown waits for SIGINT/SIGTERM signals and calls handleFunc
func HandleShutdown(handleFunc func()) {
	ch := make(chan os.Signal, 1)
	ossignal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	handleFunc()
}
