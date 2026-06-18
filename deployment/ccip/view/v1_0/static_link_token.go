package v1_0

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/cld-changesets/pkg/cldfutil"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	linkcontracts "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/contracts/link"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
)

type StaticLinkTokenView struct {
	cldfutil.ContractMetaData
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
		ContractMetaData: cldfutil.ContractMetaData{
			TypeAndVersion: cldf.TypeAndVersion{
				Type:    linkcontracts.StaticLinkToken,
				Version: *semver.MustParse("1.0.0"),
			}.String(),
			Address: lt.Address(),
			// No owner.
		},
		Decimals: decimals,
		Supply:   totalSupply,
	}, nil
}
