package vault

import (
	"encoding/hex"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
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

		// return early if no DON has the vault capability
		if !input.DonTopology.AnyDonHasCapability(flag) {
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

		for _, don := range input.DonTopology.Dons.List() {
			if !don.HasFlag(flag) {
				continue
			}

			// create job specs for the worker nodes
			workerNodes, wErr := don.Workers()
			if wErr != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			bootstrapNode, isBootstrap := input.DonTopology.Bootstrap()
			if !isBootstrap {
				return nil, errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
			}

			_, ocrPeeringData, peeringErr := cre.PeeringCfgs(bootstrapNode)
			if peeringErr != nil {
				return nil, errors.Wrap(peeringErr, "failed to find peering data")
			}

			// create job specs for the bootstrap node
			donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.BootstrapOCR3(bootstrapNode.JobDistributorDetails.NodeID, "vault-capability", vaultCapabilityAddress.Address, chainID))

			for _, workerNode := range workerNodes {
				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}

				evmOCR2KeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM]
				if !ok {
					return nil, errors.New("key bundle ID for evm family is not found")
				}

				donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.WorkerVaultOCR3(workerNode.JobDistributorDetails.NodeID, vaultCapabilityAddress.Address, dkgAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringData, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}

func handlerConfig(don *cre.DON) (cre.HandlerTypeToConfig, error) {
	if !don.HasFlag(flag) {
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
