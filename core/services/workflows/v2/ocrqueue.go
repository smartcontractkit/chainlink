package v2

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

var errOCRQueueNotImplemented = errors.New("OCRQueue: draft, use NewOCRQueueWithInnerQueue for delegating implementation")

// OCRQueueDeps holds dependencies for NewOCRQueue.
// Delegate owns the oracle; OCRQueue wraps the queue the transmitter feeds.
type OCRQueueDeps[T Identifiable] struct {
	Inner  limits.QueueLimiter[T]
	Buffer *ObservationBuffer[T]
}

// OCRQueue wraps a QueueLimiter. Put buffers to ObservationBuffer (feeds Observation).
// Get/Wait read from Inner (fed by transmitter when consensus events arrive).
type OCRQueue[T Identifiable] struct {
	inner  limits.QueueLimiter[T]
	buffer *ObservationBuffer[T]
}

// NewOCRQueue creates an OCRQueue from the planned dependencies.
func NewOCRQueue[T Identifiable](deps OCRQueueDeps[T]) (limits.QueueLimiter[T], error) {
	if deps.Inner == nil {
		return nil, errOCRQueueNotImplemented
	}
	if deps.Buffer == nil {
		return nil, errors.New("OCRQueue requires Buffer")
	}
	return &OCRQueue[T]{inner: deps.Inner, buffer: deps.Buffer}, nil
}

// NewOCRQueueWithInnerQueue returns a constructor with the same signature as NewOCRQueue.
func NewOCRQueueWithInnerQueue[T Identifiable](inner limits.QueueLimiter[T], buffer *ObservationBuffer[T]) func(OCRQueueDeps[T]) (limits.QueueLimiter[T], error) {
	return func(_ OCRQueueDeps[T]) (limits.QueueLimiter[T], error) {
		return &OCRQueue[T]{inner: inner, buffer: buffer}, nil
	}
}

func (q *OCRQueue[T]) Limit(ctx context.Context) (int, error) {
	return q.inner.Limit(ctx)
}

func (q *OCRQueue[T]) Len(ctx context.Context) (int, error) {
	return q.inner.Len(ctx)
}

func (q *OCRQueue[T]) Put(ctx context.Context, event T) error {
	q.buffer.Add(event)
	return nil
}

func (q *OCRQueue[T]) Get(ctx context.Context) (T, error) {
	return q.inner.Get(ctx)
}

func (q *OCRQueue[T]) Wait(ctx context.Context) (T, error) {
	return q.inner.Wait(ctx)
}

func (q *OCRQueue[T]) Close() error {
	return q.inner.Close()
}
