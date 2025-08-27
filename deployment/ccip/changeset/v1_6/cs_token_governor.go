package v1_6

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/latest/token_governor"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSet[TokenGovernorChangesetConfig] = DeployTokenGovernor

type TokenGovernor struct {
	Token common.Address
	// If nil, the default value will be used which is 0.
	InitialDelay *big.Int
	// If nil, the default value will be used which is the deployer's address.
	InitialDefaultAdmin common.Address
}

type TokenGovernorChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]TokenGovernor
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

			if governor.InitialDefaultAdmin == utils.ZeroAddress {
				return errors.New("initial default admin must be defined")
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

			if _, ok := chainState.TokenGovernor[token]; ok {
				return fmt.Errorf("token governor already exists for %s", governor.Token)
			}

			if err := validateTokenSymbol(context.Background(), chain, governor.Token, token); err != nil {
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

// DeployTokenGovernor deploys the token governor contracts on the given chains.
func DeployTokenGovernor(env cldf.Environment, c TokenGovernorChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenGovernorChangesetConfig: %w", err)
	}

	newAddresses := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := env.BlockChains.EVMChains()[chainSelector]

		for _, governor := range tokens {
			if governor.InitialDelay == nil {
				governor.InitialDelay = big.NewInt(0)
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
