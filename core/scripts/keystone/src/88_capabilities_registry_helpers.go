package src

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type CapabilityRegistryProvisioner struct {
	reg kcr.CapabilitiesRegistryInterface
	env helpers.Environment
}

func NewCapabilityRegistryProvisioner(reg kcr.CapabilitiesRegistryInterface, env helpers.Environment) *CapabilityRegistryProvisioner {
	return &CapabilityRegistryProvisioner{reg: reg, env: env}
}

func extractRevertReason(errData string, a abi.ABI) (string, string, error) {
	data, err := hex.DecodeString(errData[2:])
	if err != nil {
		return "", "", err
	}

	for errName, abiError := range a.Errors {
		if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
			// Found a matching error
			v, err := abiError.Unpack(data)
			if err != nil {
				return "", "", err
			}
			b, err := json.Marshal(v)
			if err != nil {
				return "", "", err
			}
			return errName, string(b), nil
		}
	}
	return "", "", errors.New("revert Reason could not be found for given abistring")
}

func (c *CapabilityRegistryProvisioner) testCallContract(method string, args ...interface{}) error {
	abi := evmtypes.MustGetABI(kcr.CapabilitiesRegistryABI)
	data, err := abi.Pack(method, args...)
	helpers.PanicErr(err)
	cAddress := c.reg.Address()
	gasPrice, err := c.env.Ec.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)

	msg := ethereum.CallMsg{
		From:     c.env.Owner.From,
		To:       &cAddress,
		Data:     data,
		Gas:      10_000_000,
		GasPrice: gasPrice,
	}
	_, err = c.env.Ec.CallContract(context.Background(), msg, nil)
	if err != nil {
		if err.Error() == "execution reverted" {
			rpcError, ierr := evmclient.ExtractRPCError(err)
			if ierr != nil {
				return ierr
			}

			reason, abiErr, ierr := extractRevertReason(rpcError.Data.(string), abi)
			if ierr != nil {
				return ierr
			}

			e := fmt.Errorf("failed to call %s: reason: %s reasonargs: %s", method, reason, abiErr)
			return e
		}

		return err
	}

	return nil
}

// AddCapabilities takes a capability set and provisions it in the registry.
func (c *CapabilityRegistryProvisioner) AddCapabilities(ctx context.Context, capSet CapabilitySet) {
	fmt.Printf("Adding capabilities to registry: %s\n", capSet.IDs())
	tx, err := c.reg.AddCapabilities(c.env.Owner, capSet.Capabilities())

	helpers.PanicErr(err)
	helpers.ConfirmTXMined(ctx, c.env.Ec, tx, c.env.ChainID)
}

// AddNodeOperator takes a node operator and provisions it in the registry.
//
// A node operator is a group of nodes that are all controlled by the same entity. The admin address is the
// address that controls the node operator.
//
// The name is a human-readable name for the node operator.
//
// The node operator is then added to the registry, and the registry will issue an ID for the node operator.
// The ID is then used when adding nodes to the registry such that the registry knows which nodes belong to which
// node operator.
func (c *CapabilityRegistryProvisioner) AddNodeOperator(ctx context.Context, nop *NodeOperator) {
	fmt.Printf("Adding NodeOperator to registry: %s\n", nop.Name)
	nop.BindToRegistry(c.reg)

	nops, err := c.reg.GetNodeOperators(&bind.CallOpts{})
	if err != nil {
		log.Printf("failed to GetNodeOperators: %s", err)
	}
	for _, n := range nops {
		if n.Admin == nop.Admin {
			log.Printf("NodeOperator with admin address %s already exists", n.Admin.Hex())
			return
		}
	}

	tx, err := c.reg.AddNodeOperators(c.env.Owner, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: nop.Admin,
			Name:  nop.Name,
		},
	})
	if err != nil {
		log.Printf("failed to AddNodeOperators: %s", err)
	}

	receipt := helpers.ConfirmTXMined(ctx, c.env.Ec, tx, c.env.ChainID)
	nop.SetCapabilityRegistryIssuedID(receipt)
}

