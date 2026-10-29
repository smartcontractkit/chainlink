//go:build coverage

package core

import (
	"os"
	"os/signal"
	"runtime/coverage"
	"syscall"
	"time"
)

// init installs a SIGUSR1 handler and a periodic flush used by E2E coverage
// collection. Flushing while the process is running means coverage is not lost
// when a container is later stopped with SIGKILL.
func init() {
	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)
	go func() {
		for range ch {
			_ = coverage.WriteCountersDir(dir)
		}
	}()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = coverage.WriteCountersDir(dir)
		}
	}()
}
