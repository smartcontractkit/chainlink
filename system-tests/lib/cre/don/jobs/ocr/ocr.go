package ocr

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func GenerateJobSpecsForStandardCapabilityWithOCR(
	donTopology *cre.DonTopology,
	ds datastore.DataStore,
	nodeSetInput []*cre.CapabilitiesAwareNodeSet,
	infraInput *infra.Input,
	flag cre.CapabilityFlag,
	contractNamer ContractNamer,
	capabilityEnabler CapabilityEnabler,
	enabledChainsProvider EnabledChainsProvider,
	jobConfigGenerator JobConfigGenerator,
	configMerger ConfigMerger,
	capabilitiesConfig map[string]cre.CapabilityConfig,
) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	if infraInput == nil {
		return nil, errors.New("infra input is nil")
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

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		if !capabilityEnabler(nodeSetInput[donIdx], flag) {
			continue
		}

		capabilityConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("evm config not found in capabilities config")
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		binaryPath := filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath))

		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// look for boostrap node and then for required values in its labels
		bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if bootErr != nil {
			// if there is no bootstrap node in this DON, we need to use the global bootstrap node
			found := false
			for _, don := range donTopology.DonsWithMetadata {
				for _, n := range don.NodesMetadata {
					p2pValue, p2pErr := node.FindLabelValue(n, node.NodeP2PIDKey)
					if p2pErr != nil {
						continue
					}

					if strings.Contains(p2pValue, donTopology.OCRPeeringData.OCRBootstraperPeerID) {
						bootstrapNode = n
						found = true
						break
					}
				}
			}

			if !found {
				return nil, errors.New("failed to find global OCR bootstrap node")
			}
		}

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		internalHostsBS, err := getBoostrapWorkflowNames(bootstrapNode, donWithMetadata, infraInput)
		if err != nil {
			return nil, fmt.Errorf("no bootstrap node found for DON %s", donWithMetadata.Name)
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

			ocr3Key := datastore.NewAddressRefKey(
				cs,
				datastore.ContractType(keystone_changeset.OCR3Capability.String()),
				semver.MustParse("1.0.0"),
				contractName,
			)

			ocr3ConfigContractAddress, err := ds.Addresses().Get(ocr3Key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get EVM capability address")
			}

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, contractName, ocr3ConfigContractAddress.Address, chainIDUint64))
			logger.Debug().Msgf("Found deployed '%s' OCR3 contract on chain %d at %s", contractName, chainIDUint64, ocr3ConfigContractAddress.Address)

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				transmitterAddress, tErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(chain.Selector))
				if tErr != nil {
					return nil, errors.Wrap(tErr, "failed to get transmitter address from bootstrap node labels")
				}

				keyBundle, kErr := node.FindLabelValue(workerNode, node.NodeOCR2KeyBundleIDKey)
				if kErr != nil {
					return nil, errors.Wrap(kErr, "failed to get key bundle id from worker node labels")
				}

				keyNodeAddress := node.AddressKeyFromSelector(chain.Selector)
				nodeAddress, nodeAddressErr := node.FindLabelValue(workerNode, keyNodeAddress)
				if nodeAddressErr != nil {
					return nil, errors.Wrap(nodeAddressErr, "failed to get node address from labels")
				}
				logger.Debug().Msgf("Deployed node on chain %d/%d at %s", chainIDUint64, chain.Selector, nodeAddress)

				bootstrapNodeP2pKeyID, pErr := node.FindLabelValue(bootstrapNode, node.NodeP2PIDKey)
				if pErr != nil {
					return nil, errors.Wrap(pErr, "failed to get p2p key id from bootstrap node labels")
				}
				// remove the prefix if it exists, to match the expected format
				bootstrapNodeP2pKeyID = strings.TrimPrefix(bootstrapNodeP2pKeyID, "p2p_")
				bootstrapPeers := make([]string, len(internalHostsBS))
				for i, workflowName := range internalHostsBS {
					bootstrapPeers[i] = fmt.Sprintf("%s@%s:5001", bootstrapNodeP2pKeyID, workflowName)
				}

				oracleFactoryConfigInstance := job.OracleFactoryConfig{
					Enabled:            true,
					ChainID:            chainIDStr,
					BootstrapPeers:     bootstrapPeers,
					OCRContractAddress: ocr3ConfigContractAddress.Address,
					OCRKeyBundleID:     keyBundle,
					TransmitterID:      transmitterAddress,
					OnchainSigning: job.OnchainSigningStrategy{
						StrategyName: "single-chain",
						Config:       map[string]string{"evm": keyBundle},
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

				logger.Debug().Msgf("Creating %s Capability job spec for chainID: %d, selector: %d, DON:%q, node:%q", flag, chainIDUint64, chain.Selector, donWithMetadata.Name, nodeID)

				jobConfig, cErr := jobConfigGenerator(logger, chainIDUint64, nodeAddress, mergedConfig)
				if cErr != nil {
					return nil, errors.Wrap(cErr, "failed to generate job config")
				}

				jobName := contractName
				if chainIDUint64 != 0 {
					jobName = jobName + "-" + strconv.FormatUint(chainIDUint64, 10)
				}

				jobSpec := jobs.WorkerStandardCapability(nodeID, jobName, binaryPath, jobConfig, oracleStr)

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}

func getBoostrapWorkflowNames(bootstrapNode *cre.NodeMetadata, donWithMetadata *cre.DonWithMetadata, infraInput *infra.Input) ([]string, error) {
	nodeIndexStr, nErr := node.FindLabelValue(bootstrapNode, node.IndexKey)
	if nErr != nil {
		return nil, errors.Wrap(nErr, "failed to find index label")
	}

	nodeIndex, nIErr := strconv.Atoi(nodeIndexStr)
	if nIErr != nil {
		return nil, errors.Wrap(nIErr, "failed to convert index label value to int")
	}

	internalHostBS := don.InternalHost(nodeIndex, cre.BootstrapNode, donWithMetadata.Name, *infraInput)
	return []string{internalHostBS}, nil
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
