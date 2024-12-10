package v1_5

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
)

type OnRampView struct {
	types.ContractMetaData
	StaticConfig  evm_2_evm_onramp.EVM2EVMOnRampStaticConfig  `json:"staticConfig"`
	DynamicConfig evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig `json:"dynamicConfig"`
}

func GenerateOnRampView(r *evm_2_evm_onramp.EVM2EVMOnRamp) (OnRampView, error) {
	if r == nil {
		return OnRampView{}, fmt.Errorf("cannot generate view for nil OnRamp")
	}
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to generate contract metadata for OnRamp: %w", err)
	}
	staticConfig, err := r.GetStaticConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get static config for OnRamp: %w", err)
	}
	dynamicConfig, err := r.GetDynamicConfig(nil)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get dynamic config for OnRamp: %w", err)
	}

	// Add billing if needed, maybe not required for legacy contract?
	return OnRampView{
		ContractMetaData: meta,
		StaticConfig:     staticConfig,
		DynamicConfig:    dynamicConfig,
	}, nil
}
