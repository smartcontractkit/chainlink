// Package reconciler provides a reconciliation tool for configuring Chainlink nodes
// deployed to Kubernetes via Griddle. It deploys smart contracts, registers
// capabilities on-chain, creates jobs via JD, and injects node TOML config —
// all driven by a declarative desired-state TOML file.
package reconciler

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/ui"
)

// Version identifies the build. Overridden at build time via:
//
//	go build -ldflags "-X github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler.Version=$(git rev-parse --short HEAD)-$(date -u +%Y%m%d%H%M%S)"
//
// Defaults to "dev-unknown" for builds that don't set it, so a stale binary
// is always distinguishable from a freshly built one.
var Version = "dev-unknown"

// RootCmd is the cobra root command for the CRE reconciler CLI.
var RootCmd = &cobra.Command{
	Use:     "reconciler",
	Version: Version,
	Short:   "Reconcile Chainlink node configurations for Griddle-deployed nodesets",
	Long: `A standalone CLI tool that reconciles a declarative desired-state TOML file
against actual state (on-chain + JD + Kubernetes) for Chainlink nodes deployed
to Kubernetes via Griddle.

It deploys smart contracts, registers capabilities on-chain, creates jobs via JD,
and injects node TOML config — reusing the existing changeset library from
chainlink/deployment.

Usage:
  reconciler apply --desired cre/desired.toml --state cre/state.toml
  reconciler status --state cre/state.toml
  reconciler diff --desired cre/desired.toml --state cre/state.toml`,
}

// CLI flags shared across commands.
var (
	flagDesired           string
	flagState             string
	flagKubeconfig        string
	flagEnv               string
	flagAddr              string
	flagChartDir          string
	flagConfirm           bool
	flagRestartWorkerPods bool
	flagWaitAtBreakpoint  bool
	flagDeployerKey       string
)

// CLI flags for the workflow-deploy command.
var (
	flagWFFile      string
	flagWFName      string
	flagWFOwner     string
	flagWFDonFamily string
	flagWFConfig    string
	flagWFTag       string
	flagWFRemoteDir string
	flagWFContainer string
	flagWFNamespace string
	flagWFPods      []string
)

// CLI flags for the workflow-trigger command.
var (
	flagTGGatewayURL   string
	flagTGWorkflowName string
	flagTGOwner        string
	flagTGWorkflowTag  string
	flagTGWorkflowID   string
	flagTGPrivateKey   string
	flagTGInput        string
	flagTGTimeout      time.Duration
	flagTGPollInterval time.Duration
)

func init() {
	RootCmd.AddCommand(applyCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(diffCmd)
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(workflowDeployCmd)
	RootCmd.AddCommand(workflowTriggerCmd)

	RootCmd.PersistentPreRun = trackCommandPreRun
}

var applyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Run the reconcile loop to converge actual state with desired state",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApply(cmd.Context())
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current state and what would change",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus(cmd.Context())
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what would change without executing anything",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDiff(cmd.Context())
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web UI for visual DON builder and status monitoring",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServe(cmd.Context())
	},
}

var workflowDeployCmd = &cobra.Command{
	Use:   "workflow-deploy",
	Short: "Compile a workflow and push it to a private file-based workflow registry on running pods",
	Long: `Compiles a workflow to WASM, brotli-compresses it, computes its workflow ID,
builds a single-entry private file-registry JSON, and copies both files into
the same directory on every target pod. No pod restart is required — the
node's v2 file workflow source re-reads the registry and re-fetches the
binary on its next poll. The node must already have a TOML AdditionalSources
entry pointing file:// at that directory (see toml_patch.go / GenerateNodeTOML).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorkflowDeploy(cmd.Context())
	},
}

var workflowTriggerCmd = &cobra.Command{
	Use:   "workflow-trigger",
	Short: "Trigger a deployed HTTP-triggered workflow via a gateway",
	Long: `Builds a workflows.execute JSON-RPC request, signs it with an ECDSA private key
(matching one of the workflow's AuthorizedKeys), and POSTs it to a gateway's external HTTP
trigger endpoint — the same flow used by system-tests/tests/smoke/cre/http_trigger_action_test.go.
Retries on failure (e.g. the workflow not yet loaded on the node) until --timeout elapses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorkflowTrigger(cmd.Context())
	},
}

