package v1_0_0

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	EXEC_REPORT_ACCEPTS = "Exec report accepts"
	ReportAccepted      = "ReportAccepted"
)

var _ ccipdata.CommitStoreReader = &CommitStore{}

type CommitStore struct {
	// Static config
	commitStore               *commit_store_1_0_0.CommitStore
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
	gasPriceEstimator prices.ExecGasPriceEstimator
	offchainConfig    ccipdata.CommitOffchainConfig
}

func (c *CommitStore) GetCommitStoreStaticConfig(ctx context.Context) (ccipdata.CommitStoreStaticConfig, error) {
	legacyConfig, err := c.commitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return ccipdata.CommitStoreStaticConfig{}, errors.New("Could not get commitStore static config")
	}
	return ccipdata.CommitStoreStaticConfig{
		ChainSelector:       legacyConfig.ChainSelector,
		SourceChainSelector: legacyConfig.SourceChainSelector,
		OnRamp:              legacyConfig.OnRamp,
		ArmProxy:            legacyConfig.ArmProxy,
	}, nil
}

func (c *CommitStore) EncodeCommitReport(report ccipdata.CommitStoreReport) ([]byte, error) {
	return encodeCommitReport(c.commitReportArgs, report)
}

func encodeCommitReport(commitReportArgs abi.Arguments, report ccipdata.CommitStoreReport) ([]byte, error) {
	var tokenPriceUpdates []commit_store_1_0_0.InternalTokenPriceUpdate
	for _, tokenPriceUpdate := range report.TokenPrices {
		tokenPriceUpdates = append(tokenPriceUpdates, commit_store_1_0_0.InternalTokenPriceUpdate{
			SourceToken: tokenPriceUpdate.Token,
			UsdPerToken: tokenPriceUpdate.Value,
		})
	}
	var usdPerUnitGas = big.NewInt(0)
	var destChainSelector = uint64(0)
	if len(report.GasPrices) > 1 {
		return []byte{}, errors.Errorf("CommitStore V1_0_0 can only accept 1 gas price, received: %d", len(report.GasPrices))
	}
	if len(report.GasPrices) > 0 {
		usdPerUnitGas = report.GasPrices[0].Value
		destChainSelector = report.GasPrices[0].DestChainSelector
	}
	rep := commit_store_1_0_0.CommitStoreCommitReport{
		PriceUpdates: commit_store_1_0_0.InternalPriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			UsdPerUnitGas:     usdPerUnitGas,
			DestChainSelector: destChainSelector,
		},
		Interval:   commit_store_1_0_0.CommitStoreInterval{Min: report.Interval.Min, Max: report.Interval.Max},
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
			DestChainSelector uint64   `json:"destChainSelector"`
			UsdPerUnitGas     *big.Int `json:"usdPerUnitGas"`
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
	if commitReport.PriceUpdates.DestChainSelector != 0 {
		// No gas price update	{
		gasPrices = append(gasPrices, ccipdata.GasPrice{
			DestChainSelector: commitReport.PriceUpdates.DestChainSelector,
			Value:             commitReport.PriceUpdates.UsdPerUnitGas,
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

// CommitOffchainConfig is a legacy version of CommitOffchainConfig, used for CommitStore version 1.0.0 and 1.1.0
type CommitOffchainConfig struct {
	SourceFinalityDepth   uint32
	DestFinalityDepth     uint32
	FeeUpdateHeartBeat    config.Duration
	FeeUpdateDeviationPPB uint32
	MaxGasPrice           uint64
	InflightCacheExpiry   config.Duration
}

func (c CommitOffchainConfig) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.FeeUpdateHeartBeat.Duration() == 0 {
		return errors.New("must set FeeUpdateHeartBeat")
	}
	if c.FeeUpdateDeviationPPB == 0 {
		return errors.New("must set FeeUpdateDeviationPPB")
	}
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}

	return nil
}

func (c *CommitStore) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](onchainConfig)
	if err != nil {
		return common.Address{}, err
	}

	offchainConfigV1, err := ccipconfig.DecodeOffchainConfig[CommitOffchainConfig](offchainConfig)
	if err != nil {
		return common.Address{}, err
	}
	c.configMu.Lock()
	c.gasPriceEstimator = prices.NewExecGasPriceEstimator(
		c.estimator,
		big.NewInt(int64(offchainConfigV1.MaxGasPrice)),
		int64(offchainConfigV1.FeeUpdateDeviationPPB))
	c.offchainConfig = ccipdata.NewCommitOffchainConfig(
		offchainConfigV1.FeeUpdateDeviationPPB,
		offchainConfigV1.FeeUpdateHeartBeat.Duration(),
		offchainConfigV1.FeeUpdateDeviationPPB,
		offchainConfigV1.FeeUpdateHeartBeat.Duration(),
		offchainConfigV1.InflightCacheExpiry.Duration())
	c.configMu.Unlock()
	c.lggr.Infow("ChangeConfig",
		"offchainConfig", offchainConfigV1,
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
	return &ccipdata.CommitStoreReport{
		TokenPrices: tokenPrices,
		GasPrices:   []ccipdata.GasPrice{{DestChainSelector: repAccepted.Report.PriceUpdates.DestChainSelector, Value: repAccepted.Report.PriceUpdates.UsdPerUnitGas}},
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
	commitStore, err := commit_store_1_0_0.NewCommitStore(addr, ec)
	if err != nil {
		return nil, err
	}
	commitStoreABI := abihelpers.MustParseABI(commit_store_1_0_0.CommitStoreABI)
	eventSig := abihelpers.MustGetEventID(ReportAccepted, commitStoreABI)
	commitReportArgs := abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI)
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(EXEC_REPORT_ACCEPTS, addr.String()),
			EventSigs: []common.Hash{eventSig},
			Addresses: []common.Address{addr},
		},
	}
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
		gasPriceEstimator: prices.ExecGasPriceEstimator{},
	}, nil
}
