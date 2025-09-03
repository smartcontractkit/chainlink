package config

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func Set(t *testing.T, nodeInput *cre.CapabilitiesAwareNodeSet, bc *blockchain.Output) (*cre.WrappedNodeOutput, error) {
	nodeset, err := ns.UpgradeNodeSet(t, nodeInput.Input, bc, 5*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upgrade node set")
	}

	return &cre.WrappedNodeOutput{Output: nodeset, NodeSetName: nodeInput.Name, Capabilities: nodeInput.ComputedCapabilities}, nil
}

func Generate(input cre.GenerateConfigsInput, nodeConfigFns []cre.NodeConfigTransformerFn) (cre.NodeIndexToConfigOverride, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	configOverrides := make(cre.NodeIndexToConfigOverride)

	// if it's only a gateway DON, we don't need to generate any extra configuration, the default one will do
	if flags.HasFlag(input.Flags, cre.GatewayDON) && (!flags.HasFlag(input.Flags, cre.WorkflowDON) && !flags.HasFlag(input.Flags, cre.CapabilitiesDON)) {
		return configOverrides, nil
	}

	homeChainID, homeErr := chain_selectors.ChainIdFromSelector(input.HomeChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	// prepare chains, we need chainIDs and URLs
	evmChains := findEVMChains(input)
	solanaChain, solErr := findOneSolanaChain(input)
	if solErr != nil {
		return nil, errors.Wrap(solErr, "failed to find Solana chain in the environment configuration")
	}

	// find contract addresses
	capabilitiesRegistryAddress, capErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capErr != nil {
		return nil, errors.Wrap(capErr, "failed to find CapabilitiesRegistry address")
	}

	// find bootstrap node for the Don
	var donBootstrapNodeHost string
	var donBootstrapNodePeerID string

	bootstrapNodes, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap nodes")
	}

	switch len(bootstrapNodes) {
	case 0:
		// if DON doesn't have bootstrap node, we need to use the global bootstrap node
		donBootstrapNodeHost = input.OCRPeeringData.OCRBootstraperHost
		donBootstrapNodePeerID = input.OCRPeeringData.OCRBootstraperPeerID
	case 1:
		bootstrapNode := bootstrapNodes[0]

		donBootstrapNodePeerID, err = node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
		}

		for _, label := range bootstrapNode.Labels {
			if label.Key == node.HostLabelKey {
				donBootstrapNodeHost = label.Value
				break
			}
		}

		if donBootstrapNodeHost == "" {
			return nil, errors.New("failed to get bootstrap node host from labels")
		}

		var nodeIndex int
		for _, label := range bootstrapNode.Labels {
			if label.Key == node.IndexKey {
				nodeIndex, err = strconv.Atoi(label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
				break
			}
		}

		// generate configuration for the bootstrap node
		configOverrides[nodeIndex] = BootstrapEVM(donBootstrapNodePeerID, homeChainID, capabilitiesRegistryAddress, evmChains)
		if flags.HasFlag(input.Flags, cre.WorkflowDON) {
			configOverrides[nodeIndex] += BoostrapDon2DonPeering(input.CapabilitiesPeeringData)
		}
	default:
		return nil, errors.New("multiple bootstrap nodes within a DON found, expected only one")
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.IndexKey {
				nodeIndex, err = strconv.Atoi(label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
			}
		}

		// connect worker nodes to all the chains, add chain ID for registry (home chain)
		var workerErr error
		configOverrides[nodeIndex], workerErr = WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.OCRPeeringData, input.CapabilitiesPeeringData, capabilitiesRegistryAddress, homeChainID, evmChains)
		if workerErr != nil {
			return nil, errors.Wrap(workerErr, "failed to generate worker [EVM.Workflow] config")
		}
		configOverrides[nodeIndex] += WorkerSolana(solanaChain)
	}

	// execute capability-provided functions that transform the node config (currently: write-evm, write-solana)
	// these functions must return whole node configs after transforming them, instead of just returning configuration parts
	// that need to be merged into the existing config
	for _, configFn := range nodeConfigFns {
		if configFn == nil {
			continue
		}

		newOverrides, err := configFn(input, configOverrides)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}

		maps.Copy(configOverrides, newOverrides)
	}

	return configOverrides, nil
}

func findEVMChains(input cre.GenerateConfigsInput) []*EVMChain {
	evmChains := make([]*EVMChain, 0)
	for chainSelector, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain != nil {
			continue
		}

		// if the DON doesn't support the chain, we skip it; if slice is empty, it means that the DON supports all chains
		// TODO: review if we really need this SupportedChains functionality
		if len(input.DonMetadata.SupportedChains) > 0 && !slices.Contains(input.DonMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		evmChains = append(evmChains, &EVMChain{
			Name:    fmt.Sprintf("node-%d", chainSelector),
			ChainID: bcOut.ChainID,
			HTTPRPC: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			WSRPC:   bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
		})
	}
	return evmChains
}

func findOneSolanaChain(input cre.GenerateConfigsInput) (*SolanaChain, error) {
	var solanaChain *SolanaChain
	chainsFound := 0

	for _, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain == nil {
			continue
		}

		chainsFound++
		if chainsFound > 1 {
			return nil, errors.New("multiple Solana chains found, expected only one")
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		chainID, err := bcOut.SolClient.GetGenesisHash(ctx)
		if err != nil {
			cancelFn()
			return nil, errors.Wrap(err, "failed to get chainID for Solana")
		}
		cancelFn()

		solanaChain = &SolanaChain{
			Name:    fmt.Sprintf("node-%d", bcOut.SolChain.ChainSelector),
			ChainID: chainID.String(),
			NodeURL: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
		}
	}

	return solanaChain, nil
}
