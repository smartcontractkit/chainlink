package internal

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

// DeployMCMSOption is a function that modifies a TypeAndVersion before or after deployment.
type DeployMCMSOption func(*deployment.TypeAndVersion)

// WithLabel is a functional option that sets a label on the TypeAndVersion.
func WithLabel(label string) DeployMCMSOption {
	return func(tv *deployment.TypeAndVersion) {
		tv.AddLabel(label)
	}
}

func DeployMCMSWithConfig(
	contractType deployment.ContractType,
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	mcmConfig config.Config,
	options ...DeployMCMSOption,
) (*deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig], error) {
	groupQuorums, groupParents, signerAddresses, signerGroups := mcmConfig.ExtractSetConfigInputs()
	mcm, err := deployment.DeployContract[*owner_helpers.ManyChainMultiSig](lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := owner_helpers.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)

			tv := deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0)
			for _, option := range options {
				option(&tv)
			}

			return deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]{
				Address: mcmAddr, Contract: mcm, Tx: tx, Tv: tv, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy mcm", "chain", chain.String(), "err", err)
		return mcm, err
	}
	mcmsTx, err := mcm.Contract.SetConfig(chain.DeployerKey,
		signerAddresses,
		// Signer 1 is int group 0 (root group) with quorum 1.
		signerGroups,
		groupQuorums,
		groupParents,
		false,
	)
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, mcmsTx, owner_helpers.ManyChainMultiSigABI, err); err != nil {
		lggr.Errorw("Failed to confirm mcm config", "chain", chain.String(), "err", err)
		return mcm, err
	}
	return mcm, nil
}

// MCMSWithTimelockDeploy holds a bundle of MCMS contract deploys.
type MCMSWithTimelockDeploy struct {
	Canceller *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Bypasser  *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Proposer  *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Timelock  *deployment.ContractDeploy[*owner_helpers.RBACTimelock]
	CallProxy *deployment.ContractDeploy[*owner_helpers.CallProxy]
}

func DeployMCMSWithTimelockContractsBatch(
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	ab deployment.AddressBook,
	cfgByChain map[uint64]types.MCMSWithTimelockConfig,
) error {
	for chainSel, cfg := range cfgByChain {
		_, err := DeployMCMSWithTimelockContracts(lggr, chains[chainSel], ab, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeployMCMSWithTimelockContracts deploys an MCMS for
// each of the timelock roles Bypasser, ProposerMcm, Canceller.
// MCMS contracts for the given configuration
// as well as the timelock. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func DeployMCMSWithTimelockContracts(
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	config types.MCMSWithTimelockConfig,
) (*MCMSWithTimelockDeploy, error) {
	opts := []DeployMCMSOption{}
	if config.Label != nil {
		opts = append(opts, WithLabel(*config.Label))
	}

	bypasser, err := DeployMCMSWithConfig(types.BypasserManyChainMultisig, lggr, chain, ab, config.Bypasser, opts...)
	if err != nil {
		return nil, err
	}
	canceller, err := DeployMCMSWithConfig(types.CancellerManyChainMultisig, lggr, chain, ab, config.Canceller, opts...)
	if err != nil {
		return nil, err
	}
	proposer, err := DeployMCMSWithConfig(types.ProposerManyChainMultisig, lggr, chain, ab, config.Proposer, opts...)
	if err != nil {
		return nil, err
	}

	timelock, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*owner_helpers.RBACTimelock] {
			timelock, tx2, cc, err2 := owner_helpers.DeployRBACTimelock(
				chain.DeployerKey,
				chain.Client,
				config.TimelockMinDelay,
				// Deployer is the initial admin.
				// TODO: Could expose this as config?
				// Or keep this enforced to follow the same pattern?
				chain.DeployerKey.From,
				[]common.Address{proposer.Address}, // proposers
				// Executors field is empty here because we grant the executor role to the call proxy later
				// and the call proxy cannot be deployed before the timelock.
				[]common.Address{},
				[]common.Address{canceller.Address, proposer.Address, bypasser.Address}, // cancellers
				[]common.Address{bypasser.Address},                                      // bypassers
			)

			tv := deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
			if config.Label != nil {
				tv.AddLabel(*config.Label)
			}

			return deployment.ContractDeploy[*owner_helpers.RBACTimelock]{
				Address: timelock, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy timelock", "chain", chain.String(), "err", err)
		return nil, err
	}

	callProxy, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*owner_helpers.CallProxy] {
			callProxy, tx2, cc, err2 := owner_helpers.DeployCallProxy(
				chain.DeployerKey,
				chain.Client,
				timelock.Address,
			)

			tv := deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
			if config.Label != nil {
				tv.AddLabel(*config.Label)
			}

			return deployment.ContractDeploy[*owner_helpers.CallProxy]{
				Address: callProxy, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy call proxy", "chain", chain.String(), "err", err)
		return nil, err
	}

	grantRoleTx, err := timelock.Contract.GrantRole(
		chain.DeployerKey,
		v1_0.EXECUTOR_ROLE.ID,
		callProxy.Address,
	)
	if err != nil {
		lggr.Errorw("Failed to grant timelock executor role", "chain", chain.String(), "err", err)
		return nil, err
	}

	if _, err := deployment.ConfirmIfNoError(chain, grantRoleTx, err); err != nil {
		lggr.Errorw("Failed to grant timelock executor role", "chain", chain.String(), "err", err)
		return nil, err
	}
	// We grant the timelock the admin role on the MCMS contracts.
	tx, err := timelock.Contract.GrantRole(chain.DeployerKey,
		v1_0.ADMIN_ROLE.ID, timelock.Address)
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to grant timelock admin role", "chain", chain.String(), "err", err)
		return nil, err
	}
	// After the proposer cycle is validated,
	// we can remove the deployer as an admin.
	return &MCMSWithTimelockDeploy{
		Canceller: canceller,
		Bypasser:  bypasser,
		Proposer:  proposer,
		Timelock:  timelock,
		CallProxy: callProxy,
	}, nil
}
