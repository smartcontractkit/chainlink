package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type OnRampView struct {
	aptosCommon.ContractMetaData

	StaticConfig          OnRampStaticConfig
	DynamicConfig         OnRampDynamicConfig
	SourceTokenToPool     map[string]string
	DestChainSpecificData map[uint64]DestChainSpecificData
}

type OnRampStaticConfig struct {
	ChainSelector uint64
}

type OnRampDynamicConfig struct {
	FeeAggregator  string
	AllowlistAdmin string
}

type DestChainSpecificData struct {
	AllowedSendersList []string
	DestChainConfig    OnRampDestChainConfig
	ExpectedNextSeqNum uint64
}

type OnRampDestChainConfig struct {
	SequenceNumber   uint64
	AllowlistEnabled bool
	Router           string
}

func GenerateOnRampView(
	chain cldf_aptos.Chain,
	onRampAddress aptos.AccountAddress,
	routerAddress aptos.AccountAddress,
	ccipAddress aptos.AccountAddress,
) (OnRampView, error) {
	boundOnRamp := ccip_onramp.Bind(onRampAddress, chain.Client)
	boundRouter := ccip_router.Bind(routerAddress, chain.Client)
	boundCCIP := ccip.Bind(ccipAddress, chain.Client)

	destinationChainSelectors, err := boundRouter.Router().GetDestChains(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get dest chain selectors: %w", err)
	}

	owner, err := boundOnRamp.Onramp().Owner(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get owner for onramp %v: %w", onRampAddress.StringLong(), err)
	}
	staticConfig, err := boundOnRamp.Onramp().GetStaticConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get static config for onramp %v: %w", onRampAddress.StringLong(), err)
	}
	dynamicConfig, err := boundOnRamp.Onramp().GetDynamicConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get dynamic config for onramp %v: %w", onRampAddress.StringLong(), err)
	}

	nodeInfo, err := chain.Client.Info()
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get node info from Aptos client: %w", err)
	}
	ledgerVersion := nodeInfo.LedgerVersion()
	callOpts := &bind.CallOpts{
		LedgerVersion: &ledgerVersion,
	}
	var (
		sourceTokens       []aptos.AccountAddress
		nextKey            aptos.AccountAddress
		next               = true
		sourceTokensToPool = make(map[string]string)
	)
	for next {
		sourceTokens, nextKey, next, err = boundCCIP.TokenAdminRegistry().GetAllConfiguredTokens(callOpts, nextKey, 100)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get all configured tokens from tokenAdminRegistry %v: %w", ccipAddress.StringLong(), err)
		}
		pools, err := boundCCIP.TokenAdminRegistry().GetPools(callOpts, sourceTokens)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get pools from tokenAdminRegistry %v: %w", ccipAddress.StringLong(), err)
		}
		for i, pool := range pools {
			sourceTokensToPool[sourceTokens[i].StringLong()] = pool.StringLong()
		}
	}

	destChainSpecificData := make(map[uint64]DestChainSpecificData, len(destinationChainSelectors))
	for _, selector := range destinationChainSelectors {
		expectedNextSequenceNumber, err := boundOnRamp.Onramp().GetExpectedNextSequenceNumber(nil, selector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get expected next sequence number for selector %v: %w", selector, err)
		}
		sequenceNumber, allowlistEnabled, routerAddr, err := boundOnRamp.Onramp().GetDestChainConfig(nil, selector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get dest chain config for selector %v: %w", selector, err)
		}
		_, allowedSenders, err := boundOnRamp.Onramp().GetAllowedSendersList(nil, selector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get allowlist for selector %v: %w", selector, err)
		}
		allowedSenderStrings := make([]string, len(allowedSenders))
		for i, allowedSender := range allowedSenders {
			allowedSenderStrings[i] = allowedSender.StringLong()
		}
		destChainSpecificData[selector]  = DestChainSpecificData{
			AllowedSendersList: allowedSenderStrings,
			DestChainConfig:    OnRampDestChainConfig{
				SequenceNumber:   sequenceNumber,
				AllowlistEnabled: allowlistEnabled,
				Router:           routerAddr.StringLong(),
			},
			ExpectedNextSeqNum: expectedNextSequenceNumber,
		}
	}

	return OnRampView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: onRampAddress.StringLong(),
			Owner:   owner.StringLong(),
		},
		StaticConfig: OnRampStaticConfig{
			ChainSelector: staticConfig.ChainSelector,
		},
		DynamicConfig: OnRampDynamicConfig{
			FeeAggregator:  dynamicConfig.FeeAggregator.StringLong(),
			AllowlistAdmin: dynamicConfig.AllowlistAdmin.StringLong(),
		},
		SourceTokenToPool: sourceTokensToPool,
		DestChainSpecificData:destChainSpecificData,
	}, nil
}
