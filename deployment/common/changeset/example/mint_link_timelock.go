package example

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type MintLinkTimelockConfig struct {
	Amount        *big.Int
	ChainSelector uint64
}

var _ deployment.ChangeSet[*MintLinkTimelockConfig] = MintLinkTimelock

// MintLinkTimelock mints LINK to the timelock contract.
func MintLinkTimelock(e deployment.Environment, req *MintLinkTimelockConfig) (deployment.ChangesetOutput, error) {

	chain := e.Chains[req.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(req.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	linkState, err := changeset.MaybeLoadLinkTokenState(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	timelockState, err := changeset.MaybeLoadMCMSWithTimelockState(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, timelockState.Timelock.Address(), req.Amount)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil

}
