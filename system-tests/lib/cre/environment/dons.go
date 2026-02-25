package environment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type StartedDON struct {
	NodeSetOutput *cre.NodeSetOutput
	DON           *cre.Don
}

type StartedDONs []*StartedDON

func (s *StartedDONs) NodeOutputs() []*cre.NodeSetOutput {
	outputs := make([]*cre.NodeSetOutput, len(*s))
	for idx, don := range *s {
		outputs[idx] = don.NodeSetOutput
	}
	return outputs
}

func (s *StartedDONs) DONs() []*cre.Don {
	dons := make([]*cre.Don, len(*s))
	for idx, don := range *s {
		dons[idx] = don.DON
	}
	return dons
}

func StartDONs(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	infraInput infra.Provider,
	registryChainBlockchainOutput *blockchain.Output,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
	nodeSets []*cre.NodeSet,
	remoteRuntime *resolvedRemoteRuntime,
) (*StartedDONs, error) {
	if err := verifyRemoteToLocalBootstrapReachability(ctx, lggr, topology); err != nil {
		return nil, pkgerrors.Wrap(err, "bootstrap reachability sanity check failed")
	}

	switch {
	case infraInput.IsKubernetes():
		return startDONsKubernetes(ctx, lggr, topology, infraInput, nodeSets)
	default:
		return startDONsContainerized(
			ctx,
			lggr,
			topology,
			infraInput,
			registryChainBlockchainOutput,
			capabilityConfigs,
			copyCapabilityBinaries,
			nodeSets,
			remoteRuntime,
		)
	}
}

func startDONsKubernetes(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	infraInput infra.Provider,
	nodeSets []*cre.NodeSet,
) (*StartedDONs, error) {
	lggr.Info().Msg("Generating Kubernetes service URLs for DONs (already running in cluster)")
	for idx, nodeSet := range nodeSets {
		donMetadata := topology.DonsMetadata.List()[idx]

		// Extract bootstrap flags for each node.
		nodeMetadataRoles := make([]bool, len(donMetadata.NodesMetadata))
		for i, nodeMeta := range donMetadata.NodesMetadata {
			nodeMetadataRoles[i] = nodeMeta.HasRole(cre.BootstrapNode)
		}

		creds := infra.GetNodeCredentials(&infraInput)
		nodeSet.Out = infra.GenerateKubernetesNodeSetOutput(&infraInput, nodeSet.Name, nodeSet.Nodes, nodeMetadataRoles, creds, lggr)
	}
	if err := applyNodeSetEnvVars(topology, nodeSets); err != nil {
		return nil, err
	}

	return buildDONsConcurrently(ctx, lggr, false, nodeSets, func(configuredIndex int, configuredNodeSet *cre.NodeSet) (*StartedDON, error) {
		lggr.Info().Msgf("Kubernetes mode: using existing DON named %s", configuredNodeSet.Name)
		return buildStartedDON(ctx, topology, configuredIndex, configuredNodeSet, configuredNodeSet.Out)
	})
}

func startDONsContainerized(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	infraInput infra.Provider,
	registryChainBlockchainOutput *blockchain.Output,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
	nodeSets []*cre.NodeSet,
	remoteRuntime *resolvedRemoteRuntime,
) (*StartedDONs, error) {
	// Skip binary operations for remote DONs.
	if infraInput.IsDocker() {
		for donIdx, donMetadata := range topology.DonsMetadata.List() {
			if !copyCapabilityBinaries {
				continue
			}
			if donMetadata.MustNodeSet().Placement == string(config.PlacementRemote) {
				continue
			}

			customBinariesPaths := make(map[cre.CapabilityFlag]string)
			for flag, config := range capabilityConfigs {
				if flags.HasFlagForAnyChain(donMetadata.Flags, flag) && config.BinaryPath != "" {
					customBinariesPaths[flag] = config.BinaryPath
				}
			}

			executableErr := crecapabilities.MakeBinariesExecutable(customBinariesPaths)
			if executableErr != nil {
				return nil, pkgerrors.Wrap(executableErr, "failed to make binaries executable")
			}

			ns, err := crecapabilities.AppendBinariesPathsNodeSpec(nodeSets[donIdx], donMetadata, customBinariesPaths)
			if err != nil {
				return nil, pkgerrors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
			}
			nodeSets[donIdx] = ns
		}
	}
	if err := applyNodeSetEnvVars(topology, nodeSets); err != nil {
		return nil, err
	}

	return buildDONsConcurrently(ctx, lggr, true, nodeSets, func(configuredIndex int, configuredNodeSet *cre.NodeSet) (*StartedDON, error) {
		return startDON(
			ctx,
			lggr,
			topology,
			configuredIndex,
			configuredNodeSet,
			registryChainBlockchainOutput,
			remoteRuntime,
		)
	})
}

