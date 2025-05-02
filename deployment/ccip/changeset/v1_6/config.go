package v1_6

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	DefaultOCRParamsForCommitForNonETH = CCIPOCRParams{
		OCRParameters:        globals.CommitOCRParams,
		CommitOffChainConfig: &globals.DefaultCommitOffChainCfg,
	}

	DefaultOCRParamsForCommitForETH = CCIPOCRParams{
		OCRParameters:        globals.CommitOCRParamsForEthereum,
		CommitOffChainConfig: &globals.DefaultCommitOffChainCfg,
	}

	DefaultOCRParamsForExecForNonETH = CCIPOCRParams{
		OCRParameters:         globals.ExecOCRParams,
		ExecuteOffChainConfig: &globals.DefaultExecuteOffChainCfg,
	}

	DefaultOCRParamsForExecForETH = CCIPOCRParams{
		OCRParameters:         globals.ExecOCRParamsForEthereum,
		ExecuteOffChainConfig: &globals.DefaultExecuteOffChainCfg,
	}

	// Used for only testing with simulated chains
	OcrParamsForTest = CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           30 * time.Second, // Lower DeltaProgress can lead to timeouts when running tests locally
			DeltaResend:                             10 * time.Second,
			DeltaInitial:                            20 * time.Second,
			DeltaRound:                              2 * time.Second,
			DeltaGrace:                              2 * time.Second,
			DeltaCertifiedCommitRequest:             10 * time.Second,
			DeltaStage:                              10 * time.Second,
			Rmax:                                    50,
			MaxDurationQuery:                        10 * time.Second,
			MaxDurationObservation:                  10 * time.Second,
			MaxDurationShouldAcceptAttestedReport:   10 * time.Second,
			MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
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

type OCRConfigChainType int

const (
	Default OCRConfigChainType = iota + 1
	Ethereum
	// SimulationTest is kept only for backward compatibility. Tests probably should
	// migrate to using Default or Ethereum
	SimulationTest
)

func (c OCRConfigChainType) CommitOCRParams() CCIPOCRParams {
	switch c {
	case Ethereum:
		return DefaultOCRParamsForCommitForETH.Copy()
	case Default:
		return DefaultOCRParamsForCommitForNonETH.Copy()
	case SimulationTest:
		return OcrParamsForTest.Copy()
	default:
		panic("unknown OCRConfigChainType")
	}
}

func (c OCRConfigChainType) ExecuteOCRParams() CCIPOCRParams {
	switch c {
	case Ethereum:
		return DefaultOCRParamsForExecForETH.Copy()
	case Default:
		return DefaultOCRParamsForExecForNonETH.Copy()
	case SimulationTest:
		return OcrParamsForTest.Copy()
	default:
		panic("unknown OCRConfigChainType")
	}
}

func DeriveOCRParamsForCommit(
	ocrChainType OCRConfigChainType,
	feedChain uint64,
	feeTokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	params := ocrChainType.CommitOCRParams()
	params.CommitOffChainConfig.TokenInfo = feeTokenInfo
	params.CommitOffChainConfig.PriceFeedChainSelector = ccipocr3.ChainSelector(feedChain)
	if override == nil {
		return params
	}
	return override(params)
}

func DeriveOCRParamsForExec(
	ocrChainType OCRConfigChainType,
	observerConfig []pluginconfig.TokenDataObserverConfig,
	override func(params CCIPOCRParams) CCIPOCRParams,
) CCIPOCRParams {
	params := ocrChainType.ExecuteOCRParams()
	params.ExecuteOffChainConfig.TokenDataObservers = observerConfig
	if override == nil {
		return params
	}
	return override(params)
}
