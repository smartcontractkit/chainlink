package v2

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type SubjectQueueLimiter[T any] interface {
	limits.QueueLimiter[T]

	// Run returns when the passed context is canceled.
	Run(context.Context)

	// Register sets a function to handle queue events with.  Only one handler
	// is allowed.
	Register(ObserverFunc[T])
}

type ObserverFunc[T any] func(context.Context, T)

// localQueue is a wrapped QueueLimiter whose state is local to each independent
// node.
type localQueue[T any] struct {
	limits.QueueLimiter[T]

	mu sync.RWMutex

	// observe handles each observed value pulled from a queue
	observe ObserverFunc[T]
}

func newLocalQueue[T any](ql limits.QueueLimiter[T]) *localQueue[T] {
	return &localQueue[T]{
		QueueLimiter: ql,
	}
}

// Run executes a wait loop that waits for the latest head of the queue then
// handles the value with the registered ObserverFunc.  The event is not
// handled if no func is registered.
func (q *localQueue[T]) Run(ctx context.Context) {
	for {
		qh, err := q.Wait(ctx)
		if err != nil {
			return
		}

		q.mu.RLock()
		defer q.mu.RUnlock()
		if q.observe != nil {
			q.observe(ctx, qh)
		}
	}
}

func (q *localQueue[T]) Register(fn ObserverFunc[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.observe = fn
}
