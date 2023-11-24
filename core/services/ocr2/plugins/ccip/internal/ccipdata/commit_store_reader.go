package ccipdata

import (
	"context"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	SourceFinalityDepth    uint32
	GasPriceDeviationPPB   uint32
	GasPriceHeartBeat      time.Duration
	TokenPriceDeviationPPB uint32
	TokenPriceHeartBeat    time.Duration
	InflightCacheExpiry    time.Duration
	DestFinalityDepth      uint32
}

type CommitStoreStaticConfig struct {
	ChainSelector       uint64
	SourceChainSelector uint64
	OnRamp              common.Address
	ArmProxy            common.Address
}

func NewCommitOffchainConfig(
	sourceFinalityDepth uint32,
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
	destFinalityDepth uint32,
) CommitOffchainConfig {
	return CommitOffchainConfig{
		SourceFinalityDepth:    sourceFinalityDepth,
		GasPriceDeviationPPB:   gasPriceDeviationPPB,
		GasPriceHeartBeat:      gasPriceHeartBeat,
		TokenPriceDeviationPPB: tokenPriceDeviationPPB,
		TokenPriceHeartBeat:    tokenPriceHeartBeat,
		InflightCacheExpiry:    inflightCacheExpiry,
		DestFinalityDepth:      destFinalityDepth,
	}
}

//go:generate mockery --quiet --name CommitStoreReader --filename commit_store_reader_mock.go --case=underscore
type CommitStoreReader interface {
	Closer
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
	RegisterFilters(qopts ...pg.QOpt) error
}

func NewCommitStoreReader(lggr logger.Logger, address common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (CommitStoreReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(address, ec)
	if err != nil {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOnRamp, contractType)
	}
	switch version.String() {
	case V1_0_0, V1_1_0:
		// Versions are identical
		return NewCommitStoreV1_0_0(lggr, address, ec, lp, estimator)
	case V1_2_0:
		return NewCommitStoreV1_2_0(lggr, address, ec, lp, estimator)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
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

// EncodeCommitReport is only used in tests
// TODO should remove it and update tests to use Reader interface.
func EncodeCommitReport(report CommitStoreReport) ([]byte, error) {
	commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
	return encodeCommitReportV1_2_0(abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI), report)
}

func CommitReportToEthTxMeta(typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	if typ != ccipconfig.CommitStore {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.CommitStore, typ)
	}
	switch ver.String() {
	case V1_0_0, V1_1_0:
		commitStoreABI := abihelpers.MustParseABI(commit_store_1_0_0.CommitStoreABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			commitReport, err := decodeCommitReportV1_0_0(abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI), report)
			if err != nil {
				return nil, err
			}
			return commitReportToEthTxMeta(commitReport)
		}, nil
	case V1_2_0:
		commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			commitReport, err := decodeCommitReportV1_2_0(abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI), report)
			if err != nil {
				return nil, err
			}
			return commitReportToEthTxMeta(commitReport)
		}, nil
	default:
		return nil, errors.Errorf("got unexpected version %v", ver.String())
	}
}

// CommitReportToEthTxMeta generates a txmgr.EthTxMeta from the given commit report.
// sequence numbers of the committed messages will be added to tx metadata
func commitReportToEthTxMeta(commitReport CommitStoreReport) (*txmgr.TxMeta, error) {
	n := uint64(commitReport.Interval.Max-commitReport.Interval.Min) + 1
	seqRange := make([]uint64, n)
	for i := uint64(0); i < n; i++ {
		seqRange[i] = i + commitReport.Interval.Min
	}
	return &txmgr.TxMeta{
		SeqNumbers: seqRange,
	}, nil
}
