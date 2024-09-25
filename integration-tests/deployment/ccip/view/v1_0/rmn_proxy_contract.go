package v1_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
)

type RMNProxyView struct {
	types.ContractMetaData
	RMN common.Address `json:"rmn"`
}

func GenerateRMNProxyView(r *rmn_proxy_contract.RMNProxyContract) (RMNProxyView, error) {
	if r == nil {
		return RMNProxyView{}, fmt.Errorf("cannot generate view for nil RMNProxy")
	}
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return RMNProxyView{}, fmt.Errorf("failed to generate contract metadata for RMNProxy: %w", err)
	}
	rmn, err := r.GetARM(nil)
	if err != nil {
		return RMNProxyView{}, fmt.Errorf("failed to get ARM: %w", err)
	}
	return RMNProxyView{
		ContractMetaData: meta,
		RMN:              rmn,
	}, nil
}
