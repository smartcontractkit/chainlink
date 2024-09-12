package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

type OnRamp struct {
	view.Contract
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
	onRampReader *onramp.OnRamp,
	destChainSelectors []uint64,
	sourceTokens []common.Address,
) (OnRamp, error) {
	tv, err := onRampReader.TypeAndVersion(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	dynamicConfig, err := onRampReader.GetDynamicConfig(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	staticConfig, err := onRampReader.GetStaticConfig(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get static config: %w", err)
	}

	owner, err := onRampReader.Owner(nil)
	if err != nil {
		return OnRamp{}, fmt.Errorf("failed to get owner: %w", err)
	}

	sourceTokenToPool := make(map[common.Address]common.Address)
	for _, sourceToken := range sourceTokens {
		pool, err := onRampReader.GetPoolBySourceToken(nil, 0, sourceToken)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get pool by source token: %w", err)
		}
		sourceTokenToPool[sourceToken] = pool
	}

	destChainSpecificData := make(map[uint64]DestChainSpecificData)
	for _, destChainSelector := range destChainSelectors {
		allowedSendersList, err := onRampReader.GetAllowedSendersList(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get allowed senders list: %w", err)
		}
		destChainConfig, err := onRampReader.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get dest chain config: %w", err)
		}
		expectedNextSeqNum, err := onRampReader.GetExpectedNextSequenceNumber(nil, destChainSelector)
		if err != nil {
			return OnRamp{}, fmt.Errorf("failed to get expected next sequence number: %w", err)
		}
		router, err := onRampReader.GetRouter(nil, destChainSelector)
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
		Contract: view.Contract{
			TypeAndVersion: tv,
			Address:        onRampReader.Address().Hex(),
		},
		DynamicConfig:         dynamicConfig,
		StaticConfig:          staticConfig,
		Owner:                 owner,
		SourceTokenToPool:     sourceTokenToPool,
		DestChainSpecificData: destChainSpecificData,
	}, nil
}
