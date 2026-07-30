package reconciler

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/discovery"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/nodeconfig"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/onchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// Reconciler orchestrates the full reconcile flow.
type Reconciler struct {
	desired           *domain.DesiredState
	state             *domain.StateFile
	statePath         string
	cv                *domain.ChartValues
	k8s               K8sAPI
	nodeDialer        NodeDialer
	deployerKey       string
	log               zerolog.Logger
	confirm           bool
	restartWorkerPods bool
	waitAtBreakpoint  bool
}

// NewReconciler creates a Reconciler from the CLI flags.
func NewReconciler(desiredPath, statePath, chartDir, kubeconfig, env string, confirm, restartWorkerPods, waitAtBreakpoint bool, deployerKey string, log zerolog.Logger) (*Reconciler, error) {
	log.Info().Str("desired", desiredPath).Str("state", statePath).Msg("Loading desired state")

	ds, err := domain.LoadDesiredState(desiredPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load desired state")
	}

	log.Info().Str("chartDir", chartDir).Str("env", env).Msg("Loading chart values from griddle.yaml")
	cv, err := domain.LoadChartValues(chartDir, env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load chart values")
	}
	log.Info().Int("nodes", len(cv.Nodes)).Str("namespace", cv.Namespace).Msg("Chart values loaded")

	// Cross-validate: each DON's name and node membership must exactly match
	// what Griddle registered with JD (chainlink-node.registerNodes.labels.don-name),
	// since that's the label JD job-proposal filters actually match against.
	if err = validateDONNamesAgainstChart(ds, cv); err != nil {
		return nil, errors.Wrap(err, "desired.toml does not match Griddle chart values")
	}

	// Load state
	state, err := domain.LoadState(statePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load state")
	}
	if state == nil {
		state = &domain.StateFile{}
		log.Info().Msg("No existing state file — starting fresh")
	} else {
		log.Info().Str("phase", string(state.Phase)).Msg("Resuming from state")
	}

	// Create K8s client
	var k8s *infra.K8sClient
	if kubeconfig == "" {
		// Try KUBECONFIG env, then ~/.kube/config
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	k8s, err = infra.NewK8sClient(kubeconfig, cv.Namespace, log)
	if err != nil {
		log.Warn().Err(err).Msg("K8s client creation failed — node discovery will be limited")
	} else {
		log.Info().Str("namespace", cv.Namespace).Msg("K8s client connected")
	}

	if deployerKey == "" {
		deployerKey = blockchain.DefaultAnvilPrivateKey
	}

	return &Reconciler{
		desired:           ds,
		state:             state,
		statePath:         statePath,
		cv:                cv,
		k8s:               k8s,
		nodeDialer:        clNodeDialer{},
		deployerKey:       deployerKey,
		log:               log,
		confirm:           confirm,
		restartWorkerPods: restartWorkerPods,
		waitAtBreakpoint:  waitAtBreakpoint,
	}, nil
}

// onChainComplete reports whether the on-chain deploy+configure phase actually
// finished (contracts deployed, DON IDs resolved, workflow registry configured).
// Phase alone is not reliable because older runs could mark phase=done before
// on-chain was implemented.
func (r *Reconciler) onChainComplete() bool {
	return r.state.WorkflowReg != nil &&
		r.state.HasAddress(keystone_changeset.CapabilitiesRegistry.String()) &&
		len(r.state.DONIDs) > 0
}

// Run executes the full reconcile flow: discovery → provisioning → config-write → job-sync.
// It persists state after each phase and returns ErrBreakpoint if the TOML breakpoint
// is reached and --wait-at-breakpoint=false.
func (r *Reconciler) Run(ctx context.Context) error {
	if err := r.requireJDAccessToken(); err != nil {
		return err
	}

	if err := r.runDiscoveryPhase(ctx); err != nil {
		return errors.Wrap(err, "discovery failed")
	}

	// Runs unconditionally, every invocation — including once on-chain work is
	// already complete and Apply would otherwise be skipped below — so a node
	// label that broke after the initial apply is still caught.
	if err := onchain.ValidateNodeLabels(r.desired, r.cv, r.state); err != nil {
		return errors.Wrap(err, "node label preflight failed")
	}

	d := onchain.NewDeployer(r.k8s, r.deployerKey, r.log, r.confirmStep)

	var onChainChanged bool
	if !r.onChainComplete() {
		if err := r.requireDiscoveryComplete(); err != nil {
			return err
		}
		changed, err := d.Apply(ctx, r.desired, r.cv, r.state, r.persistState)
		if err != nil {
			return errors.Wrap(err, "on-chain provisioning failed")
		}
		onChainChanged = changed
		r.persistState()
	}

	// injectTOML runs every time (no one-shot gate): it diffs freshly generated
	// TOML against the chart's current layer per node and is a no-op — no patch,
	// no breakpoint — when nothing differs.
	if err := r.injectTOML(ctx); err != nil {
		return err // ErrBreakpoint in the exit-42 flow only
	}
	r.persistState()

	// SyncJobs first restores the persisted gateway service configs (Feature-added handlers) onto the freshly
	// built topology, then cancels/deletes jobs, creates gateway jobs, and runs PostEnvStartup. A change earlier
	// in the on-chain chain forces Jobs to run too, regardless of its own hash.
	if err := d.SyncJobs(ctx, r.desired, r.cv, r.state, onChainChanged); err != nil {
		return errors.Wrap(err, "job sync failed")
	}
	r.state.Phase = domain.PhaseDone
	r.persistState()

	if r.restartWorkerPods {
		r.log.Info().Msg("Restarting worker node pods (--restart-worker-pods)")
		if err := r.restartWorkerNodePods(ctx); err != nil {
			return errors.Wrap(err, "failed to restart worker node pods")
		}
	}
	return nil
}

// restartWorkerNodePods restarts every worker/standard node's pod(s),
// excluding bootstrap and gateway nodes. Ad-hoc workaround for capabilities
// that don't clean up state correctly after job cancellation (deleteAllJobsForDons).
func (r *Reconciler) restartWorkerNodePods(ctx context.Context) error {
	for _, n := range r.cv.Nodes {
		if n.NodeType == domain.RoleBootstrap || n.NodeType == domain.RoleGateway {
			continue
		}
		if err := r.k8s.RestartNodePods(ctx, n.Name, n.Namespace); err != nil {
			return errors.Wrapf(err, "failed to restart pod for node %s", n.Name)
		}
		r.log.Info().Str("node", n.Name).Msg("Restarted worker node pod")
	}
	return nil
}

// validateDONNamesAgainstChart requires each DON's name to match a real
// Griddle node-set's don-name label, and that no two DONs claim the same
// node-set. Node membership itself is always chart-derived
// (cv.NodeNamesForDONName) — desired.toml carries no node list to reconcile.
func validateDONNamesAgainstChart(ds *domain.DesiredState, cv *domain.ChartValues) error {
	seen := make(map[string]string) // chart DON name -> desired.toml DON name that claimed it
	for i, don := range ds.DONs {
		members := cv.NodeNamesForDONName(don.Name)
		if len(members) == 0 {
			return fmt.Errorf("dons[%d] (%s): no Griddle node-set has don-name %q — desired.toml DON names must match a chart node-set", i, don.Name, don.Name)
		}

		if other, ok := seen[don.Name]; ok {
			return fmt.Errorf("dons[%d] (%s): node-set %q is already claimed by dons entry %q", i, don.Name, don.Name, other)
		}
		seen[don.Name] = don.Name
	}
	return nil
}

// confirmStep prints the step details and asks for user confirmation if
// --confirm is enabled. Returns nil if the step should proceed, or
// ErrStepDeclined if the user declines. Matches onchain.ConfirmFunc's signature
// so it can be passed directly to onchain.NewDeployer.
func (r *Reconciler) confirmStep(title, details string) error {
	if !r.confirm {
		return nil
	}
	fmt.Println()
	fmt.Printf("┌─ %s\n", title)
	// Print details with a border
	for line := range strings.SplitSeq(details, "\n") {
		fmt.Printf("│ %s\n", line)
	}
	fmt.Print("└─ Proceed? [y/N]: ")
	var resp string
	_, _ = fmt.Scanln(&resp)
	if strings.ToLower(strings.TrimSpace(resp)) == "y" {
		return nil
	}
	r.log.Info().Str("step", title).Msg("Step skipped by user")
	return ErrStepDeclined
}

// runDiscoveryPhase collects node runtime data, advances phase when appropriate,
// and persists state. On failure the phase is left unchanged.
func (r *Reconciler) runDiscoveryPhase(ctx context.Context) error {
	if r.state.Phase == domain.PhaseDiscovery && r.discoveryComplete() {
		r.log.Info().
			Str("phase", string(r.state.Phase)).
			Int("nodes", len(r.state.NodeRuntime)).
			Msg("Discovery already complete — continuing to on-chain")
		return nil
	}

	// Past discovery: keep cached runtime, do not re-prompt.
	if r.state.Phase != "" && r.state.Phase != domain.PhaseDiscovery {
		return nil
	}

	if err := r.discover(ctx); err != nil {
		return err
	}
	if err := r.refreshNodeConfigFiles(); err != nil {
		return errors.Wrap(err, "failed to build node config file mapping")
	}
	advancePhaseAfterDiscovery(r.state)
	r.persistState()
	r.log.Info().Str("phase", string(r.state.Phase)).Int("nodes", len(r.state.NodeRuntime)).Msg("Discovery phase complete")
	return nil
}

func (r *Reconciler) discoveryComplete() bool {
	if r.state.NodeRuntime == nil || len(r.cv.Nodes) == 0 {
		return false
	}
	for _, node := range r.cv.Nodes {
		info, ok := r.state.NodeRuntime[node.Name]
		if !ok || info.APIURL == "" {
			return false
		}
	}
	return true
}

func advancePhaseAfterDiscovery(state *domain.StateFile) {
	if state.Phase == "" || state.Phase == domain.PhaseDiscovery {
		state.Phase = domain.PhaseDiscovery
	}
}

func (r *Reconciler) requireDiscoveryComplete() error {
	if len(r.state.NodeRuntime) == 0 {
		return errors.New("cannot start on-chain: discovery data missing (node_runtime is empty)")
	}
	return nil
}

// requireJDAccessToken enforces that if JD is configured (jd.grpc is set), the JD access
// token env var is also set. Checked before any work in Run starts — including K8s
// discovery — so a missing token fails immediately instead of after wasted work or with a
// confusing downstream "JD client is required" error once the on-chain/job phases hit JD.
func (r *Reconciler) requireJDAccessToken() error {
	if r.desired.JD.GRPC != "" && infra.JDAccessToken() == "" {
		return fmt.Errorf("JD configured (jd.grpc set) but %s is not set — export the JD access token", infra.JDAccessTokenEnv)
	}
	return nil
}

// k8sAdapter wraps K8sAPI to satisfy discovery.K8sClient (allows access to infra.NodeAPIInfo).
type k8sAdapter struct {
	k8s K8sAPI
}

func (a *k8sAdapter) GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*infra.NodeAPIInfo, error) {
	return a.k8s.GetNodeAPIInfo(ctx, nodeName, namespace)
}

