package ccipdata

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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ CommitStoreReader = &CommitStoreV1_2_0{}

type CommitStoreV1_2_0 struct {
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
	offchainConfig    CommitOffchainConfig
}

func (c *CommitStoreV1_2_0) EncodeCommitReport(report CommitStoreReport) ([]byte, error) {
	return encodeCommitReportV1_2_0(c.commitReportArgs, report)
}

func encodeCommitReportV1_2_0(commitReportArgs abi.Arguments, report CommitStoreReport) ([]byte, error) {
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

func decodeCommitReportV1_2_0(commitReportArgs abi.Arguments, report []byte) (CommitStoreReport, error) {
	unpacked, err := commitReportArgs.Unpack(report)
	if err != nil {
		return CommitStoreReport{}, err
	}
	if len(unpacked) != 1 {
		return CommitStoreReport{}, errors.New("expected single struct value")
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
		return CommitStoreReport{}, errors.Errorf("invalid commit report got %T", unpacked[0])
	}

	var tokenPriceUpdates []TokenPrice
	for _, u := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, TokenPrice{
			Token: u.SourceToken,
			Value: u.UsdPerToken,
		})
	}

	var gasPrices []GasPrice
	for _, u := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, GasPrice{
			DestChainSelector: u.DestChainSelector,
			Value:             u.UsdPerUnitGas,
		})
	}

	return CommitStoreReport{
		TokenPrices: tokenPriceUpdates,
		GasPrices:   gasPrices,
		Interval: CommitStoreInterval{
			Min: commitReport.Interval.Min,
			Max: commitReport.Interval.Max,
		},
		MerkleRoot: commitReport.MerkleRoot,
	}, nil
}

func (c *CommitStoreV1_2_0) DecodeCommitReport(report []byte) (CommitStoreReport, error) {
	return decodeCommitReportV1_2_0(c.commitReportArgs, report)
}

func (c *CommitStoreV1_2_0) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return c.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, root)
}

func (c *CommitStoreV1_2_0) OffchainConfig() CommitOffchainConfig {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.offchainConfig
}

func (c *CommitStoreV1_2_0) GasPriceEstimator() prices.GasPriceEstimatorCommit {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.gasPriceEstimator
}

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type CommitOffchainConfigV1_2_0 struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        models.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      models.Duration
	TokenPriceDeviationPPB   uint32
	MaxGasPrice              uint64
	InflightCacheExpiry      models.Duration
}

func (c CommitOffchainConfigV1_2_0) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.GasPriceHeartBeat.Duration() == 0 {
		return errors.New("must set GasPriceHeartBeat")
	}
	if c.DAGasPriceDeviationPPB == 0 {
		return errors.New("must set DAGasPriceDeviationPPB")
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
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}

	return nil
}

func (c *CommitStoreV1_2_0) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[CommitOnchainConfig](onchainConfig)
	if err != nil {
		return common.Address{}, err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[CommitOffchainConfigV1_2_0](offchainConfig)
	if err != nil {
		return common.Address{}, err
	}
	c.configMu.Lock()

	c.lggr.Infow("Initializing NewDAGasPriceEstimator", "estimator", c.estimator, "l1Oracle", c.estimator.L1Oracle())
	c.gasPriceEstimator = prices.NewDAGasPriceEstimator(
		c.estimator,
		big.NewInt(int64(offchainConfigParsed.MaxGasPrice)),
		int64(offchainConfigParsed.ExecGasPriceDeviationPPB),
		int64(offchainConfigParsed.DAGasPriceDeviationPPB),
	)
	c.offchainConfig = NewCommitOffchainConfig(
		offchainConfigParsed.SourceFinalityDepth,
		offchainConfigParsed.ExecGasPriceDeviationPPB,
		offchainConfigParsed.GasPriceHeartBeat.Duration(),
		offchainConfigParsed.TokenPriceDeviationPPB,
		offchainConfigParsed.TokenPriceHeartBeat.Duration(),
		offchainConfigParsed.InflightCacheExpiry.Duration(),
		offchainConfigParsed.DestFinalityDepth,
	)
	c.configMu.Unlock()

	c.lggr.Infow("ChangeConfig",
		"offchainConfig", offchainConfigParsed,
		"onchainConfig", onchainConfigParsed,
	)
	return onchainConfigParsed.PriceRegistry, nil
}

