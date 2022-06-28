package config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func ChainSpecificConfigDefaultsAsV2() map[int64]v2.Chain {
	m := map[int64]v2.Chain{}
	for id, set := range chainSpecificConfigDefaultSets {
		m[id] = set.asV2()
	}
	return m
}

func FallbackDefaultsAsV2() v2.Chain {
	return fallbackDefaultSet.asV2()
}

func (set chainSpecificConfigDefaultSet) asV2() v2.Chain {
	c := v2.Chain{
		BlockBackfillDepth:       nil,
		BlockBackfillSkip:        nil,
		ChainType:                ptr(string(set.chainType)),
		FinalityDepth:            ptr(set.finalityDepth),
		FlagsContractAddress:     asEIP155Address(set.flagsContractAddress),
		LinkContractAddress:      asEIP155Address(set.linkContractAddress),
		LogBackfillBatchSize:     ptr(set.logBackfillBatchSize),
		LogPollInterval:          models.MustNewDuration(set.logPollInterval),
		MaxInFlightTransactions:  ptr(set.maxInFlightTransactions),
		MaxQueuedTransactions:    ptr(uint32(set.maxQueuedTransactions)),
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
		GasEstimator: &v2.GasEstimator{
			Mode:               ptr(set.gasEstimatorMode),
			EIP1559DynamicFees: ptr(set.eip1559DynamicFees),
			BumpMin:            utils.NewWei(&set.gasBumpWei),
			BumpPercent:        ptr(set.gasBumpPercent),
			BumpThreshold:      ptr(uint32(set.gasBumpThreshold)),
			BumpTxDepth:        ptr(set.gasBumpTxDepth),
			FeeCapDefault:      utils.NewWei(&set.gasFeeCapDefault),
			LimitDefault:       ptr(uint32(set.gasLimitDefault)),
			LimitMultiplier:    ptr(decimal.NewFromFloat32(set.gasLimitMultiplier)),
			LimitTransfer:      ptr(uint32(set.gasLimitTransfer)),
			TipCapDefault:      utils.NewWei(&set.gasTipCapDefault),
			TipCapMinimum:      utils.NewWei(&set.gasTipCapMinimum),
			PriceDefault:       utils.NewWei(&set.gasPriceDefault),
			PriceMax:           utils.NewWei(&set.maxGasPriceWei),
			PriceMin:           utils.NewWei(&set.minGasPriceWei),
			BlockHistory: &v2.BlockHistoryEstimator{
				BatchSize:                 ptr(set.blockHistoryEstimatorBatchSize),
				BlockDelay:                ptr(set.blockHistoryEstimatorBlockDelay),
				BlockHistorySize:          ptr(set.blockHistoryEstimatorBlockHistorySize),
				EIP1559FeeCapBufferBlocks: set.blockHistoryEstimatorEIP1559FeeCapBufferBlocks,
				TransactionPercentile:     ptr(set.blockHistoryEstimatorTransactionPercentile),
			},
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
	}
	if *c.ChainType == "" {
		c.ChainType = nil
	}
	if isZeroPtr(c.BalanceMonitor) {
		c.BalanceMonitor = nil
	}
	if isZeroPtr(c.GasEstimator.BlockHistory) {
		c.GasEstimator.BlockHistory = nil
	}
	if isZeroPtr(c.GasEstimator) {
		c.GasEstimator = nil
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
