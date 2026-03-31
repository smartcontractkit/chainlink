package v2

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// OCRQueueDeps holds dependencies for NewOCRQueue (shared across all workflow engines on this node).
type OCRQueueDeps struct {
	Inner  limits.QueueLimiter[EnqueuedTriggerEvent]
	Buffer *ObservationBuffer[EnqueuedTriggerEvent]
}

// OCRQueue is the shared node-wide queue for the OCR trigger POC:
//   - Put appends to ObservationBuffer (inputs for OCR 3.1 Observation round)
//   - Get/Wait/Close delegate to Inner (fed after consensus; POC wiring only)
//   - Per-engine trigger handling registers an ObserverFunc via ocrTriggerSubjectQueue.Run;
//     after consensus, Transmit calls DispatchConsensusEvent which invokes the observer for that workflowID.
type OCRQueue struct {
	inner  limits.QueueLimiter[EnqueuedTriggerEvent]
	buffer *ObservationBuffer[EnqueuedTriggerEvent]

	mu        sync.RWMutex
	observers map[string]ObserverFunc[EnqueuedTriggerEvent]
}

// NewOCRQueue builds a shared OCRQueue for the node. Inner must be non-nil.
func NewOCRQueue(deps OCRQueueDeps) (*OCRQueue, error) {
	if deps.Inner == nil {
		return nil, errors.New("OCRQueue requires Inner")
	}
	if deps.Buffer == nil {
		return nil, errors.New("OCRQueue requires Buffer")
	}
	return &OCRQueue{inner: deps.Inner, buffer: deps.Buffer}, nil
}

// NewSharedOCRTriggerQueueForPOC constructs one OCRQueue + buffer using the same limit settings as NewLimiters.
func NewSharedOCRTriggerQueueForPOC(lf limits.Factory, cfgFn func(*cresettings.Workflows)) (*OCRQueue, *ObservationBuffer[EnqueuedTriggerEvent], error) {
	cfg := cresettings.Default.PerWorkflow
	if cfgFn != nil {
		cfgFn(&cfg)
	}
	inner, err := limits.MakeQueueLimiter[EnqueuedTriggerEvent](lf, cfg.TriggerEventQueueLimit)
	if err != nil {
		return nil, nil, err
	}
	lamport := &LamportCounter{}
	buffer := NewObservationBuffer[EnqueuedTriggerEvent](lamport)
	q, err := NewOCRQueue(OCRQueueDeps{Inner: inner, Buffer: buffer})
	if err != nil {
		return nil, nil, err
	}
	return q, buffer, nil
}

func (q *OCRQueue) Limit(ctx context.Context) (int, error) {
	return q.inner.Limit(ctx)
}

func (q *OCRQueue) Len(ctx context.Context) (int, error) {
	return q.inner.Len(ctx)
}

func (q *OCRQueue) Put(ctx context.Context, event EnqueuedTriggerEvent) error {
	q.buffer.Add(event)
	return nil
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

// TriggerSubjectQueueForWorkflow implements PerWorkflowTriggerSubjectQueue.
func (q *OCRQueue) TriggerSubjectQueueForWorkflow(workflowID string) SubjectQueueLimiter[EnqueuedTriggerEvent] {
	return NewOCRTriggerSubjectQueue(q, workflowID)
}

// RegisterObserver stores the callback used when DispatchConsensusEvent receives an event for wid.
func (q *OCRQueue) RegisterObserver(wid string, fn ObserverFunc[EnqueuedTriggerEvent]) {
	if wid == "" {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.observers == nil {
		q.observers = make(map[string]ObserverFunc[EnqueuedTriggerEvent])
	}
	q.observers[wid] = fn
}

// UnregisterObserver removes the callback for wid.
func (q *OCRQueue) UnregisterObserver(wid string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.observers, wid)
}

// DispatchConsensusEvent delivers a post-consensus event to the registered observer for ev.WorkflowID().
func (q *OCRQueue) DispatchConsensusEvent(ctx context.Context, ev EnqueuedTriggerEvent) {
	q.mu.RLock()
	fn := q.observers[ev.WorkflowID()]
	q.mu.RUnlock()
	if fn != nil {
		fn(ctx, ev)
	}
}

// Buffer returns the observation buffer used by the OCR 3.1 plugin factory.
func (q *OCRQueue) Buffer() *ObservationBuffer[EnqueuedTriggerEvent] {
	return q.buffer
}

// ocrTriggerSubjectQueue wraps a shared OCRQueue for one workflow: Run registers the observer; queue ops delegate.
type ocrTriggerSubjectQueue struct {
	shared     *OCRQueue
	workflowID string
}

// NewOCRTriggerSubjectQueue returns a SubjectQueueLimiter that shares one OCRQueue but binds Run to workflowID.
func NewOCRTriggerSubjectQueue(shared *OCRQueue, workflowID string) SubjectQueueLimiter[EnqueuedTriggerEvent] {
	return &ocrTriggerSubjectQueue{shared: shared, workflowID: workflowID}
}

func (q *ocrTriggerSubjectQueue) Limit(ctx context.Context) (int, error) {
	return q.shared.Limit(ctx)
}

func (q *ocrTriggerSubjectQueue) Len(ctx context.Context) (int, error) {
	return q.shared.Len(ctx)
}

func (q *ocrTriggerSubjectQueue) Put(ctx context.Context, event EnqueuedTriggerEvent) error {
	return q.shared.Put(ctx, event)
}

func (q *ocrTriggerSubjectQueue) Get(ctx context.Context) (EnqueuedTriggerEvent, error) {
	return q.shared.Get(ctx)
}

func (q *ocrTriggerSubjectQueue) Wait(ctx context.Context) (EnqueuedTriggerEvent, error) {
	return q.shared.Wait(ctx)
}

func (q *ocrTriggerSubjectQueue) Close() error {
	return nil
}

func (q *ocrTriggerSubjectQueue) Run(ctx context.Context, observeFn ObserverFunc[EnqueuedTriggerEvent]) {
	q.shared.RegisterObserver(q.workflowID, observeFn)
	defer q.shared.UnregisterObserver(q.workflowID)
	<-ctx.Done()
}

func (q *ocrTriggerSubjectQueue) EvictTenant(tenant string) error {
	if tenant == q.workflowID {
		q.shared.UnregisterObserver(tenant)
	}

	return limits.TryEvictTenant(q.shared, tenant)
}
