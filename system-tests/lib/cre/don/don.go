package don

import (
	"context"
	"slices"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cre.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func ValidateTopology(nodeSetInput []*cre.CapabilitiesAwareNodeSet, infraInput infra.Input) error {
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
		if nodeSet.BootstrapNodeIndex != -1 && slices.Contains(nodeSet.DONTypes, cre.WorkflowDON) {
			workflowDONHasBootstrapNode = true
			break
		}
	}

	if !workflowDONHasBootstrapNode {
		return errors.New("due to the limitations of our implementation, workflow DON must always have a bootstrap node")
	}

	return nil
}

func BuildTopology(nodeSetInput []*cre.CapabilitiesAwareNodeSet, infraInput infra.Input, homeChainSelector uint64) (*cre.Topology, error) {
	topology := &cre.Topology{}
	donsWithMetadata := make([]*cre.DonMetadata, len(nodeSetInput))

	for i := range nodeSetInput {
		flags, err := flags.NodeSetFlags(nodeSetInput[i])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get flags for nodeset %s", nodeSetInput[i].Name)
		}

		donsWithMetadata[i] = &cre.DonMetadata{
			ID:              libc.MustSafeUint64FromInt(i + 1),
			Flags:           flags,
			NodesMetadata:   make([]*cre.NodeMetadata, len(nodeSetInput[i].NodeSpecs)),
			Name:            nodeSetInput[i].Name,
			SupportedChains: nodeSetInput[i].SupportedChains,
		}
	}

	for donIdx, donMetadata := range donsWithMetadata {
		for nodeIdx := range donMetadata.NodesMetadata {
			nodeWithLabels := cre.NodeMetadata{}
			nodeType := cre.WorkerNode
			if nodeSetInput[donIdx].BootstrapNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].BootstrapNodeIndex {
				nodeType = cre.BootstrapNode
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.NodeTypeKey,
				Value: nodeType,
			})

			// TODO think whether it would make sense for infraInput to also hold functions that resolve hostnames for various infra and node types
			// and use it with some default, so that we can easily modify it with little effort
			internalHost := InternalHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

			if flags.HasFlag(donMetadata.Flags, cre.GatewayDON) {
				if nodeSetInput[donIdx].GatewayNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].GatewayNodeIndex {
					nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
						Key:   node.ExtraRolesKey,
						Value: cre.GatewayNode,
					})

					gatewayInternalHost := InternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

					if topology.GatewayConnectorOutput == nil {
						topology.GatewayConnectorOutput = &cre.GatewayConnectorOutput{
							Configurations: make([]*cre.GatewayConfiguration, 0),
						}
					}

					topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, &cre.GatewayConfiguration{
						Outgoing: cre.Outgoing{
							Path: "/node",
							Port: GatewayOutgoingPort,
							Host: gatewayInternalHost,
						},
						Incoming: cre.Incoming{
							Protocol:     "http",
							Path:         "/",
							InternalPort: GatewayIncomingPort,
							ExternalPort: ExternalGatewayPort(infraInput),
							Host:         ExternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput),
						},
						AuthGatewayID: "cre-gateway",
						// do not set gateway connector dons, they will be resolved automatically
					})
				}
			}

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.IndexKey,
				Value: strconv.Itoa(nodeIdx),
			})

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.HostLabelKey,
				Value: internalHost,
			})

			donsWithMetadata[donIdx].NodesMetadata[nodeIdx] = &nodeWithLabels
		}
	}

	maybeID, err := flags.OneDonMetadataWithFlag(donsWithMetadata, cre.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	topology.DonsMetadata = donsWithMetadata
	topology.WorkflowDONID = maybeID.ID
	topology.HomeChainSelector = homeChainSelector

	return topology, nil
}

func AnyDonHasCapability(donMetadata []*cre.DonMetadata, capability cre.CapabilityFlag) bool {
	for _, don := range donMetadata {
		if flags.HasFlagForAnyChain(don.Flags, capability) {
			return true
		}
	}

	return false
}

func NodeNeedsAnyGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability) ||
		flags.HasFlag(nodeFlags, cre.VaultCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPActionCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPTriggerCapability)
}

func NodeNeedsWebAPIGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability)
}
