package updater

import (
	"flag"
	"fmt"
)

const (
	DefaultRepoRemote  = "origin"
	DefaultBranchTrunk = "develop"
	DefaultOrgName     = "smartcontractkit"
	DefaultRepoName    = "chainlink"
)

type Config struct {
	RepoRemote  string
	BranchTrunk string
	DryRun      bool
	ShowVersion bool
	OrgName     string
	RepoName    string
}

func ParseFlags(args []string, version string) (*Config, error) {
	flags := flag.NewFlagSet("default", flag.ContinueOnError)

	cfg := &Config{}

	flags.StringVar(&cfg.RepoRemote, "repo-remote", DefaultRepoRemote, "Git remote to use")
	flags.StringVar(&cfg.BranchTrunk, "branch-trunk", DefaultBranchTrunk, "Branch to get SHA from")
	flags.BoolVar(&cfg.DryRun, "dry-run", false, "Preview changes without applying them")
	flags.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")
	flags.StringVar(&cfg.OrgName, "org-name", DefaultOrgName, "GitHub organization name")
	flags.StringVar(&cfg.RepoName, "repo-name", DefaultRepoName, "GitHub repository name")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ShowVersion {
		return nil
	}

	if c.OrgName == "" {
		return fmt.Errorf("%w: org name must be provided", ErrInvalidConfig)
	}
	if c.RepoName == "" {
		return fmt.Errorf("%w: repo name must be provided", ErrInvalidConfig)
	}
	if c.RepoRemote == "" {
		return fmt.Errorf("%w: repo remote must be provided", ErrInvalidConfig)
	}
	if c.BranchTrunk == "" {
		return fmt.Errorf("%w: trunk branch must be provided", ErrInvalidConfig)
	}

	if !gitRemoteRE.MatchString(c.RepoRemote) {
		return fmt.Errorf("%w: git remote '%s' contains invalid characters", ErrInvalidConfig, c.RepoRemote)
	}
	if !gitBranchRE.MatchString(c.BranchTrunk) {
		return fmt.Errorf("%w: git branch '%s' contains invalid characters", ErrInvalidConfig, c.BranchTrunk)
	}

	return nil
}
