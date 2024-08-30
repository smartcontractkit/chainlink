package keystone

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	v1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type DeployRequest struct {
	RegistryChain uint64
	Env           deployment.MultiDonEnvironment

	//DonToCapType      map[string]uint8
	DonToCapabilities map[string][]kcr.CapabilitiesRegistryCapability                   // from external source
	NodeIDToNop       map[string]capabilities_registry.CapabilitiesRegistryNodeOperator // maybe should be derivable from JD interface but doesn't seem to be notion of NOP in JD
}

type DeployResponse struct {
	AddressBook deployment.AddressBook
}

// ocr2Node is a subset of the node configuration that is needed to register a node
// with the capabilities registry. Signer and P2PKey are chain agnostic.
type ocr2Node struct {
	ID     string
	Signer [32]byte
	P2PKey p2pkey.PeerID
}

func newOcr2Node(id string, cfg *v1.OCR2Config) (*ocr2Node, error) {
	if cfg == nil {
		return nil, errors.New("nil ocr2config")
	}
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(cfg.P2PKeyBundle.PeerId); err != nil {
		return nil, fmt.Errorf("failed to unmarshal peer id %s: %w", cfg.P2PKeyBundle.PeerId, err)
	}

	signer := cfg.OcrKeyBundle.OnchainSigningAddress
	if len(signer) != 40 {
		return nil, fmt.Errorf("invalid onchain signing address %s", cfg.OcrKeyBundle.OnchainSigningAddress)
	}
	signerB, err := hex.DecodeString(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer %s: %w", signer, err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return &ocr2Node{
		ID:     id,
		Signer: sigb,
		P2PKey: p,
	}, nil
}

func nodeChainConfigsToOcr2Node(resp *v1.ListNodeChainConfigsResponse, id string) (*ocr2Node, error) {
	if len(resp.ChainConfigs) == 0 {
		return nil, errors.New("no chain configs")
	}
	cfg := resp.ChainConfigs[0]
	return newOcr2Node(id, cfg.Ocr2Config)
}

