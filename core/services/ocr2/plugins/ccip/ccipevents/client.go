package ccipevents

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
)

type Event[T any] struct {
	Data T
	BlockMeta
}

type BlockMeta struct {
	BlockTimestamp time.Time
	BlockNumber    int64
}

// Client can be used to fetch CCIP related parsed on-chain events.
type Client interface {
	// GetSendRequestsGteSeqNum returns all the message send requests with sequence number greater than or equal to the provided.
	// If checkFinalityTags is set to true then confs param is ignored, the latest finalized block is used in the query.
	GetSendRequestsGteSeqNum(ctx context.Context, onRamp common.Address, seqNum uint64, checkFinalityTags bool, confs int) ([]Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested], error)

	// GetSendRequestsBetweenSeqNums returns all the message send requests in the provided sequence numbers range (inclusive).
	GetSendRequestsBetweenSeqNums(ctx context.Context, onRamp common.Address, seqNumMin, seqNumMax uint64, confs int) ([]Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested], error)

	// GetTokenPriceUpdatesCreatedAfter returns all the token price updates that happened after the provided timestamp.
	GetTokenPriceUpdatesCreatedAfter(ctx context.Context, priceRegistry common.Address, ts time.Time, confs int) ([]Event[price_registry.PriceRegistryUsdPerTokenUpdated], error)

	// GetGasPriceUpdatesCreatedAfter returns all the gas price updates that happened after the provided timestamp.
	GetGasPriceUpdatesCreatedAfter(ctx context.Context, priceRegistry common.Address, chainSelector uint64, ts time.Time, confs int) ([]Event[price_registry.PriceRegistryUsdPerUnitGasUpdated], error)

	// GetExecutionStateChangesBetweenSeqNums returns all the execution state change events for the provided message sequence numbers (inclusive).
	GetExecutionStateChangesBetweenSeqNums(ctx context.Context, offRamp common.Address, seqNumMin, seqNumMax uint64, confs int) ([]Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged], error)

	// LatestBlock returns the latest known/parsed block of the underlying implementation.
	LatestBlock(ctx context.Context) (int64, error)
}
