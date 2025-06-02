package example

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type AddMintersBurnersLinkConfig struct {
	ChainSelector uint64
	Minters       []common.Address
	Burners       []common.Address
}

var _ cldf.ChangeSet[*AddMintersBurnersLinkConfig] = AddMintersBurnersLink

// AddMintersBurnersLink grants the minter / burner role to the provided addresses.
func AddMintersBurnersLink(e cldf.Environment, cfg *AddMintersBurnersLinkConfig) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.EVMChains()[cfg.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for _, minter := range cfg.Minters {
		// check if minter is already a minter
		isMinter, err := linkState.LinkToken.IsMinter(&bind.CallOpts{Context: e.GetContext()}, minter)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if isMinter {
			continue
		}
		tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, minter)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		_, err = cldf.ConfirmIfNoError(chain, tx, err)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	for _, burner := range cfg.Burners {
		// check if burner is already a burner
		isBurner, err := linkState.LinkToken.IsBurner(&bind.CallOpts{Context: e.GetContext()}, burner)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if isBurner {
			continue
		}
		tx, err := linkState.LinkToken.GrantBurnRole(chain.DeployerKey, burner)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		_, err = cldf.ConfirmIfNoError(chain, tx, err)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	return cldf.ChangesetOutput{}, nil
}