// dialerAdapter wraps NodeDialer + NodeClient to satisfy discovery interfaces.
type dialerAdapter struct {
	dialer NodeDialer
}

func (a *dialerAdapter) Dial(apiURL, email, password string) (discovery.Client, error) {
	nc, err := a.dialer.Dial(apiURL, email, password)
	if err != nil {
		return nil, err
	}
	return &clientAdapter{client: nc}, nil
}

type clientAdapter struct {
	client NodeClient
}

func (a *clientAdapter) ReadCSAKey() (string, error) {
	return a.client.ReadCSAKey()
}

func (a *clientAdapter) ReadPeerID() (string, error) {
	return a.client.ReadPeerID()
}

func (a *clientAdapter) ReadEVMAddresses() (map[string]string, error) {
	return a.client.ReadEVMAddresses()
}

func (a *clientAdapter) ReadOCR2BundleIDs() (map[string]string, error) {
	return a.client.ReadOCR2BundleIDs()
}

func (a *clientAdapter) ReadAptosKeys() (string, error) {
	return a.client.ReadAptosKeys()
}

func (a *clientAdapter) ReadSolanaKeys() (string, error) {
	return a.client.ReadSolanaKeys()
}

// nonEVMFamiliesByNode maps worker node name -> non-EVM chain families ("aptos"/"solana")
// that node's DON declares a capability for, so discovery only attempts a Read*Keys call
// where a key is actually expected to exist.
func nonEVMFamiliesByNode(desired *domain.DesiredState, cv *domain.ChartValues) map[string][]string {
	out := make(map[string][]string)
	for i := range desired.DONs {
		don := &desired.DONs[i]
		families := don.NonEVMFamilies()
		if len(families) == 0 {
			continue
		}
		for _, nodeName := range don.WorkerNodes(cv) {
			out[nodeName] = append(out[nodeName], families...)
		}
	}
	return out
}