// AddNodes takes a node operators nodes, along with a capability set, then configures the registry such that
// each node is assigned the same capability set. The registry will then know that each node supports each of the
// capabilities in the set.
//
// This is a simplified version of the actual implementation, which is more flexible. The actual implementation
// allows for the ability to add different capability sets to different nodes, _and_ lets you add nodes from different
// node operators to the same capability set. This is not yet implemented here.
//
// Note that the registry must already have the capability set added via `AddCapabilities`, you cannot
// add capabilities that the registry is not yet aware of.
//
// Note that in terms of the provisioning process, this is not the last step. A capability is only active once
// there is a DON servicing it. This is done via `AddDON`.
func (c *CapabilityRegistryProvisioner) AddNodes(ctx context.Context, nop *NodeOperator, donNames ...string) {
	fmt.Printf("Adding nodes to registry for NodeOperator %s with DONs: %v\n", nop.Name, donNames)
	var params []kcr.CapabilitiesRegistryNodeParams
	for _, donName := range donNames {
		don, exists := nop.DONs[donName]
		if !exists {
			log.Fatalf("DON with name %s does not exist in NodeOperator %s", donName, nop.Name)
		}
		capSet := don.CapabilitySet
		for i, peer := range don.Peers {
			node, innerErr := peerToNode(nop.id, peer)
			if innerErr != nil {
				panic(innerErr)
			}
			node.HashedCapabilityIds = capSet.HashedIDs(c.reg)
			node.EncryptionPublicKey = [32]byte{2: byte(i + 1)}
			fmt.Printf("Adding node %s to registry with capabilities: %s\n", peer.PeerID, capSet.IDs())
			params = append(params, node)
		}
	}

	err := c.testCallContract("addNodes", params)
	PanicErr(err)

	tx, err := c.reg.AddNodes(c.env.Owner, params)
	if err != nil {
		log.Printf("failed to AddNodes: %s", err)
	}
	helpers.ConfirmTXMined(ctx, c.env.Ec, tx, c.env.ChainID)
}

// AddDON takes a node operator then provisions a DON with the given capabilities.
//
// A DON is a group of nodes that all support the same capability set. This set can be a subset of the
// capabilities that the nodes support. In other words, each node within the node set can support
// a different, possibly overlapping, set of capabilities, but a DON is a subgroup of those nodes that all support
// the same set of capabilities.
//
// A node can belong to multiple DONs, but it must belong to one and only one workflow DON.
//
// A DON can be a capability DON or a workflow DON, or both.
//
// When you want to add solely a workflow DON, you should set `acceptsWorkflows` to true and
// `isPublic` to false.
// This means that the DON can service workflow requests and will not service external capability requests.
//
// If you want to add solely a capability DON, you should set `acceptsWorkflows` to false and `isPublic` to true. This means that the DON
// will service external capability requests and reject workflow requests.
//
// If you want to add a DON that services both capabilities and workflows, you should set both `acceptsWorkflows` and `isPublic` to true.
//
// Another important distinction is that DON can comprise of nodes from different node operators, but for now, we're keeping it simple and restricting it to a single node operator. We also hard code F to 1.
func (c *CapabilityRegistryProvisioner) AddDON(ctx context.Context, nop *NodeOperator, donName string, isPublic bool, acceptsWorkflows bool) {
	fmt.Printf("Adding DON %s to registry for NodeOperator %s with isPublic: %t and acceptsWorkflows: %t\n", donName, nop.Name, isPublic, acceptsWorkflows)
	don, exists := nop.DONs[donName]
	if !exists {
		log.Fatalf("DON with name %s does not exist in NodeOperator %s", donName, nop.Name)
	}
	configs := don.CapabilitySet.Configs(c.reg)

	err := c.testCallContract("addDON", don.MustGetPeerIDs(), configs, isPublic, acceptsWorkflows, don.F)
	PanicErr(err)

	tx, err := c.reg.AddDON(c.env.Owner, don.MustGetPeerIDs(), configs, isPublic, acceptsWorkflows, don.F)

	if err != nil {
		log.Printf("failed to AddDON: %s", err)
	}
	helpers.ConfirmTXMined(ctx, c.env.Ec, tx, c.env.ChainID)
}

/*
 *
 * Capabilities
 *
 *
 */
const ( // Taken from https://github.com/smartcontractkit/chainlink/blob/29117850e9be1be1993dbf8f21cf13cbb6af9d24/core/capabilities/integration_tests/keystone_contracts_setup.go#L43
	CapabilityTypeTrigger   = uint8(0)
	CapabilityTypeAction    = uint8(1)
	CapabilityTypeConsensus = uint8(2)
	CapabilityTypeTarget    = uint8(3)
)

type CapabillityProvisioner interface {
	Config() kcr.CapabilitiesRegistryCapabilityConfiguration
	Capability() kcr.CapabilitiesRegistryCapability
	BindToRegistry(reg kcr.CapabilitiesRegistryInterface)
	GetHashedCID() [32]byte
}

