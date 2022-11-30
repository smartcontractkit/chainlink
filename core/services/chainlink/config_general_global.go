package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

func (g *generalConfig) GlobalBalanceMonitorEnabled() (bool, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorCheckInclusionBlocks() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorCheckInclusionPercentile() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalChainType() (string, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmEIP1559DynamicFees() (bool, bool)      { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmFinalityDepth() (uint32, bool)         { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasBumpPercent() (uint16, bool)        { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasBumpThreshold() (uint64, bool)      { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasBumpTxDepth() (uint16, bool)        { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasBumpWei() (*assets.Wei, bool)       { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasFeeCapDefault() (*assets.Wei, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitDefault() (uint32, bool)       { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitMax() (uint32, bool)           { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitMultiplier() (float32, bool)   { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitTransfer() (uint32, bool)      { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasPriceDefault() (*assets.Wei, bool)  { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasTipCapDefault() (*assets.Wei, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasTipCapMinimum() (*assets.Wei, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmLogPollInterval() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmLogKeepBlocksDepth() (uint32, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmMaxGasPriceWei() (*assets.Wei, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmMinGasPriceWei() (*assets.Wei, bool)   { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmNonceAutoSync() (bool, bool)           { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmUseForwarders() (bool, bool)           { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool)   { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalFlagsContractAddress() (string, bool)     { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalGasEstimatorMode() (string, bool)         { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalLinkContractAddress() (string, bool)      { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalOperatorFactoryAddress() (string, bool)   { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalMinIncomingConfirmations() (uint32, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalNodeNoNewHeadsThreshold() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalNodePollFailureThreshold() (uint32, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalNodePollInterval() (time.Duration, bool)  { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalNodeSelectionMode() (string, bool)        { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalNodeSyncThreshold() (uint32, bool)        { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalOCRContractConfirmations() (uint16, bool) { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalOCRDatabaseTimeout() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}
func (g *generalConfig) GlobalOCRObservationGracePeriod() (time.Duration, bool) {
	panic(v2.ErrUnsupported)
}

func (g *generalConfig) GlobalOCR2AutomationGasLimit() (uint32, bool)   { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitOCRJobType() (uint32, bool)    { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitDRJobType() (uint32, bool)     { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitVRFJobType() (uint32, bool)    { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitFMJobType() (uint32, bool)     { panic(v2.ErrUnsupported) }
func (g *generalConfig) GlobalEvmGasLimitKeeperJobType() (uint32, bool) { panic(v2.ErrUnsupported) }
