// Package reconciler provides a reconciliation tool for configuring Chainlink nodes
// deployed to Kubernetes via Griddle. It deploys smart contracts, registers
// capabilities on-chain, creates jobs via JD, and injects node TOML config —
// all driven by a declarative desired-state TOML file.
package reconciler

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
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

func init() {
	RootCmd.AddCommand(applyCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(diffCmd)
	RootCmd.AddCommand(serveCmd)
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
