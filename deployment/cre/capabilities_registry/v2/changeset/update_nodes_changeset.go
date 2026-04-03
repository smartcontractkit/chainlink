package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

var _ cldf.ChangeSetV2[UpdateNodesInput] = UpdateNodes{}

type UpdateNodesInput struct {
	RegistryChainSel  uint64                           `json:"registryChainSel" yaml:"registryChainSel"`
	RegistryQualifier string                           `json:"registryQualifier" yaml:"registryQualifier"`
	MCMSConfig        *crecontracts.MCMSConfig         `json:"mcmsConfig,omitempty" yaml:"mcmsConfig,omitempty"`
	Nodes             []CapabilitiesRegistryNodeParams `json:"nodes" yaml:"nodes"`
}

type UpdateNodes struct{}

func (u UpdateNodes) VerifyPreconditions(_ cldf.Environment, config UpdateNodesInput) error {
	if len(config.Nodes) == 0 {
		return errors.New("nodes list cannot be empty")
	}
	for i, node := range config.Nodes {
		if node.P2pID == "" {
			return fmt.Errorf("node at index %d has empty P2pID", i)
		}
		if node.NOP == "" {
			return fmt.Errorf("node at index %d (%s) has empty NOP", i, node.P2pID)
		}
		if node.Signer == "" {
			return fmt.Errorf("node at index %d (%s) has empty Signer", i, node.P2pID)
		}
		if node.EncryptionPublicKey == "" {
			return fmt.Errorf("node at index %d (%s) has empty EncryptionPublicKey", i, node.P2pID)
		}
		if node.CsaKey == "" {
			return fmt.Errorf("node at index %d (%s) has empty CsaKey", i, node.P2pID)
		}
	}
	return nil
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

	nodeUpdates := make(map[string]contracts.NodeConfig, len(config.Nodes))
	for _, node := range config.Nodes {
		wrapper, err := node.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert node %s: %w", node.P2pID, err)
		}

		var nopID uint32
		if node.NOP != "" {
			id, exists := nopNameToID[node.NOP]
			if !exists {
				return cldf.ChangesetOutput{}, fmt.Errorf("node operator %q not found in contract", node.NOP)
			}
			nopID = id
		}

		nodeUpdates[node.P2pID] = contracts.NodeConfig{
			Signer:              wrapper.Signer,
			EncryptionPublicKey: node.EncryptionPublicKey,
			CSAKey:              node.CsaKey,
			NodeOperatorID:      nopID,
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
