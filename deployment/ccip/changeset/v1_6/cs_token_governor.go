package v1_6

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/latest/token_governor"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type TokenGovernorRole uint8

const (
	RoleMinter TokenGovernorRole = iota
	RoleBridgerMinterOrBurner
	RoleBurner
	RoleFreezer
	RoleUnfreezer
	RolePauser
	RoleUnpauser
	RoleRecovery
	RoleCheckerAdmin
	RoleDefaultAdmin
)

var (
	_ cldf.ChangeSet[TokenGovernorChangesetConfig]     = DeployTokenGovernor
	_ cldf.ChangeSet[TokenGovernorRoleChangesetConfig] = GrantRoleTokenGovernor
	_ cldf.ChangeSet[TokenGovernorRoleChangesetConfig] = RenounceRoleTokenGovernor
	_ cldf.ChangeSet[TokenGovernorRoleChangesetConfig] = TransferOwnershipTokenGovernor
	_ cldf.ChangeSet[TokenGovernorChangesetConfig]     = AcceptOwnershipTokenGovernor
	_ cldf.ChangeSet[TokenGovernorRoleChangesetConfig] = BeingDefaultAdminTransferTokenGovernor
	_ cldf.ChangeSet[TokenGovernorChangesetConfig]     = AcceptDefaultAdminTransferTokenGovernor
)

type TokenGovernor struct {
	Token common.Address
	// If nil, the default value will be used which is 0.
	InitialDelay        *big.Int
	InitialDefaultAdmin common.Address
}

type TokenGovernorChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]TokenGovernor
	MCMS   *proposalutils.TimelockConfig
}

type TokenGovernorGrantRole struct {
	Role    TokenGovernorRole
	Account common.Address
}

type TokenGovernorRoleChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]TokenGovernorGrantRole
	MCMS   *proposalutils.TimelockConfig
}

// String returns the string representation of the TokenGovernorRole.
func (r TokenGovernorRole) String() string {
	switch r {
	case RoleMinter:
		return "MINTER_ROLE"
	case RoleBridgerMinterOrBurner:
		return "BRIDGE_MINTER_OR_BURNER_ROLE"
	case RoleBurner:
		return "BURNER_ROLE"
	case RoleFreezer:
		return "FREEZER_ROLE"
	case RoleUnfreezer:
		return "UNFREEZER_ROLE"
	case RolePauser:
		return "PAUSER_ROLE"
	case RoleUnpauser:
		return "UNPAUSER_ROLE"
	case RoleRecovery:
		return "RECOVERY_ROLE"
	case RoleCheckerAdmin:
		return "CHECKER_ADMIN_ROLE"
	case RoleDefaultAdmin:
		return "DEFAULT_ADMIN_ROLE"
	default:
		return "UNKNOWN"
	}
}

// Validate validates the TokenGovernorChangesetConfig.
func (c TokenGovernorChangesetConfig) Validate(env cldf.Environment) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, governor := range tokens {
			if token == "" {
				return errors.New("token must be defined")
			}

			if err := stateview.ValidateChain(env, state, chainSelector, nil); err != nil {
				return fmt.Errorf("failed to validate chain with selector %d: %w", chainSelector, err)
			}

			chain, ok := env.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			if err := validateTokenSymbol(env.GetContext(), chain, governor.Token, token); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateTokenSymbol validates that the token at the given address has the given symbol.
func validateTokenSymbol(ctx context.Context, chain cldf_evm.Chain, address common.Address, targetSymbol shared.TokenSymbol) error {
	token, err := erc20.NewERC20(address, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with erc20 bindings: %w", address, err)
	}

	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch symbol from token with address %s: %w", address, err)
	}

	if symbol != string(targetSymbol) {
		return fmt.Errorf("symbol of token with address %s (%s) does not match expected symbol (%s)", address, symbol, targetSymbol)
	}

	return nil
}

// Validate validates the TokenGovernorRoleChangesetConfig.
func (c TokenGovernorRoleChangesetConfig) Validate(env cldf.Environment) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token := range tokens {
			if token == "" {
				return errors.New("token must be defined")
			}

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

			tokenGovernor, ok := chainState.TokenGovernor[token]
			if !ok {
				return fmt.Errorf("token governor does not exist for %s", token)
			}

			tokenAddress, err := tokenGovernor.GetToken(&bind.CallOpts{Context: env.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to fetch token from token governor: %w", err)
			}

			if err := validateTokenSymbol(env.GetContext(), chain, tokenAddress, token); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetRoleFromTokenGovernor returns the role bytes32 from the token governor contract.
func GetRoleFromTokenGovernor(ctx context.Context, tokenGovernor *token_governor.TokenGovernor, role TokenGovernorRole) ([32]byte, error) {
	if tokenGovernor == nil {
		return [32]byte{}, errors.New("token governor is nil")
	}

	switch role {
	case RoleMinter:
		r, err := tokenGovernor.MINTERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch minter role: %w", err)
		}
		return r, nil
	case RoleBridgerMinterOrBurner:
		r, err := tokenGovernor.BRIDGEMINTERORBURNERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch bridge minter or burner role: %w", err)
		}
		return r, nil
	case RoleBurner:
		r, err := tokenGovernor.BURNERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch burner role: %w", err)
		}
		return r, nil
	case RoleFreezer:
		r, err := tokenGovernor.FREEZERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch freezer role: %w", err)
		}
		return r, nil
	case RoleUnfreezer:
		r, err := tokenGovernor.UNFREEZERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch unfreezer role: %w", err)
		}
		return r, nil
	case RolePauser:
		r, err := tokenGovernor.PAUSERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch pauser role: %w", err)
		}
		return r, nil
	case RoleUnpauser:
		r, err := tokenGovernor.UNPAUSERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch unpauser role: %w", err)
		}
		return r, nil
	case RoleRecovery:
		r, err := tokenGovernor.RECOVERYROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch recovery role: %w", err)
		}
		return r, nil
	case RoleCheckerAdmin:
		r, err := tokenGovernor.CHECKERADMINROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch checker admin role: %w", err)
		}
		return r, nil
	case RoleDefaultAdmin:
		r, err := tokenGovernor.DEFAULTADMINROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to fetch default admin role: %w", err)
		}
		return r, nil
	}

	return [32]byte{}, nil
}

