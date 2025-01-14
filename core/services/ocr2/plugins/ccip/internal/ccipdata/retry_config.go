package ccipdata

import "time"

// RetryConfig configures an initial delay between retries, a max delay between retries, and a maximum number of
// times to retry
type RetryConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	MaxRetries   uint
}