// D1: Discover node runtime info from K8s + node API.
func (r *Reconciler) discover(ctx context.Context) error {
	r.log.Info().Msg("=== D1: Discovery ===")

	if r.k8s == nil {
		return errors.New("K8s client not available — cannot discover nodes")
	}

	// Build a summary of what will be discovered
	var summary strings.Builder
	fmt.Fprintf(&summary, "Namespace: %s\n", r.cv.Namespace)
	fmt.Fprintf(&summary, "Nodes to discover: %d\n", len(r.cv.Nodes))
	for _, node := range r.cv.Nodes {
		fmt.Fprintf(&summary, "  • %s (%s)\n", node.Name, node.NodeType)
	}
	summary.WriteString("\nFor each node: read API URL from HTTPRoute, get creds from K8s secret,\nconnect to node API, fetch CSA key + P2P PeerID.\n")

	if err := r.confirmStep("D1: Discover node runtime info from K8s + node API", summary.String()); err != nil {
		return err
	}

	if r.state.NodeRuntime == nil {
		r.state.NodeRuntime = make(map[string]domain.NodeRuntimeInfo)
	}

	// Use discovery.Run to parallelize K8s and node API discovery
	results, err := discovery.Run(ctx, r.log, r.cv.Nodes, r.cv, &k8sAdapter{r.k8s}, &dialerAdapter{r.nodeDialer}, nonEVMFamiliesByNode(r.desired, r.cv))
	if err != nil {
		return errors.Wrap(err, "discovery failed")
	}

	// Merge results into state (successful discoveries only)
	for nodeName, info := range results {
		r.state.NodeRuntime[nodeName] = info
		r.log.Info().
			Str("node", nodeName).
			Str("apiURL", info.APIURL).
			Str("csaKey", truncStr(info.CSAKey, 20)).
			Str("peerID", truncStr(info.PeerID, 20)).
			Msg("Discovered")
	}

	return nil
}

