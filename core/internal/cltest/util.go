package cltest

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

func MustHexDecodeString(s string) []byte {
	a, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return a
}

// CommitLoop creates a goroutine that repeatedly commits and returns a
// cancellation function to avoid leaking goroutines
func CommitLoop(b *backends.SimulatedBackend) func() {
	tick := time.NewTicker(500 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-tick.C:
				b.Commit()
			case <-done:
				return
			}
		}
	}()
	return func() {
		tick.Stop()
		close(done)
	}
}
