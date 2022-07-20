package shutdown

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandleShutdown(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)

	tests := map[string]os.Signal{
		"SIGINT":  syscall.SIGINT,
		"SIGTERM": syscall.SIGTERM,
	}

	for name, sig := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			go HandleShutdown(func(string) {
				cancel()
			})

			// have to wait for ossignal.Notify
			time.Sleep(time.Second)

			err = proc.Signal(sig)
			require.NoError(t, err)

			select {
			case <-ctx.Done():
				// all good
			case <-time.After(3 * time.Second):
				require.Fail(t, "context is not cancelled within 3 seconds")
			}
		})
	}
}