// DeployTokenGovernor deploys the token governor contracts on the given chains.
func DeployTokenGovernor(env cldf.Environment, c TokenGovernorChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	newAddresses := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := env.BlockChains.EVMChains()[chainSelector]

		for token, governor := range tokens {
			if governor.InitialDelay == nil {
				governor.InitialDelay = big.NewInt(0)
			}

			if governor.InitialDefaultAdmin == utils.ZeroAddress {
				return cldf.ChangesetOutput{}, errors.New("initial default admin must be defined")
			}

			chainState, _ := state.EVMChainState(chainSelector)
			if _, ok := chainState.TokenGovernor[token]; ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("token governor already exists for %s", governor.Token)
			}

			_, err := cldf.DeployContract(env.Logger, chain, newAddresses,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_governor.TokenGovernor] {
					tgAddress, tx, tokenGovernor, err := token_governor.DeployTokenGovernor(chain.DeployerKey, chain.Client, governor.Token, governor.InitialDelay, governor.InitialDefaultAdmin)
					return cldf.ContractDeploy[*token_governor.TokenGovernor]{
						Address:  tgAddress,
						Contract: tokenGovernor,
						Tv:       cldf.NewTypeAndVersion(shared.TokenGovernor, deployment.Version1_6_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)

			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token governor on %s: %w", chain, err)
			}
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

// GrantRoleTokenGovernor grants the given role to the given account on the given chains.
func GrantRoleTokenGovernor(env cldf.Environment, c TokenGovernorRoleChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorRoleChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor role grant")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token, governor := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]

			role, err := GetRoleFromTokenGovernor(env.GetContext(), tokenGovernor, governor.Role)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			hasRole, err := tokenGovernor.HasRole(&bind.CallOpts{Context: env.GetContext()}, role, governor.Account)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if account has role: %w", err)
			}

			if hasRole {
				return cldf.ChangesetOutput{}, fmt.Errorf("account %s already has role %s", governor.Account, governor.Role)
			}

			if _, err := tokenGovernor.GrantRole(opts, role, governor.Account); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to grant role %s to account %s on chain %d: %w", governor.Role, governor.Account, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// RenounceRoleTokenGovernor renounces the given role from the given account on the given chains.
func RenounceRoleTokenGovernor(env cldf.Environment, c TokenGovernorRoleChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorRoleChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor role renounce")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token, governor := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]

			role, err := GetRoleFromTokenGovernor(env.GetContext(), tokenGovernor, governor.Role)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			hasRole, err := tokenGovernor.HasRole(&bind.CallOpts{Context: env.GetContext()}, role, governor.Account)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if account has role: %w", err)
			}

			if !hasRole {
				return cldf.ChangesetOutput{}, fmt.Errorf("account %s does not have role %s", governor.Account, governor.Role)
			}

			if _, err := tokenGovernor.RenounceRole(opts, role, governor.Account); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to renounce role %s from account %s on chain %d: %w", governor.Role, governor.Account, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// TransferOwnershipTokenGovernor transfers ownership of the token governor contracts to the given account on the given chains.
// It is assumed that the deployer has DEFAULT_ADMIN_ROLE.
func TransferOwnershipTokenGovernor(env cldf.Environment, c TokenGovernorRoleChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorRoleChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor transfer ownership")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token, governor := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]

			if _, err := tokenGovernor.TransferOwnership(opts, governor.Account); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to account %s on chain %d: %w", governor.Account, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// AcceptOwnershipTokenGovernor accepts ownership of the token governor contracts on the given chains.
// It is assumed that the deployer has DEFAULT_ADMIN_ROLE.
func AcceptOwnershipTokenGovernor(env cldf.Environment, c TokenGovernorChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor accept ownership")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]
			if _, err := tokenGovernor.AcceptOwnership(opts); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept ownership on chain %d: %w", chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// BeingDefaultAdminTransferTokenGovernor transfers ownership of the token governor contracts to the given account on the given chains.
func BeingDefaultAdminTransferTokenGovernor(env cldf.Environment, c TokenGovernorRoleChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorRoleChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor transfer adming ownership")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token, governor := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]

			if _, err := tokenGovernor.BeginDefaultAdminTransfer(opts, governor.Account); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to account %s on chain %d: %w", governor.Account, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// AcceptDefaultAdminTransferTokenGovernor accepts ownership of the token governor contracts on the given chains.
func AcceptDefaultAdminTransferTokenGovernor(env cldf.Environment, c TokenGovernorChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(env)
	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("token governor accept admin ownership")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)

		for token := range tokens {
			tokenGovernor := chainState.TokenGovernor[token]
			if _, err := tokenGovernor.AcceptDefaultAdminTransfer(opts); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept ownership on chain %d: %w", chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}
