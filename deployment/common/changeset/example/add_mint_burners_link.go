package example

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type AddMintersBurnersLinkConfig struct {
	ChainSelector uint64
	Minters       []common.Address
	Burners       []common.Address
}

var _ deployment.ChangeSet[*AddMintersBurnersLinkConfig] = AddMintersBurnersLink

// AddMintersBurnersLink grants the minter / burner role to the provided addresses.
func AddMintersBurnersLink(e deployment.Environment, cfg *AddMintersBurnersLinkConfig) (deployment.ChangesetOutput, error) {
	chain := e.Chains[cfg.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for _, minter := range cfg.Minters {
		// check if minter is already a minter
		isMinter, err := linkState.LinkToken.IsMinter(&bind.CallOpts{Context: e.GetContext()}, minter)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if isMinter {
			continue
		}
		tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, minter)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		_, err = deployment.ConfirmIfNoError(chain, tx, err)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}
	for _, burner := range cfg.Burners {
		// check if burner is already a burner
		isBurner, err := linkState.LinkToken.IsBurner(&bind.CallOpts{Context: e.GetContext()}, burner)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if isBurner {
			continue
		}
		tx, err := linkState.LinkToken.GrantBurnRole(chain.DeployerKey, burner)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		_, err = deployment.ConfirmIfNoError(chain, tx, err)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}
	return deployment.ChangesetOutput{}, nil
}