func applyNodeSetEnvVars(topology *cre.Topology, nodeSets []*cre.NodeSet) error {
	// Add env vars, which were provided programmatically, to the node specs
	// or fail, if node specs already had some env vars set in the TOML config
	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		hasEnvVarsInTomlConfig := false
		for nodeIdx, nodeSpec := range nodeSets[donIdx].NodeSpecs {
			if len(nodeSpec.Node.EnvVars) > 0 {
				hasEnvVarsInTomlConfig = true
				break
			}

			nodeSets[donIdx].NodeSpecs[nodeIdx].Node.EnvVars = nodeSets[donIdx].EnvVars
		}

		if hasEnvVarsInTomlConfig && len(nodeSets[donIdx].EnvVars) > 0 {
			return fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}
	return nil
}

func buildDONsConcurrently(
	ctx context.Context,
	lggr zerolog.Logger,
	printFailedContainerLogs bool,
	nodeSets []*cre.NodeSet,
	startFn func(configuredIndex int, configuredNodeSet *cre.NodeSet) (*StartedDON, error),
) (*StartedDONs, error) {
	errGroup, _ := errgroup.WithContext(ctx)
	startedDONs := make(StartedDONs, len(nodeSets))

	for idx, nodeSet := range nodeSets {
		configuredIndex := idx
		configuredNodeSet := nodeSet
		errGroup.Go(func() error {
			startedDON, startErr := startFn(configuredIndex, configuredNodeSet)
			if startErr != nil {
				return startErr
			}
			startedDONs[configuredIndex] = startedDON
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		if printFailedContainerLogs {
			infra.PrintFailedContainerLogs(lggr, 30)
		}
		return nil, err
	}

	return &startedDONs, nil
}

func startDON(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	configuredIndex int,
	nodeSet *cre.NodeSet,
	registryChainBlockchainOutput *blockchain.Output,
	remoteRuntime *resolvedRemoteRuntime,
) (*StartedDON, error) {
	if nodeSet == nil {
		return nil, errors.New("nodeSet is nil")
	}
	startTime := time.Now()
	lggr.Info().Msgf("Starting DON named %s", nodeSet.Name)

	nodeset, err := startNodeSet(ctx, lggr, topology, configuredIndex, nodeSet, registryChainBlockchainOutput, remoteRuntime)
	if err != nil {
		return nil, err
	}

	startedDON, buildErr := buildStartedDON(ctx, topology, configuredIndex, nodeSet, nodeset)
	if buildErr != nil {
		return nil, buildErr
	}

	lggr.Info().Msgf("DON %s started in %.2f seconds", nodeSet.Name, time.Since(startTime).Seconds())
	return startedDON, nil
}

func buildStartedDON(
	ctx context.Context,
	topology *cre.Topology,
	configuredIndex int,
	nodeSet *cre.NodeSet,
	nodeset *ns.Output,
) (*StartedDON, error) {
	if nodeSet == nil {
		return nil, errors.New("nodeSet is nil")
	}
	if nodeset == nil {
		return nil, fmt.Errorf("nodeSet output is nil for DON %s", nodeSet.Name)
	}

	donsMetadata := topology.DonsMetadata.List()
	if configuredIndex < 0 || configuredIndex >= len(donsMetadata) {
		return nil, fmt.Errorf("configured index %d out of bounds for dons metadata", configuredIndex)
	}
	don, donErr := cre.NewDON(ctx, donsMetadata[configuredIndex], nodeset.CLNodes)
	if donErr != nil {
		return nil, pkgerrors.Wrapf(donErr, "failed to create DON from node set named %s", nodeSet.Name)
	}

	return &StartedDON{
		NodeSetOutput: &cre.NodeSetOutput{
			Output:       nodeset,
			NodeSetName:  nodeSet.Name,
			Capabilities: nodeSet.Capabilities,
		},
		DON: don,
	}, nil
}
func startNodeSet(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	configuredIndex int,
	nodeSet *cre.NodeSet,
	registryChainBlockchainOutput *blockchain.Output,
	remoteRuntime *resolvedRemoteRuntime,
) (*ns.Output, error) {
	// If output is already set (Kubernetes or cached), use it.
	if nodeSet.Out != nil {
		lggr.Info().Msgf("Using pre-configured node URLs for DON %s", nodeSet.Name)
		return nodeSet.Out, nil
	}

	if strings.TrimSpace(nodeSet.Placement) == string(config.PlacementRemote) {
		if remoteRuntime == nil {
			return nil, errors.New("remote runtime is required for remote nodeset placement")
		}
		registryChainPayload, err := agent.EncodeForTransport(registryChainBlockchainOutput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to encode registry blockchain payload for remote nodeset start")
		}
		remoteInput, err := buildRemoteNodeSetInput(nodeSet)
		if err != nil {
			return nil, err
		}
		payload := agent.StartComponentPayload{
			ComponentType:      componentTypeNodeSet,
			NodeSet:            remoteInput,
			RegistryBlockchain: registryChainPayload,
			ReusePolicy:        nodeSetRemoteStartPolicy(nodeSet),
		}
		nodeset, err := startRemoteComponent[ns.Output](
			ctx,
			lggr,
			remoteRuntime.Client,
			payload,
			componentTypeNodeSet,
		)
		if err != nil {
			return nil, err
		}
		if err := rewriteRemoteNodeSetOutputForLocalAccess(topology, configuredIndex, nodeSet, nodeset, remoteRuntime.EC2HostIP); err != nil {
			return nil, err
		}
		return nodeset, nil
	}

	// For Docker, start the nodes.
	nodeSet.Input.NodeSpecs = nodeSet.ExtractCTFInputs()
	nodeset, err := ns.NewSharedDBNodeSetWithContext(ctx, nodeSet.Input, registryChainBlockchainOutput)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to start nodeSet named %s", nodeSet.Name)
	}
	return nodeset, nil
}

func nodeSetRemoteStartPolicy(nodeSet *cre.NodeSet) string {
	if nodeSet == nil || strings.TrimSpace(nodeSet.RemoteStartPolicy) == "" {
		return string(config.RemoteStartPolicyReuseIfIdentical)
	}
	return nodeSet.RemoteStartPolicy
}

func buildRemoteNodeSetInput(nodeSet *cre.NodeSet) (*ns.Input, error) {
	if nodeSet == nil || nodeSet.Input == nil {
		return nil, pkgerrors.New("nodeset input is nil for remote target")
	}
	inputCopy := *nodeSet.Input
	inputCopy.NodeSpecs = nodeSet.ExtractCTFInputs()
	if err := validateRemoteNodeSetNodeSpecs(inputCopy.Name, inputCopy.NodeSpecs); err != nil {
		return nil, err
	}
	return &inputCopy, nil
}

func validateRemoteNodeSetNodeSpecs(nodeSetName string, specs []*clnode.Input) error {
	for idx, spec := range specs {
		if spec == nil || spec.Node == nil {
			return fmt.Errorf("remote nodeset %q node_specs[%d] is nil", nodeSetName, idx)
		}
		hasImage := strings.TrimSpace(spec.Node.Image) != ""
		hasBuildConfig := strings.TrimSpace(spec.Node.DockerContext) != "" ||
			strings.TrimSpace(spec.Node.DockerFilePath) != "" ||
			len(spec.Node.DockerBuildArgs) > 0
		if hasImage && hasBuildConfig {
			return fmt.Errorf(
				"remote nodeset %q node_specs[%d] must configure either node.image or docker build fields (docker_ctx/docker_file/docker_build_args), not both",
				nodeSetName,
				idx,
			)
		}
		if !hasImage && !hasBuildConfig {
			return fmt.Errorf(
				"remote nodeset %q node_specs[%d] must set node.image or docker build fields (docker_ctx/docker_file/docker_build_args)",
				nodeSetName,
				idx,
			)
		}
	}
	return nil
}

func rewriteRemoteNodeSetOutputForLocalAccess(topology *cre.Topology, configuredIndex int, nodeSet *cre.NodeSet, output *ns.Output, ec2HostIP string) error {
	if output == nil && (nodeSet == nil || nodeSet.DbInput == nil || nodeSet.DbInput.Port == 0) {
		return nil
	}
	if err := rewriteNodeSetForDirectAccess(output, ec2HostIP); err != nil {
		return err
	}
	rewriteGatewayIncomingForDirectAccess(topology, configuredIndex, ec2HostIP)
	return nil
}

func rewriteNodeSetForDirectAccess(output *ns.Output, ec2HostIP string) error {
	if output == nil {
		return nil
	}
	for idx := range output.CLNodes {
		rawURL := output.CLNodes[idx].Node.ExternalURL
		if strings.TrimSpace(rawURL) == "" {
			continue
		}
		rewritten, err := rewriteURLHost(rawURL, ec2HostIP)
		if err != nil {
			return err
		}
		output.CLNodes[idx].Node.ExternalURL = rewritten
	}
	return nil
}

func rewriteGatewayIncomingForDirectAccess(topology *cre.Topology, configuredIndex int, ec2HostIP string) {
	if topology == nil || topology.GatewayConnectors == nil || len(topology.GatewayConnectors.Configurations) == 0 {
		return
	}
	if configuredIndex < 0 || configuredIndex >= len(topology.DonsMetadata.List()) {
		return
	}
	donMeta := topology.DonsMetadata.List()[configuredIndex]
	gatewayNode, hasGateway := donMeta.Gateway()
	if !hasGateway {
		return
	}
	for _, cfg := range topology.GatewayConnectors.Configurations {
		if cfg == nil || cfg.GatewayConfiguration == nil || cfg.NodeUUID != gatewayNode.UUID {
			continue
		}
		cfg.Incoming.Host = ec2HostIP
	}
}

func FundNodes(ctx context.Context, testLogger zerolog.Logger, dons *cre.Dons, blockchains []blockchains.Blockchain, fundingAmountPerChainFamily map[string]uint64) error {
	for _, don := range dons.List() {
		testLogger.Info().Msgf("Funding nodes for DON %s", don.Name)
		for _, bc := range blockchains {
			if !flags.RequiresForwarderContract(don.Flags, bc.ChainID()) && !bc.IsFamily(chainselectors.FamilySolana) { // for now, we can only write to solana, so we consider forwarder is always present
				continue
			}

			chainFamily := bc.CtfOutput().Family
			fundingAmount, ok := fundingAmountPerChainFamily[chainFamily]
			if !ok {
				return fmt.Errorf("missing funding amount for chain family %s", chainFamily)
			}

			for _, node := range don.Nodes {
				address, addrErr := nodeAddress(node, chainFamily, bc)
				if addrErr != nil {
					return pkgerrors.Wrapf(addrErr, "failed to get address for node %s on chain family %s and chain %d", node.Name, chainFamily, bc.ChainID())
				}

				if address == "" {
					testLogger.Info().Msgf("No key for chainID %d found for node %s. Skipping funding", bc.ChainID(), node.Name)
					continue // Skip nodes without keys for this chain
				}

				err := bc.Fund(ctx, address, fundingAmount)
				if err != nil {
					return err
				}
			}
		}

		testLogger.Info().Msgf("Funded nodes for DON %s", don.Name)
	}

	return nil
}

func nodeAddress(node *cre.Node, chainFamily string, bc blockchains.Blockchain) (string, error) {
	switch chainFamily {
	case chainselectors.FamilyEVM, chainselectors.FamilyTron:
		evmKey, ok := node.Keys.EVM[bc.ChainID()]
		if !ok {
			return "", nil // Skip nodes without EVM keys for this chain
		}

		return evmKey.PublicAddress.String(), nil
	case chainselectors.FamilySolana:
		solBc := bc.(*solana.Blockchain)
		solKey, ok := node.Keys.Solana[solBc.SolanaChainID]
		if !ok {
			return "", nil // Skip nodes without Solana keys for this chain
		}
		return solKey.PublicAddress.String(), nil
	default:
		return "", fmt.Errorf("unsupported chain family %s", chainFamily)
	}
}
