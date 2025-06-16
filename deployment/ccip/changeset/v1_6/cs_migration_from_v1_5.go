package v1_6

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	// InitChainUpgratesOnTestRoutersChangeset sets candidates for the commit and exec DONs for multiple destination chains.
	// It then identifies all existing 1.5.0 source chains for each chain in the batch.
	// For each 1.5.0 OnRamp connecting to a destination, configuration gets translated to the 1.6.0 FeeQuoter.
	// In addition, OnRamps are connected to destination chains via test routers.
	// We do NOT connect the destinations back to the source chains, as DONs are not guaranteed to exist for sources.
	// This changeset is NOT IDEMPOTENT - if AddDON is called more than once for the same chain it will revert.
	InitChainUpgradesOnTestRoutersChangeset = cldf.CreateChangeSet(
		initChainUpgradesOnTestRoutersLogic,
		initChainUpgradesOnTestRoutersPrecondition,
	)
	// PromoteChainUpgradesToMainRoutersChangeset promotes the commit and exec DON candidates for multiple destination chains.
	// It then connects the source chains to the destination chains via main routers.
	// Before running PromoteChainUpgradesToMainRoutersChangeset for a batch, you must run InitChainUpgradesOnTestRoutersChangeset followed by SetOCR3OffRampChangeset.
	// SetOCR3OffRampChangeset should be run with ConfigType set to candidate, since the config won't be promoted until this changeset is run.
	// This changeset is NOT IDEMPOTENT - re-promoting will result in clearing the active digest, which is not desired.
	PromoteChainUpgradesToMainRoutersChangeset = cldf.CreateChangeSet(
		promoteChainUpgradesToMainRoutersLogic,
		promoteChainUpgradesToMainRoutersPrecondition,
	)
)

type InitChainUpgradesOnTestRoutersConfig struct{}

func initChainUpgradesOnTestRoutersPrecondition(e cldf.Environment, config InitChainUpgradesOnTestRoutersConfig) error {
	return nil
}

func initChainUpgradesOnTestRoutersLogic(e cldf.Environment, config InitChainUpgradesOnTestRoutersConfig) (cldf.ChangesetOutput, error) {
	return cldf.ChangesetOutput{}, nil
}

type PromoteChainUpgradesToMainRoutersConfig struct{}

func promoteChainUpgradesToMainRoutersPrecondition(e cldf.Environment, config PromoteChainUpgradesToMainRoutersConfig) error {
	return nil
}

func promoteChainUpgradesToMainRoutersLogic(e cldf.Environment, config PromoteChainUpgradesToMainRoutersConfig) (cldf.ChangesetOutput, error) {
	return cldf.ChangesetOutput{}, nil
}
