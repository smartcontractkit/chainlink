package environment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
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
	tunnelManager tunnel.Manager,
) (*StartedDONs, error) {
	if infraInput.IsKubernetes() {
		// For Kubernetes, DONs are already running in the cluster, generate service URLs
		lggr.Info().Msg("Generating Kubernetes service URLs for DONs (already running in cluster)")
		for idx, nodeSet := range nodeSets {
			donMetadata := topology.DonsMetadata.List()[idx]

			// Extract bootstrap flags for each node
			nodeMetadataRoles := make([]bool, len(donMetadata.NodesMetadata))
			for i, nodeMeta := range donMetadata.NodesMetadata {
				nodeMetadataRoles[i] = nodeMeta.HasRole(cre.BootstrapNode)
			}

			creds := infra.GetNodeCredentials(&infraInput)
			nodeSet.Out = infra.GenerateKubernetesNodeSetOutput(&infraInput, nodeSet.Name, nodeSet.Nodes, nodeMetadataRoles, creds, lggr)
		}
	}

	// Skip binary operations for Kubernetes (binaries are in the cluster images)
	if infraInput.IsDocker() && !hasRemoteNodeSets(nodeSets) {
		// TODO in the future check here if don is remote and skip if it is instead of !hasRemoteNodeSets()
		for donIdx, donMetadata := range topology.DonsMetadata.List() {
			if !copyCapabilityBinaries {
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

			var err error
			ns, err := crecapabilities.AppendBinariesPathsNodeSpec(nodeSets[donIdx], donMetadata, customBinariesPaths)
			if err != nil {
				return nil, pkgerrors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
			}
			nodeSets[donIdx] = ns
		}
	}

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
			return nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}

	errGroup, _ := errgroup.WithContext(ctx)
	var resultMap sync.Map
	var startClient componentClient
	if hasRemoteNodeSets(nodeSets) {
		client, clientErr := newStartComponentClient(lggr, tunnelManager)
		if clientErr != nil {
			return nil, clientErr
		}
		startClient = client
	}

	for idx, nodeSet := range nodeSets {
		errGroup.Go(func() error {
			startTime := time.Now()
			lggr.Info().Msgf("Starting DON named %s", nodeSet.Name)

			var nodeset *ns.Output
			var nodesetErr error

			// If output is already set (Kubernetes or cached), use it
			if nodeSet.Out != nil {
				lggr.Info().Msgf("Using pre-configured node URLs for DON %s", nodeSet.Name)
				nodeset = nodeSet.Out
			} else if strings.TrimSpace(nodeSet.Target) == string(config.TargetRemote) {
				registryChainPayload, err := agent.EncodeForTransport(registryChainBlockchainOutput)
				if err != nil {
					return pkgerrors.Wrap(err, "failed to encode registry blockchain payload for remote nodeset start")
				}
				remoteInput, err := buildRemoteNodeSetInput(nodeSet)
				if err != nil {
					return err
				}
				payload := startComponentRequest{
					ComponentType:      componentTypeNodeSet,
					NodeSet:            remoteInput,
					RegistryBlockchain: registryChainPayload,
					ReusePolicy:        nodeSetRemoteStartPolicy(nodeSet),
				}
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return pkgerrors.Wrap(err, "failed to encode nodeset payload")
				}
				response, err := startClient.StartComponent(ctx, startComponentEnvelope{
					SchemaVersion: agent.SchemaVersionV1,
					Operation:     agent.OperationStartComponent,
					Payload:       payloadBytes,
				})
				if err != nil {
					return err
				}
				if response.ComponentType != componentTypeNodeSet {
					return fmt.Errorf("unexpected component type in start response: %s", response.ComponentType)
				}
				for _, logLine := range response.AgentLogs {
					pretty := prettifyAgentLogLine(logLine)
					if pretty == "" {
						continue
					}
					lggr.Info().Msgf("[agent] %s", pretty)
				}
				nodeset, err = agent.DecodeFromTransport[ns.Output](response.Output)
				if err != nil {
					return pkgerrors.Wrap(err, "failed to decode nodeset transport payload")
				}
				if err := rewriteRemoteNodeSetOutputForLocalAccess(ctx, lggr, tunnelManager, topology, idx, nodeSet, nodeset); err != nil {
					return err
				}
			} else {
				// For Docker, start the nodes
				nodeSet.Input.NodeSpecs = nodeSet.ExtractCTFInputs()
				nodeset, nodesetErr = ns.NewSharedDBNodeSetWithContext(ctx, nodeSet.Input, registryChainBlockchainOutput)
				if nodesetErr != nil {
					return pkgerrors.Wrapf(nodesetErr, "failed to start nodeSet named %s", nodeSet.Name)
				}
			}

			// For Kubernetes, we still need to create clients to register nodes with JD
			don, donErr := cre.NewDON(ctx, topology.DonsMetadata.List()[idx], nodeset.CLNodes)
			if donErr != nil {
				return pkgerrors.Wrapf(donErr, "failed to create DON from node set named %s", nodeSet.Name)
			}

			resultMap.Store(idx, &StartedDON{
				NodeSetOutput: &cre.NodeSetOutput{
					Output:       nodeset,
					NodeSetName:  nodeSet.Name,
					Capabilities: nodeSet.Capabilities,
				},
				DON: don,
			})

			lggr.Info().Msgf("DON %s started in %.2f seconds", nodeSet.Name, time.Since(startTime).Seconds())

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		if !infraInput.IsKubernetes() {
			infra.PrintFailedContainerLogs(lggr, 30)
		}
		return nil, err
	}

	startedDONs := make(StartedDONs, len(nodeSets))
	resultMap.Range(func(key, value any) bool {
		// key is index in the original slice
		startedDONs[key.(int)] = value.(*StartedDON)
		return true
	})

	return &startedDONs, nil
}

func hasRemoteNodeSets(nodeSets []*cre.NodeSet) bool {
	for _, nodeSet := range nodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Target) == string(config.TargetRemote) {
			return true
		}
	}
	return false
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

