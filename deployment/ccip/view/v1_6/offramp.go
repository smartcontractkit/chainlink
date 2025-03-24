package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
)

type OffRampView struct {
	types.ContractMetaData
	DynamicConfig             offramp.OffRampDynamicConfig        `json:"dynamicConfig"`
	LatestPriceSequenceNumber uint64                              `json:"latestPriceSequenceNumber"`
	SourceChainConfigs        map[uint64]OffRampSourceChainConfig `json:"sourceChainConfigs"`
	StaticConfig              offramp.OffRampStaticConfig         `json:"staticConfig"`
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
	sourceChainConfigs := make(map[uint64]OffRampSourceChainConfig)
	for _, sourceChainSelector := range sourceChainSelectors {
		sourceChainConfig, err := offRampContract.GetSourceChainConfig(nil, sourceChainSelector)
		if err != nil {
			return OffRampView{}, fmt.Errorf("failed to get source chain config: %w", err)
		}
		sourceChainConfigs[sourceChainSelector] = OffRampSourceChainConfig{
			Router:                    sourceChainConfig.Router,
			IsEnabled:                 sourceChainConfig.IsEnabled,
			MinSeqNr:                  sourceChainConfig.MinSeqNr,
			IsRMNVerificationDisabled: sourceChainConfig.IsRMNVerificationDisabled,
			OnRamp:                    ccipocr3.UnknownAddress(sourceChainConfig.OnRamp).String(),
		}
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
