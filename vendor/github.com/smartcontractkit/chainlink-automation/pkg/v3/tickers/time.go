package tickers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/smartcontractkit/chainlink-automation/internal/util"
)

type observer[T any] interface {
	Process(context.Context, Tick[T]) error
}

type getterFunc[T any] func(context.Context, time.Time) (Tick[T], error)

type timeTicker[T any] struct {
	closer util.Closer

	interval time.Duration
	observer observer[T]
	getterFn getterFunc[T]
	logger   *log.Logger
}

func NewTimeTicker[T any](interval time.Duration, observer observer[T], getterFn getterFunc[T], logger *log.Logger) *timeTicker[T] {
	t := &timeTicker[T]{
		interval: interval,
		observer: observer,
		getterFn: getterFn,
		logger:   logger,
	}

	return t
}

// Start uses the provided context for each call to the getter function with the
// configured interval as a timeout. This function blocks until Close is called
// or the parent context is cancelled.
func (t *timeTicker[T]) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	if !t.closer.Store(cancel) {
		return fmt.Errorf("already running")
	}

	t.logger.Printf("starting ticker service")
	defer t.logger.Printf("ticker service stopped")

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case tm := <-ticker.C:
			if t.getterFn == nil {
				continue
			}
			tick, err := t.getterFn(ctx, tm)
			if err != nil {
				t.logger.Printf("error fetching tick: %s", err.Error())
			}
			// observer.Process can be a heavy call taking upto ObservationProcessLimit seconds
			// so it is run in a separate goroutine to not block further ticks
			// Exploratory: Add some control to limit the number of goroutines spawned
			go func(c context.Context, t Tick[T], o observer[T], l *log.Logger) {
				if err := o.Process(c, t); err != nil {
					l.Printf("error processing observer: %s", err.Error())
				}
			}(ctx, tick, t.observer, t.logger)
		case <-ctx.Done():
			return nil
		}
	}
}

func (t *timeTicker[T]) Close() error {
	_ = t.closer.Close()
	return nil
}
