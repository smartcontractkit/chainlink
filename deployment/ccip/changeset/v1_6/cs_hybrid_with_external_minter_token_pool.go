package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/hybrid_with_external_minter_token_pool"

	cldfevm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var _ cldf.ChangeSet[ConfigureHybridWithExternalMinterTokenPoolConfig] = UpdateGroupsOnHybridWithExternalMinterTokenPool

// Group values as defined by HybridTokenPoolAbstract.Group.
const (
	groupLockAndRelease uint8 = 0
	groupBurnAndMint    uint8 = 1
)

type HybridGroupConfig struct {
	// Type is the type of the token pool.
	Type cldf.ContractType `json:"type"`

	// Version is the version of the token pool.
	Version semver.Version `json:"version"`

	Updates []hybrid_with_external_minter_token_pool.HybridTokenPoolAbstractGroupUpdate

	// RemoteSupplies is keyed by remote chain selector and is required for every update
	// carrying a non-zero RemoteChainSupply.
	RemoteSupplies map[uint64]shared.RemoteSupply `json:"remoteSupplies"`
}

type ConfigureHybridWithExternalMinterTokenPoolConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *cldfproposalutils.TimelockConfig

	// Symbol is the symbol of the token of interest.
	TokenSymbol shared.TokenSymbol

	// Address targets specific pool on-chain without looking it up based on the provided token.
	Address common.Address

	// Updates the group on hybrid token pools. Can only be called by the owner.
	GroupUpdates map[uint64]HybridGroupConfig
}

func (c ConfigureHybridWithExternalMinterTokenPoolConfig) Validate(env cldf.Environment) error {
	if c.Address == (common.Address{}) && c.TokenSymbol == "" {
		return errors.New("address or token symbol must be defined")
	}

	if c.Address != (common.Address{}) && c.TokenSymbol != "" {
		return errors.New("address and token symbol cannot both be defined")
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, poolUpdate := range c.GroupUpdates {
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

		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}

		pool, err := c.resolvePool(chainState, chain)
		if err != nil {
			return fmt.Errorf("failed to resolve pool on %s: %w", chain.String(), err)
		}

		if err := poolUpdate.Validate(env, pool, chain.String()); err != nil {
			return fmt.Errorf("invalid pool update on %s: %w", chain.String(), err)
		}
	}

	return nil
}

// resolvePool returns the pool addressed by c, looked up by token symbol in state or bound
// directly to the configured address.
func (c ConfigureHybridWithExternalMinterTokenPoolConfig) resolvePool(
	chainState evm.CCIPChainState,
	chain cldfevm.Chain,
) (*hybrid_with_external_minter_token_pool.HybridWithExternalMinterTokenPool, error) {
	if c.TokenSymbol != "" {
		pool, ok := chainState.HybridWithExternalMinterTokenPool[c.TokenSymbol][deployment.Version1_6_0]
		if !ok {
			return nil, fmt.Errorf("token pool does not exist with symbol %s", c.TokenSymbol)
		}
		return pool, nil
	}

	pool, err := hybrid_with_external_minter_token_pool.NewHybridWithExternalMinterTokenPool(c.Address, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create hybrid with external minter pool: %w", err)
	}
	return pool, nil
}

