package devmachine

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
)

// UpOptions captures all user supplied inputs for `devmachine up`.
type UpOptions struct {
	ConfigPath string

	Name       string
	Namespace  string
	Team       string
	CostCenter string

	InstanceType string
	VolumeSize   int
	AZ           string
	InstanceID   string
	Region       string
	AWSProfile   string

	ChainlinkRepo      string
	JobDistributorRepo string
	SSHKey             string

	ChainlinkImage string
	JDImage        string
	JDTag          string

	EnvWorkdir string
	EnvConfig  string

	DryRun  bool
	Verbose bool
	Yes     bool
	From    string
	Only    string
}

// Orchestrator drives the provisioning phases.
type Orchestrator struct {
	cfg      *Config
	opts     *UpOptions
	run      *Runner
	rep      Reporter
	prompter *Prompter
	cacheDir string
}

// NewOrchestrator wires up an Orchestrator.
func NewOrchestrator(cfg *Config, opts *UpOptions, run *Runner, rep Reporter, prompter *Prompter) *Orchestrator {
	return &Orchestrator{cfg: cfg, opts: opts, run: run, rep: rep, prompter: prompter}
}

// save persists the config unless running in dry-run mode.
func (o *Orchestrator) save() error {
	if o.run.DryRun {
		return nil
	}
	return o.cfg.Save(o.opts.ConfigPath)
}

// Up runs the provisioning phases in order, honouring --from/--only.
func (o *Orchestrator) Up(ctx context.Context) error {
	only, err := o.parseSelector(o.opts.Only)
	if err != nil {
		return err
	}
	from, err := o.parseSelector(o.opts.From)
	if err != nil {
		return err
	}
	return o.runPhases(ctx, from, only)
}

type phaseStep struct {
	id phaseID
	fn func(context.Context) error
}

func (o *Orchestrator) phaseSteps() []phaseStep {
	return []phaseStep{
		{phasePreflight, o.Preflight},
		{phaseMachine, o.Machine},
		{phaseBootstrap, o.Bootstrap},
		{phaseArtifacts, o.Artifacts},
		{phaseTransfer, o.Transfer},
		{phaseLoad, o.Load},
		{phaseStartAgent, o.StartAgent},
		{phaseHandoff, o.Handoff},
	}
}

func (o *Orchestrator) runPhases(ctx context.Context, from, only phaseID) error {
	for _, s := range o.phaseSteps() {
		if !shouldRun(s.id, only, from) {
			continue
		}
		if err := s.fn(ctx); err != nil {
			return errors.Wrapf(err, "phase %s (%s) failed", s.id, phaseTitle(s.id))
		}
	}
	return nil
}

// runPhaseRange executes phases from..to inclusive. Used by the TUI mode which
// handles preflight and handoff outside the interactive screen.
func (o *Orchestrator) runPhaseRange(ctx context.Context, from, to phaseID) error {
	for _, s := range o.phaseSteps() {
		if phaseIndex(s.id) < phaseIndex(from) || phaseIndex(s.id) > phaseIndex(to) {
			continue
		}
		if err := s.fn(ctx); err != nil {
			return errors.Wrapf(err, "phase %s (%s) failed", s.id, phaseTitle(s.id))
		}
	}
	return nil
}

func (o *Orchestrator) parseSelector(s string) (phaseID, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}
	return parsePhase(s)
}

// DownOptions captures inputs for `devmachine down`.
type DownOptions struct {
	ConfigPath  string
	PurgeCache  bool
	KeepMachine bool
	DryRun      bool
	Verbose     bool
}

