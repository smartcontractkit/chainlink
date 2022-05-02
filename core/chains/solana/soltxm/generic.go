package soltxm

import (
	"fmt"
	"sync"
)

type LazyLoad[T any] struct {
	f     func() (T, error)
	state T
	lock  sync.RWMutex
	once  sync.Once
}

func NewLazyLoad[T any](f func() (T, error)) *LazyLoad[T] {
	return &LazyLoad[T]{
		f: f,
	}
}

func (l *LazyLoad[T]) Get() (out T, err error) {

	// fetch only once (or whenever cleared)
	l.lock.Lock()
	l.once.Do(func() {
		l.state, err = l.f()
	})
	l.lock.Unlock()

	// if err, clear so next get will retry
	if err != nil {
		l.Reset()
	}

	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.state, err
}

func (l *LazyLoad[T]) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.once = sync.Once{}
}

// BatchSplit splits an array into an array of arrays with a maximum length
func BatchSplit[T any](list []T, max int) (out [][]T, err error) {
	if max == 0 {
		return out, fmt.Errorf("max batch length cannot be 0")
	}

	// batch list into no more than max each
	for len(list) > max {
		// assign to list: remaining after taking slice from beginning
		// append to out: max length slice from beginning of list
		list, out = list[max:], append(out, list[:max])
	}
	out = append(out, list) // append remaining to list (slice len < max)
	return out, nil
}
