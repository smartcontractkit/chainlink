package chainlink

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
)

func (l *legacyGeneralConfig) EVMRPCEnabled() bool { panic("unimplemented") }

func (l *legacyGeneralConfig) GlobalBalanceMonitorEnabled() (bool, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalChainType() (string, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmEIP1559DynamicFees() (bool, bool)    { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmFinalityDepth() (uint32, bool)       { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasBumpPercent() (uint16, bool)      { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasBumpThreshold() (uint64, bool)    { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasBumpTxDepth() (uint16, bool)      { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasBumpWei() (*big.Int, bool)        { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasFeeCapDefault() (*big.Int, bool)  { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitDefault() (uint64, bool)     { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitMultiplier() (float32, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitTransfer() (uint64, bool)    { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasPriceDefault() (*big.Int, bool)   { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasTipCapDefault() (*big.Int, bool)  { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasTipCapMinimum() (*big.Int, bool)  { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmLogPollInterval() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmMaxGasPriceWei() (*big.Int, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmMinGasPriceWei() (*big.Int, bool)      { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmNonceAutoSync() (bool, bool)           { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmUseForwarders() (bool, bool)           { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool)   { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalFlagsContractAddress() (string, bool)     { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalGasEstimatorMode() (string, bool)         { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalLinkContractAddress() (string, bool)      { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalOperatorFactoryAddress() (string, bool)   { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalMinIncomingConfirmations() (uint32, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalNodeNoNewHeadsThreshold() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalNodePollFailureThreshold() (uint32, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalNodePollInterval() (time.Duration, bool)  { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalOCRContractConfirmations() (uint16, bool) { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalOCRDatabaseTimeout() (time.Duration, bool) {
	panic("unimplemented")
}
func (l *legacyGeneralConfig) GlobalOCRObservationGracePeriod() (time.Duration, bool) {
	panic("unimplemented")
}

func (l *legacyGeneralConfig) GlobalEvmGasLimitOCRJobType() (uint64, bool)    { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitDRJobType() (uint64, bool)     { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitVRFJobType() (uint64, bool)    { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitFMJobType() (uint64, bool)     { panic("unimplemented") }
func (l *legacyGeneralConfig) GlobalEvmGasLimitKeeperJobType() (uint64, bool) { panic("unimplemented") }
