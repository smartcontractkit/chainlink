package dbflags

import (
	"os"
	"testing"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

func resetDBFlags(t *testing.T) {
	t.Helper()
	t.Setenv("CL_DATABASE_URL", "")
	t.Cleanup(func() {
		databaseURL = ""
		postgresVersion = ""
	})
}

func parseDBFlags(t *testing.T, args ...string) {
	t.Helper()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Register(flags)
	require.NoError(t, flags.Parse(args))
}

func TestAppConfig(t *testing.T) {
	tests := []struct {
		name          string
		envDBURL      string
		args          []string
		wantDBURL     string
		wantPGVersion string
	}{
		{
			name:          "defaults",
			wantPGVersion: config.DefaultPostgresVersion,
		},
		{
			name:          "database_url",
			args:          []string{"--database-url", "postgres://example"},
			wantDBURL:     "postgres://example",
			wantPGVersion: config.DefaultPostgresVersion,
		},
		{
			name:          "database_url_from_env",
			envDBURL:      "postgres://from-env",
			wantDBURL:     "postgres://from-env",
			wantPGVersion: config.DefaultPostgresVersion,
		},
		{
			name:          "database_url_flag_overrides_env",
			envDBURL:      "postgres://from-env",
			args:          []string{"--database-url", "postgres://from-flag"},
			wantDBURL:     "postgres://from-flag",
			wantPGVersion: config.DefaultPostgresVersion,
		},
		{
			name:          "postgres_version",
			args:          []string{"--postgres-version", "15"},
			wantPGVersion: "15",
		},
		{
			name:          "empty_postgres_version_uses_default",
			args:          []string{"--postgres-version", ""},
			wantPGVersion: config.DefaultPostgresVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetDBFlags(t)
			if tt.envDBURL != "" {
				t.Setenv("CL_DATABASE_URL", tt.envDBURL)
			}
			parseDBFlags(t, tt.args...)

			dir := t.TempDir()
			t.Chdir(dir)

			conf, err := AppConfig()
			require.NoError(t, err)
			assert.Equal(t, tt.wantDBURL, conf.DatabaseURL)
			assert.Equal(t, tt.wantPGVersion, conf.PostgresVersion)
			assert.Equal(t, dir, conf.RepoRoot)
			assert.Equal(t, !term.IsTerminal(os.Stdout.Fd()), conf.AIOutput)
		})
	}
}
