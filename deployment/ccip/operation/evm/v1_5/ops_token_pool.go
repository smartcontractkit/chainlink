package v1_5

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

var (
	TokenPoolGetSupportedChainsOp = operations.NewOperation(
		"TokenPoolGetSupportedChainsOp",
		semver.MustParse("1.0.0"),
		"Gets all Supported Chains of a Token Pool",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, tokenPoolAddress common.Address) ([]uint64, error) {
			tokenPool, err := token_pool.NewTokenPool(tokenPoolAddress, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create tokenpool contract binding: chainSelector=%v, Token Pool=%s, error=%w", deps.Chain.ChainSelector(), tokenPoolAddress.Hex(), err)
			}
			supportedChains, err := tokenPool.GetSupportedChains(nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get supported chains from token pool: chainSelector=%v, Token Pool=%s, error=%w", deps.Chain.ChainSelector(), tokenPoolAddress.Hex(), err)
			}
			return supportedChains, nil
		})
)
