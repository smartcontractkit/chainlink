package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/repo"
)

const DefaultPostgresVersion = "16"

type App struct {
	DatabaseURL        string        `mapstructure:"database_url"`
	PostgresVersion    string        `mapstructure:"postgres_version"`
	RepoRoot           string        `mapstructure:"repo_root"`
	AIOutput           bool          `mapstructure:"ai_output"`
	Iterations         int           `mapstructure:"iterations"`
	ParallelIterations int           `mapstructure:"parallel_iterations"`
	SlowThreshold      time.Duration `mapstructure:"slow_threshold"`
	FailFast           bool          `mapstructure:"fail_fast"`
	FailFastOn         []string      `mapstructure:"fail_fast_on"`
	Shuffle            bool          `mapstructure:"shuffle_seed"`
}

const (
	FailFastOnAny     = "any"
	FailFastOnFailure = "failure"
	FailFastOnTimeout = "timeout"
	FailFastOnSlow    = "slow"
)

var validFailFastOn = map[string]struct{}{
	FailFastOnAny:     {},
	FailFastOnFailure: {},
	FailFastOnTimeout: {},
	FailFastOnSlow:    {},
}

// NormalizeFailFastOn validates --fail-fast-on values and returns lowercase,
// de-duplicated categories in first-seen order.
func NormalizeFailFastOn(values []string) ([]string, error) {
	var out []string
	seen := make(map[string]struct{})
	for _, value := range values {
		for part := range strings.SplitSeq(value, ",") {
			category := strings.ToLower(strings.TrimSpace(part))
			if category == "" {
				return nil, errors.New(`--fail-fast-on must contain only "any", "failure", "timeout", or "slow"; got ""`)
			}
			if _, ok := validFailFastOn[category]; !ok {
				return nil, fmt.Errorf(`--fail-fast-on must contain only "any", "failure", "timeout", or "slow"; got %q`, category)
			}
			if _, ok := seen[category]; ok {
				continue
			}
			seen[category] = struct{}{}
			out = append(out, category)
		}
	}
	return out, nil
}

// Load binds Viper to the active command's persistent flags and local flags, then unmarshals into App.
func Load(cmd *cobra.Command) (*App, error) {
	if cmd == nil {
		return nil, errors.New("command is required")
	}
	v := viper.New()

	v.SetDefault("postgres_version", DefaultPostgresVersion)
	// Enable sparse output when stdout is not a TTY (e.g. redirected or CI).
	v.SetDefault("ai_output", !term.IsTerminal(os.Stdout.Fd()))
	v.SetDefault("iterations", 1)
	v.SetDefault("parallel_iterations", 1)
	v.SetDefault("slow_threshold", 30*time.Second)
	v.SetDefault("fail_fast", false)
	v.SetDefault("fail_fast_on", []string{})
	repoRoot, err := repo.RootFromWd()
	if err != nil {
		return nil, err
	}
	v.SetDefault("repo_root", repoRoot)

	if err := bindPFlags(v, cmd.PersistentFlags()); err != nil {
		return nil, err
	}
	if err := bindPFlags(v, cmd.Flags()); err != nil {
		return nil, err
	}

	var conf App
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func bindPFlags(v *viper.Viper, flags *pflag.FlagSet) error {
	var err error
	flags.VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")
		if bindErr := v.BindPFlag(configName, f); bindErr != nil {
			err = bindErr
		}
	})
	return err
}
