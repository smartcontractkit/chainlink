package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	router1_2 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type OffRampView struct {
	types.ContractMetaData
	DynamicConfig                       offramp.OffRampDynamicConfig        `json:"dynamicConfig"`
	LatestPriceSequenceNumber           uint64                              `json:"latestPriceSequenceNumber"`
	SourceChainConfigs                  map[uint64]OffRampSourceChainConfig `json:"sourceChainConfigs"`
	SourceChainConfigsBasedOnTestRouter map[uint64]OffRampSourceChainConfig `json:"sourceChainConfigsBasedOnTestRouter"`
	StaticConfig                        offramp.OffRampStaticConfig         `json:"staticConfig"`
}

type OffRampSourceChainConfig struct {
	Router                    common.Address
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRMNVerificationDisabled bool
	OnRamp                    string
}

func GenerateOffRampView(
	offRampContract offramp.OffRampInterface,
	routerContract, testRouterContract *router1_2.Router,
) (OffRampView, error) {
	tv, err := types.NewContractMetaData(offRampContract, offRampContract.Address())
	if err != nil {
		return OffRampView{}, err
	}

	dynamicConfig, err := offRampContract.GetDynamicConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	latestPriceSequenceNumber, err := offRampContract.GetLatestPriceSequenceNumber(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get latest price sequence number: %w", err)
	}

	staticConfig, err := offRampContract.GetStaticConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}
	sourceChainConfigs, err := generateSourceChainConfigs(offRampContract, routerContract)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get source chain configs: %w", err)
	}
	var testSourceChainConfigs map[uint64]OffRampSourceChainConfig
	if testRouterContract != nil {
		testSourceChainConfigs, err = generateSourceChainConfigs(offRampContract, testRouterContract)
		if err != nil {
			return OffRampView{}, fmt.Errorf("failed to get test source chain configs: %w", err)
		}
	}
	return OffRampView{
		ContractMetaData:                    tv,
		DynamicConfig:                       dynamicConfig,
		LatestPriceSequenceNumber:           latestPriceSequenceNumber,
		SourceChainConfigs:                  sourceChainConfigs,
		SourceChainConfigsBasedOnTestRouter: testSourceChainConfigs,
		StaticConfig:                        staticConfig,
	}, nil
}

func generateSourceChainConfigs(offRampContract offramp.OffRampInterface, routerContract *router1_2.Router) (map[uint64]OffRampSourceChainConfig, error) {
	sourceChainSelectors, err := v1_2.GetRemoteChainSelectors(routerContract)
	if err != nil {
		return nil, fmt.Errorf("failed to get source chain selectors: %w", err)
	}
	sourceChainConfigs := make(map[uint64]OffRampSourceChainConfig)
	for _, sourceChainSelector := range sourceChainSelectors {
		sourceChainConfig, err := offRampContract.GetSourceChainConfig(nil, sourceChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to get source chain config: %w", err)
		}
		sourceChainConfigs[sourceChainSelector] = OffRampSourceChainConfig{
			Router:                    sourceChainConfig.Router,
			IsEnabled:                 sourceChainConfig.IsEnabled,
			MinSeqNr:                  sourceChainConfig.MinSeqNr,
			IsRMNVerificationDisabled: sourceChainConfig.IsRMNVerificationDisabled,
			OnRamp:                    shared.GetAddressFromBytes(sourceChainSelector, sourceChainConfig.OnRamp),
		}
	}
	return sourceChainConfigs, nil
}
