package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type RegisterNodesDeps struct {
	Env *cldf.Environment
}

type RegisterNodesInput struct {
	Address       string
	ChainSelector uint64
	Nodes         []capabilities_registry_v2.CapabilitiesRegistryNodeParams
}

type RegisterNodesOutput struct {
	Nodes []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded
}

// RegisterNodes is an operation that registers nodes in the V2 Capabilities Registry contract.
var RegisterNodes = operations.NewOperation[RegisterNodesInput, RegisterNodesOutput, RegisterNodesDeps](
	"register-nodes-op",
	semver.MustParse("1.0.0"),
	"Register Nodes in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNodesDeps, input RegisterNodesInput) (RegisterNodesOutput, error) {
		// Validate input
		if input.Address == "" {
			return RegisterNodesOutput{}, errors.New("Address is not set")
		}
		if len(input.Nodes) == 0 {
			return RegisterNodesOutput{}, errors.New("Nodes are not set")
		}
		if input.ChainSelector == 0 {
			return RegisterNodesOutput{}, errors.New("ChainSelector is not set")
		}

		if err := validateNodes(input.Nodes); err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("node validation failed: %w", err)
		}

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterNodesOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistryTransactor(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		tx, err := capabilityRegistryTransactor.AddNodes(chain.DeployerKey, input.Nodes)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterNodesOutput{}, fmt.Errorf("failed to call AddNodes: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to confirm AddNodes confirm transaction %s: %w", tx.Hash().String(), err)
		}

		ctx := b.GetContext()
		receipt, err := bind.WaitMined(ctx, chain.Client, tx)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to mine AddNodes confirm transaction %s: %w", tx.Hash().String(), err)
		}
		if len(receipt.Logs) != len(input.Nodes) {
			return RegisterNodesOutput{}, fmt.Errorf("expected %d log entries for AddNodes, got %d", len(input.Nodes), len(receipt.Logs))
		}

		// Get the CapabilitiesRegistryFilterer contract
		capabilityRegistryFilterer, err := capabilities_registry_v2.NewCapabilitiesRegistryFilterer(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryFilterer: %w", err)
		}

		resp := RegisterNodesOutput{
			Nodes: make([]*capabilities_registry_v2.CapabilitiesRegistryNodeAdded, 0, len(receipt.Logs)),
		}
		// Parse the logs to get the added nodes
		for i, log := range receipt.Logs {
			if log == nil {
				continue
			}

			o, err := capabilityRegistryFilterer.ParseNodeAdded(*log)
			if err != nil {
				return RegisterNodesOutput{}, fmt.Errorf("failed to parse log %d for node added: %w", i, err)
			}
			resp.Nodes = append(resp.Nodes, o)
		}

		return resp, nil
	},
)

func validateNodes(nodes []capabilities_registry_v2.CapabilitiesRegistryNodeParams) error {
	for _, node := range nodes {
		if node.NodeOperatorId == 0 {
			return errors.New("NodeOperatorId cannot be zero")
		}
		if node.Signer == [32]byte{} {
			return errors.New("Signer cannot be empty")
		}
		if node.EncryptionPublicKey == [32]byte{} {
			return errors.New("EncryptionPublicKey cannot be empty")
		}
		if node.P2pId == [32]byte{} {
			return errors.New("P2pId cannot be empty")
		}
		if node.CsaKey == [32]byte{} {
			return errors.New("CsaKey cannot be empty")
		}
		if len(node.CapabilityIds) == 0 {
			return errors.New("CapabilityIds cannot be empty")
		}
	}
	return nil
}
