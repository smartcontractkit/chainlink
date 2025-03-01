package executable

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type task struct {
	fnCtx *context.Context
	fn    func(ctx context.Context)
}

type parallelExecutor struct {
	tasksCh  chan task
	wg       sync.WaitGroup
	stopChan services.StopChan
}

func newParallelExecutor(maxParallelTasks int) *parallelExecutor {
	executor := &parallelExecutor{
		tasksCh:  make(chan task),
		stopChan: make(services.StopChan),
		wg:       sync.WaitGroup{},
	}

	for i := 0; i < maxParallelTasks; i++ {
		executor.wg.Add(1)
		go func() {
			for t := range executor.tasksCh {
				ctxWithStop, cancel := executor.stopChan.Ctx(*t.fnCtx)
				t.fn(ctxWithStop)
				cancel()
			}
			executor.wg.Done()
		}()
	}

	return executor
}

// ExecuteTask executes a task in parallel up to the maximum allowed parallel executions.  If the maximum execute limit
// is reached, the function will block until a slot is available or the context is cancelled.
func (t *parallelExecutor) ExecuteTask(ctx context.Context, fn func(ctx context.Context)) error {
	select {
	case t.tasksCh <- task{fnCtx: &ctx, fn: fn}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *parallelExecutor) Close() {
	close(t.tasksCh)
	close(t.stopChan)
	t.wg.Wait()
}
