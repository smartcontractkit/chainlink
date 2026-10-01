package triggers

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// _buffer is the observed batch buffer; a full channel drops observed batches
const _buffer = 1

// drainerHooks are the drainer's lifecycle hooks.
// Unset hooks default to no-ops, so the drainer never nil-checks.
type drainerHooks struct {
	// OnObservedEvents receives each non-empty batch the drainer observes. It
	// is fired asynchronously on a dedicated worker goroutine and never blocks
	// the drain loop.
	// Batches are dropped rather than blocking. It should honor ctx, which is
	// cancelled when the drainer is closed; Close waits for an in-flight call to return, and a
	// batch still buffered at that point is discarded.
	OnObservedEvents func(ctx context.Context, events []CoordinatedEvent)
}

// set all to non-nil so the drainer doesn't have to check before each call
func (h *drainerHooks) setDefaultHooks() {
	if h.OnObservedEvents == nil {
		h.OnObservedEvents = func(context.Context, []CoordinatedEvent) {}
	}
}

// drainer is a consumer of the Queue; on each tick it observes the queue.  Events
// may be handled by defining an OnObservedEvents hook.  The handler runs async
// to the drain loop.
type drainer struct {
	services.Service
	eng *services.Engine

	queue    Queue
	interval time.Duration
	hooks    drainerHooks

	// observed batches awaiting OnObservedEvents
	batches chan []CoordinatedEvent
	lggr    logger.Logger
}

// NewDrainer returns a drainer that calls queue.Observe every interval.
func NewDrainer(queue Queue, interval time.Duration, hooks drainerHooks, lggr logger.Logger) *drainer {
	hooks.setDefaultHooks()
	d := &drainer{
		queue:    queue,
		interval: interval,
		hooks:    hooks,
		batches:  make(chan []CoordinatedEvent, _buffer),
		lggr:     logger.Named(lggr, "TriggerQueueDrainer"),
	}
	d.Service, d.eng = services.Config{
		Name:  "TriggerQueueDrainer",
		Start: d.start,
	}.NewServiceEngine(d.lggr)
	return d
}

// start launches the hook worker and the tick loop. Both run until the
// engine's context is cancelled, which happens when the service is closed;
// Close then waits for both goroutines to exit. The ticker is stopped when the
// loop returns.
func (d *drainer) start(context.Context) error {
	d.eng.Go(d.runHooks)
	d.eng.GoTick(services.NewTicker(d.interval), d.drain)
	return nil
}

// runHooks delivers queued batches to OnObservedEvents until ctx is cancelled.
func (d *drainer) runHooks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case events := <-d.batches:
			d.hooks.OnObservedEvents(ctx, events)
		}
	}
}

// drain observes the queue once and hands the returned events to the hook
// worker without blocking. An Observe error is logged and does not stop the
// loop. The hook is not called for an empty batch.
func (d *drainer) drain(ctx context.Context) {
	events, err := d.queue.Observe(ctx)
	if err != nil {
		d.lggr.Errorw("Failed to observe trigger queue", "err", err)
		return
	}
	if len(events) == 0 {
		return
	}
	d.lggr.Debugw("Observed trigger events", "count", len(events))
	select {
	case d.batches <- events:
	default:
		// drop this batch rather than block the drain loop
	}
}
