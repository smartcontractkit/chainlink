package v1_5

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type TokenAdminRegistryGetAllConfiguredTokensIn struct {
	Address       common.Address
	ChainSelector uint64
}

var (
	TokenAdminRegistryGetAllConfiguredTokensOp = operations.NewOperation(
		"TokenAdminRegistryGetAllConfiguredTokensOp",
		semver.MustParse("1.0.0"),
		"Gets all configured tokens from the TokenAdminRegistry",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input TokenAdminRegistryGetAllConfiguredTokensIn) ([]common.Address, error) {
			tokenAdminReg, err := token_admin_registry.NewTokenAdminRegistry(input.Address, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create TokenAdminRegistry contract binding: chainSelector=%v, TokenAdminRegistry Address=%s, error=%w", deps.Chain.ChainSelector(), input.Address.Hex(), err)
			}

			allTransferTokens := []common.Address{}
			var offset uint64 = 0
			const pageSize uint64 = 1000
			for {
				pageTokens, err := tokenAdminReg.GetAllConfiguredTokens(nil, offset, pageSize)
				if err != nil {
					return nil, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry at offset %d: %w", offset, err)
				}

				if len(pageTokens) == 0 { // No more tokens to fetch
					break
				}

				allTransferTokens = append(allTransferTokens, pageTokens...)
				offset += pageSize
			}
			return allTransferTokens, nil
		})
)
