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

func (d *Deployer) deployContracts(env *cldf.Environment, chainSelector uint64, state *domain.StateFile) error {
	d.log.Info().Msg("P1: Deploying Keystone v2 contracts")

	if state.HasAddress(keystone_changeset.CapabilitiesRegistry.String()) {
		capRegAddr := state.GetAddress(keystone_changeset.CapabilitiesRegistry.String())
		wfRegAddr := state.GetAddress(keystone_changeset.WorkflowRegistry.String())
		d.log.Info().Str("capReg", capRegAddr).Str("wfReg", wfRegAddr).Msg("Contracts already deployed — skipping P1")

		memDs, err := hydrateMemoryDataStoreFromState(state, chainSelector)
		if err != nil {
			return errors.Wrap(err, "failed to hydrate datastore from state")
		}
		env.DataStore = memDs.Seal()
		return nil
	}

	report, err := operations.ExecuteSequence(
		env.OperationsBundle,
		ks_contracts_op.DeployV2RegistryContractsSequence,
		ks_contracts_op.DeployContractsSequenceDeps{Env: env},
		ks_contracts_op.DeployRegistryContractsSequenceInput{RegistryChainSelector: chainSelector},
	)
	if err != nil {
		return errors.Wrap(err, "failed to deploy contracts")
	}

	memDs := datastore.NewMemoryDataStore()
	if err := memDs.Merge(report.Output.Datastore); err != nil {
		return errors.Wrap(err, "failed to merge deployed contract addresses")
	}

	capRegAddr := crecontracts.MustGetAddressFromMemoryDataStore(
		memDs, chainSelector, keystone_changeset.CapabilitiesRegistry.String(), crecontracts.V2Version, "",
	)
	wfRegAddr := crecontracts.MustGetAddressFromMemoryDataStore(
		memDs, chainSelector, keystone_changeset.WorkflowRegistry.String(), crecontracts.V2Version, "",
	)

	env.DataStore = memDs.Seal()

	d.log.Info().Str("capReg", capRegAddr.Hex()).Str("wfReg", wfRegAddr.Hex()).Msg("Contracts deployed")
	return nil
}
