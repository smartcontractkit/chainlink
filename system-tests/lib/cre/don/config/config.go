package config

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func Set(t *testing.T, nodeInput *types.CapabilitiesAwareNodeSet, bc *blockchain.Output) (*types.WrappedNodeOutput, error) {
	nodeset, err := ns.UpgradeNodeSet(t, nodeInput.Input, bc, 5*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upgrade node set")
	}

	return &types.WrappedNodeOutput{Output: nodeset, NodeSetName: nodeInput.Name, Capabilities: nodeInput.Capabilities}, nil
}

func Generate(input types.GenerateConfigsInput, factoryFns []types.ConfigFactoryFn) (types.NodeIndexToConfigOverride, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	configOverrides := make(types.NodeIndexToConfigOverride)

	// if it's only a gateway DON, we don't need to generate any extra configuration, the default one will do
	if flags.HasFlag(input.Flags, types.GatewayDON) && (!flags.HasFlag(input.Flags, types.WorkflowDON) && !flags.HasFlag(input.Flags, types.CapabilitiesDON)) {
		return configOverrides, nil
	}

	homeChainID, homeErr := chain_selectors.ChainIdFromSelector(input.HomeChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	// prepare chains, we need chainIDs, URLs and selectors to get contracts from AddressBook
	workerEVMInputs := make([]*WorkerEVMInput, 0)
	for chainSelector, bcOut := range input.BlockchainOutput {
		cID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
		if err != nil {
			return configOverrides, errors.Wrapf(err, "failed to parse chain ID %s", bcOut.ChainID)
		}
		c, exists := chain_selectors.ChainByEvmChainID(cID)
		if !exists {
			return configOverrides, errors.Errorf("failed to find selector for chain ID %d", cID)
		}
		workerEVMInputs = append(workerEVMInputs, &WorkerEVMInput{
			Name:          fmt.Sprintf("node-%d", chainSelector),
			ChainID:       bcOut.ChainID,
			ChainSelector: c.Selector,
			HTTPRPC:       bcOut.Nodes[0].InternalHTTPUrl,
			WSRPC:         bcOut.Nodes[0].InternalWSUrl,
		})
	}

	// find contract addresses
	capabilitiesRegistryAddress, capErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capErr != nil {
		return nil, errors.Wrap(capErr, "failed to find CapabilitiesRegistry address")
	}

	// find bootstrap node for the Don
	var donBootstrapNodeHost string
	var donBootstrapNodePeerID string

	bootstrapNodes, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.BootstrapNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap nodes")
	}

	switch len(bootstrapNodes) {
	case 0:
		// if DON doesn't have bootstrap node, we need to use the global bootstrap node
		donBootstrapNodeHost = input.PeeringData.GlobalBootstraperHost
		donBootstrapNodePeerID = input.PeeringData.GlobalBootstraperPeerID
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
		configOverrides[nodeIndex] = BootstrapEVM(donBootstrapNodePeerID, homeChainID, capabilitiesRegistryAddress, workerEVMInputs)

		if flags.HasFlag(input.Flags, types.WorkflowDON) {
			configOverrides[nodeIndex] += BoostrapDon2DonPeering(input.PeeringData)
		}
	default:
		return nil, errors.New("multiple bootstrap nodes within a DON found, expected only one")
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
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

		// get all the forwarders and add workflow config for each node ETH key + Forwarder for that chain
		for _, wi := range workerEVMInputs {
			addrsForChains, err := input.AddressBook.AddressesForChain(wi.ChainSelector)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get addresses from address book")
			}
			for addr, addrValue := range addrsForChains {
				if addrValue.Type == "KeystoneForwarder" {
					wi.ForwarderAddress = addr
					expectedAddressKey := node.AddressKeyFromSelector(wi.ChainSelector)
					for _, label := range workflowNodeSet[i].Labels {
						if label.Key == expectedAddressKey {
							if label.Value == "" {
								return nil, errors.Errorf("%s label value is empty", expectedAddressKey)
							}
							wi.FromAddress = common.HexToAddress(label.Value)
							break
						}
					}
					if wi.FromAddress == (common.Address{}) {
						return nil, errors.Errorf("failed to get from address for chain %d", wi.ChainSelector)
					}
				}
			}
		}

		// connect worker nodes to all the chains, add chain ID for registry (home chain)
		// we configure both EVM chains, nodes and EVM.Workflow with Forwarder
		configOverrides[nodeIndex] = WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.PeeringData, capabilitiesRegistryAddress, homeChainID, workerEVMInputs)
	}

	for _, factoryFn := range factoryFns {
		newOverrides, err := factoryFn(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}
		for nodeIndex, override := range newOverrides {
			configOverrides[nodeIndex] += override
		}
	}

	return configOverrides, nil
}
