package view

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

type OffRamp struct {
	Contract
	DynamicConfig             offramp.OffRampDynamicConfig                `json:"dynamicConfig"`
	LatestPriceSequenceNumber uint64                                      `json:"latestPriceSequenceNumber"`
	SourceChainConfigs        map[uint64]offramp.OffRampSourceChainConfig `json:"sourceChainConfigs"`
	StaticConfig              offramp.OffRampStaticConfig                 `json:"staticConfig"`
	Owner                     common.Address                              `json:"owner"`
}

type OffRampReader interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Address() common.Address
	GetDynamicConfig(opts *bind.CallOpts) (offramp.OffRampDynamicConfig, error)
	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)
	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (offramp.OffRampSourceChainConfig, error)
	GetStaticConfig(opts *bind.CallOpts) (offramp.OffRampStaticConfig, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
}

func OffRampSnapshot(
	offRampReader OffRampReader,
	sourceChainSelectors []uint64,
) (OffRamp, error) {
	tv, err := offRampReader.TypeAndVersion(nil)
	if err != nil {
		return OffRamp{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	dynamicConfig, err := offRampReader.GetDynamicConfig(nil)
	if err != nil {
		return OffRamp{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	latestPriceSequenceNumber, err := offRampReader.GetLatestPriceSequenceNumber(nil)
	if err != nil {
		return OffRamp{}, fmt.Errorf("failed to get latest price sequence number: %w", err)
	}

	sourceChainConfigs := make(map[uint64]offramp.OffRampSourceChainConfig)
	for _, sourceChainSelector := range sourceChainSelectors {
		sourceChainConfig, err := offRampReader.GetSourceChainConfig(nil, sourceChainSelector)
		if err != nil {
			return OffRamp{}, fmt.Errorf("failed to get source chain config: %w", err)
		}
		sourceChainConfigs[sourceChainSelector] = sourceChainConfig
	}

	staticConfig, err := offRampReader.GetStaticConfig(nil)
	if err != nil {
		return OffRamp{}, fmt.Errorf("failed to get static config: %w", err)
	}

	owner, err := offRampReader.Owner(nil)
	if err != nil {
		return OffRamp{}, fmt.Errorf("failed to get owner: %w", err)
	}

	return OffRamp{
		Contract: Contract{
			TypeAndVersion: tv,
			Address:        offRampReader.Address().Hex(),
		},
		DynamicConfig:             dynamicConfig,
		LatestPriceSequenceNumber: latestPriceSequenceNumber,
		SourceChainConfigs:        sourceChainConfigs,
		StaticConfig:              staticConfig,
		Owner:                     owner,
	}, nil
}
