package v1_2_0

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

var _ ccipdata.CommitStoreReader = &CommitStore{}

type CommitStore struct {
	// Static config
	commitStore               *commit_store_1_2_0.CommitStore
	lggr                      logger.Logger
	lp                        logpoller.LogPoller
	address                   common.Address
	estimator                 *gas.EvmFeeEstimator
	sourceMaxGasPrice         *big.Int
	filters                   []logpoller.Filter
	reportAcceptedSig         common.Hash
	reportAcceptedMaxSeqIndex int
	commitReportArgs          abi.Arguments

	// Dynamic config
	configMu           sync.RWMutex
	gasPriceEstimator  *prices.DAGasPriceEstimator
	offchainConfig     cciptypes.CommitOffchainConfig
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader
}

func (c *CommitStore) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	staticConfig, err := c.commitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return cciptypes.CommitStoreStaticConfig{}, err
	}
	return cciptypes.CommitStoreStaticConfig{
		ChainSelector:       staticConfig.ChainSelector,
		SourceChainSelector: staticConfig.SourceChainSelector,
		OnRamp:              cciptypes.Address(staticConfig.OnRamp.String()),
		ArmProxy:            cciptypes.Address(staticConfig.ArmProxy.String()),
	}, nil
}

func (c *CommitStore) EncodeCommitReport(_ context.Context, report cciptypes.CommitStoreReport) ([]byte, error) {
	return EncodeCommitReport(c.commitReportArgs, report)
}

func EncodeCommitReport(commitReportArgs abi.Arguments, report cciptypes.CommitStoreReport) ([]byte, error) {
	var tokenPriceUpdates []commit_store_1_2_0.InternalTokenPriceUpdate
	for _, tokenPriceUpdate := range report.TokenPrices {
		tokenAddressEvm, err := ccipcalc.GenericAddrToEvm(tokenPriceUpdate.Token)
		if err != nil {
			return nil, fmt.Errorf("token price update address to evm: %w", err)
		}

		tokenPriceUpdates = append(tokenPriceUpdates, commit_store_1_2_0.InternalTokenPriceUpdate{
			SourceToken: tokenAddressEvm,
			UsdPerToken: tokenPriceUpdate.Value,
		})
	}

	var gasPriceUpdates []commit_store_1_2_0.InternalGasPriceUpdate
	for _, gasPriceUpdate := range report.GasPrices {
		gasPriceUpdates = append(gasPriceUpdates, commit_store_1_2_0.InternalGasPriceUpdate{
			DestChainSelector: gasPriceUpdate.DestChainSelector,
			UsdPerUnitGas:     gasPriceUpdate.Value,
		})
	}

	rep := commit_store_1_2_0.CommitStoreCommitReport{
		PriceUpdates: commit_store_1_2_0.InternalPriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		Interval:   commit_store_1_2_0.CommitStoreInterval{Min: report.Interval.Min, Max: report.Interval.Max},
		MerkleRoot: report.MerkleRoot,
	}
	return commitReportArgs.PackValues([]interface{}{rep})
}

func DecodeCommitReport(commitReportArgs abi.Arguments, report []byte) (cciptypes.CommitStoreReport, error) {
	unpacked, err := commitReportArgs.Unpack(report)
	if err != nil {
		return cciptypes.CommitStoreReport{}, err
	}
	if len(unpacked) != 1 {
		return cciptypes.CommitStoreReport{}, errors.New("expected single struct value")
	}

	commitReport, ok := unpacked[0].(struct {
		PriceUpdates struct {
			TokenPriceUpdates []struct {
				SourceToken common.Address `json:"sourceToken"`
				UsdPerToken *big.Int       `json:"usdPerToken"`
			} `json:"tokenPriceUpdates"`
			GasPriceUpdates []struct {
				DestChainSelector uint64   `json:"destChainSelector"`
				UsdPerUnitGas     *big.Int `json:"usdPerUnitGas"`
			} `json:"gasPriceUpdates"`
		} `json:"priceUpdates"`
		Interval struct {
			Min uint64 `json:"min"`
			Max uint64 `json:"max"`
		} `json:"interval"`
		MerkleRoot [32]byte `json:"merkleRoot"`
	})
	if !ok {
		return cciptypes.CommitStoreReport{}, errors.Errorf("invalid commit report got %T", unpacked[0])
	}

	var tokenPriceUpdates []cciptypes.TokenPrice
	for _, u := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			Token: cciptypes.Address(u.SourceToken.String()),
			Value: u.UsdPerToken,
		})
	}

	var gasPrices []cciptypes.GasPrice
	for _, u := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, cciptypes.GasPrice{
			DestChainSelector: u.DestChainSelector,
			Value:             u.UsdPerUnitGas,
		})
	}

	return cciptypes.CommitStoreReport{
		TokenPrices: tokenPriceUpdates,
		GasPrices:   gasPrices,
		Interval: cciptypes.CommitStoreInterval{
			Min: commitReport.Interval.Min,
			Max: commitReport.Interval.Max,
		},
		MerkleRoot: commitReport.MerkleRoot,
	}, nil
}

