package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

type OwnershipTransferrer interface {
	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*gethtypes.Transaction, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
}

type TransferOwnershipConfig struct {
	// OwnersPerChain is a mapping from chain selector to the owner's contract address on that chain.
	OwnersPerChain map[uint64]common.Address

	// Contracts is a mapping from chain selector to the ownership transferrers on that chain.
	Contracts map[uint64][]OwnershipTransferrer
}

func (t TransferOwnershipConfig) Validate() error {
	// check that we have owners for the chains in the Contracts field.
	for chainSelector := range t.Contracts {
		if _, ok := t.OwnersPerChain[chainSelector]; !ok {
			return fmt.Errorf("missing owners for chain %d", chainSelector)
		}
	}

	return nil
}

var _ deployment.ChangeSet[TransferOwnershipConfig] = NewTransferOwnershipChangeset

// NewTransferOwnershipChangeset creates a changeset that transfers ownership of all the
// contracts in the provided configuration to correct owner on that chain.
// If the owner is already the provided address, no transaction is sent.
func NewTransferOwnershipChangeset(
	e deployment.Environment,
	cfg TransferOwnershipConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSelector, contracts := range cfg.Contracts {
		ownerAddress := cfg.OwnersPerChain[chainSelector]
		for _, contract := range contracts {
			owner, err := contract.Owner(nil)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get owner of contract %T: %v", contract, err)
			}
			if owner != ownerAddress {
				tx, err := contract.TransferOwnership(e.Chains[chainSelector].DeployerKey, ownerAddress)
				_, err = deployment.ConfirmIfNoError(e.Chains[chainSelector], tx, err)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of contract %T: %v", contract, err)
				}
			}
		}
	}

	// no new addresses or proposals or jobspecs, so changeset output is empty.
	// NOTE: onchain state has technically changed for above contracts, maybe that should
	// be captured?
	return deployment.ChangesetOutput{}, nil
}
