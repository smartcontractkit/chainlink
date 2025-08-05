package v1_5_1

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	hybrid_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/hybrid_with_external_minter_fast_transfer_token_pool"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	HybridTokenPoolUpdateGroupsChangeset = cldf.CreateChangeSet(hybridTokenPoolUpdateGroupsLogic, hybridTokenPoolUpdateGroupsPrecondition)
)

// Group represents the token pool group type
type Group uint8

const (
	LockAndRelease Group = 0
	BurnAndMint    Group = 1
)

// GroupUpdateConfig represents a group update for a specific remote chain
type GroupUpdateConfig struct {
	RemoteChainSelector uint64
	Group               Group
	RemoteChainSupply   *big.Int
}

func (g *GroupUpdateConfig) Validate(contract *hybrid_external.HybridWithExternalMinterFastTransferTokenPool) error {
	if err := cldf.IsValidChainSelector(g.RemoteChainSelector); err != nil {
		return fmt.Errorf("invalid remote chain selector %d: %w", g.RemoteChainSelector, err)
	}

	if g.Group != LockAndRelease && g.Group != BurnAndMint {
		return fmt.Errorf("invalid group %d, must be 0 (LOCK_AND_RELEASE) or 1 (BURN_AND_MINT)", g.Group)
	}

	if g.RemoteChainSupply == nil {
		g.RemoteChainSupply = big.NewInt(0)
	}

	if g.RemoteChainSupply.Sign() < 0 {
		return errors.New("remote chain supply cannot be negative")
	}

	// Check if the remote chain is supported
	supported, err := contract.IsSupportedChain(nil, g.RemoteChainSelector)
	if err != nil {
		return fmt.Errorf("failed to check if chain %d is supported: %w", g.RemoteChainSelector, err)
	}
	if !supported {
		return fmt.Errorf("remote chain %d is not supported by the token pool", g.RemoteChainSelector)
	}

	// Check current group to prevent no-op updates
	currentGroup, err := contract.GetGroup(nil, g.RemoteChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get current group for chain %d: %w", g.RemoteChainSelector, err)
	}
	if Group(currentGroup) == g.Group {
		return fmt.Errorf("remote chain %d is already in group %d", g.RemoteChainSelector, g.Group)
	}

	return nil
}

type HybridTokenPoolUpdateGroupsConfig struct {
	TokenSymbol     shared.TokenSymbol
	ContractType    cldf.ContractType
	ContractVersion semver.Version
	Updates         map[uint64][]GroupUpdateConfig // chain selector -> group updates
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func (c HybridTokenPoolUpdateGroupsConfig) Validate(env cldf.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}

	if c.ContractType != shared.HybridWithExternalMinterFastTransferTokenPool {
		return fmt.Errorf("unsupported contract type %s for hybrid token pool group updates", c.ContractType)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, groupUpdates := range c.Updates {
		err := cldf.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}

		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}

		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain.String())
		}

		if err := validateHybridTokenPoolExists(chainState, c.TokenSymbol, c.ContractType, c.ContractVersion, chain.String()); err != nil {
			return err
		}

		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}

		if len(groupUpdates) == 0 {
			return fmt.Errorf("no group updates specified for chain %d", chainSelector)
		}

		// Get the hybrid token pool contract for validation
		pool, err := getHybridTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get hybrid token pool contract for %s token on chain %d: %w", c.TokenSymbol, chainSelector, err)
		}

		for _, update := range groupUpdates {
			err := update.Validate(pool)
			if err != nil {
				return fmt.Errorf("failed to validate group update for chain selector %d: %w", chainSelector, err)
			}
		}
	}
	return nil
}

func validateHybridTokenPoolExists(chainState evm.CCIPChainState, tokenSymbol shared.TokenSymbol, contractType cldf.ContractType, contractVersion semver.Version, chainString string) error {
	if contractType != shared.HybridWithExternalMinterFastTransferTokenPool {
		return fmt.Errorf("unsupported contract type %s for hybrid token pools", contractType)
	}

	if _, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol]; !ok {
		return fmt.Errorf("token %s does not have a hybrid token pool on %s", tokenSymbol, chainString)
	}
	if _, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]; !ok {
		return fmt.Errorf("token %s does not have a hybrid token pool with version %s on %s", tokenSymbol, contractVersion.String(), chainString)
	}
	return nil
}

func getHybridTokenPoolContract(env cldf.Environment, tokenSymbol shared.TokenSymbol, contractType cldf.ContractType, contractVersion semver.Version, chainSelector uint64) (*hybrid_external.HybridWithExternalMinterFastTransferTokenPool, error) {
	if contractType != shared.HybridWithExternalMinterFastTransferTokenPool {
		return nil, fmt.Errorf("unsupported contract type %s for hybrid token pools", contractType)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain %d does not exist in state", chainSelector)
	}

	poolAddress, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]
	if !ok {
		return nil, fmt.Errorf("hybrid token pool for token %s version %s not found on chain %d", tokenSymbol, contractVersion, chainSelector)
	}

	chain := env.BlockChains.EVMChains()[chainSelector]
	return hybrid_external.NewHybridWithExternalMinterFastTransferTokenPool(poolAddress.Address(), chain.Client)
}

func hybridTokenPoolUpdateGroupsPrecondition(env cldf.Environment, c HybridTokenPoolUpdateGroupsConfig) error {
	return c.Validate(env)
}

func hybridTokenPoolUpdateGroupsLogic(env cldf.Environment, c HybridTokenPoolUpdateGroupsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid HybridTokenPoolUpdateGroupsConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Build the sequence input for multi-chain updates
	updatesByChain := make(map[uint64]opsutil.EVMCallInput[ccipops.UpdateGroupsInput])

	for chainSelector, groupUpdates := range c.Updates {
		pool, err := getHybridTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get hybrid token pool contract for %s token on chain %d: %w", c.TokenSymbol, chainSelector, err)
		}

		// Convert to operation input format
		var opGroupUpdates []ccipops.GroupUpdate
		for _, update := range groupUpdates {
			if update.RemoteChainSupply == nil {
				update.RemoteChainSupply = big.NewInt(0)
			}

			opGroupUpdates = append(opGroupUpdates, ccipops.GroupUpdate{
				RemoteChainSelector: update.RemoteChainSelector,
				Group:               uint8(update.Group),
				RemoteChainSupply:   update.RemoteChainSupply,
			})
		}

		updatesByChain[chainSelector] = opsutil.EVMCallInput[ccipops.UpdateGroupsInput]{
			Address:       pool.Address(),
			ChainSelector: chainSelector,
			CallInput: ccipops.UpdateGroupsInput{
				GroupUpdates: opGroupUpdates,
			},
			NoSend: c.MCMS != nil, // Use NoSend for MCMS proposals
		}
	}

	// Execute the sequence
	seqInput := ccipseq.HybridTokenPoolUpdateGroupsSequenceInput{
		ContractType:   c.ContractType,
		UpdatesByChain: updatesByChain,
	}

	seqReport, err := operations.ExecuteSequence(env.OperationsBundle, ccipseq.HybridTokenPoolUpdateGroupsSequence, env.BlockChains.EVMChains(), seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute hybrid token pool update groups sequence: %w", err)
	}

	return opsutil.AddEVMCallSequenceToCSOutput(
		env,
		cldf.ChangesetOutput{},
		seqReport,
		err,
		state.EVMMCMSStateByChain(),
		c.MCMS,
		fmt.Sprintf("Update %s hybrid token pool groups", c.TokenSymbol),
	)
}
