package contracts

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type RegisterNodesDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
}

type RegisterNodesInput struct {
	Address       string
	ChainSelector uint64
	Nodes         []NodesInput

	// AllNOPsInContract Optional: map of all NOP names to their IDs (or expected IDs) in the contract.
	// Useful when using MCMS, since the NOPs won't be present in the contract yet, but we can assume their IDs based on the existing NOPs.
	// If not provided, the operation will fail if it encounters a NOP that is not in the contract.
	AllNOPsInContract map[string]int
	MCMSConfig        *contracts.MCMSConfig
}

type NodesInput struct {
	NOP                 string
	Signer              [32]byte
	P2pID               [32]byte
	EncryptionPublicKey [32]byte
	CsaKey              [32]byte
	CapabilityIDs       []string
}

type RegisterNodesOutput struct {
	Nodes     []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	Operation *mcmstypes.BatchOperation
}

// RegisterNodes is an operation that registers nodes in the V2 Capabilities Registry contract.
var RegisterNodes = operations.NewOperation[RegisterNodesInput, RegisterNodesOutput, RegisterNodesDeps](
	"register-nodes-op",
	semver.MustParse("1.0.0"),
	"Register Nodes in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNodesDeps, input RegisterNodesInput) (RegisterNodesOutput, error) {
		// Validate input
		if input.Address == "" {
			return RegisterNodesOutput{}, errors.New("address is not set")
		}
		if len(input.Nodes) == 0 {
			// The contract allows to pass an empty array of nodes.
			return RegisterNodesOutput{
				Nodes: []capabilities_registry_v2.CapabilitiesRegistryNodeParams{},
			}, nil
		}
		if input.ChainSelector == 0 {
			return RegisterNodesOutput{}, errors.New("chainSelector is not set")
		}

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterNodesOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistry contract
		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}

		dedupedNodes, err := dedupNodes(deps.Env.Logger, input.Nodes, capReg)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to deduplicate nodes: %w", err)
		}
		if len(dedupedNodes) == 0 {
			deps.Env.Logger.Info("All nodes are already registered in the contract, nothing to do")
			return RegisterNodesOutput{
				Nodes: []capabilities_registry_v2.CapabilitiesRegistryNodeParams{},
			}, nil
		}

		allNOPsNamesInContract := input.AllNOPsInContract
		if len(allNOPsNamesInContract) == 0 {
			contractNOPs, err := pkg.GetNodeOperators(nil, capReg)
			if err != nil {
				return RegisterNodesOutput{}, fmt.Errorf("failed to fetch node operators from contract: %w", err)
			}

			allNOPsNamesInContract = make(map[string]int)
			for i, nop := range contractNOPs {
				// NodeOperatorId is 1-based and returned in order from the contract.
				// So the ID is the index + 1
				// See the implementation of `AddNodeOperators` in the contract for reference:
				// https://github.com/smartcontractkit/chainlink-evm/blob/develop/contracts/src/v0.8/workflow/v2/CapabilitiesRegistry.sol#L568
				allNOPsNamesInContract[nop.Name] = i + 1
			}
		}

		var nodes []capabilities_registry_v2.CapabilitiesRegistryNodeParams
		for _, node := range dedupedNodes {
			index, exists := allNOPsNamesInContract[node.NOP]
			if !exists {
				return RegisterNodesOutput{}, fmt.Errorf("node operator `%s` not found in contract", node.NOP)
			}

			nodes = append(nodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      uint32(index),
				Signer:              node.Signer,
				EncryptionPublicKey: node.EncryptionPublicKey,
				P2pId:               node.P2pID,
				CsaKey:              node.CsaKey,
				CapabilityIds:       node.CapabilityIDs,
			})
		}

		if err := validateNodes(nodes); err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("node validation failed: %w", err)
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return capReg.AddNodes(opts, nodes)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterNodesOutput{}, fmt.Errorf("failed to execute AddNodes: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterNodes on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d nodes on chain %d", len(nodes), input.ChainSelector)
		}

		return RegisterNodesOutput{
			Nodes:     nodes,
			Operation: operation,
		}, nil
	},
)

func validateNodes(nodes []capabilities_registry_v2.CapabilitiesRegistryNodeParams) error {
	for _, node := range nodes {
		if node.NodeOperatorId == 0 {
			return errors.New("nodeOperatorId cannot be zero")
		}
		if node.Signer == [32]byte{} {
			return errors.New("signer cannot be empty")
		}
		if node.EncryptionPublicKey == [32]byte{} {
			return errors.New("encryptionPublicKey cannot be empty")
		}
		if node.P2pId == [32]byte{} {
			return errors.New("p2pId cannot be empty")
		}
		if node.CsaKey == [32]byte{} {
			return errors.New("csaKey cannot be empty")
		}
		if len(node.CapabilityIds) == 0 {
			return errors.New("capabilityIds cannot be empty")
		}
	}
	return nil
}

func dedupNodes(lggr logger.Logger, inputNodes []NodesInput, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]NodesInput, error) {
	contractNodes, err := pkg.GetNodes(nil, capReg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes from contract: %w", err)
	}
	contractNodesMap := make(map[[32]byte]struct{})
	for _, node := range contractNodes {
		contractNodesMap[node.P2pId] = struct{}{}
	}

	var dedupedNodes []NodesInput
	for i, node := range inputNodes {
		if _, exists := contractNodesMap[node.P2pID]; exists {
			lggr.Infof("Node with P2P ID %s already registered in contract, skipping", hex.EncodeToString(node.P2pID[:]))
			continue
		}

		dedupedNodes = append(dedupedNodes, inputNodes[i])
	}

	return dedupedNodes, nil
}
