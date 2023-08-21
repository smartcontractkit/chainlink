package logprovider

import (
	"time"

	"golang.org/x/time/rate"
)

// LogEventProviderOptions holds the options for the log event provider.
type LogEventProviderOptions struct {
	// LookbackBlocks is the number of blocks to look back for logs.
	LookbackBlocks int64
	// ReorgBuffer is the number of blocks to add as a buffer to the lookback.
	ReorgBuffer int64
	// BlockRateLimit is the rate limit on the range of blocks the we fetch logs for.
	BlockRateLimit rate.Limit
	// BlockLimitBurst is the burst upper limit on the range of blocks the we fetch logs for.
	BlockLimitBurst int
	// ReadInterval is the interval to fetch logs in the background.
	ReadInterval time.Duration
	// ReadBatchSize is the max number of items in one read batch / partition.
	ReadBatchSize int
	// Readers is the number of reader workers to spawn.
	Readers int
}

// Defaults sets the default values for the options.
func (o *LogEventProviderOptions) Defaults() {
	if o.LookbackBlocks == 0 {
		// TODO: Ensure lookback blocks is at least as large as Finality Depth to
		// ensure recoverer does not go beyond finality depth
		o.LookbackBlocks = 200
	}
	if o.ReorgBuffer == 0 {
		o.ReorgBuffer = 32
	}
	if o.BlockRateLimit == 0 {
		o.BlockRateLimit = rate.Every(time.Second)
	}
	if o.BlockLimitBurst == 0 {
		o.BlockLimitBurst = int(o.LookbackBlocks)
	}
	if o.ReadInterval == 0 {
		o.ReadInterval = time.Second
	}
	if o.ReadBatchSize == 0 {
		o.ReadBatchSize = 32
	}
	if o.Readers == 0 {
		o.Readers = 4
	}
}
