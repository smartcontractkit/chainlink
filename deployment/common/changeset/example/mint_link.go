package example

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type MintLinkConfig struct {
	Amount        *big.Int
	ChainSelector uint64
	To            common.Address
}

var _ cldf.ChangeSet[*MintLinkConfig] = MintLink

// MintLink mints LINK to the provided contract.
func MintLink(e cldf.Environment, cfg *MintLinkConfig) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.EVMChains()[cfg.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	tx, err := linkState.LinkToken.Mint(chain.DeployerKey, cfg.To, cfg.Amount)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{}, nil
}
