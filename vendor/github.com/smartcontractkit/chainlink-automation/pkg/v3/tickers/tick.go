package tickers

import (
	"context"
)

// Tick is the container for the individual tick
type Tick[T any] interface {
	// Value provides data scoped to the tick
	Value(ctx context.Context) (T, error)
}
