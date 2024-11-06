package chans

import "sync"

// OrDone returns a channel that forwards values from `in` until `done` is closed.
// It stops processing and closes the output channel if `done` is closed.
func OrDone[T any](done <-chan struct{}, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()
	return out
}

// Tee splits an input channel into two output channels, each receiving every
// value from the input channel until `done` is closed. It guarantees that each
// value from `in` is sent once to each output channel.
func Tee[T any](done <-chan struct{}, in <-chan T) (out1, out2 <-chan T) {
	o1 := make(chan T)
	o2 := make(chan T)

	go func() {
		defer close(o1)
		defer close(o2)
		for val := range OrDone(done, in) {
			var ch1, ch2 = o1, o2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case ch1 <- val:
					ch1 = nil
				case ch2 <- val:
					ch2 = nil
				}
			}
		}
	}()

	return o1, o2
}

// Merge takes multiple input channels and merges their values into a single output channel.
// It stops merging if `done` is closed, and closes the output channel once all input channels are closed.
func Merge[T any](done <-chan struct{}, channels ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, ch := range channels {
			go func(ch <-chan T) {
				for val := range OrDone(done, ch) {
					select {
					case <-done:
						return
					case out <- val:
					}
				}
			}(ch)
		}
	}()
	return out
}

// AllClosed returns a channel that closes when either:
// - the `stop` channel is closed, or
// - all channels in `channels` are closed.
func AllClosed[T any](stop <-chan struct{}, channels ...<-chan T) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		ch := ch // capture channel variable for each goroutine
		go func() {
			defer wg.Done()
			for range ch {
				// Drain the channel until it closes
			}
		}()
	}

	go func() {
		defer close(done)
		// Wait for either all channels to close or the stop signal
		doneCh := make(chan struct{})
		go func() {
			wg.Wait()
			close(doneCh)
		}()

		select {
		case <-stop:
		case <-doneCh:
		}
	}()

	return done
}

// Transform applies a function to each value from the input channel and sends the result to the output channel.
func Transform[T any, U any](stop <-chan struct{}, fn func(T) U, in <-chan T) <-chan U {
	out := make(chan U)
	go func() {
		defer close(out)
		for {
			select {
			case <-stop:
				return
			case item, open := <-in:
				if !open {
					return
				}
				out <- fn(item)
			}
		}
	}()
	return out
}
