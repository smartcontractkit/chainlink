package deployment

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func UBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

func E18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(UBigInt(amount), UBigInt(1e18))
}

type OCRConfig struct {
	OffchainPublicKey types2.OffchainPublicKey
	// For EVM-chains, this an *address*.
	OnchainPublicKey          types2.OnchainPublicKey
	PeerID                    p2pkey.PeerID
	TransmitAccount           types2.Account
	ConfigEncryptionPublicKey types3.ConfigEncryptionPublicKey
	KeyBundleID               string
}

func (ocrCfg OCRConfig) JDOCR2KeyBundle() *nodev1.OCR2Config_OCRKeyBundle {
	return &nodev1.OCR2Config_OCRKeyBundle{
		OffchainPublicKey:     hex.EncodeToString(ocrCfg.OffchainPublicKey[:]),
		OnchainSigningAddress: hex.EncodeToString(ocrCfg.OnchainPublicKey),
		ConfigPublicKey:       hex.EncodeToString(ocrCfg.ConfigEncryptionPublicKey[:]),
		BundleId:              ocrCfg.KeyBundleID,
	}
}

// Nodes includes is a group CL nodes.
type Nodes []Node

// PeerIDs returns peerIDs in a sorted list
func (n Nodes) PeerIDs() [][32]byte {
	var peerIDs [][32]byte
	for _, node := range n {
		peerIDs = append(peerIDs, node.PeerID)
	}
	sort.Slice(peerIDs, func(i, j int) bool {
		return bytes.Compare(peerIDs[i][:], peerIDs[j][:]) < 0
	})
	return peerIDs
}

func (n Nodes) NonBootstraps() Nodes {
	var nonBootstraps Nodes
	for _, node := range n {
		if node.IsBootstrap {
			continue
		}
		nonBootstraps = append(nonBootstraps, node)
	}
	return nonBootstraps
}

func (n Nodes) DefaultF() uint8 {
	return uint8(len(n) / 3)
}

func (n Nodes) IDs() []string {
	var ids []string
	for _, node := range n {
		ids = append(ids, node.NodeID)
	}
	return ids
}

func (n Nodes) BootstrapLocators() []string {
	bootstrapMp := make(map[string]struct{})
	for _, node := range n {
		if node.IsBootstrap {
			key := node.MultiAddr
			// compatibility with legacy code. unclear what code path is setting half baked node.MultiAddr
			if !isValidMultiAddr(key) {
				key = fmt.Sprintf("%s@%s", strings.TrimPrefix(node.PeerID.String(), "p2p_"), node.MultiAddr)
			}
			bootstrapMp[key] = struct{}{}
		}
	}
	var locators []string
	for b := range bootstrapMp {
		locators = append(locators, b)
	}
	return locators
}

// P2PIDsPresentInJD - For a given p2pIDs, check if the nodes are present in JD.
func (n Nodes) P2PIDsPresentInJD(p2pIDs [][32]byte) error {
	var allErrs error
	for _, p2pID := range p2pIDs {
		p2pIDString := "p2p_" + libocrtypes.PeerID(p2pID).String()
		if !slices.ContainsFunc(n, func(n Node) bool {
			return p2pIDString == n.PeerID.String()
		}) {
			allErrs = multierror.Append(allErrs, fmt.Errorf("node with p2pID %s not found in JD", p2pIDString))
		}
	}
	return allErrs
}

