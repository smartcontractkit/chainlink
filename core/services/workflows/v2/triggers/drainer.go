package triggers

import (
	"context"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// _buffer is the observed batch buffer; a full channel drops observed batches
const _buffer = 1

// minDrainInterval is the shortest wait between drains. Shorter configured intervals are raised to it.
const minDrainInterval = 200 * time.Millisecond

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
//
// drainer exists to be simple.  Do not attempt to overload OnObservedEvents and
// ensure that hook definitions respect the passed context.
type drainer struct {
	services.Service
	eng *services.Engine

	queue    Queue
	interval limits.BoundLimiter[time.Duration]
	clock    clockwork.Clock
	hooks    drainerHooks

	// observed batches awaiting OnObservedEvents
	batches chan []CoordinatedEvent
	metrics *monitoring.TriggerDrainerMetricLabeler
	lggr    logger.Logger
}

// NewDrainer returns a drainer that calls queue.Drain every interval. The interval is read again before
// every wait, so a settings change applies once the current wait ends. Values below 200ms are raised to
// 200ms, and a failed read uses the CentralTriggerQueueDrainInterval default. metrics must be non-nil.
func NewDrainer(queue Queue, interval limits.BoundLimiter[time.Duration], hooks drainerHooks, metrics *monitoring.TriggerDrainerMetricLabeler, clock clockwork.Clock, lggr logger.Logger) *drainer {
	hooks.setDefaultHooks()
	d := &drainer{
		queue:    queue,
		interval: interval,
		clock:    clock,
		hooks:    hooks,
		batches:  make(chan []CoordinatedEvent, _buffer),
		metrics:  metrics,
		lggr:     logger.Named(lggr, "TriggerQueueDrainer"),
	}
	d.Service, d.eng = services.Config{
		Name:  "TriggerQueueDrainer",
		Start: d.start,
	}.NewServiceEngine(d.lggr)
	return d
}

// start launches the hook worker and the drain loop. Both run until the
// engine's context is cancelled, which happens when the service is closed;
// Close then waits for both goroutines to exit.
func (d *drainer) start(context.Context) error {
	d.eng.Go(d.runHooks)
	d.eng.Go(d.run)
	return nil
}

// run drains the queue after each wait until ctx is cancelled. The interval is read before every wait, so a
// changed setting applies once the current wait ends.
func (d *drainer) run(ctx context.Context) {
	for {
		timer := d.clock.NewTimer(d.nextInterval(ctx))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.Chan():
		}
		d.drain(ctx)
	}
}

// nextInterval reads the interval, using the setting's default when the read fails and minDrainInterval when
// the value is shorter.
func (d *drainer) nextInterval(ctx context.Context) time.Duration {
	interval, err := d.interval.Limit(ctx)
	if err != nil {
		d.lggr.Warnw("Failed to read the drain interval, using the default", "err", err)
		interval = cresettings.Default.CentralTriggerQueueDrainInterval.DefaultValue
	}
	return max(interval, minDrainInterval)
}

// runHooks delivers queued batches to OnObservedEvents until ctx is cancelled.
func (d *drainer) runHooks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case events := <-d.batches:
			start := time.Now()
			d.hooks.OnObservedEvents(ctx, events)
			d.metrics.RecordHookDurationSeconds(ctx, time.Since(start).Seconds())
		}
	}
}

// drain empties the queue once and hands the returned events to the hook
// worker without blocking. A Drain error is logged and does not stop the
// loop. The hook is not called for an empty batch.
func (d *drainer) drain(ctx context.Context) {
	events, err := d.queue.Drain(ctx)
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
		d.metrics.IncrementDroppedEventsCounter(ctx, monitoring.TriggerDrainerDropReasonHookBusy, len(events))
		d.lggr.Debugw("Events handler is blocked, dropping the batch", "count", len(events))
	}
}
