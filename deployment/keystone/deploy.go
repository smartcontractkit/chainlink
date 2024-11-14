package keystone

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/exp/maps"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	chainsel "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kf "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	kocr3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type ConfigureContractsRequest struct {
	RegistryChainSel uint64
	Env              *deployment.Environment

	Dons       []DonCapabilities        // externally sourced based on the environment
	OCR3Config *OracleConfigWithSecrets // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config

	DoContractDeploy bool // if false, the contracts are assumed to be deployed and the address book is used
}

func (r ConfigureContractsRequest) Validate() error {
	if r.OCR3Config == nil {
		return errors.New("OCR3Config is nil")
	}
	if r.Env == nil {
		return errors.New("environment is nil")
	}
	for _, don := range r.Dons {
		if err := don.Validate(); err != nil {
			return fmt.Errorf("don validation failed for '%s': %w", don.Name, err)
		}
	}
	_, ok := chainsel.ChainBySelector(r.RegistryChainSel)
	if !ok {
		return fmt.Errorf("chain %d not found in environment", r.RegistryChainSel)
	}
	return nil
}

type ConfigureContractsResponse struct {
	Changeset *deployment.ChangesetOutput
	DonInfos  map[string]kcr.CapabilitiesRegistryDONInfo
}

// ConfigureContracts configures contracts them with the given DONS and their capabilities. It optionally deploys the contracts
// but best practice is to deploy them separately and pass the address book in the request
func ConfigureContracts(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest) (*ConfigureContractsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	addrBook := req.Env.ExistingAddresses
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

	donInfos, err := DonInfos(req.Dons, req.Env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get don infos: %w", err)
	}

	// now we have the capability registry set up we need to configure the forwarder contracts and the OCR3 contract
	dons, err := joinInfoAndNodes(cfgRegistryResp.DonInfos, donInfos, req.RegistryChainSel)
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
		lggr.Infow("deploying contracts", "chain", chain.Selector)
		deployResp, err := deployContractsToChain(lggr, deployContractsRequest{
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

// DonInfo is DonCapabilities, but expanded to contain node information
type DonInfo struct {
	Name         string
	Nodes        []Node
	Capabilities []kcr.CapabilitiesRegistryCapability // every capability is hosted on each node
}

// TODO: merge with deployment/environment.go Node
type Node struct {
	ID           string
	P2PID        string
	Name         string
	PublicKey    *string
	ChainConfigs []*nodev1.ChainConfig
}

// TODO: merge with deployment/environment.go NodeInfo, we currently lookup based on p2p_id, and chain-selectors needs non-EVM support
func NodesFromJD(name string, nodeIDs []string, jd deployment.OffchainClient) ([]Node, error) {
	// lookup nodes based on p2p_ids
	var nodes []Node
	selector := strings.Join(nodeIDs, ",")
	nodesFromJD, err := jd.ListNodes(context.Background(), &nodev1.ListNodesRequest{
		Filter: &nodev1.ListNodesRequest_Filter{
			Enabled: 1,
			Selectors: []*ptypes.Selector{
				{
					Key:   "p2p_id",
					Op:    ptypes.SelectorOp_IN,
					Value: &selector,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes '%s': %w", name, err)
	}

	for _, id := range nodeIDs {
		idx := slices.IndexFunc(nodesFromJD.GetNodes(), func(node *nodev1.Node) bool {
			return slices.ContainsFunc(node.Labels, func(label *ptypes.Label) bool {
				return label.Key == "p2p_id" && *label.Value == id
			})
		})
		if idx < 0 {
			var got []string
			for _, node := range nodesFromJD.GetNodes() {
				for _, label := range node.Labels {
					if label.Key == "p2p_id" {
						got = append(got, *label.Value)
					}
				}
			}
			return nil, fmt.Errorf("node id %s not found in list '%s'", id, strings.Join(got, ","))
		}

		jdNode := nodesFromJD.Nodes[idx]
		// TODO: Filter should accept multiple nodes
		nodeChainConfigs, err := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{jdNode.Id}, // must use the jd-specific internal node id
		}})
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, Node{
			ID:           jdNode.Id,
			P2PID:        id,
			Name:         name,
			PublicKey:    &jdNode.PublicKey,
			ChainConfigs: nodeChainConfigs.GetChainConfigs(),
		})
	}
	return nodes, nil
}

func DonInfos(dons []DonCapabilities, jd deployment.OffchainClient) ([]DonInfo, error) {
	var donInfos []DonInfo
	for _, don := range dons {
		var nodeIDs []string
		for _, nop := range don.Nops {
			nodeIDs = append(nodeIDs, nop.Nodes...)
		}
		nodes, err := NodesFromJD(don.Name, nodeIDs, jd)
		if err != nil {
			return nil, err
		}
		donInfos = append(donInfos, DonInfo{
			Name:         don.Name,
			Nodes:        nodes,
			Capabilities: don.Capabilities,
		})
	}
	return donInfos, nil
}

// ConfigureRegistry configures the registry contract with the given DONS and their capabilities
// the address book is required to contain the addresses of the deployed registry contract
func ConfigureRegistry(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest, addrBook deployment.AddressBook) (*ConfigureContractsResponse, error) {
	registryChain, ok := req.Env.Chains[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found in environment", req.RegistryChainSel)
	}

	contractSetsResp, err := GetContractSets(req.Env.Logger, &GetContractSetsRequest{
		Chains:      req.Env.Chains,
		AddressBook: addrBook,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	donInfos, err := DonInfos(req.Dons, req.Env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get don infos: %w", err)
	}

	// ensure registry is deployed and get the registry contract and chain
	var registry *kcr.CapabilitiesRegistry
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
	donToOcr2Nodes, err := mapDonsToNodes(donInfos, true, req.RegistryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to map dons to nodes: %w", err)
	}

	// TODO: we can remove this abstractions and refactor the functions that accept them to accept []DonInfos/DonCapabilities
	// they are unnecessary indirection
	donToCapabilities := mapDonsToCaps(donInfos)
	nopsToNodeIDs, err := nopsToNodes(donInfos, req.Dons, req.RegistryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to map nops to nodes: %w", err)
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
	nopsList := maps.Keys(nopsToNodeIDs)
	nopsResp, err := RegisterNOPS(ctx, lggr, RegisterNOPSRequest{
		Chain:    registryChain,
		Registry: registry,
		Nops:     nopsList,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register node operators: %w", err)
	}
	lggr.Infow("registered node operators", "nops", nopsResp.Nops)

	// register nodes
	nodesResp, err := registerNodes(lggr, &registerNodesRequest{
		registry:          registry,
		chain:             registryChain,
		nopToNodeIDs:      nopsToNodeIDs,
		donToOcr2Nodes:    donToOcr2Nodes,
		donToCapabilities: capabilitiesResp.donToCapabilities,
		nops:              nopsResp.Nops,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register nodes: %w", err)
	}
	lggr.Infow("registered nodes", "nodes", nodesResp.nodeIDToParams)

	// TODO: annotate nodes with node_operator_id in JD?

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
	lggr.Infow("registered DONs", "dons", len(donsResp.donInfos))

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
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
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

		err := configureForwarder(env.Logger, chain, fwrd, dons)
		if err != nil {
			return fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}
	}
	return nil
}

// ocr3 contract on the registry chain for the wf dons
func ConfigureOCR3Contract(env *deployment.Environment, chainSel uint64, dons []RegisteredDon, addrBook deployment.AddressBook, cfg *OracleConfigWithSecrets) error {
	registryChain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", chainSel)
	}

	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
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
			return fmt.Errorf("no ocr3 contract found for chain %d", chainSel)
		}

		_, err := configureOCR3contract(configureOCR3Request{
			cfg:      cfg,
			chain:    registryChain,
			contract: contract,
			nodes:    don.Nodes,
		})
		if err != nil {
			return fmt.Errorf("failed to configure OCR3 contract for don %s: %w", don.Name, err)
		}
	}
	return nil
}

func ConfigureOCR3ContractFromJD(env *deployment.Environment, chainSel uint64, nodeIDs []string, addrBook deployment.AddressBook, cfg *OracleConfigWithSecrets) error {
	registryChain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", chainSel)
	}
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: addrBook,
	})
	if err != nil {
		return fmt.Errorf("failed to get contract sets: %w", err)
	}
	contracts, ok := contractSetsResp.ContractSets[chainSel]
	if !ok {
		return fmt.Errorf("failed to get contract set for chain %d", chainSel)
	}
	contract := contracts.OCR3
	if contract == nil {
		return fmt.Errorf("no ocr3 contract found for chain %d", chainSel)
	}
	nodes, err := NodesFromJD("nodes", nodeIDs, env.Offchain)
	if err != nil {
		return err
	}
	var ocr2nodes []*ocr2Node
	for _, node := range nodes {
		n, err := newOcr2NodeFromJD(&node, chainSel)
		if err != nil {
			return fmt.Errorf("failed to create ocr2 node from clo node: %w", err)
		}
		ocr2nodes = append(ocr2nodes, n)
	}
	_, err = configureOCR3contract(configureOCR3Request{
		cfg:      cfg,
		chain:    registryChain,
		contract: contract,
		nodes:    ocr2nodes,
	})
	return err
}

type registerCapabilitiesRequest struct {
	chain             deployment.Chain
	registry          *kcr.CapabilitiesRegistry
	donToCapabilities map[string][]kcr.CapabilitiesRegistryCapability
}

type registerCapabilitiesResponse struct {
	donToCapabilities map[string][]RegisteredCapability
}

type RegisteredCapability struct {
	kcr.CapabilitiesRegistryCapability
	ID [32]byte
}

// registerCapabilities add computes the capability id, adds it to the registry and associates the registered capabilities with appropriate don(s)
func registerCapabilities(lggr logger.Logger, req registerCapabilitiesRequest) (*registerCapabilitiesResponse, error) {
	if len(req.donToCapabilities) == 0 {
		return nil, fmt.Errorf("no capabilities to register")
	}
	lggr.Infow("registering capabilities...", "len", len(req.donToCapabilities))
	resp := &registerCapabilitiesResponse{
		donToCapabilities: make(map[string][]RegisteredCapability),
	}

	// capability could be hosted on multiple dons. need to deduplicate
	uniqueCaps := make(map[kcr.CapabilitiesRegistryCapability][32]byte)
	for don, caps := range req.donToCapabilities {
		var registerCaps []RegisteredCapability
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
			registerCap := RegisteredCapability{
				CapabilitiesRegistryCapability: cap,
				ID:                             id,
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

	err := AddCapabilities(lggr, req.registry, req.chain, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}
	return resp, nil
}

type RegisterNOPSRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry
	Nops     []kcr.CapabilitiesRegistryNodeOperator
}

type RegisterNOPSResponse struct {
	Nops []*kcr.CapabilitiesRegistryNodeOperatorAdded
}

func RegisterNOPS(ctx context.Context, lggr logger.Logger, req RegisterNOPSRequest) (*RegisterNOPSResponse, error) {
	lggr.Infow("registering node operators...", "len", len(req.Nops))
	existingNops, err := req.Registry.GetNodeOperators(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	existingNopsAddrToID := make(map[capabilities_registry.CapabilitiesRegistryNodeOperator]uint32)
	for id, nop := range existingNops {
		existingNopsAddrToID[nop] = uint32(id)
	}
	lggr.Infow("fetched existing node operators", "len", len(existingNopsAddrToID))
	resp := &RegisterNOPSResponse{
		Nops: []*kcr.CapabilitiesRegistryNodeOperatorAdded{},
	}
	nops := []kcr.CapabilitiesRegistryNodeOperator{}
	for _, nop := range req.Nops {
		if id, ok := existingNopsAddrToID[nop]; !ok {
			nops = append(nops, nop)
		} else {
			lggr.Debugw("node operator already exists", "name", nop.Name, "admin", nop.Admin.String(), "id", id)
			resp.Nops = append(resp.Nops, &kcr.CapabilitiesRegistryNodeOperatorAdded{
				NodeOperatorId: id,
				Name:           nop.Name,
				Admin:          nop.Admin,
			})
		}
	}
	if len(nops) == 0 {
		lggr.Debug("no new node operators to register")
		return resp, nil
	}
	tx, err := req.Registry.AddNodeOperators(req.Chain.DeployerKey, nops)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
	}
	// for some reason that i don't understand, the confirm must be called before the WaitMined or the latter will hang
	// (at least for a simulated backend chain)
	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
	}

	receipt, err := bind.WaitMined(ctx, req.Chain.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to mine AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
	}
	if len(receipt.Logs) != len(nops) {
		return nil, fmt.Errorf("expected %d log entries for AddNodeOperators, got %d", len(nops), len(receipt.Logs))
	}
	for i, log := range receipt.Logs {
		o, err := req.Registry.ParseNodeOperatorAdded(*log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log %d for operator added: %w", i, err)
		}
		resp.Nops = append(resp.Nops, o)
	}

	return resp, nil
}

func DefaultCapConfig(capType uint8, nNodes int) *capabilitiespb.CapabilityConfig {
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

func DecodeErr(encodedABI string, err error) error {
	if err == nil {
		return nil
	}

	//revive:disable
	var d rpc.DataError
	ok := errors.As(err, &d)
	if ok {
		encErr, ok := d.ErrorData().(string)
		if !ok {
			return fmt.Errorf("error without error data: %s", d.Error())
		}
		errStr, parseErr := deployment.ParseErrorFromABI(encErr, encodedABI)
		if parseErr != nil {
			return fmt.Errorf("failed to decode error '%s' with abi: %w", encErr, parseErr)
		}
		return fmt.Errorf("contract error: %s", errStr)

	}
	return fmt.Errorf("cannot decode error with abi: %w", err)
}

// register nodes
type registerNodesRequest struct {
	registry          *kcr.CapabilitiesRegistry
	chain             deployment.Chain
	nopToNodeIDs      map[kcr.CapabilitiesRegistryNodeOperator][]string
	donToOcr2Nodes    map[string][]*ocr2Node
	donToCapabilities map[string][]RegisteredCapability
	nops              []*kcr.CapabilitiesRegistryNodeOperatorAdded
}
type registerNodesResponse struct {
	nodeIDToParams map[string]kcr.CapabilitiesRegistryNodeParams
}

// registerNodes registers the nodes with the registry. it assumes that the deployer key in the Chain
// can sign the transactions update the contract state
// TODO: 467 refactor to support MCMS. Specifically need to separate the call data generation from the actual contract call
func registerNodes(lggr logger.Logger, req *registerNodesRequest) (*registerNodesResponse, error) {
	var count int
	for _, nodes := range req.nopToNodeIDs {
		count += len(nodes)
	}
	lggr.Infow("registering nodes...", "len", count)
	nodeToRegisterNop := make(map[string]*kcr.CapabilitiesRegistryNodeOperatorAdded)
	for _, nop := range req.nops {
		n := kcr.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: nop.Admin,
		}
		nodeIDs := req.nopToNodeIDs[n]
		for _, nodeID := range nodeIDs {
			_, exists := nodeToRegisterNop[nodeID]
			if !exists {
				nodeToRegisterNop[nodeID] = nop
			}
		}
	}

	nodeIDToParams := make(map[string]kcr.CapabilitiesRegistryNodeParams)
	for don, ocr2nodes := range req.donToOcr2Nodes {
		caps, ok := req.donToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		var hashedCapabilityIds [][32]byte
		for _, cap := range caps {
			hashedCapabilityIds = append(hashedCapabilityIds, cap.ID)
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", hashedCapabilityIds)

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
				params = kcr.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      nop.NodeOperatorId,
					Signer:              n.Signer,
					P2pId:               n.P2PKey,
					EncryptionPublicKey: n.EncryptionPublicKey,
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

	var uniqueNodeParams []kcr.CapabilitiesRegistryNodeParams
	for _, v := range nodeIDToParams {
		uniqueNodeParams = append(uniqueNodeParams, v)
	}
	lggr.Debugw("unique node params to add", "count", len(uniqueNodeParams))
	tx, err := req.registry.AddNodes(req.chain.DeployerKey, uniqueNodeParams)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		// no typed errors in the abi, so we have to do string matching
		// try to add all nodes in one go, if that fails, fall back to 1-by-1
		if !strings.Contains(err.Error(), "NodeAlreadyExists") {
			return nil, fmt.Errorf("failed to call AddNodes for bulk add nodes: %w", err)
		}
		lggr.Warn("nodes already exist, falling back to 1-by-1")
		for _, singleNodeParams := range uniqueNodeParams {
			tx, err = req.registry.AddNodes(req.chain.DeployerKey, []kcr.CapabilitiesRegistryNodeParams{singleNodeParams})
			if err != nil {
				err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
				if strings.Contains(err.Error(), "NodeAlreadyExists") {
					lggr.Warnw("node already exists, skipping", "p2pid", hex.EncodeToString(singleNodeParams.P2pId[:]))
					continue
				}
				return nil, fmt.Errorf("failed to call AddNode for node with p2pid %v: %w", singleNodeParams.P2pId, err)
			}
			// 1-by-1 tx is pending and we need to wait for it to be mined
			_, err = req.chain.Confirm(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to confirm AddNode of p2pid node %v transaction %s: %w", singleNodeParams.P2pId, tx.Hash().String(), err)
			}
			lggr.Debugw("registered node", "p2pid", singleNodeParams.P2pId)
		}
	} else {
		// the bulk add tx is pending and we need to wait for it to be mined
		_, err = req.chain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash().String(), err)
		}
	}
	return &registerNodesResponse{
		nodeIDToParams: nodeIDToParams,
	}, nil
}

