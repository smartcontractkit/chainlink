package v2

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var errOCRQueueNotImplemented = errors.New("OCRQueue: draft, use NewOCRQueueWithInnerQueue for delegating implementation")

// DONInfo holds DON membership info (members, F, bootstrap peers) for OCR.
// Placeholder for design; concrete type TBD when wiring from handler.
type DONInfo struct{}

// OCR3_1PluginFactory creates OCR 3.1 ReportingPlugins (Blob + KV).
// Placeholder for design; use triggerqueue.Factory when implemented.
type OCR3_1PluginFactory interface{}

// OCRQueueDeps holds the planned dependencies for NewOCRQueue from the design.
// See ocr-trigger-queue.md. Draft: DonInfo and PluginFactory may be nil when using
// Inner as fallback.
type OCRQueueDeps struct {
	Lf            limits.Factory
	Cfg           *cresettings.Workflows
	DonInfo       *DONInfo
	PluginFactory OCR3_1PluginFactory
	Lggr          logger.Logger

	// Inner is an optional fallback queue for the draft. When set, NewOCRQueue returns
	// an OCRQueue that delegates to it. When nil, NewOCRQueue returns errOCRQueueNotImplemented.
	Inner limits.QueueLimiter[EnqueuedTriggerEvent]
}

// OCRQueue wraps a standard QueueLimiter and delegates all operations to it.
// Draft: created via NewOCRQueue with deps.Inner set; functionally identical to the standard queue.
// Future: will run OCR internally and feed from report callbacks.
type OCRQueue struct {
	inner limits.QueueLimiter[EnqueuedTriggerEvent]
}

// NewOCRQueue creates an OCRQueue from the planned dependencies.
// Draft: when deps.Inner is set, returns an OCRQueue delegating to it. Otherwise returns
// errOCRQueueNotImplemented. The delegate calls this with Inner from NewStandardTriggerQueue.
func NewOCRQueue(deps OCRQueueDeps) (limits.QueueLimiter[EnqueuedTriggerEvent], error) {
	if deps.Inner != nil {
		return &OCRQueue{inner: deps.Inner}, nil
	}
	return nil, errOCRQueueNotImplemented
}

// NewOCRQueueWithInnerQueue returns a constructor with the same signature as NewOCRQueue.
// The returned function ignores deps and returns an OCRQueue that delegates to inner.
// Use this for the draft to match current behavior (delegate to standard queue).
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
