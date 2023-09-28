package ccipdata

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
)

type Event[T any] struct {
	Data T
	Meta
}

type Meta struct {
	BlockTimestamp time.Time
	BlockNumber    int64
	TxHash         common.Hash
	LogIndex       uint
}

// Client can be used to fetch CCIP related parsed on-chain data.
//
//go:generate mockery --quiet --name Reader --output . --filename reader_mock.go --inpackage --case=underscore
type Reader interface {
	// GetTokenPriceUpdatesCreatedAfter returns all the token price updates that happened after the provided timestamp.
	GetTokenPriceUpdatesCreatedAfter(ctx context.Context, priceRegistry common.Address, ts time.Time, confs int) ([]Event[price_registry.PriceRegistryUsdPerTokenUpdated], error)

	// GetGasPriceUpdatesCreatedAfter returns all the gas price updates that happened after the provided timestamp.
	GetGasPriceUpdatesCreatedAfter(ctx context.Context, priceRegistry common.Address, chainSelector uint64, ts time.Time, confs int) ([]Event[price_registry.PriceRegistryUsdPerUnitGasUpdated], error)

	// GetExecutionStateChangesBetweenSeqNums returns all the execution state change events for the provided message sequence numbers (inclusive).
	GetExecutionStateChangesBetweenSeqNums(ctx context.Context, offRamp common.Address, seqNumMin, seqNumMax uint64, confs int) ([]Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged], error)

	// GetAcceptedCommitReportsGteSeqNum returns all the accepted commit reports that have sequence number greater than or equal to the provided.
	GetAcceptedCommitReportsGteSeqNum(ctx context.Context, commitStoreAddress common.Address, seqNum uint64, confs int) ([]Event[commit_store.CommitStoreReportAccepted], error)

	// GetAcceptedCommitReportsGteTimestamp returns all the commit reports with timestamp greater than or equal to the provided.
	GetAcceptedCommitReportsGteTimestamp(ctx context.Context, commitStoreAddress common.Address, ts time.Time, confs int) ([]Event[commit_store.CommitStoreReportAccepted], error)

	// LatestBlock returns the latest known/parsed block of the underlying implementation.
	LatestBlock(ctx context.Context) (int64, error)
}
