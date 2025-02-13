package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/nonce_manager"
)

type NonceManagerView struct {
	types.ContractMetaData
	AuthorizedCallers []common.Address `json:"authorizedCallers,omitempty"`
}

func GenerateNonceManagerView(nm *nonce_manager.NonceManager) (NonceManagerView, error) {
	authorizedCallers, err := nm.GetAllAuthorizedCallers(nil)
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("view error for nonce manager: %w", err)
	}
	nmMeta, err := types.NewContractMetaData(nm, nm.Address())
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("metadata error for nonce manager: %w", err)
	}
	return NonceManagerView{
		ContractMetaData: nmMeta,
		// TODO: these can be resolved using an address book
		AuthorizedCallers: authorizedCallers,
	}, nil
}
