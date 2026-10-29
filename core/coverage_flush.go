//go:build coverage

package core

import (
	"os"
	"os/signal"
	"runtime/coverage"
	"syscall"
)

// init installs a SIGUSR1 handler used by E2E coverage collection. The handler
// flushes coverage counters while the process is still running, so coverage is
// not lost when a container is later stopped with SIGKILL.
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
}
