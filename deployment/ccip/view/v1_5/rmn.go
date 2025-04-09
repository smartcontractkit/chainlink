package v1_5

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type RMNView struct {
	types.ContractMetaData
	ConfigDetails rmn_contract.GetConfigDetails `json:"configDetails"`
}

func GenerateRMNView(r *rmn_contract.RMNContract) (RMNView, error) {
	if r == nil {
		return RMNView{}, errors.New("cannot generate view for nil RMN")
	}
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return RMNView{}, fmt.Errorf("failed to generate contract metadata for RMN: %w", err)
	}
	config, err := r.GetConfigDetails(nil)
	if err != nil {
		return RMNView{}, fmt.Errorf("failed to get config details for RMN: %w", err)
	}
	return RMNView{
		ContractMetaData: meta,
		ConfigDetails:    config,
	}, nil
}
