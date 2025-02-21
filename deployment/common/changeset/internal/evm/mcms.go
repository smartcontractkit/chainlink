package mcmsnew

import (
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	evmMcms "github.com/smartcontractkit/mcms/sdk/evm"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
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

// MCMSWithTimelockEVMDeploy holds a bundle of MCMS contract deploys.
type MCMSWithTimelockEVMDeploy struct {
	Canceller *deployment.ContractDeploy[*bindings.ManyChainMultiSig]
	Bypasser  *deployment.ContractDeploy[*bindings.ManyChainMultiSig]
	Proposer  *deployment.ContractDeploy[*bindings.ManyChainMultiSig]
	Timelock  *deployment.ContractDeploy[*bindings.RBACTimelock]
	CallProxy *deployment.ContractDeploy[*bindings.CallProxy]
}

func deployMCMSWithConfigEVM(
	contractType deployment.ContractType,
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	mcmConfig mcmsTypes.Config,
	options ...DeployMCMSOption,
) (*deployment.ContractDeploy[*bindings.ManyChainMultiSig], error) {
	groupQuorums, groupParents, signerAddresses, signerGroups, err := evmMcms.ExtractSetConfigInputs(&mcmConfig)
	if err != nil {
		lggr.Errorw("Failed to extract set config inputs", "chain", chain.String(), "err", err)
		return nil, err
	}
	mcm, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*bindings.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := bindings.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)

			tv := deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0)
			for _, option := range options {
				option(&tv)
			}

			return deployment.ContractDeploy[*bindings.ManyChainMultiSig]{
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
	if _, err := deployment.ConfirmIfNoError(chain, mcmsTx, err); err != nil {
		lggr.Errorw("Failed to confirm mcm config", "chain", chain.String(), "err", err)
		return mcm, err
	}
	return mcm, nil
}

// DeployMCMSWithTimelockContractsEVM deploys an MCMS for
// each of the timelock roles Bypasser, ProposerMcm, Canceller on an EVM chain.
// MCMS contracts for the given configuration
// as well as the timelock. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func DeployMCMSWithTimelockContractsEVM(
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	config commontypes.MCMSWithTimelockConfigV2,
) (*MCMSWithTimelockEVMDeploy, error) {
	opts := []DeployMCMSOption{}
	if config.Label != nil {
		opts = append(opts, WithLabel(*config.Label))
	}

	bypasser, err := deployMCMSWithConfigEVM(commontypes.BypasserManyChainMultisig, lggr, chain, ab, config.Bypasser, opts...)
	if err != nil {
		return nil, err
	}
	canceller, err := deployMCMSWithConfigEVM(commontypes.CancellerManyChainMultisig, lggr, chain, ab, config.Canceller, opts...)
	if err != nil {
		return nil, err
	}
	proposer, err := deployMCMSWithConfigEVM(commontypes.ProposerManyChainMultisig, lggr, chain, ab, config.Proposer, opts...)
	if err != nil {
		return nil, err
	}

	timelock, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*bindings.RBACTimelock] {
			timelock, tx2, cc, err2 := bindings.DeployRBACTimelock(
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

			tv := deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0)
			if config.Label != nil {
				tv.AddLabel(*config.Label)
			}

			return deployment.ContractDeploy[*bindings.RBACTimelock]{
				Address: timelock, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy timelock", "chain", chain.String(), "err", err)
		return nil, err
	}

	callProxy, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*bindings.CallProxy] {
			callProxy, tx2, cc, err2 := bindings.DeployCallProxy(
				chain.DeployerKey,
				chain.Client,
				timelock.Address,
			)

			tv := deployment.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0)
			if config.Label != nil {
				tv.AddLabel(*config.Label)
			}

			return deployment.ContractDeploy[*bindings.CallProxy]{
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
	return &MCMSWithTimelockEVMDeploy{
		Canceller: canceller,
		Bypasser:  bypasser,
		Proposer:  proposer,
		Timelock:  timelock,
		CallProxy: callProxy,
	}, nil
}
