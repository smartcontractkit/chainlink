package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

type TokenAdminRegistry struct {
	types.ContractMetaData
	Tokens []common.Address `json:"tokens"`
}

func (ta *TokenAdminRegistry) Snapshot(contractMeta types.ContractMetaData, _ []types.ContractMetaData, client bind.ContractBackend) error {
	ta.ContractMetaData = contractMeta
	if err := ta.ContractMetaData.Validate(); err != nil {
		return fmt.Errorf("snapshot error for TokenAdminRegistry: %w", err)
	}
	taContract, err := token_admin_registry.NewTokenAdminRegistry(ta.Address, client)
	if err != nil {
		return fmt.Errorf("failed to get token admin registry contract: %w", err)
	}
	// TODO : CCIP-3416 : get all tokens here instead of just 10
	tokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
	if err != nil {
		return err
	}
	ta.Tokens = tokens
	return nil
}
