package v2

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var errOCRQueueNotImplemented = errors.New("OCRQueue: draft, use NewOCRQueueWithInnerQueue for delegating implementation")

// OCR3_1PluginFactory creates OCR 3.1 ReportingPlugins (Blob + KV).
// Placeholder for design; use triggerqueue.Factory when implemented.
type OCR3_1PluginFactory any

// OCRQueueDeps holds the planned dependencies for NewOCRQueue from the design.
type OCRQueueDeps struct {
	Lf            limits.Factory
	Cfg           *cresettings.Workflows
	PluginFactory OCR3_1PluginFactory
	Lggr          logger.Logger

	// Subscribe to receive capabilities.DON (members, F, etc.)
	// on each update; reconfigure OCR oracle accordingly.
	//
	// TODO: implement DON sync
	DonSubscriber capabilities.DonSubscriber

	// Inner is the fallback implementation for draft/mock mode to get things running
	Inner limits.QueueLimiter[EnqueuedTriggerEvent]
}

// OCRQueue wraps a standard QueueLimiter and delegates all operations to it.
type OCRQueue struct {
	inner limits.QueueLimiter[EnqueuedTriggerEvent]
}

// NewOCRQueue creates an OCRQueue from the planned dependencies.
func NewOCRQueue(deps OCRQueueDeps) (limits.QueueLimiter[EnqueuedTriggerEvent], error) {
	if deps.Inner == nil {
		return nil, errOCRQueueNotImplemented
	}
	return &OCRQueue{inner: deps.Inner}, nil
}

// NewOCRQueueWithInnerQueue returns a constructor with the same signature as NewOCRQueue.
func NewOCRQueueWithInnerQueue(inner limits.QueueLimiter[EnqueuedTriggerEvent]) func(OCRQueueDeps) (limits.QueueLimiter[EnqueuedTriggerEvent], error) {
	return func(_ OCRQueueDeps) (limits.QueueLimiter[EnqueuedTriggerEvent], error) {
		return &OCRQueue{inner: inner}, nil
	}
}

func (q *OCRQueue) Limit(ctx context.Context) (int, error) {
	return q.inner.Limit(ctx)
}

func (q *OCRQueue) Len(ctx context.Context) (int, error) {
	return q.inner.Len(ctx)
}

func (q *OCRQueue) Put(ctx context.Context, event EnqueuedTriggerEvent) error {
	return q.inner.Put(ctx, event)
}

func (q *OCRQueue) Get(ctx context.Context) (EnqueuedTriggerEvent, error) {
	return q.inner.Get(ctx)
}

func (q *OCRQueue) Wait(ctx context.Context) (EnqueuedTriggerEvent, error) {
	return q.inner.Wait(ctx)
}

func (q *OCRQueue) Close() error {
	return q.inner.Close()
}
