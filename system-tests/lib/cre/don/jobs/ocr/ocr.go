package ocr

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	ptypes "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func GenerateJobSpecsForStandardCapabilityWithOCR(
	donTopology *cre.DonTopology,
	ds datastore.DataStore,
	nodeSetInput []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Provider,
	flag cre.CapabilityFlag,
	contractNamer ContractNamer,
	dataStoreOCR3ContractKeyProvider DataStoreOCR3ContractKeyProvider,
	capabilityEnabler CapabilityEnabler,
	enabledChainsProvider EnabledChainsProvider,
	jobConfigGenerator JobConfigGenerator,
	configMerger ConfigMerger,
	capabilitiesConfig map[string]cre.CapabilityConfig,
) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	if configMerger == nil {
		return nil, errors.New("config merger is nil")
	}
	if jobConfigGenerator == nil {
		return nil, errors.New("job config generator is nil")
	}
	if contractNamer == nil {
		return nil, errors.New("contract namer is nil")
	}
	if capabilityEnabler == nil {
		return nil, errors.New("capability enabler is nil")
	}
	if enabledChainsProvider == nil {
		return nil, errors.New("enabled chains provider is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	logger := framework.L

	for donIdx, donMetadata := range donTopology.ToDonMetadata() {
		if !capabilityEnabler(nodeSetInput[donIdx], flag) {
			continue
		}

		capabilityConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, fmt.Errorf("%s config not found in capabilities config: %v", flag, capabilitiesConfig)
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		binaryPath := filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath))

		workerNodes, wErr := donMetadata.WorkerNodes()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		bootstrapNode, err := donTopology.BootstrapNode()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get bootstrap node from DON metadata")
		}

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		chainIDs, err := enabledChainsProvider(donTopology, nodeSetInput[donIdx], flag)
		if err != nil {
			return nil, fmt.Errorf("failed to get enabled chains %w", err)
		}

		for _, chainIDUint64 := range chainIDs {
			chainIDStr := strconv.FormatUint(chainIDUint64, 10)
			chain, ok := chainsel.ChainByEvmChainID(chainIDUint64)
			if !ok {
				return nil, fmt.Errorf("failed to get chain selector for chain ID %d", chainIDUint64)
			}

			mergedConfig, enabled, rErr := configMerger(flag, nodeSetInput[donIdx], chainIDUint64, capabilityConfig)
			if rErr != nil {
				return nil, errors.Wrap(rErr, "failed to merge capability config")
			}

			// if the capability is not enabled for this chain, skip
			if !enabled {
				continue
			}

			cs, ok := chainsel.EvmChainIdToChainSelector()[chainIDUint64]
			if !ok {
				return nil, fmt.Errorf("chain selector not found for chainID: %d", chainIDUint64)
			}

			contractName := contractNamer(chainIDUint64)
			ocr3Key := dataStoreOCR3ContractKeyProvider(contractName, cs)
			ocr3ConfigContractAddress, err := ds.Addresses().Get(ocr3Key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get EVM capability address")
			}

			if _, ok := donToJobSpecs[donMetadata.ID]; !ok {
				donToJobSpecs[donMetadata.ID] = make(cre.DonJobs, 0)
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, contractName, ocr3ConfigContractAddress.Address, chainIDUint64))
			logger.Debug().Msgf("Found deployed '%s' OCR3 contract on chain %d at %s", contractName, chainIDUint64, ocr3ConfigContractAddress.Address)

			for _, workerNode := range workerNodes {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				ethKey, ok := workerNode.Keys.EVM[chainIDUint64]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainIDUint64, workerNode.Index)
				}
				transmitterAddress := ethKey.PublicAddress.Hex()

				bundlesPerFamily, kbErr := node.ExtractBundleKeysPerFamily(workerNode)
				if kbErr != nil {
					return nil, errors.Wrap(kbErr, "failed to get ocr families bundle id from worker node labels")
				}

				keyBundle, ok := bundlesPerFamily[chainsel.FamilyEVM] // we can always expect evm bundle key id present since evm is the registry chain
				if !ok {
					return nil, errors.New("failed to get key bundle id for evm family")
				}

				nodeAddress := transmitterAddress
				logger.Debug().Msgf("Deployed node on chain %d/%d at %s", chainIDUint64, chain.Selector, nodeAddress)

				bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", bootstrapNode.Keys.CleansedPeerID(), bootstrapNode.Host, cre.OCRPeeringPort)}

				strategyName := "single-chain"
				if len(bundlesPerFamily) > 1 {
					strategyName = "multi-chain"
				}
				oracleFactoryConfigInstance := job.OracleFactoryConfig{
					Enabled:            true,
					ChainID:            chainIDStr,
					BootstrapPeers:     bootstrapPeers,
					OCRContractAddress: ocr3ConfigContractAddress.Address,
					OCRKeyBundleID:     keyBundle,
					TransmitterID:      transmitterAddress,
					OnchainSigning: job.OnchainSigningStrategy{
						StrategyName: strategyName,
						Config:       bundlesPerFamily,
					},
				}

				// TODO: merge with jobConfig?
				type OracleFactoryConfigWrapper struct {
					OracleFactory job.OracleFactoryConfig `toml:"oracle_factory"`
				}
				wrapper := OracleFactoryConfigWrapper{OracleFactory: oracleFactoryConfigInstance}

				var oracleBuffer bytes.Buffer
				if errEncoder := toml.NewEncoder(&oracleBuffer).Encode(wrapper); errEncoder != nil {
					return nil, errors.Wrap(errEncoder, "failed to encode oracle factory config to TOML")
				}
				oracleStr := strings.ReplaceAll(oracleBuffer.String(), "\n", "\n\t")

				logger.Debug().Msgf("Creating %s Capability job spec for chainID: %d, selector: %d, DON:%q, node:%q", flag, chainIDUint64, chain.Selector, donMetadata.Name, nodeID)

				jobConfig, cErr := jobConfigGenerator(logger, chainIDUint64, nodeAddress, mergedConfig)
				if cErr != nil {
					return nil, errors.Wrap(cErr, "failed to generate job config")
				}

				jobName := contractName
				if chainIDUint64 != 0 {
					jobName = jobName + "-" + strconv.FormatUint(chainIDUint64, 10)
				}

				jobSpec := jobs.WorkerStandardCapability(nodeID, jobName, binaryPath, jobConfig, oracleStr)
				jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: &flag}}

				if _, ok := donToJobSpecs[donMetadata.ID]; !ok {
					donToJobSpecs[donMetadata.ID] = make(cre.DonJobs, 0)
				}

				donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}

// ConfigMerger merges default config with overrides (either on DON or chain level)
type ConfigMerger func(flag cre.CapabilityFlag, nodeSetInput *cre.CapabilitiesAwareNodeSet, chainIDUint64 uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error)

// JobConfigGenerator constains the logic that generates the job-specific part of the job spec
type JobConfigGenerator = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error)

// CapabilityEnabler determines if a capability is enabled for a given DON
type CapabilityEnabler func(nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool

// EnabledChainsProvider provides the list of enabled chains for a given capability
type EnabledChainsProvider func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) ([]uint64, error)

// ContractNamer is a function that returns the name of the OCR3 contract  used in the datastore
type ContractNamer func(chainID uint64) string

type DataStoreOCR3ContractKeyProvider func(contractName string, chainSelector uint64) datastore.AddressRefKey
