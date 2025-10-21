package evm

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
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

func ConfigureEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, ocr3DON *cre.Don) (*forwarder.Config, error) {
	forwarderCfg := forwarder.DonConfiguration{
		Name:    ocr3DON.Name,
		ID:      libc.MustSafeUint32FromUint64(ocr3DON.ID),
		F:       ocr3DON.F,
		Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point, can we get that version from forwarder contract?
		NodeIDs: ocr3DON.KeystoneDONConfig().NodeIDs,
	}

	chainsWithForwarders := make(map[uint64]struct{})
	for _, selector := range chainSelectors {
		chainsWithForwarders[selector] = struct{}{}
	}

	fout, err3 := operations.ExecuteSequence(
		cldfEnv.OperationsBundle,
		forwarder.ConfigureSeq,
		forwarder.ConfigureSeqDeps{
			Env: cldfEnv,
		},
		forwarder.ConfigureSeqInput{
			DON:    forwarderCfg,
			Chains: chainsWithForwarders,
		},
	)
	if err3 != nil {
		return nil, errors.Wrap(err3, "failed to configure forwarders")
	}

	return &fout.Output.Config, nil
}

func ChainsWithForwarders(blockchains []blockchains.Blockchain, nodeSets []cre.NodeSetWithCapabilityConfigs) map[string][]uint64 {
	chainsWithForwarders := make(map[string][]uint64)

	for _, bcOut := range blockchains {
		for _, nodeSet := range nodeSets {
			if chainSelectors, familyExists := chainsWithForwarders[bcOut.ChainFamily()]; familyExists {
				if slices.Contains(chainSelectors, bcOut.ChainSelector()) {
					continue
				}
			}

			if !bcOut.IsFamily(chainselectors.FamilyEVM) && !bcOut.IsFamily(chainselectors.FamilyTron) {
				continue
			}

			if flags.RequiresForwarderContract(nodeSet.GetCapabilityFlags(), bcOut.ChainID()) {
				if _, exists := chainsWithForwarders[bcOut.ChainFamily()]; !exists {
					chainsWithForwarders[bcOut.ChainFamily()] = []uint64{}
				}
				chainsWithForwarders[bcOut.ChainFamily()] = append(chainsWithForwarders[bcOut.ChainFamily()], bcOut.ChainSelector())
			}
		}
	}

	return chainsWithForwarders
}
