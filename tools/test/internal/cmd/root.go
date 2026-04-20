package cmd

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
)

var dbHandle *db.Handle

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "Run Chainlink Go tests with a single command",
	Long: `Run Chainlink Go tests with a single command. Use an ephemeral Postgres database or an existing one.

Modes:

- test: Run tests using vanilla go test command and arguments
- gotestsum: Run tests using gotestsum for those that prefer its output and tools
- survey: Run tests multiple times to collect statistics, debug logs, and more to help find and diagnose flakes, races, panics, timeouts, and other issues`,
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "go -C ./tools/test run .",
	},
	Example: `# Use vanilla go test commands
go -C ./tools/test run . test -v -count=1 -p 4 ./core/...
# Use gotestsum as the runner
go -C ./tools/test run . gotestsum --format=dots -- -count=1 ./core/...
# Run the full core test suite 10 times and collect statistics, debug logs, and more
go -C ./tools/test run . survey --iterations 10 --timeout=15m ./core/...`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.Load(cmd)
		if err != nil {
			return err
		}

		dbHandle, err = db.Ensure(cmd.Context(), conf)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().String("database-url", "", "Provide a PostgreSQL connection string to use an existing database instead of an ephemeral one")
	rootCmd.PersistentFlags().String("postgres-version", config.DefaultPostgresVersion, "PostgreSQL version to run tests against")
	rootCmd.PersistentFlags().Bool("ai-output", !term.IsTerminal(os.Stdout.Fd()), "Use sparse output for agent tooling (and robotic humans)")

	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(gotestsumCmd)
	rootCmd.AddCommand(surveyCmd)
}

// Execute runs the root command.
func Execute() {
	opts := []fang.Option{fang.WithoutCompletions()}
	if err := fang.Execute(context.Background(), rootCmd, opts...); err != nil {
		os.Exit(1)
	}
}
