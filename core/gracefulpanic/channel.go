package gracefulpanic

import "sync"

var ch = make(chan struct{})
var panicOnce sync.Once

func Panic() {
	panicOnce.Do(func() {
		go close(ch)
	})
}

func Wait() <-chan struct{} {
	return ch
}
