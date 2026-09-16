package v1_5

import (
	"fmt"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

// PlannedRefs returns every datastore key this changeset will write.
func (c BurnMintERC20PausableFreezableTransparentChangesetConfig) PlannedRefs() []datastore.AddressRef {
	version := deployment.Version1_5_0
	refs := make([]datastore.AddressRef, 0)
	for chainSelector, tokens := range c.Tokens {
		for _, token := range tokens {
			refs = append(refs, datastore.AddressRef{
				ChainSelector: chainSelector,
				Type:          datastore.ContractType(shared.BurnMintERC20PausableFreezableTransparentToken),
				Version:       &version,
				Qualifier:     token,
			})
		}
	}
	return refs
}

func (c BurnMintERC20PausableFreezableTransparentChangesetConfig) reserveRefs() (*datastore.MemoryDataStore, error) {
	return shared.ReserveRefs(c.PlannedRefs())
}

func (c BurnMintERC20PausableFreezableTransparentChangesetConfig) Validate(e cldf.Environment) error {
	if _, err := c.reserveRefs(); err != nil {
		return fmt.Errorf("pausable token datastore refs conflict: %w", err)
	}
	if err := shared.ValidateAddressRefsStrict(e, c.PlannedRefs()); err != nil {
		return fmt.Errorf("pausable token datastore refs conflict: %w", err)
	}

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
	ds, err := c.reserveRefs()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("pausable token datastore refs conflict: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		for _, token := range tokens {
			_, err := shared.DeployContractAndRecord(e.Logger, chain, addressBook, ds, cldf.NewTypeAndVersion(shared.BurnMintERC20PausableFreezableTransparentToken, deployment.Version1_5_0), token,
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

	return cldf.ChangesetOutput{
		AddressBook: addressBook, //nolint:staticcheck // SA1019 AddressBook is deprecated
		DataStore:   ds,
	}, nil
}
