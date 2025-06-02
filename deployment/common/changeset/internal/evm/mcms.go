package mcmsnew

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	evmMcms "github.com/smartcontractkit/mcms/sdk/evm"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

// DeployMCMSOption is a function that modifies a TypeAndVersion before or after deployment.
type DeployMCMSOption func(*cldf.TypeAndVersion)

// WithLabel is a functional option that sets a label on the TypeAndVersion.
func WithLabel(label string) DeployMCMSOption {
	return func(tv *cldf.TypeAndVersion) {
		tv.AddLabel(label)
	}
}

// MCMSWithTimelockEVMDeploy holds a bundle of MCMS contract deploys.
type MCMSWithTimelockEVMDeploy struct {
	Canceller *cldf.ContractDeploy[*bindings.ManyChainMultiSig]
	Bypasser  *cldf.ContractDeploy[*bindings.ManyChainMultiSig]
	Proposer  *cldf.ContractDeploy[*bindings.ManyChainMultiSig]
	Timelock  *cldf.ContractDeploy[*bindings.RBACTimelock]
	CallProxy *cldf.ContractDeploy[*bindings.CallProxy]
}

func DeployMCMSWithConfigEVM(
	contractType cldf.ContractType,
	lggr logger.Logger,
	chain cldf_evm.Chain,
	ab cldf.AddressBook,
	mcmConfig mcmsTypes.Config,
	options ...DeployMCMSOption,
) (*cldf.ContractDeploy[*bindings.ManyChainMultiSig], error) {
	groupQuorums, groupParents, signerAddresses, signerGroups, err := evmMcms.ExtractSetConfigInputs(&mcmConfig)
	if err != nil {
		lggr.Errorw("Failed to extract set config inputs", "chain", chain.String(), "err", err)
		return nil, err
	}
	mcm, err := cldf.DeployContract(lggr, chain, ab,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := bindings.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)

			tv := cldf.NewTypeAndVersion(contractType, deployment.Version1_0_0)
			for _, option := range options {
				option(&tv)
			}

			return cldf.ContractDeploy[*bindings.ManyChainMultiSig]{
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
	if _, err := cldf.ConfirmIfNoError(chain, mcmsTx, err); err != nil {
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
	ctx context.Context,
	lggr logger.Logger,
	chain cldf_evm.Chain,
	ab cldf.AddressBook,
	config commontypes.MCMSWithTimelockConfigV2,
	state *state.MCMSWithTimelockState,
) (*proposalutils.MCMSWithTimelockContracts, error) {
	opts := []DeployMCMSOption{}
	if config.Label != nil {
		opts = append(opts, WithLabel(*config.Label))
	}
	var bypasser, proposer, canceller *bindings.ManyChainMultiSig
	var timelock *bindings.RBACTimelock
	var callProxy *bindings.CallProxy
	if state != nil {
		bypasser = state.BypasserMcm
		proposer = state.ProposerMcm
		canceller = state.CancellerMcm
		timelock = state.Timelock
		callProxy = state.CallProxy
	}
	if bypasser == nil {
		bypasserC, err := DeployMCMSWithConfigEVM(commontypes.BypasserManyChainMultisig, lggr, chain, ab, config.Bypasser, opts...)
		if err != nil {
			return nil, err
		}
		bypasser = bypasserC.Contract
		lggr.Infow("Bypasser MCMS deployed", "chain", chain.String(), "address", bypasser.Address)
	} else {
		lggr.Infow("Bypasser MCMS already deployed", "chain", chain.String(), "address", bypasser.Address)
	}

	if canceller == nil {
		cancellerC, err := DeployMCMSWithConfigEVM(commontypes.CancellerManyChainMultisig, lggr, chain, ab, config.Canceller, opts...)
		if err != nil {
			return nil, err
		}
		canceller = cancellerC.Contract
		lggr.Infow("Canceller MCMS deployed", "chain", chain.String(), "address", canceller.Address)
	} else {
		lggr.Infow("Canceller MCMS already deployed", "chain", chain.String(), "address", canceller.Address)
	}

	if proposer == nil {
		proposerC, err := DeployMCMSWithConfigEVM(commontypes.ProposerManyChainMultisig, lggr, chain, ab, config.Proposer, opts...)
		if err != nil {
			return nil, err
		}
		proposer = proposerC.Contract
		lggr.Infow("Proposer MCMS deployed", "chain", chain.String(), "address", proposer.Address)
	} else {
		lggr.Infow("Proposer MCMS already deployed", "chain", chain.String(), "address", proposer.Address)
	}

	if timelock == nil {
		timelockC, err := cldf.DeployContract(lggr, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.RBACTimelock] {
				timelock, tx2, cc, err2 := bindings.DeployRBACTimelock(
					chain.DeployerKey,
					chain.Client,
					config.TimelockMinDelay,
					// Deployer is the initial admin.
					// TODO: Could expose this as config?
					// Or keep this enforced to follow the same pattern?
					chain.DeployerKey.From,
					[]common.Address{proposer.Address()}, // proposers
					// Executors field is empty here because we grant the executor role to the call proxy later
					// and the call proxy cannot be deployed before the timelock.
					[]common.Address{},
					[]common.Address{canceller.Address(), proposer.Address(), bypasser.Address()}, // cancellers
					[]common.Address{bypasser.Address()},                                          // bypassers
				)

				tv := cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0)
				if config.Label != nil {
					tv.AddLabel(*config.Label)
				}

				return cldf.ContractDeploy[*bindings.RBACTimelock]{
					Address: timelock, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy timelock", "chain", chain.String(), "err", err)
			return nil, err
		}
		timelock = timelockC.Contract
		lggr.Infow("Timelock deployed", "chain", chain.String(), "address", timelock.Address)
	} else {
		lggr.Infow("Timelock already deployed", "chain", chain.String(), "address", timelock.Address)
	}

	if callProxy == nil {
		callProxyC, err := cldf.DeployContract(lggr, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.CallProxy] {
				callProxy, tx2, cc, err2 := bindings.DeployCallProxy(
					chain.DeployerKey,
					chain.Client,
					timelock.Address(),
				)

				tv := cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0)
				if config.Label != nil {
					tv.AddLabel(*config.Label)
				}

				return cldf.ContractDeploy[*bindings.CallProxy]{
					Address: callProxy, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy call proxy", "chain", chain.String(), "err", err)
			return nil, err
		}
		callProxy = callProxyC.Contract
		lggr.Infow("CallProxy deployed", "chain", chain.String(), "address", callProxy.Address)
	} else {
		lggr.Infow("CallProxy already deployed", "chain", chain.String(), "address", callProxy.Address)
	}
	timelockContracts := &proposalutils.MCMSWithTimelockContracts{
		BypasserMcm:  bypasser,
		ProposerMcm:  proposer,
		CancellerMcm: canceller,
		Timelock:     timelock,
		CallProxy:    callProxy,
	}
	// grant roles for timelock
	// this is called only if deployer key is an admin in timelock
	_, err := GrantRolesForTimelock(ctx, lggr, chain, timelockContracts, true)
	if err != nil {
		return nil, err
	}
	// After the proposer cycle is validated,
	// we can remove the deployer as an admin.
	return timelockContracts, nil
}

// TODO: delete this function after it is available in timelock Inspector
func getAdminAddresses(ctx context.Context, timelock *bindings.RBACTimelock) ([]string, error) {
	numAddresses, err := timelock.GetRoleMemberCount(&bind.CallOpts{
		Context: ctx,
	}, v1_0.ADMIN_ROLE.ID)
	if err != nil {
		return nil, err
	}
	adminAddresses := make([]string, 0, numAddresses.Uint64())
	for i := range numAddresses.Uint64() {
		if i > math.MaxUint32 {
			return nil, fmt.Errorf("value %d exceeds uint32 range", i)
		}
		idx, err := cast.ToInt64E(i)
		if err != nil {
			return nil, err
		}
		address, err := timelock.GetRoleMember(&bind.CallOpts{
			Context: ctx,
		}, v1_0.ADMIN_ROLE.ID, big.NewInt(idx))
		if err != nil {
			return nil, err
		}
		adminAddresses = append(adminAddresses, address.String())
	}
	return adminAddresses, nil
}

func GrantRolesForTimelock(
	ctx context.Context,
	lggr logger.Logger,
	chain cldf_evm.Chain,
	timelockContracts *proposalutils.MCMSWithTimelockContracts,
	skipIfDeployerKeyNotAdmin bool, // If true, skip role grants if the deployer key is not an admin.
) ([]mcmsTypes.Transaction, error) {
	if timelockContracts == nil {
		lggr.Errorw("Timelock contracts not found", "chain", chain.String())
		return nil, fmt.Errorf("timelock contracts not found for chain %s", chain.String())
	}

	timelock := timelockContracts.Timelock
	proposer := timelockContracts.ProposerMcm
	canceller := timelockContracts.CancellerMcm
	bypasser := timelockContracts.BypasserMcm
	callProxy := timelockContracts.CallProxy

	// get admin addresses
	adminAddresses, err := getAdminAddresses(ctx, timelock)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin addresses: %w", err)
	}
	isDeployerKeyAdmin := slices.Contains(adminAddresses, chain.DeployerKey.From.String())
	isTimelockAdmin := slices.Contains(adminAddresses, timelock.Address().String())
	if !isDeployerKeyAdmin && skipIfDeployerKeyNotAdmin {
		lggr.Infow("Deployer key is not admin, skipping role grants", "chain", chain.String())
		return nil, nil
	}
	if !isDeployerKeyAdmin && !isTimelockAdmin {
		return nil, errors.New("neither deployer key nor timelock is admin, cannot grant roles")
	}

	var mcmsTxs []mcmsTypes.Transaction

	timelockInspector := evmMcms.NewTimelockInspector(chain.Client)
	proposerAddresses, err := timelockInspector.GetProposers(ctx, timelock.Address().String())
	if err != nil {
		return nil, fmt.Errorf("failed to get proposers for chain %s: %w", chain.String(), err)
	}
	if !slices.Contains(proposerAddresses, proposer.Address().String()) {
		tx, err := grantRoleTx(lggr, timelock, chain, isDeployerKeyAdmin, v1_0.PROPOSER_ROLE.ID, proposer.Address())
		if err != nil {
			lggr.Errorw("Failed to grant timelock proposer role", "chain", chain.String(), "err", err)
			return nil, err
		}
		if !isDeployerKeyAdmin {
			mcmsTxs = append(mcmsTxs, tx)
		} else {
			lggr.Infow("Proposer role granted", "chain", chain.String(), "address", proposer.Address())
		}
	}
	cancellerAddresses, err := timelockInspector.GetCancellers(ctx, timelock.Address().String())
	if err != nil {
		return nil, fmt.Errorf("failed to get proposers for chain %s: %w", chain.String(), err)
	}
	for _, addr := range []common.Address{proposer.Address(), canceller.Address(), bypasser.Address()} {
		if !slices.Contains(cancellerAddresses, addr.String()) {
			tx, err := grantRoleTx(lggr, timelock, chain, isDeployerKeyAdmin, v1_0.CANCELLER_ROLE.ID, addr)
			if err != nil {
				lggr.Errorw("Failed to grant timelock canceller role", "chain", chain.String(), "err", err)
				return nil, err
			}
			if !isDeployerKeyAdmin {
				mcmsTxs = append(mcmsTxs, tx)
			} else {
				lggr.Infow("Canceller role granted", "chain", chain.String(), "address", addr)
			}
		}
	}

	bypasserAddresses, err := timelockInspector.GetBypassers(ctx, timelock.Address().String())
	if err != nil {
		return nil, fmt.Errorf("failed to get bypassers for chain %s: %w", chain.String(), err)
	}
	if !slices.Contains(bypasserAddresses, bypasser.Address().String()) {
		tx, err := grantRoleTx(lggr, timelock, chain, isDeployerKeyAdmin, v1_0.BYPASSER_ROLE.ID, bypasser.Address())
		if err != nil {
			lggr.Errorw("Failed to grant timelock bypasser role", "chain", chain.String(), "err", err)
			return nil, err
		}
		if !isDeployerKeyAdmin {
			mcmsTxs = append(mcmsTxs, tx)
		} else {
			lggr.Infow("Bypasser role granted", "chain", chain.String(), "address", bypasser.Address())
		}
	}
	executorAddresses, err := timelockInspector.GetExecutors(ctx, timelock.Address().String())
	if err != nil {
		return nil, fmt.Errorf("failed to get executors for chain %s: %w", chain.String(), err)
	}
	if !slices.Contains(executorAddresses, callProxy.Address().String()) {
		tx, err := grantRoleTx(lggr, timelock, chain, isDeployerKeyAdmin, v1_0.EXECUTOR_ROLE.ID, callProxy.Address())
		if err != nil {
			lggr.Errorw("Failed to grant timelock executor role", "chain", chain.String(), "err", err)
			return nil, err
		}
		if !isDeployerKeyAdmin {
			mcmsTxs = append(mcmsTxs, tx)
		} else {
			lggr.Infow("Executor role granted", "chain", chain.String(), "address", callProxy.Address())
		}
	}

	if !isTimelockAdmin {
		// We grant the timelock the admin role on the MCMS contracts.
		tx, err := grantRoleTx(lggr, timelock, chain, isDeployerKeyAdmin, v1_0.ADMIN_ROLE.ID, timelock.Address())
		if err != nil {
			lggr.Errorw("Failed to grant timelock admin role", "chain", chain.String(), "err", err)
			return nil, err
		}
		if !isDeployerKeyAdmin {
			mcmsTxs = append(mcmsTxs, tx)
		} else {
			lggr.Infow("Admin role granted", "chain", chain.String(), "address", timelock.Address())
		}
	}
	return mcmsTxs, nil
}

func grantRoleTx(
	lggr logger.Logger,
	timelock *bindings.RBACTimelock,
	chain cldf_evm.Chain,
	isDeployerKeyAdmin bool,
	roleID [32]byte,
	address common.Address,
) (mcmsTypes.Transaction, error) {
	txOpts := cldf.SimTransactOpts()
	if isDeployerKeyAdmin {
		txOpts = chain.DeployerKey
	}
	grantRoleTx, err := timelock.GrantRole(
		txOpts, roleID, address,
	)
	if isDeployerKeyAdmin {
		if _, err2 := cldf.ConfirmIfNoErrorWithABI(chain, grantRoleTx, bindings.RBACTimelockABI, err); err != nil {
			lggr.Errorw("Failed to grant timelock role",
				"chain", chain.String(),
				"timelock", timelock.Address().Hex(),
				"Address to grant role", address.Hex(),
				"err", err2)
			return mcmsTypes.Transaction{}, err2
		}
		return mcmsTypes.Transaction{}, err
	}
	if err != nil {
		lggr.Errorw("Failed to grant timelock role",
			"chain", chain.String(),
			"timelock", timelock.Address().Hex(),
			"Address to grant role", address.Hex(),
			"err", err)
		return mcmsTypes.Transaction{}, err
	}
	tx, err := proposalutils.TransactionForChain(chain.Selector, timelock.Address().Hex(), grantRoleTx.Data(),
		big.NewInt(0), commontypes.RBACTimelock.String(), []string{})
	if err != nil {
		lggr.Errorw("Failed to create transaction for chain",
			"chain", chain.String(),
			"timelock", timelock.Address().Hex(),
			"Address to grant role", address.Hex(),
			"err", err)
		return mcmsTypes.Transaction{}, err
	}
	return tx, nil
}
