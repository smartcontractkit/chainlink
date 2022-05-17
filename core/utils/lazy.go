package utils

import (
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
