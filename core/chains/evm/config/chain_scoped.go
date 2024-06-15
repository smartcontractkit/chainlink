package config

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func NewTOMLChainScopedConfig(tomlConfig *toml.EVMConfig, lggr logger.Logger) *ChainScoped {
	return &ChainScoped{
		evmConfig: &EVMConfig{C: tomlConfig},
		lggr:      lggr}
}

// ChainScoped implements config.ChainScopedConfig with a gencfg.BasicConfig and EVMConfig.
type ChainScoped struct {
	lggr logger.Logger

	evmConfig *EVMConfig
}

func (c *ChainScoped) EVM() EVM {
	return c.evmConfig
}

func (c *ChainScoped) Nodes() toml.EVMNodes {
	return c.evmConfig.C.Nodes
}

func (c *ChainScoped) BlockEmissionIdleWarningThreshold() time.Duration {
	return c.EVM().NodeNoNewHeadsThreshold()
}

type EVMConfig struct {
	C *toml.EVMConfig
}

func (e *EVMConfig) IsEnabled() bool {
	return e.C.IsEnabled()
}

func (e *EVMConfig) TOMLString() (string, error) {
	return e.C.TOMLString()
}

func (e *EVMConfig) BalanceMonitor() BalanceMonitor {
	return &balanceMonitorConfig{c: e.C.BalanceMonitor}
}

func (e *EVMConfig) Transactions() Transactions {
	return &transactionsConfig{c: e.C.Transactions}
}

func (e *EVMConfig) HeadTracker() HeadTracker {
	return &headTrackerConfig{c: e.C.HeadTracker}
}

func (e *EVMConfig) OCR() OCR {
	return &ocrConfig{c: e.C.OCR}
}

func (e *EVMConfig) OCR2() OCR2 {
	return &ocr2Config{c: e.C.OCR2}
}

func (e *EVMConfig) Workflow() Workflow {
	return &workflowConfig{c: e.C.Workflow}
}

func (e *EVMConfig) GasEstimator() GasEstimator {
	return &gasEstimatorConfig{c: e.C.GasEstimator, blockDelay: e.C.RPCBlockQueryDelay, transactionsMaxInFlight: e.C.Transactions.MaxInFlight, k: e.C.KeySpecific}
}

func (e *EVMConfig) AutoCreateKey() bool {
	return *e.C.AutoCreateKey
}

func (e *EVMConfig) BlockBackfillDepth() uint64 {
	return uint64(*e.C.BlockBackfillDepth)
}

func (e *EVMConfig) BlockBackfillSkip() bool {
	return *e.C.BlockBackfillSkip
}

func (e *EVMConfig) LogBackfillBatchSize() uint32 {
	return *e.C.LogBackfillBatchSize
}

func (e *EVMConfig) LogPollInterval() time.Duration {
	return e.C.LogPollInterval.Duration()
}

func (e *EVMConfig) FinalityDepth() uint32 {
	return *e.C.FinalityDepth
}

func (e *EVMConfig) FinalityTagEnabled() bool {
	return *e.C.FinalityTagEnabled
}

func (e *EVMConfig) LogKeepBlocksDepth() uint32 {
	return *e.C.LogKeepBlocksDepth
}

func (e *EVMConfig) BackupLogPollerBlockDelay() uint64 {
	return *e.C.BackupLogPollerBlockDelay
}

func (e *EVMConfig) NonceAutoSync() bool {
	return *e.C.NonceAutoSync
}

func (e *EVMConfig) RPCDefaultBatchSize() uint32 {
	return *e.C.RPCDefaultBatchSize
}

func (e *EVMConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	return e.C.NoNewHeadsThreshold.Duration()
}

func (e *EVMConfig) ChainType() chaintype.ChainType {
	if e.C.ChainType == nil {
		return ""
	}
	return e.C.ChainType.ChainType()
}

func (e *EVMConfig) ChainID() *big.Int {
	return e.C.ChainID.ToInt()
}

func (e *EVMConfig) MinIncomingConfirmations() uint32 {
	return *e.C.MinIncomingConfirmations
}

func (e *EVMConfig) NodePool() NodePool {
	return &NodePoolConfig{C: e.C.NodePool}
}

func (e *EVMConfig) ClientErrors() ClientErrors {
	return &clientErrorsConfig{c: e.C.NodePool.Errors}
}

func (e *EVMConfig) NodeNoNewHeadsThreshold() time.Duration {
	return e.C.NoNewHeadsThreshold.Duration()
}

func (e *EVMConfig) MinContractPayment() *assets.Link {
	return e.C.MinContractPayment
}

func (e *EVMConfig) FlagsContractAddress() string {
	if e.C.FlagsContractAddress == nil {
		return ""
	}
	return e.C.FlagsContractAddress.String()
}

func (e *EVMConfig) LinkContractAddress() string {
	if e.C.LinkContractAddress == nil {
		return ""
	}
	return e.C.LinkContractAddress.String()
}

func (e *EVMConfig) OperatorFactoryAddress() string {
	if e.C.OperatorFactoryAddress == nil {
		return ""
	}
	return e.C.OperatorFactoryAddress.String()
}

func (e *EVMConfig) LogPrunePageSize() uint32 {
	return *e.C.LogPrunePageSize
}
