package don

import (
	"context"
	"slices"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cretypes.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, don.DON, don.Flags, jobSpecs)
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
				return errors.New("you must use at least 2 nodeSets when using CRIB and gateway DON. Gateway DON must be in a separate nodeSet and it must be named 'gateway'. Try using 'full' topology by passing '-t full' to the CLI")
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

func BuildTopology(nodeSetInput []*cretypes.CapabilitiesAwareNodeSet, infraInput types.InfraInput, homeChainSelector uint64) (*cretypes.Topology, error) {
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
			internalHost := infra.InternalHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

			if flags.HasFlag(donMetadata.Flags, cretypes.GatewayDON) {
				if nodeSetInput[donIdx].GatewayNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].GatewayNodeIndex {
					nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
						Key:   node.ExtraRolesKey,
						Value: cretypes.GatewayNode,
					})

					gatewayInternalHost := infra.InternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

					topology.GatewayConnectorOutput = &cretypes.GatewayConnectorOutput{
						Outgoing: cretypes.Outgoing{
							Path: "/node",
							Port: 5003,
							Host: gatewayInternalHost,
						},
						Incoming: cretypes.Incoming{
							Protocol:     "http",
							Path:         "/",
							InternalPort: 5002,
							ExternalPort: infra.ExternalGatewayPort(infraInput),
							Host:         infra.ExternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput),
						},
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
				Value: internalHost,
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
	topology.HomeChainSelector = homeChainSelector

	return topology, nil
}

func NodeNeedsGateway(nodeFlags []cretypes.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cretypes.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cretypes.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cretypes.WebAPITargetCapability)
}
