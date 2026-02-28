package v2

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// StepRetryConfig configures retry behavior for capability step execution.
type StepRetryConfig struct {
	// MaxAttempts is the maximum number of attempts (including the first).
	MaxAttempts int
	// InitialBackoff is the delay before the first retry.
	InitialBackoff time.Duration
	// MaxBackoff is the maximum delay between retries.
	MaxBackoff time.Duration
	// BackoffMultiplier is the factor by which backoff increases each retry.
	BackoffMultiplier float64
}

// DefaultStepRetryConfig returns the default retry configuration.
func DefaultStepRetryConfig() StepRetryConfig {
	return StepRetryConfig{
		MaxAttempts:       3,
		InitialBackoff:   2 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// IsRetryableError determines if an error from a capability execution is transient
// and should be retried. Non-retryable errors (validation failures, application-level
// contract reverts) should fail immediately.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()

	// Non-retryable: validation errors
	if strings.Contains(errMsg, "VALIDATION_FAILED") {
		return false
	}
	// Non-retryable: user/application errors from the capability
	if strings.Contains(errMsg, "user error") {
		return false
	}
	// Non-retryable: config errors
	if strings.Contains(errMsg, "not found") && strings.Contains(errMsg, "capability") {
		return false
	}
	// Non-retryable: limit exceeded errors
	if strings.Contains(errMsg, "limit") && strings.Contains(errMsg, "exceeded") {
		return false
	}

	// Retryable: DON unreachable / timeout / transient network issues
	if strings.Contains(errMsg, "context deadline exceeded") {
		return true
	}
	if strings.Contains(errMsg, "context canceled") {
		return true
	}
	if strings.Contains(errMsg, "CAPABILITY_NOT_FOUND") {
		// During brief restart window, capability may not be registered yet
		return true
	}
	if strings.Contains(errMsg, "dispatcher not ready") {
		return true
	}
	if strings.Contains(errMsg, "request expired") {
		return true
	}
	if strings.Contains(errMsg, "received") && strings.Contains(errMsg, "errors") {
		// "received N errors" from client_request.go
		return true
	}

	// Default: treat unknown errors as retryable (conservative for F=0)
	return true
}

// ExecuteWithRetry wraps a capability execution function with retry logic.
// It retries transient errors with exponential backoff.
func ExecuteWithRetry[T any](
	ctx context.Context,
	lggr logger.Logger,
	cfg StepRetryConfig,
	capabilityID string,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var lastErr error
	var zero T

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		result, err := fn(ctx)
		if err == nil {
			if attempt > 1 {
				lggr.Infow("Capability execution succeeded after retry",
					"capabilityID", capabilityID,
					"attempt", attempt)
			}
			return result, nil
		}

		lastErr = err

		// Don't retry non-retryable errors
		if !IsRetryableError(err) {
			lggr.Debugw("Capability error is non-retryable, failing immediately",
				"capabilityID", capabilityID,
				"attempt", attempt,
				"err", err)
			return zero, err
		}

		// Don't retry if we've exhausted attempts
		if attempt >= cfg.MaxAttempts {
			lggr.Warnw("Capability execution failed after all retry attempts",
				"capabilityID", capabilityID,
				"attempts", attempt,
				"lastErr", err)
			break
		}

		// Calculate backoff
		backoff := time.Duration(float64(cfg.InitialBackoff) * math.Pow(cfg.BackoffMultiplier, float64(attempt-1)))
		if backoff > cfg.MaxBackoff {
			backoff = cfg.MaxBackoff
		}

		lggr.Warnw("Capability execution failed, retrying",
			"capabilityID", capabilityID,
			"attempt", attempt,
			"maxAttempts", cfg.MaxAttempts,
			"backoff", backoff,
			"err", err)

		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled during retry backoff: %w", errors.Join(ctx.Err(), lastErr))
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return zero, fmt.Errorf("capability execution failed after %d attempts: %w", cfg.MaxAttempts, lastErr)
}
