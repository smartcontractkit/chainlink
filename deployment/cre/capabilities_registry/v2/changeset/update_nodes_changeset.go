package changeset

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

var _ cldf.ChangeSetV2[UpdateNodesInput] = UpdateNodes{}

type UpdateNodesInput struct {
	RegistryChainSel  uint64                           `json:"registryChainSel" yaml:"registryChainSel"`
	RegistryQualifier string                           `json:"registryQualifier" yaml:"registryQualifier"`
	MCMSConfig        *crecontracts.MCMSConfig         `json:"mcmsConfig,omitempty" yaml:"mcmsConfig,omitempty"`
	Nodes             []CapabilitiesRegistryNodeParams `json:"nodes" yaml:"nodes"`
}

type UpdateNodes struct{}

func (u UpdateNodes) VerifyPreconditions(e cldf.Environment, config UpdateNodesInput) error {
	if len(config.Nodes) == 0 {
		return errors.New("nodes list cannot be empty")
	}
	for i, node := range config.Nodes {
		if node.NOP == "" {
			return fmt.Errorf("node at index %d has empty NOP", i)
		}
		if node.CsaKey == "" {
			return fmt.Errorf("node at index %d has empty CsaKey", i)
		}
	}
	_, err := resolveNodesFromJD(e, config.RegistryChainSel, config.Nodes)
	return err
}

func (u UpdateNodes) Apply(e cldf.Environment, config UpdateNodesInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	chain, ok := e.BlockChains.EVMChains()[config.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", config.RegistryChainSel)
	}

	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(registryAddressRef.Address), chain.Client,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
	}

	nopNameToID, err := buildNOPNameToIDMap(capReg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build NOP name-to-ID map: %w", err)
	}

	resolvedNodes, err := resolveNodesFromJD(e, config.RegistryChainSel, config.Nodes)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to resolve node keys from JD: %w", err)
	}

	nodeUpdates := make(map[string]contracts.NodeConfig, len(resolvedNodes))
	for _, node := range resolvedNodes {
		wrapper, err := node.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert node %s: %w", node.P2pID, err)
		}

		id, exists := nopNameToID[node.NOP]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("node operator %q not found in contract", node.NOP)
		}

		nodeUpdates[node.P2pID] = contracts.NodeConfig{
			Signer:              wrapper.Signer,
			EncryptionPublicKey: node.EncryptionPublicKey,
			CSAKey:              node.CsaKey,
			NodeOperatorID:      id,
		}
	}

	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		capReg.Address(),
		contracts.UpdateNodesDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	updateReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateNodes,
		contracts.UpdateNodesDeps{
			Env:                  &e,
			Strategy:             strategy,
			CapabilitiesRegistry: capReg,
		},
		contracts.UpdateNodesInput{
			ChainSelector: config.RegistryChainSel,
			NodesUpdates:  nodeUpdates,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update nodes: %w", err)
	}

	var proposals []mcmslib.TimelockProposal
	if updateReport.Output.Operation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]mcmstypes.BatchOperation{*updateReport.Output.Operation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for UpdateNodes on chain %d: %w", config.RegistryChainSel, mcmsErr)
		}
		proposals = append(proposals, *proposal)
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{updateReport.ToGenericReport()},
		MCMSTimelockProposals: proposals,
	}, nil
}

