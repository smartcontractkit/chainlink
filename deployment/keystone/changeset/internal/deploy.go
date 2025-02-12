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
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	kf "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
)

type ConfigureContractsRequest struct {
	RegistryChainSel uint64
	Env              *deployment.Environment

	Dons       []DonCapabilities // externally sourced based on the environment
	OCR3Config *OracleConfig     // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config
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

	contractSetsResp, err := GetContractSets(lggr, &GetContractSetsRequest{
		Chains:      req.Env.Chains,
		AddressBook: req.Env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	registry := contractSetsResp.ContractSets[req.RegistryChainSel].CapabilitiesRegistry
	cfgRegistryResp, err := ConfigureRegistry(ctx, lggr, &ConfigureRegistryRequest{
		ConfigureContractsRequest: req,
		CapabilitiesRegistry:      registry,
	}, req.Env.ExistingAddresses)
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

// DonInfo is DonCapabilities, but expanded to contain node information
type DonInfo struct {
	Name         string
	F            uint8
	Nodes        []deployment.Node
	Capabilities []DONCapabilityWithConfig // every capability is hosted on each node
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

type ConfigureRegistryRequest struct {
	ConfigureContractsRequest
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
}

func (r *ConfigureRegistryRequest) Validate() error {
	if r.CapabilitiesRegistry == nil {
		return errors.New("capabilities registry is nil")
	}
	return r.ConfigureContractsRequest.Validate()
}

// ConfigureRegistry configures the registry contract with the given DONS and their capabilities
// the address book is required to contain the addresses of the deployed registry contract
func ConfigureRegistry(ctx context.Context, lggr logger.Logger, req *ConfigureRegistryRequest, addrBook deployment.AddressBook) (*ConfigureContractsResponse, error) {
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
	donToCapabilities, err := mapDonsToCaps(req.CapabilitiesRegistry, donInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to map dons to capabilities: %w", err)
	}
	nopsToNodeIDs, err := nopsToNodes(donInfos, req.Dons, req.RegistryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to map nops to nodes: %w", err)
	}
	var capabilities []capabilities_registry.CapabilitiesRegistryCapability
	for _, don := range donToCapabilities {
		for _, cap := range don {
			capabilities = append(capabilities, cap.CapabilitiesRegistryCapability)
		}
	}
	_, err = AddCapabilities(lggr, req.CapabilitiesRegistry, req.Env.Chains[req.RegistryChainSel], capabilities, false)

	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities to registry: %w", err)
	}

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
		DonToCapabilities:     donToCapabilities,
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
		DonToCapabilities:     donToCapabilities,
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
			cfg:        cfg,
			chain:      registryChain,
			contract:   contract,
			nodes:      don.Nodes,
			ocrSecrets: env.OCRSecrets,
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
		cfg:        cfg.OCR3Config,
		chain:      registryChain,
		contract:   contract,
		nodes:      nodes,
		dryRun:     cfg.DryRun,
		useMCMS:    cfg.UseMCMS,
		ocrSecrets: env.OCRSecrets,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureOCR3Resp{
		OCR2OracleConfig: r.ocrConfig,
		Ops:              r.ops,
	}, nil
}

type RegisteredCapability struct {
	capabilities_registry.CapabilitiesRegistryCapability
	ID     [32]byte
	Config *capabilitiespb.CapabilityConfig
}

func FromCapabilitiesRegistryCapability(capReg *capabilities_registry.CapabilitiesRegistryCapability, cfg *capabilitiespb.CapabilityConfig, e deployment.Environment, registryChainSelector uint64) (*RegisteredCapability, error) {
	registry, _, err := GetRegistryContract(&e, registryChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}
	id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, capReg.LabelledName, capReg.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetHashedCapabilityId for capability %v: %w", capReg, err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required for capability %v", capReg)
	}
	return &RegisteredCapability{
		CapabilitiesRegistryCapability: *capReg,
		ID:                             id,
		Config:                         cfg,
	}, nil
}

type RegisterNOPSRequest struct {
	Env                   *deployment.Environment
	RegistryChainSelector uint64
	Nops                  []capabilities_registry.CapabilitiesRegistryNodeOperator
	UseMCMS               bool
}

type RegisterNOPSResponse struct {
	Nops []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded // if UseMCMS is false, a list of added node operators is returned
	Ops  *timelock.BatchChainOperation                                  // if UseMCMS is true, a batch proposal is returned and no transaction is confirmed on chain.
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

	nodeIDToParams := make(map[string]capabilities_registry.CapabilitiesRegistryNodeParams)
	for don, nodes := range req.DonToNodes {
		caps, ok := req.DonToCapabilities[don]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for don %s", don)
		}
		var (
			hashedCapabilityIDs [][32]byte
			capIDs              []string
		)
		for _, cap := range caps {
			hashedCapabilityIDs = append(hashedCapabilityIDs, cap.ID)
			capIDs = append(capIDs, hex.EncodeToString(cap.ID[:]))
		}
		lggr.Debugw("hashed capability ids", "don", don, "ids", capIDs)

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
				signer, enc, err := extractSignerEncryptionKeys(n, registryChain.Selector)
				if err != nil {
					return nil, fmt.Errorf("failed to extract signer and encryption keys for node %s: %w", n.NodeID, err)
				}
				params = capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      nop.NodeOperatorId,
					Signer:              signer,
					P2pId:               n.PeerID,
					EncryptionPublicKey: enc,
					HashedCapabilityIds: hashedCapabilityIDs,
				}
			} else {
				// when we have a node operator, we need to dedup capabilities against the existing ones
				var newCapIDs [][32]byte
				for _, proposedCapID := range hashedCapabilityIDs {
					shouldAdd := true
					for _, existingCapID := range params.HashedCapabilityIds {
						if existingCapID == proposedCapID {
							shouldAdd = false
							break
						}
					}
					if shouldAdd {
						newCapIDs = append(newCapIDs, proposedCapID)
					}
				}
				params.HashedCapabilityIds = append(params.HashedCapabilityIds, newCapIDs...)
			}
			nodeIDToParams[n.NodeID] = params
		}
	}
	addResp, err := AddNodes(lggr, &AddNodesRequest{
		RegistryChain:        registryChain,
		NodeParams:           nodeIDToParams,
		UseMCMS:              req.UseMCMS,
		CapabilitiesRegistry: registry,
	})
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNodes: %w", err)
	}

	return &RegisterNodesResponse{
		nodeIDToParams: nodeIDToParams,
		Ops:            addResp.Ops,
	}, nil
}

