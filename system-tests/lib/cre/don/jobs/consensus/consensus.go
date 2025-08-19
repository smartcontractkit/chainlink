package consensus

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

var ConsensusJobSpecFactoryFn = func(chainID uint64) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			input.CldEnvironment.DataStore,
			chainID,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, ds datastore.DataStore, chainID uint64) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	ocr3Key := datastore.NewAddressRefKey(
		donTopology.HomeChainSelector,
		datastore.ContractType(keystone_changeset.OCR3Capability.String()),
		semver.MustParse("1.0.0"),
		"capability_ocr3",
	)
	ocr3CapabilityAddress, err := ds.Addresses().Get(ocr3Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Vault capability address")
	}

	donTimeKey := datastore.NewAddressRefKey(
		donTopology.HomeChainSelector,
		datastore.ContractType(keystone_changeset.OCR3Capability.String()),
		semver.MustParse("1.0.0"),
		"DONTime",
	)
	donTimeAddress, err := ds.Addresses().Get(donTimeKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DON Time address")
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.OCR3Capability) {
			continue
		}

		// create job specs for the worker nodes
		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// look for boostrap node and then for required values in its labels
		bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if bootErr != nil {
			return nil, errors.Wrap(bootErr, "failed to find bootstrap node")
		}

		donBootstrapNodePeerID, pIDErr := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
		if pIDErr != nil {
			return nil, errors.Wrap(pIDErr, "failed to get bootstrap node peer ID")
		}

		donBootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
		if hostErr != nil {
			return nil, errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
		}

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		// create job specs for the bootstrap node
		donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "ocr3-capability", ocr3CapabilityAddress.Address, chainID))

		ocrPeeringData := cre.OCRPeeringData{
			OCRBootstraperPeerID: donBootstrapNodePeerID,
			OCRBootstraperHost:   donBootstrapNodeHost,
			Port:                 config.OCRPeeringPort,
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			nodeEthAddr, ethErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(donTopology.HomeChainSelector))
			if ethErr != nil {
				return nil, errors.Wrap(ethErr, "failed to get eth address from labels")
			}

			ocr2KeyBundleID, ocr2Err := node.FindLabelValue(workerNode, node.NodeOCR2KeyBundleIDKey)
			if ocr2Err != nil {
				return nil, errors.Wrap(ocr2Err, "failed to get ocr2 key bundle id from labels")
			}
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerOCR3(nodeID, ocr3CapabilityAddress.Address, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainID))
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.DonTimeJob(nodeID, donTimeAddress.Address, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainID))
		}
	}

	return donToJobSpecs, nil
}

var ConsensusV2JobSpecFactoryFn = func(logger zerolog.Logger,
	chainID uint64,
	config map[string]any,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Input,
	evmBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		configGen := func(nodeAddress string) (string, error) {
			return fmt.Sprintf(
				`'{"chainId":%d,"network":"evm","nodeAddress":"%s"}'`,
				chainID,
				nodeAddress,
			), nil
		}

		return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
			logger,
			input.DonTopology,
			input.CldEnvironment.DataStore,
			chainID,
			capabilitiesAwareNodeSets,
			infraInput,
			evmBinaryPath,
			"capability_consensus",
			cre.ConsensusCapability,
			fmt.Sprintf("consensus-capability-%d", chainID),
			configGen,
		)
	}
}
