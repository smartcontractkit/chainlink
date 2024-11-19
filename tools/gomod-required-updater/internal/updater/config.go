package updater

import (
	"flag"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type TOMLConfig struct {
	Modules []string `toml:"modules"`
}

type Config struct {
	ModulesToUpdate []string
	RepoRemote     string
	BranchTrunk    string
	DryRun         bool
	ConfigFile     string
	RootPath       string
	ShowVersion    bool
	modulesSource  string
	UpdateOrgModules bool   // Update modules from same org/repo with local replaces
	OrgName         string // GitHub organization name
	RepoName        string // Repository name
}

func (c *Config) Validate() error {
    if !c.ShowVersion {
        // Skip module validation if using UpdateOrgModules
        if !c.UpdateOrgModules && len(c.ModulesToUpdate) == 0 {
            return fmt.Errorf("%w: no modules specified to update (use -module flag or config file)", ErrInvalidConfig)
        }
        if c.RepoRemote == "" {
            return fmt.Errorf("%w: repo remote cannot be empty", ErrInvalidConfig)
        }
        if c.BranchTrunk == "" {
            return fmt.Errorf("%w: branch trunk cannot be empty", ErrInvalidConfig)
        }
    }
    return nil
}

func ParseFlags(args []string, version string) (*Config, error) {
	flags := flag.NewFlagSet("gomod-required-updater", flag.ContinueOnError)
	
	cfg := &Config{}
	var cliModules arrayFlags // Define custom flag type for multiple -module flags
	
	flags.StringVar(&cfg.RepoRemote, "repo-remote", "origin", "The name of the repo remote")
	flags.StringVar(&cfg.BranchTrunk, "branch-trunk", "develop", "The name of the trunk branch")
	flags.BoolVar(&cfg.DryRun, "dry-run", false, "Print what would be done without making changes")
	flags.StringVar(&cfg.ConfigFile, "config", "", "Path to TOML config file (optional if using -module)")
	flags.StringVar(&cfg.RootPath, "root", ".", "Root path to start scanning for go.mod files")
	flags.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")
	flags.Var(&cliModules, "module", "Module to update (can be specified multiple times)")
	flags.BoolVar(&cfg.UpdateOrgModules, "update-org-modules", false, "Update modules from same org/repo that have local replace directives")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	// CLI modules take precedence
	if len(cliModules) > 0 {
		cfg.ModulesToUpdate = cliModules
		cfg.modulesSource = "command line"
	} else if cfg.ConfigFile != "" {
		// Only load config file if no CLI modules specified
		modules, err := loadTOMLConfig(cfg.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
		cfg.ModulesToUpdate = modules
		cfg.modulesSource = "config file"
	}

	if cfg.UpdateOrgModules {
		gitOp := NewGitOperator()
		org, repo, err := gitOp.GetRepoInfo(cfg.RepoRemote)
		if err != nil {
			return nil, fmt.Errorf("failed to get repo info: %w", err)
		}
		cfg.OrgName = org
		cfg.RepoName = repo
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// arrayFlags allows for repeated flag values
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func loadTOMLConfig(path string) ([]string, error) {
	var cfg TOMLConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode TOML: %w", err)
	}
	return cfg.Modules, nil
}