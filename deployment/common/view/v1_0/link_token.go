package v1_0

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type LinkTokenView struct {
	types.ContractMetaData
	Decimals uint8    `json:"decimals"`
	Supply   *big.Int `json:"supply"`
}

func GenerateLinkTokenView(lt *link_token.LinkToken) (LinkTokenView, error) {
	owner, err := lt.Owner(nil)
	if err != nil {
		return LinkTokenView{}, err
	}
	decimals, err := lt.Decimals(nil)
	if err != nil {
		return LinkTokenView{}, err
	}
	totalSupply, err := lt.TotalSupply(nil)
	if err != nil {
		return LinkTokenView{}, err
	}
	return LinkTokenView{
		ContractMetaData: types.ContractMetaData{
			TypeAndVersion: deployment.TypeAndVersion{
				commontypes.LinkToken,
				deployment.Version1_0_0,
			}.String(),
			Address: lt.Address(),
			Owner:   owner,
		},
		Decimals: decimals,
		Supply:   totalSupply,
	}, nil
}
