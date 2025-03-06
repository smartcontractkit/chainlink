package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
)

type OnRampView struct {
	types.ContractMetaData
	DynamicConfig         onramp.OnRampDynamicConfig        `json:"dynamicConfig"`
	StaticConfig          onramp.OnRampStaticConfig         `json:"staticConfig"`
	SourceTokenToPool     map[common.Address]common.Address `json:"sourceTokenToPool"`
	DestChainSpecificData map[uint64]DestChainSpecificData  `json:"destChainSpecificData"`
}

type DestChainSpecificData struct {
	AllowedSendersList []common.Address          `json:"allowedSendersList"`
	DestChainConfig    onramp.GetDestChainConfig `json:"destChainConfig"`
	ExpectedNextSeqNum uint64                    `json:"expectedNextSeqNum"`
}

func GenerateOnRampView(
	onRampContract onramp.OnRampInterface,
	routerContract *router1_2.Router,
	taContract *token_admin_registry.TokenAdminRegistry,
) (OnRampView, error) {
	tv, err := types.NewContractMetaData(onRampContract, onRampContract.Address())
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get contract metadata: %w", err)
	}
	dynamicConfig, err := onRampContract.GetDynamicConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	staticConfig, err := onRampContract.GetStaticConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}

	// populate destChainSelectors from router
	destChainSelectors, err := v1_2.GetRemoteChainSelectors(routerContract)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get destination selectors: %w", err)
	}
	// populate sourceTokens from token admin registry contract
	sourceTokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get all configured tokens: %w", err)
	}
	sourceTokenToPool := make(map[common.Address]common.Address)
	for _, sourceToken := range sourceTokens {
		pool, err := onRampContract.GetPoolBySourceToken(nil, 0, sourceToken)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get pool by source token: %w", err)
		}
		sourceTokenToPool[sourceToken] = pool
	}

	destChainSpecificData := make(map[uint64]DestChainSpecificData)
	for _, destChainSelector := range destChainSelectors {
		allowListInformation, err := onRampContract.GetAllowedSendersList(nil, destChainSelector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get allowed senders list: %w", err)
		}
		destChainConfig, err := onRampContract.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get dest chain config: %w", err)
		}
		expectedNextSeqNum, err := onRampContract.GetExpectedNextSequenceNumber(nil, destChainSelector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get expected next sequence number: %w", err)
		}
		destChainSpecificData[destChainSelector] = DestChainSpecificData{
			AllowedSendersList: allowListInformation.ConfiguredAddresses,
			DestChainConfig:    destChainConfig,
			ExpectedNextSeqNum: expectedNextSeqNum,
		}
	}

	return OnRampView{
		ContractMetaData:      tv,
		DynamicConfig:         dynamicConfig,
		StaticConfig:          staticConfig,
		SourceTokenToPool:     sourceTokenToPool,
		DestChainSpecificData: destChainSpecificData,
	}, nil
}
