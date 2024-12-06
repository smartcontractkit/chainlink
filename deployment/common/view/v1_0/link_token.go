package v1_0

import (
	"fmt"
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
		return LinkTokenView{}, fmt.Errorf("view error to get link token owner addr %s: %w", lt.Address().String(), err)
	}
	decimals, err := lt.Decimals(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("view error to get link token decimals addr %s: %w", lt.Address().String(), err)
	}
	totalSupply, err := lt.TotalSupply(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("view error to get link token total supply addr %s: %w", lt.Address().String(), err)
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
