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
		// moved from global, so setting that default here
		BlockBackfillDepth: ptr[uint32](10),
		BlockBackfillSkip:  ptr(false),

		ChainType:                ptr(string(set.chainType)),
		FinalityDepth:            ptr(set.finalityDepth),
		FlagsContractAddress:     asEIP155Address(set.flagsContractAddress),
		LinkContractAddress:      asEIP155Address(set.linkContractAddress),
		LogBackfillBatchSize:     ptr(set.logBackfillBatchSize),
		LogPollInterval:          models.MustNewDuration(set.logPollInterval),
		LogKeepBlocksDepth:       ptr(set.logKeepBlocksDepth),
		MinIncomingConfirmations: ptr(set.minIncomingConfirmations),
		MinContractPayment:       set.minimumContractPayment,
		NonceAutoSync:            ptr(set.nonceAutoSync),
		NoNewHeadsThreshold:      models.MustNewDuration(set.nodeDeadAfterNoNewHeadersThreshold),
		OperatorFactoryAddress:   asEIP155Address(set.operatorFactoryAddress),
		RPCDefaultBatchSize:      ptr(set.rpcDefaultBatchSize),
		RPCBlockQueryDelay:       ptr(set.blockHistoryEstimatorBlockDelay),
		Transactions: &v2.Transactions{
			ForwardersEnabled:    ptr(set.useForwarders),
			MaxInFlight:          ptr(set.maxInFlightTransactions),
			MaxQueued:            ptr(uint32(set.maxQueuedTransactions)),
			ReaperInterval:       models.MustNewDuration(set.ethTxReaperInterval),
			ReaperThreshold:      models.MustNewDuration(set.ethTxReaperThreshold),
			ResendAfterThreshold: models.MustNewDuration(set.ethTxResendAfterThreshold),
		},
		BalanceMonitor: &v2.BalanceMonitor{
			Enabled: ptr(set.balanceMonitorEnabled),
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
			LimitMax:           ptr(uint32(set.gasLimitMax)),
			LimitMultiplier:    ptr(decimal.NewFromFloat32(set.gasLimitMultiplier)),
			LimitTransfer:      ptr(uint32(set.gasLimitTransfer)),
			TipCapDefault:      utils.NewWei(&set.gasTipCapDefault),
			TipCapMin:          utils.NewWei(&set.gasTipCapMinimum),
			PriceDefault:       utils.NewWei(&set.gasPriceDefault),
			PriceMax:           utils.NewWei(&set.maxGasPriceWei),
			PriceMin:           utils.NewWei(&set.minGasPriceWei),
			LimitJobType: &v2.GasLimitJobType{
				OCR:    set.gasLimitOCRJobType,
				DR:     set.gasLimitDRJobType,
				VRF:    set.gasLimitVRFJobType,
				FM:     set.gasLimitFMJobType,
				Keeper: set.gasLimitKeeperJobType,
			},
			BlockHistory: &v2.BlockHistoryEstimator{
				BatchSize:             ptr(set.blockHistoryEstimatorBatchSize),
				BlockHistorySize:      ptr(set.blockHistoryEstimatorBlockHistorySize),
				TransactionPercentile: ptr(set.blockHistoryEstimatorTransactionPercentile),
			},
		},
		HeadTracker: &v2.HeadTracker{
			HistoryDepth:     ptr(set.headTrackerHistoryDepth),
			MaxBufferSize:    ptr(set.headTrackerMaxBufferSize),
			SamplingInterval: models.MustNewDuration(set.headTrackerSamplingInterval),
		},
		KeySpecific: nil,
		NodePool: &v2.NodePool{
			PollFailureThreshold: ptr(set.nodePollFailureThreshold),
			PollInterval:         models.MustNewDuration(set.nodePollInterval),
			SelectionMode:        ptr(set.nodeSelectionMode),
		},
		OCR: &v2.OCR{
			ContractConfirmations:              ptr(set.ocrContractConfirmations),
			ContractTransmitterTransmitTimeout: models.MustNewDuration(set.ocrContractTransmitterTransmitTimeout),
			DatabaseTimeout:                    models.MustNewDuration(set.ocrDatabaseTimeout),
			ObservationGracePeriod:             models.MustNewDuration(set.ocrObservationGracePeriod),
		},
	}
	if *c.ChainType == "" {
		c.ChainType = nil
	}
	if isZeroPtr(c.BalanceMonitor) {
		c.BalanceMonitor = nil
	}
	if isZeroPtr(c.GasEstimator.LimitJobType) {
		c.GasEstimator.LimitJobType = nil
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
