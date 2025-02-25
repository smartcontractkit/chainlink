package v1_6

import (
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

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
		CommitOffChainConfig: globals.DefaultCommitOffChainCfg,
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
		ExecuteOffChainConfig: globals.DefaultExecuteOffChainCfg,
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
		CommitOffChainConfig: globals.DefaultCommitOffChainCfg,
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
		ExecuteOffChainConfig: globals.DefaultExecuteOffChainCfg,
	}
)

func DeriveOCRParamsForCommit(
	chainsel uint64,
	feedChain uint64,
	feeTokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	params := DefaultOCRParamsForCommitForNonETH
	if chainsel == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector ||
		chainsel == chain_selectors.ETHEREUM_MAINNET.Selector {
		params = DefaultOCRParamsForCommitForETH
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
	observerConfig []pluginconfig.TokenDataObserverConfig,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	params := DefaultOCRParamsForExecForNonETH
	if chainsel == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector ||
		chainsel == chain_selectors.ETHEREUM_MAINNET.Selector {
		params = DefaultOCRParamsForExecForETH
	}
	params.ExecuteOffChainConfig.TokenDataObservers = observerConfig
	if override == nil {
		return params
	}

	return override(params)
}
