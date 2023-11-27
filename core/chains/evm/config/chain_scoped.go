package config

import (
	"math/big"
	"time"

	"go.uber.org/multierr"

	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

func NewTOMLChainScopedConfig(appCfg config.AppConfig, tomlConfig *toml.EVMConfig, lggr logger.Logger) *ChainScoped {
	return &ChainScoped{
		AppConfig: appCfg,
		evmConfig: &evmConfig{c: tomlConfig},
		lggr:      lggr}
}

// ChainScoped implements config.ChainScopedConfig with a gencfg.BasicConfig and EVMConfig.
type ChainScoped struct {
	config.AppConfig
	lggr logger.Logger

	evmConfig *evmConfig
}

func (c *ChainScoped) EVM() EVM {
	return c.evmConfig
}

func (c *ChainScoped) Nodes() toml.EVMNodes {
	return c.evmConfig.c.Nodes
}

func (c *ChainScoped) BlockEmissionIdleWarningThreshold() time.Duration {
	return c.EVM().NodeNoNewHeadsThreshold()
}

func (c *ChainScoped) Validate() (err error) {
	// Most per-chain validation is done on startup, but this combines globals as well.
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      c.OCR().BlockchainTimeout(),
		ContractConfigConfirmations:            c.EVM().OCR().ContractConfirmations(),
		ContractConfigTrackerPollInterval:      c.OCR().ContractPollInterval(),
		ContractConfigTrackerSubscribeInterval: c.OCR().ContractSubscribeInterval(),
		ContractTransmitterTransmitTimeout:     c.EVM().OCR().ContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        c.EVM().OCR().DatabaseTimeout(),
		DataSourceTimeout:                      c.OCR().ObservationTimeout(),
		DataSourceGracePeriod:                  c.EVM().OCR().ObservationGracePeriod(),
	}
	if ocrerr := ocr.SanityCheckLocalConfig(lc); ocrerr != nil {
		err = multierr.Append(err, ocrerr)
	}
	return
}

type evmConfig struct {
	c *toml.EVMConfig
}

func (e *evmConfig) IsEnabled() bool {
	return e.c.IsEnabled()
}

func (e *evmConfig) TOMLString() (string, error) {
	return e.c.TOMLString()
}

func (e *evmConfig) BalanceMonitor() BalanceMonitor {
	return &balanceMonitorConfig{c: e.c.BalanceMonitor}
}

func (e *evmConfig) Transactions() Transactions {
	return &transactionsConfig{c: e.c.Transactions}
}

func (e *evmConfig) HeadTracker() HeadTracker {
	return &headTrackerConfig{c: e.c.HeadTracker}
}

func (e *evmConfig) OCR() OCR {
	return &ocrConfig{c: e.c.OCR}
}

func (e *evmConfig) OCR2() OCR2 {
	return &ocr2Config{c: e.c.OCR2}
}

func (e *evmConfig) GasEstimator() GasEstimator {
	return &gasEstimatorConfig{c: e.c.GasEstimator, blockDelay: e.c.RPCBlockQueryDelay, transactionsMaxInFlight: e.c.Transactions.MaxInFlight, k: e.c.KeySpecific}
}

func (e *evmConfig) AutoCreateKey() bool {
	return *e.c.AutoCreateKey
}

func (e *evmConfig) BlockBackfillDepth() uint64 {
	return uint64(*e.c.BlockBackfillDepth)
}

func (e *evmConfig) BlockBackfillSkip() bool {
	return *e.c.BlockBackfillSkip
}

func (e *evmConfig) LogBackfillBatchSize() uint32 {
	return *e.c.LogBackfillBatchSize
}

func (e *evmConfig) LogPollInterval() time.Duration {
	return e.c.LogPollInterval.Duration()
}

func (e *evmConfig) FinalityDepth() uint32 {
	return *e.c.FinalityDepth
}

func (e *evmConfig) FinalityTagEnabled() bool {
	return *e.c.FinalityTagEnabled
}

func (e *evmConfig) LogKeepBlocksDepth() uint32 {
	return *e.c.LogKeepBlocksDepth
}

func (e *evmConfig) NonceAutoSync() bool {
	return *e.c.NonceAutoSync
}

func (e *evmConfig) RPCDefaultBatchSize() uint32 {
	return *e.c.RPCDefaultBatchSize
}

func (e *evmConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	return e.c.NoNewHeadsThreshold.Duration()
}

func (e *evmConfig) ChainType() commonconfig.ChainType {
	if e.c.ChainType == nil {
		return ""
	}
	return commonconfig.ChainType(*e.c.ChainType)
}

func (e *evmConfig) ChainID() *big.Int {
	return e.c.ChainID.ToInt()
}

func (e *evmConfig) MinIncomingConfirmations() uint32 {
	return *e.c.MinIncomingConfirmations
}

func (e *evmConfig) NodePool() NodePool {
	return &nodePoolConfig{c: e.c.NodePool}
}

func (e *evmConfig) NodeNoNewHeadsThreshold() time.Duration {
	return e.c.NoNewHeadsThreshold.Duration()
}

func (e *evmConfig) MinContractPayment() *assets.Link {
	return e.c.MinContractPayment
}

func (e *evmConfig) FlagsContractAddress() string {
	if e.c.FlagsContractAddress == nil {
		return ""
	}
	return e.c.FlagsContractAddress.String()
}

func (e *evmConfig) LinkContractAddress() string {
	if e.c.LinkContractAddress == nil {
		return ""
	}
	return e.c.LinkContractAddress.String()
}

func (e *evmConfig) OperatorFactoryAddress() string {
	if e.c.OperatorFactoryAddress == nil {
		return ""
	}
	return e.c.OperatorFactoryAddress.String()
}
