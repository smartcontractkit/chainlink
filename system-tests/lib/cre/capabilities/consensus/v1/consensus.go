package v1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
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
			contracts.OCR3ContractQualifier,
		)
		ocr3CapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(ocr3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
		}

		donTimeKey := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			contracts.DONTimeContractQualifier,
		)
		donTimeAddress, err := input.CldEnvironment.DataStore.Addresses().Get(donTimeKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get DON Time address")
		}

		for _, donMetadata := range input.DonTopology.ToDonMetadata() {
			if !flags.HasFlag(donMetadata.Flags, flag) {
				continue
			}

			workerNodes, wErr := donMetadata.WorkerNodes()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			bootstrapNode, err := input.DonTopology.BootstrapNode()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get bootstrap node from DON metadata")
			}

			_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrapNode)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get peering configs")
			}

			bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
			}

			donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "ocr3-capability", ocr3CapabilityAddress.Address, chainID))

			for _, workerNode := range workerNodes {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}

				ocr2KeyBundlesPerFamily, ocr2kbErr := node.ExtractBundleKeysPerFamily(workerNode)
				if ocr2kbErr != nil {
					return nil, errors.Wrap(ocr2kbErr, "failed to get ocr2 key bundle id from labels")
				}

				// we need the OCR2 key bundle for the EVM chain, because OCR jobs currently run only on EVM chains
				evmOCR2KeyBundle, ok := ocr2KeyBundlesPerFamily[chainselectors.FamilyEVM]
				if !ok {
					return nil, fmt.Errorf("node %s does not have OCR2 key bundle for evm", nodeID)
				}

				// we pass here bundles for all chains to enable multi-chain signing
				donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobs.WorkerOCR3(nodeID, ocr3CapabilityAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocr2KeyBundlesPerFamily, ocrPeeringCfg, chainID))
				donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobs.DonTimeJob(nodeID, donTimeAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringCfg, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}
