package vault

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

const flag = cre.VaultCapability

func New(chainID uint64) (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec(chainID)),
		capabilities.WithGatewayJobHandlerConfigFn(handlerConfig),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		capabilities.WithValidateFn(func(c *capabilities.Capability) error {
			if chainID == 0 {
				return fmt.Errorf("chainID is required, got %d", chainID)
			}
			return nil
		}),
	)
}

func EncryptSecret(secret, masterPublicKeyStr string) (string, error) {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(masterPublicKeyStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode master public key")
	}
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal master public key")
	}
	cipher, err := tdh2easy.Encrypt(&masterPublicKey, []byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt secret")
	}
	cipherBytes, err := cipher.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal encrypted secrets to bytes")
	}
	return hex.EncodeToString(cipherBytes), nil
}

func jobSpec(chainID uint64) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}
		donToJobSpecs := make(cre.DonsToJobSpecs)

		donMetadata := make([]*cre.DonMetadata, 0)
		for _, don := range input.DonTopology.DonsWithMetadata {
			donMetadata = append(donMetadata, don.DonMetadata)
		}

		// return early if no DON has the vault capability
		if !don.AnyDonHasCapability(donMetadata, flag) {
			return donToJobSpecs, nil
		}

		vaultOCR3Key := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_vault_plugin",
		)
		vaultCapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(vaultOCR3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
		}

		dkgKey := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_vault_dkg",
		)
		dkgAddress, err := input.CldEnvironment.DataStore.Addresses().Get(dkgKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get DKG address")
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
				// if there is no bootstrap node in this DON, we need to use the global bootstrap node
				for _, don := range input.DonTopology.DonsWithMetadata {
					for _, n := range don.NodesMetadata {
						p2pValue, p2pErr := node.FindLabelValue(n, node.NodeP2PIDKey)
						if p2pErr != nil {
							continue
						}

						if strings.Contains(p2pValue, input.DonTopology.OCRPeeringData.OCRBootstraperPeerID) {
							bootstrapNode = n
							break
						}
					}
				}
			}

			bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "vault-capability", vaultCapabilityAddress.Address, chainID))

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				nodeEthAddr, ethErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(input.DonTopology.HomeChainSelector))
				if ethErr != nil {
					return nil, errors.Wrap(ethErr, "failed to get eth address from labels")
				}

				ocr2KeyBundlesPerFamily, ocr2kbErr := node.ExtractBundleKeysPerFamily(workerNode)
				if ocr2kbErr != nil {
					return nil, errors.Wrap(ocr2kbErr, "failed to get ocr2 key bundle id from labels")
				}
				offchainKeyBundleID, ok := ocr2KeyBundlesPerFamily["evm"]
				if !ok {
					return nil, errors.New("key bundle ID for evm family is not found")
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerVaultOCR3(nodeID, vaultCapabilityAddress.Address, dkgAddress.Address, nodeEthAddr, offchainKeyBundleID, input.DonTopology.OCRPeeringData, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}

func handlerConfig(donMetadata *cre.DonMetadata) (cre.HandlerTypeToConfig, error) {
	if !flags.HasFlag(donMetadata.Flags, flag) {
		return nil, nil
	}

	return map[string]string{coregateway.VaultHandlerType: `
ServiceName = "vault"
[gatewayConfig.Dons.Handlers.Config]
requestTimeoutSec = 70
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`}, nil
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "vault",
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
