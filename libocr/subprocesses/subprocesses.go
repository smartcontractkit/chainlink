package subprocesses

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Subprocesses struct {
	wg sync.WaitGroup
}

func (s *Subprocesses) Wait() {
	s.wg.Wait()
}

func (s *Subprocesses) Go(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

func (s *Subprocesses) BlockForAtMost(ctx context.Context, d time.Duration, f func(context.Context)) (ok bool) {
	done := make(chan struct{})
	childCtx, childCancel := context.WithTimeout(ctx, d)
	defer childCancel()
	s.Go(func() {
		f(childCtx)
		close(done)
	})

	select {
	case <-done:
		return true
	case <-childCtx.Done():
		return false
	}
}

func (s *Subprocesses) RepeatWithCancel(name string, interval time.Duration, ctx context.Context, f func()) {
	s.Go(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("canceling", name) 				return
			case <-ticker.C:
				f()
			}
		}
	})
}
