package logprovider

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

// New creates a new log event provider and recoverer.
// using default values for the options.
func New(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, stateStore core.UpkeepStateReader, finalityDepth uint32) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker()
	opts := NewOptions(int64(finalityDepth))
	provider := NewLogProvider(lggr, poller, packer, filterStore, opts)
	recoverer := NewLogRecoverer(lggr, poller, c, stateStore, packer, filterStore, opts)

	return provider, recoverer
}

// LogTriggersOptions holds the options for the log trigger components.
type LogTriggersOptions struct {
	// LookbackBlocks is the number of blocks the provider will look back for logs.
	// The recoverer will scan for logs up to this depth.
	// NOTE: MUST be set to a greater-or-equal to the chain's finality depth.
	LookbackBlocks int64
	// ReadInterval is the interval to fetch logs in the background.
	ReadInterval time.Duration
	// BlockRateLimit is the rate limit on the range of blocks the we fetch logs for.
	BlockRateLimit rate.Limit
	// blockLimitBurst is the burst upper limit on the range of blocks the we fetch logs for.
	BlockLimitBurst int
	// Finality depth is the number of blocks to wait before considering a block final.
	FinalityDepth int64
}

func NewOptions(finalityDepth int64) LogTriggersOptions {
	opts := new(LogTriggersOptions)
	opts.Defaults(finalityDepth)
	return *opts
}

// Defaults sets the default values for the options.
// NOTE: o.LookbackBlocks should be set only from within tests
func (o *LogTriggersOptions) Defaults(finalityDepth int64) {
	if o.LookbackBlocks == 0 {
		lookbackBlocks := int64(200)
		if lookbackBlocks < finalityDepth {
			lookbackBlocks = finalityDepth
		}
		o.LookbackBlocks = lookbackBlocks
	}
	if o.ReadInterval == 0 {
		o.ReadInterval = time.Second
	}
	if o.BlockLimitBurst == 0 {
		o.BlockLimitBurst = int(o.LookbackBlocks)
	}
	if o.BlockRateLimit == 0 {
		o.BlockRateLimit = rate.Every(o.ReadInterval)
	}
	if o.FinalityDepth == 0 {
		o.FinalityDepth = finalityDepth
	}
}
