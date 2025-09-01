package vault

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/tdh2/go/tdh2/lib/group/nist"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault/sanmarinodkg/dummydkg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault/sanmarinodkg/tdh2shim"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

const flag = cre.VaultCapability

// MasterPublicKeyStr is the public key used for the vault keys. It is set during environment setup
var MasterPublicKeyStr = "7b2247726f7570223a2250323536222c22475f626172223a22424c634f345a2b314e7571746f4b656f5264705672334d6e62547762637a31644e524b526472614a68744b426534534a46686a5049346e6e4243692b544c43556e64565968364961623661774b496e45527a4d7137614d3d222c2248223a2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c22484172726179223a5b2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d222c2242496c68454d4f43513569436c5437537976614d4a6e3250417269714a70396e672f52462f78623549664965655a43704633364a543933756b674c336641754468327637692f693035676a6451396776344c556357324d3d225d7d"

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

func EncryptSecret(secret string) (string, error) {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(MasterPublicKeyStr)
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
			"capability_vault",
		)
		vaultCapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(vaultOCR3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
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

			// Create Vault keys in a deterministic manner by setting n = 8
			// and t = 1, so that we can use the same keys across all tests,
			// irrespective of the number of nodes in the workflow DON.
			// We just need to ensure that the number of workflow nodes is at most 8, because dkgKeys()
			// will generate only 8 shares. If workflow nodes are increased more than 8,
			// we need to update the n value in dkgKeys() and the value of MasterPublicKeyStr accordingly.
			if len(workflowNodeSet) > 8 {
				return nil, errors.New("workflow node set must not exceed 8 nodes for vault capability, please update dkgKeys() and MasterPublicKeyStr accordingly if you want to increase the number of workflow nodes")
			}
			pk, sks, err := dkgKeys(8, 1)
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate DKG keys")
			}

			for idx, workerNode := range workflowNodeSet {
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

				encryptedShare, encErr := encryptPrivateShare(input.CldEnvironment.Offchain, nodeID, sks[idx])
				if err != nil {
					return nil, errors.Wrap(encErr, "failed to encrypt private share")
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerVaultOCR3(nodeID, vaultCapabilityAddress.Address, nodeEthAddr, ocr2KeyBundleID, input.DonTopology.OCRPeeringData, chainID, pk, encryptedShare))
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

func encryptPrivateShare(offchain cldf_offchain.Client, nodeID string, sk *tdh2easy.PrivateShare) (string, error) {
	nodeResp, err := offchain.GetNode(context.Background(), &nodev1.GetNodeRequest{
		Id: nodeID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get node from jd")
	}
	wk := nodeResp.GetNode().GetWorkflowKey()
	if wk == "" {
		return "", errors.New("node must contain a workflow key")
	}

	wkb, err := hex.DecodeString(wk)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode workflow key from hex")
	}

	skb, err := sk.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal private share")
	}

	wkbSized := [32]byte(wkb)
	sealed, err := box.SealAnonymous(nil, skb, &wkbSized, cryptorand.Reader)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt private share")
	}

	return hex.EncodeToString(sealed), nil
}

func dkgKeys(n, t int) (string, []*tdh2easy.PrivateShare, error) {
	instanceID, recipCfg, recipSecKeys, err := dummydkg.NewDKGSetup(n, t, "REPLACE_ME_WITH_RANDOM_SEED")
	if err != nil {
		return "", nil, err
	}

	group := nist.NewP256()
	result, err := dummydkg.NewDKGResult(instanceID, recipCfg, group)
	if err != nil {
		return "", nil, err
	}

	shares := []*tdh2easy.PrivateShare{}
	for _, share := range recipSecKeys {
		s, ierr := tdh2shim.TDH2PrivateShareFromDKGResult(result, share)
		if ierr != nil {
			return "", nil, errors.Wrap(ierr, "failed to convert DKG share to TDH2 share")
		}

		shares = append(shares, s)
	}

	pk, err := tdh2shim.TDH2PublicKeyFromDKGResult(result)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to convert DKG result to TDH2 public key")
	}

	pkb, err := pk.Marshal()
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to marshal TDH2 public key")
	}

	pks := hex.EncodeToString(pkb)
	framework.L.Info().Msg("Generated MasterPublicKey for n=" + strconv.Itoa(n) + ", t=" + strconv.Itoa(t) + ". Key = " + pks)
	if MasterPublicKeyStr != pks {
		framework.L.Warn().Msgf("Generated MasterPublicKeyStr is not same as the expected one: %s != %s", MasterPublicKeyStr, pks)
		return "", nil, errors.New("generated MasterPublicKeyStr does not match the expected one")
	}
	return pks, shares, nil
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
