package executable

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type parallelExecutor struct {
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
		t.wg.Add(1)
		go func() {
			ctxWithStop, cancel := t.stopChan.Ctx(ctx)
			fn(ctxWithStop)
			cancel()
			<-t.taskSemaphore
			t.wg.Done()
		}()
	case <-t.stopChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (t *parallelExecutor) Close() {
	close(t.stopChan)
	t.wg.Wait()
}
