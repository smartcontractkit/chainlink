package onchain

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

// contractsFullyDeployed reports whether both registry contracts this step is
// responsible for are already present in state. DeployV2RegistryContractsSequence
// deploys CapabilitiesRegistry and WorkflowRegistry together and can't deploy just
// one of them, so a partial deploy (only one present, e.g. from a crashed prior
// run) is not skippable — it's redeployed in full, see deployContracts.
func contractsFullyDeployed(state *domain.StateFile) bool {
	return state.HasAddress(keystone_changeset.CapabilitiesRegistry.String()) &&
		state.HasAddress(keystone_changeset.WorkflowRegistry.String())
}

// deployContracts deploys the registry contracts if not already fully present,
// and reports whether it actually deployed (vs. skipped). A fresh deploy means
// the new contract(s) start unconfigured on-chain regardless of any stored
// phase hash, so callers must treat a true return as forcing the downstream
// hash-gated phases to (re)run — see Deployer.Apply.
func (d *Deployer) deployContracts(env *cldf.Environment, chainSelector uint64, state *domain.StateFile) (bool, error) {
	d.log.Info().Msg("Deploying Keystone v2 contracts")

	if contractsFullyDeployed(state) {
		capRegAddr := state.GetAddress(keystone_changeset.CapabilitiesRegistry.String())
		wfRegAddr := state.GetAddress(keystone_changeset.WorkflowRegistry.String())
		d.log.Info().Str("capReg", capRegAddr).Str("wfReg", wfRegAddr).Msg("Contracts already deployed — skipping")

		memDs, err := hydrateMemoryDataStoreFromState(state, chainSelector)
		if err != nil {
			return false, errors.Wrap(err, "failed to hydrate datastore from state")
		}
		env.DataStore = memDs.Seal()
		return false, nil
	}

	report, err := operations.ExecuteSequence(
		env.OperationsBundle,
		ks_contracts_op.DeployV2RegistryContractsSequence,
		ks_contracts_op.DeployContractsSequenceDeps{Env: env},
		ks_contracts_op.DeployRegistryContractsSequenceInput{RegistryChainSelector: chainSelector},
	)
	if err != nil {
		return false, errors.Wrap(err, "failed to deploy contracts")
	}

	memDs := datastore.NewMemoryDataStore()
	if err := memDs.Merge(report.Output.Datastore); err != nil {
		return false, errors.Wrap(err, "failed to merge deployed contract addresses")
	}

	capRegAddr := crecontracts.MustGetAddressFromMemoryDataStore(
		memDs, chainSelector, keystone_changeset.CapabilitiesRegistry.String(), crecontracts.V2Version, "",
	)
	wfRegAddr := crecontracts.MustGetAddressFromMemoryDataStore(
		memDs, chainSelector, keystone_changeset.WorkflowRegistry.String(), crecontracts.V2Version, "",
	)

	env.DataStore = memDs.Seal()

	d.log.Info().Str("capReg", capRegAddr.Hex()).Str("wfReg", wfRegAddr.Hex()).Msg("Contracts deployed")
	return true, nil
}
