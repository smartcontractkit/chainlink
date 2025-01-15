package example

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type MintLinkConfig struct {
	Amount        *big.Int
	ChainSelector uint64
	To            common.Address
}

var _ deployment.ChangeSet[*MintLinkConfig] = MintLink

// MintLink mints LINK to the provided contract.
func MintLink(e deployment.Environment, cfg *MintLinkConfig) (deployment.ChangesetOutput, error) {
	chain := e.Chains[cfg.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	tx, err := linkState.LinkToken.Mint(chain.DeployerKey, cfg.To, cfg.Amount)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}