func Deploy(ctx context.Context, lggr logger.Logger, req DeployRequest) (*DeployResponse, error) {
	ad := deployment.NewMemoryAddressBook()
	resp := &DeployResponse{
		AddressBook: ad,
	}
	var registry *capabilities_registry.CapabilitiesRegistry
	var registryChain deployment.Chain

	for _, chain := range req.Env.ListChains() {
		lggr.Info("deploying contracts", "chain", chain)
		deployResp, err := deployContracts(req.Env.Logger, deployContractsRequest{
			chain:           chain,
			isRegistryChain: chain.Selector == req.RegistryChain,
		},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contracts: %w", err)
		}
		err = ad.Merge(deployResp.AddressBook)
		if err != nil {
			return nil, fmt.Errorf("failed to merge address book: %w", err)
		}
		if chain.Selector == req.RegistryChain {
			registry = deployResp.capabilitiesRegistryDeployer.contract
			registryChain = chain
		}
	}

	// todo cleanup types
	//donNodes := make(map[string]map[string][]*v1.ChainConfig) // donid -> nodeid -> chainconfigs
	donToOcr2Nodes := make(map[string][]*ocr2Node)
	for donName, env := range req.Env.Environments {
		donNodeSet, err := env.Offchain.ListNodes(ctx, &v1.ListNodesRequest{}) // each env is a don
		if err != nil {
			return nil, fmt.Errorf("failed to list nodes: %w", err)
		}
		if len(donNodeSet.Nodes) == 0 {
			return nil, fmt.Errorf("no nodes found")
		}
		// each node in the nodeset may support mulitple chains
		nodeCfgs := make(map[string][]*v1.ChainConfig)
		for _, node := range donNodeSet.Nodes {
			cfgResp, err := env.Offchain.ListNodeChainConfigs(ctx, &v1.ListNodeChainConfigsRequest{
				Filter: &v1.ListNodeChainConfigsRequest_Filter{NodeIds: []string{node.Id}},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list node chain configs for node %s: %w", node.Id, err)
			}
			nodeCfgs[node.Id] = cfgResp.ChainConfigs
			// convert to ocr2 node and store
			ocr2n, err := nodeChainConfigsToOcr2Node(cfgResp, node.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to convert node chain configs to ocr2 node for id %s: %w", node.Id, err)
			}
			if _, ok := donToOcr2Nodes[donName]; !ok {
				donToOcr2Nodes[donName] = make([]*ocr2Node, 0)
			}
			donToOcr2Nodes[donName] = append(donToOcr2Nodes[donName], ocr2n)
		}
		if len(nodeCfgs) == 0 {
			return nil, fmt.Errorf("no node chain configs found for don %s", donName)
		}
		//donNodes[donName] = nodeCfgs
		req.Env.Logger.Infow("node chain configs", "don", donName, "configs", nodeCfgs)
	}

	donToNodeIDs := make(map[string][]string)
	for donName, ocr2nodes := range donToOcr2Nodes {
		ids := make([]string, 0)
		for _, n := range ocr2nodes {
			ids = append(ids, n.ID)
		}
		donToNodeIDs[donName] = ids
	}

	// register capabilities
	capabilitiesResp, err := registerCapabilities(lggr, registerCapabilitiesRequest{
		chain:                registryChain,
		registry:             registry,
		donToCapabilities:    req.DonToCapabilities,
		donToOffchainNodeIDs: donToNodeIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register capabilities: %w", err)
	}
	lggr.Infow("registered capabilities", "capabilities", capabilitiesResp.donToCapabilities)

	// register node operators
	var nops []capabilities_registry.CapabilitiesRegistryNodeOperator
	for _, nop := range req.NodeIDToNop {
		nops = append(nops, nop)
	}
	nopsResp, err := registerNOPS(ctx, registerNOPSRequest{
		chain:    registryChain,
		registry: registry,
		nops:     nops,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register node operators: %w", err)
	}
	lggr.Infow("registered node operators", "nops", nopsResp.nops)

	nopToNodeIDs := make(map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string)
	for nodeID, nop := range req.NodeIDToNop {
		if _, ok := nopToNodeIDs[nop]; !ok {
			nopToNodeIDs[nop] = make([]string, 0)
		}
		nopToNodeIDs[nop] = append(nopToNodeIDs[nop], nodeID)
	}
	nodeToRegisterNop := make(map[string]*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded)
	for _, nop := range nopsResp.nops {
		n := capabilities_registry.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: nop.Admin,
		}
		nodeIDs := nopToNodeIDs[n]
		for _, nodeID := range nodeIDs {
			_, exists := nodeToRegisterNop[nodeID]
			if !exists {
				nodeToRegisterNop[nodeID] = nop
			}
		}
	}

	//	nopOnchainIDtoParams := make(map[uint32]capabilities_registry.CapabilitiesRegistryNodeParams)
	nodeIDToParams := make(map[string]capabilities_registry.CapabilitiesRegistryNodeParams)
	for don, ocrNodes := range donToOcr2Nodes {
		caps, ok := capabilitiesResp.donToCapabilities[don]
		var hashedCapabilityIds [][32]byte
		for _, cap := range caps {
			hashedCapabilityIds = append(hashedCapabilityIds, cap.id)
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", hashedCapabilityIds)
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		for _, n := range ocrNodes {
			nop, ok := nodeToRegisterNop[n.ID]
			if !ok {
				return nil, fmt.Errorf("node operator not found for node %s", n.ID)
			}
			params, ok := nodeIDToParams[n.ID]

			if !ok {
				params = capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      nop.NodeOperatorId,
					Signer:              n.Signer,
					P2pId:               n.P2PKey,
					HashedCapabilityIds: hashedCapabilityIds,
				}
			} else {
				// when we have a node operator, we need to add any new capabilities
				var newCapIds [][32]byte
				for _, proposedCapId := range hashedCapabilityIds {
					shouldAdd := true
					for _, existingCapId := range params.HashedCapabilityIds {
						if existingCapId == proposedCapId {
							shouldAdd = false
							break
						}
					}
					if shouldAdd {
						newCapIds = append(newCapIds, proposedCapId)
					}
				}
				params.HashedCapabilityIds = append(params.HashedCapabilityIds, newCapIds...)
			}
			//nopOnchainIDtoParams[nop.NodeOperatorId] = params
			nodeIDToParams[n.ID] = params
		}
	}
	lggr.Debugw("node params", "params", nodeIDToParams)
	var uniqueNodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, v := range nodeIDToParams {
		uniqueNodeParams = append(uniqueNodeParams, v)
	}
	lggr.Debug("unique node params", "params", uniqueNodeParams)
	tx, err := registry.AddNodes(registryChain.DeployerKey, uniqueNodeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to call AddNode: %w", err)
	}
	_, err = registryChain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash(), err)
	}

	lggr.Infow("registered nodes", "nodes", nodeIDToParams)

	for don, ocrNodes := range donToOcr2Nodes {
		var p2pIds [][32]byte
		for _, n := range ocrNodes {
			params, ok := nodeIDToParams[n.ID]
			if !ok {
				return nil, fmt.Errorf("node params not found for node %s", n)
			}
			p2pIds = append(p2pIds, params.P2pId)
		}

		caps, ok := capabilitiesResp.donToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		wfSupported := false
		var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
		for _, cap := range caps {
			if cap.CapabilityType == 2 { // OCR3
				wfSupported = true
			}
			capCfg := defaultCapConfig(cap.CapabilityType, len(p2pIds))
			cfgb, err := proto.Marshal(capCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal capability config for %v: %w", cap, err)
			}
			cfgs = append(cfgs, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: cap.id,
				Config:       cfgb,
			})
		}

		f := len(p2pIds) / 3 // assuming n=3f+1
		tx, err = registry.AddDON(registryChain.DeployerKey, p2pIds, cfgs, true, wfSupported, uint8(f))
		if err != nil {
			err = decodeErr(kcr.CapabilitiesRegistryABI, err)
			return nil, fmt.Errorf("failed to call AddDON for don %s p2p2Ids %v capability %v: %w", don, p2pIds, cfgs, err)
		}
		_, err = registryChain.Confirm(tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddDON transaction %s for don %s: %w", tx.Hash(), don, err)
		}
		lggr.Debugw("registered DON", "don", don, "p2pids", p2pIds, "cgs", cfgs, "wfSupported", wfSupported, "f", f)
	}

	lggr.Infow("registered DONS")

	return resp, err
}

