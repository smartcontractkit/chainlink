package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

func TestRootCommandPathShowsGoCInvocation(t *testing.T) {
	t.Parallel()

	if got := rootCmd.CommandPath(); got != "go -C tools/test run ." {
		t.Fatalf("root CommandPath (help / errors): got %q want %q", got, "go -C tools/test run .")
	}
	if got := rootCmd.DisplayName(); got != "go -C tools/test run ." {
		t.Fatalf("DisplayName: got %q want %q", got, "go -C tools/test run .")
	}
	if got := rootCmd.Name(); got != "test" {
		t.Fatalf("internal Name (subcommand paths use CommandPath + Name): got %q want %q", got, "test")
	}
}

func TestSubcommandCommandPaths(t *testing.T) {
	t.Parallel()

	var gotestsum *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "gotestsum" {
			gotestsum = c
			break
		}
	}
	if gotestsum == nil {
		t.Fatal("gotestsum subcommand not found")
	}
	want := "go -C tools/test run . gotestsum"
	if got := gotestsum.CommandPath(); got != want {
		t.Fatalf("gotestsum CommandPath: got %q want %q", got, want)
	}
}

func TestValidateDiagnoseConfigParallelIterations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		conf    *config.App
		wantErr string
	}{
		{
			name: "valid default",
			conf: &config.App{Iterations: 1, ParallelIterations: 1},
		},
		{
			name:    "invalid iterations",
			conf:    &config.App{Iterations: 0, ParallelIterations: 1},
			wantErr: "--iterations must be >= 1",
		},
		{
			name:    "invalid parallel iterations",
			conf:    &config.App{Iterations: 1, ParallelIterations: 0},
			wantErr: "--parallel-iterations must be >= 1",
		},
		{
			name:    "parallel iterations cannot exceed iterations",
			conf:    &config.App{Iterations: 2, ParallelIterations: 3},
			wantErr: "--parallel-iterations must be <= --iterations",
		},
		{
			name:    "external database rejected for parallel",
			conf:    &config.App{Iterations: 10, ParallelIterations: 2, DatabaseURL: "postgres://example/db"},
			wantErr: "--parallel-iterations > 1 cannot be used with --database-url",
		},
		{
			name:    "invalid fail fast category",
			conf:    &config.App{Iterations: 1, ParallelIterations: 1, FailFastOn: []string{"timeout", "banana"}},
			wantErr: `--fail-fast-on must contain only "any", "failure", "timeout", or "slow"; got "banana"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validateDiagnoseConfig(tc.conf)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}
