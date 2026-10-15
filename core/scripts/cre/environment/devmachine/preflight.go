package devmachine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// requiredTools are the external binaries the tool shells out to.
var requiredTools = []struct {
	bin  string
	hint string
}{
	{"aws", "install AWS CLI v2 and run 'aws sso login'"},
	{"kubectl", "install kubectl and select a context that can reach the DevMachine CRD"},
	{"docker", "install Docker Desktop / docker engine with the buildx plugin"},
	{"go", "install Go (needed to build the remote agent)"},
	{"rsync", "install rsync (macOS ships it; AL2023 remote gets it in P2)"},
	{"ssh", "install an OpenSSH client"},
}

// Preflight performs phase P0: verify prerequisites, resolve the AWS profile,
// cluster access, repository locations and machine naming, and persist the
// merged configuration.
func (o *Orchestrator) Preflight(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phasePreflight))

	if err := o.checkTools(); err != nil {
		return err
	}
	if err := o.resolveAWS(ctx); err != nil {
		return err
	}
	if err := o.resolveCluster(ctx); err != nil {
		return err
	}
	if err := o.resolveRepos(); err != nil {
		return err
	}
	if err := o.resolveNaming(); err != nil {
		return err
	}
	return o.save()
}

func (o *Orchestrator) checkTools() error {
	var missing []string
	for _, t := range requiredTools {
		if _, err := o.run.LookPath(t.bin); err != nil {
			missing = append(missing, fmt.Sprintf("  - %s: %s", t.bin, t.hint))
		}
	}
	if len(missing) > 0 {
		return errors.Errorf("missing required tools:\n%s", strings.Join(missing, "\n"))
	}
	// buildx is a docker plugin, verify it is present.
	if _, err := o.run.Run(context.Background(), "docker", "buildx", "version"); err != nil {
		return errors.New("docker buildx plugin not available; install it to build amd64 images")
	}
	o.rep.Info("all required tools present")
	return nil
}

func (o *Orchestrator) resolveAWS(ctx context.Context) error {
	if o.opts.AWSProfile != "" {
		o.cfg.AWS.Profile = o.opts.AWSProfile
	}
	if o.opts.Region != "" {
		o.cfg.AWS.Region = o.opts.Region
	}

	args := []string{"sts", "get-caller-identity", "--region", o.cfg.AWS.Region, "--output", "json"}
	if o.cfg.AWS.Profile != "" {
		args = append(args, "--profile", o.cfg.AWS.Profile)
	}
	if o.run.DryRun {
		fmt.Fprintf(o.run.Out, "[dry-run] aws %s\n", strings.Join(args, " "))
		return nil
	}
	out, err := o.run.Run(ctx, "aws", args...)
	if err != nil {
		hint := "run 'aws sso login"
		if o.cfg.AWS.Profile != "" {
			hint += " --profile " + o.cfg.AWS.Profile
		}
		hint += "'"
		return errors.Wrapf(err, "AWS credentials are not valid; %s", hint)
	}
	o.rep.Info("AWS identity: %s", oneLine(out))
	return nil
}

func (o *Orchestrator) resolveCluster(ctx context.Context) error {
	if o.opts.Namespace != "" {
		o.cfg.Cluster.Namespace = o.opts.Namespace
	}

	if !o.run.DryRun {
		out, err := o.run.Run(ctx, "kubectl", "config", "current-context")
		if err != nil {
			return errors.Wrap(err, "failed to read current kubectl context")
		}
		o.rep.Info("kubectl context: %s", out)
	}

	// Determine whether the caller may create DevMachines.
	out, err := o.run.Run(ctx, "kubectl", "auth", "can-i", "create", "devmachines", "-n", o.cfg.Cluster.Namespace)
	privileged := err == nil && outputSaysYes(out)
	if o.run.DryRun {
		privileged = o.cfg.Cluster.Privileged
	}
	o.cfg.Cluster.Privileged = privileged
	if privileged {
		o.rep.Info("permission to create DevMachines: yes")
	} else {
		o.rep.Info("permission to create DevMachines: no (will use an existing instance)")
	}
	return nil
}

func (o *Orchestrator) resolveRepos() error {
	chainlink, err := o.resolveRepo(
		"chainlink",
		"CRE_CHAINLINK_REPO",
		o.opts.ChainlinkRepo,
		o.cfg.Repos.Chainlink,
		detectChainlinkRepo,
	)
	if err != nil {
		return err
	}
	o.cfg.Repos.Chainlink = chainlink

	jd, err := o.resolveRepo(
		"job-distributor",
		"CRE_JOB_DISTRIBUTOR_REPO",
		o.opts.JobDistributorRepo,
		o.cfg.Repos.JobDistributor,
		func() string { return detectJobDistributorRepo(chainlink) },
	)
	if err != nil {
		return err
	}
	o.cfg.Repos.JobDistributor = jd

	// env workdir defaults relative to the chainlink repo.
	if o.opts.EnvWorkdir != "" {
		o.cfg.Env.Workdir = o.opts.EnvWorkdir
	} else if o.cfg.Env.Workdir == DefaultEnvWorkdir {
		o.cfg.Env.Workdir = filepath.Join(chainlink, DefaultEnvWorkdir)
	}
	if o.opts.EnvConfig != "" {
		o.cfg.Env.Config = o.opts.EnvConfig
	}
	return nil
}