// T1: Generate TOML patch for all nodes and write into chart values.
func (r *Reconciler) injectTOML(_ context.Context) error {
	r.log.Info().Msg("=== T1: TOML injection ===")

	if err := r.refreshNodeConfigFiles(); err != nil {
		return errors.Wrap(err, "failed to refresh node config file mapping")
	}

	globalBootstrap, err := r.resolveGlobalBootstrapForTOML()
	if err != nil {
		return errors.Wrap(err, "failed to resolve global bootstrap for TOML")
	}
	r.log.Info().
		Str("bootstrap", globalBootstrap.nodeName).
		Str("peerID", truncStr(globalBootstrap.peerID, 20)).
		Str("host", globalBootstrap.host).
		Msg("Resolved global bootstrap for TOML")

	// Generate TOML for each node
	var nodePatches []nodeTOMLPatch
	capRegAddr := r.state.GetAddress("CapabilitiesRegistry")
	wfRegAddr := r.state.GetAddress("WorkflowRegistry")
	registryChain, _ := r.desired.RegistryChain()
	registryChainID := registryChain.ChainID

	for _, node := range r.cv.Nodes {
		_, ok := r.state.NodeRuntime[node.Name]
		if !ok {
			r.log.Warn().Str("node", node.Name).Msg("No runtime info, skipping TOML")
			continue
		}

		// Find which DON this node belongs to via its chart don-name label.
		var donName string
		var don *domain.DON
		if chartNode := r.cv.GetNode(node.Name); chartNode != nil {
			donName = chartNode.DONName
			don = r.desired.DONByName(donName)
		}

		inputs := nodeconfig.Inputs{
			CapRegAddress:      capRegAddr,
			WorkflowRegAddress: wfRegAddr,
			RegistryChainID:    registryChainID,
		}

		// Global bootstrap peer/host for P2P bootstrapping
		switch node.NodeType {
		case domain.RoleBootstrap:
			inputs.IsBootstrapNode = true
			inputs.BootstrapPeerID = globalBootstrap.peerID
			r.log.Info().Str("node", node.Name).Msg("Generating bootstrap node TOML")
		case domain.RoleGateway:
			inputs.IsGatewayNode = true
			r.log.Info().Str("node", node.Name).Msg("Generating gateway node TOML")
		default:
			if globalBootstrap.peerID == "" || globalBootstrap.host == "" {
				return fmt.Errorf(
					"node %s: missing global bootstrap info for worker TOML (bootstrap=%q peerID=%q host=%q)",
					node.Name, globalBootstrap.nodeName, globalBootstrap.peerID, globalBootstrap.host,
				)
			}
			inputs.BootstrapPeerID = globalBootstrap.peerID
			inputs.BootstrapHost = globalBootstrap.host

			if don != nil && donName != "" && r.desired.NeedsGateway() &&
				(don.IsWorkflowDon() || don.NeedsGatewayAccess()) {
				gateways := r.buildWorkerGateways()
				if len(gateways) > 0 {
					evmAddr := ""
					if addrs := r.state.NodeRuntime[node.Name].EVMAddress; addrs != nil {
						evmAddr = addrs[strconv.FormatUint(registryChainID, 10)]
					}
					inputs.GatewayConnector = &nodeconfig.GatewayConnectorConfig{
						DonID:             donName,
						ChainIDForNodeKey: strconv.FormatUint(registryChainID, 10),
						NodeAddress:       evmAddr,
						Gateways:          gateways,
					}
					r.log.Info().
						Str("node", node.Name).
						Int("gateways", len(gateways)).
						Str("nodeAddress", truncStr(evmAddr, 20)).
						Msg("Generating worker TOML with GatewayConnector")
				}
			} else if don != nil && donName != "" {
				r.log.Info().Str("node", node.Name).Str("don", donName).Msg("Generating worker TOML")
			}
		}

		if don != nil {
			inputs.Allowlist = don.RegistryBasedAllowlist
		}

		toml, genErr := nodeconfig.Generate(inputs)
		if genErr != nil {
			return errors.Wrapf(genErr, "node %s: failed to generate config", node.Name)
		}
		nodePatches = append(nodePatches, nodeTOMLPatch{
			Namespace: node.Namespace,
			Name:      node.Name,
			TOML:      toml,
		})
		r.log.Info().
			Str("node", node.Name).
			Str("namespace", node.Namespace).
			Int("tomlLen", len(toml)).
			Msg("Generated TOML")
	}

	patchesByFile, err := groupPatchesByConfigFile(r.state.NodeConfigFiles, nodePatches)
	if err != nil {
		return errors.Wrap(err, "failed to resolve TOML patch targets")
	}

	// Diff each node's freshly generated TOML against its current chart-layer value;
	// drop anything unchanged so we only patch/confirm/breakpoint for real changes.
	patchesByFile, nodePatches, err = filterUnchangedTOMLPatches(patchesByFile, nodePatches, r.state.NodeConfigFiles)
	if err != nil {
		return errors.Wrap(err, "failed to diff generated TOML against chart values")
	}
	if len(patchesByFile) == 0 {
		r.log.Info().Msg("No TOML changes detected — nothing to patch")
		r.state.Phase = domain.PhaseConfigWritten
		return nil
	}

	// Build a preview of what will be written
	var preview strings.Builder
	fmt.Fprintf(&preview, "Target files: %d\n", len(patchesByFile))
	preview.WriteString("Layer name: 30-cre\n")
	fmt.Fprintf(&preview, "Nodes to patch: %d\n", len(nodePatches))
	preview.WriteString("\n--- TOML content per node ---\n\n")
	for _, patch := range nodePatches {
		role := "worker"
		if n := r.cv.GetNodeInNamespace(patch.Namespace, patch.Name); n != nil {
			role = string(n.NodeType)
		}
		file := r.state.NodeConfigFiles[domain.NodeIdentity(patch.Namespace, patch.Name)]
		fmt.Fprintf(&preview, "▼ %s/%s (%s) -> %s\n%s\n\n",
			patch.Namespace, patch.Name, role, file, strings.TrimSpace(patch.TOML))
	}

	for filePath := range patchesByFile {
		currentContent, _ := os.ReadFile(filePath)
		if len(currentContent) == 0 {
			continue
		}
		fmt.Fprintf(&preview, "--- Current file %s (first 500 chars) ---\n", filePath)
		head := string(currentContent)
		if len(head) > 500 {
			head = head[:500] + "..."
		}
		preview.WriteString(head)
		preview.WriteString("\n")
	}

	if err := r.confirmStep("T1: Write TOML patch into chart values", preview.String()); err != nil {
		r.log.Info().Msg("TOML injection skipped by user")
		return nil
	}

	for filePath, patches := range patchesByFile {
		if err := infra.PatchChartValues(filePath, patches); err != nil {
			return errors.Wrapf(err, "failed to patch chart values file %s", filePath)
		}
		r.log.Info().Str("yamlPath", filePath).Int("nodes", len(patches)).Msg("TOML patch written")
	}

	r.state.Phase = domain.PhaseConfigWritten

	r.printBreakpointInstructions()
	if r.waitAtBreakpoint {
		fmt.Print("  4. Press Enter here to continue once the nodes have re-rolled (Ctrl-C to abort)... ")
		var resp string
		_, _ = fmt.Scanln(&resp)
		return nil // continue in-process to syncJobs (which restores gateway handlers)
	}
	fmt.Println("  4. Re-run: reconciler apply")
	return ErrBreakpoint // exit-42 handoff, raised in cmd.go
}