// chainSel must be an evm chain; it ought to be the registry chain
// the signer is the onchain public key
// the enc is the encryption public key
func extractSignerEncryptionKeys(n deployment.Node, chainSel uint64) (signer [32]byte, enc [32]byte, err error) {
	chainID, err := chainsel.ChainIdFromSelector(chainSel)
	if err != nil {
		return signer, enc, fmt.Errorf("error getting chain id for selector %d: %w", chainSel, err)
	}
	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chainID, 10), chainsel.FamilyEVM)
	if err != nil {
		return signer, enc, fmt.Errorf("error getting chain details for selector %d, chain id %d: %w", chainSel, chainID, err)
	}

	evmCC, exists := n.SelToOCRConfig[chainDetails]
	if !exists {
		return signer, enc, fmt.Errorf("config for selector %v not found on node (id: %s, name: %s)", chainSel, n.NodeID, n.Name)
	}
	copy(signer[:], evmCC.OnchainPublicKey)
	copy(enc[:], evmCC.ConfigEncryptionPublicKey[:])
	return signer, enc, nil
}

type AddNodesResponse struct {
	// AddedNodes is only populated if UseMCMS is false
	// It may be empty if no nodes were added
	AddedNodes []capabilities_registry.CapabilitiesRegistryNodeParams
	// Ops is only populated if UseMCMS is true
	Ops *timelock.BatchChainOperation
}
type AddNodesRequest struct {
	RegistryChain        deployment.Chain
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	NodeParams           map[string]capabilities_registry.CapabilitiesRegistryNodeParams // the node id to the node params may be any node-unique identifier such as p2p, csa, etc

	// If true, a batch proposal is created and returned. No transaction is confirmed on chain.
	UseMCMS bool
}

