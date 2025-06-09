package consensus

import (
	"github.com/pkg/errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var ConsensusJobSpecFactoryFn = func(chainID uint64) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			input.AddressBook,
			chainID,
		)
	}
}

func GenerateJobSpecs(donTopology *types.DonTopology, addressBook cldf.AddressBook, chainID uint64) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	oCR3CapabilityAddress, ocr3err := crecontracts.FindAddressesForChain(addressBook, donTopology.HomeChainSelector, keystone_changeset.OCR3Capability.String())
	if ocr3err != nil {
		return nil, errors.Wrap(ocr3err, "failed to get OCR3 capability address")
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// create job specs for the worker nodes
		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
		if err != nil {
			// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		if flags.HasFlag(donWithMetadata.Flags, types.OCR3Capability) {
			// look for boostrap node and then for required values in its labels
			bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.BootstrapNode}, node.EqualLabels)
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
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, oCR3CapabilityAddress, chainID))

			ocrPeeringData := types.OCRPeeringData{
				OCRBootstraperPeerID: donBootstrapNodePeerID,
				OCRBootstraperHost:   donBootstrapNodeHost,
				Port:                 5001,
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
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerOCR3(nodeID, oCR3CapabilityAddress, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainID))
			}
		}
	}

	return donToJobSpecs, nil
}
