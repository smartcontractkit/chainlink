package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
)

type ownershipTransferrer interface {
	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*gethtypes.Transaction, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
}

type TransferOwnershipConfig struct {
	State             CCIPOnChainState
	ChainSelectors    []uint64
	HomeChainSelector uint64
}

var _ deployment.ChangeSet[TransferOwnershipConfig] = NewTransferOwnershipChangeset

// NewTransferOwnershipChangeset creates a changeset that transfers ownership of all the
// ccip chain contracts deployed on the given chain selectors.
// New chain contracts are:
// * OnRamp
// * OffRamp
// * FeeQuoter
// * NonceManager
// * RMNRemote
// Home chain contracts are:
// * CCIPHome
// * RMNHome
// * CapabilityRegistry
// This can be composed with NewAcceptOwnershipChangeset in order to fully transfer
// ownership of all the contracts listed above.
func NewTransferOwnershipChangeset(
	e deployment.Environment,
	cfg TransferOwnershipConfig,
) (deployment.ChangesetOutput, error) {
	// basic validation
	if len(cfg.ChainSelectors) == 0 || cfg.HomeChainSelector == 0 {
		return deployment.ChangesetOutput{}, fmt.Errorf("no chain selectors provided")
	}

	if len(cfg.State.Chains) == 0 {
		return deployment.ChangesetOutput{}, fmt.Errorf("no chains in state")
	}

	// transfer ownership of chain contracts
	// these are assumed to be owned by the deployer configured in the given
	// environment.
	for _, chain := range cfg.ChainSelectors {
		for _, contract := range []ownershipTransferrer{
			cfg.State.Chains[chain].OnRamp,
			cfg.State.Chains[chain].OffRamp,
			cfg.State.Chains[chain].FeeQuoter,
			cfg.State.Chains[chain].NonceManager,
			cfg.State.Chains[chain].RMNRemote,
		} {
			tx, err := contract.TransferOwnership(
				e.Chains[chain].DeployerKey,
				cfg.State.Chains[chain].Timelock.Address(),
			)
			_, err = deployment.ConfirmIfNoError(e.Chains[chain], tx, err)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	// transfer ownership of home chain contracts
	homeChainTimelockAddress := cfg.State.Chains[cfg.HomeChainSelector].Timelock.Address()
	for _, contract := range []ownershipTransferrer{
		cfg.State.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		cfg.State.Chains[cfg.HomeChainSelector].CCIPHome,
		cfg.State.Chains[cfg.HomeChainSelector].RMNHome,
	} {
		tx, err := contract.TransferOwnership(
			e.Chains[cfg.HomeChainSelector].DeployerKey,
			homeChainTimelockAddress,
		)
		_, err = deployment.ConfirmIfNoError(e.Chains[cfg.HomeChainSelector], tx, err)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	// no new addresses or proposals or jobspecs, so changeset output is empty.
	// NOTE: onchain state has technically changed for above contracts, maybe that should
	// be captured?
	return deployment.ChangesetOutput{}, nil
}