type registerCapabilitiesRequest struct {
	chain                deployment.Chain
	registry             *capabilities_registry.CapabilitiesRegistry
	donToCapabilities    map[string][]kcr.CapabilitiesRegistryCapability
	donToOffchainNodeIDs map[string][]string
}

type registerCapabilitiesResponse struct {
	donToCapabilities          map[string][]registeredCapability
	offchainNodeIdToCapability map[string][]registeredCapability
}

type registeredCapability struct {
	capabilities_registry.CapabilitiesRegistryCapability
	id [32]byte
}

type capabilityId [32]byte

// func registerCapabilities(reg *capabilities_registry.CapabilitiesRegistry, chain deployment.Chain) error {
func registerCapabilities(lggr logger.Logger, req registerCapabilitiesRequest) (*registerCapabilitiesResponse, error) {
	if len(req.donToCapabilities) == 0 {
		return nil, fmt.Errorf("no capabilities to register")
	}
	resp := &registerCapabilitiesResponse{
		donToCapabilities:          make(map[string][]registeredCapability),
		offchainNodeIdToCapability: make(map[string][]registeredCapability),
	}

	uniqueOffchainNodeIDs := make(map[string]struct{})
	for _, nodeIDs := range req.donToOffchainNodeIDs {
		for _, nodeID := range nodeIDs {
			uniqueOffchainNodeIDs[nodeID] = struct{}{}
		}
	}
	offchainNodeIDs := make([]string, 0)
	for nodeID := range uniqueOffchainNodeIDs {
		offchainNodeIDs = append(offchainNodeIDs, nodeID)
	}

	// capability could be hosted on multiple dons. need to deduplicate
	uniqueCaps := make(map[kcr.CapabilitiesRegistryCapability][32]byte)
	for don, caps := range req.donToCapabilities {
		nodeIds, nodesExist := req.donToOffchainNodeIDs[don]
		if !nodesExist {
			return nil, fmt.Errorf("no offchain nodes found for don %s", don)
		}
		var registerCaps []registeredCapability
		for _, cap := range caps {
			id, ok := uniqueCaps[cap]
			if !ok {
				var err error
				id, err = req.registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
				if err != nil {
					return nil, fmt.Errorf("failed to call GetHashedCapabilityId for capability %v: %w", cap, err)
				}
				uniqueCaps[cap] = id
			}
			registerCap := registeredCapability{
				CapabilitiesRegistryCapability: cap,
				id:                             id,
			}
			lggr.Debugw("hashed capability id", "capability", cap, "id", id)
			registerCaps = append(registerCaps, registerCap)
		}
		for _, nodeID := range nodeIds {
			nodeCaps, ok := resp.offchainNodeIdToCapability[nodeID]
			if !ok {
				nodeCaps = make([]registeredCapability, 0)
			}
			// only add new capabilities
			var newCaps []registeredCapability
			for _, cap := range registerCaps {
				shouldAdd := true
				for _, existingCap := range nodeCaps {
					if existingCap.id == cap.id {
						shouldAdd = false
						break
					}
				}
				if shouldAdd {
					newCaps = append(newCaps, cap)
				}
			}
			resp.offchainNodeIdToCapability[nodeID] = append(resp.offchainNodeIdToCapability[nodeID], newCaps...)
		}
		resp.donToCapabilities[don] = registerCaps
	}

	var capabilities []kcr.CapabilitiesRegistryCapability
	for cap := range uniqueCaps {
		capabilities = append(capabilities, cap)
	}

	tx, err := req.registry.AddCapabilities(req.chain.DeployerKey, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to call AddCapabilities: %w", err)
	}
	_, err = req.chain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash(), err)
	}
	lggr.Info("registered capabilities", "capabilities", capabilities)
	return resp, nil
}

