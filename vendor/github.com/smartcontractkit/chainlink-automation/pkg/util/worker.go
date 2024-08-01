package util

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
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

type GroupedItem[T any] struct {
	Group int
	Item  WorkItem[T]
}

type WorkerGroup[T any] struct {
	maxWorkers    int
	activeWorkers int
	workers       chan *worker[T]

	queue         *Queue[GroupedItem[T]]
	input         chan GroupedItem[T]
	chInputNotify chan struct{}
	mu            sync.Mutex
	resultData    map[int][]WorkItemResult[T]
	resultNotify  map[int]chan struct{}

	// channels used to stop processing
	chStopInputs     chan struct{}
	chStopProcessing chan struct{}
	queueClosed      atomic.Bool

	// service state management
	svcChStop services.StopChan
	once      sync.Once
}

func NewWorkerGroup[T any](workers int, queue int) *WorkerGroup[T] {
	wg := &WorkerGroup[T]{
		maxWorkers:       workers,
		workers:          make(chan *worker[T], workers),
		queue:            &Queue[GroupedItem[T]]{},
		input:            make(chan GroupedItem[T], 1),
		chInputNotify:    make(chan struct{}, 1),
		resultData:       map[int][]WorkItemResult[T]{},
		resultNotify:     map[int]chan struct{}{},
		chStopInputs:     make(chan struct{}),
		chStopProcessing: make(chan struct{}),
		svcChStop:        make(chan struct{}),
	}

	go wg.run()

	return wg
}

// Do adds a new work item onto the work queue. This function blocks until
// the work queue clears up or the context is cancelled.
func (wg *WorkerGroup[T]) Do(ctx context.Context, w WorkItem[T], group int) error {

	if ctx.Err() != nil {
		return fmt.Errorf("%w; work not added to queue", ErrContextCancelled)
	}

	if wg.queueClosed.Load() {
		return fmt.Errorf("%w; work not added to queue", ErrProcessStopped)
	}

	gi := GroupedItem[T]{
		Group: group,
		Item:  w,
	}

	wg.mu.Lock()
	if _, ok := wg.resultData[group]; !ok {
		wg.resultData[group] = make([]WorkItemResult[T], 0)
	}

	if _, ok := wg.resultNotify[group]; !ok {
		wg.resultNotify[group] = make(chan struct{}, 1)
	}
	wg.mu.Unlock()

	select {
	case wg.input <- gi:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w; work not added to queue", ErrContextCancelled)
	case <-wg.svcChStop:
		return fmt.Errorf("%w; work not added to queue", ErrProcessStopped)
	}
}

func (wg *WorkerGroup[T]) NotifyResult(group int) <-chan struct{} {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	ch, ok := wg.resultNotify[group]
	if !ok {
		// if a channel isn't found for the group, create it
		wg.resultNotify[group] = make(chan struct{}, 1)

		return wg.resultNotify[group]
	}

	return ch
}

func (wg *WorkerGroup[T]) Results(group int) []WorkItemResult[T] {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	resultData, ok := wg.resultData[group]
	if !ok {
		wg.resultData[group] = []WorkItemResult[T]{}

		return wg.resultData[group]
	}

	wg.resultData[group] = []WorkItemResult[T]{}

	// results are stored as latest first
	// switch the order to provide oldest first
	if len(resultData) > 1 {
		for i, j := 0, len(resultData)-1; i < j; i, j = i+1, j-1 {
			resultData[i], resultData[j] = resultData[j], resultData[i]
		}
	}

	return resultData
}

func (wg *WorkerGroup[T]) RemoveGroup(group int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	delete(wg.resultData, group)
	delete(wg.resultNotify, group)
}

func (wg *WorkerGroup[T]) Stop() {
	wg.once.Do(func() {
		close(wg.svcChStop)
		wg.queueClosed.Store(true)
		wg.chStopInputs <- struct{}{}
	})
}

func (wg *WorkerGroup[T]) processQueue() {
	for {
		if wg.queue.Len() == 0 {
			break
		}

		value, err := wg.queue.Pop()

		// an error from pop means there is nothing to pop
		// the length check above should protect from that, but just in case
		// this error also breaks the loop
		if err != nil {
			break
		}

		wg.doJob(value)
	}
}

