package worker

import (
	"context"
	"fmt"

	"github.com/alitto/pond/v2"
)

type Pool struct {
	pool   pond.Pool
	ctx    context.Context //nolint:containedctx // Pool needs to manage a shared context for all tasks
	cancel context.CancelCauseFunc
}

func New(ctx context.Context, maxConcurrency int, opts ...pond.Option) *Pool {
	poolCtx, cancel := context.WithCancelCause(ctx)
	return &Pool{
		pool:   pond.NewPool(maxConcurrency, opts...),
		ctx:    poolCtx,
		cancel: cancel,
	}
}

func (p *Pool) StopAndWait() { p.pool.StopAndWait() }

type FutureAny struct {
	ch <-chan struct {
		value any
		err   error
	}
}

func (p *Pool) SubmitErr(fn func(context.Context) error) FutureAny {
	return p.SubmitAny(func(ctx context.Context) (any, error) {
		return nil, fn(ctx)
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

func (p *Pool) SubmitAny(fn func(context.Context) (any, error)) FutureAny {
	ch := make(chan struct {
		value any
		err   error
	}, 1)

	task := p.pool.Submit(func() {
		// If there was no panic, execute the function and send the result to the channel (value or error)
		value, err := fn(p.ctx)
		ch <- struct {
			value any
			err   error
		}{value, err}

		// Cancel pool context if there was an error
		if err != nil {
			p.cancel(err)
		}
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

			// Cancel pool context on panic
			p.cancel(taskErr)
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
