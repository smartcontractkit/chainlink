package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

type TokenAdminRegistryView struct {
	types.ContractMetaData
	Tokens []common.Address `json:"tokens"`
}

func GenerateTokenAdminRegistryView(taContract *token_admin_registry.TokenAdminRegistry) (TokenAdminRegistryView, error) {
	tokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("view error for token admin registry: %w", err)
	}
	tvMeta, err := types.NewContractMetaData(taContract, taContract.Address())
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("metadata error for token admin registry: %w", err)
	}
	return TokenAdminRegistryView{
		ContractMetaData: tvMeta,
		Tokens:           tokens,
	}, nil
}