func (wg *WorkerGroup[T]) runQueuing() {
	for {
		select {
		case item := <-wg.input:
			wg.queue.Add(item)

			// notify that new work item came in
			// drop if notification channel is full
			select {
			case wg.chInputNotify <- struct{}{}:
			default:
			}
		case <-wg.chStopInputs:
			wg.chStopProcessing <- struct{}{}
			return
		}
	}
}

func (wg *WorkerGroup[T]) runProcessing() {
	for {
		select {
		// watch notification channel and begin processing queue
		// when notification occurs
		case <-wg.chInputNotify:
			wg.processQueue()
		case <-wg.chStopProcessing:
			return
		}
	}
}

func (wg *WorkerGroup[T]) run() {
	// start listening on the input channel for new jobs
	go wg.runQueuing()

	// main run loop for queued jobs
	wg.runProcessing()

	// run the job queue one more time just in case some
	// new work items snuck in
	wg.processQueue()
}

func (wg *WorkerGroup[T]) doJob(item GroupedItem[T]) {
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
	go func() {
		ctx, cancel := wg.svcChStop.NewCtx()
		defer cancel()
		wkr.Do(ctx, wg.storeResult(item.Group), item.Item)
	}()
}

func (wg *WorkerGroup[T]) storeResult(group int) func(result WorkItemResult[T]) {
	return func(result WorkItemResult[T]) {
		wg.mu.Lock()
		defer wg.mu.Unlock()

		_, ok := wg.resultData[group]
		if !ok {
			wg.resultData[group] = make([]WorkItemResult[T], 0)
		}

		_, ok = wg.resultNotify[group]
		if !ok {
			wg.resultNotify[group] = make(chan struct{}, 1)
		}

		wg.resultData[group] = append([]WorkItemResult[T]{result}, wg.resultData[group]...)

		select {
		case wg.resultNotify[group] <- struct{}{}:
		default:
		}
	}
}

type JobFunc[T, K any] func(context.Context, T) (K, error)
type JobResultFunc[T any] func(T, error)

func RunJobs[T, K any](ctx context.Context, wg *WorkerGroup[T], jobs []K, jobFunc JobFunc[K, T], resFunc JobResultFunc[T]) {
	var wait sync.WaitGroup
	end := make(chan struct{}, 1)

	group := rand.Intn(1_000_000_000)

	go func(g *WorkerGroup[T], w *sync.WaitGroup, ch chan struct{}) {
		for {
			select {
			case <-g.NotifyResult(group):
				//fmt.Println("NotifyResult")
				for _, r := range g.Results(group) {
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

		if err := wg.Do(ctx, makeJobFunc(ctx, job, jobFunc), group); err != nil {
			// the makeJobFunc will exit early if the context passed to it has
			// already completed or if the worker process has been stopped
			wait.Done()
			break
		}
	}

	// wait for all results to be read
	wait.Wait()

	// clean up run group resources
	wg.RemoveGroup(group)

	// close the results reader process to clean up resources
	close(end)
}

func makeJobFunc[T, K any](jobCtx context.Context, value T, jobFunc JobFunc[T, K]) WorkItem[K] {
	return func(svcCtx context.Context) (K, error) {
		// the jobFunc should exit in the case that either the job context
		// cancels or the worker service context cancels.
		ctx, cancel := context.WithCancel(jobCtx)
		defer cancel()
		stop := context.AfterFunc(svcCtx, cancel)
		defer stop()
		return jobFunc(ctx, value)
	}
}

type Queue[T any] struct {
	mu     sync.RWMutex
	values []T
}

func (q *Queue[T]) Add(values ...T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.values = append(q.values, values...)
}

func (q *Queue[T]) Pop() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.values) == 0 {
		return getZero[T](), fmt.Errorf("no values to return")
	}

	val := q.values[0]

	if len(q.values) > 1 {
		q.values = q.values[1:]
	} else {
		q.values = []T{}
	}

	return val, nil
}

func (q *Queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.values)
}
