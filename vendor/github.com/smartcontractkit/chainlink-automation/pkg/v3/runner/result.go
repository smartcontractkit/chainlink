package runner

import (
	"sync"
)

type result[T any] struct {
	// this struct type isn't expressly defined to run in a single thread or
	// multiple threads so internally a mutex provides the thread safety
	// guarantees in the case it is used in a multi-threaded way
	mu        sync.RWMutex
	successes int
	failures  int
	err       error
	values    []T
}

func newResult[T any]() *result[T] {
	return &result[T]{
		values: make([]T, 0),
	}
}

func (r *result[T]) Successes() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.successes
}

func (r *result[T]) AddSuccesses(v int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.successes += v
}

func (r *result[T]) Failures() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.failures
}

func (r *result[T]) AddFailures(v int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.failures += v
}

func (r *result[T]) Err() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.err
}

func (r *result[T]) SetErr(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.err = err
}

func (r *result[T]) Total() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.successes + r.failures
}

func (r *result[T]) unsafeTotal() int {
	return r.successes + r.failures
}

func (r *result[T]) SuccessRate() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.unsafeTotal() == 0 {
		return 0
	}

	return float64(r.successes) / float64(r.unsafeTotal())
}

func (r *result[T]) FailureRate() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.unsafeTotal() == 0 {
		return 0
	}

	return float64(r.failures) / float64(r.unsafeTotal())
}

func (r *result[T]) Add(res T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.values = append(r.values, res)
}

func (r *result[T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.values
}