type nodeTOMLPatch struct {
	Namespace string
	Name      string
	TOML      string
}

func (r *Reconciler) refreshNodeConfigFiles() error {
	mapping, err := r.cv.BuildNodeConfigFileMap()
	if err != nil {
		return err
	}
	r.state.NodeConfigFiles = mapping
	return nil
}

func groupPatchesByConfigFile(mapping map[string]string, patches []nodeTOMLPatch) (map[string]map[string]string, error) {
	if len(patches) == 0 {
		return nil, errors.New("no TOML patches to apply")
	}
	if len(mapping) == 0 {
		return nil, errors.New("node config file mapping is empty")
	}

	out := make(map[string]map[string]string)
	for _, patch := range patches {
		key := domain.NodeIdentity(patch.Namespace, patch.Name)
		filePath, ok := mapping[key]
		if !ok {
			return nil, fmt.Errorf("node %s: no config file mapping in state", key)
		}
		if _, err := os.Stat(filePath); err != nil {
			return nil, fmt.Errorf("node %s: config file %s does not exist: %w", key, filePath, err)
		}
		if out[filePath] == nil {
			out[filePath] = make(map[string]string)
		}
		if existing, dup := out[filePath][patch.Name]; dup && existing != patch.TOML {
			return nil, fmt.Errorf("node %s: duplicate patch for node name %q in file %s", key, patch.Name, filePath)
		}
		out[filePath][patch.Name] = patch.TOML
	}
	return out, nil
}

