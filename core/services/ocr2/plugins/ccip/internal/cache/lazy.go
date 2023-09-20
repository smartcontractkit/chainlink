package cache

import "sync"

type LazyFunction[T any] func() (T, error)

// LazyFetch caches the results during the first call and then returns the cached value
// on each consecutive call.
func LazyFetch[T any](fun LazyFunction[T]) LazyFunction[T] {
	var result T
	var err error
	var once sync.Once

	return func() (T, error) {
		once.Do(func() {
			result, err = fun()
		})
		return result, err
	}
}
