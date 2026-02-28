package v2

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestIsRetryableError(t *testing.T) {
	t.Run("nil error is not retryable", func(t *testing.T) {
		assert.False(t, IsRetryableError(nil))
	})

	// Non-retryable errors
	t.Run("VALIDATION_FAILED is not retryable", func(t *testing.T) {
		assert.False(t, IsRetryableError(errors.New("VALIDATION_FAILED: bad request")))
	})

	t.Run("user error is not retryable", func(t *testing.T) {
		assert.False(t, IsRetryableError(errors.New("capability execution failed with user error: contract reverted")))
	})

	t.Run("capability not found config error is not retryable", func(t *testing.T) {
		assert.False(t, IsRetryableError(errors.New("action capability not found: no such capability")))
	})

	t.Run("limit exceeded is not retryable", func(t *testing.T) {
		assert.False(t, IsRetryableError(errors.New("limit exceeded for capability calls")))
	})

	// Retryable errors
	t.Run("context deadline exceeded is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("context deadline exceeded")))
	})

	t.Run("context canceled is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("context canceled")))
	})

	t.Run("CAPABILITY_NOT_FOUND from dispatcher is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("CAPABILITY_NOT_FOUND: capability not registered yet")))
	})

	t.Run("dispatcher not ready is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("dispatcher not ready: peer connection lost")))
	})

	t.Run("request expired is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("request expired by executable client")))
	})

	t.Run("received N errors is retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("received 1 errors, last error INTERNAL : node unavailable")))
	})

	t.Run("unknown error defaults to retryable", func(t *testing.T) {
		assert.True(t, IsRetryableError(errors.New("some unknown transient error")))
	})
}

func TestExecuteWithRetry(t *testing.T) {
	lggr := logger.Test(t)
	ctx := context.Background()

	t.Run("succeeds on first attempt", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       3,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		}

		callCount := 0
		result, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "success", nil
			},
		)

		require.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 1, callCount)
	})

	t.Run("retries on transient error and succeeds", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       3,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		}

		callCount := 0
		result, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				if callCount < 3 {
					return "", errors.New("request expired by executable client")
				}
				return "recovered", nil
			},
		)

		require.NoError(t, err)
		assert.Equal(t, "recovered", result)
		assert.Equal(t, 3, callCount)
	})

	t.Run("fails immediately on non-retryable error", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       3,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		}

		callCount := 0
		_, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "", errors.New("VALIDATION_FAILED: invalid input")
			},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "VALIDATION_FAILED")
		assert.Equal(t, 1, callCount, "should not retry non-retryable errors")
	})

	t.Run("exhausts all attempts on persistent transient error", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       3,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		}

		callCount := 0
		_, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "", errors.New("dispatcher not ready: connection refused")
			},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed after 3 attempts")
		assert.Equal(t, 3, callCount)
	})

	t.Run("respects context cancellation during backoff", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       5,
			InitialBackoff:   5 * time.Second, // long backoff
			MaxBackoff:        10 * time.Second,
			BackoffMultiplier: 2.0,
		}

		cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		callCount := 0
		start := time.Now()
		_, err := ExecuteWithRetry(cancelCtx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "", errors.New("dispatcher not ready")
			},
		)

		elapsed := time.Since(start)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context")
		assert.Equal(t, 1, callCount, "should only attempt once before context cancels during backoff")
		assert.Less(t, elapsed, 2*time.Second, "should not wait for full backoff")
	})

	t.Run("backoff does not exceed MaxBackoff", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       4,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        20 * time.Millisecond,
			BackoffMultiplier: 10.0, // aggressive multiplier
		}

		callCount := 0
		start := time.Now()
		_, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "", errors.New("dispatcher not ready")
			},
		)
		elapsed := time.Since(start)

		require.Error(t, err)
		assert.Equal(t, 4, callCount)
		// With max backoff of 20ms and 3 retries, total should be well under 1s
		assert.Less(t, elapsed, 1*time.Second)
	})

	t.Run("works with single attempt config", func(t *testing.T) {
		cfg := StepRetryConfig{
			MaxAttempts:       1,
			InitialBackoff:   10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		}

		callCount := 0
		_, err := ExecuteWithRetry(ctx, lggr, cfg, "test-cap",
			func(ctx context.Context) (string, error) {
				callCount++
				return "", errors.New("dispatcher not ready")
			},
		)

		require.Error(t, err)
		assert.Equal(t, 1, callCount)
	})
}

func TestDefaultStepRetryConfig(t *testing.T) {
	cfg := DefaultStepRetryConfig()
	assert.Equal(t, 3, cfg.MaxAttempts)
	assert.Equal(t, 2*time.Second, cfg.InitialBackoff)
	assert.Equal(t, 30*time.Second, cfg.MaxBackoff)
	assert.Equal(t, 2.0, cfg.BackoffMultiplier)
}
