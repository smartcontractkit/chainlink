package v1_6

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"
)

// ConnectNewChainChangeset activates connects a new chain with other chains by updating onRamp, offRamp, and router contracts.
// When running this changeset, the onRamp, offRamp, and router contracts on the new chain should NOT be owned by MCMS yet.
// If connecting to production routers, this changeset will ensure that ALL chain contracts on the new chain are transferred to MCMS.
// This changeset enforces that the onRamp, offRamp, and router contracts on other chains are already owned by MCMS, regardless of the desired router.
var ConnectNewChainChangeset = deployment.CreateChangeSet(connectNewChainLogic, connectNewChainPrecondition)

// ConnectionConfig defines how a chain should connect with other chains
type ConnectionConfig struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages from this chain.
	RMNVerificationDisabled bool
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages to this chain.
	AllowListEnabled bool
}

// ConnectNewChainConfig is a configuration struct for ConnectNewChainChangeset.
type ConnectNewChainConfig struct {
	// NewChainSelector is the selector of the new chain to connect.
	NewChainSelector uint64
	// NewChainConnectionConfig defines how the new chain should connect with other chains.
	NewChainConnectionConfig ConnectionConfig
	// RemoteChains are the chains to connect the new chain to.
	RemoteChains map[uint64]ConnectionConfig
	// TestRouter is true if we want to connect via test routers.
	TestRouter *bool
	// MCMSConfig is the MCMS configuration.
	MCMSConfig *changeset.MCMSConfig
}

// ValidateNewChain validates the new chain.
func (c ConnectNewChainConfig) ValidateNewChain(env deployment.Environment, state changeset.CCIPOnChainState) error {
	err := deployment.IsValidChainSelector(c.NewChainSelector)
	if err != nil {
		return fmt.Errorf("chain selector is invalid: %w", err)
	}

	chainState, ok := state.Chains[c.NewChainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d not found", c.NewChainSelector)
	}

	err = c.validateChain(env.GetContext(), chainState, env.Chains[c.NewChainSelector].DeployerKey.From, false)
	if err != nil {
		return fmt.Errorf("failed to validate chain with selector %d: %w", c.NewChainSelector, err)
	}

	return nil
}

// ValidateRemoteChains validates the remote chains.
func (c ConnectNewChainConfig) ValidateRemoteChains(env deployment.Environment, state changeset.CCIPOnChainState) error {
	for remoteChainSelector := range c.RemoteChains {
		err := deployment.IsValidChainSelector(remoteChainSelector)
		if err != nil {
			return fmt.Errorf("chain selector is invalid: %w", err)
		}

		chainState, ok := state.Chains[remoteChainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d not found", remoteChainSelector)
		}

		err = c.validateChain(env.GetContext(), chainState, env.Chains[remoteChainSelector].DeployerKey.From, true)
		if err != nil {
			return fmt.Errorf("failed to validate chain with selector %d: %w", remoteChainSelector, err)
		}
	}

	return nil
}

func (c ConnectNewChainConfig) validateChain(ctx context.Context, state changeset.CCIPChainState, deployerKey common.Address, ownedByMCMS bool) error {
	if state.OnRamp == nil {
		return errors.New("onRamp contract not found")
	}
	if state.OffRamp == nil {
		return errors.New("offRamp contract not found")
	}
	if state.Router == nil {
		return errors.New("router contract not found")
	}
	if state.TestRouter == nil {
		return errors.New("test router contract not found")
	}
	if state.ProposerMcm == nil {
		return errors.New("proposerMcm contract not found")
	}
	if state.Timelock == nil {
		return errors.New("timelock contract not found")
	}

	err := commoncs.ValidateOwnership(ctx, ownedByMCMS, deployerKey, state.Timelock.Address(), state.OnRamp)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of onRamp: %w", err)
	}
	err = commoncs.ValidateOwnership(ctx, ownedByMCMS, deployerKey, state.Timelock.Address(), state.OffRamp)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of offRamp: %w", err)
	}
	err = commoncs.ValidateOwnership(ctx, ownedByMCMS, deployerKey, state.Timelock.Address(), state.Router)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of router: %w", err)
	}
	// Test router should always be owned by deployer key
	err = commoncs.ValidateOwnership(ctx, false, deployerKey, state.Timelock.Address(), state.TestRouter)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of test router: %w", err)
	}

	return nil
}

