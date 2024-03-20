package workflows

import (
	"context"
	"sync"
)

type stepRequest struct {
	stepRef string
	state   executionState
}

type queue[T any] struct {
	enqueue  chan T
	dequeue chan T
}

func (q *queue[T]) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// NOTE: Should there be a max size for the queue?
	qData := []T{}

	for {
		select {
		case <-ctx.Done():
			return
		case inc := <-q.enqueue:
			qData = append(qData, inc)
		default:
			if len(qData) > 0 {
				popped := qData[0]
				select {
				case q.dequeue <- popped:
					qData = qData[1:]
				default:
				}
			}
		}

	}
}

func (q *queue[T]) start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go q.worker(ctx, wg)
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		enqueue:  make(chan T),
		dequeue: make(chan T),
	}
}
