package v1_0

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type StaticLinkTokenView struct {
	types.ContractMetaData
	Decimals uint8    `json:"decimals"`
	Supply   *big.Int `json:"supply"`
}

func GenerateStaticLinkTokenView(lt *link_token_interface.LinkToken) (StaticLinkTokenView, error) {
	decimals, err := lt.Decimals(nil)
	if err != nil {
		return StaticLinkTokenView{}, fmt.Errorf("failed to get decimals %s: %w", lt.Address(), err)
	}
	totalSupply, err := lt.TotalSupply(nil)
	if err != nil {
		return StaticLinkTokenView{}, fmt.Errorf("failed to get total supply %s: %w", lt.Address(), err)
	}
	return StaticLinkTokenView{
		ContractMetaData: types.ContractMetaData{
			TypeAndVersion: cldf.TypeAndVersion{
				Type:    commontypes.StaticLinkToken,
				Version: deployment.Version1_0_0,
			}.String(),
			Address: lt.Address(),
			// No owner.
		},
		Decimals: decimals,
		Supply:   totalSupply,
	}, nil
}
