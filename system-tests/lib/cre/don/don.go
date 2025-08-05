package don

import (
	"context"
	"slices"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

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
						// we need to call the DonID "vault" because it is used in two-fold manner:
						// - to authenticate the caller with the gateway, and since each node can only have 1 gateway connector, it uses the same DonID for all gateways.
						// - to specify which handler should be used to handle request (for "vault" it needs to be "vault", for "web-api" anything else)
						// And that introduces an unfortunate cupling. If the node is connected to "vault" gateway, then "DonID" can be only be "vault".
						topology.GatewayConnectorOutput = initGatewayConnectorOutput("vault")
					}

					topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, &cre.GatewayConfiguration{
						Outgoing: cre.Outgoing{
							Path: "/node",
							Port: 5003,
							Host: gatewayInternalHost,
						},
						Incoming: cre.Incoming{
							Protocol:     "http",
							Path:         "/",
							InternalPort: 5002,
							ExternalPort: ExternalGatewayPort(infraInput, coregateway.WebAPICapabilitiesType),
							Host:         ExternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput),
						},
						HandlerType:   coregateway.WebAPICapabilitiesType,
						AuthGatewayID: "web-api-gateway",
						// do not set gateway connector dons, they will be resolved automatically
					})

					if AnyDonHasCapability(donsWithMetadata, cre.VaultCapability) {
						topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, &cre.GatewayConfiguration{
							Outgoing: cre.Outgoing{
								Path: "/node",
								Port: 15003,
								Host: gatewayInternalHost,
							},
							Incoming: cre.Incoming{
								Protocol:     "http",
								Path:         "/",
								InternalPort: 15002,
								ExternalPort: ExternalGatewayPort(infraInput, coregateway.VaultHandlerType),
								Host:         ExternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput),
							},
							HandlerType:   coregateway.VaultHandlerType,
							AuthGatewayID: "vault-gateway",
							// do not set gateway connector dons, they will be resolved automatically
						})
					}
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
		if slices.Contains(don.Flags, capability) {
			return true
		}
	}

	return false
}

func NodeNeedsGateway(handlerType coregateway.HandlerType, nodeFlags []cre.CapabilityFlag) bool {
	switch handlerType {
	case coregateway.WebAPICapabilitiesType:
		return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
			flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
			flags.HasFlag(nodeFlags, cre.WebAPITargetCapability)
	case coregateway.VaultHandlerType:
		return flags.HasFlag(nodeFlags, cre.VaultCapability)
	}
	return false
}

func GatewayConfigurationsForHandler(handlerType coregateway.HandlerType, gatewayConnectorOutput *cre.GatewayConnectorOutput) []*cre.GatewayConfiguration {
	gatewayConfigurations := make([]*cre.GatewayConfiguration, 0)
	for _, configuration := range gatewayConnectorOutput.Configurations {
		if configuration.HandlerType == handlerType {
			gatewayConfigurations = append(gatewayConfigurations, configuration)
		}
	}

	return gatewayConfigurations
}

func initGatewayConnectorOutput(donID string) *cre.GatewayConnectorOutput {
	return &cre.GatewayConnectorOutput{
		DonID:          donID,
		Configurations: make([]*cre.GatewayConfiguration, 0),
	}
}
