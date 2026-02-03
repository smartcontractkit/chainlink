package v1_5

import (
	"fmt"
	"math/big"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/1_5_0/burn_mint_erc20_transparent"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSet[BurnMintERC20TransparentChangesetConfig] = DeployBurnMintERC20Transparent

type BurnMintERC20Transparent struct {
	MaxSupply *big.Int
	PreMint   *big.Int
}

type BurnMintERC20TransparentChangesetConfig struct {
	Tokens map[uint64]map[string]BurnMintERC20Transparent
}

func (c BurnMintERC20TransparentChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, config := range tokens {
			if config.PreMint == nil {
				config.PreMint = big.NewInt(0)
				tokens[token] = config
			}

			if config.MaxSupply == nil {
				return fmt.Errorf("max supply is required for %s token on chain %d", token, chainSelector)
			}

			chain, ok := e.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			chainState, ok := state.EVMChainState(chainSelector)
			if !ok {
				return fmt.Errorf("%s does not exist in state", chain)
			}

			if _, ok := chainState.BurnMintERC20Transparent[shared.TokenSymbol(token)]; ok {
				return fmt.Errorf("BurnMintERC20Transparent already exists for %s token on %s", token, chain)
			}
		}
	}

	return nil
}

// DeployBurnMintERC20Transparent deploys BurnMintERC20Transparent contracts for the specified tokens on each chain.
func DeployBurnMintERC20Transparent(e cldf.Environment, c BurnMintERC20TransparentChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid BurnMintERC20TransparentChangesetConfig: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		for token := range tokens {
			_, err := cldf.DeployContract(e.Logger, chain, addressBook,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc20_transparent.BurnMintERC20Transparent] {
					address, tx, transparent, err := burn_mint_erc20_transparent.DeployBurnMintERC20Transparent(chain.DeployerKey, chain.Client)
					return cldf.ContractDeploy[*burn_mint_erc20_transparent.BurnMintERC20Transparent]{
						Address:  address,
						Contract: transparent,
						Tv:       cldf.NewTypeAndVersion(shared.BurnMintERC20TransparentToken, deployment.Version1_5_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)

			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BurnMintERC20Transparent for %s token on %s: %w", token, chain, err)
			}
		}
	}

	ds, err := shared.PopulateDataStore(addressBook)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
		DataStore:   ds,
	}, nil
}