// Down tears down the machine (if the tool created it) and the managed ssh block.
func Down(ctx context.Context, cfg *Config, opts *DownOptions, run *Runner, rep Reporter) error {
	deleted := false
	if cfg.Machine.CreatedByTool && !opts.KeepMachine {
		if cfg.Machine.Name == "" {
			return errors.New("no machine name recorded; nothing to delete")
		}
		rep.Info("deleting DevMachine %q", cfg.Machine.Name)
		if _, err := run.Run(ctx, "kubectl", "delete", "devmachine", cfg.Machine.Name, "-n", cfg.Cluster.Namespace); err != nil {
			return errors.Wrap(err, "failed to delete DevMachine")
		}
		deleted = true
	} else if !cfg.Machine.CreatedByTool {
		rep.Info("machine was not created by this tool; skipping deletion")
	}

	if !run.DryRun && cfg.SSH.Alias != "" {
		if err := removeSSHConfig(cfg.SSH.Alias); err != nil {
			return err
		}
		rep.Info("removed managed ssh config block for %q", cfg.SSH.Alias)
	}

	// Clear the recorded machine identity so a subsequent `up` creates a fresh
	// machine instead of adopting the one we just deleted.
	if deleted && !opts.KeepMachine {
		cfg.Machine.InstanceID = ""
		cfg.Machine.PrivateIP = ""
		cfg.Machine.CreatedByTool = false
		if !run.DryRun {
			if err := cfg.Save(opts.ConfigPath); err != nil {
				return errors.Wrap(err, "failed to update config after deletion")
			}
		}
	}

	if opts.PurgeCache {
		cache, err := CacheDir()
		if err != nil {
			return err
		}
		if !run.DryRun {
			if err := os.RemoveAll(cache); err != nil {
				return errors.Wrap(err, "failed to purge cache")
			}
		}
		rep.Info("purged cache at %s", cache)
	}
	return nil
}

// StatusOptions captures inputs for `devmachine status`.
type StatusOptions struct {
	ConfigPath string
}

// Status prints the recorded configuration and performs lightweight live checks.
func Status(ctx context.Context, cfg *Config, opts *StatusOptions, run *Runner, rep Reporter) error {
	fmt.Fprintf(run.Out, "Config file:     %s\n", opts.ConfigPath)
	fmt.Fprintf(run.Out, "Machine:         %s\n", cfg.Machine.Name)
	fmt.Fprintf(run.Out, "Instance ID:     %s\n", orDash(cfg.Machine.InstanceID))
	fmt.Fprintf(run.Out, "Private IP:      %s\n", orDash(cfg.Machine.PrivateIP))
	fmt.Fprintf(run.Out, "Namespace:       %s\n", cfg.Cluster.Namespace)
	fmt.Fprintf(run.Out, "AWS profile:     %s\n", orDash(cfg.AWS.Profile))
	fmt.Fprintf(run.Out, "Region:          %s\n", cfg.AWS.Region)
	fmt.Fprintf(run.Out, "SSH alias:       %s\n", cfg.SSH.Alias)
	fmt.Fprintf(run.Out, "Created by tool: %t\n", cfg.Machine.CreatedByTool)

	if cfg.Machine.InstanceID == "" {
		rep.Warn("no instance id recorded yet")
		return nil
	}

	args := []string{"ec2", "describe-instances",
		"--instance-ids", cfg.Machine.InstanceID,
		"--region", cfg.AWS.Region,
		"--query", "Reservations[0].Instances[0].[State.Name,PrivateIpAddress]",
		"--output", "text",
	}
	if cfg.AWS.Profile != "" {
		args = append(args, "--profile", cfg.AWS.Profile)
	}
	if out, err := run.Run(ctx, "aws", args...); err == nil {
		fmt.Fprintf(run.Out, "EC2 state:       %s\n", strings.Join(strings.Fields(out), " "))
	} else {
		rep.Warn("could not describe instance: %v", err)
	}

	if _, err := run.Run(ctx, "ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=10", cfg.SSH.Alias, "true"); err == nil {
		fmt.Fprintln(run.Out, "SSH:             reachable")
	} else {
		rep.Warn("ssh not reachable (run 'devmachine up' or check VPN)")
	}

	if cfg.Machine.PrivateIP != "" {
		if line, err := probeAgent(ctx, cfg); err != nil {
			rep.Warn("remote agent not reachable: %v", err)
		} else {
			fmt.Fprintf(run.Out, "Remote agent:    %s\n", line)
		}
	}
	return nil
}

// probeAgent resolves the remote agent runtime from the recorded machine data
// and queries its status endpoint.
func probeAgent(ctx context.Context, cfg *Config) (string, error) {
	runtime, err := remoteclient.ResolveRuntimeWithInput(zerolog.Nop(), remoteclient.RuntimeInput{
		RemoteHostIP: cfg.Machine.PrivateIP,
		AgentPort:    cfg.Agent.Port,
	})
	if err != nil {
		return "", err
	}
	probeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	status, err := remoteclient.GetAgentStatus(probeCtx, runtime)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("reachable at %s (version %q, %d runtime component(s), uptime %ds)",
		runtime.AgentBaseURL, status.AgentVersion, len(status.RuntimeComponents), status.UptimeSeconds), nil
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
