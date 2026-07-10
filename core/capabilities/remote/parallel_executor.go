package remote

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// ParallelExecutor runs tasks concurrently up to a configured limit.
type ParallelExecutor struct {
	services.StateMachine
	wg       sync.WaitGroup
	stopChan services.StopChan

	taskSemaphore  chan struct{}
	name           string
	slotUsageGauge metric.Float64Gauge
	slotUsageAttrs []attribute.KeyValue
}

// ParallelExecutorOption configures a ParallelExecutor.
type ParallelExecutorOption func(*ParallelExecutor)

// WithExecutorName sets the service name returned by Name().
func WithExecutorName(name string) ParallelExecutorOption {
	return func(e *ParallelExecutor) {
		e.name = name
	}
}

// WithSlotUsageMetric records occupied/max slot ratio on the gauge when tasks start and finish.
func WithSlotUsageMetric(gauge metric.Float64Gauge, attrs ...attribute.KeyValue) ParallelExecutorOption {
	return func(e *ParallelExecutor) {
		e.slotUsageGauge = gauge
		e.slotUsageAttrs = attrs
	}
}

// NewParallelExecutor creates an executor that allows at most maxParallelTasks in-flight tasks.
func NewParallelExecutor(maxParallelTasks int, opts ...ParallelExecutorOption) *ParallelExecutor {
	e := &ParallelExecutor{
		stopChan:      make(services.StopChan),
		wg:            sync.WaitGroup{},
		taskSemaphore: make(chan struct{}, maxParallelTasks),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// SetSlotUsageGauge binds the gauge used by RecordSlotUsage. Attributes are set via WithSlotUsageMetric.
func (t *ParallelExecutor) SetSlotUsageGauge(gauge metric.Float64Gauge) {
	t.slotUsageGauge = gauge
}

// RecordSlotUsage publishes current slot utilization to the configured gauge, if any.
func (t *ParallelExecutor) RecordSlotUsage(ctx context.Context) {
	if t.slotUsageGauge == nil {
		return
	}
	maxSlots := t.MaxSlots()
	if maxSlots == 0 {
		return
	}
	usage := float64(t.OccupiedSlots()) / float64(maxSlots)
	t.slotUsageGauge.Record(ctx, usage, metric.WithAttributes(t.slotUsageAttrs...))
}

// OccupiedSlots returns the number of in-flight tasks holding executor slots.
func (t *ParallelExecutor) OccupiedSlots() int {
	return len(t.taskSemaphore)
}

// MaxSlots returns the maximum number of concurrent tasks allowed.
func (t *ParallelExecutor) MaxSlots() int {
	return cap(t.taskSemaphore)
}

// ExecuteTask executes a task in parallel up to the maximum allowed parallel executions. If the
// maximum execute limit is reached, the function will block until a slot is available or the
// context is cancelled.
func (t *ParallelExecutor) ExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.taskSemaphore <- struct{}{}:
		t.RecordSlotUsage(ctx)
		stopped := !t.IfNotStopped(func() {
			t.wg.Go(func() {
				ctxWithStop, cancel := t.stopChan.Ctx(ctx)
				defer func() {
					cancel()
					<-t.taskSemaphore
					t.RecordSlotUsage(ctxWithStop)
				}()
				fn(ctxWithStop)
			})
		})

		if stopped {
			return errors.New("executor stopped")
		}
	case <-t.stopChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// Start starts the executor.
func (t *ParallelExecutor) Start(ctx context.Context) error {
	return t.StartOnce(t.Name(), func() error {
		return nil
	})
}

// Close stops the executor and waits for in-flight tasks to finish.
func (t *ParallelExecutor) Close() error {
	return t.StopOnce(t.Name(), func() error {
		close(t.stopChan)
		t.wg.Wait()
		return nil
	})
}

// Name returns the service name.
func (t *ParallelExecutor) Name() string {
	if t.name != "" {
		return t.name
	}
	return "ParallelExecutor"
}
