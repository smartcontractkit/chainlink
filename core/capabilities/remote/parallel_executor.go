package remote

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// ErrNoSlotAvailable is returned by TryExecuteTask when all executor slots are in use.
var ErrNoSlotAvailable = errors.New("no parallel executor slot available")

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

// NewParallelExecutor creates an executor that allows at most maxParallelTasks in-flight tasks.
func NewParallelExecutor(maxParallelTasks int, name string, attrs ...attribute.KeyValue) *ParallelExecutor {
	attrs = append(attrs, attribute.String("parallel_executor_name", name))
	e := &ParallelExecutor{
		stopChan:       make(services.StopChan),
		wg:             sync.WaitGroup{},
		taskSemaphore:  make(chan struct{}, maxParallelTasks),
		name:           name,
		slotUsageAttrs: attrs,
	}
	return e
}

// recordSlotUsage publishes current slot utilization to the configured gauge, if any.
func (t *ParallelExecutor) recordSlotUsage(ctx context.Context) {
	if t.slotUsageGauge == nil {
		return
	}
	maxSlots := t.maxSlots()
	if maxSlots == 0 {
		return
	}
	usage := float64(t.occupiedSlots()) / float64(maxSlots)
	t.slotUsageGauge.Record(ctx, usage, metric.WithAttributes(t.slotUsageAttrs...))
}

// occupiedSlots returns the number of in-flight tasks holding executor slots.
func (t *ParallelExecutor) occupiedSlots() int {
	return len(t.taskSemaphore)
}

// maxSlots returns the maximum number of concurrent tasks allowed.
func (t *ParallelExecutor) maxSlots() int {
	return cap(t.taskSemaphore)
}

// ExecuteTask executes a task in parallel up to the maximum allowed parallel executions. If the
// maximum execute limit is reached, the function will block until a slot is available or the
// context is cancelled.
func (t *ParallelExecutor) ExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	return t.ExecuteTaskWithWaitContext(ctx, ctx, fn)
}

// ExecuteTaskWithWaitContext is like ExecuteTask, but waits for a free slot only until waitCtx is
// done, while the task itself runs with ctx. Cancelling waitCtx after the task started has no effect on it.
func (t *ParallelExecutor) ExecuteTaskWithWaitContext(waitCtx, ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.taskSemaphore <- struct{}{}:
		return t.runInAcquiredSlot(ctx, fn)
	case <-t.stopChan:
		return nil
	case <-waitCtx.Done():
		return waitCtx.Err()
	}
}

// TryExecuteTask is a non-blocking variant of ExecuteTask. If all slots are in use, it returns
// ErrNoSlotAvailable immediately instead of waiting for a slot to free up.
func (t *ParallelExecutor) TryExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.taskSemaphore <- struct{}{}:
		return t.runInAcquiredSlot(ctx, fn)
	default:
		return ErrNoSlotAvailable
	}
}

// runInAcquiredSlot runs fn in a new goroutine that releases the slot when done. The caller
// must have already acquired a slot.
func (t *ParallelExecutor) runInAcquiredSlot(ctx context.Context, fn func(ctx context.Context)) error {
	t.recordSlotUsage(ctx)
	stopped := !t.IfNotStopped(func() {
		t.wg.Go(func() {
			ctxWithStop, cancel := t.stopChan.Ctx(ctx)
			defer func() {
				<-t.taskSemaphore
				t.recordSlotUsage(ctxWithStop)
				cancel()
			}()
			fn(ctxWithStop)
		})
	})

	if stopped {
		<-t.taskSemaphore
		return errors.New("executor stopped")
	}
	return nil
}

// Start starts the executor.
func (t *ParallelExecutor) Start(ctx context.Context) error {
	var err error
	t.slotUsageGauge, err = beholder.GetMeter().Float64Gauge("platform_parallel_executor_slot_usage")
	if err != nil {
		return fmt.Errorf("failed to register platform_parallel_executor_slot_usage: %w", err)
	}
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