func rewriteRemoteNodeSetOutputForLocalAccess(
	ctx context.Context,
	lggr zerolog.Logger,
	tunnelManager tunnel.Manager,
	topology *cre.Topology,
	configuredIndex int,
	nodeSet *cre.NodeSet,
	output *ns.Output,
) error {
	if output == nil && (nodeSet == nil || nodeSet.DbInput == nil || nodeSet.DbInput.Port == 0) {
		return nil
	}
	if isRemoteAccessDirectMode() {
		hostIP, err := resolveDirectAccessHostIP()
		if err != nil {
			return err
		}
		if err := rewriteNodeSetForDirectAccess(output, hostIP); err != nil {
			return err
		}
		rewriteGatewayIncomingForDirectAccess(topology, configuredIndex, hostIP)
		return nil
	}
	componentID := tunnel.CanonicalComponentID(tunnel.KindNodeSet, configuredIndex, nodeSet.Name)
	refs, err := describeNodeSetEndpoints(componentID, nodeSet, output)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to describe nodeset tunnel endpoints")
	}
	bindings, err := tunnelManager.Start(ctx, refs)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to start tunnels for nodeset output")
	}
	for _, binding := range bindings {
		lggr.Info().
			Str("componentID", binding.ComponentID).
			Str("endpointName", binding.EndpointName).
			Str("originalURL", binding.OriginalURL).
			Str("localURL", binding.LocalURL).
			Msg("Established endpoint tunnel")
	}
	rewriteGatewayIncomingForNodeSetBindings(topology, configuredIndex, nodeSet, bindings)
	return rewriteNodeSetWithBindings(output, nodeSet, bindings)
}

func rewriteNodeSetForDirectAccess(output *ns.Output, hostIP string) error {
	if output == nil {
		return nil
	}
	for idx := range output.CLNodes {
		rawURL := output.CLNodes[idx].Node.ExternalURL
		if strings.TrimSpace(rawURL) == "" {
			continue
		}
		rewritten, err := rewriteURLHost(rawURL, hostIP)
		if err != nil {
			return err
		}
		output.CLNodes[idx].Node.ExternalURL = rewritten
	}
	return nil
}

const nodeSetDBEndpointName = "nodeset-db"

func describeNodeSetEndpoints(componentID string, nodeSet *cre.NodeSet, output *ns.Output) ([]tunnel.EndpointRef, error) {
	sizeHint := 1
	if output != nil {
		sizeHint += len(output.CLNodes)
	}
	if nodeSet != nil {
		for _, spec := range nodeSet.NodeSpecs {
			if spec == nil || spec.Node == nil {
				continue
			}
			sizeHint += len(spec.Node.CustomPorts)
		}
	}
	refs := make([]tunnel.EndpointRef, 0, sizeHint)
	if output != nil {
		for idx := range output.CLNodes {
			endpointName := fmt.Sprintf("node-%d-api", idx)
			rawURL := output.CLNodes[idx].Node.ExternalURL
			ref, err := nodeSetEndpointFromURL(componentID, endpointName, rawURL)
			if err != nil {
				return nil, err
			}
			if ref != nil {
				refs = append(refs, *ref)
			}
		}
	}
	if nodeSet != nil {
		for nodeIdx, spec := range nodeSet.NodeSpecs {
			customRefs, err := nodeSetCustomPortEndpointRefs(componentID, nodeIdx, spec)
			if err != nil {
				return nil, err
			}
			refs = append(refs, customRefs...)
		}
	}
	dbRef, err := nodeSetDBEndpointRef(componentID, nodeSet)
	if err != nil {
		return nil, err
	}
	if dbRef != nil {
		refs = append(refs, *dbRef)
	}
	return refs, nil
}