type registerDonsRequest struct {
	registry *kcr.CapabilitiesRegistry
	chain    deployment.Chain

	nodeIDToParams    map[string]kcr.CapabilitiesRegistryNodeParams
	donToCapabilities map[string][]RegisteredCapability
	donToOcr2Nodes    map[string][]*ocr2Node
}

type registerDonsResponse struct {
	donInfos map[string]kcr.CapabilitiesRegistryDONInfo
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
	lggr.Infow("registering DONs...", "len", len(req.donToOcr2Nodes))
	// track hash of sorted p2pids to don name because the registry return value does not include the don name
	// and we need to map it back to the don name to access the other mapping data such as the don's capabilities & nodes
	p2pIdsToDon := make(map[string]string)
	var addedDons = 0

	donInfos, err := req.registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	existingDONs := make(map[string]struct{})
	for _, donInfo := range donInfos {
		existingDONs[sortedHash(donInfo.NodeP2PIds)] = struct{}{}
	}
	lggr.Infow("fetched existing DONs...", "len", len(donInfos), "lenByNodesHash", len(existingDONs))

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

		if _, ok := existingDONs[p2pSortedHash]; ok {
			lggr.Debugw("don already exists, ignoring", "don", don, "p2p sorted hash", p2pSortedHash)
			continue
		}

		caps, ok := req.donToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node operator %s", don)
		}
		wfSupported := false
		var cfgs []kcr.CapabilitiesRegistryCapabilityConfiguration
		for _, cap := range caps {
			if cap.CapabilityType == 2 { // OCR3 capability => WF supported
				wfSupported = true
			}
			// TODO: accept configuration from external source for each (don,capability)
			capCfg := DefaultCapConfig(cap.CapabilityType, len(p2pIds))
			cfgb, err := proto.Marshal(capCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal capability config for %v: %w", cap, err)
			}
			cfgs = append(cfgs, kcr.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: cap.ID,
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
		addedDons++
	}
	lggr.Debugf("Registered all DONs (new=%d), waiting for registry to update", addedDons)

	// occasionally the registry does not return the expected number of DONS immediately after the txns above
	// so we retry a few times. while crude, it is effective
	foundAll := false
	for i := 0; i < 10; i++ {
		lggr.Debugw("attempting to get DONs from registry", "attempt#", i)
		donInfos, err = req.registry.GetDONs(&bind.CallOpts{})
		if !containsAllDONs(donInfos, p2pIdsToDon) {
			lggr.Debugw("some expected dons not registered yet, re-checking after a delay ...")
			time.Sleep(2 * time.Second)
		} else {
			foundAll = true
			break
		}
	}
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	if !foundAll {
		return nil, fmt.Errorf("did not find all desired DONS")
	}

	resp := registerDonsResponse{
		donInfos: make(map[string]kcr.CapabilitiesRegistryDONInfo),
	}
	for i, donInfo := range donInfos {
		donName, ok := p2pIdsToDon[sortedHash(donInfo.NodeP2PIds)]
		if !ok {
			lggr.Debugw("irrelevant DON found in the registry, ignoring", "p2p sorted hash", sortedHash(donInfo.NodeP2PIds))
			continue
		}
		lggr.Debugw("adding don info to the reponse (keyed by DON name)", "don", donName)
		resp.donInfos[donName] = donInfos[i]
	}
	return &resp, nil
}

