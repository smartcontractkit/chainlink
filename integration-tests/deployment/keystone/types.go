package keystone

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	v1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var (
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	KeystoneForwarder    deployment.ContractType = "KeystoneForwarder"
	OCR3Capability       deployment.ContractType = "OCR3Capability"
)

type deployResponse struct {
	Address common.Address
	Tx      common.Hash // todo: chain agnostic
	Tv      deployment.TypeAndVersion
}

type deployRequest struct {
	Chain deployment.Chain
}

type DonNode struct {
	Don  string
	Node string // not unique across environments
}

type CapabilityHost struct {
	NodeID       string // globally unique
	Capabilities []capabilities_registry.CapabilitiesRegistryCapability
}

type Nop struct {
	capabilities_registry.CapabilitiesRegistryNodeOperator
	NodeIDs []string // nodes run by this operator
}

// ocr2Node is a subset of the node configuration that is needed to register a node
// with the capabilities registry. Signer and P2PKey are chain agnostic.
type ocr2Node struct {
	ID         string
	Signer     [32]byte // note that in capabilities registry we need a [32]byte, but in the forwarder we need a common.Address [20]byte
	P2PKey     p2pkey.PeerID
	IsBoostrap bool
	// useful when have to register the ocr3 contract config
	p2pKeyBundle   *v1.OCR2Config_P2PKeyBundle
	ocrKeyBundle   *v1.OCR2Config_OCRKeyBundle
	accountAddress string
}

func (o *ocr2Node) signerAddress() common.Address {
	return common.BytesToAddress(o.Signer[:])
}

func (o *ocr2Node) toNodeKeys() NodeKeys {
	return NodeKeys{
		EthAddress:            o.accountAddress,
		P2PPeerID:             o.p2pKeyBundle.PeerId,
		OCR2BundleID:          o.ocrKeyBundle.BundleId,
		OCR2OnchainPublicKey:  o.ocrKeyBundle.OnchainSigningAddress,
		OCR2OffchainPublicKey: o.ocrKeyBundle.OffchainPublicKey,
		OCR2ConfigPublicKey:   o.ocrKeyBundle.ConfigPublicKey,
	}
}

func newOcr2Node(id string, ccfg *v1.ChainConfig) (*ocr2Node, error) {
	if ccfg == nil {
		return nil, errors.New("nil ocr2config")
	}
	ocfg := ccfg.Ocr2Config
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(ocfg.P2PKeyBundle.PeerId); err != nil {
		return nil, fmt.Errorf("failed to unmarshal peer id %s: %w", ocfg.P2PKeyBundle.PeerId, err)
	}

	signer := ocfg.OcrKeyBundle.OnchainSigningAddress
	if len(signer) != 40 {
		return nil, fmt.Errorf("invalid onchain signing address %s", ocfg.OcrKeyBundle.OnchainSigningAddress)
	}
	signerB, err := hex.DecodeString(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer %s: %w", signer, err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return &ocr2Node{
		ID:             id,
		Signer:         sigb,
		P2PKey:         p,
		IsBoostrap:     ocfg.IsBootstrap,
		p2pKeyBundle:   ocfg.P2PKeyBundle,
		ocrKeyBundle:   ocfg.OcrKeyBundle,
		accountAddress: ccfg.AccountAddress,
	}, nil
}

func makeNodeKeysSlice(nodes []*ocr2Node) []NodeKeys {
	var out []NodeKeys
	for _, n := range nodes {
		out = append(out, n.toNodeKeys())
	}
	return out
}

// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convienent form to handle the CLO representation of the nop data
type DonCapabilities struct {
	Name         string
	Nops         []*models.NodeOperator               // each nop is a node operator and may have multiple nodes
	Capabilities []kcr.CapabilitiesRegistryCapability // every capability is hosted on each nop
}

// map the node id to the NOP
func (dc DonCapabilities) nodeIdToNop(cs uint64) (map[string]capabilities_registry.CapabilitiesRegistryNodeOperator, error) {
	cid, err := chainsel.ChainIdFromSelector(cs)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id from selector %d: %w", cs, err)
	}
	cidStr := strconv.FormatUint(cid, 10)
	out := make(map[string]capabilities_registry.CapabilitiesRegistryNodeOperator)
	for _, nop := range dc.Nops {
		for _, node := range nop.Nodes {
			found := false
			for _, chain := range node.ChainConfigs {
				if chain.Network.ChainID == cidStr {
					found = true
					out[node.ID] = capabilities_registry.CapabilitiesRegistryNodeOperator{
						Name:  nop.Name,
						Admin: adminAddr(chain.AdminAddress),
					}
				}
			}
			if !found {
				return nil, fmt.Errorf("node '%s' %s does not support chain %d", node.Name, node.ID, cid)
			}
		}
	}
	return out, nil
}

// helpers to maintain compatibility with the existing registration functions
// NodesToNops converts a list of DonCapabilities to a map of node id to NOP
func NodesToNops(dons []DonCapabilities, chainSel uint64) (map[string]capabilities_registry.CapabilitiesRegistryNodeOperator, error) {
	out := make(map[string]capabilities_registry.CapabilitiesRegistryNodeOperator)
	for _, don := range dons {
		nops, err := don.nodeIdToNop(chainSel)
		if err != nil {
			return nil, fmt.Errorf("failed to get registry NOPs for don %s: %w", don.Name, err)
		}
		for donName, nop := range nops {
			_, exists := out[donName]
			if exists {
				continue
			}
			out[donName] = nop
		}
	}
	return out, nil
}

// MapDonsToCaps converts a list of DonCapabilities to a map of don name to capabilities
func MapDonsToCaps(dons []DonCapabilities) map[string][]kcr.CapabilitiesRegistryCapability {
	out := make(map[string][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		out[don.Name] = don.Capabilities
	}
	return out
}

var emptyAddr = "0x0000000000000000000000000000000000000000"

// compute the admin address from the string. If the address is empty, replaces the 0s with fs
// contract registry disallows 0x0 as an admin address, but our test net nops use it
func adminAddr(addr string) common.Address {
	needsFixing := addr == emptyAddr
	addr = strings.TrimPrefix(addr, "0x")
	if needsFixing {
		addr = strings.ReplaceAll(addr, "0", "f")
	}
	return common.HexToAddress(strings.TrimPrefix(addr, "0x"))
}