type registerNOPSRequest struct {
	chain    deployment.Chain
	registry *capabilities_registry.CapabilitiesRegistry
	//nodesets map[string][]liveNode // capabilities_registry.CapabilitiesRegistryNodeOperator // don id to node operators
	nops []capabilities_registry.CapabilitiesRegistryNodeOperator
}

type registerNOPSResponse struct {
	nops []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
}

func registerNOPS(ctx context.Context, req registerNOPSRequest) (*registerNOPSResponse, error) {
	//nops := uniqueNodeOperators(req.nodesets)
	nops := req.nops
	tx, err := req.registry.AddNodeOperators(req.chain.DeployerKey, nops)
	if err != nil {
		return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
	}
	// for some reason that i don't understand, the confirm must be called before the WaitMined or the latter will hang
	// (at least for a simulated backend chain)
	_, err = req.chain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNodeOperators confirm transaction %s: %w", tx.Hash(), err)
	}

	receipt, err := bind.WaitMined(ctx, req.chain.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to mine AddNodeOperators confirm transaction %s: %w", tx.Hash(), err)
	}
	if len(receipt.Logs) != len(nops) {
		return nil, fmt.Errorf("expected %d log entries for AddNodeOperators, got %d", len(nops), len(receipt.Logs))
	}
	resp := &registerNOPSResponse{
		nops: make([]*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded, len(receipt.Logs)),
	}
	for i, log := range receipt.Logs {
		o, err := req.registry.ParseNodeOperatorAdded(*log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log %d for operator added: %w", i, err)
		}
		resp.nops[i] = o
	}

	return resp, nil
}

