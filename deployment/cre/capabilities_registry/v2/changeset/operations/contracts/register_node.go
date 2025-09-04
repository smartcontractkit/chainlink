package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type RegisterNodesDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type RegisterNodesInput struct {
	Address       string
	ChainSelector uint64
	Nodes         []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	MCMSConfig    *ocr3.MCMSConfig
}

type RegisterNodesOutput struct {
	Nodes     []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded
	Proposals []mcmslib.TimelockProposal
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
				Nodes: []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded{},
			}, nil
		}
		if input.ChainSelector == 0 {
			return RegisterNodesOutput{}, errors.New("chainSelector is not set")
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

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(input.Address),
			RegisterNodesDescription,
		)
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var resultNodes []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := capabilityRegistryTransactor.AddNodes(opts, input.Nodes)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call AddNodes: %w", err)
			}

			// For direct execution, we can get the receipt and parse logs
			if input.MCMSConfig == nil {
				// Confirm transaction and get receipt
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm AddNodes transaction %s: %w", tx.Hash().String(), err)
				}

				ctx := b.GetContext()
				receipt, err := bind.WaitMined(ctx, chain.Client, tx)
				if err != nil {
					return nil, fmt.Errorf("failed to mine AddNodes transaction %s: %w", tx.Hash().String(), err)
				}

				// Get the CapabilitiesRegistryFilterer contract for parsing logs
				capabilityRegistryFilterer, err := capabilities_registry_v2.NewCapabilitiesRegistryFilterer(
					common.HexToAddress(input.Address),
					chain.Client,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create CapabilitiesRegistryFilterer: %w", err)
				}

				// Parse the logs to get the added nodes
				resultNodes = make([]*capabilities_registry_v2.CapabilitiesRegistryNodeAdded, 0, len(receipt.Logs))
				for i, log := range receipt.Logs {
					if log == nil {
						continue
					}

					o, err := capabilityRegistryFilterer.ParseNodeAdded(*log)
					if err != nil {
						return nil, fmt.Errorf("failed to parse log %d for node added: %w", i, err)
					}
					resultNodes = append(resultNodes, o)
				}
			}

			return tx, nil
		})
		if err != nil {
			return RegisterNodesOutput{}, fmt.Errorf("failed to execute AddNodes: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterNodes on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d nodes on chain %d", len(resultNodes), input.ChainSelector)
		}

		return RegisterNodesOutput{
			Nodes:     resultNodes,
			Proposals: proposals,
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
