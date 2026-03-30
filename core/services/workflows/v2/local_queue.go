package v2

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type SubjectQueueLimiter[T any] interface {
	limits.QueueLimiter[T]

	// Run returns when the passed context is canceled.
	Run(context.Context, ObserverFunc[T])
}

type ObserverFunc[T any] func(context.Context, T)

// localQueue is a wrapped QueueLimiter whose state is local to each independent
// node.
type localQueue[T any] struct {
	limits.QueueLimiter[T]
}

func newLocalQueue[T any](ql limits.QueueLimiter[T]) *localQueue[T] {
	return &localQueue[T]{
		QueueLimiter: ql,
	}
}

// Run executes a wait loop that waits for the latest head of the queue then
// handles the value with the passed ObserverFunc.  The event is not
// handled if the func is nil.
func (q *localQueue[T]) Run(ctx context.Context, observeFn ObserverFunc[T]) {
	for {
		qh, err := q.Wait(ctx)
		if err != nil {
			return
		}

		if observeFn != nil {
			observeFn(ctx, qh)
		}
	}
}
