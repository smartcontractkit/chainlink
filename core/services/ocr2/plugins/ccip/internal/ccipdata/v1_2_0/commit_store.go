package v1_2_0

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ ccipdata.CommitStoreReader = &CommitStore{}

type CommitStore struct {
	// Static config
	commitStore               *commit_store.CommitStore
	lggr                      logger.Logger
	lp                        logpoller.LogPoller
	address                   common.Address
	estimator                 gas.EvmFeeEstimator
	filters                   []logpoller.Filter
	reportAcceptedSig         common.Hash
	reportAcceptedMaxSeqIndex int
	commitReportArgs          abi.Arguments

	// Dynamic config
	configMu          sync.RWMutex
	gasPriceEstimator prices.DAGasPriceEstimator
	offchainConfig    ccipdata.CommitOffchainConfig
}

func (c *CommitStore) GetCommitStoreStaticConfig(ctx context.Context) (ccipdata.CommitStoreStaticConfig, error) {
	config, err := c.commitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return ccipdata.CommitStoreStaticConfig{}, err
	}
	return ccipdata.CommitStoreStaticConfig{
		ChainSelector:       config.ChainSelector,
		SourceChainSelector: config.SourceChainSelector,
		OnRamp:              config.OnRamp,
		ArmProxy:            config.ArmProxy,
	}, nil
}

func (c *CommitStore) EncodeCommitReport(report ccipdata.CommitStoreReport) ([]byte, error) {
	return EncodeCommitReport(c.commitReportArgs, report)
}

func EncodeCommitReport(commitReportArgs abi.Arguments, report ccipdata.CommitStoreReport) ([]byte, error) {
	var tokenPriceUpdates []commit_store.InternalTokenPriceUpdate
	for _, tokenPriceUpdate := range report.TokenPrices {
		tokenPriceUpdates = append(tokenPriceUpdates, commit_store.InternalTokenPriceUpdate{
			SourceToken: tokenPriceUpdate.Token,
			UsdPerToken: tokenPriceUpdate.Value,
		})
	}

	var gasPriceUpdates []commit_store.InternalGasPriceUpdate
	for _, gasPriceUpdate := range report.GasPrices {
		gasPriceUpdates = append(gasPriceUpdates, commit_store.InternalGasPriceUpdate{
			DestChainSelector: gasPriceUpdate.DestChainSelector,
			UsdPerUnitGas:     gasPriceUpdate.Value,
		})
	}

	rep := commit_store.CommitStoreCommitReport{
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		Interval:   commit_store.CommitStoreInterval{Min: report.Interval.Min, Max: report.Interval.Max},
		MerkleRoot: report.MerkleRoot,
	}
	return commitReportArgs.PackValues([]interface{}{rep})
}

func DecodeCommitReport(commitReportArgs abi.Arguments, report []byte) (ccipdata.CommitStoreReport, error) {
	unpacked, err := commitReportArgs.Unpack(report)
	if err != nil {
		return ccipdata.CommitStoreReport{}, err
	}
	if len(unpacked) != 1 {
		return ccipdata.CommitStoreReport{}, errors.New("expected single struct value")
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
		return ccipdata.CommitStoreReport{}, errors.Errorf("invalid commit report got %T", unpacked[0])
	}

	var tokenPriceUpdates []ccipdata.TokenPrice
	for _, u := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, ccipdata.TokenPrice{
			Token: u.SourceToken,
			Value: u.UsdPerToken,
		})
	}

	var gasPrices []ccipdata.GasPrice
	for _, u := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, ccipdata.GasPrice{
			DestChainSelector: u.DestChainSelector,
			Value:             u.UsdPerUnitGas,
		})
	}

	return ccipdata.CommitStoreReport{
		TokenPrices: tokenPriceUpdates,
		GasPrices:   gasPrices,
		Interval: ccipdata.CommitStoreInterval{
			Min: commitReport.Interval.Min,
			Max: commitReport.Interval.Max,
		},
		MerkleRoot: commitReport.MerkleRoot,
	}, nil
}

func (c *CommitStore) DecodeCommitReport(report []byte) (ccipdata.CommitStoreReport, error) {
	return DecodeCommitReport(c.commitReportArgs, report)
}

func (c *CommitStore) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return c.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, root)
}

func (c *CommitStore) OffchainConfig() ccipdata.CommitOffchainConfig {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.offchainConfig
}

func (c *CommitStore) GasPriceEstimator() prices.GasPriceEstimatorCommit {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.gasPriceEstimator
}

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type JSONCommitOffchainConfig struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	MaxGasPrice              uint64
	SourceMaxGasPrice        uint64
	InflightCacheExpiry      config.Duration
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
	if c.SourceMaxGasPrice == 0 && c.MaxGasPrice == 0 {
		return errors.New("must set SourceMaxGasPrice")
	}
	if c.SourceMaxGasPrice != 0 && c.MaxGasPrice != 0 {
		return errors.New("cannot set both MaxGasPrice and SourceMaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	// DAGasPriceDeviationPPB is not validated because it can be 0 on non-rollups

	return nil
}

func (c *JSONCommitOffchainConfig) ComputeSourceMaxGasPrice() uint64 {
	if c.SourceMaxGasPrice != 0 {
		return c.SourceMaxGasPrice
	}
	return c.MaxGasPrice
}

func (c *CommitStore) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](onchainConfig)
	if err != nil {
		return common.Address{}, err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[JSONCommitOffchainConfig](offchainConfig)
	if err != nil {
		return common.Address{}, err
	}
	c.configMu.Lock()

	c.lggr.Infow("Initializing NewDAGasPriceEstimator", "estimator", c.estimator, "l1Oracle", c.estimator.L1Oracle())
	c.gasPriceEstimator = prices.NewDAGasPriceEstimator(
		c.estimator,
		big.NewInt(int64(offchainConfigParsed.ComputeSourceMaxGasPrice())),
		int64(offchainConfigParsed.ExecGasPriceDeviationPPB),
		int64(offchainConfigParsed.DAGasPriceDeviationPPB),
	)
	c.offchainConfig = ccipdata.NewCommitOffchainConfig(
		offchainConfigParsed.ExecGasPriceDeviationPPB,
		offchainConfigParsed.GasPriceHeartBeat.Duration(),
		offchainConfigParsed.TokenPriceDeviationPPB,
		offchainConfigParsed.TokenPriceHeartBeat.Duration(),
		offchainConfigParsed.InflightCacheExpiry.Duration(),
	)
	c.configMu.Unlock()

	c.lggr.Infow("ChangeConfig",
		"offchainConfig", offchainConfigParsed,
		"onchainConfig", onchainConfigParsed,
	)
	return onchainConfigParsed.PriceRegistry, nil
}

