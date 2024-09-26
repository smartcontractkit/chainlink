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
	kocr3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ksscript "github.com/smartcontractkit/chainlink/core/scripts/keystone/src"
)

type ConfigureContractsRequest struct {
	RegistryChainSel uint64
	Env              *deployment.Environment

	Dons       []DonCapabilities            // externally sourced based on the environment
	OCR3Config *ksscript.OracleConfigSource // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config

	AddressBook      deployment.AddressBook
	DoContractDeploy bool // if false, the contracts are assumed to be deployed and the address book is used
}

type ConfigureContractsResponse struct {
	Changeset *deployment.ChangesetOutput
	DonInfos  map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

// ConfigureContracts configures contracts them with the given DONS and their capabilities. It optionally deploys the contracts
// but best practice is to deploy them separately and pass the address book in the request
func ConfigureContracts(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest) (*ConfigureContractsResponse, error) {
	if req.OCR3Config == nil {
		return nil, errors.New("OCR3Config is nil")
	}
	_, ok := req.Env.Chains[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found in environment", req.RegistryChainSel)
	}

	addrBook := req.AddressBook
	if req.DoContractDeploy {
		contractDeployCS, err := DeployContracts(lggr, req.Env, req.RegistryChainSel)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contracts: %w", err)
		}
		addrBook = contractDeployCS.AddressBook
	} else {
		lggr.Debug("skipping contract deployment")
	}
	if addrBook == nil {
		return nil, errors.New("address book is nil")
	}

	cfgRegistryResp, err := ConfigureRegistry(ctx, lggr, req, addrBook)
	if err != nil {
		return nil, fmt.Errorf("failed to configure registry: %w", err)
	}

	// now we have the capability registry set up we need to configure the forwarder contracts and the OCR3 contract
	dons, err := joinInfoAndNodes(cfgRegistryResp.DonInfos, req.Dons)
	if err != nil {
		return nil, fmt.Errorf("failed to assimilate registry to Dons: %w", err)
	}
	err = ConfigureForwardContracts(req.Env, dons, addrBook)
	if err != nil {
		return nil, fmt.Errorf("failed to configure forwarder contracts: %w", err)
	}

	err = ConfigureOCR3Contract(req.Env, req.RegistryChainSel, dons, addrBook, req.OCR3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OCR3 contract: %w", err)
	}

	return &ConfigureContractsResponse{
		Changeset: &deployment.ChangesetOutput{
			AddressBook: addrBook,
		},
		DonInfos: cfgRegistryResp.DonInfos,
	}, nil
}

// DeployContracts deploys the all the keystone contracts on all chains and returns the address book in the changeset
func DeployContracts(lggr logger.Logger, e *deployment.Environment, chainSel uint64) (*deployment.ChangesetOutput, error) {
	adbook := deployment.NewMemoryAddressBook()
	// deploy contracts on all chains and track the registry and ocr3 contracts
	for _, chain := range e.Chains {
		lggr.Infow("deploying contracts", "chain", chain)
		deployResp, err := deployContracts(lggr, deployContractsRequest{
			chain:           chain,
			isRegistryChain: chain.Selector == chainSel,
		},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contracts: %w", err)
		}
		err = adbook.Merge(deployResp.AddressBook)
		if err != nil {
			return nil, fmt.Errorf("failed to merge address book: %w", err)
		}
	}
	return &deployment.ChangesetOutput{
		AddressBook: adbook,
	}, nil
}