// filterUnchangedTOMLPatches drops any node whose freshly generated TOML
// matches its current "30-cre" chart-layer value, so injectTOML only
// patches/confirms/breaks for nodes that actually changed.
func filterUnchangedTOMLPatches(
	patchesByFile map[string]map[string]string,
	nodePatches []nodeTOMLPatch,
	nodeConfigFiles map[string]string,
) (map[string]map[string]string, []nodeTOMLPatch, error) {
	filteredByFile := make(map[string]map[string]string, len(patchesByFile))
	for filePath, byNode := range patchesByFile {
		for nodeName, tomlContent := range byNode {
			current, exists, err := infra.ReadLayerValue(filePath, nodeName, "30-cre")
			if err != nil {
				return nil, nil, errors.Wrapf(err, "node %s: failed to read current chart layer", nodeName)
			}
			if exists && current == tomlContent {
				continue
			}
			if filteredByFile[filePath] == nil {
				filteredByFile[filePath] = make(map[string]string)
			}
			filteredByFile[filePath][nodeName] = tomlContent
		}
	}

	var filteredPatches []nodeTOMLPatch
	for _, patch := range nodePatches {
		if byNode, ok := filteredByFile[nodeConfigFiles[domain.NodeIdentity(patch.Namespace, patch.Name)]]; ok {
			if _, ok := byNode[patch.Name]; ok {
				filteredPatches = append(filteredPatches, patch)
			}
		}
	}

	return filteredByFile, filteredPatches, nil
}

