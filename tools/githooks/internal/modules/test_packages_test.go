package modules_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func TestFindTestModules(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module root\n"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "tools/githooks"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "tools/githooks/go.mod"), []byte("module githooks\n"), 0o600))

	tests := []struct {
		name     string
		files    []string
		expected []modules.ModulePackages
	}{
		{
			name: "files across root module and tools/githooks submodule",
			files: []string{
				"core/services/gateway/gateway.go",
				"tools/githooks/cmd/test.go",
			},
			expected: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./core/services/gateway"},
				},
				{
					Module:   "tools/githooks",
					Packages: []string{"./cmd"},
				},
			},
		},
		{
			name: "skips excluded E2E modules",
			files: []string{
				"deployment/environment.go",
				"system-tests/lib/suite.go",
				"core/logger/logger.go",
			},
			expected: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./core/logger"},
				},
			},
		},
		{
			name: "root module dependency change runs tests on entire root module",
			files: []string{
				"go.mod",
			},
			expected: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./..."},
				},
			},
		},
		{
			name: "submodule dependency change runs tests on entire submodule",
			files: []string{
				"tools/githooks/go.sum",
			},
			expected: []modules.ModulePackages{
				{
					Module:   "tools/githooks",
					Packages: []string{"./..."},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mods, err := modules.FindTestModules(tmpDir, tc.files)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, mods)
		})
	}
}
