package v1_5

import (
	"fmt"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/1_5_0/burn_mint_erc20_pausable_freezable_transparent"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSet[BurnMintERC20PausableFreezableTransparentChangesetConfig] = DeployBurnMintERC20PausableFreezableTransparent

type BurnMintERC20PausableFreezableTransparentChangesetConfig struct {
	Tokens map[uint64][]string
}

func (c BurnMintERC20PausableFreezableTransparentChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for _, token := range tokens {
			chain, ok := e.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			chainState, ok := state.EVMChainState(chainSelector)
			if !ok {
				return fmt.Errorf("%s does not exist in state", chain)
			}

			if _, ok := chainState.BurnMintERC20PausableFreezableTransparent[shared.TokenSymbol(token)]; ok {
				return fmt.Errorf("BurnMintERC20PausableFreezableTransparent already exists for %s token on %s", token, chain)
			}
		}
	}

	return nil
}

// DeployBurnMintERC20PausableFreezableTransparent deploys a BurnMintERC20PausableFreezableTransparent contract for each token specified in the config on the respective chain.
func DeployBurnMintERC20PausableFreezableTransparent(e cldf.Environment, c BurnMintERC20PausableFreezableTransparentChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid BurnMintERC20PausableFreezableTransparentChangesetConfig: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		for _, token := range tokens {
			_, err := cldf.DeployContract(e.Logger, chain, addressBook,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc20_pausable_freezable_transparent.BurnMintERC20PausableFreezableTransparent] {
					address, tx, transparent, err := burn_mint_erc20_pausable_freezable_transparent.DeployBurnMintERC20PausableFreezableTransparent(chain.DeployerKey, chain.Client)
					return cldf.ContractDeploy[*burn_mint_erc20_pausable_freezable_transparent.BurnMintERC20PausableFreezableTransparent]{
						Address:  address,
						Contract: transparent,
						Tv:       cldf.NewTypeAndVersion(shared.BurnMintERC20PausableFreezableTransparentToken, deployment.Version1_5_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)

			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BurnMintERC20PausableFreezableTransparent for %s token on %s: %w", token, chain, err)
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
