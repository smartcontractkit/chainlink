package logprovider

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

// New creates a new log event provider and recoverer.
// using default values for the options.
func New(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, stateStore core.UpkeepStateReader, finalityDepth uint32, chainID *big.Int) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker()
	opts := NewOptions(int64(finalityDepth), chainID)

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

	// LogLimit is the minimum number of logs to process in a single block window.
	LogLimit uint32
	// BlockRate determines the block window for log processing.
	BlockRate uint32
}

func NewOptions(finalityDepth int64, chainID *big.Int) LogTriggersOptions {
	opts := new(LogTriggersOptions)
	opts.chainID = chainID
	opts.Defaults(finalityDepth)
	return *opts
}

// Defaults sets the default values for the options.
// NOTE: o.LookbackBlocks should be set only from within tests
func (o *LogTriggersOptions) Defaults(finalityDepth int64) {
	if o.LookbackBlocks == 0 {
		lookbackBlocks := int64(100)
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
		o.BlockRate = o.defaultBlockRate()
	}
	if o.LogLimit == 0 {
		o.LogLimit = o.defaultLogLimit()
	}
}

func (o *LogTriggersOptions) defaultBlockRate() uint32 {
	switch o.chainID.Int64() {
	case 42161, 421613, 421614: // Arbitrum, Arb Goerli, Arb Sepolia
		return 2
	default:
		return 1
	}
}

func (o *LogTriggersOptions) defaultLogLimit() uint32 {
	switch o.chainID.Int64() {
	case 1, 4, 5, 42, 11155111: // Eth, Rinkeby, Goerli, Kovan, Sepolia
		return 20
	case 10, 420, 11155420, 56, 97, 137, 80001, 80002, 43114, 43113, 8453, 84531, 84532: // Optimism, OP Goerli, OP Sepolia, BSC, BSC Test, Polygon, Mumbai, Amoy, Avax, Avax Fuji, Base, Base Goerli, Base Sepolia
		return 4
	default:
		return 1
	}
}
