package v1_6

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

type OffRampView struct {
	types.ContractMetaData
	DynamicConfig             offramp.OffRampDynamicConfig                `json:"dynamicConfig"`
	LatestPriceSequenceNumber uint64                                      `json:"latestPriceSequenceNumber"`
	SourceChainConfigs        map[uint64]offramp.OffRampSourceChainConfig `json:"sourceChainConfigs"`
	StaticConfig              offramp.OffRampStaticConfig                 `json:"staticConfig"`
}

func GenerateOffRampView(
	offRampContract *offramp.OffRamp,
	routerContract *router1_2.Router,
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

	sourceChainSelectors, err := v1_2.GetRemoteChainSelectors(routerContract)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get source chain selectors: %w", err)
	}
	sourceChainConfigs := make(map[uint64]offramp.OffRampSourceChainConfig)
	for _, sourceChainSelector := range sourceChainSelectors {
		sourceChainConfig, err := offRampContract.GetSourceChainConfig(nil, sourceChainSelector)
		if err != nil {
			return OffRampView{}, fmt.Errorf("failed to get source chain config: %w", err)
		}
		sourceChainConfigs[sourceChainSelector] = sourceChainConfig
	}

	staticConfig, err := offRampContract.GetStaticConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}

	return OffRampView{
		ContractMetaData:          tv,
		DynamicConfig:             dynamicConfig,
		LatestPriceSequenceNumber: latestPriceSequenceNumber,
		SourceChainConfigs:        sourceChainConfigs,
		StaticConfig:              staticConfig,
	}, nil
}
