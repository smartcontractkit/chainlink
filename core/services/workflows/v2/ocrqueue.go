package v2

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

var errOCRQueueNotImplemented = errors.New("OCRQueue: draft, use NewOCRQueueWithInnerQueue for delegating implementation")

// OCRQueueDeps holds dependencies for NewOCRQueue.
// Delegate owns the oracle; OCRQueue wraps the queue the transmitter feeds.
type OCRQueueDeps struct {
	Inner limits.QueueLimiter[EnqueuedTriggerEvent]
}

// OCRQueue wraps a QueueLimiter and delegates all operations to it.
// Delegate owns the oracle; transmitter decodes reports and Puts into Inner.
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
