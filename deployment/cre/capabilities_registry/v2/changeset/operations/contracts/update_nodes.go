package contracts

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type UpdateNodesDeps struct {
	Env                  *cldf.Environment
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry
	MCMSContracts        *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type NodeConfig struct {
	EncryptionPublicKey string
	NodeOperatorID      uint32
	Signer              [32]byte
	CSAKey              string

	Capabilities []capabilities_registry_v2.CapabilitiesRegistryCapability
}

type UpdateNodesInput struct {
	ChainSelector uint64

	// NodesUpdates is a map of p2p key to NodeConfig
	NodesUpdates map[string]NodeConfig
	MCMSConfig   *ocr3.MCMSConfig
}

type UpdateNodesOutput struct {
	UpdatedNodes []*capabilities_registry_v2.CapabilitiesRegistryNodeUpdated
	Proposals    []mcmslib.TimelockProposal
}

var UpdateNodes = operations.NewOperation[UpdateNodesInput, UpdateNodesOutput, UpdateNodesDeps](
	"update-nodes-op",
	semver.MustParse("1.0.0"),
	"Update Nodes in Capabilities Registry",
	func(b operations.Bundle, deps UpdateNodesDeps, input UpdateNodesInput) (UpdateNodesOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return UpdateNodesOutput{}, cldf.ErrChainNotFound
		}

		nodeParams, err := makeNodeParams(deps.CapabilitiesRegistry, input.NodesUpdates)
		if err != nil {
			return UpdateNodesOutput{}, fmt.Errorf("failed to make node params: %w", err)
		}

		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			deps.CapabilitiesRegistry.Address(),
			UpdateNodesDescription,
		)
		if err != nil {
			return UpdateNodesOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var resultNodes []*capabilities_registry_v2.CapabilitiesRegistryNodeUpdated

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.CapabilitiesRegistry.UpdateNodes(opts, nodeParams)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call UpdateNodes: %w", err)
			}

			// For direct execution, we can get the receipt and parse logs
			if input.MCMSConfig == nil {
				// Confirm transaction and get receipt
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm UpdateNodes transaction %s: %w", tx.Hash().String(), err)
				}

				ctx := b.GetContext()
				receipt, err := bind.WaitMined(ctx, chain.Client, tx)
				if err != nil {
					return nil, fmt.Errorf("failed to mine UpdateNodes transaction %s: %w", tx.Hash().String(), err)
				}

				// Parse the logs to get the updated nodes
				resultNodes = make([]*capabilities_registry_v2.CapabilitiesRegistryNodeUpdated, 0, len(receipt.Logs))
				for i, log := range receipt.Logs {
					if log == nil {
						continue
					}

					o, err := deps.CapabilitiesRegistry.ParseNodeUpdated(*log)
					if err != nil {
						return nil, fmt.Errorf("failed to parse log %d for node updated: %w", i, err)
					}
					resultNodes = append(resultNodes, o)
				}
			}

			return tx, nil
		})
		if err != nil {
			return UpdateNodesOutput{}, fmt.Errorf("failed to execute UpdateNodes: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpdateNodes on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully updated %d nodes on chain %d", len(resultNodes), input.ChainSelector)
		}

		return UpdateNodesOutput{
			UpdatedNodes: resultNodes,
			Proposals:    proposals,
		}, nil
	},
)

func makeNodeParams(
	registry *capabilities_registry_v2.CapabilitiesRegistry,
	p2pToUpdates map[string]NodeConfig,
) ([]capabilities_registry_v2.CapabilitiesRegistryNodeParams, error) {
	var p2pIDs []p2pkey.PeerID
	for p2pIDStr := range p2pToUpdates {
		p2pID, err := p2pkey.MakePeerID(p2pIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p id %s: %w", p2pIDStr, err)
		}
		p2pIDs = append(p2pIDs, p2pID)
	}

	var out []capabilities_registry_v2.CapabilitiesRegistryNodeParams

	nodes, err := registry.GetNodesByP2PIds(&bind.CallOpts{}, pkg.PeerIDsToBytes(p2pIDs))
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to get nodes by p2p ids: %w", err)
	}

	for _, node := range nodes {
		p2pIDStr := p2pkey.PeerID(node.P2pId).String()
		updates, ok := p2pToUpdates[p2pIDStr]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node %s", p2pIDStr)
		}

		// We merge the already existing capabilities IDs with the new ones, to make sure that capabilities required by the DON
		// are still supported.
		ids := node.CapabilityIds
		for _, capability := range updates.Capabilities {
			ids = append(ids, capability.CapabilityId)
		}

		encryptionKey := node.EncryptionPublicKey
		if updates.EncryptionPublicKey != "" {
			pk, err := hex.DecodeString(updates.EncryptionPublicKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode encryption public key: %w", err)
			}
			encryptionKey = [32]byte(pk)
		}

		signer := node.Signer
		var zero [32]byte
		if !bytes.Equal(updates.Signer[:], zero[:]) {
			signer = updates.Signer
		}

		nodeOperatorID := node.NodeOperatorId
		if updates.NodeOperatorID != 0 {
			nodeOperatorID = updates.NodeOperatorID
		}

		csaKey := node.CsaKey
		if updates.CSAKey != "" {
			k, err := hex.DecodeString(updates.CSAKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode csa key: %w", err)
			}
			csaKey = [32]byte(k)
		}

		out = append(out, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      nodeOperatorID,
			P2pId:               node.P2pId,
			CapabilityIds:       ids,
			EncryptionPublicKey: encryptionKey,
			Signer:              signer,
			CsaKey:              csaKey,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].NodeOperatorId == out[j].NodeOperatorId {
			return bytes.Compare(out[i].P2pId[:], out[j].P2pId[:]) < 0
		}
		return out[i].NodeOperatorId < out[j].NodeOperatorId
	})

	return out, nil
}
