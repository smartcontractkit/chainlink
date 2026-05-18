package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// retryable is a helper function that retries a function until it succeeds.
//
// It will retry every `retryMs` milliseconds, up to `maxRetries` times.
//
// If `maxRetries` is 0, it will retry indefinitely.
//
// retryable will return an error in the following conditions:
//   - the context is cancelled: the error returned is the context error
//   - the retry limit has been hit: the error returned is the last error returned by `fn`
func retryable(ctx context.Context, lggr logger.Logger, retryMs int, maxRetries int, fn func() error) error {
	ticker := time.NewTicker(time.Duration(retryMs) * time.Millisecond)
	defer ticker.Stop()

	// immediately try once
	err := fn()
	if err == nil {
		return nil
	}
	retries := 0

	for {
		// if maxRetries is 0, we'll retry indefinitely
		if maxRetries > 0 {
			if retries >= maxRetries {
				lggr.Errorf("%s", err)
				return fmt.Errorf("max retries reached, aborting")
			}
		}
		lggr.Errorf("%s, retrying in %.2fs", err, float64(retryMs)/1000)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err = fn()
			if err == nil {
				return nil
			}
		}

		retries++
	}
}
