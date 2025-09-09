package v2

import "context"

type semaphore[T any] struct {
	sem chan struct{}
}

func NewSemaphore[T any](capacity uint16) *semaphore[T] {
	return &semaphore[T]{
		sem: make(chan struct{}, capacity),
	}
}

func (s *semaphore[T]) WhenAcquired(ctx context.Context, doFn func() (T, error)) (T, error) {
	select {
	case s.sem <- struct{}{}:
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}

	defer func() { <-s.sem }()
	return doFn()
}

func (s *semaphore[T]) Len() int {
	return len(s.sem)
}
