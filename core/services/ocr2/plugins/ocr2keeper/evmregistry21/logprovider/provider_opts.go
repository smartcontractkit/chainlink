package logprovider

import (
	"time"

	"golang.org/x/time/rate"
)

// LogEventProviderOptions holds the options for the log event provider.
type LogEventProviderOptions struct {
	// LogRetention is the amount of time to retain logs for.
	LogRetention time.Duration
	// AllowedLogsPerBlock is the maximum number of logs allowed per block in the buffer.
	BufferMaxBlockSize int
	// LogBufferSize is the number of blocks in the buffer.
	LogBufferSize int
	// AllowedLogsPerBlock is the maximum number of logs allowed per block & upkeep in the buffer.
	AllowedLogsPerBlock int
	// LogBlocksLookback is the number of blocks to look back for logs.
	LogBlocksLookback int64
	// LookbackBuffer is the number of blocks to add as a buffer to the lookback.
	LookbackBuffer int64
	// BlockRateLimit is the rate limit for fetching logs per block.
	BlockRateLimit rate.Limit
	// BlockLimitBurst is the burst limit for fetching logs per block.
	BlockLimitBurst int
	// ReadInterval is the interval to fetch logs in the background.
	ReadInterval time.Duration
	// ReadMaxBatchSize is the max number of items in one read batch / partition.
	ReadMaxBatchSize int
	// Readers is the number of reader workers to spawn.
	Readers int
}

// Defaults sets the default values for the options.
func (o *LogEventProviderOptions) Defaults() {
	if o.LogRetention == 0 {
		o.LogRetention = 24 * time.Hour
	}
	if o.BufferMaxBlockSize == 0 {
		o.BufferMaxBlockSize = 1024
	}
	if o.AllowedLogsPerBlock == 0 {
		o.AllowedLogsPerBlock = 128
	}
	if o.LogBlocksLookback == 0 {
		o.LogBlocksLookback = 512
	}
	if o.LogBufferSize == 0 {
		o.LogBufferSize = int(o.LogBlocksLookback * 3)
	}
	if o.LookbackBuffer == 0 {
		o.LookbackBuffer = 32
	}
	if o.BlockRateLimit == 0 {
		o.BlockRateLimit = rate.Every(time.Second)
	}
	if o.BlockLimitBurst == 0 {
		o.BlockLimitBurst = int(o.LogBlocksLookback)
	}
	if o.ReadInterval == 0 {
		o.ReadInterval = time.Second
	}
	if o.ReadMaxBatchSize == 0 {
		o.ReadMaxBatchSize = 32
	}
	if o.Readers == 0 {
		o.Readers = 2
	}
}
