package cmd

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/repo"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const defaultPostgresVersion = "16"

var (
	databaseURL     string
	postgresVersion string
	repoRoot        string
	aiOutput        bool
)

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "Run Chainlink Go tests with optional Postgres from testcontainers",
	Long:  `test enables you to run chainlink tests with a single command and ephemeral Postgres database.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		repoRoot, err = repo.RootFromWd()
		if err != nil {
			return err
		}
		cleanup, err := db.Ensure(cmd.Context(), databaseURL, postgresVersion, repoRoot)
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(cmd.Context(), "cleanup", cleanup))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "Provide a PostgreSQL connection string to use an existing database instead of an ephemeral one")
	rootCmd.PersistentFlags().StringVar(&postgresVersion, "postgres-version", defaultPostgresVersion, "PostgreSQL version to run tests against")
	rootCmd.PersistentFlags().BoolVar(&aiOutput, "ai-output", !term.IsTerminal(int(os.Stdout.Fd())), "Use sparse output for agent tooling (and robotic humans)")

	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(gotestsumCmd)
	rootCmd.AddCommand(surveyCmd)
}

// Execute runs the root command.
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