func init() {
	// Apply flags
	applyCmd.Flags().StringVarP(&flagDesired, "desired", "d", "cre/desired.toml", "path to desired-state TOML")
	applyCmd.Flags().StringVarP(&flagState, "state", "s", "cre/state.toml", "path to state file")
	applyCmd.Flags().StringVar(&flagKubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to KUBECONFIG env or ~/.kube/config)")
	applyCmd.Flags().StringVarP(&flagEnv, "env", "e", "dev", "environment name (e.g. dev, stage, prod)")
	applyCmd.Flags().BoolVar(&flagConfirm, "confirm", false, "prompt for confirmation before each step (shows full details)")
	applyCmd.Flags().BoolVar(&flagRestartWorkerPods, "restart-worker-pods", false,
		"restart worker/standard node pods (excluding bootstrap and gateway) as the last step, working around capabilities that don't clean up state after job cancellation")
	applyCmd.Flags().BoolVar(&flagWaitAtBreakpoint, "wait-at-breakpoint", true,
		"pause in-process at the TOML breakpoint and continue after the user presses Enter; set false to exit with code 42 for the two-invocation workflow (both paths restore persisted gateway handlers)")
	applyCmd.Flags().StringVar(&flagDeployerKey, "deployer-key", "",
		"hex private key for the on-chain deployer (defaults to the Anvil dev account)")

	// Status flags
	statusCmd.Flags().StringVarP(&flagState, "state", "s", "cre/state.toml", "path to state file")
	statusCmd.Flags().StringVar(&flagKubeconfig, "kubeconfig", "", "path to kubeconfig")
	statusCmd.Flags().StringVarP(&flagEnv, "env", "e", "dev", "environment name")

	// Diff flags
	diffCmd.Flags().StringVarP(&flagDesired, "desired", "d", "cre/desired.toml", "path to desired-state TOML")
	diffCmd.Flags().StringVarP(&flagState, "state", "s", "cre/state.toml", "path to state file")
	diffCmd.Flags().StringVar(&flagKubeconfig, "kubeconfig", "", "path to kubeconfig")
	diffCmd.Flags().StringVarP(&flagEnv, "env", "e", "dev", "environment name")

	// Serve flags
	serveCmd.Flags().StringVarP(&flagDesired, "desired", "d", "cre/desired.toml", "path to desired-state TOML")
	serveCmd.Flags().StringVarP(&flagState, "state", "s", "cre/state.toml", "path to state file")
	serveCmd.Flags().StringVar(&flagChartDir, "chart-dir", ".", "path to repo root containing griddle.yaml")
	serveCmd.Flags().StringVarP(&flagEnv, "env", "e", "dev", "environment name")
	serveCmd.Flags().StringVar(&flagAddr, "addr", "localhost:8089", "address to serve the web UI on")
	serveCmd.Flags().StringVar(&flagKubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to KUBECONFIG env or ~/.kube/config)")

	// workflow-deploy flags
	workflowDeployCmd.Flags().StringVar(&flagWFFile, "workflow-file", "", "path to the workflow source file (.go or .ts) (required)")
	workflowDeployCmd.Flags().StringVar(&flagWFName, "workflow-name", "", "workflow name, at least 10 characters (required)")
	workflowDeployCmd.Flags().StringVar(&flagWFOwner, "owner", "", "hex-encoded workflow owner address, with or without 0x (required)")
	workflowDeployCmd.Flags().StringVar(&flagWFDonFamily, "don-family", "workflow", "DON family this workflow belongs to")
	workflowDeployCmd.Flags().StringVar(&flagWFConfig, "config", "", "optional path to a workflow config file")
	workflowDeployCmd.Flags().StringVar(&flagWFTag, "tag", "v1.0.0", "registry entry tag/version")
	workflowDeployCmd.Flags().StringVar(&flagWFRemoteDir, "remote-dir", "/home/chainlink/workflows", "directory on the pod both the registry file and binary are copied into (must match the node's file:// TOML paths)")
	workflowDeployCmd.Flags().StringVar(&flagWFContainer, "container", "chainlink-node", "container name to exec into inside each pod")
	workflowDeployCmd.Flags().StringVar(&flagWFNamespace, "namespace", "", "default namespace for --pod values given without one")
	workflowDeployCmd.Flags().StringVar(&flagKubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to KUBECONFIG env or ~/.kube/config)")
	workflowDeployCmd.Flags().StringArrayVar(&flagWFPods, "pod", nil, "target pod as namespace/podName (repeatable, required)")

	// workflow-trigger flags
	workflowTriggerCmd.Flags().StringVar(&flagTGGatewayURL, "gateway-url", "", "full external gateway URL to send the trigger request to (required)")
	workflowTriggerCmd.Flags().StringVar(&flagTGWorkflowName, "workflow-name", "", "workflow name (required)")
	workflowTriggerCmd.Flags().StringVar(&flagTGOwner, "owner", "", "hex-encoded workflow owner address, with or without 0x (required)")
	workflowTriggerCmd.Flags().StringVar(&flagTGWorkflowTag, "tag", "", "workflow tag/version (optional)")
	workflowTriggerCmd.Flags().StringVar(&flagTGWorkflowID, "workflow-id", "", "hex-encoded workflow ID (optional)")
	workflowTriggerCmd.Flags().StringVar(&flagTGPrivateKey, "private-key", "", "hex-encoded ECDSA private key used to sign the request, with or without 0x (required)")
	workflowTriggerCmd.Flags().StringVar(&flagTGInput, "input", "{}", "JSON payload for the trigger; prefix with @ to read from a file")
	workflowTriggerCmd.Flags().DurationVar(&flagTGTimeout, "timeout", 2*time.Minute, "how long to keep retrying before giving up")
	workflowTriggerCmd.Flags().DurationVar(&flagTGPollInterval, "poll-interval", 5*time.Second, "delay between retries")
}

