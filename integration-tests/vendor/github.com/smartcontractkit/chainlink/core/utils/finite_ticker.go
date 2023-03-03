package utils

import (
	"sync"
	"time"
)

// FiniteTicker starts a goroutine to execute the given function periodically, until the returned function is called.
func FiniteTicker(period time.Duration, onTick func()) func() {
	tick := time.NewTicker(period)
	chStop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-tick.C:
				onTick()
			case <-chStop:
				return
			}
		}
	}()

	// NOTE: tick.Stop does not close the ticker channel,
	// so we still need another way of returning (chStop).
	return func() {
		tick.Stop()
		close(chStop)
		wg.Wait()
	}
}
