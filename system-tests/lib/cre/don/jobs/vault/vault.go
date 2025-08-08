package vault

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/tdh2/go/tdh2/lib/group/nist"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/vault/sanmarinodkg/dummydkg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/vault/sanmarinodkg/tdh2shim"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var VaultJobSpecFactoryFn = func(chainID uint64) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.CldEnvironment.Offchain,
			input.DonTopology,
			input.CldEnvironment.DataStore,
			chainID,
		)
	}
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
	return pks, shares, nil
}

func GenerateJobSpecs(offchainClient cldf_offchain.Client, donTopology *cre.DonTopology, ds datastore.DataStore, chainID uint64) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	donMetadata := make([]*cre.DonMetadata, 0)
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

		// look for boostrap node and then for required values in its labels
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

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		// create job specs for the bootstrap node
		donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "vault-capability", vaultCapabilityAddress.Address, chainID))

		pk, sks, err := dkgKeys(len(workflowNodeSet), 1)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate DKG keys")
		}

		for idx, workerNode := range workflowNodeSet {
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

			encryptedShare, encErr := encryptPrivateShare(offchainClient, nodeID, sks[idx])
			if err != nil {
				return nil, errors.Wrap(encErr, "failed to encrypt private share")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerVaultOCR3(nodeID, vaultCapabilityAddress.Address, nodeEthAddr, ocr2KeyBundleID, donTopology.OCRPeeringData, chainID, pk, encryptedShare))
		}
	}

	return donToJobSpecs, nil
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
