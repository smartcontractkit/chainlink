package utils

import "time"

func FiniteTicker(period time.Duration, onTick func()) func() {
	tick := time.NewTicker(period)
	chStop := make(chan struct{})
	go func() {
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
	}
}