func (r *Reconciler) buildWorkerGateways() []nodeconfig.ConnectorGateway {
	if len(r.state.GatewayConnectors) > 0 {
		gateways := make([]nodeconfig.ConnectorGateway, 0, len(r.state.GatewayConnectors))
		for _, gc := range r.state.GatewayConnectors {
			gateways = append(gateways, nodeconfig.ConnectorGateway{
				ID:    gc.AuthGatewayID,
				URL:   gc.WebSocketURL,
				DonID: gc.GatewayDonID,
			})
		}
		return gateways
	}

	gwNodes := r.cv.FindGatewayNodes()
	gateways := make([]nodeconfig.ConnectorGateway, 0, len(gwNodes))
	for idx, gwNode := range gwNodes {
		gwDON := r.desired.GatewayDONFor(gwNode.Name)
		gwNamespace := r.cv.GetNodeNamespace(gwNode.Name)
		gwURL := fmt.Sprintf("ws://%s.%s.svc.cluster.local:5003", gwNode.Name, gwNamespace)
		gateways = append(gateways, nodeconfig.ConnectorGateway{
			ID:    fmt.Sprintf("gateway-node-%d", idx),
			URL:   gwURL,
			DonID: gwDON,
		})
	}
	return gateways
}

type globalBootstrapInfo struct {
	nodeName string
	peerID   string
	host     string
}