type baseCapability struct {
	registry   kcr.CapabilitiesRegistryInterface
	capability kcr.CapabilitiesRegistryCapability
}

func (b *baseCapability) BindToRegistry(reg kcr.CapabilitiesRegistryInterface) {
	b.registry = reg
}

func (b *baseCapability) GetHashedCID() [32]byte {
	if b.registry == nil {
		panic(errors.New("registry not bound to capability, cannot get hashed capability ID"))
	}

	return mustHashCapabilityID(b.registry, b.capability)
}

func (b *baseCapability) GetID() string {
	return fmt.Sprintf("%s@%s", b.capability.LabelledName, b.capability.Version)
}

func (b *baseCapability) config(config *capabilitiespb.CapabilityConfig) kcr.CapabilitiesRegistryCapabilityConfiguration {
	configBytes, err := proto.Marshal(config)
	if err != nil {
		panic(err)
	}

	return kcr.CapabilitiesRegistryCapabilityConfiguration{
		Config:       configBytes,
		CapabilityId: b.GetHashedCID(),
	}
}

func (b *baseCapability) Capability() kcr.CapabilitiesRegistryCapability {
	return b.capability
}

type ConsensusCapability struct {
	baseCapability
}

var _ CapabillityProvisioner = &ConsensusCapability{}

func (c *ConsensusCapability) Config() kcr.CapabilitiesRegistryCapabilityConfiguration {
	// Note that this is hard-coded for now, we'll want to support more flexible configurations in the future
	// for configuring consensus once it has more configuration options
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}

	return c.config(config)
}

// NewOCR3V1ConsensusCapability returns a new ConsensusCapability for OCR3
func NewOCR3V1ConsensusCapability() *ConsensusCapability {
	return &ConsensusCapability{
		baseCapability{
			capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: CapabilityTypeConsensus,
			},
		},
	}
}

type TargetCapability struct {
	baseCapability
}

var _ CapabillityProvisioner = &TargetCapability{}

func (t *TargetCapability) Config() kcr.CapabilitiesRegistryCapabilityConfiguration {
	// Note that this is hard-coded for now, we'll want to support more flexible configurations in the future
	// for configuring the target. This configuration is also specific to the write target
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
			RemoteTargetConfig: &capabilitiespb.RemoteTargetConfig{
				RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
			},
		},
	}

	return t.config(config)
}

func NewEthereumGethTestnetV1WriteCapability() *TargetCapability {
	return &TargetCapability{
		baseCapability{
			capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "write_geth-testnet",
				Version:        "1.0.0",
				CapabilityType: CapabilityTypeTarget,
			},
		},
	}
}

type TriggerCapability struct {
	baseCapability
}

var _ CapabillityProvisioner = &TriggerCapability{}

func (t *TriggerCapability) Config() kcr.CapabilitiesRegistryCapabilityConfiguration {
	// Note that this is hard-coded for now, we'll want to support more flexible configurations in the future
	// for configuring the trigger. This configuration is also possibly specific to the streams trigger.
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1, // We've hardcoded F + 1 here
			},
		},
	}

	return t.config(config)
}

func NewStreamsTriggerV1Capability() *TriggerCapability {
	return &TriggerCapability{
		baseCapability{
			capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "streams-trigger",
				Version:        "1.1.0",
				CapabilityType: CapabilityTypeTrigger,
			},
		},
	}
}

func mustHashCapabilityID(reg kcr.CapabilitiesRegistryInterface, capability kcr.CapabilitiesRegistryCapability) [32]byte {
	hashedCapabilityID, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, capability.LabelledName, capability.Version)
	if err != nil {
		panic(err)
	}
	return hashedCapabilityID
}

/*
 *
 * Capability Sets
 *
 *
 */
type CapabilitySet []CapabillityProvisioner

func NewCapabilitySet(capabilities ...CapabillityProvisioner) CapabilitySet {
	if len(capabilities) == 0 {
		log.Fatalf("No capabilities provided to NewCapabilitySet")
	}

	return capabilities
}

func MergeCapabilitySets(sets ...CapabilitySet) CapabilitySet {
	var merged CapabilitySet
	for _, set := range sets {
		merged = append(merged, set...)
	}

	return merged
}

