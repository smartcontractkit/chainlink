// Package dbflags registers the tools/test database flags on testrig's root
// command and builds the config the database helpers consume. The flags are
// bound to package variables so the resource provider (which runs without a
// cobra command in hand) and the persistent-db commands read the same values.
package dbflags

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/pflag"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

var (
	databaseURL     string
	postgresVersion string
)

// Register adds the --database-url and --postgres-version persistent flags. Pass
// it to testrig.WithRootFlags.
func Register(flags *pflag.FlagSet) {
	flags.StringVar(&databaseURL, "database-url", "",
		"Provide a PostgreSQL connection string to use an existing database instead of an ephemeral one")
	flags.StringVar(&postgresVersion, "postgres-version", config.DefaultPostgresVersion,
		"PostgreSQL version to run tests against")
}

// AppConfig builds the database configuration from the parsed flags and the
// current working directory. The cwd is used as the repo root for preparetest,
// so the tool must be run from the directory your package patterns are relative
// to (the repository root). Call it after cobra has parsed flags.
func AppConfig() (*config.App, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}
	version := postgresVersion
	if version == "" {
		version = config.DefaultPostgresVersion
	}
	return &config.App{
		DatabaseURL:     databaseURL,
		PostgresVersion: version,
		RepoRoot:        cwd,
		AIOutput:        !term.IsTerminal(os.Stdout.Fd()),
	}, nil
}
