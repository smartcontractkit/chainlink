package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	DeployNonceManagerOp = operations.NewOperation(
		"DeployNonceManager",
		semver.MustParse("1.0.0"),
		"Deploys NonceManager 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input uint64) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain
			nonceManager, err := cldf.DeployContract(b.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*nonce_manager.NonceManager] {
					nonceManagerAddr, tx2, nonceManager, err2 := nonce_manager.DeployNonceManager(
						chain.DeployerKey,
						chain.Client,
						[]common.Address{}, // Need to add onRamp after
					)
					return cldf.ContractDeploy[*nonce_manager.NonceManager]{
						Address: nonceManagerAddr, Contract: nonceManager, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.NonceManager, deployment.Version1_6_0), Err: err2,
					}
				})
			if err != nil {
				b.Logger.Errorw("Failed to deploy nonce manager", "chain", chain.String(), "err", err)
				return common.Address{}, err
			}
			return nonceManager.Address, nil
		})

	NonceManagerUpdateAuthorizedCallerOp = operations.NewOperation(
		"NonceManagerUpdateAuthorizedCaller",
		semver.MustParse("1.0.0"),
		"Update authorized callers in NonceManager 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.ConfigureDependencies, input NonceManagerUpdateAuthorizedCallerInput) (opsutil.OpOutput, error) {
			state := deps.CurrentState
			e := deps.Env
			err := input.Validate(e, state)
			if err != nil {
				return opsutil.OpOutput{}, err
			}
			chain := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
			chainState := state.MustGetEVMChainState(input.ChainSelector)
			deployerGroup := deployergroup.NewDeployerGroup(e, state, input.MCMS).
				WithDeploymentContext("set NonceManager authorized caller on " + chain.String())
			opts, err := deployerGroup.GetDeployer(input.ChainSelector)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
			}
			nonceManager := chainState.NonceManager
			_, err = nonceManager.ApplyAuthorizedCallerUpdates(opts, input.Callers)
			if err != nil {
				b.Logger.Errorw("Failed to apply authorized caller updates on NonceManager", "chain", chain.String(), "err", err)
				return opsutil.OpOutput{}, fmt.Errorf("failed to apply authorized caller updates on NonceManager: %w", err)
			}
			csOutput, err := deployerGroup.Enact()
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to apply authorized caller updates on NonceManager: %w", err)
			}
			return opsutil.OpOutput{
				Proposals:                  csOutput.MCMSTimelockProposals,
				DescribedTimelockProposals: csOutput.DescribedTimelockProposals,
			}, nil
		})
)

type NonceManagerUpdateAuthorizedCallerInput struct {
	ChainSelector uint64
	Callers       nonce_manager.AuthorizedCallersAuthorizedCallerArgs
	MCMS          *proposalutils.TimelockConfig
}

func (n NonceManagerUpdateAuthorizedCallerInput) Validate(env cldf.Environment, state stateview.CCIPOnChainState) error {
	err := stateview.ValidateChain(env, state, n.ChainSelector, n.MCMS)
	if err != nil {
		return err
	}
	chain := env.BlockChains.EVMChains()[n.ChainSelector]
	if state.MustGetEVMChainState(n.ChainSelector).NonceManager == nil {
		return fmt.Errorf("NonceManager not found for chain %s", chain.String())
	}
	err = commoncs.ValidateOwnership(
		env.GetContext(), n.MCMS != nil,
		chain.DeployerKey.From, state.MustGetEVMChainState(n.ChainSelector).Timelock.Address(),
		state.MustGetEVMChainState(n.ChainSelector).NonceManager,
	)
	if err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if len(n.Callers.AddedCallers) == 0 && len(n.Callers.RemovedCallers) == 0 {
		return errors.New("at least one caller is required")
	}
	return nil
}
