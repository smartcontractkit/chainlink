package ccipdata

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

type CommitStoreInterval struct {
	Min, Max uint64
}

type CommitStoreReport struct {
	TokenPrices []TokenPrice
	GasPrices   []GasPrice
	Interval    CommitStoreInterval
	MerkleRoot  [32]byte
}

// Common to all versions
type CommitOnchainConfig commit_store.CommitStoreDynamicConfig

func (d CommitOnchainConfig) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "priceRegistry", "type": "address"}
			],
			"type": "tuple"
		}
	]`
}

func (d CommitOnchainConfig) Validate() error {
	if d.PriceRegistry == (common.Address{}) {
		return errors.New("must set Price Registry address")
	}
	return nil
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
	OnRamp              common.Address
	ArmProxy            common.Address
}

func NewCommitOffchainConfig(
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
) CommitOffchainConfig {
	return CommitOffchainConfig{
		GasPriceDeviationPPB:   gasPriceDeviationPPB,
		GasPriceHeartBeat:      gasPriceHeartBeat,
		TokenPriceDeviationPPB: tokenPriceDeviationPPB,
		TokenPriceHeartBeat:    tokenPriceHeartBeat,
		InflightCacheExpiry:    inflightCacheExpiry,
	}
}

//go:generate mockery --quiet --name CommitStoreReader --filename commit_store_reader_mock.go --case=underscore
type CommitStoreReader interface {
	GetExpectedNextSequenceNumber(context context.Context) (uint64, error)
	GetLatestPriceEpochAndRound(context context.Context) (uint64, error)
	// GetCommitReportMatchingSeqNum returns accepted commit report that satisfies Interval.Min <= seqNum <= Interval.Max. Returned slice should be empty or have exactly one element
	GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confs int) ([]Event[CommitStoreReport], error)
	// GetAcceptedCommitReportsGteTimestamp returns all the commit reports with timestamp greater than or equal to the provided.
	// Returned Commit Reports have to be sorted by Interval.Min/Interval.Max in ascending order.
	GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]Event[CommitStoreReport], error)
	IsDown(ctx context.Context) (bool, error)
	IsBlessed(ctx context.Context, root [32]byte) (bool, error)
	// Notifies the reader that the config has changed onchain
	ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error)
	OffchainConfig() CommitOffchainConfig
	GasPriceEstimator() prices.GasPriceEstimatorCommit
	EncodeCommitReport(report CommitStoreReport) ([]byte, error)
	DecodeCommitReport(report []byte) (CommitStoreReport, error)
	VerifyExecutionReport(ctx context.Context, report ExecReport) (bool, error)
	GetCommitStoreStaticConfig(ctx context.Context) (CommitStoreStaticConfig, error)
}

// FetchCommitStoreStaticConfig provides access to a commitStore's static config, which is required to access the source chain ID.
func FetchCommitStoreStaticConfig(address common.Address, ec client.Client) (commit_store.CommitStoreStaticConfig, error) {
	commitStore, err := loadCommitStore(address, ec)
	if err != nil {
		return commit_store.CommitStoreStaticConfig{}, err
	}
	return commitStore.GetStaticConfig(&bind.CallOpts{})
}

func loadCommitStore(commitStoreAddress common.Address, client client.Client) (commit_store.CommitStoreInterface, error) {
	_, err := ccipconfig.VerifyTypeAndVersion(commitStoreAddress, client, ccipconfig.CommitStore)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid commitStore contract")
	}
	return commit_store.NewCommitStore(commitStoreAddress, client)
}
