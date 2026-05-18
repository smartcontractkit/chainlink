package workflows

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRetryableZeroMaxRetries(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	fn := func() error {
		return errors.New("test error")
	}

	err := retryable(ctx, logger.NullLogger, 100, 0, fn)
	assert.ErrorIs(t, err, context.DeadlineExceeded, "Expected context deadline exceeded error")
}

func TestRetryableSuccessOnFirstAttempt(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fn := func() error {
		return nil
	}

	err := retryable(ctx, logger.NullLogger, 100, 3, fn)
	require.NoError(t, err, "Expected no error as function succeeds on first attempt")
}

func TestRetryableSuccessAfterRetries(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retries := 0
	fn := func() error {
		if retries < 2 {
			retries++
			return errors.New("test error")
		}
		return nil
	}

	err := retryable(ctx, logger.NullLogger, 100, 5, fn)
	assert.NoError(t, err, "Expected no error after successful retry")
	assert.Equal(t, 2, retries, "Expected two retries before success")
}

func TestRetryableErrorOnFirstTryNoRetries(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fn := func() error {
		return errors.New("immediate failure")
	}

	err := retryable(ctx, logger.NullLogger, 100, 1, fn)
	require.Error(t, err, "Expected an error on the first try with no retries allowed")
	assert.Equal(t, "max retries reached, aborting", err.Error(), "Expected function to abort after the first try")
}

func TestRetryableErrorAfterMultipleRetries(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	attempts := 0
	fn := func() error {
		attempts++
		return errors.New("persistent error")
	}

	maxRetries := 3
	err := retryable(ctx, logger.NullLogger, 100, maxRetries, fn)
	require.Error(t, err, "Expected an error after multiple retries")
	assert.Equal(t, "max retries reached, aborting", err.Error(), "Expected the max retries reached error message")
	assert.Equal(t, maxRetries+1, attempts, "Expected the function to be executed retry + 1 times")
}

func TestRetryableCancellationHandling(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	fn := func() error {
		return errors.New("test error")
	}

	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	err := retryable(ctx, logger.NullLogger, 100, 5, fn)
	assert.ErrorIs(t, err, context.Canceled, "Expected context cancellation error")
}