// AddNodes adds the nodes to the registry
//
// It is idempotent.
//
// It is an error if there are duplicate nodes in the request.
// It is an error if the referenced capabilities or nops are not already registered. See AddCapabilities, AddNops.
func AddNodes(lggr logger.Logger, req *AddNodesRequest) (*AddNodesResponse, error) {
	registry := req.CapabilitiesRegistry
	nodeIDToParams := req.NodeParams
	nodes2Add, err := getNodesToRegister(registry, nodeIDToParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes to register: %w", err)
	}

	lggr.Debugf("found %d missing nodes", len(nodes2Add))

	if len(nodes2Add) == 0 {
		lggr.Debug("no new nodes to register")
		return &AddNodesResponse{}, nil
	}

	lggr.Debugw("unique node params to add after deduplication", "count", len(nodes2Add), "params", nodes2Add)

	if req.UseMCMS {
		ops, err := addNodesMCMSProposal(registry, nodes2Add, req.RegistryChain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to generate proposal to add nodes: %w", err)
		}

		return &AddNodesResponse{
			Ops: ops,
		}, nil
	}

	tx, err := registry.AddNodes(req.RegistryChain.DeployerKey, nodes2Add)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		// no typed errors in the abi, so we have to do string matching
		// try to add all nodes in one go, if that fails, fall back to 1-by-1
		return nil, fmt.Errorf("failed to call AddNodes for bulk add nodes: %w", err)
	}
	// the bulk add tx is pending and we need to wait for it to be mined
	_, err = req.RegistryChain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddNode confirm transaction %s: %w", tx.Hash().String(), err)
	}

	return &AddNodesResponse{
		AddedNodes: nodes2Add,
	}, nil
}

// getNodesToRegister returns the nodes that are not already registered in the registry
func getNodesToRegister(
	registry *capabilities_registry.CapabilitiesRegistry,
	nodeIDToParams map[string]capabilities_registry.CapabilitiesRegistryNodeParams,
) ([]capabilities_registry.CapabilitiesRegistryNodeParams, error) {
	nodes2Add := make([]capabilities_registry.CapabilitiesRegistryNodeParams, 0)
	for nodeID, nodeParams := range nodeIDToParams {
		var (
			ni  capabilities_registry.INodeInfoProviderNodeInfo
			err error
		)
		if ni, err = registry.GetNode(&bind.CallOpts{}, nodeParams.P2pId); err != nil {
			if err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err); strings.Contains(err.Error(), "NodeDoesNotExist") {
				nodes2Add = append(nodes2Add, nodeParams)
				continue
			}
			return nil, fmt.Errorf("failed to call GetNode for node %s: %w", nodeID, err)
		}

		// if no error, but node info is empty, then the node does not exist and should be added.
		if hex.EncodeToString(ni.P2pId[:]) != hex.EncodeToString(nodeParams.P2pId[:]) && hex.EncodeToString(ni.P2pId[:]) == "0000000000000000000000000000000000000000000000000000000000000000" {
			nodes2Add = append(nodes2Add, nodeParams)
			continue
		}
	}
	return nodes2Add, nil
}

// addNodesMCMSProposal generates a single call to AddNodes for all the node params at once.
func addNodesMCMSProposal(registry *capabilities_registry.CapabilitiesRegistry, params []capabilities_registry.CapabilitiesRegistryNodeParams, regChainSel uint64) (*timelock.BatchChainOperation, error) {
	tx, err := registry.AddNodes(deployment.SimTransactOpts(), params)
	if err != nil {
		err = deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to simulate call to AddNodes: %w", err)
	}

	return &timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(regChainSel),
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

	mcmsOps := make([]mcms.Operation, 0)
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
		regCaps, ok := req.DonToCapabilities[don.Name]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for DON %s", don.Name)
		}
		wfSupported := false
		var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
		for _, regCap := range regCaps {
			if regCap.CapabilityType == 2 { // OCR3 capability => WF supported
				wfSupported = true
			}
			if regCap.Config == nil {
				return nil, fmt.Errorf("config not found for capability %v", regCap)
			}
			cfgB, capErr := proto.Marshal(regCap.Config)
			if capErr != nil {
				return nil, fmt.Errorf("failed to marshal config for capability %v: %w", regCap, capErr)
			}
			cfgs = append(cfgs, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: regCap.ID,
				Config:       cfgB,
			})
		}

		txOpts := registryChain.DeployerKey
		if req.UseMCMS {
			txOpts = deployment.SimTransactOpts()
		}

		lggr.Debugw("calling add don", "don", don.Name, "p2p sorted hash", p2pSortedHash, "cgs", cfgs, "wfSupported", wfSupported, "f", don.F,
			"p2pids", p2pIds, "node count", len(p2pIds))
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
		if len(mcmsOps) > 0 {
			return &RegisterDonsResponse{
				Ops: &timelock.BatchChainOperation{
					ChainIdentifier: mcms.ChainIdentifier(registryChain.Selector),
					Batch:           mcmsOps,
				},
			}, nil
		}
		return &RegisterDonsResponse{}, nil
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
func configureForwarder(lggr logger.Logger, chain deployment.Chain, fwdr *kf.KeystoneForwarder, dons []RegisteredDon, useMCMS bool) (map[uint64]timelock.BatchChainOperation, error) {
	if fwdr == nil {
		return nil, errors.New("nil forwarder contract")
	}
	var (
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