func (c *CommitStore) DecodeCommitReport(_ context.Context, report []byte) (cciptypes.CommitStoreReport, error) {
	return DecodeCommitReport(c.commitReportArgs, report)
}

func (c *CommitStore) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return c.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, root)
}

func (c *CommitStore) OffchainConfig(context.Context) (cciptypes.CommitOffchainConfig, error) {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.offchainConfig, nil
}

func (c *CommitStore) GasPriceEstimator(context.Context) (cciptypes.GasPriceEstimatorCommit, error) {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.gasPriceEstimator, nil
}

func (c *CommitStore) SetGasEstimator(ctx context.Context, gpe gas.EvmFeeEstimator) error {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	c.estimator = &gpe
	return nil
}

func (c *CommitStore) SetSourceMaxGasPrice(ctx context.Context, sourceMaxGasPrice *big.Int) error {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	c.sourceMaxGasPrice = sourceMaxGasPrice
	return nil
}

// Do not change the JSON format of this struct without consulting with the RDD people first.
type JSONCommitOffchainConfig struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	InflightCacheExpiry      config.Duration
	PriceReportingDisabled   bool
}

func (c JSONCommitOffchainConfig) Validate() error {
	if c.GasPriceHeartBeat.Duration() == 0 {
		return errors.New("must set GasPriceHeartBeat")
	}
	if c.ExecGasPriceDeviationPPB == 0 {
		return errors.New("must set ExecGasPriceDeviationPPB")
	}
	if c.TokenPriceHeartBeat.Duration() == 0 {
		return errors.New("must set TokenPriceHeartBeat")
	}
	if c.TokenPriceDeviationPPB == 0 {
		return errors.New("must set TokenPriceDeviationPPB")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	// DAGasPriceDeviationPPB is not validated because it can be 0 on non-rollups

	return nil
}

func (c *CommitStore) ChangeConfig(_ context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](onchainConfig)
	if err != nil {
		return "", err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[JSONCommitOffchainConfig](offchainConfig)
	if err != nil {
		return "", err
	}
	c.configMu.Lock()
	defer c.configMu.Unlock()

	if c.estimator == nil {
		return "", fmt.Errorf("this CommitStore estimator is nil. SetGasEstimator should be called before ChangeConfig")
	}

	if c.sourceMaxGasPrice == nil {
		return "", fmt.Errorf("this CommitStore sourceMaxGasPrice is nil. SetSourceMaxGasPrice should be called before ChangeConfig")
	}

	c.gasPriceEstimator = prices.NewDAGasPriceEstimator(
		*c.estimator,
		c.sourceMaxGasPrice,
		int64(offchainConfigParsed.ExecGasPriceDeviationPPB),
		int64(offchainConfigParsed.DAGasPriceDeviationPPB),
		c.feeEstimatorConfig,
	)
	c.offchainConfig = ccipdata.NewCommitOffchainConfig(
		offchainConfigParsed.ExecGasPriceDeviationPPB,
		offchainConfigParsed.GasPriceHeartBeat.Duration(),
		offchainConfigParsed.TokenPriceDeviationPPB,
		offchainConfigParsed.TokenPriceHeartBeat.Duration(),
		offchainConfigParsed.InflightCacheExpiry.Duration(),
		offchainConfigParsed.PriceReportingDisabled,
	)

	c.lggr.Infow("ChangeConfig",
		"offchainConfig", offchainConfigParsed,
		"onchainConfig", onchainConfigParsed,
	)
	return cciptypes.Address(onchainConfigParsed.PriceRegistry.String()), nil
}

func (c *CommitStore) Close() error {
	return logpollerutil.UnregisterLpFilters(c.lp, c.filters)
}

func (c *CommitStore) parseReport(log types.Log) (*cciptypes.CommitStoreReport, error) {
	repAccepted, err := c.commitStore.ParseReportAccepted(log)
	if err != nil {
		return nil, err
	}
	// Translate to common struct.
	var tokenPrices []cciptypes.TokenPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.TokenPriceUpdates {
		tokenPrices = append(tokenPrices, cciptypes.TokenPrice{
			Token: cciptypes.Address(tpu.SourceToken.String()),
			Value: tpu.UsdPerToken,
		})
	}
	var gasPrices []cciptypes.GasPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, cciptypes.GasPrice{
			DestChainSelector: tpu.DestChainSelector,
			Value:             tpu.UsdPerUnitGas,
		})
	}

	return &cciptypes.CommitStoreReport{
		TokenPrices: tokenPrices,
		GasPrices:   gasPrices,
		MerkleRoot:  repAccepted.Report.MerkleRoot,
		Interval:    cciptypes.CommitStoreInterval{Min: repAccepted.Report.Interval.Min, Max: repAccepted.Report.Interval.Max},
	}, nil
}

