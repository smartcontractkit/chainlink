package vault

import (
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var VaultJobSpecFactoryFn = func(chainID uint64) cre.JobSpecFactoryFn {
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

	var donMetadata []*cre.DonMetadata
	for _, don := range donTopology.DonsWithMetadata {
		donMetadata = append(donMetadata, don.DonMetadata)
	}

	// return early if no DON has the vault capability
	if !don.AnyDonHasCapability(donMetadata, cre.VaultCapability) {
		return donToJobSpecs, nil
	}

	vaultOCR3Key := datastore.NewAddressRefKey(
		donTopology.HomeChainSelector,
		datastore.ContractType(keystone_changeset.OCR3Capability.String()),
		semver.MustParse("1.0.0"),
		"capability_vault",
	)
	vaultCapabilityAddress, err := ds.Addresses().Get(vaultOCR3Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Vault capability address")
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.VaultCapability) {
			continue
		}

		// create job specs for the worker nodes
		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// // look for boostrap node and then for required values in its labels
		bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if bootErr != nil {
			// if there is no bootstrap node in this DON, we need to use the global bootstrap node
			for _, don := range donTopology.DonsWithMetadata {
				for _, n := range don.NodesMetadata {
					p2pValue, p2pErr := node.FindLabelValue(n, node.NodeP2PIDKey)
					if p2pErr != nil {
						continue
					}

					if strings.Contains(p2pValue, donTopology.OCRPeeringData.OCRBootstraperPeerID) {
						bootstrapNode = n
						break
					}
				}
			}
		}

		// donBootstrapNodePeerID, pIDErr := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
		// if pIDErr != nil {
		// 	return nil, errors.Wrap(pIDErr, "failed to get bootstrap node peer ID")
		// }

		// donBootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
		// if hostErr != nil {
		// 	return nil, errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
		// }

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		// create job specs for the bootstrap node
		donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "vault-capability", vaultCapabilityAddress.Address, chainID))

		// ocrPeeringData := cre.OCRPeeringData{
		// 	OCRBootstraperPeerID: donBootstrapNodePeerID,
		// 	OCRBootstraperHost:   donBootstrapNodeHost,
		// 	Port:                 5001,
		// }

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
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerVaultOCR3(nodeID, vaultCapabilityAddress.Address, nodeEthAddr, ocr2KeyBundleID, donTopology.OCRPeeringData, chainID))
		}
	}

	return donToJobSpecs, nil
}
