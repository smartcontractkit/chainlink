package v1_6

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
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
)

type OCRConfigChainType int

const (
	Default OCRConfigChainType = iota + 1
	Ethereum
)

func DeriveOCRConfigTypeFromSelector(chainsel uint64) OCRConfigChainType {
	switch chainsel {
	case chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
		chain_selectors.ETHEREUM_TESTNET_HOLESKY.Selector,
		chain_selectors.ETHEREUM_MAINNET.Selector:
		return Ethereum
	default:
		return Default
	}
}

func (c OCRConfigChainType) CommitOCRParams() CCIPOCRParams {
	switch c {
	case Ethereum:
		return DefaultOCRParamsForCommitForETH.Copy()
	case Default:
		return DefaultOCRParamsForCommitForNonETH.Copy()
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
