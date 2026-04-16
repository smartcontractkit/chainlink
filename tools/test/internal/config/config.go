package config

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/repo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const DefaultPostgresVersion = "16"

type App struct {
	DatabaseURL     string `mapstructure:"database_url"`
	PostgresVersion string `mapstructure:"postgres_version"`
	RepoRoot        string `mapstructure:"repo_root"`
	AIOutput        bool   `mapstructure:"ai_output"`

	// Function to cleanup the database
	CleanupDB func() error
}

func Load(flags *pflag.FlagSet) (*App, error) {
	if flags == nil {
		return nil, errors.New("flags are required")
	}
	v := viper.New()

	v.SetDefault("postgres_version", DefaultPostgresVersion)
	v.SetDefault("ai_output", !term.IsTerminal(uintptr(os.Stdout.Fd()))) // If TTY (in an AI terminal), use ai-output
	repoRoot, err := repo.RootFromWd()
	if err != nil {
		return nil, err
	}
	v.SetDefault("repo_root", repoRoot)

	flags.VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")
		if bindErr := v.BindPFlag(configName, f); bindErr != nil {
			err = bindErr
			return
		}
	})
	if err != nil {
		return nil, err
	}

	var conf App
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