func nodeSetDBEndpointRef(componentID string, nodeSet *cre.NodeSet) (*tunnel.EndpointRef, error) {
	if nodeSet == nil || nodeSet.DbInput == nil || nodeSet.DbInput.Port == 0 {
		return nil, nil
	}
	if nodeSet.DbInput.Port < 0 || nodeSet.DbInput.Port > 65535 {
		return nil, fmt.Errorf("nodeset db port %d is invalid", nodeSet.DbInput.Port)
	}
	return &tunnel.EndpointRef{
		ComponentID:  componentID,
		EndpointName: nodeSetDBEndpointName,
		Scheme:       "tcp",
		Host:         "127.0.0.1",
		Port:         nodeSet.DbInput.Port,
		OriginalURL:  fmt.Sprintf("tcp://127.0.0.1:%d", nodeSet.DbInput.Port),
	}, nil
}

func rewriteNodeSetWithBindings(output *ns.Output, nodeSet *cre.NodeSet, bindings []tunnel.TunnelBinding) error {
	byName := make(map[string]tunnel.TunnelBinding, len(bindings))
	for _, binding := range bindings {
		byName[binding.EndpointName] = binding
	}
	if output != nil {
		for idx := range output.CLNodes {
			endpointName := fmt.Sprintf("node-%d-api", idx)
			rawURL := output.CLNodes[idx].Node.ExternalURL
			if rawURL == "" {
				continue
			}
			binding, ok := byName[endpointName]
			if !ok {
				return fmt.Errorf("missing tunnel binding for nodeset endpoint %s", endpointName)
			}
			output.CLNodes[idx].Node.ExternalURL = binding.LocalURL
		}
	}
	if nodeSet != nil && nodeSet.DbInput != nil && nodeSet.DbInput.Port != 0 {
		binding, ok := byName[nodeSetDBEndpointName]
		if !ok {
			return fmt.Errorf("missing tunnel binding for nodeset endpoint %s", nodeSetDBEndpointName)
		}
		nodeSet.DbInput.Port = binding.LocalPort
	}
	if nodeSet != nil {
		for nodeIdx, spec := range nodeSet.NodeSpecs {
			if spec == nil || spec.Input == nil || spec.Input.Node == nil || len(spec.Input.Node.CustomPorts) == 0 {
				continue
			}
			for portIdx, mapping := range spec.Input.Node.CustomPorts {
				_, containerPort, err := parseCustomPortMapping(mapping)
				if err != nil {
					return fmt.Errorf("invalid custom_ports entry %q for node %d: %w", mapping, nodeIdx, err)
				}
				binding, ok := byName[nodeSetCustomPortEndpointName(nodeIdx, portIdx, containerPort)]
				if !ok {
					return fmt.Errorf("missing tunnel binding for nodeset endpoint %s", nodeSetCustomPortEndpointName(nodeIdx, portIdx, containerPort))
				}
				spec.Input.Node.CustomPorts[portIdx] = rewriteCustomPortMappingHostPort(mapping, binding.LocalPort)
			}
		}
	}
	return nil
}

func nodeSetCustomPortEndpointRefs(componentID string, nodeIdx int, spec *cre.NodeSpecWithRole) ([]tunnel.EndpointRef, error) {
	if spec == nil || spec.Input == nil || spec.Input.Node == nil || len(spec.Input.Node.CustomPorts) == 0 {
		return nil, nil
	}
	refs := make([]tunnel.EndpointRef, 0, len(spec.Input.Node.CustomPorts))
	for portIdx, mapping := range spec.Input.Node.CustomPorts {
		hostPort, containerPort, err := parseCustomPortMapping(mapping)
		if err != nil {
			return nil, fmt.Errorf("invalid custom_ports entry %q for node %d: %w", mapping, nodeIdx, err)
		}
		refs = append(refs, tunnel.EndpointRef{
			ComponentID:  componentID,
			EndpointName: nodeSetCustomPortEndpointName(nodeIdx, portIdx, containerPort),
			Scheme:       "tcp",
			Host:         "127.0.0.1",
			Port:         hostPort,
			OriginalURL:  fmt.Sprintf("tcp://127.0.0.1:%d", hostPort),
		})
	}
	return refs, nil
}

