package aptos

import (
	"encoding/hex"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type OffRampView struct {
	aptosCommon.ContractMetaData

	LatestPriceSequenceNumber uint64
	StaticConfig              OffRampStaticConfig
	DynamicConfig             OffRampDynamicConfig
	SourceChainConfigs        map[uint64]OffRampSourceChainConfig
}

type OffRampStaticConfig struct {
	ChainSelector      uint64
	RMNRemote          string
	TokenAdminRegistry string
	NonceManager       string
}

type OffRampDynamicConfig struct {
	FeeQuoter                               string
	PermissionlessExecutionThresholdSeconds uint32
}

type OffRampSourceChainConfig struct {
	Router                    string
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRMNVerificationDisabled bool
	OnRamp                    string
}

func GenerateOffRampView(chain cldf_aptos.Chain, offRampAddress aptos.AccountAddress, routerAddress aptos.AccountAddress) (OffRampView, error) {
	boundOffRamp := ccip_offramp.Bind(offRampAddress, chain.Client)
	boundRouter := ccip_router.Bind(routerAddress, chain.Client)

	owner, err := boundOffRamp.Offramp().Owner(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to retrieve owner of offRamp %s: %w", routerAddress.StringLong(), err)
	}

	latestPriceSequenceNumber, err := boundOffRamp.Offramp().GetLatestPriceSequenceNumber(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get latest price sequence number for offRamp %s: %w", offRampAddress.StringLong(), err)
	}

	destChainSelectors, err := boundRouter.Router().GetDestChains(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get dest chain selectors: %w", err)
	}

	staticConfig, err := boundOffRamp.Offramp().GetStaticConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}
	dynamicConfig, err := boundOffRamp.Offramp().GetDynamicConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	sourceChainConfigs := make(map[uint64]OffRampSourceChainConfig, len(destChainSelectors))
	for _, destChainSelector := range destChainSelectors {
		sourceChainConfig, err := boundOffRamp.Offramp().GetSourceChainConfig(nil, destChainSelector)
		if err != nil {
			return OffRampView{}, fmt.Errorf("failed to get source chain config for chain %v: %w", destChainSelector, err)
		}
		sourceChainConfigs[destChainSelector] = OffRampSourceChainConfig{
			Router:                    sourceChainConfig.Router.StringLong(),
			IsEnabled:                 sourceChainConfig.IsEnabled,
			MinSeqNr:                  sourceChainConfig.MinSeqNr,
			IsRMNVerificationDisabled: sourceChainConfig.IsRMNVerificationDisabled,
			OnRamp:                    hex.EncodeToString(sourceChainConfig.OnRamp),
		}
	}

	return OffRampView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: offRampAddress.StringLong(),
			Owner:   owner.StringLong(),
		},
		LatestPriceSequenceNumber: latestPriceSequenceNumber,
		StaticConfig: OffRampStaticConfig{
			ChainSelector:      staticConfig.ChainSelector,
			RMNRemote:          staticConfig.RMNRemote.StringLong(),
			TokenAdminRegistry: staticConfig.TokenAdminRegistry.StringLong(),
			NonceManager:       staticConfig.NonceManager.StringLong(),
		},
		DynamicConfig: OffRampDynamicConfig{
			FeeQuoter:                               dynamicConfig.FeeQuoter.StringLong(),
			PermissionlessExecutionThresholdSeconds: dynamicConfig.PermissionlessExecutionThresholdSeconds,
		},
		SourceChainConfigs: sourceChainConfigs,
	}, nil
}
