package evm

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

func DeployEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[string]string) error {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(cldfEnv.DataStore)
	if mergeErr != nil {
		return fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	evmForwardersReport, deployErr := operations.ExecuteSequence(
		cldfEnv.OperationsBundle,
		forwarder.DeploySequence,
		forwarder.DeploySequenceDeps{
			Env: cldfEnv,
		},
		forwarder.DeploySequenceInput{
			Targets: chainSelectors,
		},
	)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy evm forwarder")
	}

	if err := cldfEnv.ExistingAddresses.Merge(evmForwardersReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return errors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	}

	if err := memoryDatastore.Merge(evmForwardersReport.Output.Datastore); err != nil {
		return errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	for _, selector := range chainSelectors {
		forwarderAddr := contracts.MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
		testLogger.Info().Msgf("Deployed EVM Forwarder %s contract on chain %d at %s", contractVersions[keystone_changeset.KeystoneForwarder.String()], selector, forwarderAddr)
	}

	cldfEnv.DataStore = memoryDatastore.Seal()

	return nil
}
