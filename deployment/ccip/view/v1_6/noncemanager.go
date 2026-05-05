package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	commoncldchangesets "github.com/smartcontractkit/cld-changesets/pkg/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
)

type NonceManagerView struct {
	commoncldchangesets.ContractMetaData
	AuthorizedCallers []common.Address `json:"authorizedCallers,omitempty"`
}

func GenerateNonceManagerView(nm *nonce_manager.NonceManager) (NonceManagerView, error) {
	authorizedCallers, err := nm.GetAllAuthorizedCallers(nil)
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("view error for nonce manager: %w", err)
	}
	nmMeta, err := commoncldchangesets.NewContractMetaData(nm, nm.Address())
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("metadata error for nonce manager: %w", err)
	}
	return NonceManagerView{
		ContractMetaData: nmMeta,
		// TODO: these can be resolved using an address book
		AuthorizedCallers: authorizedCallers,
	}, nil
}
