package concurrency

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// NoTaskType is a dummy type to be used when no task type is needed
type NoTaskType struct{}

type ConcurrentExecutorOpt[ResultType any, ResultChannelType ChannelWithResult[ResultType], TaskType any] func(c *ConcurrentExecutor[ResultType, ResultChannelType, TaskType])

// ConcurrentExecutor is a utility to execute tasks concurrently
type ConcurrentExecutor[ResultType any, ResultChannelType ChannelWithResult[ResultType], TaskType any] struct {
	results  []ResultType
	errors   []error
	logger   zerolog.Logger
	failFast bool
	context  context.Context //nolint: containedctx // won't fix for now
}

// NewConcurrentExecutor creates a new ConcurrentExecutor
func NewConcurrentExecutor[ResultType any, ResultChannelType ChannelWithResult[ResultType], TaskType any](logger zerolog.Logger, opts ...ConcurrentExecutorOpt[ResultType, ResultChannelType, TaskType]) *ConcurrentExecutor[ResultType, ResultChannelType, TaskType] {
	c := &ConcurrentExecutor[ResultType, ResultChannelType, TaskType]{
		logger:   logger,
		results:  []ResultType{},
		errors:   []error{},
		failFast: true,
		context:  context.Background(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// / WithContext sets the context for the executor, if not set it defaults to context.Background()
func WithContext[ResultType any, ResultChannelType ChannelWithResult[ResultType], TaskType any](context context.Context) ConcurrentExecutorOpt[ResultType, ResultChannelType, TaskType] {
	return func(c *ConcurrentExecutor[ResultType, ResultChannelType, TaskType]) {
		c.context = context
	}
}

// WithoutFailFast disables fail fast. Executor will wait for all tasks to finish even if some of them fail.
func WithoutFailFast[ResultType any, ResultChannelType ChannelWithResult[ResultType], TaskType any]() ConcurrentExecutorOpt[ResultType, ResultChannelType, TaskType] {
	return func(c *ConcurrentExecutor[ResultType, ResultChannelType, TaskType]) {
		c.failFast = false
	}
}

// TaskProcessorFn is a function that processes a task that requires a payload. It should send the result to the resultCh and any error to the errorCh. It should
// never send to both channels. The executorNum is the index of the executor that is processing the task. The payload is the task's payload. If task doesn't require
// one use SimpleTaskProcessorFn instead.
type TaskProcessorFn[ResultChannelType, TaskType any] func(resultCh chan ResultChannelType, errorCh chan error, executorNum int, payload TaskType)

// SimpleTaskProcessorFn is a function that processes a task that doesn't require a payload. It should send the result to the resultCh and any error to the errorCh. It should
// never send to both channels. The executorNum is the index of the executor that is processing the task.
type SimpleTaskProcessorFn[ResultChannelType any] func(resultCh chan ResultChannelType, errorCh chan error, executorNum int)

// ChannelWithResult is an interface that should be implemented by the result channel
type ChannelWithResult[ResultType any] interface {
	GetResult() ResultType
}

// GetErrors returns all errors that occurred during processing
func (e *ConcurrentExecutor[ResultType, ResultChannelType, TaskType]) GetErrors() []error {
	return e.errors
}

// ExecuteSimple executes a task that doesn't require a payload. It is executed repeatTimes times with given concurrency. The simpleProcessorFn is the function that processes the task.
// Executor will attempt to distribute the tasks evenly among the executors.
func (e *ConcurrentExecutor[ResultType, ResultChannelType, TaskType]) ExecuteSimple(concurrency int, repeatTimes int, simpleProcessorFn SimpleTaskProcessorFn[ResultChannelType]) ([]ResultType, error) {
	dummy := make([]TaskType, repeatTimes)
	for i := 0; i < repeatTimes; i++ {
		dummy[i] = *new(TaskType)
	}

	return e.Execute(concurrency, dummy, adaptSimpleToTaskProcessorFn[ResultChannelType, TaskType](simpleProcessorFn))
}

// Execute executes a task that requires a payload. It is executed with given concurrency. The processorFn is the function that processes the task.
// Executor will attempt to distribute the tasks evenly among the executors.
func (e *ConcurrentExecutor[ResultType, ResultChannelType, TaskType]) Execute(concurrency int, payload []TaskType, processorFn TaskProcessorFn[ResultChannelType, TaskType]) ([]ResultType, error) {
	if len(payload) == 0 {
		return []ResultType{}, nil
	}

	if concurrency <= 0 {
		e.logger.Warn().Msg("Concurrency is less than 1, setting it to 1")
		concurrency = 1
	}

	var wgProcesses sync.WaitGroup
	wgProcesses.Add(len(payload))

	canSafelyContinueCh := make(chan struct{}) // waits until listening goroutine finishes, so we can safely return from the function
	doneProcessingCh := make(chan struct{})    // signals that both result and error channels are closed
	errorCh := make(chan error, len(payload))
	resultCh := make(chan ResultChannelType, len(payload))

	// mutex to protect shared state
	mutex := sync.Mutex{}

	// atomic counter to keep track of processed tasks
	var atomicCounter atomic.Int32

	ctx, cancel := context.WithCancel(e.context)

	// listen in the background until all tasks are processed (no fail-fast)
	go func() {
		defer func() {
			e.logger.Trace().Msg("Finished listening to task processing results")
			close(canSafelyContinueCh)
		}()
		for {
			select {
			case err, ok := <-errorCh:
				if !ok {
					e.logger.Trace().Msg("Error channel closed")
					return
				}
				if err != nil {
					mutex.Lock()
					e.errors = append(e.errors, err)
					e.logger.Err(err).Msg("Error processing a task")
					mutex.Unlock()
					wgProcesses.Done()

					// cancel the context if failFast is enabled and it hasn't been cancelled yet
					if e.failFast && ctx.Err() == nil {
						cancel()
					}
				}
			case result, ok := <-resultCh:
				if !ok {
					e.logger.Trace().Msg("Result channel closed")
					return
				}

				counter := atomicCounter.Add(1)
				mutex.Lock()
				e.results = append(e.results, result.GetResult())
				e.logger.Trace().Str("Done/Total", fmt.Sprintf("%d/%d", counter, len(payload))).Msg("Finished aggregating task result")
				mutex.Unlock()
				wgProcesses.Done()
			case <-doneProcessingCh:
				e.logger.Trace().Msg("Signaling that processing is done")
				return
			}
		}
	}()

	dividedPayload := DivideSlice(payload, concurrency)

	for executorNum := 0; executorNum < concurrency; executorNum++ {
		go func(key int) {
			payloads := dividedPayload[key]

			if len(payloads) == 0 {
				return
			}

			e.logger.Debug().
				Int("Executor Index", key).
				Int("Tasks to process", len(payloads)).
				Msg("Started processing tasks")

			for i := range payloads {
				// if context is cancelled and failFast is enabled  mark all remaining tasks as finished
				if e.failFast && ctx.Err() != nil {
					e.logger.Trace().
						Int("Executor Index", key).
						Str("Cancelled/Total", fmt.Sprintf("%d/%d", (i+1), len(payloads))).
						Msg("Canelling remaining tasks")
					wgProcesses.Done()

					continue
				}

				processorFn(resultCh, errorCh, key, payloads[i])
				e.logger.Trace().
					Int("Executor Index", key).
					Str("Done/Total", fmt.Sprintf("%d/%d", (i+1), len(payloads))).
					Msg("Processed a tasks")
			}

			e.logger.Debug().
				Int("Executor Index", key).
				Msg("Finished processing tasks")
		}(executorNum)
	}

	wgProcesses.Wait()
	close(resultCh)
	close(errorCh)
	close(doneProcessingCh)
	<-canSafelyContinueCh

	if len(e.errors) > 0 {
		return []ResultType{}, fmt.Errorf("failed to process %d task(s). Check logs for more details. Look for 'Error processing a task'", len(e.errors))
	}

	return e.results, nil
}

func adaptSimpleToTaskProcessorFn[ResultChannelType any, TaskType any](
	simpleFn SimpleTaskProcessorFn[ResultChannelType],
) TaskProcessorFn[ResultChannelType, TaskType] {
	return func(resultCh chan ResultChannelType, errorCh chan error, executorNum int, _ TaskType) {
		simpleFn(resultCh, errorCh, executorNum)
	}
}
