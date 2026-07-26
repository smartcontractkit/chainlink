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
	select {
	case t.taskSemaphore <- struct{}{}:
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
