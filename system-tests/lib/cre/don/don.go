package don

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/credentials/insecure"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	ctf_jd "github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
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
	if len(nodeSetInput) == 0 {
		return errors.New("at least one nodeset is required")
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
		if nodeSet.BootstrapNodeIndex != -1 && slices.Contains(nodeSet.DONTypes, cre.WorkflowDON) {
			workflowDONHasBootstrapNode = true
			break
		}
	}

	if !workflowDONHasBootstrapNode {
		return errors.New("due to the limitations of our implementation, workflow DON must always have a bootstrap node")
	}

	isGatewayRequired := false
	for _, nodeSet := range nodeSetInput {
		if NodeNeedsAnyGateway(nodeSet.ComputedCapabilities) {
			isGatewayRequired = true
			break
		}
	}

	if !isGatewayRequired {
		return nil
	}

	anyDONHasGatewayConfigured := false
	for _, nodeSet := range nodeSetInput {
		if isGatewayRequired {
			if flags.HasFlag(nodeSet.DONTypes, cre.GatewayDON) && nodeSet.GatewayNodeIndex != -1 {
				anyDONHasGatewayConfigured = true
				break
			}
		}
	}

	if !anyDONHasGatewayConfigured {
		return errors.New("at least one DON must be configured with gateway DON type and have a gateway node index set, because at least one DON requires gateway due to its capabilities")
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
			ID:              libc.MustSafeUint64FromInt(i + 1), // optimistically set the id to the that which the capabilities registry will assign it
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

func LinkToJobDistributor(ctx context.Context, input *cre.LinkDonsToJDInput) (*cldf.Environment, []*devenv.DON, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	dons := make([]*devenv.DON, len(input.NodeSetOutput))
	var allNodesInfo []devenv.NodeInfo

	for idx, nodeOutput := range input.NodeSetOutput {
		bootstrapNodes, err := node.FindManyWithLabel(input.Topology.DonsMetadata[idx].NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		nodeInfo, err := node.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, input.Topology.DonsMetadata[idx].ID, len(bootstrapNodes))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		supportedChains, schErr := findSupportedChainsForDON(input.Topology.DonsMetadata[idx], input.BlockchainOutputs)
		if schErr != nil {
			return nil, nil, errors.Wrap(schErr, "failed to find supported chains for DON")
		}

		var regErr error
		dons[idx], regErr = configureJDForDON(ctx, nodeInfo, supportedChains, input.JdOutput)
		if regErr != nil {
			return nil, nil, fmt.Errorf("failed to configure JD for DON: %w", regErr)
		}
	}

	var nodeIDs []string
	for _, don := range dons {
		nodeIDs = append(nodeIDs, don.NodeIds()...)
	}

	dons = addOCRKeyLabelsToNodeMetadata(dons, input.Topology)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	jd, jdErr := devenv.NewJDClient(ctxWithTimeout, devenv.JDConfig{
		GRPC:     input.JdOutput.ExternalGRPCUrl,
		WSRPC:    input.JdOutput.InternalWSRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: allNodesInfo,
	})

	if jdErr != nil {
		return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	input.CldfEnvironment.Offchain = jd
	input.CldfEnvironment.NodeIDs = nodeIDs

	return input.CldfEnvironment, dons, nil
}

func configureJDForDON(ctx context.Context, nodeInfo []devenv.NodeInfo, supportedChains []devenv.ChainConfig, jdOutput *ctf_jd.Output) (*devenv.DON, error) {
	jdConfig := jd.JDConfig{
		GRPC:  jdOutput.ExternalGRPCUrl,
		WSRPC: jdOutput.InternalWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	jdClient, jdErr := jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	donJDClient := &devenv.JobDistributor{
		JobDistributor: jdClient,
	}

	don, regErr := devenv.NewRegisteredDON(ctx, nodeInfo, *donJDClient)
	if regErr != nil {
		return nil, fmt.Errorf("failed to create registered DON: %w", regErr)
	}

	if err := don.CreateSupportedChains(ctx, supportedChains, *donJDClient); err != nil {
		return nil, fmt.Errorf("failed to create supported chains: %w", err)
	}

	return don, nil
}

func findSupportedChainsForDON(donMetadata *cre.DonMetadata, blockchainOutputs []*cre.WrappedBlockchainOutput) ([]devenv.ChainConfig, error) {
	chains := make([]devenv.ChainConfig, 0)

	for chainSelector, bcOut := range blockchainOutputs {
		if len(donMetadata.SupportedChains) > 0 && !slices.Contains(donMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
		}

		chains = append(chains, cfg)
	}

	return chains, nil
}

func addOCRKeyLabelsToNodeMetadata(dons []*devenv.DON, topology *cre.Topology) []*devenv.DON {
	for i, don := range dons {
		for j, donNode := range topology.DonsMetadata[i].NodesMetadata {
			// required for job proposals, because they need to include the ID of the node in Job Distributor
			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   node.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			ocrSupportedFamilies := make([]string, 0)
			for family, key := range don.Nodes[j].ChainsOcr2KeyBundlesID {
				donNode.Labels = append(donNode.Labels, &cre.Label{
					Key:   node.CreateNodeOCR2KeyBundleIDKey(family),
					Value: key,
				})
				ocrSupportedFamilies = append(ocrSupportedFamilies, family)
			}

			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   node.NodeOCRFamiliesKey,
				Value: node.CreateNodeOCRFamiliesListValue(ocrSupportedFamilies),
			})
		}
	}

	return dons
}
