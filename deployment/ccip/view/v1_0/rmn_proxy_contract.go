package v1_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type RMNProxyView struct {
	types.ContractMetaData
	RMN common.Address `json:"rmn"`
}

func GenerateRMNProxyView(r *rmn_proxy_contract.RMNProxy) (RMNProxyView, error) {
	if r == nil {
		return RMNProxyView{}, errors.New("cannot generate view for nil RMNProxy")
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
