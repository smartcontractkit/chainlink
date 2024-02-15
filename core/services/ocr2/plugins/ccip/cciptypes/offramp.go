package cciptypes

import (
	"context"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type OffRampReader interface {
	// EncodeExecutionReport will error if messages are not a compatible version.
	EncodeExecutionReport(report ExecReport) ([]byte, error)

	// DecodeExecutionReport will error if messages are not a compatible version.
	DecodeExecutionReport(report []byte) (ExecReport, error)

	// GetExecutionStateChangesBetweenSeqNums returns all the execution state change events for the provided message sequence numbers (inclusive).
	GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confirmations int) ([]ExecutionStateChangedWithTxMeta, error)

	Address() Address

	// ChangeConfig notifies the reader that the config has changed onchain
	ChangeConfig(onchainConfig []byte, offchainConfig []byte) (Address, Address, error)

	OffchainConfig() ExecOffchainConfig

	OnchainConfig() ExecOnchainConfig

	GasPriceEstimator() GasPriceEstimatorExec

	GetSenderNonce(ctx context.Context, sender Address) (uint64, error)

	CurrentRateLimiterState(ctx context.Context) (TokenBucketRateLimit, error)

	GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error)

	GetStaticConfig(ctx context.Context) (OffRampStaticConfig, error)

	GetSourceToDestTokensMapping(ctx context.Context) (map[Address]Address, error)

	GetTokens(ctx context.Context) (OffRampTokens, error)
}

type ExecReport struct {
	Messages          []EVM2EVMMessage
	OffchainTokenData [][][]byte
	Proofs            [][32]byte
	ProofFlagBits     *big.Int
}

type ExecutionStateChangedWithTxMeta struct {
	TxMeta
	ExecutionStateChanged
}

type ExecutionStateChanged struct {
	SequenceNumber uint64
	Finalized      bool
}

type TokenBucketRateLimit struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

// ExecOffchainConfig specifies configuration for nodes executing committed messages.
type ExecOffchainConfig struct {
	// DestOptimisticConfirmations is how many confirmations to wait for the dest chain event before we consider it
	// confirmed (optimistically, need not be finalized).
	DestOptimisticConfirmations uint32
	// BatchGasLimit is the maximum sum of user callback gas we permit in one execution report.
	BatchGasLimit uint32
	// RelativeBoostPerWaitHour indicates how much to increase (artificially) the fee paid on the source chain per hour
	// of wait time, such that eventually the fee paid is greater than the execution cost, and we’ll execute it.
	// For example: if set to 0.5, that means the fee paid is increased by 50% every hour the message has been waiting.
	RelativeBoostPerWaitHour float64
	// InflightCacheExpiry indicates how long we keep a report in the plugin cache before we expire it.
	// The caching prevents us from issuing another report while one is already in flight.
	InflightCacheExpiry config.Duration
	// RootSnoozeTime is the interval at which we check roots for executable messages.
	RootSnoozeTime config.Duration
}

type ExecOnchainConfig struct {
	PermissionLessExecutionThresholdSeconds time.Duration
}

func (c ExecOnchainConfig) Validate() error {
	if c.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("must set PermissionLessExecutionThresholdSeconds")
	}
	return nil
}

type OffRampStaticConfig struct {
	CommitStore         Address
	ChainSelector       uint64
	SourceChainSelector uint64
	OnRamp              Address
	PrevOffRamp         Address
	ArmProxy            Address
}

type OffRampTokens struct {
	DestinationTokens []Address
	SourceTokens      []Address
	DestinationPool   map[Address]Address
}

// MessageExecutionState defines the execution states of CCIP messages.
type MessageExecutionState uint8

const (
	ExecutionStateUntouched MessageExecutionState = iota
	ExecutionStateInProgress
	ExecutionStateSuccess
	ExecutionStateFailure
)
