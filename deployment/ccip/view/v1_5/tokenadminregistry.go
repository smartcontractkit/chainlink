package v1_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
)

type TokenAdminRegistryView struct {
	types.ContractMetaData
	Tokens map[common.Address]TokenDetails `json:"tokens"`
}

type TokenDetails struct {
	Pool         common.Address `json:"pool"`
	Admin        common.Address `json:"admin"`
	PendingAdmin common.Address `json:"pendingAdmin"`
}

func GenerateTokenAdminRegistryView(taContract *token_admin_registry.TokenAdminRegistry) (TokenAdminRegistryView, error) {
	if taContract == nil {
		return TokenAdminRegistryView{}, errors.New("token admin registry contract is nil")
	}
	tokens, err := getAllConfiguredTokensPaginated(taContract)
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

// getAllConfiguredTokensPaginated fetches all configured tokens from the TokenAdminRegistry contract in paginated
// manner to avoid RPC timeouts since the list of configured tokens can grow to be very large over time.
func getAllConfiguredTokensPaginated(taContract *token_admin_registry.TokenAdminRegistry) (map[common.Address]TokenDetails, error) {
	tokenDetails := make(map[common.Address]TokenDetails)
	allTokens, err := shared.GetSupportedTokens(taContract)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported tokens for tokenAdminRegistry %s: %w", taContract.Address().String(), err)
	}
	for _, token := range allTokens {
		config, err := taContract.GetTokenConfig(nil, token)
		if err != nil {
			return nil, fmt.Errorf("failed to get token config for token %s tokenAdminReg %s: %w",
				token.String(), taContract.Address().String(), err)
		}
		tokenDetails[token] = TokenDetails{
			Pool:         config.TokenPool,
			Admin:        config.Administrator,
			PendingAdmin: config.PendingAdministrator,
		}
	}
	return tokenDetails, nil
}