// resolveGlobalBootstrapForTOML finds the single bootstrap node used for P2P
// DefaultBootstrappers across all worker nodes, whether the bootstrap lives in a
// bootstrap-only DON or inside a workflow DON.
func (r *Reconciler) resolveGlobalBootstrapForTOML() (globalBootstrapInfo, error) {
	bootstrapName, err := r.resolveGlobalBootstrapNodeName()
	if err != nil {
		return globalBootstrapInfo{}, err
	}

	runtime, ok := r.state.NodeRuntime[bootstrapName]
	if !ok {
		return globalBootstrapInfo{}, fmt.Errorf(
			"bootstrap node %s: no discovered runtime info — re-run discovery (D1)",
			bootstrapName,
		)
	}
	if runtime.PeerID == "" {
		return globalBootstrapInfo{}, fmt.Errorf(
			"bootstrap node %s: missing peer_id in runtime state — re-run discovery (D1)",
			bootstrapName,
		)
	}

	return globalBootstrapInfo{
		nodeName: bootstrapName,
		peerID:   runtime.PeerID,
		host:     r.cv.NodeInternalHost(bootstrapName),
	}, nil
}

func (r *Reconciler) resolveGlobalBootstrapNodeName() (string, error) {
	var bootstrapName string
	setBootstrap := func(name string) error {
		if name == "" {
			return nil
		}
		if bootstrapName != "" && bootstrapName != name {
			return fmt.Errorf("multiple bootstrap nodes found: %q and %q", bootstrapName, name)
		}
		bootstrapName = name
		return nil
	}

	// 1. Bootstrap-only DON (preferred for separate-bootstrap topology).
	for _, don := range r.desired.DONs {
		if !don.IsBootstrapOnly(r.cv) {
			continue
		}
		boot := don.ResolveBootstrap(r.cv)
		if boot == "" {
			if members := r.cv.NodeNamesForDONName(don.Name); len(members) > 0 {
				boot = members[0]
			}
		}
		if err := setBootstrap(boot); err != nil {
			return "", err
		}
	}
	if bootstrapName != "" {
		return bootstrapName, nil
	}

	// 2. Explicit bootstrap_node on any DON.
	for _, don := range r.desired.DONs {
		if err := setBootstrap(don.BootstrapNode); err != nil {
			return "", err
		}
	}
	if bootstrapName != "" {
		return bootstrapName, nil
	}

	// 3. Resolved bootstrap from DON nodes / chart node types.
	for _, don := range r.desired.DONs {
		if err := setBootstrap(don.ResolveBootstrap(r.cv)); err != nil {
			return "", err
		}
	}
	if bootstrapName != "" {
		return bootstrapName, nil
	}

	// 4. Chart nodes with bootstrap role.
	var chartBootstrap []string
	for _, node := range r.cv.Nodes {
		if node.NodeType == domain.RoleBootstrap {
			chartBootstrap = append(chartBootstrap, node.Name)
		}
	}
	if len(chartBootstrap) > 1 {
		return "", fmt.Errorf("multiple bootstrap nodes found in chart values: %v", chartBootstrap)
	}
	if len(chartBootstrap) == 1 {
		return chartBootstrap[0], nil
	}

	return "", errors.New("no bootstrap node found in desired state or chart values")
}

// --- Helpers ---

func (r *Reconciler) persistState() {
	if err := r.state.Store(r.statePath); err != nil {
		r.log.Error().Err(err).Msg("Failed to persist state")
	} else {
		r.log.Debug().Str("path", r.statePath).Str("phase", string(r.state.Phase)).Msg("State persisted")
	}
}

func (r *Reconciler) printBreakpointInstructions() {
	fmt.Println()
	fmt.Println("========================================================")
	fmt.Println("  TOML patch written! Manual steps required:")
	fmt.Println("========================================================")
	fmt.Println()
	fmt.Printf("  1. Commit the patched chart values YAML")
	fmt.Printf("  2. Run: .deploy %s\n", r.desired.JD.Environment)
	fmt.Println("  3. Wait for the nodes to re-roll")
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
