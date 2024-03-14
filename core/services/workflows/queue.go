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
	in  chan T
	out chan T
}

func (q *queue[T]) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	qData := []T{}

	for {
		select {
		case <-ctx.Done():
			return
		case inc := <-q.in:
			qData = append(qData, inc)
		default:
			if len(qData) > 0 {
				popped := qData[0]
				select {
				case q.out <- popped:
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
		in:  make(chan T),
		out: make(chan T),
	}
}
