package signalers

import "time"

func MakeTicker(stop <-chan struct{}, d time.Duration) <-chan struct{} {
	ticker := make(chan struct{})
	internalTicker := time.NewTicker(d)

	go func() {
		defer close(ticker)
		defer internalTicker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-internalTicker.C:
				ticker <- struct{}{}
			}
		}
	}()

	return ticker
}
