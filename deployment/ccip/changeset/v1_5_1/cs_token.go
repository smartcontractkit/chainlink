package v1_5_1

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSet[TokenChangesetConfig] = RenounceRoleBurnMintERC20

// TokenRole defines the roles that can be assigned to accounts on BurnMintERC20 tokens.
type TokenRole int8

const (
	RoleBurner TokenRole = iota
	RoleAdmin
	RoleMinter
)

// TokenInfo holds information about a token, including its address, role, and the account to which the role is assigned.
type TokenInfo struct {
	Address common.Address
	Role    TokenRole
	Account common.Address
}

// TokenChangesetConfig defines a configuration for changes to BurnMintERC20 tokens.
type TokenChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]TokenInfo
	MCMS   *proposalutils.TimelockConfig
}

// Validate checks if the TokenChangesetConfig is valid for the given environment.
func (c TokenChangesetConfig) Validate(env cldf.Environment) error {
	ctx := env.GetContext()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokenSymbol := range c.Tokens {
		if err := stateview.ValidateChain(env, state, chainSelector, nil); err != nil {
			return fmt.Errorf("failed to validate chain with selector %d: %w", chainSelector, err)
		}

		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain)
		}

		for symbol, tokenInfo := range tokenSymbol {
			if _, ok := chainState.BurnMintERC20[symbol]; !ok {
				return fmt.Errorf("token %s does not exist on chain %s", symbol, chain)
			}

			if tokenInfo.Address == utils.ZeroAddress {
				return fmt.Errorf("token address for %s on chain %s is missing", symbol, chain)
			}

			if tokenInfo.Account == utils.ZeroAddress {
				return fmt.Errorf("account address for %s on chain %s is missing", symbol, chain)
			}

			token, err := burn_mint_erc20.NewBurnMintERC20(tokenInfo.Address, env.BlockChains.EVMChains()[chainSelector].Client)
			if err != nil {
				return fmt.Errorf("failed to create BurnMintERC20 instance for %s on chain %s: %w", symbol, chain, err)
			}

			role, err := getRoleOnTokenByTokenRole(ctx, token, tokenInfo.Role)
			if err != nil {
				return fmt.Errorf("failed to get role %d for token %s on chain %s: %w", tokenInfo.Role, symbol, chain, err)
			}

			hasRole, err := token.HasRole(&bind.CallOpts{Context: ctx}, role, tokenInfo.Account)
			if err != nil {
				return fmt.Errorf("failed to check if account %s has role %d for token %s on chain %s: %w", tokenInfo.Account, tokenInfo.Role, symbol, chain, err)
			}
			if hasRole {
				return fmt.Errorf("account %s already has role %d for token %s on chain %s", tokenInfo.Account, tokenInfo.Role, symbol, chain)
			}
		}
	}

	return nil
}

// getRoleOnTokenByTokenRole retrieves the role address for a given token and role type.
func getRoleOnTokenByTokenRole(ctx context.Context, token *burn_mint_erc20.BurnMintERC20, role TokenRole) ([32]byte, error) {
	switch role {
	case RoleBurner:
		r, err := token.BURNERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, err
		}
		return r, nil
	case RoleAdmin:
		r, err := token.DEFAULTADMINROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, err
		}
		return r, nil
	case RoleMinter:
		r, err := token.MINTERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, err
		}
		return r, nil
	default:
		return [32]byte{}, fmt.Errorf("unknown token role: %d", role)
	}
}

// RenounceRoleBurnMintERC20 renounces roles on BurnMintERC20 tokens.
func RenounceRoleBurnMintERC20(env cldf.Environment, c TokenChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("renounce roles on BurnMintERC20 tokens")

	for chainSelector := range c.Tokens {
		chain := env.BlockChains.EVMChains()[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s: %w", chain, err)
		}

		for symbol, tokenInfo := range c.Tokens[chainSelector] {
			token, err := burn_mint_erc20.NewBurnMintERC20(tokenInfo.Address, chain.Client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create BurnMintERC20 instance for %s on chain %s: %w", symbol, chain, err)
			}

			role, err := getRoleOnTokenByTokenRole(env.GetContext(), token, tokenInfo.Role)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get role %d for token %s: %w", tokenInfo.Role, symbol, err)
			}

			if _, err := token.RenounceRole(opts, role, tokenInfo.Account); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to renounce role %d for token %s on chain %s: %w", tokenInfo.Role, symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
