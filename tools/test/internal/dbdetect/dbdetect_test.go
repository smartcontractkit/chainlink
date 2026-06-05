package dbdetect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "GNUmakefile")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

func TestNeedsPostgres(t *testing.T) {
	repoRoot := findRepoRoot(t)
	t.Logf("repoRoot: %q", repoRoot)

	tests := []struct {
		name    string
		args    []string
		want    bool
		wantErr bool
	}{
		{
			name: "short flag skips DB",
			args: []string{"-short", "./core/services/cron/..."},
			want: false,
		},
		{
			name: "double short flag skips DB",
			args: []string{"--short", "./core/services/cron/..."},
			want: false,
		},
		{
			name: "short with equals skips DB",
			args: []string{"-short=true", "./core/services/cron/..."},
			want: false,
		},
		{
			name: "no patterns defaults to true",
			args: []string{"diagnose"},
			want: true,
		},
		{
			name: "package needing DB",
			args: []string{"./core/services/cron/..."},
			want: true,
		},
		{
			name: "package NOT needing DB",
			args: []string{"./core/config/..."},
			want: false,
		},
		{
			name: "database-url value is not a package pattern",
			args: []string{
				"--database-url",
				"postgres://user:pass@localhost:5432/chainlink_test?sslmode=disable",
				"./core/config/...",
			},
			want: false,
		},
		{
			name: "custom tag unlocks pgtest dependency",
			args: []string{
				"-tags", "dbdetecttag",
				"./core/internal/testutils/dbdetectfixture",
			},
			want: true,
		},
		{
			name: "without custom tag skips pgtest dependency",
			args: []string{"./core/internal/testutils/dbdetectfixture"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NeedsPostgres(repoRoot, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