func (c HybridGroupConfig) Validate(
	env cldf.Environment,
	pool *hybrid_with_external_minter_token_pool.HybridWithExternalMinterTokenPool,
	chainName string,
) error {
	if _, ok := shared.TokenPoolTypes[c.Type]; !ok {
		return fmt.Errorf("%s is not a known token pool type", c.Type)
	}

	if _, ok := shared.TokenPoolVersions[c.Version]; !ok {
		return fmt.Errorf("%s is not a known token pool version", c.Version)
	}

	if c.Type != shared.HybridWithExternalMinterTokenPool {
		return fmt.Errorf("token pool type %s is not supported", c.Type)
	}

	if len(c.Updates) == 0 {
		return fmt.Errorf("no group updates specified for %s", chainName)
	}

	opts := &bind.CallOpts{Context: env.GetContext()}

	localDecimals, err := pool.GetTokenDecimals(opts)
	if err != nil {
		return fmt.Errorf("failed to fetch local token decimals from pool %s: %w", pool.Address(), err)
	}

	seen := make(map[uint64]struct{}, len(c.Updates))

	for _, update := range c.Updates {
		if err := cldf.IsValidChainSelector(update.RemoteChainSelector); err != nil {
			return fmt.Errorf("invalid remote chain selector %d: %w", update.RemoteChainSelector, err)
		}

		// updateGroups applies updates in order, so a duplicate would depend on the state
		// left by the first entry.
		if _, duplicate := seen[update.RemoteChainSelector]; duplicate {
			return fmt.Errorf("remote chain %d appears more than once in the group updates", update.RemoteChainSelector)
		}
		seen[update.RemoteChainSelector] = struct{}{}

		if update.Group != groupLockAndRelease && update.Group != groupBurnAndMint {
			return fmt.Errorf("invalid group %d for remote chain %d, must be %d (LOCK_AND_RELEASE) or %d (BURN_AND_MINT)",
				update.Group, update.RemoteChainSelector, groupLockAndRelease, groupBurnAndMint)
		}

		supported, err := pool.IsSupportedChain(opts, update.RemoteChainSelector)
		if err != nil {
			return fmt.Errorf("failed to check if chain %d is supported: %w", update.RemoteChainSelector, err)
		}

		if !supported {
			return fmt.Errorf("remote chain %d is not supported by the token pool", update.RemoteChainSelector)
		}

		// The contract reverts with InvalidGroupUpdate on a no-op transition.
		currentGroup, err := pool.GetGroup(opts, update.RemoteChainSelector)
		if err != nil {
			return fmt.Errorf("failed to read current group for chain %d: %w", update.RemoteChainSelector, err)
		}

		if currentGroup == update.Group {
			return fmt.Errorf("remote chain %d is already in group %d", update.RemoteChainSelector, update.Group)
		}

		// A group change that migrates no liquidity needs no remote supply.
		if update.RemoteChainSupply == nil || update.RemoteChainSupply.Sign() == 0 {
			continue
		}

		if update.RemoteChainSupply.Sign() < 0 {
			return fmt.Errorf("remote chain %d: RemoteChainSupply must not be negative", update.RemoteChainSelector)
		}

		remoteSupply, ok := c.RemoteSupplies[update.RemoteChainSelector]
		if !ok {
			return fmt.Errorf("remote chain %d: RemoteChainSupply is non-zero but no RemoteSupply was given, so the local amount cannot be derived", update.RemoteChainSelector)
		}

		remoteToken, err := pool.GetRemoteToken(opts, update.RemoteChainSelector)
		if err != nil {
			return fmt.Errorf("failed to read remote token for chain %d: %w", update.RemoteChainSelector, err)
		}

		verified, err := shared.VerifyRemoteDecimals(env.GetContext(), env, update.RemoteChainSelector, remoteToken, remoteSupply.Decimals)
		if err != nil {
			return err
		}

		if !verified {
			env.Logger.Warnw("Remote token decimals could not be verified on chain, relying on the supplied value",
				"remoteChainSelector", update.RemoteChainSelector,
				"declaredDecimals", remoteSupply.Decimals,
			)
		}

		if err := shared.ValidateRemoteChainSupply(remoteSupply, update.RemoteChainSupply, localDecimals, update.RemoteChainSelector); err != nil {
			return err
		}
	}

	return nil
}

// UpdateGroupsOnHybridWithExternalMinterTokenPool updates the groups on hybrid with external minter token pools for a given token across multiple chains.
func UpdateGroupsOnHybridWithExternalMinterTokenPool(env cldf.Environment, c ConfigureHybridWithExternalMinterTokenPoolConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid ConfigureTokenPoolContractsConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext(fmt.Sprintf("configure %s token pool groups", c.TokenSymbol))

	for chainSelector, tokenPool := range c.GroupUpdates {
		if tokenPool.Type != shared.HybridWithExternalMinterTokenPool {
			return cldf.ChangesetOutput{}, fmt.Errorf("token pool type %s is not supported", tokenPool.Type)
		}

		chain := env.BlockChains.EVMChains()[chainSelector]
		chainState, _ := state.EVMChainState(chainSelector)

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		pool, err := c.resolvePool(chainState, chain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to resolve pool on %s: %w", chain, err)
		}

		if _, err := pool.UpdateGroups(opts, tokenPool.Updates); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to update groups on token pool with address %s on %s: %w", pool.Address().String(), chain, err)
		}
	}

	return deployerGroup.Enact()
}
