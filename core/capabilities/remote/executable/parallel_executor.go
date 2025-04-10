package executable

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type parallelExecutor struct {
	services.StateMachine
	wg       sync.WaitGroup
	stopChan services.StopChan

	taskSemaphore chan struct{}
}

func newParallelExecutor(maxParallelTasks int) *parallelExecutor {
	executor := &parallelExecutor{
		stopChan:      make(services.StopChan),
		wg:            sync.WaitGroup{},
		taskSemaphore: make(chan struct{}, maxParallelTasks),
	}

	return executor
}

// ExecuteTask executes a task in parallel up to the maximum allowed parallel executions.  If the maximum execute limit
// is reached, the function will block until a slot is available or the context is cancelled.
func (t *parallelExecutor) ExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.taskSemaphore <- struct{}{}:
		stopped := !t.IfNotStopped(func() {
			t.wg.Add(1)
			go func() {
				ctxWithStop, cancel := t.stopChan.Ctx(ctx)
				fn(ctxWithStop)
				cancel()
				<-t.taskSemaphore
				t.wg.Done()
			}()
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

func (t *parallelExecutor) Start(ctx context.Context) error {
	return t.StartOnce(t.Name(), func() error {
		return nil
	})
}

func (t *parallelExecutor) Close() error {
	return t.StopOnce(t.Name(), func() error {
		close(t.stopChan)
		t.wg.Wait()
		return nil
	})
}

func (t *parallelExecutor) Name() string {
	return "ParallelExecutor"
}
