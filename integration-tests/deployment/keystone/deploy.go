package keystone

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	v1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kf "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type DeployRequest struct {
	RegistryChainSel uint64
	Menv             deployment.MultiDonEnvironment

	DonToCapabilities map[string][]kcr.CapabilitiesRegistryCapability                   // from external source
	NodeIDToNop       map[string]capabilities_registry.CapabilitiesRegistryNodeOperator // maybe should be derivable from JD interface but doesn't seem to be notion of NOP in JD
}

type DeployResponse struct {
	Changeset *deployment.ChangesetOutput
	DonInfos  map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

type deployedContracts struct {
	chain deployment.Chain
	*deployerSet
}

func Deploy(ctx context.Context, lggr logger.Logger, req DeployRequest) (*DeployResponse, error) {
	ad := deployment.NewMemoryAddressBook()
	resp := &DeployResponse{
		Changeset: &deployment.ChangesetOutput{
			AddressBook: ad,
		},
		//DonToId: make(map[string]uint32),
		DonInfos: make(map[string]capabilities_registry.CapabilitiesRegistryDONInfo),
	}

	// deploy contracts on all chains and track the registry and ocr3 contracts
	deployedBySel := make(map[uint64]*deployedContracts)
	for _, chain := range req.Menv.ListChains() {
		lggr.Info("deploying contracts", "chain", chain)
		deployResp, err := deployContracts(req.Menv.Logger, deployContractsRequest{
			chain:           chain,
			isRegistryChain: chain.Selector == req.RegistryChainSel,
		},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contracts: %w", err)
		}
		err = ad.Merge(deployResp.AddressBook)
		if err != nil {
			return nil, fmt.Errorf("failed to merge address book: %w", err)
		}
		deployedBySel[chain.Selector] = &deployedContracts{
			chain:       chain,
			deployerSet: deployResp.deployerSet,
		}
	}

	// ensure registry is deployed and get the registry contract and chain
	var registry *capabilities_registry.CapabilitiesRegistry
	var registryChain deployment.Chain
	registryChainContracts, ok := deployedBySel[req.RegistryChainSel]
	if !ok {
		var got = []uint64{}
		for k := range req.Menv.Chains() {
			got = append(got, k)
		}
		return nil, fmt.Errorf("failed to deploy registry chain contracts. expected chain %d in %v", req.RegistryChainSel, got)
	}
	registry = registryChainContracts.capabilitiesRegistryDeployer.contract
	if registry == nil {
		return nil, fmt.Errorf("no registry contract found")
	}
	registryChain = registryChainContracts.chain
	lggr.Debugf("registry contract address: %s, chain %d", registry.Address().String(), registryChain.Selector)

	// all the subsequent calls to the registry are in terms of nodes
	// compute the mapping of dons to their nodes for reuse in various registry calls
	donToOcr2Nodes, err := mapDonsToNodes(ctx, req.Menv, true)
	if err != nil {
		return nil, fmt.Errorf("failed to map dons to nodes: %w", err)
	}

	// register capabilities
	capabilitiesResp, err := registerCapabilities(lggr, registerCapabilitiesRequest{
		chain:             registryChain,
		registry:          registry,
		donToCapabilities: req.DonToCapabilities,
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

	// register nodes
	nodesResp, err := registerNodes(lggr, &registerNodesRequest{
		registry:          registry,
		chain:             registryChain,
		nodeIdToNop:       req.NodeIDToNop,
		donToOcr2Nodes:    donToOcr2Nodes,
		donToCapabilities: capabilitiesResp.donToCapabilities,
		nops:              nopsResp.nops,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register nodes: %w", err)
	}
	lggr.Infow("registered nodes", "nodes", nodesResp.nodeIDToParams)

	// register DONS
	donsResp, err := registerDons(lggr, registerDonsRequest{
		registry:          registry,
		chain:             registryChain,
		nodeIDToParams:    nodesResp.nodeIDToParams,
		donToCapabilities: capabilitiesResp.donToCapabilities,
		donToOcr2Nodes:    donToOcr2Nodes,
	})
	lggr.Infow("registered DONS", "dons", donsResp.donToId)
	//resp.DonToId = donsResp.donToId
	resp.DonInfos = donsResp.donInfos

	// now we have the capability registry set up we need to configure the forwarder contracts and the OCR3 contract
	dons, err := assimalateToDons(resp.DonInfos, donToOcr2Nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to assimilate registry to Dons: %w", err)
	}
	// TODO this should be loaded from onchain state and passed in. as it is, this will fail if re-run on a real chain
	var fixedConfigVersion = uint32(1) // config version is fixed for to 1 when setting up the forwarder the first time
	for _, contracts := range deployedBySel {
		forwarder := contracts.keystoneForwarderDeployer.contract
		if forwarder == nil {
			return nil, fmt.Errorf("no forwarder contract found for chain %d", contracts.chain.Selector)
		}
		// configure the forwarder for the WF dons
		for _, dn := range dons {
			if !dn.info.AcceptsWorkflows {
				continue
			}
			tx, err := forwarder.SetConfig(contracts.chain.DeployerKey, dn.info.Id, fixedConfigVersion, dn.info.F, dn.signers())
			if err != nil {
				err = DecodeErr(kf.KeystoneForwarderABI, err)
				return nil, fmt.Errorf("failed to call SetConfig for forwarder %s on chain %d: %w", forwarder.Address().String(), contracts.chain.Selector, err)
			}
			_, err = contracts.chain.Confirm(tx.Hash())
			if err != nil {
				err = DecodeErr(kf.KeystoneForwarderABI, err)
				return nil, fmt.Errorf("failed to confirm SetConfig for forwarder %s: %w", forwarder.Address().String(), err)
			}
		}
	}
	// ocr3 contract on the registry chain
	return resp, err
}

type registerCapabilitiesRequest struct {
	chain             deployment.Chain
	registry          *capabilities_registry.CapabilitiesRegistry
	donToCapabilities map[string][]kcr.CapabilitiesRegistryCapability
}

type registerCapabilitiesResponse struct {
	donToCapabilities map[string][]registeredCapability
}

type registeredCapability struct {
	capabilities_registry.CapabilitiesRegistryCapability
	id [32]byte
}

// registerCapabilities add computes the capability id, adds it to the registry and associates the registered capabilities with appropriate don(s)
func registerCapabilities(lggr logger.Logger, req registerCapabilitiesRequest) (*registerCapabilitiesResponse, error) {
	if len(req.donToCapabilities) == 0 {
		return nil, fmt.Errorf("no capabilities to register")
	}
	resp := &registerCapabilitiesResponse{
		donToCapabilities: make(map[string][]registeredCapability),
	}

	// capability could be hosted on multiple dons. need to deduplicate
	uniqueCaps := make(map[kcr.CapabilitiesRegistryCapability][32]byte)
	for don, caps := range req.donToCapabilities {
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
	nops     []capabilities_registry.CapabilitiesRegistryNodeOperator
}

type registerNOPSResponse struct {
	nops []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
}

func registerNOPS(ctx context.Context, req registerNOPSRequest) (*registerNOPSResponse, error) {
	nops := req.nops
	tx, err := req.registry.AddNodeOperators(req.chain.DeployerKey, nops)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
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

// encodedABI is the abi of the contract that emitted the error
// hexErrBytes is the hex encoded error data
// returns an error that is the decoded error
func parseContractErr(encodedABI string, hexErr string) error {
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

func DecodeErr(encodedABI string, err error) error {
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
func nodeChainConfigsToOcr2Node(resp *v1.ListNodeChainConfigsResponse, id string) (*ocr2Node, error) {
	if len(resp.ChainConfigs) == 0 {
		return nil, errors.New("no chain configs")
	}
	cfg := resp.ChainConfigs[0]
	return newOcr2Node(id, cfg.Ocr2Config)
}

// mapDonsToNodes returns a map of don name to simplified representation of their nodes
func mapDonsToNodes(ctx context.Context, menv deployment.MultiDonEnvironment, excludeBootstraps bool) (map[string][]*ocr2Node, error) {
	donToOcr2Nodes := make(map[string][]*ocr2Node)
	for donName, env := range menv.DonToEnv {
		donNodeSet, err := env.Offchain.ListNodes(ctx, &v1.ListNodesRequest{}) // each env is a don
		if err != nil {
			return nil, fmt.Errorf("failed to list nodes: %w", err)
		}
		if len(donNodeSet.Nodes) == 0 {
			return nil, fmt.Errorf("no nodes found")
		}
		// each node in the nodeset may support multiple chains
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
			if excludeBootstraps && ocr2n.IsBoostrap {
				continue
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

// register nodes
type registerNodesRequest struct {
	registry          *capabilities_registry.CapabilitiesRegistry
	chain             deployment.Chain
	nodeIdToNop       map[string]capabilities_registry.CapabilitiesRegistryNodeOperator
	donToOcr2Nodes    map[string][]*ocr2Node
	donToCapabilities map[string][]registeredCapability
	nops              []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
}
type registerNodesResponse struct {
	nodeIDToParams map[string]capabilities_registry.CapabilitiesRegistryNodeParams
}

func registerNodes(lggr logger.Logger, req *registerNodesRequest) (*registerNodesResponse, error) {
	nopToNodeIDs := make(map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string)
	for nodeID, nop := range req.nodeIdToNop {
		if _, ok := nopToNodeIDs[nop]; !ok {
			nopToNodeIDs[nop] = make([]string, 0)
		}
		nopToNodeIDs[nop] = append(nopToNodeIDs[nop], nodeID)
	}
	nodeToRegisterNop := make(map[string]*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded)
	for _, nop := range req.nops {
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

	nodeIDToParams := make(map[string]capabilities_registry.CapabilitiesRegistryNodeParams)
	for don, ocr2nodes := range req.donToOcr2Nodes {
		caps, ok := req.donToCapabilities[don]
		var hashedCapabilityIds [][32]byte
		for _, cap := range caps {
			hashedCapabilityIds = append(hashedCapabilityIds, cap.id)
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", hashedCapabilityIds)
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		for _, n := range ocr2nodes {
			if n.IsBoostrap { // bootstraps are part of the DON but don't host capabilities
				continue
			}
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
			nodeIDToParams[n.ID] = params
		}
	}
	lggr.Debugw("node params", "params", nodeIDToParams)
	var uniqueNodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, v := range nodeIDToParams {
		uniqueNodeParams = append(uniqueNodeParams, v)
	}
	lggr.Debug("unique node params", "params", uniqueNodeParams)
	tx, err := req.registry.AddNodes(req.chain.DeployerKey, uniqueNodeParams)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNode: %w", err)
	}
	_, err = req.chain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash(), err)
	}
	return &registerNodesResponse{
		nodeIDToParams: nodeIDToParams,
	}, nil
}

type registerDonsRequest struct {
	registry *capabilities_registry.CapabilitiesRegistry
	chain    deployment.Chain

	nodeIDToParams    map[string]capabilities_registry.CapabilitiesRegistryNodeParams
	donToCapabilities map[string][]registeredCapability
	donToOcr2Nodes    map[string][]*ocr2Node
}

type registerDonsResponse struct {
	donToId  map[string]uint32
	donInfos map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

func id(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func registerDons(lggr logger.Logger, req registerDonsRequest) (*registerDonsResponse, error) {
	resp := &registerDonsResponse{
		donToId:  make(map[string]uint32),
		donInfos: make(map[string]capabilities_registry.CapabilitiesRegistryDONInfo),
	}
	//	p2pIdsToDon := make(map[[][32]byte]string)
	p2pIdsToDon := make(map[string]string)

	donid := uint32(1) //TODO: this should be loaded from onchain state and passed in. as it is, this will fail if re-run on a real chain
	for don, ocr2nodes := range req.donToOcr2Nodes {
		var p2pIds [][32]byte
		for _, n := range ocr2nodes {
			if n.IsBoostrap {
				continue
			}
			params, ok := req.nodeIDToParams[n.ID]
			if !ok {
				return nil, fmt.Errorf("node params not found for non-bootstrap node %s", n.ID)
			}
			p2pIds = append(p2pIds, params.P2pId)
		}

		p2pIdsToDon[id(p2pIds)] = don
		caps, ok := req.donToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		wfSupported := false
		var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
		for _, cap := range caps {
			if cap.CapabilityType == 2 { // OCR3 capability => WF supported
				wfSupported = true
			}
			// TODO: accept configuration from external source for each (don,capability)
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
		tx, err := req.registry.AddDON(req.chain.DeployerKey, p2pIds, cfgs, true, wfSupported, uint8(f))
		if err != nil {
			err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
			return nil, fmt.Errorf("failed to call AddDON for don %s p2p2Ids %v capability %v: %w", don, p2pIds, cfgs, err)
		}
		_, err = req.chain.Confirm(tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddDON transaction %s for don %s: %w", tx.Hash(), don, err)
		}
		resp.donToId[don] = donid
		lggr.Debugw("registered DON", "don", don, "p2pids", p2pIds, "cgs", cfgs, "wfSupported", wfSupported, "f", f, "id", donid)
		donid++
	}
	donInfos, err := req.registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	for i, donInfo := range donInfos {
		donName, ok := p2pIdsToDon[id(donInfo.NodeP2PIds)]
		if !ok {
			return nil, fmt.Errorf("don not found for p2pids %s in %v", id(donInfo.NodeP2PIds), p2pIdsToDon)
		}
		resp.donInfos[donName] = donInfos[i]
	}
	return resp, nil
}

func assimalateToDons(donInfos map[string]kcr.CapabilitiesRegistryDONInfo, nodes map[string][]*ocr2Node) ([]don, error) {
	// all maps should have the same keys
	if len(donInfos) != len(nodes) {
		return nil, fmt.Errorf("mismatched lengths don infos %d,  nodes %d", len(donInfos), len(nodes))
	}
	var out []don
	for donName, info := range donInfos {

		ocr2nodes, ok := nodes[donName]
		if !ok {
			return nil, fmt.Errorf("nodes not found for don %s", donName)
		}
		var ocr2ns []*ocr2Node
		for _, n := range ocr2nodes {
			ocr2ns = append(ocr2ns, n)
		}
		out = append(out, don{
			name:  donName,
			info:  info,
			nodes: ocr2ns,
		})
	}

	return out, nil
}

type don struct {
	name string
	/*
		registryId   uint32
		capabilities []registeredCapability
	*/
	info   capabilities_registry.CapabilitiesRegistryDONInfo
	nodes  []*ocr2Node
	config string //TODO
}

/*
func (d don) isWf() bool {
	return d.info.AcceptsWorkflows
}

func (d don) f() uint32 {
	// assumes n = 3f+1
	return uint32(len(d.nodes) / 3)
}
*/

func (d don) signers() []common.Address {
	sort.Slice(d.nodes, func(i, j int) bool {
		return d.nodes[i].P2PKey.String() < d.nodes[j].P2PKey.String()
	})
	var out []common.Address
	for _, n := range d.nodes {
		if n.IsBoostrap {
			continue
		}
		out = append(out, n.address())
	}
	return out
}
