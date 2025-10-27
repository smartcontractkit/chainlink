package worker

import (
	"context"
	"fmt"

	"github.com/alitto/pond/v2"
)

type Pool struct {
	pool pond.Pool
}

func New(maxConcurrency int, opts ...pond.Option) *Pool {
	return &Pool{
		pool: pond.NewPool(maxConcurrency, opts...),
	}
}

func (p *Pool) StopAndWait() { p.pool.StopAndWait() }

type FutureAny struct {
	ch <-chan struct {
		value any
		err   error
	}
}

func (p *Pool) SubmitErr(fn func() error) FutureAny {
	return p.SubmitAny(func() (any, error) {
		return nil, fn()
	})
}

func AwaitErr(ctx context.Context, future FutureAny) error {
	select {
	case result := <-future.ch:
		return result.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pool) SubmitAny(fn func() (any, error)) FutureAny {
	ch := make(chan struct {
		value any
		err   error
	}, 1)

	task := p.pool.Submit(func() {
		// If there was no panic, execute the function and send the result to the channel (value or error)
		value, err := fn()
		ch <- struct {
			value any
			err   error
		}{value, err}
	})

	// Monitor the task for panics that pond caught
	go func() {
		defer close(ch)
		if taskErr := task.Wait(); taskErr != nil {
			// If there was a panic, pond caught it and the task function never sent to the channel
			// So we send the panic as error instead
			ch <- struct {
				value any
				err   error
			}{nil, taskErr}
		}
	}()

	return FutureAny{ch: ch}
}

func AwaitAs[T any](ctx context.Context, f FutureAny) (T, error) {
	var zero T
	select {
	case result := <-f.ch:
		if result.err != nil {
			return zero, result.err
		}
		value, ok := result.value.(T)
		if !ok {
			return zero, fmt.Errorf("type mismatch. Expected: %T Got: %T", zero, result.value)
		}

		return value, nil
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}
