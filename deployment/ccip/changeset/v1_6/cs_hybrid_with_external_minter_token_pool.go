package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/latest/hybrid_with_external_minter_token_pool"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSet[ConfigureHybridWithExternalMinterTokenPoolConfig] = UpdateGroupsOnHybridWithExternalMinterTokenPool

type HybridGroupConfig struct {
	// Type is the type of the token pool.
	Type cldf.ContractType `json:"type"`

	// Version is the version of the token pool.
	Version semver.Version `json:"version"`

	Updates []hybrid_with_external_minter_token_pool.HybridTokenPoolAbstractGroupUpdate
}

type ConfigureHybridWithExternalMinterTokenPoolConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig

	// Symbol is the symbol of the token of interest.
	TokenSymbol shared.TokenSymbol

	// Updates the group on hybrid token pools. Can only be called by the owner.
	GroupUpdates map[uint64]HybridGroupConfig
}

func (c ConfigureHybridWithExternalMinterTokenPoolConfig) Validate(env cldf.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
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
		if err := poolUpdate.Validate(); err != nil {
			return fmt.Errorf("invalid pool update on %s: %w", chain.String(), err)
		}
	}

	return nil
}

func (c HybridGroupConfig) Validate() error {
	if _, ok := shared.TokenPoolTypes[c.Type]; !ok {
		return fmt.Errorf("%s is not a known token pool type", c.Type)
	}

	if _, ok := shared.TokenPoolVersions[c.Version]; !ok {
		return fmt.Errorf("%s is not a known token pool version", c.Version)
	}

	if c.Type != shared.HybridWithExternalMinterTokenPool {
		return fmt.Errorf("token pool type %s is not supported", c.Type)
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

		pool, ok := chainState.HybridWithExternalMinterTokenPool[c.TokenSymbol][deployment.Version1_6_0]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("token pool does not exist on %s with symbol %s", chain, c.TokenSymbol)
		}

		if _, err := pool.UpdateGroups(opts, tokenPool.Updates); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to update groups on token pool with address %s on %s: %w", pool.Address().String(), chain, err)
		}
	}

	return deployerGroup.Enact()
}