// are all DONs from p2pIdsToDon in donInfos
func containsAllDONs(donInfos []kcr.CapabilitiesRegistryDONInfo, p2pIdsToDon map[string]string) bool {
	found := make(map[string]struct{})
	for _, donInfo := range donInfos {
		hash := sortedHash(donInfo.NodeP2PIds)
		if _, ok := p2pIdsToDon[hash]; ok {
			found[hash] = struct{}{}
		}
	}
	return len(found) == len(p2pIdsToDon)
}

// configureForwarder sets the config for the forwarder contract on the chain for all Dons that accept workflows
// dons that don't accept workflows are not registered with the forwarder
func configureForwarder(lggr logger.Logger, chain deployment.Chain, fwdr *kf.KeystoneForwarder, dons []RegisteredDon) error {
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
		lggr.Debugw("configured forwarder", "forwarder", fwdr.Address().String(), "donId", dn.Info.Id, "version", ver, "f", dn.Info.F, "signers", dn.signers())
	}
	return nil
}

type configureOCR3Request struct {
	cfg      *OracleConfigWithSecrets
	chain    deployment.Chain
	contract *kocr3.OCR3Capability
	nodes    []*ocr2Node
}
type configureOCR3Response struct {
	ocrConfig Orc2drOracleConfig
}

func configureOCR3contract(req configureOCR3Request) (*configureOCR3Response, error) {
	if req.contract == nil {
		return nil, fmt.Errorf("OCR3 contract is nil")
	}
	nks := makeNodeKeysSlice(req.nodes)
	ocrConfig, err := GenerateOCR3Config(*req.cfg, nks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}
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
