package don

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func CreateJobs(testLogger zerolog.Logger, input cretypes.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(input.CldEnv.Offchain, don.DON, don.Flags, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func ValidateTopology(nodeSetInput []*cretypes.CapabilitiesAwareNodeSet, infraInput types.InfraInput) error {
	if infraInput.InfraType == types.CRIB {
		if len(nodeSetInput) == 1 && slices.Contains(nodeSetInput[0].DONTypes, cretypes.GatewayDON) {
			if len(nodeSetInput[0].Capabilities) > 1 {
				return errors.New("you must use at least 2 nodeSets when using CRIB and gateway DON. Gateway DON must be in a separate nodeSet and it must be named 'gateway'")
			}
		}

		for _, nodeSet := range nodeSetInput {
			if infraInput.InfraType == types.CRIB && slices.Contains(nodeSetInput[0].DONTypes, cretypes.GatewayDON) && nodeSet.Name != "gateway" {
				return errors.New("when using CRIB gateway nodeSet with the Gateway DON must be named 'gateway', but got " + nodeSet.Name)
			}
		}
	}

	hasAtLeastOneBootstrapNode := false
	for _, nodeSet := range nodeSetInput {
		if nodeSet.BootstrapNodeIndex != -1 {
			hasAtLeastOneBootstrapNode = true
			break
		}
	}

	if !hasAtLeastOneBootstrapNode {
		return errors.New("at least one nodeSet must have a bootstrap node")
	}

	workflowDONHasBootstrapNode := false
	for _, nodeSet := range nodeSetInput {
		if nodeSet.BootstrapNodeIndex != -1 && slices.Contains(nodeSet.DONTypes, cretypes.WorkflowDON) {
			workflowDONHasBootstrapNode = true
			break
		}
	}

	if !workflowDONHasBootstrapNode {
		return errors.New("due to the limitations of our implementation, workflow DON must always have a bootstrap node")
	}

	return nil
}

func BuildTopology(nodeSetInput []*cretypes.CapabilitiesAwareNodeSet, infraInput types.InfraInput) (*cretypes.Topology, error) {
	topology := &cretypes.Topology{}
	donsWithMetadata := make([]*cretypes.DonMetadata, len(nodeSetInput))

	for i := range nodeSetInput {
		flags, err := flags.NodeSetFlags(nodeSetInput[i])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get flags for nodeset %s", nodeSetInput[i].Name)
		}

		donsWithMetadata[i] = &cretypes.DonMetadata{
			ID:            libc.MustSafeUint32(i + 1),
			Flags:         flags,
			NodesMetadata: make([]*cretypes.NodeMetadata, len(nodeSetInput[i].NodeSpecs)),
			Name:          nodeSetInput[i].Name,
		}
	}

	for donIdx, donMetadata := range donsWithMetadata {
		for nodeIdx := range donMetadata.NodesMetadata {
			nodeWithLabels := cretypes.NodeMetadata{}
			nodeType := cretypes.WorkerNode
			if nodeSetInput[donIdx].BootstrapNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].BootstrapNodeIndex {
				nodeType = cretypes.BootstrapNode
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.NodeTypeKey,
				Value: nodeType,
			})

			// TODO think whether it would make sense for infraInput to also hold functions that resolve hostnames for various infra and node types
			// and use it with some default, so that we can easily modify it with little effort
			host := infra.Host(nodeIdx, nodeType, donMetadata.Name, infraInput)

			if flags.HasFlag(donMetadata.Flags, cretypes.GatewayDON) {
				if nodeSetInput[donIdx].GatewayNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].GatewayNodeIndex {
					nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
						Key:   node.ExtraRolesKey,
						Value: cretypes.GatewayNode,
					})

					gatewayHost := host
					if infraInput.InfraType == types.CRIB {
						gatewayHost += "-gtwnode"
					}

					topology.GatewayConnectorOutput = &cretypes.GatewayConnectorOutput{
						Path: "/node",
						Port: 5003,
						Host: gatewayHost,
						// do not set gateway connector dons, they will be resolved automatically
					}
				}
			}

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.IndexKey,
				Value: strconv.Itoa(nodeIdx),
			})

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.HostLabelKey,
				Value: host,
			})

			donsWithMetadata[donIdx].NodesMetadata[nodeIdx] = &nodeWithLabels
		}
	}

	maybeID, err := flags.OneDonMetadataWithFlag(donsWithMetadata, cretypes.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	topology.DonsMetadata = donsWithMetadata
	topology.WorkflowDONID = maybeID.ID

	return topology, nil
}

