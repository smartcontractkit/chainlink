package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

var _ cciptypes.CommitStoreReader = (*IncompleteSourceCommitStoreReader)(nil)
var _ cciptypes.CommitStoreReader = (*IncompleteDestCommitStoreReader)(nil)

// IncompleteSourceCommitStoreReader is an implementation of CommitStoreReader with the only valid methods being
// GasPriceEstimator, ChangeConfig, and OffchainConfig
type IncompleteSourceCommitStoreReader struct {
	estimator          gas.EvmFeeEstimator
	gasPriceEstimator  *prices.DAGasPriceEstimator
	sourceMaxGasPrice  *big.Int
	offchainConfig     cciptypes.CommitOffchainConfig
	feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider
}

func NewIncompleteSourceCommitStoreReader(estimator gas.EvmFeeEstimator, sourceMaxGasPrice *big.Int, feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider) *IncompleteSourceCommitStoreReader {
	return &IncompleteSourceCommitStoreReader{
		estimator:          estimator,
		sourceMaxGasPrice:  sourceMaxGasPrice,
		feeEstimatorConfig: feeEstimatorConfig,
	}
}

func (i *IncompleteSourceCommitStoreReader) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ccip.CommitOnchainConfig](onchainConfig)
	if err != nil {
		return "", err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[ccip.JSONCommitOffchainConfigV1_2_0](offchainConfig)
	if err != nil {
		return "", err
	}

	i.gasPriceEstimator = prices.NewDAGasPriceEstimator(
		i.estimator,
		i.sourceMaxGasPrice,
		int64(offchainConfigParsed.ExecGasPriceDeviationPPB),
		int64(offchainConfigParsed.DAGasPriceDeviationPPB),
		i.feeEstimatorConfig,
	)
	i.offchainConfig = ccip.NewCommitOffchainConfig(
		offchainConfigParsed.ExecGasPriceDeviationPPB,
		offchainConfigParsed.GasPriceHeartBeat.Duration(),
		offchainConfigParsed.TokenPriceDeviationPPB,
		offchainConfigParsed.TokenPriceHeartBeat.Duration(),
		offchainConfigParsed.InflightCacheExpiry.Duration(),
		offchainConfigParsed.PriceReportingDisabled,
	)

	return cciptypes.Address(onchainConfigParsed.PriceRegistry.String()), nil
}

func (i *IncompleteSourceCommitStoreReader) DecodeCommitReport(ctx context.Context, report []byte) (cciptypes.CommitStoreReport, error) {
	return cciptypes.CommitStoreReport{}, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) EncodeCommitReport(ctx context.Context, report cciptypes.CommitStoreReport) ([]byte, error) {
	return []byte{}, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

// GasPriceEstimator returns an ExecGasPriceEstimator to satisfy the GasPriceEstimatorCommit interface,
// with deviationPPB values hardcoded to 0 when this implementation is first constructed.
// When ChangeConfig is called, another call to this method must be made to fetch a GasPriceEstimator with updated values
func (i *IncompleteSourceCommitStoreReader) GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorCommit, error) {
	return i.gasPriceEstimator, nil
}

func (i *IncompleteSourceCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return nil, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return nil, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	return cciptypes.CommitStoreStaticConfig{}, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return 0, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return 0, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return false, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) IsDestChainHealthy(ctx context.Context) (bool, error) {
	return false, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return false, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) OffchainConfig(ctx context.Context) (cciptypes.CommitOffchainConfig, error) {
	return i.offchainConfig, nil
}

func (i *IncompleteSourceCommitStoreReader) VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error) {
	return false, fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

func (i *IncompleteSourceCommitStoreReader) Close() error {
	return fmt.Errorf("invalid usage of IncompleteSourceCommitStoreReader")
}

// IncompleteDestCommitStoreReader is an implementation of CommitStoreReader with all valid methods except
// GasPriceEstimator, ChangeConfig, and OffchainConfig.
type IncompleteDestCommitStoreReader struct {
	cs cciptypes.CommitStoreReader
}

func NewIncompleteDestCommitStoreReader(
	ctx context.Context,
	lggr logger.Logger,
	versionFinder ccip.VersionFinder,
	address cciptypes.Address,
	ec client.Client,
	lp logpoller.LogPoller,
	feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider,
) (*IncompleteDestCommitStoreReader, error) {
	cs, err := ccip.NewCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
	if err != nil {
		return nil, err
	}

	return &IncompleteDestCommitStoreReader{
		cs: cs,
	}, nil
}

func (i *IncompleteDestCommitStoreReader) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error) {
	return "", fmt.Errorf("invalid usage of IncompleteDestCommitStoreReader")
}

func (i *IncompleteDestCommitStoreReader) DecodeCommitReport(ctx context.Context, report []byte) (cciptypes.CommitStoreReport, error) {
	return i.cs.DecodeCommitReport(ctx, report)
}

func (i *IncompleteDestCommitStoreReader) EncodeCommitReport(ctx context.Context, report cciptypes.CommitStoreReport) ([]byte, error) {
	return i.cs.EncodeCommitReport(ctx, report)
}

func (i *IncompleteDestCommitStoreReader) GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorCommit, error) {
	return nil, fmt.Errorf("invalid usage of IncompleteDestCommitStoreReader")
}

func (i *IncompleteDestCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return i.cs.GetAcceptedCommitReportsGteTimestamp(ctx, ts, confirmations)
}

func (i *IncompleteDestCommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return i.cs.GetCommitReportMatchingSeqNum(ctx, seqNum, confirmations)
}

func (i *IncompleteDestCommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	return i.cs.GetCommitStoreStaticConfig(ctx)
}

func (i *IncompleteDestCommitStoreReader) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return i.cs.GetExpectedNextSequenceNumber(ctx)
}

func (i *IncompleteDestCommitStoreReader) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return i.cs.GetLatestPriceEpochAndRound(ctx)
}

func (i *IncompleteDestCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return i.cs.IsBlessed(ctx, root)
}

func (i *IncompleteDestCommitStoreReader) IsDestChainHealthy(ctx context.Context) (bool, error) {
	return i.cs.IsDestChainHealthy(ctx)
}

func (i *IncompleteDestCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return i.cs.IsDown(ctx)
}

func (i *IncompleteDestCommitStoreReader) OffchainConfig(ctx context.Context) (cciptypes.CommitOffchainConfig, error) {
	return cciptypes.CommitOffchainConfig{}, fmt.Errorf("invalid usage of IncompleteDestCommitStoreReader")
}

func (i *IncompleteDestCommitStoreReader) VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error) {
	return i.cs.VerifyExecutionReport(ctx, report)
}

func (i *IncompleteDestCommitStoreReader) Close() error {
	return i.cs.Close()
}
