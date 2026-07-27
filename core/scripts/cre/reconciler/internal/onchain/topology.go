package onchain

import (
	"context"
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// buildTopology constructs cre.Topology from desired state, chart values, and
// discovered node runtime info.
func (d *Deployer) buildTopology(ctx context.Context, desired *domain.DesiredState, cv *domain.ChartValues, state *domain.StateFile) (*cre.Topology, error) {
	if d.k8s == nil {
		return nil, errors.New("k8s client is required to build topology")
	}

	provider := infra.Provider{Type: infra.Kubernetes}
	globalCapConfigs := d.buildGlobalCapabilityConfigs(desired)
	supportedEVMChains := desired.ChainIDs()
	registryChainVal, ok := desired.RegistryChain()
	if !ok {
		return nil, errors.New("no registry chain declared in desired state (exactly one [[chains]] entry must set registry = true)")
	}
	registryChain := registryChainVal.ChainID

	var nodeSets []*cre.NodeSet
	nodeNamesBySet := make([][]string, 0)
	requiredEVMChainsBySet := make([][]uint64, 0)
	bootstrapAssigned := false

	for i := range desired.DONs {
		don := &desired.DONs[i]

		if don.IsGatewayDon() {
			continue
		}

		var (
			nodeSet        *cre.NodeSet
			nodeNames      []string
			requiredChains []uint64
			err            error
		)
		if don.IsBootstrapOnly(cv) {
			nodeSet, nodeNames, requiredChains, err = d.buildBootstrapOnlyNodeSet(
				ctx, don, supportedEVMChains, registryChain, &bootstrapAssigned, cv, state,
			)
		} else {
			nodeSet, nodeNames, err = d.buildNodeSetForDON(ctx, don, supportedEVMChains, registryChain, &bootstrapAssigned, cv, state)
			requiredChains = requiredEVMChainIDsForDON(don, supportedEVMChains, registryChain)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build nodeset for DON %s", don.Name)
		}
		if id, ok := state.DONIDs[don.Name]; ok {
			nodeSet.ContractDonID = id
		} else {
			nodeSet.ContractDonID = uint64(i + 1)
		}

		nodeSets = append(nodeSets, nodeSet)
		nodeNamesBySet = append(nodeNamesBySet, nodeNames)
		requiredEVMChainsBySet = append(requiredEVMChainsBySet, requiredChains)
	}

	if desired.NeedsGateway() {
		gwNodeSets, gwNodeNamesBySet, err := d.buildGatewayNodeSets(ctx, desired, cv, state, supportedEVMChains, registryChain)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build gateway nodesets")
		}
		for i, gwNodeSet := range gwNodeSets {
			nodeSets = append(nodeSets, gwNodeSet)
			nodeNamesBySet = append(nodeNamesBySet, gwNodeNamesBySet[i])
			requiredEVMChainsBySet = append(requiredEVMChainsBySet, requiredEVMChainIDsForGatewayNode(registryChain))
		}
	}

	topology, err := cre.NewTopology(nodeSets, provider, globalCapConfigs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create topology")
	}

	for i, donMeta := range topology.DonsMetadata.List() {
		if i >= len(nodeNamesBySet) {
			continue
		}
		if err := hydrateDiscoveredEVMAddresses(
			donMeta,
			nodeNamesBySet[i],
			requiredEVMChainsBySet[i],
			state.NodeRuntime,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to hydrate EVM addresses for DON %s", donMeta.Name)
		}
		if err := hydrateDiscoveredOCR2BundleIDs(
			donMeta,
			nodeNamesBySet[i],
			state.NodeRuntime,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to hydrate OCR2 key bundle IDs for DON %s", donMeta.Name)
		}
		if err := hydrateDiscoveredHosts(
			donMeta,
			nodeNamesBySet[i],
			cv,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to hydrate node hosts for DON %s", donMeta.Name)
		}
	}

	return topology, nil
}

