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
func New(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, stateStore core.UpkeepStateReader, opts LogTriggersOptions) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker()
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
	// AllowedLogsPerUpkeep is the maximum number of logs allowed per upkeep every single call.
	AllowedLogsPerUpkeep int
	// MaxPayloads is the maximum number of payloads to return per call.
	MaxPayloads int
	// MaxLogsPerBlock is the maximum number of blocks in the buffer.
	MaxLogsPerBlock int
	// MaxLogsPerUpkeepInBlock is the maximum number of logs allowed per upkeep in a block.
	MaxLogsPerUpkeepInBlock int
	// MaxProposals is the maximum number of proposals that can be returned by GetRecoveryProposals
	MaxProposals int
}

func NewOptions(finalityDepth int64) LogTriggersOptions {
	opts := new(LogTriggersOptions)
	opts.FinalityDepth = finalityDepth
	opts.assignDefaults()
	return *opts
}

// assignDefaults sets the default values for the options.
// NOTE: o.LookbackBlocks should be set only from within tests
func (o *LogTriggersOptions) assignDefaults() {
	if o.LookbackBlocks == 0 {
		lookbackBlocks := int64(200)
		if lookbackBlocks < o.FinalityDepth { // TODO the order of assigning lookbackBlocks vs FinalityDepth is fickle
			lookbackBlocks = o.FinalityDepth
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
	if o.AllowedLogsPerUpkeep == 0 {
		o.AllowedLogsPerUpkeep = 5
	}
	if o.MaxPayloads == 0 {
		o.MaxPayloads = 100
	}
	if o.MaxLogsPerBlock == 0 {
		o.MaxLogsPerBlock = 1024
	}
	if o.MaxLogsPerUpkeepInBlock == 0 {
		o.MaxLogsPerUpkeepInBlock = 32
	}
	if o.MaxProposals == 0 {
		o.MaxProposals = 20
	}
}
