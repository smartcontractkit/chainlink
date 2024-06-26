package logprovider

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

// New creates a new log event provider and recoverer.
// using default values for the options.
func New(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, stateStore core.UpkeepStateReader, finalityDepth uint32, chainID *big.Int, blockRate, logLimit uint32) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker()
	opts := NewOptions(int64(finalityDepth), blockRate, logLimit, chainID)

	provider := NewLogProvider(lggr, poller, chainID, packer, filterStore, opts)
	recoverer := NewLogRecoverer(lggr, poller, c, stateStore, packer, filterStore, opts)

	return provider, recoverer
}

// LogTriggersOptions holds the options for the log trigger components.
type LogTriggersOptions struct {
	chainID *big.Int
	// LookbackBlocks is the number of blocks the provider will look back for logs.
	// The recoverer will scan for logs up to this depth.
	// NOTE: MUST be set to a greater-or-equal to the chain's finality depth.
	LookbackBlocks int64
	// ReadInterval is the interval to fetch logs in the background.
	ReadInterval time.Duration
	// Finality depth is the number of blocks to wait before considering a block final.
	FinalityDepth int64

	// TODO: (AUTO-9355) remove once we have a single version
	BufferVersion BufferVersion
	// LogLimit is the minimum number of logs to process in a single block window.
	LogLimit uint32
	// BlockRate determines the block window for log processing.
	BlockRate uint32
}

// BufferVersion is the version of the log buffer.
// TODO: (AUTO-9355) remove once we have a single version
type BufferVersion string

const (
	BufferVersionDefault BufferVersion = ""
	BufferVersionV1      BufferVersion = "v1"
)

func NewOptions(finalityDepth int64, blockRate, logLimit uint32, chainID *big.Int) LogTriggersOptions {
	opts := new(LogTriggersOptions)
	opts.chainID = chainID
	opts.Defaults(finalityDepth, blockRate, logLimit)
	return *opts
}

// Defaults sets the default values for the options.
// NOTE: o.LookbackBlocks should be set only from within tests
func (o *LogTriggersOptions) Defaults(finalityDepth int64, blockRate, logLimit uint32) {
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
	if o.FinalityDepth == 0 {
		o.FinalityDepth = finalityDepth
	}
	if o.BlockRate == 0 {
		o.BlockRate = blockRate
	}
	if o.LogLimit == 0 {
		o.LogLimit = logLimit
	}
}