func (d *Deployer) buildNodeSetForDON(
	ctx context.Context,
	don *domain.DON,
	supportedEVMChains []uint64,
	registryChainID uint64,
	bootstrapAssigned *bool,
	cv *domain.ChartValues,
	state *domain.StateFile,
) (*cre.NodeSet, []string, error) {
	workers := don.WorkerNodes(cv)
	if len(workers) == 0 {
		return nil, nil, fmt.Errorf("DON %s has no worker nodes", don.Name)
	}

	bootstrap := don.ResolveBootstrap(cv)
	includeBootstrap := false
	if bootstrap != "" && !*bootstrapAssigned {
		if slices.Contains(cv.NodeNamesForDONName(don.Name), bootstrap) {
			includeBootstrap = true
		}
	}
	if includeBootstrap {
		*bootstrapAssigned = true
	}

	var nodeNames []string
	if includeBootstrap {
		nodeNames = append(nodeNames, bootstrap)
	}
	nodeNames = append(nodeNames, workers...)

	nodeSpecs, err := d.buildNodeSpecs(ctx, nodeNames, bootstrap, cv)
	if err != nil {
		return nil, nil, err
	}

	requiredChains := requiredEVMChainIDsForDON(don, supportedEVMChains, registryChainID)
	if err := validateDiscoveredEVMAddresses(don.Name, nodeNames, requiredChains, state.NodeRuntime); err != nil {
		return nil, nil, err
	}

	capConfigs := make(cre.CapabilityConfigs)
	for k, v := range don.CapabilityConfigs {
		capConfigs[k] = toCreCapabilityConfig(v)
	}

	return &cre.NodeSet{
		Input: &ns.Input{
			Name:         don.Name,
			Nodes:        len(workers),
			OverrideMode: "all",
		},
		NodeSpecs:                    nodeSpecs,
		Capabilities:                 append([]string(nil), don.Capabilities...),
		DONTypes:                     append([]string(nil), don.DONTypes...),
		SupportedEVMChains:           append([]uint64(nil), requiredChains...),
		CapabilityConfigs:            capConfigs,
		ExposesRemoteCapabilities:    don.ExposesRemoteCaps,
		RegistryBasedLaunchAllowlist: append([]string(nil), don.RegistryBasedAllowlist...),
		// A workflow DON's don_family is its own name; the gateway node set(s)
		// serving it (see newGatewayNodeSet) are given the same value via
		// desired.GatewayDONFor, so GatewayConnectorsForDonFamily/
		// GatewayServiceConfigsForGateway pair them correctly.
		DonFamily: don.Name,
	}, nodeNames, nil
}

func (d *Deployer) buildBootstrapOnlyNodeSet(
	ctx context.Context,
	don *domain.DON,
	supportedEVMChains []uint64,
	registryChainID uint64,
	bootstrapAssigned *bool,
	cv *domain.ChartValues,
	state *domain.StateFile,
) (*cre.NodeSet, []string, []uint64, error) {
	if !don.IsBootstrapDon() {
		return nil, nil, nil, fmt.Errorf("DON %s is not a bootstrap DON; bootstrap-only builder cannot be used", don.Name)
	}

	bootstrap := don.ResolveBootstrap(cv)
	if bootstrap == "" {
		return nil, nil, nil, fmt.Errorf("DON %s is bootstrap-only but has no resolvable bootstrap node", don.Name)
	}

	nodeNames := []string{bootstrap}
	nodeSpecs, err := d.buildNodeSpecs(ctx, nodeNames, bootstrap, cv)
	if err != nil {
		return nil, nil, nil, err
	}

	requiredChains := requiredEVMChainIDsForDON(don, supportedEVMChains, registryChainID)
	if err := validateDiscoveredEVMAddresses(don.Name, nodeNames, requiredChains, state.NodeRuntime); err != nil {
		return nil, nil, nil, err
	}

	if bootstrapAssigned != nil {
		*bootstrapAssigned = true
	}

	return newBootstrapOnlyNodeSet(don, nodeSpecs, requiredChains), nodeNames, requiredChains, nil
}

func (d *Deployer) buildGatewayNodeSets(
	ctx context.Context,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
	supportedEVMChains []uint64,
	registryChainID uint64,
) ([]*cre.NodeSet, [][]string, error) {
	gatewayNodes := cv.FindGatewayNodes()
	if len(gatewayNodes) == 0 {
		return nil, nil, errors.New("gateway capabilities require gateway nodes in chart values")
	}

	nodeSets := make([]*cre.NodeSet, 0, len(gatewayNodes))
	nodeNamesBySet := make([][]string, 0, len(gatewayNodes))
	requiredChains := requiredEVMChainIDsForGatewayNode(registryChainID)

	for _, gw := range gatewayNodes {
		donName := gatewayDONNameForNode(desired, cv, gw.Name)
		if donName == "" {
			return nil, nil, fmt.Errorf("gateway node %s: no DON in desired state has DON type %q and includes this node", gw.Name, "gateway")
		}

		nodeNames := []string{gw.Name}
		nodeSpecs, err := d.buildNodeSpecs(ctx, nodeNames, "", cv)
		if err != nil {
			return nil, nil, err
		}
		if err := validateDiscoveredEVMAddresses(donName, nodeNames, requiredChains, state.NodeRuntime); err != nil {
			return nil, nil, err
		}

		nodeSets = append(nodeSets, newGatewayNodeSet(
			donName,
			nodeSpecs,
			desired.GatewayDONFor(gw.Name),
			requiredChains,
		))
		nodeNamesBySet = append(nodeNamesBySet, nodeNames)
	}

	return nodeSets, nodeNamesBySet, nil
}

