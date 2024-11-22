package updater

import (
	"flag"
)

type Config struct {
	ModulesToUpdate []string
	RepoRemote      string
	BranchTrunk     string
	DryRun          bool
	ShowVersion     bool
	OrgName         string
	RepoName        string
}

func ParseFlags(args []string, version string) (*Config, error) {
	flags := flag.NewFlagSet("gomod-required-updater", flag.ContinueOnError)

	cfg := &Config{
		ModulesToUpdate: make([]string, 0),
	}

	flags.StringVar(&cfg.RepoRemote, "repo-remote", "origin", "Git remote to use")
	flags.StringVar(&cfg.BranchTrunk, "branch-trunk", "develop", "Branch to get SHA from")
	flags.BoolVar(&cfg.DryRun, "dry-run", false, "Preview changes without applying them")
	flags.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")
	flags.StringVar(&cfg.OrgName, "org-name", "smartcontractkit", "GitHub organization name")
	flags.StringVar(&cfg.RepoName, "repo-name", "chainlink", "GitHub repository name")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ShowVersion {
		return nil
	}
	return nil
}
