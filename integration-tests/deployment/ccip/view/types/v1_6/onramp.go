package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

type OnRamp struct {
	types.ContractMetaData
	DynamicConfig         onramp.OnRampDynamicConfig        `json:"dynamicConfig"`
	StaticConfig          onramp.OnRampStaticConfig         `json:"staticConfig"`
	Owner                 common.Address                    `json:"owner"`
	SourceTokenToPool     map[common.Address]common.Address `json:"sourceTokenToPool"`
	DestChainSpecificData map[uint64]DestChainSpecificData  `json:"destChainSpecificData"`
}

type DestChainSpecificData struct {
	AllowedSendersList []common.Address          `json:"allowedSendersList"`
	DestChainConfig    onramp.GetDestChainConfig `json:"destChainConfig"`
	ExpectedNextSeqNum uint64                    `json:"expectedNextSeqNum"`
	Router             common.Address            `json:"router"`
}

func OnRampSnapshot(
	onRampContract *onramp.OnRamp,
	routerContract types.ContractMetaData,
	taContract types.ContractMetaData,
	client bind.ContractBackend,
) (OnRamp, error) {
	tv, err := onRampContract.TypeAndVersion(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	dynamicConfig, err := onRampContract.GetDynamicConfig(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	staticConfig, err := onRampContract.GetStaticConfig(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get static config: %w", err)
	}

	owner, err := onRampContract.Owner(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get owner: %w", err)
	}
	// populate destChainSelectors from router
	destChainSelectors := make([]uint64, 0)
	switch routerContract.TypeAndVersion {
	case types.RouterTypeAndVersionV1_2:
		router, err := router1_2.NewRouter(routerContract.Address, client)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get router: %w", err)
		}
		offRampList, err := router.GetOffRamps(nil)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get offRamps from router %w", err)
		}
		for _, offRamp := range offRampList {
			// the lanes are bidirectional, so we get the list of source chains to know which chains are supported as destinations as well
			destChainSelectors = append(destChainSelectors, offRamp.SourceChainSelector)
		}
	default:
		return OnRamp{}, fmt.Errorf("unsupported router type and version: %s", routerContract.TypeAndVersion)
	}

	// populate sourceTokens from token admin registry contract
	var sourceTokens []common.Address
	switch taContract.TypeAndVersion {
	case types.TokenAdminRegistryTypeAndVersionV1_5:
		ta, err := token_admin_registry.NewTokenAdminRegistry(taContract.Address, client)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get token admin registry: %w", err)
		}
		// TODO : CCIP-3416 : get all tokens here instead of just 10
		sourceTokens, err = ta.GetAllConfiguredTokens(nil, 0, 10)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get all configured tokens: %w", err)
		}
	default:
		return OnRamp{}, fmt.Errorf("unsupported token admin registry type and version: %s", taContract.TypeAndVersion)
	}

	sourceTokenToPool := make(map[common.Address]common.Address)
	for _, sourceToken := range sourceTokens {
		pool, err := onRampContract.GetPoolBySourceToken(nil, 0, sourceToken)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get pool by source token: %w", err)
		}
		sourceTokenToPool[sourceToken] = pool
	}

	destChainSpecificData := make(map[uint64]DestChainSpecificData)
	for _, destChainSelector := range destChainSelectors {
		allowedSendersList, err := onRampContract.GetAllowedSendersList(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get allowed senders list: %w", err)
		}
		destChainConfig, err := onRampContract.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get dest chain config: %w", err)
		}
		expectedNextSeqNum, err := onRampContract.GetExpectedNextSequenceNumber(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get expected next sequence number: %w", err)
		}
		router, err := onRampContract.GetRouter(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get router: %w", err)
		}
		destChainSpecificData[destChainSelector] = DestChainSpecificData{
			AllowedSendersList: allowedSendersList,
			DestChainConfig:    destChainConfig,
			ExpectedNextSeqNum: expectedNextSeqNum,
			Router:             router,
		}
	}

	return OnRamp{
		ContractMetaData: types.ContractMetaData{
			TypeAndVersion: tv,
			Address:        onRampContract.Address(),
		},
		DynamicConfig:         dynamicConfig,
		StaticConfig:          staticConfig,
		Owner:                 owner,
		SourceTokenToPool:     sourceTokenToPool,
		DestChainSpecificData: destChainSpecificData,
	}, nil
}
