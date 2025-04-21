package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// RunWithRetries is a helper function that retries a function until it succeeds.
//
// It will call it immediately and then retry on failure every `retryInterval`, up to `maxRetries` times.
// If `maxRetries` is 0, it will retry indefinitely.
//
// RunWithRetries will return an error in the following conditions:
//   - the context is cancelled
//   - the retry limit has been hit
func RunWithRetries(ctx context.Context, lggr logger.Logger, retryInterval time.Duration, maxRetries int, fn func() error) error {
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	// immediately try once
	err := fn()
	if err == nil {
		return nil
	}
	retries := 0

	for {
		lggr.Errorf("error: %s, retrying in %s", err, retryInterval)

		// if maxRetries is 0, we'll retry indefinitely
		if maxRetries > 0 && retries >= maxRetries {
			msg := fmt.Sprintf("max retries (%d) reached, aborting", maxRetries)
			lggr.Error(msg)
			return errors.New(msg)
		}

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