// ConfigureRegistry configures the registry contract with the given DONS and their capabilities
// the address book is required to contain the addresses of the deployed registry contract
func ConfigureRegistry(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest, addrBook deployment.AddressBook) (*ConfigureContractsResponse, error) {
	registryChain, ok := req.Env.Chains[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found in environment", req.RegistryChainSel)
	}

	contractSetsResp, err := GetContractSets(&GetContractSetsRequest{
		Chains:      req.Env.Chains,
		AddressBook: addrBook,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	// ensure registry is deployed and get the registry contract and chain
	var registry *capabilities_registry.CapabilitiesRegistry
	registryChainContracts, ok := contractSetsResp.ContractSets[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("failed to deploy registry chain contracts. expected chain %d", req.RegistryChainSel)
	}
	registry = registryChainContracts.CapabilitiesRegistry
	if registry == nil {
		return nil, fmt.Errorf("no registry contract found")
	}
	lggr.Debugf("registry contract address: %s, chain %d", registry.Address().String(), req.RegistryChainSel)

	// all the subsequent calls to the registry are in terms of nodes
	// compute the mapping of dons to their nodes for reuse in various registry calls
	donToOcr2Nodes, err := mapDonsToNodes(req.Dons, true)
	if err != nil {
		return nil, fmt.Errorf("failed to map dons to nodes: %w", err)
	}

	// TODO: we can remove this abstractions and refactor the functions that accept them to accept []DonCapabilities
	// they are unnecessary indirection
	donToCapabilities := mapDonsToCaps(req.Dons)
	nodeIdToNop, err := nodesToNops(req.Dons, req.RegistryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to map nodes to nops: %w", err)
	}

	// register capabilities
	capabilitiesResp, err := registerCapabilities(lggr, registerCapabilitiesRequest{
		chain:             registryChain,
		registry:          registry,
		donToCapabilities: donToCapabilities,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register capabilities: %w", err)
	}
	lggr.Infow("registered capabilities", "capabilities", capabilitiesResp.donToCapabilities)

	// register node operators
	var nops []capabilities_registry.CapabilitiesRegistryNodeOperator
	for _, nop := range nodeIdToNop {
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
		nodeIdToNop:       nodeIdToNop,
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
	if err != nil {
		return nil, fmt.Errorf("failed to register DONS: %w", err)
	}
	lggr.Infow("registered DONS", "dons", len(donsResp.donInfos))

	return &ConfigureContractsResponse{
		Changeset: &deployment.ChangesetOutput{
			AddressBook: addrBook,
		},
		DonInfos: donsResp.donInfos,
	}, nil
}

// ConfigureForwardContracts configures the forwarder contracts on all chains for the given DONS
// the address book is required to contain the an address of the deployed forwarder contract for every chain in the environment
func ConfigureForwardContracts(env *deployment.Environment, dons []RegisteredDon, addrBook deployment.AddressBook) error {
	contractSetsResp, err := GetContractSets(&GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: addrBook,
	})
	if err != nil {
		return fmt.Errorf("failed to get contract sets: %w", err)
	}

	// configure forwarders on all chains
	for _, chain := range env.Chains {
		// get the forwarder contract for the chain
		contracts, ok := contractSetsResp.ContractSets[chain.Selector]
		if !ok {
			return fmt.Errorf("failed to get contract set for chain %d", chain.Selector)
		}
		fwrd := contracts.Forwarder
		if fwrd == nil {
			return fmt.Errorf("no forwarder contract found for chain %d", chain.Selector)
		}

		err := configureForwarder(chain, fwrd, dons)
		if err != nil {
			return fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}
	}
	return nil
}

// ocr3 contract on the registry chain for the wf dons
func ConfigureOCR3Contract(env *deployment.Environment, chainSel uint64, dons []RegisteredDon, addrBook deployment.AddressBook, cfg *ksscript.OracleConfigSource) error {
	registryChain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", chainSel)
	}

	contractSetsResp, err := GetContractSets(&GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: addrBook,
	})
	if err != nil {
		return fmt.Errorf("failed to get contract sets: %w", err)
	}

	for _, don := range dons {
		if !don.Info.AcceptsWorkflows {
			continue
		}
		// only on the registry chain
		contracts, ok := contractSetsResp.ContractSets[chainSel]
		if !ok {
			return fmt.Errorf("failed to get contract set for chain %d", chainSel)
		}
		contract := contracts.OCR3
		if contract == nil {
			return fmt.Errorf("no forwarder contract found for chain %d", chainSel)
		}

		_, err := configureOCR3contract(configureOCR3Request{
			cfg:      cfg,
			chain:    registryChain,
			contract: contract,
			don:      don,
		})
		if err != nil {
			return fmt.Errorf("failed to configure OCR3 contract for don %s: %w", don.Name, err)
		}
	}
	return nil
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
	_, err = req.chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
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
	_, err = req.chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
	}

	receipt, err := bind.WaitMined(ctx, req.chain.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to mine AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
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
	return fmt.Errorf("not found in abi: %s", hexErr)
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
	return fmt.Errorf("cannot decode error with abi: %w", err)
}
func nodeChainConfigsToOcr2Node(resp *v1.ListNodeChainConfigsResponse, id string) (*ocr2Node, error) {
	if len(resp.ChainConfigs) == 0 {
		return nil, errors.New("no chain configs")
	}
	cfg := resp.ChainConfigs[0]
	return newOcr2Node(id, cfg)
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
	_, err = req.chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash().String(), err)
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
	donInfos map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

func sortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func registerDons(lggr logger.Logger, req registerDonsRequest) (*registerDonsResponse, error) {
	resp := &registerDonsResponse{
		donInfos: make(map[string]capabilities_registry.CapabilitiesRegistryDONInfo),
	}
	// track hash of sorted p2pids to don name because the registry return value does not include the don name
	// and we need to map it back to the don name to access the other mapping data such as the don's capabilities & nodes
	p2pIdsToDon := make(map[string]string)

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

		p2pSortedHash := sortedHash(p2pIds)
		p2pIdsToDon[p2pSortedHash] = don
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

		f := len(p2pIds) / 3 // assuming n=3f+1. TODO should come for some config.
		tx, err := req.registry.AddDON(req.chain.DeployerKey, p2pIds, cfgs, true, wfSupported, uint8(f))
		if err != nil {
			err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
			return nil, fmt.Errorf("failed to call AddDON for don '%s' p2p2Id hash %s capability %v: %w", don, p2pSortedHash, cfgs, err)
		}
		_, err = req.chain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddDON transaction %s for don %s: %w", tx.Hash().String(), don, err)
		}
		lggr.Debugw("registered DON", "don", don, "p2p sorted hash", p2pSortedHash, "cgs", cfgs, "wfSupported", wfSupported, "f", f)
	}
	donInfos, err := req.registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	for i, donInfo := range donInfos {
		donName, ok := p2pIdsToDon[sortedHash(donInfo.NodeP2PIds)]
		if !ok {
			return nil, fmt.Errorf("don not found for p2pids %s in %v", sortedHash(donInfo.NodeP2PIds), p2pIdsToDon)
		}
		resp.donInfos[donName] = donInfos[i]
	}
	return resp, nil
}

