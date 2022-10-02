package v2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/multierr"

	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	gencfg "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func NewTOMLChainScopedConfig(genCfg gencfg.BasicConfig, chain *EVMConfig, lggr logger.Logger) *ChainScoped {
	return &ChainScoped{BasicConfig: genCfg, cfg: chain, lggr: lggr}
}

// ChainScoped implements config.ChainScopedConfig with a gencfg.BasicConfig and EVMConfig.
type ChainScoped struct {
	gencfg.BasicConfig
	lggr logger.Logger

	cfg *EVMConfig
}

func (c *ChainScoped) ChainID() *big.Int {
	return c.cfg.ChainID.ToInt()
}

func (c *ChainScoped) ChainType() gencfg.ChainType {
	if c.cfg.ChainType == nil {
		return ""
	}
	return gencfg.ChainType(*c.cfg.ChainType)
}

func (c *ChainScoped) Validate() (err error) {
	// Most per-chain validation is done on startup, but this combines globals as well.
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      c.OCRBlockchainTimeout(),
		ContractConfigConfirmations:            c.OCRContractConfirmations(),
		ContractConfigTrackerPollInterval:      c.OCRContractPollInterval(),
		ContractConfigTrackerSubscribeInterval: c.OCRContractSubscribeInterval(),
		ContractTransmitterTransmitTimeout:     c.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        c.OCRDatabaseTimeout(),
		DataSourceTimeout:                      c.OCRObservationTimeout(),
		DataSourceGracePeriod:                  c.OCRObservationGracePeriod(),
	}
	if ocrerr := ocr.SanityCheckLocalConfig(lc); ocrerr != nil {
		err = multierr.Append(err, ocrerr)
	}
	return
}

func (c *ChainScoped) BlockBackfillDepth() uint64 {
	return uint64(*c.cfg.BlockBackfillDepth)
}

func (c *ChainScoped) BlockBackfillSkip() bool {
	return *c.cfg.BlockBackfillSkip
}

func (c *ChainScoped) BalanceMonitorEnabled() bool {
	return *c.cfg.BalanceMonitor.Enabled
}

func (c *ChainScoped) BlockEmissionIdleWarningThreshold() time.Duration {
	return c.NodeNoNewHeadsThreshold()
}

func (c *ChainScoped) BlockHistoryEstimatorBatchSize() (size uint32) {
	return *c.cfg.GasEstimator.BlockHistory.BatchSize
}

func (c *ChainScoped) BlockHistoryEstimatorBlockDelay() uint16 {
	return *c.cfg.RPCBlockQueryDelay
}

func (c *ChainScoped) BlockHistoryEstimatorBlockHistorySize() uint16 {
	return *c.cfg.GasEstimator.BlockHistory.BlockHistorySize
}

func (c *ChainScoped) BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16 {
	if c.cfg.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks == nil {
		return uint16(c.EvmGasBumpThreshold()) + 1
	}
	return *c.cfg.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks
}

func (c *ChainScoped) BlockHistoryEstimatorTransactionPercentile() uint16 {
	return *c.cfg.GasEstimator.BlockHistory.TransactionPercentile
}

func (c *ChainScoped) EvmEIP1559DynamicFees() bool {
	return *c.cfg.GasEstimator.EIP1559DynamicFees
}

func (c *ChainScoped) EthTxReaperInterval() time.Duration {
	return c.cfg.Transactions.ReaperInterval.Duration()
}

func (c *ChainScoped) EthTxReaperThreshold() time.Duration {
	return c.cfg.Transactions.ReaperThreshold.Duration()
}

func (c *ChainScoped) EthTxResendAfterThreshold() time.Duration {
	return c.cfg.Transactions.ResendAfterThreshold.Duration()
}

func (c *ChainScoped) EvmFinalityDepth() uint32 {
	return *c.cfg.FinalityDepth
}

func (c *ChainScoped) EvmGasBumpPercent() uint16 {
	return *c.cfg.GasEstimator.BumpPercent
}

func (c *ChainScoped) EvmGasBumpThreshold() uint64 {
	return uint64(*c.cfg.GasEstimator.BumpThreshold)
}

func (c *ChainScoped) EvmGasBumpTxDepth() uint16 {
	return *c.cfg.GasEstimator.BumpTxDepth
}

func (c *ChainScoped) EvmGasBumpWei() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.BumpMin)
}

func (c *ChainScoped) EvmGasFeeCapDefault() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.FeeCapDefault)
}

func (c *ChainScoped) EvmGasLimitDefault() uint32 {
	return *c.cfg.GasEstimator.LimitDefault
}

func (c *ChainScoped) EvmGasLimitMax() uint32 {
	return *c.cfg.GasEstimator.LimitMax
}

func (c *ChainScoped) EvmGasLimitMultiplier() float32 {
	f, _ := c.cfg.GasEstimator.LimitMultiplier.BigFloat().Float32()
	return f
}

func (c *ChainScoped) EvmGasLimitTransfer() uint32 {
	return *c.cfg.GasEstimator.LimitTransfer
}

func (c *ChainScoped) EvmGasLimitOCRJobType() *uint32 {
	if t := c.cfg.GasEstimator.LimitJobType; t != nil {
		return t.OCR
	}
	return nil
}