func AddKeysToTopology(topology *cretypes.Topology, keys *cretypes.GenerateKeysOutput) (*cretypes.Topology, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	if keys == nil {
		return nil, errors.New("keys is nil")
	}

	if len(keys.P2PKeys) != len(topology.DonsMetadata) {
		return nil, fmt.Errorf("number of P2P keys does not match the number of DONs. Expected %d, got %d", len(topology.DonsMetadata), len(keys.P2PKeys))
	}

	if len(keys.EVMKeys) != len(topology.DonsMetadata) {
		return nil, fmt.Errorf("number of EVM keys does not match the number of DONs. Expected %d, got %d", len(topology.DonsMetadata), len(keys.EVMKeys))
	}

	for _, donMetadata := range topology.DonsMetadata {
		if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
			if len(p2pKeys.PeerIDs) != len(donMetadata.NodesMetadata) {
				return nil, fmt.Errorf("number of P2P keys for DON %d does not match the number of nodes. Expected %d, got %d", donMetadata.ID, len(donMetadata.NodesMetadata), len(p2pKeys.PeerIDs))
			}
			for idx, nodeMetadata := range donMetadata.NodesMetadata {
				nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
					Key:   node.NodeP2PIDKey,
					Value: p2pKeys.PeerIDs[idx],
				})
			}
		} else {
			return nil, fmt.Errorf("no P2P keys found for DON %d", donMetadata.ID)
		}

		if chainIDsToEVMKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
			// First, verify that all chain IDs have the same EVM keys for each node
			// This is a limitation of our current testing SDK implementation, because we only have 1 label for ETH address for each node
			// If in the future we need to support multiple chain IDs, we will need to change this and prefix the label with the chain ID
			// For now let's just make sure that all the chain IDs have the same EVM keys for each node to avoid hard to debug issues
			var firstChainID int
			var firstEVMKeys *types.EVMKeys
			for chainID, evmKeys := range chainIDsToEVMKeys {
				if firstEVMKeys == nil {
					firstChainID = chainID
					firstEVMKeys = evmKeys
					continue
				}
				if len(evmKeys.PublicAddresses) != len(firstEVMKeys.PublicAddresses) {
					return nil, fmt.Errorf("number of EVM keys for DON %d differs between chain IDs %d and %d", donMetadata.ID, firstChainID, chainID)
				}
				for i := range evmKeys.PublicAddresses {
					if evmKeys.PublicAddresses[i] != firstEVMKeys.PublicAddresses[i] {
						return nil, fmt.Errorf("EVM public address mismatch for DON %d, node %d between chain IDs %d and %d", donMetadata.ID, i, firstChainID, chainID)
					}
				}
			}

			// Now add the EVM addresses to the node metadata
			for chainID, evmKeys := range chainIDsToEVMKeys {
				if len(evmKeys.PublicAddresses) != len(donMetadata.NodesMetadata) {
					return nil, fmt.Errorf("number of EVM keys for DON %d and chain ID %d does not match the number of nodes. Expected %d, got %d", donMetadata.ID, chainID, len(donMetadata.NodesMetadata), len(evmKeys.PublicAddresses))
				}
				for idx, nodeMetadata := range donMetadata.NodesMetadata {
					nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
						Key:   node.EthAddressKey,
						Value: evmKeys.PublicAddresses[idx].Hex(),
					})
				}
				// Use first ETH address for the DON metadata, because all of them are the same
				break
			}
		} else {
			return nil, fmt.Errorf("no EVM keys found for DON %d", donMetadata.ID)
		}
	}

	return topology, nil
}

func GenereteKeys(input *cretypes.GenerateKeysInput) (*cretypes.GenerateKeysOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	if input.Out != nil {
		return input.Out, nil
	}

	output := &cretypes.GenerateKeysOutput{
		EVMKeys: make(cretypes.DonsToEVMKeys),
		P2PKeys: make(cretypes.DonsToP2PKeys),
	}

	for _, donMetadata := range input.Topology.DonsMetadata {
		if input.GenerateP2PKeys {
			p2pKeys, err := crypto.GenerateP2PKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate P2P keys")
			}
			output.P2PKeys[donMetadata.ID] = p2pKeys
		}

		if len(input.GenerateEVMKeysForChainIDs) > 0 {
			evmKeys, err := crypto.GenerateEVMKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate EVM keys")
			}

			// use the same EVM keys for all the chain IDs
			for _, chainID := range input.GenerateEVMKeysForChainIDs {
				if _, ok := output.EVMKeys[donMetadata.ID]; !ok {
					output.EVMKeys[donMetadata.ID] = make(cretypes.ChainIDToEVMKeys)
				}
				output.EVMKeys[donMetadata.ID][chainID] = evmKeys
			}
		}
	}

	return output, nil
}

func NodeNeedsGateway(nodeFlags []cretypes.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cretypes.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cretypes.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cretypes.WebAPITargetCapability)
}