func (c *CommitStore) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(c.lp, c.filters, qopts...)
}

func (c *CommitStore) parseReport(log types.Log) (*ccipdata.CommitStoreReport, error) {
	repAccepted, err := c.commitStore.ParseReportAccepted(log)
	if err != nil {
		return nil, err
	}
	// Translate to common struct.
	var tokenPrices []ccipdata.TokenPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.TokenPriceUpdates {
		tokenPrices = append(tokenPrices, ccipdata.TokenPrice{
			Token: tpu.SourceToken,
			Value: tpu.UsdPerToken,
		})
	}
	var gasPrices []ccipdata.GasPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, ccipdata.GasPrice{
			DestChainSelector: tpu.DestChainSelector,
			Value:             tpu.UsdPerUnitGas,
		})
	}

	return &ccipdata.CommitStoreReport{
		TokenPrices: tokenPrices,
		GasPrices:   gasPrices,
		MerkleRoot:  repAccepted.Report.MerkleRoot,
		Interval:    ccipdata.CommitStoreInterval{Min: repAccepted.Report.Interval.Min, Max: repAccepted.Report.Interval.Max},
	}, nil
}

func (c *CommitStore) GetCommitReportMatchingSeqNum(ctx context.Context, seqNr uint64, confs int) ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
	logs, err := c.lp.LogsDataWordBetween(
		c.reportAcceptedSig,
		c.address,
		c.reportAcceptedMaxSeqIndex-1,
		c.reportAcceptedMaxSeqIndex,
		logpoller.EvmWord(seqNr),
		logpoller.Confirmations(confs),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[ccipdata.CommitStoreReport](
		logs,
		c.lggr,
		c.parseReport,
	)
	if err != nil {
		return nil, err
	}

	if len(parsedLogs) > 1 {
		c.lggr.Errorw("More than one report found for seqNr", "seqNr", seqNr, "commitReports", parsedLogs)
		return parsedLogs[:1], nil
	}
	return parsedLogs, nil
}

func (c *CommitStore) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
	logs, err := c.lp.LogsCreatedAfter(
		c.reportAcceptedSig,
		c.address,
		ts,
		logpoller.Confirmations(confs),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return ccipdata.ParseLogs[ccipdata.CommitStoreReport](
		logs,
		c.lggr,
		c.parseReport,
	)
}

func (c *CommitStore) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return c.commitStore.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: ctx})
}

func (c *CommitStore) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return c.commitStore.GetLatestPriceEpochAndRound(&bind.CallOpts{Context: ctx})
}

func (c *CommitStore) IsDown(ctx context.Context) (bool, error) {
	unPausedAndHealthy, err := c.commitStore.IsUnpausedAndARMHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		// If we cannot read the state, assume the worst
		c.lggr.Errorw("Unable to read CommitStore IsUnpausedAndARMHealthy", "err", err)
		return true, nil
	}
	return !unPausedAndHealthy, nil
}

func (c *CommitStore) VerifyExecutionReport(ctx context.Context, report ccipdata.ExecReport) (bool, error) {
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

func (c *CommitStore) RegisterFilters(qopts ...pg.QOpt) error {
	return logpollerutil.RegisterLpFilters(c.lp, c.filters, qopts...)
}

func NewCommitStore(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (*CommitStore, error) {
	commitStore, err := commit_store.NewCommitStore(addr, ec)
	if err != nil {
		return nil, err
	}
	commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
	eventSig := abihelpers.MustGetEventID(v1_0_0.ReportAccepted, commitStoreABI)
	commitReportArgs := abihelpers.MustGetEventInputs(v1_0_0.ReportAccepted, commitStoreABI)
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(v1_0_0.EXEC_REPORT_ACCEPTS, addr.String()),
			EventSigs: []common.Hash{eventSig},
			Addresses: []common.Address{addr},
		},
	}
	lggr.Infow("Initializing CommitStore with estimator", "estimator", estimator)

	return &CommitStore{
		commitStore:       commitStore,
		address:           addr,
		lggr:              lggr,
		lp:                lp,
		estimator:         estimator,
		filters:           filters,
		commitReportArgs:  commitReportArgs,
		reportAcceptedSig: eventSig,
		// offset || priceUpdatesOffset || minSeqNum || maxSeqNum || merkleRoot
		reportAcceptedMaxSeqIndex: 3,
		configMu:                  sync.RWMutex{},

		// The fields below are initially empty and set on ChangeConfig method
		offchainConfig:    ccipdata.CommitOffchainConfig{},
		gasPriceEstimator: prices.DAGasPriceEstimator{},
	}, nil
}