func (c *ChainScoped) EvmGasLimitDRJobType() *uint32 {
	if t := c.cfg.GasEstimator.LimitJobType; t != nil {
		return t.DR
	}
	return nil
}

func (c *ChainScoped) EvmGasLimitVRFJobType() *uint32 {
	if t := c.cfg.GasEstimator.LimitJobType; t != nil {
		return t.VRF
	}
	return nil
}

func (c *ChainScoped) EvmGasLimitFMJobType() *uint32 {
	if t := c.cfg.GasEstimator.LimitJobType; t != nil {
		return t.FM
	}
	return nil
}

func (c *ChainScoped) EvmGasLimitKeeperJobType() *uint32 {
	if t := c.cfg.GasEstimator.LimitJobType; t != nil {
		return t.Keeper
	}
	return nil
}

func (c *ChainScoped) EvmGasPriceDefault() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.PriceDefault)
}

func (c *ChainScoped) EvmMinGasPriceWei() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.PriceMin)
}

func (c *ChainScoped) EvmMaxGasPriceWei() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.PriceMax)
}

func (c *ChainScoped) EvmGasTipCapDefault() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.TipCapDefault)
}

func (c *ChainScoped) EvmGasTipCapMinimum() *big.Int {
	return (*big.Int)(c.cfg.GasEstimator.TipCapMin)
}

func (c *ChainScoped) EvmHeadTrackerHistoryDepth() uint32 {
	return *c.cfg.HeadTracker.HistoryDepth
}

func (c *ChainScoped) EvmHeadTrackerMaxBufferSize() uint32 {
	return *c.cfg.HeadTracker.MaxBufferSize
}

func (c *ChainScoped) EvmHeadTrackerSamplingInterval() time.Duration {
	return c.cfg.HeadTracker.SamplingInterval.Duration()
}

func (c *ChainScoped) EvmLogBackfillBatchSize() uint32 {
	return *c.cfg.LogBackfillBatchSize
}

func (c *ChainScoped) EvmLogPollInterval() time.Duration {
	return c.cfg.LogPollInterval.Duration()
}

func (c *ChainScoped) EvmMaxInFlightTransactions() uint32 {
	return *c.cfg.Transactions.MaxInFlight
}

func (c *ChainScoped) EvmMaxQueuedTransactions() uint64 {
	return uint64(*c.cfg.Transactions.MaxQueued)
}

func (c *ChainScoped) EvmNonceAutoSync() bool {
	return *c.cfg.NonceAutoSync
}

func (c *ChainScoped) EvmUseForwarders() bool {
	return *c.cfg.Transactions.ForwardersEnabled
}

func (c *ChainScoped) EvmRPCDefaultBatchSize() uint32 {
	return *c.cfg.RPCDefaultBatchSize
}

func (c *ChainScoped) FlagsContractAddress() string {
	if c.cfg.FlagsContractAddress == nil {
		return ""
	}
	return c.cfg.FlagsContractAddress.String()
}

func (c *ChainScoped) GasEstimatorMode() string {
	return *c.cfg.GasEstimator.Mode
}
func (c *ChainScoped) KeySpecificMaxGasPriceWei(addr common.Address) *big.Int {
	for i := range c.cfg.KeySpecific {
		ks := c.cfg.KeySpecific[i]
		if ks.Key.Address() == addr {
			return (*big.Int)(ks.GasEstimator.PriceMax)
		}
	}
	return (*big.Int)(c.cfg.GasEstimator.PriceMax)
}

func (c *ChainScoped) LinkContractAddress() string {
	if c.cfg.LinkContractAddress == nil {
		return ""
	}
	return c.cfg.LinkContractAddress.String()
}

func (c *ChainScoped) OperatorFactoryAddress() string {
	if c.cfg.OperatorFactoryAddress == nil {
		return ""
	}
	return c.cfg.OperatorFactoryAddress.String()
}

func (c *ChainScoped) MinIncomingConfirmations() uint32 {
	return *c.cfg.MinIncomingConfirmations
}

func (c *ChainScoped) MinimumContractPayment() *assets.Link {
	return c.cfg.MinContractPayment
}

func (c *ChainScoped) NodeNoNewHeadsThreshold() time.Duration {
	return c.cfg.NoNewHeadsThreshold.Duration()
}

func (c *ChainScoped) NodePollFailureThreshold() uint32 {
	return *c.cfg.NodePool.PollFailureThreshold
}

func (c *ChainScoped) NodePollInterval() time.Duration {
	return c.cfg.NodePool.PollInterval.Duration()
}

func (c *ChainScoped) NodeSelectionMode() string {
	return *c.cfg.NodePool.SelectionMode
}

func (c *ChainScoped) OCRContractConfirmations() uint16 {
	return *c.cfg.OCR.ContractConfirmations
}

func (c *ChainScoped) OCRContractTransmitterTransmitTimeout() time.Duration {
	return c.cfg.OCR.ContractTransmitterTransmitTimeout.Duration()
}

func (c *ChainScoped) OCRObservationGracePeriod() time.Duration {
	return c.cfg.OCR.ObservationGracePeriod.Duration()
}

func (c *ChainScoped) OCRDatabaseTimeout() time.Duration {
	return c.cfg.OCR.DatabaseTimeout.Duration()
}