type registeredNode struct {
	capabilities_registry.CapabilitiesRegistryNodeOperator
	deployment.Node
}

type addNodesRequest struct {
	chain        deployment.Chain
	registry     *capabilities_registry.CapabilitiesRegistry
	capNodeSet   map[capabilities_registry.CapabilitiesRegistryCapability][]registeredNode
	nodeToCapSet map[string][]capabilities_registry.CapabilitiesRegistryCapability
	nops         []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
	capabilities map[capabilities_registry.CapabilitiesRegistryCapability][32]byte
}

type addNodesResponse struct {
	nodes map[capabilities_registry.CapabilitiesRegistryCapability][]capabilities_registry.CapabilitiesRegistryNodeParams
}

func addNodes(req addNodesRequest) (*addNodesResponse, error) {
	return nil, nil
}

type addDonRequest struct {
	chain       deployment.Chain
	registry    *capabilities_registry.CapabilitiesRegistry
	capNodeSets map[capabilities_registry.CapabilitiesRegistryCapability][]capabilities_registry.CapabilitiesRegistryNodeParams
	capConfig   map[capabilities_registry.CapabilitiesRegistryCapability]*capabilitiespb.CapabilityConfig
	capID       map[capabilities_registry.CapabilitiesRegistryCapability][32]byte
}

type addDonResponse struct{}

func addDons(req addDonRequest) (*addDonResponse, error) {

	return &addDonResponse{}, nil
}

func defaultCapConfig(capType uint8, nNodes int) *capabilitiespb.CapabilityConfig {
	switch capType {
	// TODO: use the enum defined in ??
	case uint8(0): // trigger
		return &capabilitiespb.CapabilityConfig{
			DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
			RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
				RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
					RegistrationRefresh: durationpb.New(20 * time.Second),
					RegistrationExpiry:  durationpb.New(60 * time.Second),
					// F + 1; assuming n = 3f+1
					MinResponsesToAggregate: uint32(nNodes/3) + 1,
				},
			},
		}
	case uint8(2): // consensus
		return &capabilitiespb.CapabilityConfig{
			DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		}
	case uint8(3): // target
		return &capabilitiespb.CapabilityConfig{
			DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
			RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
				RemoteTargetConfig: &capabilitiespb.RemoteTargetConfig{
					RequestHashExcludedAttributes: []string{"signed_report.Signatures"}, // TODO: const defn in a common place
				},
			},
		}
	default:
		return &capabilitiespb.CapabilityConfig{
			DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		}
	}
}

func parseContractErr(encodedABI string, hexErr string) error {
	// abi is the abi of the contract that emitted the error
	// hexErrBytes is the hex encoded error data
	// returns an error that is the decoded error
	hexErr = strings.TrimPrefix(hexErr, "0x")
	errb, err := hex.DecodeString(hexErr)
	if err != nil {
		return fmt.Errorf("failed to decode hex error bytes: %w", err)
	}
	parsedAbi, err := abi.JSON(strings.NewReader(encodedABI))
	if err != nil {
		return fmt.Errorf("failed to parse abi: %w", err)
	}
	stringErr, err := abi.UnpackRevert(errb)
	if err == nil {
		return fmt.Errorf("string error: %s", stringErr)
	}
	for errName, abierr := range parsedAbi.Errors {
		fmt.Println("checking error", errName)
		if bytes.Equal(errb[:4], abierr.ID.Bytes()[:4]) {
			// matching error
			v, err2 := abierr.Unpack(errb)
			if err2 != nil {
				return fmt.Errorf("failed to unpack error: %w", err2)
			}
			return fmt.Errorf("error: %s, content %v", errName, v)
		}
	}
	return fmt.Errorf("error not found in abi: %s", hexErr)
}

func decodeErr(encodedABI string, err error) error {
	if err == nil {
		return nil
	}

	//revive:disable
	var d rpc.DataError
	ok := errors.As(err, &d)
	if ok {
		return parseContractErr(encodedABI, d.ErrorData().(string))
	}
	return fmt.Errorf("cannot decode err: %w", err)
}