func nodeSetCustomPortEndpointName(nodeIdx, portIdx, containerPort int) string {
	return fmt.Sprintf("node-%d-custom-%d-%d", nodeIdx, portIdx, containerPort)
}

func parseCustomPortMapping(mapping string) (hostPort int, containerPort int, err error) {
	parts := strings.Split(strings.TrimSpace(mapping), ":")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("expected hostPort:containerPort, got %q", mapping)
	}
	hostPortRaw := parts[len(parts)-2]
	containerPortRaw := parts[len(parts)-1]
	hostPort, err = strconv.Atoi(hostPortRaw)
	if err != nil || hostPort <= 0 || hostPort > 65535 {
		return 0, 0, fmt.Errorf("invalid host port %q", hostPortRaw)
	}
	containerPort, err = strconv.Atoi(containerPortRaw)
	if err != nil || containerPort <= 0 || containerPort > 65535 {
		return 0, 0, fmt.Errorf("invalid container port %q", containerPortRaw)
	}
	return hostPort, containerPort, nil
}

func rewriteCustomPortMappingHostPort(mapping string, newHostPort int) string {
	parts := strings.Split(strings.TrimSpace(mapping), ":")
	if len(parts) < 2 {
		return mapping
	}
	parts[len(parts)-2] = strconv.Itoa(newHostPort)
	return strings.Join(parts, ":")
}

func rewriteGatewayIncomingForNodeSetBindings(
	topology *cre.Topology,
	configuredIndex int,
	nodeSet *cre.NodeSet,
	bindings []tunnel.TunnelBinding,
) {
	if topology == nil || topology.GatewayConnectors == nil || len(topology.GatewayConnectors.Configurations) == 0 || nodeSet == nil {
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
	if gatewayNode.Index < 0 || gatewayNode.Index >= len(nodeSet.NodeSpecs) {
		return
	}
	spec := nodeSet.NodeSpecs[gatewayNode.Index]
	if spec == nil || spec.Input == nil || spec.Input.Node == nil || len(spec.Input.Node.CustomPorts) == 0 {
		return
	}

	for _, cfg := range topology.GatewayConnectors.Configurations {
		if cfg == nil || cfg.GatewayConfiguration == nil || cfg.NodeUUID != gatewayNode.UUID {
			continue
		}
		// Test process reaches gateway via local port (direct for local runs, tunneled for remote runs).
		cfg.Incoming.Host = "127.0.0.1"
		// Resolve tunnel by gateway container port (e.g. 5002), not by possibly stale host-side custom port.
		if localPort, ok := gatewayLocalPortFromBindings(gatewayNode.Index, cfg.Incoming.ExternalPort, bindings); ok {
			cfg.Incoming.ExternalPort = localPort
		}
	}
}

func rewriteGatewayIncomingForDirectAccess(topology *cre.Topology, configuredIndex int, hostIP string) {
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
		cfg.Incoming.Host = hostIP
	}
}

func gatewayLocalPortFromBindings(gatewayNodeIndex, gatewayContainerPort int, bindings []tunnel.TunnelBinding) (int, bool) {
	for _, binding := range bindings {
		if !strings.HasPrefix(binding.EndpointName, fmt.Sprintf("node-%d-custom-", gatewayNodeIndex)) {
			continue
		}
		if strings.HasSuffix(binding.EndpointName, fmt.Sprintf("-%d", gatewayContainerPort)) {
			return binding.LocalPort, true
		}
	}
	return 0, false
}

func nodeSetEndpointFromURL(componentID, endpointName, rawURL string) (*tunnel.EndpointRef, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint url %q: %w", rawURL, err)
	}
	host := parsed.Hostname()
	if host == "" {
		return nil, fmt.Errorf("endpoint url %q has empty hostname", rawURL)
	}
	port, err := nodeSetResolveURLPort(parsed)
	if err != nil {
		return nil, err
	}
	return &tunnel.EndpointRef{
		ComponentID:  componentID,
		EndpointName: endpointName,
		Scheme:       parsed.Scheme,
		Host:         host,
		Port:         port,
		OriginalURL:  rawURL,
	}, nil
}

func nodeSetResolveURLPort(parsed *url.URL) (int, error) {
	if parsed.Port() != "" {
		port, err := strconv.Atoi(parsed.Port())
		if err != nil || port <= 0 || port > 65535 {
			return 0, fmt.Errorf("url %q has invalid port %q", parsed.String(), parsed.Port())
		}
		return port, nil
	}
	switch parsed.Scheme {
	case "http", "ws":
		return 80, nil
	case "https", "wss":
		return 443, nil
	default:
		return 0, fmt.Errorf("url %q has unsupported scheme %q without explicit port", parsed.String(), parsed.Scheme)
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
