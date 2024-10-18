package keystone

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"

	v1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
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
// TODO: KS-466 when we migrate fully to the JD offchain client, we should be able remove this shim and use environment.Node directly
type ocr2Node struct {
	ID                  string
	Signer              [32]byte // note that in capabilities registry we need a [32]byte, but in the forwarder we need a common.Address [20]byte
	P2PKey              p2pkey.PeerID
	EncryptionPublicKey [32]byte
	IsBoostrap          bool
	// useful when have to register the ocr3 contract config
	p2pKeyBundle       *v1.OCR2Config_P2PKeyBundle
	ethOcr2KeyBundle   *v1.OCR2Config_OCRKeyBundle
	aptosOcr2KeyBundle *v1.OCR2Config_OCRKeyBundle
	csaKey             string // *v1.Node.PublicKey
	accountAddress     string
}

func (o *ocr2Node) signerAddress() common.Address {
	return common.BytesToAddress(o.Signer[:])
}

func (o *ocr2Node) toNodeKeys() NodeKeys {
	return NodeKeys{
		EthAddress:            o.accountAddress,
		P2PPeerID:             o.p2pKeyBundle.PeerId,
		OCR2BundleID:          o.ethOcr2KeyBundle.BundleId,
		OCR2OnchainPublicKey:  o.ethOcr2KeyBundle.OnchainSigningAddress,
		OCR2OffchainPublicKey: o.ethOcr2KeyBundle.OffchainPublicKey,
		OCR2ConfigPublicKey:   o.ethOcr2KeyBundle.ConfigPublicKey,
		CSAPublicKey:          o.csaKey,
		// default value of encryption public key is the CSA public key
		// TODO: DEVSVCS-760
		EncryptionPublicKey: o.csaKey,
		// TODO Aptos support. How will that be modeled in clo data?
	}
}

func newOcr2Node(id string, ccfgs map[chaintype.ChainType]*v1.ChainConfig, csaPubKey string) (*ocr2Node, error) {
	if ccfgs == nil {
		return nil, errors.New("nil ocr2config")
	}
	evmCC, exists := ccfgs[chaintype.EVM]
	if !exists {
		return nil, errors.New("no evm chain config for node id " + id)
	}

	if csaPubKey == "" {
		return nil, errors.New("empty csa public key")
	}
	// parse csapublic key to
	csaKey, err := hex.DecodeString(csaPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode csa public key %s: %w", csaPubKey, err)
	}
	if len(csaKey) != 32 {
		return nil, fmt.Errorf("invalid csa public key '%s'. expected len 32 got %d", csaPubKey, len(csaKey))
	}
	var csaKeyb [32]byte
	copy(csaKeyb[:], csaKey)

	ocfg := evmCC.Ocr2Config
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

	n := &ocr2Node{
		ID:                  id,
		Signer:              sigb,
		P2PKey:              p,
		EncryptionPublicKey: csaKeyb,
		IsBoostrap:          ocfg.IsBootstrap,
		p2pKeyBundle:        ocfg.P2PKeyBundle,
		ethOcr2KeyBundle:    evmCC.Ocr2Config.OcrKeyBundle,
		aptosOcr2KeyBundle:  nil,
		accountAddress:      evmCC.AccountAddress,
		csaKey:              csaPubKey,
	}
	// aptos chain config is optional
	if aptosCC, exists := ccfgs[chaintype.Aptos]; exists {
		n.aptosOcr2KeyBundle = aptosCC.Ocr2Config.OcrKeyBundle
	}

	return n, nil
}

func makeNodeKeysSlice(nodes []*ocr2Node) []NodeKeys {
	var out []NodeKeys
	for _, n := range nodes {
		out = append(out, n.toNodeKeys())
	}
	return out
}

// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
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
// nodesToNops converts a list of DonCapabilities to a map of node id to NOP
func nodesToNops(dons []DonCapabilities, chainSel uint64) (map[string]capabilities_registry.CapabilitiesRegistryNodeOperator, error) {
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

// mapDonsToCaps converts a list of DonCapabilities to a map of don name to capabilities
func mapDonsToCaps(dons []DonCapabilities) map[string][]kcr.CapabilitiesRegistryCapability {
	out := make(map[string][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		out[don.Name] = don.Capabilities
	}
	return out
}

// mapDonsToNodes returns a map of don name to simplified representation of their nodes
// all nodes must have evm config and ocr3 capability nodes are must also have an aptos chain config
func mapDonsToNodes(dons []DonCapabilities, excludeBootstraps bool) (map[string][]*ocr2Node, error) {
	donToOcr2Nodes := make(map[string][]*ocr2Node)
	// get the nodes for each don from the offchain client, get ocr2 config from one of the chain configs for the node b/c
	// they are equivalent, and transform to ocr2node representation

	for _, don := range dons {
		for _, nop := range don.Nops {
			for _, node := range nop.Nodes {
				csaPubKey := node.PublicKey
				if csaPubKey == nil {
					return nil, fmt.Errorf("no public key for node %s", node.ID)
				}
				// the chain configs are equivalent as far as the ocr2 config is concerned so take the first one
				if len(node.ChainConfigs) == 0 {
					return nil, fmt.Errorf("no chain configs for node %s. cannot obtain keys", node.ID)
				}
				// all nodes should have an evm chain config, specifically the registry chain
				evmCC, exists := firstChainConfigByType(node.ChainConfigs, chaintype.EVM)
				if !exists {
					return nil, fmt.Errorf("no evm chain config for node %s", node.ID)
				}
				cfgs := map[chaintype.ChainType]*v1.ChainConfig{
					chaintype.EVM: evmCC,
				}
				aptosCC, exists := firstChainConfigByType(node.ChainConfigs, chaintype.Aptos)
				if exists {
					cfgs[chaintype.Aptos] = aptosCC
				}
				ocr2n, err := newOcr2Node(node.ID, cfgs, *csaPubKey)
				if err != nil {
					return nil, fmt.Errorf("failed to create ocr2 node for node %s: %w", node.ID, err)
				}
				if excludeBootstraps && ocr2n.IsBoostrap {
					continue
				}
				if _, ok := donToOcr2Nodes[don.Name]; !ok {
					donToOcr2Nodes[don.Name] = make([]*ocr2Node, 0)
				}
				donToOcr2Nodes[don.Name] = append(donToOcr2Nodes[don.Name], ocr2n)

			}
		}
	}

	return donToOcr2Nodes, nil
}

func firstChainConfigByType(ccfgs []*models.NodeChainConfig, t chaintype.ChainType) (*v1.ChainConfig, bool) {
	for _, c := range ccfgs {
		//nolint:staticcheck //ignore EqualFold it broke ci for some reason (go version skew btw local and ci?)
		if strings.ToLower(c.Network.ChainType.String()) == strings.ToLower(string(t)) {
			return chainConfigFromClo(c), true
		}
	}
	return nil, false
}

// RegisteredDon is a representation of a don that exists in the in the capabilities registry all with the enriched node data
type RegisteredDon struct {
	Name  string
	Info  capabilities_registry.CapabilitiesRegistryDONInfo
	Nodes []*ocr2Node
}

func (d RegisteredDon) signers() []common.Address {
	sort.Slice(d.Nodes, func(i, j int) bool {
		return d.Nodes[i].P2PKey.String() < d.Nodes[j].P2PKey.String()
	})
	var out []common.Address
	for _, n := range d.Nodes {
		if n.IsBoostrap {
			continue
		}
		out = append(out, n.signerAddress())
	}
	return out
}

func joinInfoAndNodes(donInfos map[string]kcr.CapabilitiesRegistryDONInfo, dons []DonCapabilities) ([]RegisteredDon, error) {
	// all maps should have the same keys
	nodes, err := mapDonsToNodes(dons, true)
	if err != nil {
		return nil, fmt.Errorf("failed to map dons to capabilities: %w", err)
	}
	if len(donInfos) != len(nodes) {
		return nil, fmt.Errorf("mismatched lengths don infos %d,  nodes %d", len(donInfos), len(nodes))
	}
	var out []RegisteredDon
	for donName, info := range donInfos {

		ocr2nodes, ok := nodes[donName]
		if !ok {
			return nil, fmt.Errorf("nodes not found for don %s", donName)
		}
		out = append(out, RegisteredDon{
			Name:  donName,
			Info:  info,
			Nodes: ocr2nodes,
		})
	}

	return out, nil
}

func chainConfigFromClo(chain *models.NodeChainConfig) *v1.ChainConfig {
	return &v1.ChainConfig{
		Chain: &v1.Chain{
			Id:   chain.Network.ChainID,
			Type: v1.ChainType_CHAIN_TYPE_EVM, // TODO: support other chain types
		},

		AccountAddress: chain.AccountAddress,
		AdminAddress:   chain.AdminAddress,
		Ocr2Config: &v1.OCR2Config{
			Enabled: chain.Ocr2Config.Enabled,
			P2PKeyBundle: &v1.OCR2Config_P2PKeyBundle{
				PeerId:    chain.Ocr2Config.P2pKeyBundle.PeerID,
				PublicKey: chain.Ocr2Config.P2pKeyBundle.PublicKey,
			},
			OcrKeyBundle: &v1.OCR2Config_OCRKeyBundle{
				BundleId:              chain.Ocr2Config.OcrKeyBundle.BundleID,
				OnchainSigningAddress: chain.Ocr2Config.OcrKeyBundle.OnchainSigningAddress,
				OffchainPublicKey:     chain.Ocr2Config.OcrKeyBundle.OffchainPublicKey,
				ConfigPublicKey:       chain.Ocr2Config.OcrKeyBundle.ConfigPublicKey,
			},
		},
	}
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
