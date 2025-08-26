package config

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func Set(t *testing.T, nodeInput *cre.CapabilitiesAwareNodeSet, bc *blockchain.Output) (*cre.WrappedNodeOutput, error) {
	nodeset, err := ns.UpgradeNodeSet(t, nodeInput.Input, bc, 5*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upgrade node set")
	}

	return &cre.WrappedNodeOutput{Output: nodeset, NodeSetName: nodeInput.Name, Capabilities: nodeInput.ComputedCapabilities}, nil
}

func Generate(input cre.GenerateConfigsInput, nodeConfigFns []cre.NodeConfigFn) (cre.NodeIndexToConfigOverride, error) {
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

	// prepare chains, we need chainIDs, URLs and selectors to get contracts from AddressBook
	workerEVMInputs := make([]*WorkerEVMInput, 0)
	workerSolInputs := make([]*WorkerSolanaInput, 0)
	for chainSelector, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain != nil {
			chainID, err := bcOut.SolClient.GetGenesisHash(context.Background())
			if err != nil {
				return nil, errors.Wrap(err, "failed to get chainID for Solana")
			}

			// Determine write-solana enablement per chain via node-set ChainCapabilities
			hasWrite := false
			hasWrite = slices.Contains(input.NodeSet.Capabilities, cre.WriteSolanaCapability)

			workerSolInputs = append(workerSolInputs, &WorkerSolanaInput{
				ChainSelector: bcOut.SolChain.ChainSelector,
				Name:          fmt.Sprintf("node-%d", bcOut.SolChain.ChainSelector),
				ChainID:       chainID.String(),
				NodeURL:       bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
				HasWrite:      hasWrite,
			})

			continue
		}
		// if the DON doesn't support the chain, we skip it; if slice is empty, it means that the DON supports all chains
		if len(input.DonMetadata.SupportedChains) > 0 && !slices.Contains(input.DonMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		c, exists := chain_selectors.ChainByEvmChainID(bcOut.ChainID)
		if !exists {
			return configOverrides, errors.Errorf("failed to find selector for chain ID %d", bcOut.ChainID)
		}
		// Determine write-evm enablement per chain via node-set ChainCapabilities
		hasWriteEVM := false
		if input.NodeSet != nil && input.NodeSet.ChainCapabilities != nil {
			if cc, ok := input.NodeSet.ChainCapabilities[cre.WriteEVMCapability]; ok && cc != nil {
				if slices.Contains(cc.EnabledChains, bcOut.ChainID) {
					hasWriteEVM = true
				}
			}
		}
		workerEVMInputs = append(workerEVMInputs, &WorkerEVMInput{
			Name:          fmt.Sprintf("node-%d", chainSelector),
			ChainID:       bcOut.ChainID,
			ChainSelector: c.Selector,
			HTTPRPC:       bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			WSRPC:         bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
			WritesToEVM:   hasWriteEVM,
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
		configOverrides[nodeIndex] = BootstrapEVM(donBootstrapNodePeerID, homeChainID, capabilitiesRegistryAddress, workerEVMInputs)
		if flags.HasFlag(input.Flags, cre.WorkflowDON) {
			configOverrides[nodeIndex] += BoostrapDon2DonPeering(input.CapabilitiesPeeringData)
		}
		if len(workerSolInputs) > 0 {
			configOverrides[nodeIndex] += BootstrapSolana(workerSolInputs)
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

		// get all the forwarders and add workflow config (FromAddress + Forwarder) for chains that have write-evm enabled
		for _, wi := range workerEVMInputs {
			if !wi.WritesToEVM {
				continue
			}

			addrsForChains, err := input.AddressBook.AddressesForChain(wi.ChainSelector)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get addresses from address book")
			}
			for addr, addrValue := range addrsForChains {
				if addrValue.Type == keystone_changeset.KeystoneForwarder {
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

			if input.CapabilityConfigs == nil {
				return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-evm capability")
			}

			if writeEvmConfig, ok := input.CapabilityConfigs[cre.WriteEVMCapability]; ok {
				enabled, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
					cre.WriteEVMCapability,
					input.NodeSet.ChainCapabilities,
					writeEvmConfig.Config,
					wi.ChainID,
				)
				if rErr != nil {
					return nil, errors.Wrapf(rErr, "failed to resolve write-evm config for chain %d", wi.ChainID)
				}

				if !enabled {
					// This should never happen, but guard anyway. We have already checked that the capability is enabled in the chain capabilities, when we generated the workerEVMInputs.
					continue
				}

				runtimeValues := map[string]any{
					"FromAddress":      wi.FromAddress.Hex(),
					"ForwarderAddress": wi.ForwarderAddress,
				}

				var mErr error
				wi.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
				if mErr != nil {
					return nil, errors.Wrap(mErr, "failed to apply runtime values")
				}
			}
		}

		// get all sol forwarders
		for _, wi := range workerSolInputs {
			if !wi.HasWrite {
				continue
			}
			forwarders := input.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(wi.ChainSelector))
			for _, addr := range forwarders {
				if addr.Type == ks_sol.ForwarderState {
					wi.ForwarderState = addr.Address
					continue
				}
				expectedAddressKey := node.AddressKeyFromSelector(wi.ChainSelector)
				wi.ForwarderAddress = addr.Address
				for _, label := range workflowNodeSet[i].Labels {
					if label.Key == expectedAddressKey {
						if label.Value == "" {
							return nil, errors.Errorf("%s label value is empty", expectedAddressKey)
						}
						wi.FromAddress = solana.MustPublicKeyFromBase58(label.Value)
						break
					}
				}
				if wi.FromAddress.IsZero() {
					return nil, errors.Errorf("failed to get from address for Solana chain %d", wi.ChainSelector)
				}
			}
			if input.CapabilityConfigs == nil {
				return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-evm capability")
			}

			if writeSolConfig, ok := input.CapabilityConfigs[cre.WriteSolanaCapability]; ok {
				mergedConfig := envconfig.ResolveCapabilityConfigForDON(
					cre.WriteSolanaCapability,
					writeSolConfig.Config,
					nil,
				)

				runtimeValues := map[string]any{
					"FromAddress":      wi.FromAddress.String(),
					"ForwarderAddress": wi.ForwarderAddress,
					"ForwarderState":   wi.ForwarderState,
				}

				var mErr error
				wi.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
				if mErr != nil {
					return nil, errors.Wrap(mErr, "failed to apply runtime values")
				}
			}
		}

		// connect worker nodes to all the chains, add chain ID for registry (home chain)
		// we configure both EVM chains, nodes and EVM.Workflow with Forwarder
		var workerErr error
		configOverrides[nodeIndex], workerErr = WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.OCRPeeringData, input.CapabilitiesPeeringData, capabilitiesRegistryAddress, homeChainID, workerEVMInputs)
		if workerErr != nil {
			return nil, errors.Wrap(workerErr, "failed to generate worker [EVM.Workflow] config")
		}
		solOverride, solWorkerErr := WorkerSolana(workerSolInputs)
		if solWorkerErr != nil {
			return nil, errors.Wrap(workerErr, "failed to generate worker [Solana.Workflow] config")
		}

		configOverrides[nodeIndex] += solOverride
	}

	for _, configFn := range nodeConfigFns {
		if configFn == nil {
			continue
		}
		newOverrides, err := configFn(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}
		for nodeIndex, override := range newOverrides {
			configOverrides[nodeIndex] += override
		}
	}

	return configOverrides, nil
}
