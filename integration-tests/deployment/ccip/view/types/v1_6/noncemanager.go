package v1_6

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
)

type NonceManager struct {
	types.Contract
	AuthorizedCallers []common.Address `json:"authorizedCallers,omitempty"`
}

func (nm NonceManager) Address() common.Address {
	return common.HexToAddress(nm.Contract.Address)
}

func NonceManagerSnapshot(nm *nonce_manager.NonceManager) (NonceManager, error) {
	authorizedCallers, err := nm.GetAllAuthorizedCallers(nil)
	if err != nil {
		return NonceManager{}, err
	}
	tv, err := nm.TypeAndVersion(nil)
	if err != nil {
		return NonceManager{}, err
	}
	return NonceManager{
		Contract: types.Contract{
			TypeAndVersion: tv,
			Address:        nm.Address().Hex(),
		},
		// TODO: these can be resolved using an address book
		AuthorizedCallers: authorizedCallers,
	}, nil
}
