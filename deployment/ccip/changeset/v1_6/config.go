package v1_6

import (
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	DefaultOCRParamsForCommitForNonETH = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           120 * time.Second,
			DeltaResend:                             globals.OcrDeltaResend,
			DeltaInitial:                            globals.OcrDeltaInitial,
			DeltaRound:                              15 * time.Second,
			DeltaGrace:                              globals.OcrDeltaGrace,
			DeltaCertifiedCommitRequest:             globals.OcrDeltaCertifiedCommitRequest,
			DeltaStage:                              25 * time.Second, // TransmissionDelayMultiplier overrides DeltaStage
			Rmax:                                    globals.OcrRMax,
			MaxDurationQuery:                        globals.OcrMaxDurationQuery,
			MaxDurationObservation:                  13 * time.Second,
			MaxDurationShouldAcceptAttestedReport:   globals.OcrMaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
		},
		CommitOffChainConfig: &globals.DefaultCommitOffChainCfg,
	}

	DefaultOCRParamsForExecForNonETH = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           100 * time.Second,
			DeltaResend:                             globals.OcrDeltaResend,
			DeltaInitial:                            globals.OcrDeltaInitial,
			DeltaRound:                              15 * time.Second,
			DeltaGrace:                              globals.OcrDeltaGrace,
			DeltaCertifiedCommitRequest:             globals.OcrDeltaCertifiedCommitRequest,
			DeltaStage:                              25 * time.Second, // TransmissionDelayMultiplier overrides DeltaStage
			Rmax:                                    globals.OcrRMax,
			MaxDurationQuery:                        globals.OcrMaxDurationQuery,
			MaxDurationObservation:                  13 * time.Second,
			MaxDurationShouldAcceptAttestedReport:   globals.OcrMaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
		},
		ExecuteOffChainConfig: &globals.DefaultExecuteOffChainCfg,
	}

	DefaultOCRParamsForCommitForETH = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           120 * time.Second,
			DeltaResend:                             globals.OcrDeltaResend,
			DeltaInitial:                            globals.OcrDeltaInitial,
			DeltaRound:                              90 * time.Second,
			DeltaGrace:                              globals.OcrDeltaGrace,
			DeltaCertifiedCommitRequest:             globals.OcrDeltaCertifiedCommitRequest,
			DeltaStage:                              60 * time.Second, // TransmissionDelayMultiplier overrides DeltaStage
			Rmax:                                    globals.OcrRMax,
			MaxDurationQuery:                        globals.OcrMaxDurationQuery,
			MaxDurationObservation:                  35 * time.Second,
			MaxDurationShouldAcceptAttestedReport:   globals.OcrMaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
		},
		CommitOffChainConfig: &globals.DefaultCommitOffChainCfg,
	}

	DefaultOCRParamsForExecForETH = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           100 * time.Second,
			DeltaResend:                             globals.OcrDeltaResend,
			DeltaInitial:                            globals.OcrDeltaInitial,
			DeltaRound:                              90 * time.Second,
			DeltaGrace:                              globals.OcrDeltaGrace,
			DeltaCertifiedCommitRequest:             globals.OcrDeltaCertifiedCommitRequest,
			DeltaStage:                              60 * time.Second, // TransmissionDelayMultiplier overrides DeltaStage
			Rmax:                                    globals.OcrRMax,
			MaxDurationQuery:                        globals.OcrMaxDurationQuery,
			MaxDurationObservation:                  20 * time.Second,
			MaxDurationShouldAcceptAttestedReport:   globals.OcrMaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: 8 * time.Second,
		},
		ExecuteOffChainConfig: &globals.DefaultExecuteOffChainCfg,
	}

	// Used for only testing with simulated chains
	OcrParamsForTest = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           globals.DeltaProgress,
			DeltaResend:                             globals.DeltaResend,
			DeltaInitial:                            globals.DeltaInitial,
			DeltaRound:                              globals.DeltaRound,
			DeltaGrace:                              globals.DeltaGrace,
			DeltaCertifiedCommitRequest:             globals.DeltaCertifiedCommitRequest,
			DeltaStage:                              globals.DeltaStage,
			Rmax:                                    globals.Rmax,
			MaxDurationQuery:                        globals.MaxDurationQuery,
			MaxDurationObservation:                  globals.MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport:   globals.MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: globals.MaxDurationShouldTransmitAcceptedReport,
		},
		CommitOffChainConfig: &pluginconfig.CommitOffchainConfig{
			RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(globals.RemoteGasPriceBatchWriteFrequency),
			TokenPriceBatchWriteFrequency:      *config.MustNewDuration(globals.TokenPriceBatchWriteFrequency),
			NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
			MaxReportTransmissionCheckAttempts: 5,
			RMNEnabled:                         false,
			RMNSignaturesTimeout:               30 * time.Minute,
			MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
			SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
			MerkleRootAsyncObserverDisabled:    false,
			MerkleRootAsyncObserverSyncFreq:    4 * time.Second,
			MerkleRootAsyncObserverSyncTimeout: 12 * time.Second,
			ChainFeeAsyncObserverSyncFreq:      10 * time.Second,
			ChainFeeAsyncObserverSyncTimeout:   12 * time.Second,
		},
		ExecuteOffChainConfig: &pluginconfig.ExecuteOffchainConfig{
			BatchGasLimit:             globals.BatchGasLimit,
			InflightCacheExpiry:       *config.MustNewDuration(globals.InflightCacheExpiry),
			RootSnoozeTime:            *config.MustNewDuration(globals.RootSnoozeTime),
			MessageVisibilityInterval: *config.MustNewDuration(globals.PermissionLessExecutionThreshold),
			BatchingStrategyID:        globals.BatchingStrategyID,
		},
	}
)

func DeriveOCRParamsForCommit(
	chainsel uint64,
	isSimulatedChain bool,
	feedChain uint64,
	feeTokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	var params CCIPOCRParams
	if isSimulatedChain {
		params = OcrParamsForTest.Copy()
	} else if chainsel == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector ||
		chainsel == chain_selectors.ETHEREUM_MAINNET.Selector {
		params = DefaultOCRParamsForCommitForETH.Copy()
	} else {
		params = DefaultOCRParamsForCommitForNonETH.Copy()
	}
	params.CommitOffChainConfig.TokenInfo = feeTokenInfo
	params.CommitOffChainConfig.PriceFeedChainSelector = ccipocr3.ChainSelector(feedChain)
	if override == nil {
		return params
	}

	return override(params)
}

func DeriveOCRParamsForExec(
	chainsel uint64,
	isSimulatedChain bool,
	observerConfig []pluginconfig.TokenDataObserverConfig,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	var params CCIPOCRParams
	if isSimulatedChain {
		params = OcrParamsForTest.Copy()
	} else if chainsel == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector ||
		chainsel == chain_selectors.ETHEREUM_MAINNET.Selector {
		params = DefaultOCRParamsForExecForETH.Copy()
	} else {
		params = DefaultOCRParamsForExecForNonETH.Copy()
	}

	params.ExecuteOffChainConfig.TokenDataObservers = observerConfig
	if override == nil {
		return params
	}

	return override(params)
}
