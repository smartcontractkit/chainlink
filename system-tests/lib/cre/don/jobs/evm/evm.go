package evm

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

var EVMJobSpecFactoryFn = func(logger zerolog.Logger, chainID uint64, networkFamily string, config map[string]any,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Input,
	evmBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(logger, input.DonTopology, input.CldEnvironment.DataStore, chainID, networkFamily, config, capabilitiesAwareNodeSets, infraInput, evmBinaryPath)
	}
}

var jobName = func(chainID uint64) string {
	return fmt.Sprintf("evm-capability-%d", chainID)
}

func GenerateJobSpecs(logger zerolog.Logger, donTopology *cre.DonTopology,
	ds datastore.DataStore,
	chainID uint64,
	networkFamily string, config map[string]any,
	nodeSetInput []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Input,
	evmBinaryPath string) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.EVMCapability) {
			continue
		}

		evmOCR3Key := datastore.NewAddressRefKey(
			donTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_evm",
		)
		evmOCR3CapabilityAddress, err := ds.Addresses().Get(evmOCR3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get EVM capability address")
		}

		internalHostsBS := getBoostrapWorkflowNames(donWithMetadata, nodeSetInput, donIdx, infraInput)
		if len(internalHostsBS) == 0 {
			return nil, fmt.Errorf("no bootstrap node found for DON %s (there should be at least 1)", donWithMetadata.Name)
		}

		pollIntervalStr, ok := config["logTriggerPollInterval"]
		if !ok || pollIntervalStr == nil {
			pollIntervalStr = "1s"
		}
		pollInterval, err := time.ParseDuration(pollIntervalStr.(string))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse poll interval '%s'", pollIntervalStr))
		}

		receiverGasMinimumStr, ok := config["receiverGasMinimum"]
		if !ok || receiverGasMinimumStr == nil {
			receiverGasMinimumStr = "1"
		}
		receiverGasMinimum, err := strconv.ParseUint(receiverGasMinimumStr.(string), 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid receiverGasMinimum '%s'", receiverGasMinimumStr))
		}

		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// look for boostrap node and then for required values in its labels
		bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if bootErr != nil {
			return nil, errors.Wrap(bootErr, "failed to find bootstrap node")
		}

		chain, ok := chainsel.ChainByEvmChainID(chainID)
		if !ok {
			return nil, fmt.Errorf("failed to get chain selector for chain ID %d", chainID)
		}

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		// create job specs for the bootstrap node
		donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "evm-capability", evmOCR3CapabilityAddress.Address, chainID))
		logger.Debug().Msgf("Deployed EVM OCR3 contract on chain %d at %s", chainID, evmOCR3CapabilityAddress.Address)

		creForwarderKey := datastore.NewAddressRefKey(
			donTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		)
		creForwarderAddress, err := ds.Addresses().Get(creForwarderKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get CRE Forwarder address")
		}

		logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

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
			logger.Debug().Msgf("Deployed node on chain %d/%d at %s", chainID, chain.Selector, nodeAddress)

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
				ChainID:            strconv.FormatUint(chainID, 10),
				BootstrapPeers:     bootstrapPeers,
				OCRContractAddress: evmOCR3CapabilityAddress.Address,
				OCRKeyBundleID:     keyBundle,
				TransmitterID:      transmitterAddress,
				OnchainSigning: job.OnchainSigningStrategy{
					StrategyName: "single-chain",
					Config:       map[string]string{"evm": keyBundle},
				},
			}

			type OracleFactoryConfigWrapper struct {
				OracleFactory job.OracleFactoryConfig `toml:"oracle_factory"`
			}
			wrapper := OracleFactoryConfigWrapper{OracleFactory: oracleFactoryConfigInstance}

			var oracleBuffer bytes.Buffer
			if errEncoder := toml.NewEncoder(&oracleBuffer).Encode(wrapper); errEncoder != nil {
				return nil, errors.Wrap(errEncoder, "failed to encode oracle factory config to TOML")
			}
			oracleStr := strings.ReplaceAll(oracleBuffer.String(), "\n", "\n\t")

			logger.Info().Msgf("Creating EVM Capability job spec for chainID: %d, selector: %d, DON:%q, node:%q", chainID, chain.Selector, donWithMetadata.Name, nodeID)

			jobSpec := jobs.WorkerStandardCapability(nodeID, jobName(chainID), evmBinaryPath,
				fmt.Sprintf(
					`'{"chainId":%d,"network":"%s","logTriggerPollInterval":%d, "creForwarderAddress":"%s","receiverGasMinimum":%d,"nodeAddress":"%s"}'`,
					chainID,
					networkFamily,
					pollInterval.Nanoseconds(),
					creForwarderAddress.Address,
					receiverGasMinimum,
					nodeAddress,
				),
				oracleStr,
			)

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}

func getBoostrapWorkflowNames(donWithMetadata *cre.DonWithMetadata, nodeSetInput []*cre.CapabilitiesAwareNodeSet, donIdx int, infraInput infra.Input) []string {
	internalHostsBS := make([]string, 0)
	for nodeIdx := range donWithMetadata.NodesMetadata {
		if nodeSetInput[donIdx].BootstrapNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].BootstrapNodeIndex {
			internalHostBS := don.InternalHost(nodeIdx, cre.BootstrapNode, donWithMetadata.Name, infraInput)
			internalHostsBS = append(internalHostsBS, internalHostBS)
		}
	}
	return internalHostsBS
}
