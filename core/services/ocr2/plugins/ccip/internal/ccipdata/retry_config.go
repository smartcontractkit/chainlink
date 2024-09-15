package ccipdata

import "time"

// RetryConfig configures an initial delay between retries and a max delay between retries
type RetryConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
}
