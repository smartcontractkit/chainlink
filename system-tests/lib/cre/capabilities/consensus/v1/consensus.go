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

		for _, don := range input.DonTopology.Dons.List() {
			if !don.HasFlag(flag) {
				continue
			}

			workerNodes, wErr := don.Workers()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			bootstrapNode, isBootstrap := input.DonTopology.Bootstrap()
			if !isBootstrap {
				return nil, errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
			}

			_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrapNode)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get peering configs")
			}

			donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.BootstrapOCR3(bootstrapNode.JobDistributorDetails.NodeID, "ocr3-capability", ocr3CapabilityAddress.Address, chainID))

			for _, workerNode := range workerNodes {
				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}

				// we need the OCR2 key bundle for the EVM chain, because OCR jobs currently run only on EVM chains
				evmOCR2KeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM]
				if !ok {
					return nil, fmt.Errorf("node %s does not have OCR2 key bundle for evm", workerNode.Name)
				}

				// we pass here bundles for all chains to enable multi-chain signing
				donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.WorkerOCR3(workerNode.JobDistributorDetails.NodeID, ocr3CapabilityAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, workerNode.Keys.OCR2BundleIDs, ocrPeeringCfg, chainID))
				donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.DonTimeJob(workerNode.JobDistributorDetails.NodeID, donTimeAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringCfg, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}
