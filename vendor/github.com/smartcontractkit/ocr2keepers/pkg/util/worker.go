package util

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrProcessStopped   = fmt.Errorf("worker process has stopped")
	ErrContextCancelled = fmt.Errorf("worker context cancelled")
)

type WorkItemResult[T any] struct {
	Worker string
	Data   T
	Err    error
	Time   time.Duration
}

type WorkItem[T any] func(context.Context) (T, error)

type worker[T any] struct {
	Name  string
	Queue chan *worker[T]
}

func (w *worker[T]) Do(ctx context.Context, r func(WorkItemResult[T]), wrk WorkItem[T]) {
	start := time.Now()

	var data T
	var err error

	if ctx.Err() != nil {
		err = ctx.Err()
	} else {
		data, err = wrk(ctx)
	}

	r(WorkItemResult[T]{
		Worker: w.Name,
		Data:   data,
		Err:    err,
		Time:   time.Since(start),
	})

	// put itself back on the queue when done
	select {
	case w.Queue <- w:
	default:
	}
}

type WorkerGroup[T any] struct {
	maxWorkers    int
	activeWorkers int
	workers       chan *worker[T]
	queue         chan WorkItem[T]
	queuedItems   atomic.Int64
	queueClosed   atomic.Bool
	drain         chan struct{}
	svcCtx        context.Context
	svcCancel     context.CancelFunc
	resultData    []WorkItemResult[T]
	resultNotify  chan struct{}
	resultLen     atomic.Int64
	mu            sync.RWMutex
	once          sync.Once
}

func NewWorkerGroup[T any](workers int, queue int) *WorkerGroup[T] {
	svcCtx, svcCancel := context.WithCancel(context.Background())
	wg := &WorkerGroup[T]{
		maxWorkers:   workers,
		workers:      make(chan *worker[T], workers),
		queue:        make(chan WorkItem[T], queue),
		drain:        make(chan struct{}, 1),
		resultData:   make([]WorkItemResult[T], 0),
		resultNotify: make(chan struct{}, 1),
		svcCtx:       svcCtx,
		svcCancel:    svcCancel,
	}

	go wg.run()

	runtime.SetFinalizer(wg, func(g *WorkerGroup[T]) { g.Stop() })

	return wg
}

// Do adds a new work item onto the work queue. This function blocks until
// the work queue clears up or the context is cancelled.
func (wg *WorkerGroup[T]) Do(ctx context.Context, w WorkItem[T]) error {

	if ctx.Err() != nil {
		return fmt.Errorf("%w; work not added to queue", ErrContextCancelled)
	}

	if wg.queueClosed.Load() {
		return fmt.Errorf("%w; work not added to queue", ErrProcessStopped)
	}

	select {
	case wg.queue <- w:
		wg.queuedItems.Add(1)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w; work not added to queue", ErrContextCancelled)
	case <-wg.svcCtx.Done():
		return fmt.Errorf("%w; work not added to queue", ErrProcessStopped)
	}
}

func (wg *WorkerGroup[T]) NotifyResult() <-chan struct{} {
	return wg.resultNotify
}

func (wg *WorkerGroup[T]) Results() []WorkItemResult[T] {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	resultData := wg.resultData
	wg.resultData = nil
	wg.resultLen.Store(0)
	for i, j := 0, len(resultData)-1; i < j; i, j = i+1, j-1 {
		resultData[i], resultData[j] = resultData[j], resultData[i]
	}
	return resultData
}

func (wg *WorkerGroup[T]) Stop() {
	wg.once.Do(func() {
		wg.svcCancel()
		wg.queueClosed.Store(true)
		wg.drain <- struct{}{}
	})
}

func (wg *WorkerGroup[T]) run() {
	// main run loop for queued jobs
	{
	Runner:
		for {
			select {
			case item := <-wg.queue:
				wg.queuedItems.Add(-1)
				wg.doJob(item)
			case <-wg.drain:
				// if drain is called, cancel the service context
				// and break from the run loop
				break Runner
			}
		}
	}

	if wg.queuedItems.Load() == 0 {
		return
	}

	// drain the job queue before terminating the run process
	{
	Drainer:
		for item := range wg.queue {
			if wg.queuedItems.Load() == 0 {
				break Drainer
			}

			wg.queuedItems.Add(-1)
			wg.doJob(item)
		}
	}
}

func (wg *WorkerGroup[T]) doJob(item WorkItem[T]) {
	var wkr *worker[T]

	// no read or write locks on activeWorkers or maxWorkers because it's
	// assumed the job loop is a single process reading from the job queue
	if wg.activeWorkers < wg.maxWorkers {
		// create a new worker
		wkr = &worker[T]{
			Name:  fmt.Sprintf("worker-%d", wg.activeWorkers+1),
			Queue: wg.workers,
		}
		wg.activeWorkers++
	} else {
		// wait for a worker to be available
		wkr = <-wg.workers
	}

	// have worker do the work
	go wkr.Do(wg.svcCtx, wg.storeResult, item)
}

func (wg *WorkerGroup[T]) storeResult(result WorkItemResult[T]) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	wg.resultData = append([]WorkItemResult[T]{result}, wg.resultData...)
	wg.resultLen.Add(1)

	select {
	case wg.resultNotify <- struct{}{}:
	default:
	}
}

type JobFunc[T, K any] func(context.Context, T) (K, error)
type JobResultFunc[T any] func(T, error)

func RunJobs[T, K any](ctx context.Context, wg *WorkerGroup[T], jobs []K, jobFunc JobFunc[K, T], resFunc JobResultFunc[T]) {
	var wait sync.WaitGroup
	end := make(chan struct{}, 1)

	go func(g *WorkerGroup[T], w *sync.WaitGroup, ch chan struct{}) {
		for {
			select {
			case <-g.NotifyResult():
				for _, r := range g.Results() {
					resFunc(r.Data, r.Err)
					w.Done()
				}
			case <-ch:
				return
			}
		}
	}(wg, &wait, end)

	for _, job := range jobs {
		wait.Add(1)

		if err := wg.Do(ctx, makeJobFunc(ctx, job, jobFunc)); err != nil {
			// the makeJobFunc will exit early if the context passed to it has
			// already completed or if the worker process has been stopped
			wait.Done()
			break
		}
	}

	// wait for all results to be read
	wait.Wait()

	// close the results reader process to clean up resources
	close(end)
}

func makeJobFunc[T, K any](jobCtx context.Context, value T, jobFunc JobFunc[T, K]) WorkItem[K] {
	return func(svcCtx context.Context) (K, error) {
		// the jobFunc should exit in the case that either the job context
		// cancels or the worker service context cancels. To ensure we don't end
		// up with memory leaks, cancel the merged context to release resources.
		ctx, cancel := MergeContextsWithCancel(svcCtx, jobCtx)
		defer cancel()
		return jobFunc(ctx, value)
	}
}
