package chans

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
