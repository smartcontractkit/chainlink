package keystone

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"

	v1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var (
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry" // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L392
	KeystoneForwarder    deployment.ContractType = "KeystoneForwarder"    // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
	OCR3Capability       deployment.ContractType = "OCR3Capability"       // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
	FeedConsumer         deployment.ContractType = "FeedConsumer"         // no type and a version in contract https://github.com/smartcontractkit/chainlink/blob/89183a8a5d22b1aeca0ade3b76d16aa84067aa57/contracts/src/v0.8/keystone/KeystoneFeedsConsumer.sol#L1
)

type DeployResponse struct {
	Address common.Address
	Tx      common.Hash // todo: chain agnostic
	Tv      deployment.TypeAndVersion
}

type DeployRequest struct {
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
	// eth address is the first 20 bytes of the Signer
	return common.BytesToAddress(o.Signer[:20])
}

func (o *ocr2Node) toNodeKeys() NodeKeys {
	var aptosOcr2KeyBundleId string
	var aptosOnchainPublicKey string
	if o.aptosOcr2KeyBundle != nil {
		aptosOcr2KeyBundleId = o.aptosOcr2KeyBundle.BundleId
		aptosOnchainPublicKey = o.aptosOcr2KeyBundle.OnchainSigningAddress
	}
	return NodeKeys{
		EthAddress:            o.accountAddress,
		P2PPeerID:             strings.TrimPrefix(o.p2pKeyBundle.PeerId, "p2p_"),
		OCR2BundleID:          o.ethOcr2KeyBundle.BundleId,
		OCR2OnchainPublicKey:  o.ethOcr2KeyBundle.OnchainSigningAddress,
		OCR2OffchainPublicKey: o.ethOcr2KeyBundle.OffchainPublicKey,
		OCR2ConfigPublicKey:   o.ethOcr2KeyBundle.ConfigPublicKey,
		CSAPublicKey:          o.csaKey,
		// default value of encryption public key is the CSA public key
		// TODO: DEVSVCS-760
		EncryptionPublicKey: strings.TrimPrefix(o.csaKey, "csa_"),
		// TODO Aptos support. How will that be modeled in clo data?
		AptosBundleID:         aptosOcr2KeyBundleId,
		AptosOnchainPublicKey: aptosOnchainPublicKey,
	}
}
func newOcr2NodeFromJD(n *Node, registryChainSel uint64) (*ocr2Node, error) {
	if n.PublicKey == nil {
		return nil, errors.New("no public key")
	}
	// the chain configs are equivalent as far as the ocr2 config is concerned so take the first one
	if len(n.ChainConfigs) == 0 {
		return nil, errors.New("no chain configs")
	}
	// all nodes should have an evm chain config, specifically the registry chain
	evmCC, err := registryChainConfig(n.ChainConfigs, v1.ChainType_CHAIN_TYPE_EVM, registryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry chain config for sel %d: %w", registryChainSel, err)
	}
	cfgs := map[chaintype.ChainType]*v1.ChainConfig{
		chaintype.EVM: evmCC,
	}
	aptosCC, exists := firstChainConfigByType(n.ChainConfigs, v1.ChainType_CHAIN_TYPE_APTOS)
	if exists {
		cfgs[chaintype.Aptos] = aptosCC
	}
	return newOcr2Node(n.ID, cfgs, *n.PublicKey)
}