// Run is the entry point called from main.go.
func Run() {
	var trackOnce sync.Once

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		trackOnce.Do(func() { trackCommandResult("interrupted") })
		os.Exit(0)
	}()

	err := RootCmd.Execute()
	trackOnce.Do(func() {
		if err != nil {
			trackCommandResult("error")
		} else {
			trackCommandResult("success")
		}
	})
	if err == nil {
		return
	}
	if stderrors.Is(err, ErrBreakpoint) {
		os.Exit(ExitCodeTOMLPatch) // 42, handoff
	}
	// applyCmd sets SilenceErrors so cobra doesn't double-print the breakpoint's
	// noisy handoff message — but that means genuine errors get no message at
	// all unless printed here.
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

// ExitCodeTOMLPatch signals the tool wrote a TOML patch and the user needs
// to apply it via Griddle before re-running. Used as a custom exit code
// so CI can distinguish "needs handoff" from success/error.
const ExitCodeTOMLPatch = 42

func runApply(ctx context.Context) error {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}).
		Level(zerolog.DebugLevel).
		With().Timestamp().Logger()

	log.Info().Str("version", Version).Msg("reconciler starting")

	reconciler, err := NewReconciler(flagDesired, flagState, flagKubeconfig, flagEnv, flagConfirm, flagRestartWorkerPods, flagWaitAtBreakpoint, flagDeployerKey, log)
	if err != nil {
		return err
	}

	return reconciler.Run(ctx)
}

func runStatus(_ context.Context) error {
	st, err := domain.LoadState(flagState)
	if err != nil {
		return err
	}
	if st == nil {
		fmt.Println("No state file — nothing reconciled yet.")
		return nil
	}
	fmt.Printf("Phase: %s\n", st.Phase)
	fmt.Println("Contracts:")
	for _, a := range st.Addresses {
		fmt.Printf("  %-24s %s (chain %d, v%s)\n", a.Type, a.Address, uint64(a.ChainSelector), a.Version)
	}
	fmt.Println("DON IDs:")
	for name, id := range st.DONIDs {
		fmt.Printf("  %-24s %d\n", name, id)
	}
	fmt.Printf("Discovered nodes: %d\n", len(st.NodeRuntime))
	for name, n := range st.NodeRuntime {
		fmt.Printf("  %-28s role=%-9s api=%s peer=%s\n", name, n.NodeType, n.APIURL, truncStr(n.PeerID, 16))
	}
	return nil
}

func runDiff(_ context.Context) error {
	ds, err := domain.LoadDesiredState(flagDesired)
	if err != nil {
		return err
	}
	cv, err := domain.LoadChartValues(ds.Infra.ChartValues, flagEnv)
	if err != nil {
		return err
	}
	st, _ := domain.LoadState(flagState)
	if st == nil {
		st = &domain.StateFile{}
	}

	fmt.Println("=== Contracts ===")
	for _, t := range []string{
		"CapabilitiesRegistry",
		"WorkflowRegistry",
	} {
		if st.HasAddress(t) {
			fmt.Printf("  [have]    %s = %s\n", t, st.GetAddress(t))
		} else {
			fmt.Printf("  [deploy]  %s\n", t)
		}
	}
	fmt.Println("=== DONs ===")
	for _, don := range ds.DONs {
		members := cv.NodeNamesForDONName(don.Name)
		status := "configure"
		if _, ok := st.DONIDs[don.Name]; ok {
			status = "have (id " + strconv.FormatUint(st.DONIDs[don.Name], 10) + ")"
		}
		fmt.Printf("  [%s] %s: %d nodes, caps=%v\n", status, don.Name, len(members), don.Capabilities)
	}
	fmt.Println("=== Workflow Registry ===")
	if st.WorkflowReg != nil {
		fmt.Println("  [have]   configured")
	} else {
		fmt.Println("  [pending] not configured")
	}
	fmt.Println("=== Node config ===")
	if st.TOMLPatchApplied {
		fmt.Println("  [have]   30-cre layer written")
	} else {
		fmt.Println("  [pending] not written")
	}
	return nil
}