func (c *CapabilitySet) Capabilities() []kcr.CapabilitiesRegistryCapability {
	definitions := make([]kcr.CapabilitiesRegistryCapability, 0, len(*c))
	for _, cap := range *c {
		definitions = append(definitions, cap.Capability())
	}

	return definitions
}

func (c *CapabilitySet) IDs() []string {
	strings := make([]string, 0, len(*c))
	for _, cap := range *c {
		strings = append(strings, fmt.Sprintf("%s@%s", cap.Capability().LabelledName, cap.Capability().Version))
	}

	return strings
}

func (c *CapabilitySet) HashedIDs(reg kcr.CapabilitiesRegistryInterface) [][32]byte {
	ids := make([][32]byte, 0, len(*c))
	for _, cap := range *c {
		cap.BindToRegistry(reg)
		ids = append(ids, cap.GetHashedCID())
	}

	return ids
}

func (c *CapabilitySet) Configs(reg kcr.CapabilitiesRegistryInterface) []kcr.CapabilitiesRegistryCapabilityConfiguration {
	configs := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, 0, len(*c))
	for _, cap := range *c {
		cap.BindToRegistry(reg)
		configs = append(configs, cap.Config())
	}

	return configs
}

/*
 *
 * Node Operator
 *
 *
 */

// DON represents a Decentralized Oracle Network with a name, peers, and associated capabilities.
type DON struct {
	F             uint8
	Name          string
	Peers         []peer
	CapabilitySet CapabilitySet
}

// MustGetPeerIDs retrieves the peer IDs for the DON. It panics if any error occurs.
func (d *DON) MustGetPeerIDs() [][32]byte {
	ps, err := peers(d.Peers)
	if err != nil {
		panic(fmt.Errorf("failed to get peer IDs for DON %s: %w", d.Name, err))
	}
	return ps
}

// NodeOperator represents a node operator with administrative details and multiple DONs.
type NodeOperator struct {
	Admin gethCommon.Address
	Name  string
	DONs  map[string]DON

	reg kcr.CapabilitiesRegistryInterface
	// This ID is generated by the registry when the NodeOperator is added
	id uint32
}

// NewNodeOperator creates a new NodeOperator with the provided admin address, name, and DONs.
func NewNodeOperator(admin gethCommon.Address, name string, dons map[string]DON) *NodeOperator {
	return &NodeOperator{
		Admin: admin,
		Name:  name,
		DONs:  dons,
	}
}

func (n *NodeOperator) BindToRegistry(reg kcr.CapabilitiesRegistryInterface) {
	n.reg = reg
}

func (n *NodeOperator) SetCapabilityRegistryIssuedID(receipt *gethTypes.Receipt) uint32 {
	if n.reg == nil {
		panic(errors.New("registry not bound to node operator, cannot set ID"))
	}
	// We'll need more complex handling for multiple node operators
	// since we'll need to handle log ordering
	recLog, err := n.reg.ParseNodeOperatorAdded(*receipt.Logs[0])
	if err != nil {
		panic(err)
	}

	n.id = recLog.NodeOperatorId
	return n.id
}

func peerIDToB(peerID string) ([32]byte, error) {
	var peerIDB ragetypes.PeerID
	err := peerIDB.UnmarshalText([]byte(peerID))
	if err != nil {
		return [32]byte{}, err
	}

	return peerIDB, nil
}

func peers(ps []peer) ([][32]byte, error) {
	out := [][32]byte{}
	for _, p := range ps {
		b, err := peerIDToB(p.PeerID)
		if err != nil {
			return nil, err
		}

		out = append(out, b)
	}

	return out, nil
}

func peerToNode(nopID uint32, p peer) (kcr.CapabilitiesRegistryNodeParams, error) {
	peerIDB, err := peerIDToB(p.PeerID)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert peerID: %w", err)
	}

	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	epk := strings.TrimPrefix(p.EncryptionPublicKey, "0x")
	epkB, err := hex.DecodeString(epk)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert encryptionPublicKey: %w", err)
	}

	var epkb [32]byte
	copy(epkb[:], epkB)

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      nopID,
		P2pId:               peerIDB,
		Signer:              sigb,
		EncryptionPublicKey: epkb,
	}, nil
}

// newCapabilityConfig returns a new capability config with the default config set as empty.
// Override the empty default config with functional options.
func newCapabilityConfig(opts ...func(*values.Map)) *capabilitiespb.CapabilityConfig {
	dc := values.EmptyMap()
	for _, opt := range opts {
		opt(dc)
	}

	return &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.ProtoMap(dc),
	}
}