func newGatewayNodeSet(
	donName string,
	nodeSpecs []*cre.NodeSpecWithRole,
	gatewayDonID string,
	supportedEVMChains []uint64,
) *cre.NodeSet {
	for i := range nodeSpecs {
		nodeSpecs[i].Roles = []cre.NodeType{cre.GatewayNode}
	}

	return &cre.NodeSet{
		Input: &ns.Input{
			Name:         donName,
			Nodes:        1,
			OverrideMode: "all",
		},
		NodeSpecs:          nodeSpecs,
		DONTypes:           []string{"gateway"},
		SupportedEVMChains: append([]uint64(nil), supportedEVMChains...),
		GatewayDonID:       gatewayDonID,
		// gatewayDonID here is the workflow DON name this gateway serves
		// (desired.GatewayDONFor), not an on-chain ID. Using it as don_family
		// matches the served workflow DON's own don_family (its DON name, set
		// in buildNodeSetForDON), which is how system-tests' pairing logic
		// (GatewayConnectorsForDonFamily/GatewayServiceConfigsForGateway)
		// wires gateway connector config into the workflow DON's node TOML.
		DonFamily: gatewayDonID,
	}
}

func newBootstrapOnlyNodeSet(
	don *domain.DON,
	nodeSpecs []*cre.NodeSpecWithRole,
	supportedEVMChains []uint64,
) *cre.NodeSet {
	capConfigs := make(cre.CapabilityConfigs)
	for k, v := range don.CapabilityConfigs {
		capConfigs[k] = toCreCapabilityConfig(v)
	}

	return &cre.NodeSet{
		Input: &ns.Input{
			Name:         don.Name,
			Nodes:        0,
			OverrideMode: "all",
		},
		NodeSpecs:                    nodeSpecs,
		Capabilities:                 append([]string(nil), don.Capabilities...),
		DONTypes:                     append([]string(nil), don.DONTypes...),
		SupportedEVMChains:           append([]uint64(nil), supportedEVMChains...),
		CapabilityConfigs:            capConfigs,
		ExposesRemoteCapabilities:    don.ExposesRemoteCaps,
		RegistryBasedLaunchAllowlist: append([]string(nil), don.RegistryBasedAllowlist...),
		// Bootstrap DONs are excluded from don_family gateway pairing (they're
		// neither a gateway nor a workflow DON per topology_don_family.go), so
		// this only needs to be non-empty to satisfy NewDonMetadata.
		DonFamily: don.Name,
	}
}

func (d *Deployer) buildNodeSpecs(ctx context.Context, nodeNames []string, bootstrapName string, cv *domain.ChartValues) ([]*cre.NodeSpecWithRole, error) {
	specs := make([]*cre.NodeSpecWithRole, 0, len(nodeNames))
	for _, nodeName := range nodeNames {
		// 00-secrets.toml supplies imported non-EVM key material (P2P, etc.).
		// EVM public addresses are not stored there; they are discovered from the
		// node API during D1 and hydrated into node metadata before on-chain work.
		secretsToml, err := d.k8s.GetNodeSecretsToml(ctx, nodeName, cv.GetNodeNamespace(nodeName))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read secrets for node %s", nodeName)
		}

		role := cre.WorkerNode
		if nodeName == bootstrapName {
			role = cre.BootstrapNode
		}

		specs = append(specs, &cre.NodeSpecWithRole{
			Input: &clnode.Input{
				Node: &clnode.NodeInput{
					TestSecretsOverrides: secretsToml,
				},
			},
			Roles: []cre.NodeType{role},
		})
	}
	return specs, nil
}

func (d *Deployer) buildGlobalCapabilityConfigs(desired *domain.DesiredState) cre.CapabilityConfigs {
	out := make(cre.CapabilityConfigs)
	for k, v := range domain.LoadCapabilityDefaults("") {
		out[k] = toCreCapabilityConfig(v)
	}
	for k, v := range desired.CapabilityConfigs {
		out[k] = toCreCapabilityConfig(v)
	}
	return out
}

func toCreCapabilityConfig(cfg domain.CapabilityConfig) cre.CapabilityConfig {
	return cre.CapabilityConfig{
		BinaryName: cfg.BinaryName,
		Values:     cfg.Values,
	}
}

func donRequiresGatewayAccess(don *domain.DON) bool {
	for _, cap := range don.Capabilities {
		switch cap {
		case cre.VaultCapability, cre.HTTPActionCapability, cre.HTTPTriggerCapability:
			return true
		}
	}
	return false
}

func requiredEVMChainIDsForDON(don *domain.DON, supportedEVMChains []uint64, registryChainID uint64) []uint64 {
	required := make(map[uint64]struct{})
	for _, cap := range don.Capabilities {
		if chainID, ok := domain.ParseEVMChainIDFromCapability(cap); ok {
			required[chainID] = struct{}{}
		}
	}
	if slices.Contains(don.Capabilities, cre.EVMCapability) {
		for _, chainID := range supportedEVMChains {
			required[chainID] = struct{}{}
		}
	}
	if donRequiresGatewayAccess(don) && registryChainID != 0 {
		required[registryChainID] = struct{}{}
	}

	out := make([]uint64, 0, len(required))
	for chainID := range required {
		out = append(out, chainID)
	}
	slices.Sort(out)
	return out
}

func requiredEVMChainIDsForGatewayNode(registryChainID uint64) []uint64 {
	if registryChainID == 0 {
		return nil
	}
	return []uint64{registryChainID}
}