func (c *CommitStoreV1_2_0) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(c.lp, c.filters, qopts...)
}

func (c *CommitStoreV1_2_0) parseReport(log types.Log) (*CommitStoreReport, error) {
	repAccepted, err := c.commitStore.ParseReportAccepted(log)
	if err != nil {
		return nil, err
	}
	// Translate to common struct.
	var tokenPrices []TokenPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.TokenPriceUpdates {
		tokenPrices = append(tokenPrices, TokenPrice{
			Token: tpu.SourceToken,
			Value: tpu.UsdPerToken,
		})
	}
	var gasPrices []GasPrice
	for _, tpu := range repAccepted.Report.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, GasPrice{
			DestChainSelector: tpu.DestChainSelector,
			Value:             tpu.UsdPerUnitGas,
		})
	}

	return &CommitStoreReport{
		TokenPrices: tokenPrices,
		GasPrices:   gasPrices,
		MerkleRoot:  repAccepted.Report.MerkleRoot,
		Interval:    CommitStoreInterval{Min: repAccepted.Report.Interval.Min, Max: repAccepted.Report.Interval.Max},
	}, nil
}

func (c *CommitStoreV1_2_0) GetAcceptedCommitReportsGteSeqNum(ctx context.Context, seqNum uint64, confs int) ([]Event[CommitStoreReport], error) {
	logs, err := c.lp.LogsDataWordGreaterThan(
		c.reportAcceptedSig,
		c.address,
		c.reportAcceptedMaxSeqIndex,
		logpoller.EvmWord(seqNum),
		confs,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return parseLogs[CommitStoreReport](
		logs,
		c.lggr,
		c.parseReport,
	)
}

func (c *CommitStoreV1_2_0) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]Event[CommitStoreReport], error) {
	logs, err := c.lp.LogsCreatedAfter(
		c.reportAcceptedSig,
		c.address,
		ts,
		confs,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return parseLogs[CommitStoreReport](
		logs,
		c.lggr,
		c.parseReport,
	)
}

func (c *CommitStoreV1_2_0) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return c.commitStore.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: ctx})
}

func (c *CommitStoreV1_2_0) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return c.commitStore.GetLatestPriceEpochAndRound(&bind.CallOpts{Context: ctx})
}

func (c *CommitStoreV1_2_0) IsDown(ctx context.Context) (bool, error) {
	unPausedAndHealthy, err := c.commitStore.IsUnpausedAndARMHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		// If we cannot read the state, assume the worst
		c.lggr.Errorw("Unable to read CommitStore IsUnpausedAndARMHealthy", "err", err)
		return true, nil
	}
	return !unPausedAndHealthy, nil
}

func (c *CommitStoreV1_2_0) VerifyExecutionReport(ctx context.Context, report ExecReport) (bool, error) {
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

func NewCommitStoreV1_2_0(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (*CommitStoreV1_2_0, error) {
	commitStore, err := commit_store.NewCommitStore(addr, ec)
	if err != nil {
		return nil, err
	}
	commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
	eventSig := abihelpers.MustGetEventID(ReportAccepted, commitStoreABI)
	commitReportArgs := abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI)
	var filters = []logpoller.Filter{
		{
			Name:      logpoller.FilterName(EXEC_REPORT_ACCEPTS, addr.String()),
			EventSigs: []common.Hash{eventSig},
			Addresses: []common.Address{addr},
		},
	}
	if err := logpollerutil.RegisterLpFilters(lp, filters); err != nil {
		return nil, err
	}

	lggr.Infow("Initializing CommitStoreV1_2_0 with estimator", "estimator", estimator)

	return &CommitStoreV1_2_0{
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
		offchainConfig:    CommitOffchainConfig{},
		gasPriceEstimator: prices.DAGasPriceEstimator{},
	}, nil
}