func (o *Orchestrator) resolveRepo(label, envVar, flag, current string, detect func() string) (string, error) {
	// 1. explicit flag
	if flag != "" {
		return o.requireRepoDir(label, flag)
	}
	// 2. env var
	if v := strings.TrimSpace(os.Getenv(envVar)); v != "" {
		return o.requireRepoDir(label, v)
	}
	// 3. persisted config
	if current != "" {
		if isRepoDir(current) {
			return current, nil
		}
	}
	// 4. auto-detect
	if detected := detect(); detected != "" {
		o.rep.Info("detected %s repo at %s", label, detected)
		return detected, nil
	}
	// 5. prompt
	if o.opts.Yes {
		return "", errors.Errorf("could not locate the %s repository; set %s or pass --%s-repo", label, envVar, label)
	}
	val, err := o.prompter.Ask(fmt.Sprintf("Path to the %s repository", label), "")
	if err != nil {
		return "", err
	}
	return o.requireRepoDir(label, val)
}

func (o *Orchestrator) requireRepoDir(label, dir string) (string, error) {
	dir = ExpandPath(strings.TrimSpace(dir))
	if dir == "" {
		return "", errors.Errorf("%s repository path is required", label)
	}
	if !isRepoDir(dir) {
		return "", errors.Errorf("%s repository path %q does not look like a git checkout", label, dir)
	}
	return dir, nil
}

func isRepoDir(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir) //nolint:gosec // dir comes from user config, not attacker-controlled input
	if err != nil || !info.IsDir() {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, ".git")) //nolint:gosec // see above
	return err == nil
}

// detectChainlinkRepo walks up from the working directory looking for a
// chainlink checkout, then falls back to common locations.
func detectChainlinkRepo() string {
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for {
			if looksLikeChainlinkRepo(dir) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return firstExistingRepo(commonRepoNames("chainlink")...)
}

func looksLikeChainlinkRepo(dir string) bool {
	if !isRepoDir(dir) {
		return false
	}
	_, errDocker := os.Stat(filepath.Join(dir, "core", "chainlink.Dockerfile"))
	_, errLib := os.Stat(filepath.Join(dir, "system-tests", "lib", "go.mod"))
	return errDocker == nil && errLib == nil
}

// detectJobDistributorRepo looks for a job-distributor checkout, preferring a
// sibling of the chainlink repo.
func detectJobDistributorRepo(chainlinkRepo string) string {
	var candidates []string
	if chainlinkRepo != "" {
		sibling := filepath.Join(filepath.Dir(chainlinkRepo), "job-distributor")
		candidates = append(candidates, sibling)
	}
	candidates = append(candidates, commonRepoNames("job-distributor")...)
	return firstExistingRepo(candidates...)
}

func commonRepoNames(name string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	bases := []string{
		filepath.Join(home, name),
		filepath.Join(home, "Desktop", "repos", name),
		filepath.Join(home, "src", name),
		filepath.Join(home, "repos", name),
		filepath.Join(home, "code", name),
		filepath.Join(home, "projects", name),
	}
	return bases
}

func firstExistingRepo(candidates ...string) string {
	for _, c := range candidates {
		if looksLikeChainlinkRepo(c) {
			return c
		}
		if isRepoDir(c) {
			return c
		}
	}
	return ""
}

func (o *Orchestrator) resolveNaming() error {
	if o.opts.Name != "" {
		o.cfg.Machine.Name = o.opts.Name
	} else if o.cfg.Machine.Name == "" {
		user := os.Getenv("USER")
		if user == "" {
			user = os.Getenv("LOGNAME")
		}
		o.cfg.Machine.Name = MachineName(user)
	}
	if o.cfg.SSH.Alias == "" {
		o.cfg.SSH.Alias = o.cfg.Machine.Name
	}
	if o.opts.SSHKey != "" {
		o.cfg.SSH.Key = o.opts.SSHKey
	}

	if o.cfg.Cluster.Privileged {
		if o.opts.Team != "" {
			o.cfg.Machine.Team = o.opts.Team
		}
		if o.opts.CostCenter != "" {
			o.cfg.Machine.CostCenter = o.opts.CostCenter
		}
		if o.cfg.Machine.Team == "" && !o.opts.Yes {
			team, err := o.prompter.Ask("Team label (chain.link/team)", "")
			if err != nil {
				return err
			}
			o.cfg.Machine.Team = team
		}
		if o.cfg.Machine.CostCenter == "" {
			if o.cfg.Machine.Team != "" {
				o.cfg.Machine.CostCenter = o.cfg.Machine.Team
			} else if !o.opts.Yes {
				cc, err := o.prompter.Ask("Cost center label (chain.link/cost-center)", "")
				if err != nil {
					return err
				}
				o.cfg.Machine.CostCenter = cc
			}
		}
	}
	return nil
}

func oneLine(s string) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(s) > 200 {
		s = s[:200] + "..."
	}
	return s
}

// outputSaysYes reports whether any line of the output is exactly "yes". This
// tolerates advisory lines that some clusters print alongside the answer.
func outputSaysYes(out string) bool {
	for line := range strings.SplitSeq(out, "\n") {
		if strings.TrimSpace(line) == "yes" {
			return true
		}
	}
	return false
}
