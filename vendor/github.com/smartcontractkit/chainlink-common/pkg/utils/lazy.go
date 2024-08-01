package utils

import (
	"sync"
)

type LazyLoad[T any] struct {
	f     func() (T, error)
	state T
	ok    bool
	lock  sync.Mutex
}

func NewLazyLoad[T any](f func() (T, error)) *LazyLoad[T] {
	return &LazyLoad[T]{
		f: f,
	}
}

func (l *LazyLoad[T]) Get() (out T, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.ok {
		return l.state, nil
	}
	l.state, err = l.f()
	l.ok = err == nil
	return l.state, err
}

func (l *LazyLoad[T]) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.ok = false
}