// configureForwarder sets the config for the forwarder contract on the chain for all Dons that accept workflows
// dons that don't accept workflows are not registered with the forwarder
func configureForwarder(chain deployment.Chain, fwdr *kf.KeystoneForwarder, dons []RegisteredDon) error {
	if fwdr == nil {
		return errors.New("nil forwarder contract")
	}
	for _, dn := range dons {
		if !dn.Info.AcceptsWorkflows {
			continue
		}
		ver := dn.Info.ConfigCount // note config count on the don info is the version on the forwarder
		tx, err := fwdr.SetConfig(chain.DeployerKey, dn.Info.Id, ver, dn.Info.F, dn.signers())
		if err != nil {
			err = DecodeErr(kf.KeystoneForwarderABI, err)
			return fmt.Errorf("failed to call SetConfig for forwarder %s on chain %d: %w", fwdr.Address().String(), chain.Selector, err)
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			err = DecodeErr(kf.KeystoneForwarderABI, err)
			return fmt.Errorf("failed to confirm SetConfig for forwarder %s: %w", fwdr.Address().String(), err)
		}
	}
	return nil
}

type configureOCR3Request struct {
	cfg      *ksscript.OracleConfigSource
	chain    deployment.Chain
	contract *kocr3.OCR3Capability
	don      RegisteredDon
}
type configureOCR3Response struct {
	ocrConfig ksscript.Orc2drOracleConfig
}

func configureOCR3contract(req configureOCR3Request) (*configureOCR3Response, error) {
	if req.contract == nil {
		return nil, fmt.Errorf("OCR3 contract is nil")
	}
	nks := makeNodeKeysSlice(req.don.Nodes)
	ocrConfig := ksscript.GenerateOCR3Config(*req.cfg, nks)
	/*ocrConfig, err := generateOCR3Config(*req.cfg, nks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}
	*/
	tx, err := req.contract.SetConfig(req.chain.DeployerKey,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	if err != nil {
		err = DecodeErr(kocr3.OCR3CapabilityABI, err)
		return nil, fmt.Errorf("failed to call SetConfig for OCR3 contract %s: %w", req.contract.Address().String(), err)
	}
	_, err = req.chain.Confirm(tx)
	if err != nil {
		err = DecodeErr(kocr3.OCR3CapabilityABI, err)
		return nil, fmt.Errorf("failed to confirm SetConfig for OCR3 contract %s: %w", req.contract.Address().String(), err)
	}
	return &configureOCR3Response{ocrConfig}, nil
}
