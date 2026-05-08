package evm

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
)

func deployEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[cre.ContractType]*semver.Version) error {
	memoryDatastore, mErr := contracts.NewDataStoreFromExisting(cldfEnv.DataStore)
	if mErr != nil {
		return fmt.Errorf("failed to create memory datastore: %w", mErr)
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

func configureEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, ocr3DON *cre.Don) (*forwarder.Config, error) {
	forwarderCfg := forwarder.DonConfiguration{
		Name:    ocr3DON.Name,
		ID:      libc.MustSafeUint32FromUint64(ocr3DON.ID),
		F:       ocr3DON.F,
		Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point, can we get that version from forwarder contract?
		NodeIDs: ocr3DON.KeystoneDONConfig().NodeIDs,
	}

	if len(chainSelectors) == 0 {
		for _, chain := range cldfEnv.BlockChains.EVMChains() {
			chainSelectors = append(chainSelectors, chain.Selector)
		}
	}

	chainsByQualifier := make(map[string]map[uint64]struct{})
	for _, selector := range chainSelectors {
		refs := cldfEnv.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(selector),
			datastore.AddressRefByType(datastore.ContractType(keystone_changeset.KeystoneForwarder.String())),
		)
		if len(refs) == 0 {
			return nil, fmt.Errorf("failed to resolve deployed forwarder for chain selector %d", selector)
		}

		for _, ref := range refs {
			if chainsByQualifier[ref.Qualifier] == nil {
				chainsByQualifier[ref.Qualifier] = make(map[uint64]struct{})
			}
			chainsByQualifier[ref.Qualifier][selector] = struct{}{}
		}
	}

	qualifiers := make([]string, 0, len(chainsByQualifier))
	for qualifier := range chainsByQualifier {
		qualifiers = append(qualifiers, qualifier)
	}
	sort.Strings(qualifiers)

	var configuredConfig forwarder.Config
	for _, qualifier := range qualifiers {
		fout, err := operations.ExecuteSequence(
			cldfEnv.OperationsBundle,
			forwarder.ConfigureSeq,
			forwarder.ConfigureSeqDeps{
				Env: cldfEnv,
			},
			forwarder.ConfigureSeqInput{
				DON:       forwarderCfg,
				Qualifier: qualifier,
				Chains:    chainsByQualifier[qualifier],
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to configure forwarders with qualifier %q", qualifier)
		}
		configuredConfig = fout.Output.Config
	}

	return &configuredConfig, nil
}

func chainsWithForwarders(blockchains []blockchains.Blockchain, nodeSets []cre.NodeSetWithCapabilityConfigs) map[string][]uint64 {
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
