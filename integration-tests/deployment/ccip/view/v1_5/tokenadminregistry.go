package v1_5

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

type TokenAdminRegistry struct {
	view.Contract
	Tokens []common.Address `json:"tokens"`
}

func (ta TokenAdminRegistry) Address() common.Address {
	return common.HexToAddress(ta.Contract.Address)
}

func TokenAdminRegistrySnapshot(taContract *token_admin_registry.TokenAdminRegistry) (TokenAdminRegistry, error) {
	tokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
	if err != nil {
		return TokenAdminRegistry{}, err
	}
	tv, err := taContract.TypeAndVersion(nil)
	if err != nil {
		return TokenAdminRegistry{}, err
	}
	return TokenAdminRegistry{
		Contract: view.Contract{
			TypeAndVersion: tv,
			Address:        taContract.Address().Hex(),
		},
		Tokens: tokens,
	}, nil
}
