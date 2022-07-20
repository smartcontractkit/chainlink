package chainlink

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
)

func (g *generalConfig) EVMRPCEnabled() bool { panic("unimplemented") }

func (g *generalConfig) GlobalBalanceMonitorEnabled() (bool, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalChainType() (string, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmEIP1559DynamicFees() (bool, bool)    { panic("unimplemented") }
func (g *generalConfig) GlobalEvmFinalityDepth() (uint32, bool)       { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasBumpPercent() (uint16, bool)      { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasBumpThreshold() (uint64, bool)    { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasBumpTxDepth() (uint16, bool)      { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasBumpWei() (*big.Int, bool)        { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasFeeCapDefault() (*big.Int, bool)  { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitDefault() (uint64, bool)     { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitMultiplier() (float32, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitTransfer() (uint64, bool)    { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasPriceDefault() (*big.Int, bool)   { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasTipCapDefault() (*big.Int, bool)  { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasTipCapMinimum() (*big.Int, bool)  { panic("unimplemented") }
func (g *generalConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalEvmLogPollInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmMaxGasPriceWei() (*big.Int, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalEvmMinGasPriceWei() (*big.Int, bool)      { panic("unimplemented") }
func (g *generalConfig) GlobalEvmNonceAutoSync() (bool, bool)           { panic("unimplemented") }
func (g *generalConfig) GlobalEvmUseForwarders() (bool, bool)           { panic("unimplemented") }
func (g *generalConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool)   { panic("unimplemented") }
func (g *generalConfig) GlobalFlagsContractAddress() (string, bool)     { panic("unimplemented") }
func (g *generalConfig) GlobalGasEstimatorMode() (string, bool)         { panic("unimplemented") }
func (g *generalConfig) GlobalLinkContractAddress() (string, bool)      { panic("unimplemented") }
func (g *generalConfig) GlobalOperatorFactoryAddress() (string, bool)   { panic("unimplemented") }
func (g *generalConfig) GlobalMinIncomingConfirmations() (uint32, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalNodeNoNewHeadsThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalNodePollFailureThreshold() (uint32, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalNodePollInterval() (time.Duration, bool)  { panic("unimplemented") }
func (g *generalConfig) GlobalOCRContractConfirmations() (uint16, bool) { panic("unimplemented") }
func (g *generalConfig) GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalOCRDatabaseTimeout() (time.Duration, bool) {
	panic("unimplemented")
}
func (g *generalConfig) GlobalOCRObservationGracePeriod() (time.Duration, bool) {
	panic("unimplemented")
}

func (g *generalConfig) GlobalEvmGasLimitOCRJobType() (uint64, bool)    { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitDRJobType() (uint64, bool)     { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitVRFJobType() (uint64, bool)    { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitFMJobType() (uint64, bool)     { panic("unimplemented") }
func (g *generalConfig) GlobalEvmGasLimitKeeperJobType() (uint64, bool) { panic("unimplemented") }