func connectNewChainPrecondition(env deployment.Environment, c ConnectNewChainConfig) error {
	if c.MCMSConfig == nil {
		return fmt.Errorf("mcms config is required")
	}

	if c.TestRouter == nil {
		return fmt.Errorf("must define whether to use the test router")
	}

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = c.ValidateNewChain(env, state)
	if err != nil {
		return fmt.Errorf("failed to validate new chain: %w", err)
	}

	err = c.ValidateRemoteChains(env, state)
	if err != nil {
		return fmt.Errorf("failed to validate remote chains: %w", err)
	}

	return nil
}

func connectNewChainLogic(env deployment.Environment, c ConnectNewChainConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	readOpts := &bind.CallOpts{Context: env.GetContext()}

	var ownershipTransferProposals []timelock.MCMSWithTimelockProposal
	if !*c.TestRouter {
		// If using the production router, transfer ownership of all contracts on the new chain to MCMS.
		allContracts := []commoncs.Ownable{
			state.Chains[c.NewChainSelector].OnRamp,
			state.Chains[c.NewChainSelector].OffRamp,
			state.Chains[c.NewChainSelector].FeeQuoter,
			state.Chains[c.NewChainSelector].RMNProxy,
			state.Chains[c.NewChainSelector].NonceManager,
			state.Chains[c.NewChainSelector].TokenAdminRegistry,
			state.Chains[c.NewChainSelector].Router,
			state.Chains[c.NewChainSelector].RMNRemote,
		}
		addressesToTransfer := make([]common.Address, 0, len(allContracts))
		for _, contract := range allContracts {
			if contract == nil {
				continue
			}
			owner, err := contract.Owner(readOpts)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get owner of contract %s: %w", contract.Address().Hex(), err)
			}
			if owner == env.Chains[c.NewChainSelector].DeployerKey.From {
				addressesToTransfer = append(addressesToTransfer, contract.Address())
			}
		}
		out, err := commoncs.TransferToMCMSWithTimelock(env, commoncs.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				c.NewChainSelector: addressesToTransfer,
			},
			MinDelay: c.MCMSConfig.MinDelay,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run TransferToMCMSWithTimelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		ownershipTransferProposals = out.Proposals

		// Also, renounce the admin role on the Timelock (if not already done).
		adminRole, err := state.Chains[c.NewChainSelector].Timelock.ADMINROLE(readOpts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get admin role of timelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		hasRole, err := state.Chains[c.NewChainSelector].Timelock.HasRole(readOpts, adminRole, env.Chains[c.NewChainSelector].DeployerKey.From)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if deployer key has admin role on timelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		if hasRole {
			out, err = commoncs.RenounceTimelockDeployer(env, commoncs.RenounceTimelockDeployerConfig{
				ChainSel: c.NewChainSelector,
			})
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to run RenounceTimelockDeployer on chain with selector %d: %w", c.NewChainSelector, err)
			}
		}
	}

	// Enable the production router on [new chain -> each remote chain] and [each remote chain -> new chain].
	var allEnablementProposals []mcmslib.TimelockProposal
	var mcmsConfig *changeset.MCMSConfig
	if !*c.TestRouter {
		mcmsConfig = c.MCMSConfig
	}
	allEnablementProposals, err = connectRampsAndRouters(env, c.NewChainSelector, c.RemoteChains, mcmsConfig, *c.TestRouter, allEnablementProposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", c.NewChainSelector, err)
	}
	for remoteChainSelector := range c.RemoteChains {
		allEnablementProposals, err = connectRampsAndRouters(env, remoteChainSelector, map[uint64]ConnectionConfig{c.NewChainSelector: c.NewChainConnectionConfig}, c.MCMSConfig, *c.TestRouter, allEnablementProposals)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", remoteChainSelector, err)
		}
	}

	// Aggregate all proposals.
	// First, add ownership transfer proposals to the aggregate.
	var batches []types.BatchOperation
	for _, proposal := range ownershipTransferProposals {
		for _, batchTransaction := range proposal.Transactions {
			for _, transaction := range batchTransaction.Batch {
				batchOperation, err := proposalutils.BatchOperationForChain(
					uint64(batchTransaction.ChainIdentifier),
					transaction.To.Hex(),
					transaction.Data,
					big.NewInt(0),
					transaction.ContractType,
					transaction.Tags,
				)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to create batch operation on chain with selector %d: %w", batchTransaction.ChainIdentifier, err)
				}
				batches = append(batches, batchOperation)
			}
		}
	}
	// Add proposals that follow the new format to the aggregate.
	for _, proposal := range allEnablementProposals {
		batches = append(batches, proposal.Operations...)
	}

	// Store the timelocks, proposers, and inspectors for each chain.
	timelocks := make(map[uint64]string)
	proposers := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	for _, op := range batches {
		chainSel := uint64(op.ChainSelector)
		timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
		proposers[chainSel] = state.Chains[chainSel].ProposerMcm.Address().Hex()
		inspectors[chainSel], err = proposalutils.McmsInspectorForChain(env, chainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get MCMS inspector for chain with selector %d: %w", chainSel, err)
		}
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		timelocks,
		proposers,
		inspectors,
		batches,
		fmt.Sprintf("Connect chain with selector %d to other chains", c.NewChainSelector),
		c.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

// connectRampsAndRouters updates the onRamp and offRamp to point at the router for the given remote chains.
// It also sets the onRamp and offRamp on the router for the given remote chains.
// This function will add the proposals required to make these changes to the proposalAggregate slice.
func connectRampsAndRouters(
	e deployment.Environment,
	chainSelector uint64,
	remoteChains map[uint64]ConnectionConfig,
	mcmsConfig *changeset.MCMSConfig,
	testRouter bool,
	proposalAggregate []mcmslib.TimelockProposal,
) ([]mcmslib.TimelockProposal, error) {
	// Update offRamp sources on the new chain.
	offRampUpdatesOnNew := make(map[uint64]OffRampSourceUpdate, len(remoteChains))
	for remoteChainSelector, remoteChain := range remoteChains {
		offRampUpdatesOnNew[remoteChainSelector] = OffRampSourceUpdate{
			TestRouter:                testRouter,
			IsRMNVerificationDisabled: remoteChain.RMNVerificationDisabled,
			IsEnabled:                 true,
		}
	}
	out, err := UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
		UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
			chainSelector: offRampUpdatesOnNew,
		},
		MCMS:               mcmsConfig,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update onRamp destinations on the new chain.
	onRampUpdatesOnNew := make(map[uint64]OnRampDestinationUpdate, len(remoteChains))
	for remoteChainSelector, remoteChain := range remoteChains {
		onRampUpdatesOnNew[remoteChainSelector] = OnRampDestinationUpdate{
			TestRouter:       testRouter,
			AllowListEnabled: remoteChain.AllowListEnabled,
			IsEnabled:        true,
		}
	}
	out, err = UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
			chainSelector: onRampUpdatesOnNew,
		},
		MCMS:               mcmsConfig,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update router ramps on the new chain.
	offRampUpdates := make(map[uint64]bool, len(remoteChains))
	onRampUpdates := make(map[uint64]bool, len(remoteChains))
	for remoteChainSelector := range remoteChains {
		offRampUpdates[remoteChainSelector] = true
		onRampUpdates[remoteChainSelector] = true
	}
	cfg := mcmsConfig
	if testRouter {
		cfg = nil
	}
	out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
		TestRouter: testRouter,
		UpdatesByChain: map[uint64]RouterUpdates{
			chainSelector: RouterUpdates{
				OnRampUpdates:  onRampUpdates,
				OffRampUpdates: offRampUpdates,
			},
		},
		MCMS:               cfg,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateRouterRampsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	return proposalAggregate, nil
}
