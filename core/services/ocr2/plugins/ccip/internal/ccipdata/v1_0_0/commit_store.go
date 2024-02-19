package v1_0_0

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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
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
	offchainConfig    cciptypes.CommitOffchainConfig
}

func (c *CommitStore) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	legacyConfig, err := c.commitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return cciptypes.CommitStoreStaticConfig{}, errors.New("Could not get commitStore static config")
	}
	return cciptypes.CommitStoreStaticConfig{
		ChainSelector:       legacyConfig.ChainSelector,
		SourceChainSelector: legacyConfig.SourceChainSelector,
		OnRamp:              ccipcalc.EvmAddrToGeneric(legacyConfig.OnRamp),
		ArmProxy:            ccipcalc.EvmAddrToGeneric(legacyConfig.ArmProxy),
	}, nil
}

func (c *CommitStore) EncodeCommitReport(report cciptypes.CommitStoreReport) ([]byte, error) {
	return encodeCommitReport(c.commitReportArgs, report)
}

func encodeCommitReport(commitReportArgs abi.Arguments, report cciptypes.CommitStoreReport) ([]byte, error) {
	var tokenPriceUpdates []commit_store_1_0_0.InternalTokenPriceUpdate
	for _, tokenPriceUpdate := range report.TokenPrices {
		sourceTokenEvmAddr, err := ccipcalc.GenericAddrToEvm(tokenPriceUpdate.Token)
		if err != nil {
			return nil, err
		}
		tokenPriceUpdates = append(tokenPriceUpdates, commit_store_1_0_0.InternalTokenPriceUpdate{
			SourceToken: sourceTokenEvmAddr,
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
	if commitReport.PriceUpdates.DestChainSelector != 0 {
		// No gas price update	{
		gasPrices = append(gasPrices, cciptypes.GasPrice{
			DestChainSelector: commitReport.PriceUpdates.DestChainSelector,
			Value:             commitReport.PriceUpdates.UsdPerUnitGas,
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

func (c *CommitStore) DecodeCommitReport(report []byte) (cciptypes.CommitStoreReport, error) {
	return DecodeCommitReport(c.commitReportArgs, report)
}

func (c *CommitStore) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return c.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, root)
}

func (c *CommitStore) OffchainConfig() cciptypes.CommitOffchainConfig {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.offchainConfig
}

func (c *CommitStore) GasPriceEstimator() cciptypes.GasPriceEstimatorCommit {
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

func (c *CommitStore) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ccipdata.CommitOnchainConfig](onchainConfig)
	if err != nil {
		return "", err
	}

	offchainConfigV1, err := ccipconfig.DecodeOffchainConfig[CommitOffchainConfig](offchainConfig)
	if err != nil {
		return "", err
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
	return cciptypes.Address(onchainConfigParsed.PriceRegistry.String()), nil
}

func (c *CommitStore) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(c.lp, c.filters, qopts...)
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
	return &cciptypes.CommitStoreReport{
		TokenPrices: tokenPrices,
		GasPrices:   []cciptypes.GasPrice{{DestChainSelector: repAccepted.Report.PriceUpdates.DestChainSelector, Value: repAccepted.Report.PriceUpdates.UsdPerUnitGas}},
		MerkleRoot:  repAccepted.Report.MerkleRoot,
		Interval:    cciptypes.CommitStoreInterval{Min: repAccepted.Report.Interval.Min, Max: repAccepted.Report.Interval.Max},
	}, nil
}

func (c *CommitStore) GetCommitReportMatchingSeqNum(ctx context.Context, seqNr uint64, confs int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
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

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.CommitStoreReport](logs, c.lggr, c.parseReport)
	if err != nil {
		return nil, fmt.Errorf("parse logs: %w", err)
	}

	parsedReports := make([]cciptypes.CommitStoreReportWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		parsedReports = append(parsedReports, cciptypes.CommitStoreReportWithTxMeta{
			TxMeta:            log.TxMeta,
			CommitStoreReport: log.Data,
		})
	}

	return parsedReports, nil
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
		offchainConfig:    cciptypes.CommitOffchainConfig{},
		gasPriceEstimator: prices.ExecGasPriceEstimator{},
	}, nil
}
