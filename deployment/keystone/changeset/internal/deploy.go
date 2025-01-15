package internal

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	chainsel "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	kf "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type ConfigureContractsRequest struct {
	RegistryChainSel uint64
	Env              *deployment.Environment

	Dons       []DonCapabilities // externally sourced based on the environment
	OCR3Config *OracleConfig     // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config

	// TODO rm this option; unused
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
	DonInfos  map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

// ConfigureContracts configures contracts them with the given DONS and their capabilities. It optionally deploys the contracts
// but best practice is to deploy them separately and pass the address book in the request
func ConfigureContracts(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest) (*ConfigureContractsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	cfgRegistryResp, err := ConfigureRegistry(ctx, lggr, req, req.Env.ExistingAddresses)
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
	// ignore response because we are not using mcms here and therefore no proposals are returned
	_, err = ConfigureForwardContracts(req.Env, ConfigureForwarderContractsRequest{
		Dons: dons,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure forwarder contracts: %w", err)
	}

	err = ConfigureOCR3Contract(req.Env, req.RegistryChainSel, dons, req.OCR3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OCR3 contract: %w", err)
	}

	return &ConfigureContractsResponse{
		Changeset: &deployment.ChangesetOutput{}, // no new addresses, proposals etc
		DonInfos:  cfgRegistryResp.DonInfos,
	}, nil
}

// DeployContracts deploys the all the keystone contracts on all chains and returns the address book in the changeset
func DeployContracts(e *deployment.Environment, chainSel uint64) (*deployment.ChangesetOutput, error) {
	lggr := e.Logger
	adbook := deployment.NewMemoryAddressBook()
	// deploy contracts on all chains and track the registry and ocr3 contracts
	for _, chain := range e.Chains {
		lggr.Infow("deploying contracts", "chain", chain.Selector)
		deployResp, err := deployContractsToChain(deployContractsRequest{
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
	F            uint8
	Nodes        []deployment.Node
	Capabilities []capabilities_registry.CapabilitiesRegistryCapability // every capability is hosted on each node
}

func DonInfos(dons []DonCapabilities, jd deployment.OffchainClient) ([]DonInfo, error) {
	var donInfos []DonInfo
	for _, don := range dons {
		var nodeIDs []string
		for _, nop := range don.Nops {
			nodeIDs = append(nodeIDs, nop.Nodes...)
		}
		nodes, err := deployment.NodeInfo(nodeIDs, jd)
		if err != nil {
			return nil, err
		}
		donInfos = append(donInfos, DonInfo{
			Name:         don.Name,
			F:            don.F,
			Nodes:        nodes,
			Capabilities: don.Capabilities,
		})
	}
	return donInfos, nil
}

func GetRegistryContract(e *deployment.Environment, registryChainSel uint64) (*capabilities_registry.CapabilitiesRegistry, deployment.Chain, error) {
	registryChain, ok := e.Chains[registryChainSel]
	if !ok {
		return nil, deployment.Chain{}, fmt.Errorf("chain %d not found in environment", registryChainSel)
	}

	contractSetsResp, err := GetContractSets(e.Logger, &GetContractSetsRequest{
		Chains:      e.Chains,
		AddressBook: e.ExistingAddresses,
	})
	if err != nil {
		return nil, deployment.Chain{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	// ensure registry is deployed and get the registry contract and chain
	var registry *capabilities_registry.CapabilitiesRegistry
	registryChainContracts, ok := contractSetsResp.ContractSets[registryChainSel]
	if !ok {
		return nil, deployment.Chain{}, fmt.Errorf("failed to deploy registry chain contracts. expected chain %d", registryChainSel)
	}
	registry = registryChainContracts.CapabilitiesRegistry
	if registry == nil {
		return nil, deployment.Chain{}, errors.New("no registry contract found")
	}
	e.Logger.Debugf("registry contract address: %s, chain %d", registry.Address().String(), registryChainSel)
	return registry, registryChain, nil
}

// ConfigureRegistry configures the registry contract with the given DONS and their capabilities
// the address book is required to contain the addresses of the deployed registry contract
func ConfigureRegistry(ctx context.Context, lggr logger.Logger, req ConfigureContractsRequest, addrBook deployment.AddressBook) (*ConfigureContractsResponse, error) {
	donInfos, err := DonInfos(req.Dons, req.Env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get don infos: %w", err)
	}

	// all the subsequent calls to the registry are in terms of nodes
	// compute the mapping of dons to their nodes for reuse in various registry calls
	donToNodes, err := mapDonsToNodes(donInfos, true, req.RegistryChainSel)
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
	capabilitiesResp, err := RegisterCapabilities(lggr, RegisterCapabilitiesRequest{
		Env:                   req.Env,
		RegistryChainSelector: req.RegistryChainSel,
		DonToCapabilities:     donToCapabilities,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register capabilities: %w", err)
	}
	lggr.Infow("registered capabilities", "capabilities", capabilitiesResp.DonToCapabilities)

	// register node operators
	nopsList := maps.Keys(nopsToNodeIDs)
	nopsResp, err := RegisterNOPS(ctx, lggr, RegisterNOPSRequest{
		Env:                   req.Env,
		RegistryChainSelector: req.RegistryChainSel,
		Nops:                  nopsList,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register node operators: %w", err)
	}
	lggr.Infow("registered node operators", "nops", nopsResp.Nops)

	// register nodes
	nodesResp, err := RegisterNodes(lggr, &RegisterNodesRequest{
		Env:                   req.Env,
		RegistryChainSelector: req.RegistryChainSel,
		NopToNodeIDs:          nopsToNodeIDs,
		DonToNodes:            donToNodes,
		DonToCapabilities:     capabilitiesResp.DonToCapabilities,
		Nops:                  nopsResp.Nops,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register nodes: %w", err)
	}
	lggr.Infow("registered nodes", "nodes", nodesResp.nodeIDToParams)

	// TODO: annotate nodes with node_operator_id in JD?

	donsToRegister := []DONToRegister{}
	for _, don := range req.Dons {
		nodes, ok := donToNodes[don.Name]
		if !ok {
			return nil, fmt.Errorf("nodes not found for don %s", don.Name)
		}
		f := don.F
		if f == 0 {
			// TODO: fallback to a default value for compatibility - change to error
			f = uint8(len(nodes) / 3)
			lggr.Warnw("F not set for don - falling back to default", "don", don.Name, "f", f)
		}
		donsToRegister = append(donsToRegister, DONToRegister{
			Name:  don.Name,
			F:     f,
			Nodes: nodes,
		})
	}

	nodeIdToP2PID := map[string][32]byte{}
	for nodeID, params := range nodesResp.nodeIDToParams {
		nodeIdToP2PID[nodeID] = params.P2pId
	}
	// register DONS
	donsResp, err := RegisterDons(lggr, RegisterDonsRequest{
		Env:                   req.Env,
		RegistryChainSelector: req.RegistryChainSel,
		NodeIDToP2PID:         nodeIdToP2PID,
		DonToCapabilities:     capabilitiesResp.DonToCapabilities,
		DonsToRegister:        donsToRegister,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register DONS: %w", err)
	}
	lggr.Infow("registered DONs", "dons", len(donsResp.DonInfos))

	return &ConfigureContractsResponse{
		Changeset: &deployment.ChangesetOutput{}, // no new addresses, proposals etc
		DonInfos:  donsResp.DonInfos,
	}, nil
}

// Depreciated: use changeset.ConfigureOCR3Contract instead
// ocr3 contract on the registry chain for the wf dons
func ConfigureOCR3Contract(env *deployment.Environment, chainSel uint64, dons []RegisteredDon, cfg *OracleConfig) error {
	registryChain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", chainSel)
	}

	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
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

		contract, err := contracts.GetOCR3Contract(nil)
		if err != nil {
			env.Logger.Errorf("failed to get OCR3 contract: %s", err)
			return fmt.Errorf("failed to get OCR3 contract: %w", err)
		}

		_, err = configureOCR3contract(configureOCR3Request{
			cfg:         cfg,
			chain:       registryChain,
			contract:    contract,
			nodes:       don.Nodes,
			contractSet: &contracts,
			ocrSecrets:  env.OCRSecrets,
		})
		if err != nil {
			return fmt.Errorf("failed to configure OCR3 contract for don %s: %w", don.Name, err)
		}
	}
	return nil
}

type ConfigureOCR3Resp struct {
	OCR2OracleConfig
	Ops *timelock.BatchChainOperation
}

type ConfigureOCR3Config struct {
	ChainSel   uint64
	NodeIDs    []string
	Address    *common.Address // address of the OCR3 contract to configure
	OCR3Config *OracleConfig
	DryRun     bool

	UseMCMS bool
}

// Depreciated: use changeset.ConfigureOCR3Contract instead
func ConfigureOCR3ContractFromJD(env *deployment.Environment, cfg ConfigureOCR3Config) (*ConfigureOCR3Resp, error) {
	prefix := ""
	if cfg.DryRun {
		prefix = "DRY RUN: "
	}
	env.Logger.Infof("%sconfiguring OCR3 contract for chain %d", prefix, cfg.ChainSel)
	registryChain, ok := env.Chains[cfg.ChainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found in environment", cfg.ChainSel)
	}
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contracts, ok := contractSetsResp.ContractSets[cfg.ChainSel]
	if !ok {
		return nil, fmt.Errorf("failed to get contract set for chain %d", cfg.ChainSel)
	}

	contract, err := contracts.GetOCR3Contract(cfg.Address)
	if err != nil {
		env.Logger.Errorf("%sfailed to get OCR3 contract at %s : %s", prefix, cfg.Address, err)
		return nil, fmt.Errorf("failed to get OCR3 contract: %w", err)
	}

	nodes, err := deployment.NodeInfo(cfg.NodeIDs, env.Offchain)
	if err != nil {
		return nil, err
	}
	r, err := configureOCR3contract(configureOCR3Request{
		cfg:         cfg.OCR3Config,
		chain:       registryChain,
		contract:    contract,
		nodes:       nodes,
		dryRun:      cfg.DryRun,
		contractSet: &contracts,
		useMCMS:     cfg.UseMCMS,
		ocrSecrets:  env.OCRSecrets,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureOCR3Resp{
		OCR2OracleConfig: r.ocrConfig,
		Ops:              r.ops,
	}, nil
}

type RegisterCapabilitiesRequest struct {
	Env                   *deployment.Environment
	RegistryChainSelector uint64
	DonToCapabilities     map[string][]capabilities_registry.CapabilitiesRegistryCapability

	// if UseMCMS is true, a batch proposal is returned and no transaction is confirmed on chain.
	UseMCMS bool
}

type RegisterCapabilitiesResponse struct {
	DonToCapabilities map[string][]RegisteredCapability
	Ops               *timelock.BatchChainOperation
}

type RegisteredCapability struct {
	capabilities_registry.CapabilitiesRegistryCapability
	ID [32]byte
}

func FromCapabilitiesRegistryCapability(cap *capabilities_registry.CapabilitiesRegistryCapability, e deployment.Environment, registryChainSelector uint64) (*RegisteredCapability, error) {
	registry, _, err := GetRegistryContract(&e, registryChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}
	id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetHashedCapabilityId for capability %v: %w", cap, err)
	}
	return &RegisteredCapability{
		CapabilitiesRegistryCapability: *cap,
		ID:                             id,
	}, nil
}

// RegisterCapabilities add computes the capability id, adds it to the registry and associates the registered capabilities with appropriate don(s)
func RegisterCapabilities(lggr logger.Logger, req RegisterCapabilitiesRequest) (*RegisterCapabilitiesResponse, error) {
	if len(req.DonToCapabilities) == 0 {
		return nil, errors.New("no capabilities to register")
	}
	cresp, err := GetContractSets(req.Env.Logger, &GetContractSetsRequest{
		Chains:      req.Env.Chains,
		AddressBook: req.Env.ExistingAddresses,
	})
	contracts := cresp.ContractSets[req.RegistryChainSelector]
	registry := contracts.CapabilitiesRegistry
	registryChain := req.Env.Chains[req.RegistryChainSelector]

	lggr.Infow("registering capabilities...", "len", len(req.DonToCapabilities))
	resp := &RegisterCapabilitiesResponse{
		DonToCapabilities: make(map[string][]RegisteredCapability),
	}

	// capability could be hosted on multiple dons. need to deduplicate
	uniqueCaps := make(map[capabilities_registry.CapabilitiesRegistryCapability][32]byte)
	for don, caps := range req.DonToCapabilities {
		var registerCaps []RegisteredCapability
		for _, cap := range caps {
			id, ok := uniqueCaps[cap]
			if !ok {
				var err error
				id, err = registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
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
		resp.DonToCapabilities[don] = registerCaps
	}

	var capabilities []capabilities_registry.CapabilitiesRegistryCapability
	for cap := range uniqueCaps {
		capabilities = append(capabilities, cap)
	}
	if len(capabilities) == 0 {
		lggr.Warn("no new capabilities to register")
		return &RegisterCapabilitiesResponse{}, nil
	}

	ops, err := AddCapabilities(lggr, contracts.CapabilitiesRegistry, registryChain, capabilities, req.UseMCMS)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	resp.Ops = ops

	return resp, nil
}

type RegisterNOPSRequest struct {
	Env                   *deployment.Environment
	RegistryChainSelector uint64
	Nops                  []capabilities_registry.CapabilitiesRegistryNodeOperator
	UseMCMS               bool
}

type RegisterNOPSResponse struct {
	Nops []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
	Ops  *timelock.BatchChainOperation
}

func RegisterNOPS(ctx context.Context, lggr logger.Logger, req RegisterNOPSRequest) (*RegisterNOPSResponse, error) {
	registry, registryChain, err := GetRegistryContract(req.Env, req.RegistryChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}
	lggr.Infow("registering node operators...", "len", len(req.Nops))
	existingNops, err := registry.GetNodeOperators(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	existingNopsAddrToID := make(map[capabilities_registry.CapabilitiesRegistryNodeOperator]uint32)
	for id, nop := range existingNops {
		existingNopsAddrToID[nop] = uint32(id)
	}
	lggr.Infow("fetched existing node operators", "len", len(existingNopsAddrToID))
	resp := &RegisterNOPSResponse{
		Nops: []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded{},
	}
	nops := []capabilities_registry.CapabilitiesRegistryNodeOperator{}
	for _, nop := range req.Nops {
		if id, ok := existingNopsAddrToID[nop]; !ok {
			nops = append(nops, nop)
		} else {
			lggr.Debugw("node operator already exists", "name", nop.Name, "admin", nop.Admin.String(), "id", id)
			resp.Nops = append(resp.Nops, &capabilities_registry.CapabilitiesRegistryNodeOperatorAdded{
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

	if req.UseMCMS {
		ops, err := addNOPsMCMSProposal(registry, nops, registryChain)
		if err != nil {
			return nil, fmt.Errorf("failed to generate proposal to add node operators: %w", err)
		}

		resp.Ops = ops
		return resp, nil
	}

	tx, err := registry.AddNodeOperators(registryChain.DeployerKey, nops)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
	}

	// for some reason that i don't understand, the confirm must be called before the WaitMined or the latter will hang
	// (at least for a simulated backend chain)
	_, err = registryChain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
	}

	receipt, err := bind.WaitMined(ctx, registryChain.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to mine AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
	}
	if len(receipt.Logs) != len(nops) {
		return nil, fmt.Errorf("expected %d log entries for AddNodeOperators, got %d", len(nops), len(receipt.Logs))
	}
	for i, log := range receipt.Logs {
		o, err := registry.ParseNodeOperatorAdded(*log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log %d for operator added: %w", i, err)
		}
		resp.Nops = append(resp.Nops, o)
	}

	return resp, nil
}

func addNOPsMCMSProposal(registry *capabilities_registry.CapabilitiesRegistry, nops []capabilities_registry.CapabilitiesRegistryNodeOperator, regChain deployment.Chain) (*timelock.BatchChainOperation, error) {
	tx, err := registry.AddNodeOperators(deployment.SimTransactOpts(), nops)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
	}

	return &timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(regChain.Selector),
		Batch: []mcms.Operation{
			{
				To:    registry.Address(),
				Data:  tx.Data(),
				Value: big.NewInt(0),
			},
		},
	}, nil
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

// register nodes
type RegisterNodesRequest struct {
	Env                   *deployment.Environment
	RegistryChainSelector uint64
	NopToNodeIDs          map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string
	DonToNodes            map[string][]deployment.Node
	DonToCapabilities     map[string][]RegisteredCapability
	Nops                  []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
	UseMCMS               bool
}
type RegisterNodesResponse struct {
	nodeIDToParams map[string]capabilities_registry.CapabilitiesRegistryNodeParams
	Ops            *timelock.BatchChainOperation
}

// registerNodes registers the nodes with the registry. it assumes that the deployer key in the Chain
// can sign the transactions update the contract state
func RegisterNodes(lggr logger.Logger, req *RegisterNodesRequest) (*RegisterNodesResponse, error) {
	registry, registryChain, err := GetRegistryContract(req.Env, req.RegistryChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}

	var count int
	for _, nodes := range req.NopToNodeIDs {
		count += len(nodes)
	}
	lggr.Infow("registering nodes...", "len", count)
	nodeToRegisterNop := make(map[string]*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded)
	for _, nop := range req.Nops {
		n := capabilities_registry.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: nop.Admin,
		}
		nodeIDs := req.NopToNodeIDs[n]
		for _, nodeID := range nodeIDs {
			_, exists := nodeToRegisterNop[nodeID]
			if !exists {
				nodeToRegisterNop[nodeID] = nop
			}
		}
	}

	// TODO: deduplicate everywhere
	registryChainID, err := chainsel.ChainIdFromSelector(registryChain.Selector)
	if err != nil {
		return nil, err
	}
	registryChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(registryChainID)), chainsel.FamilyEVM)
	if err != nil {
		return nil, err
	}

	nodeIDToParams := make(map[string]capabilities_registry.CapabilitiesRegistryNodeParams)
	for don, nodes := range req.DonToNodes {
		caps, ok := req.DonToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for don %s", don)
		}
		var hashedCapabilityIds [][32]byte
		for _, cap := range caps {
			hashedCapabilityIds = append(hashedCapabilityIds, cap.ID)
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", hashedCapabilityIds)

		for _, n := range nodes {
			if n.IsBootstrap { // bootstraps are part of the DON but don't host capabilities
				continue
			}
			nop, ok := nodeToRegisterNop[n.NodeID]
			if !ok {
				return nil, fmt.Errorf("node operator not found for node %s", n.NodeID)
			}
			params, ok := nodeIDToParams[n.NodeID]

			if !ok {
				evmCC, exists := n.SelToOCRConfig[registryChainDetails]
				if !exists {
					return nil, fmt.Errorf("config for selector %v not found on node (id: %s, name: %s)", registryChain.Selector, n.NodeID, n.Name)
				}
				var signer [32]byte
				copy(signer[:], evmCC.OnchainPublicKey)
				var csakey [32]byte
				copy(csakey[:], evmCC.ConfigEncryptionPublicKey[:])
				params = capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      nop.NodeOperatorId,
					Signer:              signer,
					P2pId:               n.PeerID,
					EncryptionPublicKey: csakey,
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
			nodeIDToParams[n.NodeID] = params
		}
	}

	var uniqueNodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, v := range nodeIDToParams {
		uniqueNodeParams = append(uniqueNodeParams, v)
	}
	lggr.Debugw("unique node params to add", "count", len(uniqueNodeParams), "params", uniqueNodeParams)

	if req.UseMCMS {
		ops, err := addNodesMCMSProposal(registry, uniqueNodeParams, registryChain)
		if err != nil {
			return nil, fmt.Errorf("failed to generate proposal to add nodes: %w", err)
		}

		return &RegisterNodesResponse{
			nodeIDToParams: nodeIDToParams,
			Ops:            ops,
		}, nil
	}

	tx, err := registry.AddNodes(registryChain.DeployerKey, uniqueNodeParams)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		// no typed errors in the abi, so we have to do string matching
		// try to add all nodes in one go, if that fails, fall back to 1-by-1
		if !strings.Contains(err.Error(), "NodeAlreadyExists") {
			return nil, fmt.Errorf("failed to call AddNodes for bulk add nodes: %w", err)
		}
		lggr.Warn("nodes already exist, falling back to 1-by-1")
		for _, singleNodeParams := range uniqueNodeParams {
			tx, err = registry.AddNodes(registryChain.DeployerKey, []capabilities_registry.CapabilitiesRegistryNodeParams{singleNodeParams})
			if err != nil {
				err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
				if strings.Contains(err.Error(), "NodeAlreadyExists") {
					lggr.Warnw("node already exists, skipping", "p2pid", hex.EncodeToString(singleNodeParams.P2pId[:]))
					continue
				}
				return nil, fmt.Errorf("failed to call AddNode for node with p2pid %v: %w", singleNodeParams.P2pId, err)
			}
			// 1-by-1 tx is pending and we need to wait for it to be mined
			_, err = registryChain.Confirm(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to confirm AddNode of p2pid node %v transaction %s: %w", singleNodeParams.P2pId, tx.Hash().String(), err)
			}
			lggr.Debugw("registered node", "p2pid", singleNodeParams.P2pId)
		}
	} else {
		// the bulk add tx is pending and we need to wait for it to be mined
		_, err = registryChain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash().String(), err)
		}
	}

	return &RegisterNodesResponse{
		nodeIDToParams: nodeIDToParams,
	}, nil
}

// addNodesMCMSProposal generates a single call to AddNodes for all the node params at once.
func addNodesMCMSProposal(registry *capabilities_registry.CapabilitiesRegistry, params []capabilities_registry.CapabilitiesRegistryNodeParams, regChain deployment.Chain) (*timelock.BatchChainOperation, error) {
	tx, err := registry.AddNodes(deployment.SimTransactOpts(), params)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to simulate call to AddNodes: %w", err)
	}

	return &timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(regChain.Selector),
		Batch: []mcms.Operation{
			{
				To:    registry.Address(),
				Data:  tx.Data(),
				Value: big.NewInt(0),
			},
		},
	}, nil
}

type DONToRegister struct {
	Name  string
	F     uint8
	Nodes []deployment.Node
}

type RegisterDonsRequest struct {
	Env                   *deployment.Environment
	RegistryChainSelector uint64

	NodeIDToP2PID     map[string][32]byte
	DonToCapabilities map[string][]RegisteredCapability
	DonsToRegister    []DONToRegister
	UseMCMS           bool
}

type RegisterDonsResponse struct {
	DonInfos map[string]capabilities_registry.CapabilitiesRegistryDONInfo
	Ops      *timelock.BatchChainOperation
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

func RegisterDons(lggr logger.Logger, req RegisterDonsRequest) (*RegisterDonsResponse, error) {
	registry, registryChain, err := GetRegistryContract(req.Env, req.RegistryChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}
	lggr.Infow("registering DONs...", "len", len(req.DonsToRegister))
	// track hash of sorted p2pids to don name because the registry return value does not include the don name
	// and we need to map it back to the don name to access the other mapping data such as the don's capabilities & nodes
	p2pIdsToDon := make(map[string]string)
	var addedDons = 0

	donInfos, err := registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	existingDONs := make(map[string]struct{})
	for _, donInfo := range donInfos {
		existingDONs[sortedHash(donInfo.NodeP2PIds)] = struct{}{}
	}
	lggr.Infow("fetched existing DONs...", "len", len(donInfos), "lenByNodesHash", len(existingDONs))

	mcmsOps := make([]mcms.Operation, 0, len(req.DonsToRegister))
	for _, don := range req.DonsToRegister {
		var p2pIds [][32]byte
		for _, n := range don.Nodes {
			if n.IsBootstrap {
				continue
			}
			p2pID, ok := req.NodeIDToP2PID[n.NodeID]
			if !ok {
				return nil, fmt.Errorf("node params not found for non-bootstrap node %s", n.NodeID)
			}
			p2pIds = append(p2pIds, p2pID)
		}

		p2pSortedHash := sortedHash(p2pIds)
		p2pIdsToDon[p2pSortedHash] = don.Name

		if _, ok := existingDONs[p2pSortedHash]; ok {
			lggr.Debugw("don already exists, ignoring", "don", don.Name, "p2p sorted hash", p2pSortedHash)
			continue
		}

		lggr.Debugw("registering DON", "don", don.Name, "p2p sorted hash", p2pSortedHash)

		caps, ok := req.DonToCapabilities[don.Name]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for DON %s", don.Name)
		}
		wfSupported := false
		var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
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
			cfgs = append(cfgs, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: cap.ID,
				Config:       cfgb,
			})
		}

		txOpts := registryChain.DeployerKey
		if req.UseMCMS {
			txOpts = deployment.SimTransactOpts()
		}

		tx, err := registry.AddDON(txOpts, p2pIds, cfgs, true, wfSupported, don.F)
		if err != nil {
			err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
			return nil, fmt.Errorf("failed to call AddDON for don '%s' p2p2Id hash %s capability %v: %w", don.Name, p2pSortedHash, cfgs, err)
		}

		if req.UseMCMS {
			lggr.Debugw("adding mcms op for DON", "don", don.Name)
			op := mcms.Operation{
				To:    registry.Address(),
				Data:  tx.Data(),
				Value: big.NewInt(0),
			}
			mcmsOps = append(mcmsOps, op)
			continue
		}

		_, err = registryChain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm AddDON transaction %s for don %s: %w", tx.Hash().String(), don.Name, err)
		}
		lggr.Debugw("registered DON", "don", don.Name, "p2p sorted hash", p2pSortedHash, "cgs", cfgs, "wfSupported", wfSupported, "f", don.F)
		addedDons++
	}

	if req.UseMCMS {
		return &RegisterDonsResponse{
			Ops: &timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(registryChain.Selector),
				Batch:           mcmsOps,
			},
		}, nil
	}

	lggr.Debugf("Registered all DONs (new=%d), waiting for registry to update", addedDons)

	// occasionally the registry does not return the expected number of DONS immediately after the txns above
	// so we retry a few times. while crude, it is effective
	foundAll := false
	for i := 0; i < 10; i++ {
		lggr.Debugw("attempting to get DONs from registry", "attempt#", i)
		donInfos, err = registry.GetDONs(&bind.CallOpts{})
		if !containsAllDONs(donInfos, p2pIdsToDon) {
			lggr.Debugw("some expected dons not registered yet, re-checking after a delay ...")
			time.Sleep(2 * time.Second)
		} else {
			foundAll = true
			break
		}
	}
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call GetDONs: %w", err)
	}
	if !foundAll {
		return nil, errors.New("did not find all desired DONS")
	}

	resp := RegisterDonsResponse{
		DonInfos: make(map[string]capabilities_registry.CapabilitiesRegistryDONInfo),
	}
	for i, donInfo := range donInfos {
		donName, ok := p2pIdsToDon[sortedHash(donInfo.NodeP2PIds)]
		if !ok {
			lggr.Debugw("irrelevant DON found in the registry, ignoring", "p2p sorted hash", sortedHash(donInfo.NodeP2PIds))
			continue
		}
		lggr.Debugw("adding don info to the response (keyed by DON name)", "don", donName)
		resp.DonInfos[donName] = donInfos[i]
	}

	return &resp, nil
}

// are all DONs from p2pIdsToDon in donInfos
func containsAllDONs(donInfos []capabilities_registry.CapabilitiesRegistryDONInfo, p2pIdsToDon map[string]string) bool {
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
func configureForwarder(lggr logger.Logger, chain deployment.Chain, contractSet ContractSet, dons []RegisteredDon, useMCMS bool) (map[uint64]timelock.BatchChainOperation, error) {
	if contractSet.Forwarder == nil {
		return nil, errors.New("nil forwarder contract")
	}
	var (
		fwdr  = contractSet.Forwarder
		opMap = make(map[uint64]timelock.BatchChainOperation)
	)
	for _, dn := range dons {
		if !dn.Info.AcceptsWorkflows {
			continue
		}
		ver := dn.Info.ConfigCount // note config count on the don info is the version on the forwarder
		signers := dn.Signers(chainsel.FamilyEVM)
		txOpts := chain.DeployerKey
		if useMCMS {
			txOpts = deployment.SimTransactOpts()
		}
		tx, err := fwdr.SetConfig(txOpts, dn.Info.Id, ver, dn.Info.F, signers)
		if err != nil {
			err = deployment.DecodeErr(kf.KeystoneForwarderABI, err)
			return nil, fmt.Errorf("failed to call SetConfig for forwarder %s on chain %d: %w", fwdr.Address().String(), chain.Selector, err)
		}
		if !useMCMS {
			_, err = chain.Confirm(tx)
			if err != nil {
				err = deployment.DecodeErr(kf.KeystoneForwarderABI, err)
				return nil, fmt.Errorf("failed to confirm SetConfig for forwarder %s: %w", fwdr.Address().String(), err)
			}
		} else {
			// create the mcms proposals
			ops := timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chain.Selector),
				Batch: []mcms.Operation{
					{
						To:    fwdr.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			}
			opMap[chain.Selector] = ops
		}
		lggr.Debugw("configured forwarder", "forwarder", fwdr.Address().String(), "donId", dn.Info.Id, "version", ver, "f", dn.Info.F, "signers", signers)
	}
	return opMap, nil
}
