package keystone

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	v1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type DeployRequest struct {
	RegistryChain uint64
	Menv          deployment.MultiDonEnvironment

	DonToCapabilities map[string][]kcr.CapabilitiesRegistryCapability                   // from external source
	NodeIDToNop       map[string]capabilities_registry.CapabilitiesRegistryNodeOperator // maybe should be derivable from JD interface but doesn't seem to be notion of NOP in JD
}

type DeployResponse struct {
	AddressBook deployment.AddressBook
}

func Deploy(ctx context.Context, lggr logger.Logger, req DeployRequest) (*DeployResponse, error) {
	ad := deployment.NewMemoryAddressBook()
	resp := &DeployResponse{
		AddressBook: ad,
	}
	var registry *capabilities_registry.CapabilitiesRegistry
	var registryChain deployment.Chain

	for _, chain := range req.Menv.ListChains() {
		lggr.Info("deploying contracts", "chain", chain)
		deployResp, err := deployContracts(req.Menv.Logger, deployContractsRequest{
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

	donToOcr2Nodes, err := mapDonsToNodes(ctx, req.Menv)

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

	nopOnchainIDtoParams := make(map[uint32]capabilities_registry.CapabilitiesRegistryNodeParams)
	nodeIDToParams := make(map[string]capabilities_registry.CapabilitiesRegistryNodeParams)
	for don, ocr2nodes := range donToOcr2Nodes {
		caps, ok := capabilitiesResp.donToCapabilities[don]
		var hashedCapabilityIds [][32]byte
		for _, cap := range caps {
			hashedCapabilityIds = append(hashedCapabilityIds, cap.id)
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", hashedCapabilityIds)
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		for _, n := range ocr2nodes {
			nop, ok := nodeToRegisterNop[n.ID]
			if !ok {
				return nil, fmt.Errorf("node operator not found for node %s", n.ID)
			}
			params, ok := nopOnchainIDtoParams[nop.NodeOperatorId]
			// the first time we see a node operator, we need to fetch its p2p id and signer
			if !ok {
				params = capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      nop.NodeOperatorId,
					Signer:              n.Signer,
					P2pId:               n.P2PKey,
					HashedCapabilityIds: hashedCapabilityIds,
				}
			} else {
				// when we have a node operator, we need to dedup capabilities against the existing ones
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
			nopOnchainIDtoParams[nop.NodeOperatorId] = params
			nodeIDToParams[n.ID] = params
		}
	}
	lggr.Debugw("node params", "params", nopOnchainIDtoParams)
	var uniqueNodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, v := range nopOnchainIDtoParams {
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

	lggr.Infow("registered nodes", "nodes", nopOnchainIDtoParams)

	for don, ocr2nodes := range donToOcr2Nodes {
		var p2pIds [][32]byte
		for _, n := range ocr2nodes {
			params, ok := nodeIDToParams[n.ID]
			if !ok {
				return nil, fmt.Errorf("node params not found for node %s", n.ID)
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
			return nil, fmt.Errorf("failed to call AddDON for capability %v: %w", cfgs, err)
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
	donToCapabilities  map[string][]registeredCapability
	nodeIdToCapability map[string][]registeredCapability
}

type registeredCapability struct {
	capabilities_registry.CapabilitiesRegistryCapability
	id [32]byte
}

// func registerCapabilities(reg *capabilities_registry.CapabilitiesRegistry, chain deployment.Chain) error {
func registerCapabilities(lggr logger.Logger, req registerCapabilitiesRequest) (*registerCapabilitiesResponse, error) {
	if len(req.donToCapabilities) == 0 {
		return nil, fmt.Errorf("no capabilities to register")
	}
	resp := &registerCapabilitiesResponse{
		donToCapabilities:  make(map[string][]registeredCapability),
		nodeIdToCapability: make(map[string][]registeredCapability),
	}

	uniqueOffchainNodeIDs := make(map[string]struct{})
	for _, nodeIDs := range req.donToOffchainNodeIDs {
		for _, nodeID := range nodeIDs {
			uniqueOffchainNodeIDs[nodeID] = struct{}{}
		}
	}
	nodeIDs := make([]string, 0)
	for nodeID := range uniqueOffchainNodeIDs {
		nodeIDs = append(nodeIDs, nodeID)
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
			nodeCaps, ok := resp.nodeIdToCapability[nodeID]
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
			resp.nodeIdToCapability[nodeID] = append(resp.nodeIdToCapability[nodeID], newCaps...)
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

func nodeChainConfigsToOcr2Node(resp *v1.ListNodeChainConfigsResponse, id string) (*ocr2Node, error) {
	if len(resp.ChainConfigs) == 0 {
		return nil, errors.New("no chain configs")
	}
	cfg := resp.ChainConfigs[0]
	return newOcr2Node(id, cfg.Ocr2Config)
}

func mapDonsToNodes(ctx context.Context, menv deployment.MultiDonEnvironment) (map[string][]*ocr2Node, error) {
	// todo cleanup types
	donToOcr2Nodes := make(map[string][]*ocr2Node)
	for donName, env := range menv.DonToEnv {
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
				Filter: &v1.ListNodeChainConfigsRequest_Filter{NodeId: node.Id},
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
		menv.Logger.Infow("node chain configs", "don", donName, "configs", nodeCfgs)
	}
	return donToOcr2Nodes, nil
}
