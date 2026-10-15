package devmachine

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Cmd is the `devmachine` command group.
var Cmd = &cobra.Command{
	Use:   "devmachine",
	Short: "Provision a remote EC2 DevMachine for hybrid CRE runs",
	Long: `Provision a remote EC2 DevMachine (Docker installed, images loaded, SSH
alias configured) so a hybrid CRE environment can run components remotely.

Privileged users can create a machine; users without permission can pass an
existing instance id and the tool will prepare it.`,
}

func init() {
	Cmd.AddCommand(upCmd())
	Cmd.AddCommand(downCmd())
	Cmd.AddCommand(statusCmd())
	Cmd.AddCommand(handoffCmd())
}

func upCmd() *cobra.Command {
	opts := &UpOptions{}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Create (or adopt) a DevMachine and prepare it for remote CRE",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUp(cmd.Context(), opts)
		},
	}

	f := cmd.Flags()
	f.StringVar(&opts.ConfigPath, "config", "", "path to config file (default ~/.config/cre/devmachine.yaml)")
	f.StringVar(&opts.Name, "name", "", "DevMachine name (default hybrid-local-cre-<user>)")
	f.StringVar(&opts.Namespace, "namespace", "", "Kubernetes namespace for the DevMachine claim")
	f.StringVar(&opts.Team, "team", "", "value for the chain.link/team label")
	f.StringVar(&opts.CostCenter, "cost-center", "", "value for the chain.link/cost-center label")
	f.StringVar(&opts.InstanceType, "instance-type", "", "EC2 instance type (default "+DefaultInstanceType+")")
	f.IntVar(&opts.VolumeSize, "volume-size", 0, "EC2 volume size in GiB (default 200)")
	f.StringVar(&opts.AZ, "az", "", "availability zone suffix (default a)")
	f.StringVar(&opts.InstanceID, "instance-id", "", "use an existing instance instead of creating one")
	f.StringVar(&opts.Region, "region", "", "AWS region (default "+DefaultRegion+")")
	f.StringVar(&opts.AWSProfile, "aws-profile", "", "AWS profile to use")
	f.StringVar(&opts.ChainlinkRepo, "chainlink-repo", "", "path to the chainlink repository")
	f.StringVar(&opts.JobDistributorRepo, "job-distributor-repo", "", "path to the job-distributor repository")
	f.StringVar(&opts.SSHKey, "ssh-key", "", "SSH private key to install (default ~/.ssh/id_ed25519)")
	f.StringVar(&opts.ChainlinkImage, "chainlink-image", "", "chainlink image tag to build/save")
	f.StringVar(&opts.JDImage, "jd-image", "", "job-distributor image tag to build/save")
	f.StringVar(&opts.JDTag, "jd-tag", "", "job-distributor git tag to build (default v0.22.1)")
	f.StringVar(&opts.EnvWorkdir, "env-workdir", "", "directory to run 'go run . env start' from")
	f.StringVar(&opts.EnvConfig, "env-config", "", "CTF_CONFIGS value used by the environment")
	f.BoolVar(&opts.DryRun, "dry-run", false, "print commands without executing them")
	f.BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose command logging")
	f.BoolVar(&opts.Yes, "yes", false, "non-interactive; fail instead of prompting")
	f.StringVar(&opts.From, "from", "", "resume from a phase (P0..P6)")
	f.StringVar(&opts.Only, "only", "", "run a single phase (P0..P6)")

	return cmd
}

func runUp(ctx context.Context, opts *UpOptions) error {
	cfg, err := LoadConfig(opts.ConfigPath)
	if err != nil {
		return err
	}
	applyUpOverrides(cfg, opts)

	run := NewRunner(nil, nil, opts.DryRun, opts.Verbose)

	// Interactive TUI mode: run preflight (which may prompt) in plain mode,
	// then drive the remaining phases inside a Bubble Tea program and print the
	// handoff afterwards.
	if useTUI(opts) {
		o := NewOrchestrator(cfg, opts, run, NewPlainReporter(run.Out), NewPrompter(run.Out))
		if err := o.Preflight(ctx); err != nil {
			return err
		}
		if err := runTUI(ctx, o, phaseMachine, phaseStartAgent); err != nil {
			return err
		}
		fmt.Fprint(run.Out, handoffText(o.cfg))
		return nil
	}

	rep := NewPlainReporter(run.Out)
	prompter := NewPrompter(run.Out)
	o := NewOrchestrator(cfg, opts, run, rep, prompter)
	return o.Up(ctx)
}