func ExtractKeys(n *Node, registerChainSel uint64) (p2p p2pkey.PeerID, signer [32]byte, encPubKey [32]byte, err error) {
	orc2n, err := newOcr2NodeFromJD(n, registerChainSel)
	if err != nil {
		return p2p, signer, encPubKey, fmt.Errorf("failed to create ocr2 node for node %s: %w", n.ID, err)
	}
	return orc2n.P2PKey, orc2n.Signer, orc2n.EncryptionPublicKey, nil
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

type NOP struct {
	Name  string
	Nodes []string // peerID
}

func (v NOP) Validate() error {
	if v.Name == "" {
		return errors.New("name is empty")
	}
	if len(v.Nodes) == 0 {
		return errors.New("no nodes")
	}
	for i, n := range v.Nodes {
		_, err := p2pkey.MakePeerID(n)
		if err != nil {
			return fmt.Errorf("failed to nop %s: node %d is not valid peer id %s: %w", v.Name, i, n, err)
		}
	}

	return nil
}

// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
type DonCapabilities struct {
	Name         string
	Nops         []NOP
	Capabilities []kcr.CapabilitiesRegistryCapability // every capability is hosted on each nop
}

func (v DonCapabilities) Validate() error {
	if v.Name == "" {
		return errors.New("name is empty")
	}
	if len(v.Nops) == 0 {
		return errors.New("no nops")
	}
	for i, n := range v.Nops {
		if err := n.Validate(); err != nil {
			return fmt.Errorf("failed to validate nop %d '%s': %w", i, n.Name, err)
		}
	}
	if len(v.Capabilities) == 0 {
		return errors.New("no capabilities")
	}
	return nil
}

func NodeOperator(name string, adminAddress string) capabilities_registry.CapabilitiesRegistryNodeOperator {
	return capabilities_registry.CapabilitiesRegistryNodeOperator{
		Name:  name,
		Admin: adminAddr(adminAddress),
	}
}

func AdminAddress(n *Node, chainSel uint64) (string, error) {
	cid, err := chainsel.ChainIdFromSelector(chainSel)
	if err != nil {
		return "", fmt.Errorf("failed to get chain id from selector %d: %w", chainSel, err)
	}
	cidStr := strconv.FormatUint(cid, 10)
	for _, chain := range n.ChainConfigs {
		//TODO validate chainType field
		if chain.Chain.Id == cidStr {
			return chain.AdminAddress, nil
		}
	}
	return "", fmt.Errorf("no chain config for chain %d", cid)
}

func nopsToNodes(donInfos []DonInfo, dons []DonCapabilities, chainSelector uint64) (map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string, error) {
	out := make(map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string)
	for _, don := range dons {
		for _, nop := range don.Nops {
			idx := slices.IndexFunc(donInfos, func(donInfo DonInfo) bool {
				return donInfo.Name == don.Name
			})
			if idx < 0 {
				return nil, fmt.Errorf("couldn't find donInfo for %v", don.Name)
			}
			donInfo := donInfos[idx]
			idx = slices.IndexFunc(donInfo.Nodes, func(node Node) bool {
				return node.P2PID == nop.Nodes[0]
			})
			if idx < 0 {
				return nil, fmt.Errorf("couldn't find node with p2p_id %v", nop.Nodes[0])
			}
			node := donInfo.Nodes[idx]
			a, err := AdminAddress(&node, chainSelector)
			if err != nil {
				return nil, fmt.Errorf("failed to get admin address for node %s: %w", node.ID, err)
			}
			nodeOperator := NodeOperator(nop.Name, a)
			for _, node := range nop.Nodes {

				idx = slices.IndexFunc(donInfo.Nodes, func(n Node) bool {
					return n.P2PID == node
				})
				if idx < 0 {
					return nil, fmt.Errorf("couldn't find node with p2p_id %v", node)
				}
				out[nodeOperator] = append(out[nodeOperator], donInfo.Nodes[idx].ID)

			}
		}
	}

	return out, nil
}

// mapDonsToCaps converts a list of DonCapabilities to a map of don name to capabilities
func mapDonsToCaps(dons []DonInfo) map[string][]kcr.CapabilitiesRegistryCapability {
	out := make(map[string][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		out[don.Name] = don.Capabilities
	}
	return out
}

// mapDonsToNodes returns a map of don name to simplified representation of their nodes
// all nodes must have evm config and ocr3 capability nodes are must also have an aptos chain config
func mapDonsToNodes(dons []DonInfo, excludeBootstraps bool, registryChainSel uint64) (map[string][]*ocr2Node, error) {
	donToOcr2Nodes := make(map[string][]*ocr2Node)
	// get the nodes for each don from the offchain client, get ocr2 config from one of the chain configs for the node b/c
	// they are equivalent, and transform to ocr2node representation

	for _, don := range dons {
		for _, node := range don.Nodes {
			ocr2n, err := newOcr2NodeFromJD(&node, registryChainSel)
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

	return donToOcr2Nodes, nil
}

func firstChainConfigByType(ccfgs []*v1.ChainConfig, t v1.ChainType) (*v1.ChainConfig, bool) {
	for _, c := range ccfgs {
		if c.Chain.Type == t {
			return c, true
		}
	}
	return nil, false
}

func registryChainConfig(ccfgs []*v1.ChainConfig, t v1.ChainType, sel uint64) (*v1.ChainConfig, error) {
	chainId, err := chainsel.ChainIdFromSelector(sel)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id from selector %d: %w", sel, err)
	}
	chainIdStr := strconv.FormatUint(chainId, 10)
	for _, c := range ccfgs {
		if c.Chain.Type == t && c.Chain.Id == chainIdStr {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no chain config for chain %d", chainId)
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

func joinInfoAndNodes(donInfos map[string]kcr.CapabilitiesRegistryDONInfo, dons []DonInfo, registryChainSel uint64) ([]RegisteredDon, error) {
	// all maps should have the same keys
	nodes, err := mapDonsToNodes(dons, true, registryChainSel)
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