func runWorkflowDeploy(ctx context.Context) error {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}).
		Level(zerolog.DebugLevel).
		With().Timestamp().Logger()

	if flagWFFile == "" {
		return errors.New("--workflow-file is required")
	}
	if flagWFName == "" {
		return errors.New("--workflow-name is required")
	}
	if flagWFOwner == "" {
		return errors.New("--owner is required")
	}
	if len(flagWFPods) == 0 {
		return errors.New("at least one --pod is required")
	}

	pods := make([]PodTarget, 0, len(flagWFPods))
	for _, p := range flagWFPods {
		if !strings.Contains(p, "/") {
			if flagWFNamespace == "" {
				return fmt.Errorf("--pod %q has no namespace and --namespace was not set", p)
			}
			pods = append(pods, PodTarget{Namespace: flagWFNamespace, PodName: p})
			continue
		}
		target, err := ParsePodTarget(p)
		if err != nil {
			return err
		}
		pods = append(pods, target)
	}

	namespace := flagWFNamespace
	if namespace == "" && len(pods) > 0 {
		namespace = pods[0].Namespace
	}

	k8s, err := infra.NewK8sClient(flagKubeconfig, namespace, log)
	if err != nil {
		return errors.Wrap(err, "failed to create k8s client")
	}

	result, err := WorkflowDeploy(ctx, k8s, WorkflowDeployInputs{
		WorkflowFilePath: flagWFFile,
		WorkflowName:     flagWFName,
		Owner:            flagWFOwner,
		DonFamily:        flagWFDonFamily,
		ConfigPath:       flagWFConfig,
		Tag:              flagWFTag,
		RemoteDir:        flagWFRemoteDir,
		Container:        flagWFContainer,
		Pods:             pods,
	}, log)
	if result != nil {
		fmt.Printf("Workflow ID: %s\n", result.WorkflowID)
	}
	return err
}

func runWorkflowTrigger(ctx context.Context) error {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}).
		Level(zerolog.DebugLevel).
		With().Timestamp().Logger()

	if flagTGGatewayURL == "" {
		return errors.New("--gateway-url is required")
	}
	if flagTGWorkflowName == "" {
		return errors.New("--workflow-name is required")
	}
	if flagTGOwner == "" {
		return errors.New("--owner is required")
	}
	if flagTGPrivateKey == "" {
		return errors.New("--private-key is required")
	}

	result, err := WorkflowTrigger(ctx, WorkflowTriggerInputs{
		GatewayURL:    flagTGGatewayURL,
		WorkflowName:  flagTGWorkflowName,
		WorkflowOwner: flagTGOwner,
		WorkflowTag:   flagTGWorkflowTag,
		WorkflowID:    flagTGWorkflowID,
		PrivateKeyHex: flagTGPrivateKey,
		Input:         flagTGInput,
		Timeout:       flagTGTimeout,
		PollInterval:  flagTGPollInterval,
	}, log)
	if err != nil {
		return err
	}

	responseJSON, err := json.MarshalIndent(result.Response, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal response")
	}
	fmt.Println(string(responseJSON))
	return nil
}

func runServe(_ context.Context) error {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Try to resolve namespace from desired state
	namespace := ""
	if ds, err := domain.LoadDesiredState(flagDesired); err == nil {
		namespace = ds.Infra.Namespace
	}

	// Try to resolve namespace from griddle.yaml
	if namespace == "" {
		if cv, err := domain.LoadChartValues(flagChartDir, flagEnv); err == nil {
			namespace = cv.Namespace
		}
	}
	if namespace == "" {
		namespace = "default"
	}

	// Ensure files exist (create empty desired state if not present)
	if _, err := os.Stat(flagDesired); os.IsNotExist(err) {
		dir := filepath.Dir(flagDesired)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errors.Wrapf(err, "failed to create dir %s", dir)
			}
		}
		if err := os.WriteFile(flagDesired, []byte{}, 0600); err != nil {
			return errors.Wrapf(err, "failed to create desired state file %s", flagDesired)
		}
		log.Info().Msgf("Created empty desired state file: %s", flagDesired)
	}

	server := ui.NewServer(flagDesired, flagState, flagChartDir, flagEnv, namespace, flagKubeconfig, log)

	log.Info().Msgf("Desired state: %s", flagDesired)
	log.Info().Msgf("State file: %s", flagState)
	log.Info().Msgf("Repo root (chart-dir): %s", flagChartDir)
	log.Info().Msgf("Environment: %s", flagEnv)

	if err := server.ListenAndServe(flagAddr); err != nil {
		return errors.Wrap(err, "web server failed")
	}

	return nil
}
