package cache

import (
	"context"
	"sync"
)

type OnceCtxFunction[T any] func(ctx context.Context) (T, error)

// CallOnceOnNoError returns a new function that wraps the given function f with caching capabilities.
// If f returns an error, the result is not cached, allowing f to be retried on subsequent calls.
// Use case for that is to avoid caching an error forever in case of transient errors (e.g. flaky RPC)
func CallOnceOnNoError[T any](f OnceCtxFunction[T]) OnceCtxFunction[T] {
	var (
		mu     sync.Mutex
		value  T
		err    error
		called bool
	)

	return func(ctx context.Context) (T, error) {
		mu.Lock()
		defer mu.Unlock()

		// If the function has been called successfully before, return the cached result.
		if called && err == nil {
			return value, nil
		}

		// Call the function and cache the result only if there is no error.
		value, err = f(ctx)
		if err == nil {
			called = true
		}

		return value, err
	}
}
