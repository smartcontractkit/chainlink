package utils

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func NewRealClock() Clock {
	return &realClock{}
}

func (realClock) Now() time.Time {
	return time.Now()
}

type fixedClock struct {
	now time.Time
}

func NewFixedClock(now time.Time) Clock {
	return &fixedClock{now: now}
}

func (fc fixedClock) Now() time.Time {
	return fc.now
}
