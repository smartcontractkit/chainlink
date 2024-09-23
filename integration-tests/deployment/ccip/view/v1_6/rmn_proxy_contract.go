package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
)

type RMNProxyView struct {
	types.ContractMetaData
	ARM common.Address `json:"arm"`
}

func GenerateRMNProxyView(r *rmn_proxy_contract.RMNProxyContract) (RMNProxyView, error) {
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return RMNProxyView{}, fmt.Errorf("failed to generate contract metadata for RMNProxy: %w", err)
	}
	arm, err := r.GetARM(nil)
	if err != nil {
		return RMNProxyView{}, fmt.Errorf("failed to get ARM: %w", err)
	}
	return RMNProxyView{
		ContractMetaData: meta,
		ARM:              arm,
	}, nil
}
