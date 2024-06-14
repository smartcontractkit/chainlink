package utils

import (
	"sync"
)

type LazyLoad[T any] struct {
	f     func() (T, error)
	state T
	lock  sync.Mutex
	once  sync.Once
}

func NewLazyLoad[T any](f func() (T, error)) *LazyLoad[T] {
	return &LazyLoad[T]{
		f: f,
	}
}

func (l *LazyLoad[T]) Get() (out T, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// fetch only once (or whenever cleared)
	l.once.Do(func() {
		l.state, err = l.f()
	})
	// if err, clear so next get will retry
	if err != nil {
		l.once = sync.Once{}
	}
	out = l.state
	return
}

func (l *LazyLoad[T]) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.once = sync.Once{}
}
