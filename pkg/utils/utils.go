package utils

import (
	"context"
	"fmt"
	"math"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// WithJitter adds +/- 10% to a duration
func WithJitter(d time.Duration) time.Duration {
	// #nosec
	if d == 0 {
		return 0
	}
	// ensure non-zero arg to Intn to avoid panic
	max := math.Max(float64(d.Abs())/5.0, 1.)
	// #nosec - non critical randomness
	jitter := mrand.Intn(int(max))
	jitter = jitter - (jitter / 2)
	return time.Duration(int(d) + jitter)
}

// ContextFromChan creates a context that finishes when the provided channel
// receives or is closed.
// When channel closes, the ctx.Err() will always be context.Canceled
// NOTE: Spins up a goroutine that exits on cancellation.
// REMEMBER TO CALL CANCEL OTHERWISE IT CAN LEAD TO MEMORY BCF-3067 LEAKS
func ContextFromChan(chStop <-chan struct{}) (context.Context, context.CancelFunc) {
	return services.StopRChan(chStop).NewCtx()
}

// ContextWithDeadlineFn returns a copy of the parent context with the deadline modified by deadlineFn.
// deadlineFn will only be called if the parent has a deadline.
// The new deadline must be sooner than the old to have an effect.
func ContextWithDeadlineFn(ctx context.Context, deadlineFn func(orig time.Time) time.Time) (context.Context, context.CancelFunc) {
	cancel := func() {}
	if d, ok := ctx.Deadline(); ok {
		if m := deadlineFn(d); m.Before(d) {
			ctx, cancel = context.WithDeadline(ctx, m)
		}
	}
	return ctx, cancel
}

func IsZero[C comparable](val C) bool {
	var zero C
	return zero == val
}

// JustError takes a tuple and returns the last entry, the error.
func JustError(_ interface{}, err error) error {
	return err
}

// WrapIfError decorates an error with the given message.  It is intended to
// be used with `defer` statements, like so:
//
//	func SomeFunction() (err error) {
//	    defer WrapIfError(&err, "error in SomeFunction:")
//
//	    ...
//	}
func WrapIfError(err *error, msg string) {
	if *err != nil {
		*err = fmt.Errorf("%s: %w", msg, *err)
	}
}

// AllEqual returns true iff all the provided elements are equal to each other.
func AllEqual[T comparable](elems ...T) bool {
	for i := 1; i < len(elems); i++ {
		if elems[i] != elems[0] {
			return false
		}
	}
	return true
}

// WaitGroupChan creates a channel that closes when the provided sync.WaitGroup is done.
func WaitGroupChan(wg *sync.WaitGroup) <-chan struct{} {
	chAwait := make(chan struct{})
	go func() {
		defer close(chAwait)
		wg.Wait()
	}()
	return chAwait
}

// DependentAwaiter contains Dependent funcs
type DependentAwaiter interface {
	AwaitDependents() <-chan struct{}
	AddDependents(n int)
	DependentReady()
}

type dependentAwaiter struct {
	wg *sync.WaitGroup
	ch <-chan struct{}
}

// NewDependentAwaiter creates a new DependentAwaiter
func NewDependentAwaiter() DependentAwaiter {
	return &dependentAwaiter{
		wg: &sync.WaitGroup{},
	}
}

func (da *dependentAwaiter) AwaitDependents() <-chan struct{} {
	if da.ch == nil {
		da.ch = WaitGroupChan(da.wg)
	}
	return da.ch
}

func (da *dependentAwaiter) AddDependents(n int) {
	da.wg.Add(n)
}

func (da *dependentAwaiter) DependentReady() {
	da.wg.Done()
}
