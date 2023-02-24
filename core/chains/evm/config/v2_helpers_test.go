package config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		BlockBackfillDepth: testutils.Ptr[uint32](10),
		BlockBackfillSkip:  testutils.Ptr(false),

		ChainType:                testutils.Ptr(string(set.chainType)),
		FinalityDepth:            testutils.Ptr(set.finalityDepth),
		FlagsContractAddress:     asEIP155Address(set.flagsContractAddress),
		LinkContractAddress:      asEIP155Address(set.linkContractAddress),
		LogBackfillBatchSize:     testutils.Ptr(set.logBackfillBatchSize),
		LogPollInterval:          models.MustNewDuration(set.logPollInterval),
		LogKeepBlocksDepth:       testutils.Ptr(set.logKeepBlocksDepth),
		MinIncomingConfirmations: testutils.Ptr(set.minIncomingConfirmations),
		MinContractPayment:       set.minimumContractPayment,
		NonceAutoSync:            testutils.Ptr(set.nonceAutoSync),
		NoNewHeadsThreshold:      models.MustNewDuration(set.nodeDeadAfterNoNewHeadersThreshold),
		OperatorFactoryAddress:   asEIP155Address(set.operatorFactoryAddress),
		RPCDefaultBatchSize:      testutils.Ptr(set.rpcDefaultBatchSize),
		RPCBlockQueryDelay:       testutils.Ptr(set.blockHistoryEstimatorBlockDelay),
		Transactions: v2.Transactions{
			ForwardersEnabled:    testutils.Ptr(set.useForwarders),
			MaxInFlight:          testutils.Ptr(set.maxInFlightTransactions),
			MaxQueued:            testutils.Ptr(uint32(set.maxQueuedTransactions)),
			ReaperInterval:       models.MustNewDuration(set.ethTxReaperInterval),
			ReaperThreshold:      models.MustNewDuration(set.ethTxReaperThreshold),
			ResendAfterThreshold: models.MustNewDuration(set.ethTxResendAfterThreshold),
		},
		BalanceMonitor: v2.BalanceMonitor{
			Enabled: testutils.Ptr(set.balanceMonitorEnabled),
		},
		GasEstimator: v2.GasEstimator{
			Mode:               testutils.Ptr(set.gasEstimatorMode),
			EIP1559DynamicFees: testutils.Ptr(set.eip1559DynamicFees),
			BumpMin:            &set.gasBumpWei,
			BumpPercent:        testutils.Ptr(set.gasBumpPercent),
			BumpThreshold:      testutils.Ptr(uint32(set.gasBumpThreshold)),
			BumpTxDepth:        testutils.Ptr(set.gasBumpTxDepth),
			FeeCapDefault:      &set.gasFeeCapDefault,
			LimitDefault:       testutils.Ptr(uint32(set.gasLimitDefault)),
			LimitMax:           testutils.Ptr(uint32(set.gasLimitMax)),
			LimitMultiplier:    testutils.Ptr(decimal.NewFromFloat32(set.gasLimitMultiplier)),
			LimitTransfer:      testutils.Ptr(uint32(set.gasLimitTransfer)),
			TipCapDefault:      &set.gasTipCapDefault,
			TipCapMin:          &set.gasTipCapMinimum,
			PriceDefault:       &set.gasPriceDefault,
			PriceMax:           &set.maxGasPriceWei,
			PriceMin:           &set.minGasPriceWei,
			LimitJobType: v2.GasLimitJobType{
				OCR:    set.gasLimitOCRJobType,
				DR:     set.gasLimitDRJobType,
				VRF:    set.gasLimitVRFJobType,
				FM:     set.gasLimitFMJobType,
				Keeper: set.gasLimitKeeperJobType,
			},
			BlockHistory: v2.BlockHistoryEstimator{
				BatchSize:                testutils.Ptr(set.blockHistoryEstimatorBatchSize),
				BlockHistorySize:         testutils.Ptr(set.blockHistoryEstimatorBlockHistorySize),
				CheckInclusionBlocks:     testutils.Ptr(set.blockHistoryEstimatorCheckInclusionBlocks),
				CheckInclusionPercentile: testutils.Ptr(set.blockHistoryEstimatorCheckInclusionPercentile),
				TransactionPercentile:    testutils.Ptr(set.blockHistoryEstimatorTransactionPercentile),
			},
		},
		HeadTracker: v2.HeadTracker{
			HistoryDepth:     testutils.Ptr(set.headTrackerHistoryDepth),
			MaxBufferSize:    testutils.Ptr(set.headTrackerMaxBufferSize),
			SamplingInterval: models.MustNewDuration(set.headTrackerSamplingInterval),
		},
		KeySpecific: nil,
		NodePool: v2.NodePool{
			PollFailureThreshold: testutils.Ptr(set.nodePollFailureThreshold),
			PollInterval:         models.MustNewDuration(set.nodePollInterval),
			SelectionMode:        testutils.Ptr(set.nodeSelectionMode),
			SyncThreshold:        testutils.Ptr(set.nodeSyncThreshold),
		},
		OCR: v2.OCR{
			ContractConfirmations:              testutils.Ptr(set.ocrContractConfirmations),
			ContractTransmitterTransmitTimeout: models.MustNewDuration(set.ocrContractTransmitterTransmitTimeout),
			DatabaseTimeout:                    models.MustNewDuration(set.ocrDatabaseTimeout),
			ObservationGracePeriod:             models.MustNewDuration(set.ocrObservationGracePeriod),
		},
		OCR2: v2.OCR2{
			Automation: v2.Automation{
				GasLimit: testutils.Ptr(set.ocr2AutomationGasLimit),
			},
		},
	}
	if *c.ChainType == "" {
		c.ChainType = nil
	}
	return c
}

func asEIP155Address(s string) *ethkey.EIP55Address {
	if s == "" {
		return nil
	}
	a := ethkey.EIP55AddressFromAddress(common.HexToAddress(s))
	return &a
}
