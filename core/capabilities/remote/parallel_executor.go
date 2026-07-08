package remote

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// ParallelExecutor runs tasks concurrently up to a configured limit.
type ParallelExecutor struct {
	services.StateMachine
	wg       sync.WaitGroup
	stopChan services.StopChan

	taskSemaphore chan struct{}
}

// NewParallelExecutor creates an executor that allows at most maxParallelTasks in-flight tasks.
func NewParallelExecutor(maxParallelTasks int) *ParallelExecutor {
	return &ParallelExecutor{
		stopChan:      make(services.StopChan),
		wg:            sync.WaitGroup{},
		taskSemaphore: make(chan struct{}, maxParallelTasks),
	}
}

// ExecuteTask executes a task in parallel up to the maximum allowed parallel executions. If the
// maximum execute limit is reached, the function will block until a slot is available or the
// context is cancelled.
func (t *ParallelExecutor) ExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.taskSemaphore <- struct{}{}:
		stopped := !t.IfNotStopped(func() {
			t.wg.Go(func() {
				ctxWithStop, cancel := t.stopChan.Ctx(ctx)
				fn(ctxWithStop)
				cancel()
				<-t.taskSemaphore
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
	return "ParallelExecutor"
}
