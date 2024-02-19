package cciptypes

import (
	"context"
	"math/big"
	"time"
)

type CommitStoreReader interface {
	GetExpectedNextSequenceNumber(context context.Context) (uint64, error)

	GetLatestPriceEpochAndRound(context context.Context) (uint64, error)

	// GetCommitReportMatchingSeqNum returns accepted commit report that satisfies Interval.Min <= seqNum <= Interval.Max. Returned slice should be empty or have exactly one element
	GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]CommitStoreReportWithTxMeta, error)

	// GetAcceptedCommitReportsGteTimestamp returns all the commit reports with timestamp greater than or equal to the provided.
	// Returned Commit Reports have to be sorted by Interval.Min/Interval.Max in ascending order.
	GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]CommitStoreReportWithTxMeta, error)

	IsDown(ctx context.Context) (bool, error)

	IsBlessed(ctx context.Context, root [32]byte) (bool, error)

	// ChangeConfig notifies the reader that the config has changed onchain
	ChangeConfig(onchainConfig []byte, offchainConfig []byte) (Address, error)

	OffchainConfig() CommitOffchainConfig

	GasPriceEstimator() GasPriceEstimatorCommit

	EncodeCommitReport(report CommitStoreReport) ([]byte, error)

	DecodeCommitReport(report []byte) (CommitStoreReport, error)

	VerifyExecutionReport(ctx context.Context, report ExecReport) (bool, error)

	GetCommitStoreStaticConfig(ctx context.Context) (CommitStoreStaticConfig, error)
}

type CommitStoreReportWithTxMeta struct {
	TxMeta
	CommitStoreReport
}

type CommitStoreReport struct {
	TokenPrices []TokenPrice
	GasPrices   []GasPrice
	Interval    CommitStoreInterval
	MerkleRoot  [32]byte
}

type TokenPrice struct {
	Token Address
	Value *big.Int
}

type GasPrice struct {
	DestChainSelector uint64
	Value             *big.Int
}

type CommitOffchainConfig struct {
	GasPriceDeviationPPB   uint32
	GasPriceHeartBeat      time.Duration
	TokenPriceDeviationPPB uint32
	TokenPriceHeartBeat    time.Duration
	InflightCacheExpiry    time.Duration
}

type CommitStoreStaticConfig struct {
	ChainSelector       uint64
	SourceChainSelector uint64
	OnRamp              Address
	ArmProxy            Address
}

type CommitStoreInterval struct {
	Min uint64
	Max uint64
}
