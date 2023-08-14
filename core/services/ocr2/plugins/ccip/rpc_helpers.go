package ccip

import "sync"

type LazyFunction[T any] func() (T, error)

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
