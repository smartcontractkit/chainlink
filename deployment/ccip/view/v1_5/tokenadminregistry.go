package v1_5

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
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
	tokenDetailsSyncMap := sync.Map{}
	grp := errgroup.Group{}
	// try to get all token details in parallel
	for _, token := range allTokens {
		token := token
		grp.Go(func() error {
			config, err := taContract.GetTokenConfig(nil, token)
			if err != nil {
				return fmt.Errorf("failed to get token config for token %s tokenAdminReg %s: %w",
					token.String(), taContract.Address().String(), err)
			}
			tokenDetailsSyncMap.Store(token, TokenDetails{
				Pool:         config.TokenPool,
				Admin:        config.Administrator,
				PendingAdmin: config.PendingAdministrator,
			})
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get token details for tokenAdminRegistry %s: %w", taContract.Address().String(), err)
	}
	// convert sync map to regular map
	tokenDetailsSyncMap.Range(func(key, value interface{}) bool {
		tokenDetails[key.(common.Address)] = value.(TokenDetails)
		return true
	})
	return tokenDetails, nil
}