func (c *CommitStore) GetCommitReportMatchingSeqNum(ctx context.Context, seqNr uint64, confs int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	logs, err := c.lp.LogsDataWordBetween(
		ctx,
		c.reportAcceptedSig,
		c.address,
		c.reportAcceptedMaxSeqIndex-1,
		c.reportAcceptedMaxSeqIndex,
		logpoller.EvmWord(seqNr),
		evmtypes.Confirmations(confs),
	)
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.CommitStoreReport](
		logs,
		c.lggr,
		c.parseReport,
	)
	if err != nil {
		return nil, err
	}

	res := make([]cciptypes.CommitStoreReportWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.CommitStoreReportWithTxMeta{
			TxMeta:            log.TxMeta,
			CommitStoreReport: log.Data,
		})
	}

	if len(res) > 1 {
		c.lggr.Errorw("More than one report found for seqNr", "seqNr", seqNr, "commitReports", parsedLogs)
		return res[:1], nil
	}
	return res, nil
}

func (c *CommitStore) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	latestBlock, err := c.lp.LatestBlock(ctx)
	if err != nil {
		return nil, err
	}

	reportsQuery, err := query.Where(
		c.address.String(),
		logpoller.NewAddressFilter(c.address),
		logpoller.NewEventSigFilter(c.reportAcceptedSig),
		query.Timestamp(uint64(ts.Unix()), primitives.Gte),
		logpoller.NewConfirmationsFilter(evmtypes.Confirmations(confs)),
	)
	if err != nil {
		return nil, err
	}

	logs, err := c.lp.FilteredLogs(
		ctx,
		reportsQuery.Expressions,
		query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc)),
		"GetAcceptedCommitReportsGteTimestamp",
	)
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.CommitStoreReport](logs, c.lggr, c.parseReport)
	if err != nil {
		return nil, fmt.Errorf("parse logs: %w", err)
	}

	res := make([]cciptypes.CommitStoreReportWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.CommitStoreReportWithTxMeta{
			TxMeta:            log.TxMeta.WithFinalityStatus(uint64(latestBlock.FinalizedBlockNumber)),
			CommitStoreReport: log.Data,
		})
	}
	return res, nil
}

func (c *CommitStore) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return c.commitStore.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: ctx})
}

func (c *CommitStore) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return c.commitStore.GetLatestPriceEpochAndRound(&bind.CallOpts{Context: ctx})
}

func (c *CommitStore) IsDestChainHealthy(context.Context) (bool, error) {
	if err := c.lp.Healthy(); err != nil {
		return false, nil
	}
	return true, nil
}

func (c *CommitStore) IsDown(ctx context.Context) (bool, error) {
	unPausedAndHealthy, err := c.commitStore.IsUnpausedAndARMHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		return true, err
	}
	return !unPausedAndHealthy, nil
}

func (c *CommitStore) VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error) {
	var hashes [][32]byte
	for _, msg := range report.Messages {
		hashes = append(hashes, msg.Hash)
	}
	res, err := c.commitStore.Verify(&bind.CallOpts{Context: ctx}, hashes, report.Proofs, report.ProofFlagBits)
	if err != nil {
		c.lggr.Errorw("Unable to call verify", "messages", report.Messages, "err", err)
		return false, nil
	}
	// No timestamp, means failed to verify root.
	if res.Cmp(big.NewInt(0)) == 0 {
		c.lggr.Errorw("Root does not verify", "messages", report.Messages)
		return false, nil
	}
	return true, nil
}

func (c *CommitStore) RegisterFilters() error {
	return logpollerutil.RegisterLpFilters(c.lp, c.filters)
}

func NewCommitStore(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (*CommitStore, error) {
	commitStore, err := commit_store_1_2_0.NewCommitStore(addr, ec)
	if err != nil {
		return nil, err
	}
	commitStoreABI := abihelpers.MustParseABI(commit_store_1_2_0.CommitStoreABI)
	eventSig := abihelpers.MustGetEventID(v1_0_0.ReportAccepted, commitStoreABI)
	commitReportArgs := abihelpers.MustGetEventInputs(v1_0_0.ReportAccepted, commitStoreABI)
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(v1_0_0.EXEC_REPORT_ACCEPTS, addr.String()),
			EventSigs: []common.Hash{eventSig},
			Addresses: []common.Address{addr},
			Retention: ccipdata.CommitExecLogsRetention,
		},
	}

	return &CommitStore{
		commitStore: commitStore,
		address:     addr,
		lggr:        lggr,
		lp:          lp,

		// Note that sourceMaxGasPrice and estimator now have explicit setters (CCIP-2493)

		filters:           filters,
		commitReportArgs:  commitReportArgs,
		reportAcceptedSig: eventSig,
		// offset || priceUpdatesOffset || minSeqNum || maxSeqNum || merkleRoot
		reportAcceptedMaxSeqIndex: 3,
		configMu:                  sync.RWMutex{},

		// The fields below are initially empty and set on ChangeConfig method
		offchainConfig:     cciptypes.CommitOffchainConfig{},
		gasPriceEstimator:  nil,
		feeEstimatorConfig: feeEstimatorConfig,
	}, nil
}
