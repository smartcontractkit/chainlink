package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//TODO doc
func V2Defaults() map[int64]v2.Chain {
	m := map[int64]v2.Chain{}
	for id, set := range chainSpecificConfigDefaultSets {
		m[id] = v2Defaults(set)
	}
	return m
}

func V2Fallback() v2.Chain {
	return v2Defaults(fallbackDefaultSet)
}

func v2Defaults(set chainSpecificConfigDefaultSet) v2.Chain {
	c := v2.Chain{
		BlockBackfillDepth:       nil,
		BlockBackfillSkip:        nil,
		ChainType:                ptr(string(set.chainType)),
		EIP1559DynamicFees:       ptr(set.eip1559DynamicFees),
		FinalityDepth:            ptr(set.finalityDepth),
		FlagsContractAddress:     asEIP155Address(set.flagsContractAddress),
		GasBumpPercent:           ptr(set.gasBumpPercent),
		GasBumpThreshold:         utils.NewWei(new(big.Int).SetUint64(set.gasBumpThreshold)),
		GasBumpTxDepth:           ptr(set.gasBumpTxDepth),
		GasBumpWei:               utils.NewWei(&set.gasBumpWei),
		GasEstimatorMode:         ptr(set.gasEstimatorMode),
		GasFeeCapDefault:         utils.NewWei(&set.gasFeeCapDefault),
		GasLimitDefault:          utils.NewBig(new(big.Int).SetUint64(set.gasLimitDefault)),
		GasLimitMultiplier:       ptr(decimal.NewFromFloat32(set.gasLimitMultiplier)),
		GasLimitTransfer:         utils.NewBig(new(big.Int).SetUint64(set.gasLimitTransfer)),
		GasPriceDefault:          utils.NewWei(&set.gasPriceDefault),
		GasTipCapDefault:         utils.NewWei(&set.gasTipCapDefault),
		GasTipCapMinimum:         utils.NewWei(&set.gasTipCapMinimum),
		LinkContractAddress:      asEIP155Address(set.linkContractAddress),
		LogBackfillBatchSize:     ptr(set.logBackfillBatchSize),
		LogPollInterval:          models.MustNewDuration(set.logPollInterval),
		MaxGasPriceWei:           utils.NewWei(&set.maxGasPriceWei),
		MaxInFlightTransactions:  ptr(set.maxInFlightTransactions),
		MaxQueuedTransactions:    ptr(uint32(set.maxQueuedTransactions)),
		MinGasPriceWei:           utils.NewWei(&set.minGasPriceWei),
		MinIncomingConfirmations: ptr(set.minIncomingConfirmations),
		MinimumContractPayment:   set.minimumContractPayment,
		NonceAutoSync:            ptr(set.nonceAutoSync),
		OperatorFactoryAddress:   asEIP155Address(set.operatorFactoryAddress),
		RPCDefaultBatchSize:      ptr(set.rpcDefaultBatchSize),
		TxReaperInterval:         models.MustNewDuration(set.ethTxReaperInterval),
		TxReaperThreshold:        models.MustNewDuration(set.ethTxReaperThreshold),
		TxResendAfterThreshold:   models.MustNewDuration(set.ethTxResendAfterThreshold),
		UseForwarders:            ptr(set.useForwarders),
		BalanceMonitor: &v2.BalanceMonitor{
			Enabled:    ptr(set.balanceMonitorEnabled),
			BlockDelay: ptr(set.balanceMonitorBlockDelay),
		},
		BlockHistoryEstimator: &v2.BlockHistoryEstimator{
			BatchSize:                 ptr(set.blockHistoryEstimatorBatchSize),
			BlockDelay:                ptr(set.blockHistoryEstimatorBlockDelay),
			BlockHistorySize:          ptr(set.blockHistoryEstimatorBlockHistorySize),
			EIP1559FeeCapBufferBlocks: set.blockHistoryEstimatorEIP1559FeeCapBufferBlocks,
			TransactionPercentile:     ptr(set.blockHistoryEstimatorTransactionPercentile),
		},
		HeadTracker: &v2.HeadTracker{
			BlockEmissionIdleWarningThreshold: models.MustNewDuration(set.blockEmissionIdleWarningThreshold),
			HistoryDepth:                      ptr(set.headTrackerHistoryDepth),
			MaxBufferSize:                     ptr(set.headTrackerMaxBufferSize),
			SamplingInterval:                  models.MustNewDuration(set.headTrackerSamplingInterval),
		},
		KeySpecific: nil,
		NodePool: &v2.NodePool{
			NoNewHeadsThreshold:  models.MustNewDuration(set.nodeDeadAfterNoNewHeadersThreshold),
			PollFailureThreshold: ptr(set.nodePollFailureThreshold),
			PollInterval:         models.MustNewDuration(set.nodePollInterval),
		},
		OCR: &v2.OCR{
			ContractConfirmations:              ptr(set.ocrContractConfirmations),
			ContractTransmitterTransmitTimeout: models.MustNewDuration(set.ocrContractTransmitterTransmitTimeout),
			DatabaseTimeout:                    models.MustNewDuration(set.ocrDatabaseTimeout),
			ObservationTimeout:                 nil,
			ObservationGracePeriod:             models.MustNewDuration(set.ocrObservationGracePeriod),
		},
		OCR2: &v2.OCR2{
			ContractConfirmations: nil,
		},
	}
	if *c.ChainType == "" {
		c.ChainType = nil
	}
	if isZeroPtr(c.BalanceMonitor) {
		c.BalanceMonitor = nil
	}
	if isZeroPtr(c.BlockHistoryEstimator) {
		c.BlockHistoryEstimator = nil
	}
	if isZeroPtr(c.HeadTracker) {
		c.HeadTracker = nil
	}
	if isZeroPtr(c.NodePool) {
		c.NodePool = nil
	}
	if isZeroPtr(c.OCR) {
		c.OCR = nil
	}
	if isZeroPtr(c.OCR2) {
		c.OCR2 = nil
	}
	return c
}

func ptr[T any](v T) *T {
	return &v
}

func isZeroPtr[T comparable](p *T) bool {
	var t T
	return p == nil || *p == t
}

func asEIP155Address(s string) *ethkey.EIP55Address {
	if s == "" {
		return nil
	}
	a := ethkey.EIP55AddressFromAddress(common.HexToAddress(s))
	return &a
}
