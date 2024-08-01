package services

import "context"

// A StopChan signals when some work should stop.
// Use StopChanR if you already have a read only <-chan.
type StopChan chan struct{}

// NewCtx returns a background [context.Context] that is cancelled when StopChan is closed.
func (s StopChan) NewCtx() (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).NewCtx()
}

// Ctx cancels a [context.Context] when StopChan is closed.
func (s StopChan) Ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).Ctx(ctx)
}

// CtxCancel cancels a [context.Context] when StopChan is closed.
// Returns ctx and cancel unmodified, for convenience.
func (s StopChan) CtxCancel(ctx context.Context, cancel context.CancelFunc) (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).CtxCancel(ctx, cancel)
}

// A StopRChan signals when some work should stop.
// This is a receive-only version of StopChan, for casting an existing <-chan.
type StopRChan <-chan struct{}

// NewCtx returns a background [context.Context] that is cancelled when StopChan is closed.
func (s StopRChan) NewCtx() (context.Context, context.CancelFunc) {
	return s.Ctx(context.Background())
}

// Ctx cancels a [context.Context] when StopChan is closed.
func (s StopRChan) Ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	return s.CtxCancel(context.WithCancel(ctx))
}

// CtxCancel cancels a [context.Context] when StopChan is closed.
// Returns ctx and cancel unmodified, for convenience.
func (s StopRChan) CtxCancel(ctx context.Context, cancel context.CancelFunc) (context.Context, context.CancelFunc) {
	go func() {
		select {
		case <-s:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}
