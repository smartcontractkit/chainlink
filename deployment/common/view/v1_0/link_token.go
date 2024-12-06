package v1_0

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type LinkTokenView struct {
	types.ContractMetaData
	Decimals uint8            `json:"decimals"`
	Supply   *big.Int         `json:"supply"`
	Minters  []common.Address `json:"minters"`
	Burners  []common.Address `json:"burners"`
}

func GenerateLinkTokenView(lt *link_token.LinkToken) (LinkTokenView, error) {
	owner, err := lt.Owner(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("failed to get owner: %w", err)
	}
	decimals, err := lt.Decimals(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("failed to get decimals: %w", err)
	}
	totalSupply, err := lt.TotalSupply(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("failed to get total supply: %w", err)
	}
	minters, err := lt.GetMinters(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("failed to get minters: %w", err)
	}
	burners, err := lt.GetBurners(nil)
	if err != nil {
		return LinkTokenView{}, fmt.Errorf("failed to get burners: %w", err)
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
		Minters:  minters,
		Burners:  burners,
	}, nil
}
