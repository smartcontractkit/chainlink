package utils

import (
	"context"
	"sync"
	"time"
)

func ExampleStopChan() {
	stopCh := make(StopChan)

	work := func(context.Context) {}

	a := func(ctx context.Context, done func()) {
		defer done()
		ctx, cancel := stopCh.Ctx(ctx)
		defer cancel()
		work(ctx)
	}

	b := func(ctx context.Context, done func()) {
		defer done()
		ctx, cancel := stopCh.CtxCancel(context.WithTimeout(ctx, time.Minute))
		defer cancel()
		work(ctx)
	}

	c := func(ctx context.Context, done func()) {
		defer done()
		ctx, cancel := stopCh.CtxCancel(context.WithDeadline(ctx, time.Now().Add(5*time.Minute)))
		defer cancel()
		work(ctx)
	}

	ctx, cancel := stopCh.NewCtx()
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)
	go a(ctx, wg.Done)
	go b(ctx, wg.Done)
	go c(ctx, wg.Done)

	time.AfterFunc(time.Second, func() { close(stopCh) })

	wg.Wait()
	// Output:
}