// useTUI reports whether the interactive TUI should be used. It is disabled for
// non-interactive, dry-run and phase-selective runs.
func useTUI(opts *UpOptions) bool {
	if opts.Yes || opts.DryRun || opts.From != "" || opts.Only != "" {
		return false
	}
	return isTerminal(os.Stdout)
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func applyUpOverrides(cfg *Config, opts *UpOptions) {
	if opts.Name != "" {
		cfg.Machine.Name = opts.Name
	}
	if opts.Namespace != "" {
		cfg.Cluster.Namespace = opts.Namespace
	}
	if opts.InstanceType != "" {
		cfg.Machine.InstanceType = opts.InstanceType
	}
	if opts.VolumeSize != 0 {
		cfg.Machine.VolumeSize = opts.VolumeSize
	}
	if opts.AZ != "" {
		cfg.Machine.AvailabilityZ = opts.AZ
	}
	if opts.InstanceID != "" {
		cfg.Machine.InstanceID = opts.InstanceID
	}
	if opts.Region != "" {
		cfg.AWS.Region = opts.Region
	}
	if opts.AWSProfile != "" {
		cfg.AWS.Profile = opts.AWSProfile
	}
	if opts.SSHKey != "" {
		cfg.SSH.Key = opts.SSHKey
	}
	if opts.ChainlinkRepo != "" {
		cfg.Repos.Chainlink = opts.ChainlinkRepo
	}
	if opts.JobDistributorRepo != "" {
		cfg.Repos.JobDistributor = opts.JobDistributorRepo
	}
	if opts.ChainlinkImage != "" {
		cfg.Images.Chainlink = opts.ChainlinkImage
	}
	if opts.JDImage != "" {
		cfg.Images.JobDistributor = opts.JDImage
	}
	if opts.JDTag != "" {
		cfg.Images.JDTag = opts.JDTag
	}
	if opts.EnvWorkdir != "" {
		cfg.Env.Workdir = opts.EnvWorkdir
	}
	if opts.EnvConfig != "" {
		cfg.Env.Config = opts.EnvConfig
	}
}

func downCmd() *cobra.Command {
	opts := &DownOptions{}
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Delete the DevMachine (if created by this tool) and clean up",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := LoadConfig(opts.ConfigPath)
			if err != nil {
				return err
			}
			run := NewRunner(nil, nil, opts.DryRun, opts.Verbose)
			rep := NewPlainReporter(run.Out)
			return Down(cmd.Context(), cfg, opts, run, rep)
		},
	}
	f := cmd.Flags()
	f.StringVar(&opts.ConfigPath, "config", "", "path to config file")
	f.BoolVar(&opts.PurgeCache, "purge-cache", false, "also remove the local artifact cache")
	f.BoolVar(&opts.KeepMachine, "keep-machine", false, "do not delete the DevMachine")
	f.BoolVar(&opts.DryRun, "dry-run", false, "print commands without executing them")
	f.BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose command logging")
	return cmd
}

func statusCmd() *cobra.Command {
	opts := &StatusOptions{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show configuration and live status of the DevMachine",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := LoadConfig(opts.ConfigPath)
			if err != nil {
				return err
			}
			run := NewRunner(nil, nil, false, false)
			rep := NewPlainReporter(run.Out)
			return Status(cmd.Context(), cfg, opts, run, rep)
		},
	}
	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "path to config file")
	return cmd
}

func handoffCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "handoff",
		Short: "Print the manual commands to start the agent and environment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := LoadConfig(configPath)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), handoffText(cfg))
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "path to config file")
	return cmd
}