func isValidMultiAddr(s string) bool {
	// Define the regular expression pattern
	pattern := `^(.+)@(.+):(\d+)$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 4 { // 4 because the entire match + 3 submatches
		return false
	}

	_, err := p2pkey.MakePeerID("p2p_" + matches[1])
	return err == nil
}

type Node struct {
	NodeID         string
	Name           string
	CSAKey         string
	WorkflowKey    string
	SelToOCRConfig map[chain_selectors.ChainDetails]OCRConfig
	PeerID         p2pkey.PeerID
	IsBootstrap    bool
	MultiAddr      string
	AdminAddr      string
	Labels         []*ptypes.Label
}

func (n Node) OCRConfigForChainDetails(details chain_selectors.ChainDetails) (OCRConfig, bool) {
	c, ok := n.SelToOCRConfig[details]
	return c, ok
}

func (n Node) AllOCRConfigs() map[chain_selectors.ChainDetails]OCRConfig {
	return n.SelToOCRConfig
}

func (n Node) OCRConfigForChainSelector(chainSel uint64) (OCRConfig, bool) {
	fam, err := chain_selectors.GetSelectorFamily(chainSel)
	if err != nil {
		return OCRConfig{}, false
	}

	id, err := chain_selectors.GetChainIDFromSelector(chainSel)
	if err != nil {
		return OCRConfig{}, false
	}

	want, err := chain_selectors.GetChainDetailsByChainIDAndFamily(id, fam)
	if err != nil {
		return OCRConfig{}, false
	}
	// only applicable for test related simulated chains, the chains don't have a name
	if want.ChainName == "" {
		want.ChainName = strconv.FormatUint(want.ChainSelector, 10)
	}
	c, ok := n.SelToOCRConfig[want]
	return c, ok
}

// ChainConfigs returns the chain configs for this node
// in the format required by JD
//
// WARNING: this is a lossy conversion because the Node abstraction
// is not as rich as the JD abstraction
func (n Node) ChainConfigs() ([]*nodev1.ChainConfig, error) {
	var out []*nodev1.ChainConfig
	for details, ocrCfg := range n.SelToOCRConfig {
		c, err := detailsToChain(details)
		if err != nil {
			return nil, fmt.Errorf("failed to get convert chain details: %w", err)
		}
		out = append(out, &nodev1.ChainConfig{
			Chain: c,
			// only have ocr2 in Node
			Ocr2Config: &nodev1.OCR2Config{
				OcrKeyBundle: ocrCfg.JDOCR2KeyBundle(),
				P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
					PeerId: n.PeerID.String(),
					// note: we don't have the public key in the OCRConfig struct
				},
				IsBootstrap: n.IsBootstrap,
				Multiaddr:   n.MultiAddr,
			},
			AccountAddress: string(ocrCfg.TransmitAccount),
			AdminAddress:   n.AdminAddr,
			NodeId:         n.NodeID,
		})
	}
	return out, nil
}

func MustPeerIDFromString(s string) p2pkey.PeerID {
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(s); err != nil {
		panic(err)
	}
	return p
}

type NodeChainConfigsLister interface {
	ListNodes(ctx context.Context, in *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error)
	ListNodeChainConfigs(ctx context.Context, in *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error)
}

var ErrMissingNodeMetadata = errors.New("missing node metadata")

// Gathers all the node info through JD required to be able to set
// OCR config for example. nodeIDs can be JD IDs or PeerIDs
//
// It is optimistic execution and will attempt to return an element for all
// nodes in the input list that exists in JD
//
// If some subset of nodes cannot have all their metadata returned, the error with be
// [ErrMissingNodeMetadata] and the caller can choose to handle or continue.
func NodeInfo(nodeIDs []string, oc NodeChainConfigsLister) (Nodes, error) {
	if len(nodeIDs) == 0 {
		return nil, nil
	}
	// if nodeIDs starts with `p2p_` lookup by p2p_id instead
	filterByPeerIDs := strings.HasPrefix(nodeIDs[0], "p2p_")
	var filter *nodev1.ListNodesRequest_Filter
	if filterByPeerIDs {
		selector := strings.Join(nodeIDs, ",")
		filter = &nodev1.ListNodesRequest_Filter{
			Enabled: 1,
			Selectors: []*ptypes.Selector{
				{
					Key:   "p2p_id",
					Op:    ptypes.SelectorOp_IN,
					Value: &selector,
				},
			},
		}
	} else {
		filter = &nodev1.ListNodesRequest_Filter{
			Enabled: 1,
			Ids:     nodeIDs,
		}
	}
	nodesFromJD, err := oc.ListNodes(context.Background(), &nodev1.ListNodesRequest{
		Filter: filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var nodes []Node
	onlyMissingEVMChain := true
	var xerr error
	for _, node := range nodesFromJD.GetNodes() {
		// TODO: Filter should accept multiple nodes
		nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{node.Id},
		}})
		if err != nil {
			return nil, fmt.Errorf("failed to list node chain configs for node %s id %s: %w", node.Name, node.Id, err)
		}
		n, err := NewNodeFromJD(node, nodeChainConfigs.ChainConfigs)
		if err != nil {
			xerr = errors.Join(xerr, fmt.Errorf("failed to get node metadata for node %s id %s: %w", node.Name, node.Id, err))
			if !errors.Is(err, ErrMissingEVMChain) {
				onlyMissingEVMChain = false
			}
		}
		nodes = append(nodes, *n)
	}
	if xerr != nil && onlyMissingEVMChain {
		xerr = errors.Join(ErrMissingNodeMetadata, xerr)
	}
	return nodes, xerr
}

var ErrMissingEVMChain = errors.New("no EVM chain found")

// NewNodeFromJD creates a Node from a JD Node. Populating all the fields requires an enabled
// EVM chain and OCR2 config. If this does not exist, the Node will be returned with
// the minimal fields populated and return a [ErrMissingEVMChain] error.
func NewNodeFromJD(jdNode *nodev1.Node, chainConfigs []*nodev1.ChainConfig) (*Node, error) {
	// the protobuf does not map well to the domain model
	// we have to infer the p2p key, bootstrap and multiaddr from some chain config
	// arbitrarily pick the first EVM chain config
	// we use EVM because the home or registry chain is always EVM
	emptyNode := &Node{
		NodeID:         jdNode.Id,
		Name:           jdNode.Name,
		CSAKey:         jdNode.PublicKey,
		WorkflowKey:    jdNode.GetWorkflowKey(),
		SelToOCRConfig: make(map[chain_selectors.ChainDetails]OCRConfig),
	}
	var goldenConfig *nodev1.ChainConfig
	for _, chainConfig := range chainConfigs {
		if chainConfig.Chain.Type == nodev1.ChainType_CHAIN_TYPE_EVM {
			goldenConfig = chainConfig
			break
		}
	}
	if goldenConfig == nil {
		return emptyNode, fmt.Errorf("node '%s', id '%s', csa '%s': %w", jdNode.Name, jdNode.Id, jdNode.PublicKey, ErrMissingEVMChain)
	}
	selToOCRConfig := make(map[chain_selectors.ChainDetails]OCRConfig)
	bootstrap := goldenConfig.Ocr2Config.IsBootstrap
	if !bootstrap { // no ocr config on bootstrap
		var err error
		selToOCRConfig, err = ChainConfigsToOCRConfig(chainConfigs)
		if err != nil {
			return emptyNode, fmt.Errorf("failed to get chain to ocr config: %w", err)
		}
	}
	return &Node{
		NodeID:         jdNode.Id,
		Name:           jdNode.Name,
		CSAKey:         jdNode.PublicKey,
		WorkflowKey:    jdNode.GetWorkflowKey(),
		SelToOCRConfig: selToOCRConfig,
		IsBootstrap:    bootstrap,
		PeerID:         MustPeerIDFromString(goldenConfig.Ocr2Config.P2PKeyBundle.PeerId),
		MultiAddr:      goldenConfig.Ocr2Config.Multiaddr,
		AdminAddr:      goldenConfig.AdminAddress,
		Labels:         jdNode.Labels,
	}, nil
}

func ChainConfigsToOCRConfig(chainConfigs []*nodev1.ChainConfig) (map[chain_selectors.ChainDetails]OCRConfig, error) {
	selToOCRConfig := make(map[chain_selectors.ChainDetails]OCRConfig)
	for _, chainConfig := range chainConfigs {
		b := common.Hex2Bytes(chainConfig.Ocr2Config.OcrKeyBundle.OffchainPublicKey)
		var opk types2.OffchainPublicKey
		copy(opk[:], b)

		b = common.Hex2Bytes(chainConfig.Ocr2Config.OcrKeyBundle.ConfigPublicKey)
		var cpk types3.ConfigEncryptionPublicKey
		copy(cpk[:], b)

		var pubkey types3.OnchainPublicKey
		if chainConfig.Chain.Type == nodev1.ChainType_CHAIN_TYPE_EVM {
			// convert from pubkey to address
			pubkey = common.HexToAddress(chainConfig.Ocr2Config.OcrKeyBundle.OnchainSigningAddress).Bytes()
		} else {
			pubkey = common.Hex2Bytes(chainConfig.Ocr2Config.OcrKeyBundle.OnchainSigningAddress)
		}

		details, err := chainToDetails(chainConfig.Chain)
		if err != nil {
			return nil, err
		}

		selToOCRConfig[details] = OCRConfig{
			OffchainPublicKey:         opk,
			OnchainPublicKey:          pubkey,
			PeerID:                    MustPeerIDFromString(chainConfig.Ocr2Config.P2PKeyBundle.PeerId),
			TransmitAccount:           types2.Account(chainConfig.AccountAddress),
			ConfigEncryptionPublicKey: cpk,
			KeyBundleID:               chainConfig.Ocr2Config.OcrKeyBundle.BundleId,
		}
	}
	return selToOCRConfig, nil
}

func chainToDetails(c *nodev1.Chain) (chain_selectors.ChainDetails, error) {
	var family string
	switch c.Type {
	case nodev1.ChainType_CHAIN_TYPE_EVM:
		family = chain_selectors.FamilyEVM
	case nodev1.ChainType_CHAIN_TYPE_APTOS:
		family = chain_selectors.FamilyAptos
	case nodev1.ChainType_CHAIN_TYPE_SOLANA:
		family = chain_selectors.FamilySolana
	case nodev1.ChainType_CHAIN_TYPE_STARKNET:
		family = chain_selectors.FamilyStarknet
	default:
		return chain_selectors.ChainDetails{}, fmt.Errorf("unsupported chain type %s", c.Type)
	}
	if family == chain_selectors.FamilySolana {
		// Temporary workaround to handle cases when solana chainId was not using the standard genesis hash,
		// but using old strings mainnet/testnet/devnet.
		switch c.Id {
		case "mainnet":
			c.Id = "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
		case "devnet":
			c.Id = "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG"
		case "testnet":
			c.Id = "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY"
		}
	}
	details, err := chain_selectors.GetChainDetailsByChainIDAndFamily(c.Id, family)
	if err != nil {
		return chain_selectors.ChainDetails{}, err
	}
	// only applicable for test related simulated chains, the chains don't have a name
	if details.ChainName == "" {
		details.ChainName = strconv.FormatUint(details.ChainSelector, 10)
	}
	return details, nil
}

func detailsToChain(details chain_selectors.ChainDetails) (*nodev1.Chain, error) {
	family, err := chain_selectors.GetSelectorFamily(details.ChainSelector)
	if err != nil {
		return nil, err
	}

	var t nodev1.ChainType
	switch family {
	case chain_selectors.FamilyEVM:
		t = nodev1.ChainType_CHAIN_TYPE_EVM
	case chain_selectors.FamilyAptos:
		t = nodev1.ChainType_CHAIN_TYPE_APTOS
	case chain_selectors.FamilySolana:
		t = nodev1.ChainType_CHAIN_TYPE_SOLANA
	case chain_selectors.FamilyStarknet:
		t = nodev1.ChainType_CHAIN_TYPE_STARKNET
	default:
		return nil, fmt.Errorf("unsupported chain family %s", family)
	}

	id, err := chain_selectors.GetChainIDFromSelector(details.ChainSelector)
	if err != nil {
		return nil, err
	}

	return &nodev1.Chain{
		Type: t,
		Id:   id,
	}, nil
}

type CapabilityRegistryConfig struct {
	EVMChainID  uint64         // chain id of the chain the CR is deployed on
	Contract    common.Address // address of the CR contract
	NetworkType string         // network type of the chain
}
