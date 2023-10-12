package ccipdata

import (
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

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

//go:generate mockery --quiet --name CommitStoreReader --output . --filename commit_store_reader_mock.go --inpackage --case=underscore
type CommitStoreReader interface {
	Closer
	GetExpectedNextSequenceNumber(context context.Context) (uint64, error)
	GetLatestPriceEpochAndRound(context context.Context) (uint64, error)
	// GetAcceptedCommitReportsGteSeqNum returns all the accepted commit reports that have sequence number greater than or equal to the provided.
	GetAcceptedCommitReportsGteSeqNum(ctx context.Context, seqNum uint64, confs int) ([]Event[CommitStoreReport], error)
	// GetAcceptedCommitReportsGteTimestamp returns all the commit reports with timestamp greater than or equal to the provided.
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
}

func NewCommitStoreReader(lggr logger.Logger, address common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (CommitStoreReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(address, ec)
	if err != nil {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOnRamp, contractType)
	}
	switch version.String() {
	case v1_0_0, v1_1_0:
		return NewCommitStoreV1_0_0(lggr, address, ec, lp, estimator)
	case v1_2_0:
		return NewCommitStoreV1_2_0(lggr, address, ec, lp, estimator)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
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
	case v1_0_0, v1_1_0:
		commitStoreABI := abihelpers.MustParseABI(commit_store_1_0_0.CommitStoreABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			commitReport, err := decodeCommitReportV1_0_0(abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI), report)
			if err != nil {
				return nil, err
			}
			return commitReportToEthTxMeta(commitReport)
		}, nil
	case v1_2_0:
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