// resolveNodesFromJD uses JD as the source of truth, keying lookup on CsaKey.
// It populates P2pID, Signer, and EncryptionPublicKey from JD for any node
// that leaves those fields empty, and validates that fields already set in the
// input match the JD values. CsaKey is required in the input and validated
// against JD. JD must be available and must contain every node in the list.
//
// Field mappings from JD:
//   - P2pID               → PeerID
//   - Signer              → OCR2 OnchainPublicKey (20-byte EVM address), right-padded with zeros to 32 bytes
//   - EncryptionPublicKey → WorkflowKey
//   - CsaKey              → CSAKey (PublicKey in JD)
func resolveNodesFromJD(e cldf.Environment, chainSel uint64, nodes []CapabilitiesRegistryNodeParams) ([]CapabilitiesRegistryNodeParams, error) {
	if e.Offchain == nil {
		return nil, errors.New("JD client is required but not configured")
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	// Use CSA keys as the JD lookup identifier. Strip any leading "0x" so the
	// key matches the plain-hex format JD stores in PublicKey.
	csaKeys := make([]string, 0, len(nodes))
	for _, node := range nodes {
		csaKeys = append(csaKeys, strings.ToLower(strings.TrimPrefix(node.CsaKey, "0x")))
	}

	jdNodes, err := deployment.NodeInfo(csaKeys, e.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node info from JD: %w", err)
	}

	// Index by normalized CSA key (strip "csa_" prefix JD may add).
	jdNodeByCSA := make(map[string]deployment.Node, len(jdNodes))
	for _, n := range jdNodes {
		jdNodeByCSA[strings.ToLower(strings.TrimPrefix(n.CSAKey, "csa_"))] = n
	}

	resolved := make([]CapabilitiesRegistryNodeParams, len(nodes))
	var errs []error
	for i, node := range nodes {
		resolved[i] = node
		inputCSA := strings.ToLower(strings.TrimPrefix(node.CsaKey, "0x"))

		jdNode, ok := jdNodeByCSA[inputCSA]
		if !ok {
			errs = append(errs, fmt.Errorf("node with CSA key %s not found in JD", inputCSA))
			continue
		}

		// P2pID
		jdP2pID := jdNode.PeerID.String()
		if node.P2pID == "" {
			resolved[i].P2pID = jdP2pID
		} else if node.P2pID != jdP2pID {
			errs = append(errs, fmt.Errorf("node %s: P2P ID mismatch: config=%s JD=%s",
				inputCSA, node.P2pID, jdP2pID))
		}

		// Signer
		evmCC, exists := jdNode.OCRConfigForChainSelector(chainSel)
		if !exists {
			errs = append(errs, fmt.Errorf("node %s: no OCR config in JD for chain %d", inputCSA, chainSel))
		} else {
			// OnchainPublicKey is a 20-byte EVM address; right-pad to 32 bytes to match the
			// bytes32 signer stored in the capabilities registry.
			var padded [32]byte
			copy(padded[:], evmCC.OnchainPublicKey)
			jdSigner := hex.EncodeToString(padded[:])
			if node.Signer == "" {
				resolved[i].Signer = jdSigner
			} else {
				configSigner := strings.ToLower(strings.TrimPrefix(node.Signer, "0x"))
				if configSigner != jdSigner {
					errs = append(errs, fmt.Errorf("node %s: signer mismatch: config=%s JD=%s",
						inputCSA, configSigner, jdSigner))
				}
			}
		}

		// EncryptionPublicKey
		if jdNode.WorkflowKey == "" {
			errs = append(errs, fmt.Errorf("node %s: no WorkflowKey in JD", inputCSA))
		} else {
			jdEncKey := strings.ToLower(jdNode.WorkflowKey)
			if node.EncryptionPublicKey == "" {
				resolved[i].EncryptionPublicKey = jdEncKey
			} else {
				configEncKey := strings.ToLower(strings.TrimPrefix(node.EncryptionPublicKey, "0x"))
				if configEncKey != jdEncKey {
					errs = append(errs, fmt.Errorf("node %s: encryption public key mismatch: config=%s JD=%s",
						inputCSA, configEncKey, jdEncKey))
				}
			}
		}

		// CsaKey: validate input against JD
		if jdNode.CSAKey == "" {
			errs = append(errs, fmt.Errorf("node %s: no CSA key in JD", inputCSA))
		} else {
			// JD may store the CSA key with or without a "csa_" prefix.
			jdCsaKey := strings.ToLower(strings.TrimPrefix(jdNode.CSAKey, "csa_"))
			if inputCSA != jdCsaKey {
				errs = append(errs, fmt.Errorf("node %s: CSA key mismatch: config=%s JD=%s",
					inputCSA, inputCSA, jdCsaKey))
			}
		}
	}

	return resolved, errors.Join(errs...)
}

func buildNOPNameToIDMap(capReg *capabilities_registry_v2.CapabilitiesRegistry) (map[string]uint32, error) {
	contractNOPs, err := pkg.GetNodeOperators(nil, capReg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node operators from contract: %w", err)
	}

	nopNameToID := make(map[string]uint32, len(contractNOPs))
	for i, nop := range contractNOPs {
		nopNameToID[nop.Name] = uint32(i) + 1 //nolint:gosec // i is bounded by the contract's NOP list length
	}
	return nopNameToID, nil
}
