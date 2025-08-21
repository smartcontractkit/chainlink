package v1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.ConsensusCapability

func New(chainID uint64) (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec(chainID)),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		capabilities.WithValidateFn(func(c *capabilities.Capability) error {
			if chainID == 0 {
				return fmt.Errorf("chainID is required, got %d", chainID)
			}
			return nil
		}),
	)
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}

func jobSpec(chainID uint64) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}
		donToJobSpecs := make(cre.DonsToJobSpecs)

		ocr3Key := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_ocr3",
		)
		ocr3CapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(ocr3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
		}

		donTimeKey := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"DONTime",
		)
		donTimeAddress, err := input.CldEnvironment.DataStore.Addresses().Get(donTimeKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get DON Time address")
		}

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			if !flags.HasFlag(donWithMetadata.Flags, flag) {
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
				Port:                 don.OCRPeeringPort,
			}

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				nodeEthAddr, ethErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(input.DonTopology.HomeChainSelector))
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
}
