package v1_5

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type OffRampView struct {
	types.ContractMetaData
	StaticConfig  evm_2_evm_offramp.EVM2EVMOffRampStaticConfig  `json:"staticConfig"`
	DynamicConfig evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig `json:"dynamicConfig"`
}

func GenerateOffRampView(r *evm_2_evm_offramp.EVM2EVMOffRamp) (OffRampView, error) {
	if r == nil {
		return OffRampView{}, errors.New("cannot generate view for nil OffRamp")
	}
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to generate contract metadata for OffRamp %s: %w", r.Address(), err)
	}
	staticConfig, err := r.GetStaticConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get static config for OffRamp %s: %w", r.Address(), err)
	}
	dynamicConfig, err := r.GetDynamicConfig(nil)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get dynamic config for OffRamp %s: %w", r.Address(), err)
	}

	// TODO: If needed, we can filter logs to get the OCR config.
	// May not be required for the legacy contracts.
	return OffRampView{
		ContractMetaData: meta,
		StaticConfig:     staticConfig,
		DynamicConfig:    dynamicConfig,
	}, nil
}
