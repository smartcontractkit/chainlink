package modules_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func TestFindAffectedModules(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	modDirs := []string{
		"", // root module
		"deployment",
		"integration-tests",
		"system-tests/lib",
		"core/scripts/cre/workflows/cron",
	}

	for _, mod := range modDirs {
		modPath := filepath.Join(tmpDir, mod)
		require.NoError(t, os.MkdirAll(modPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(modPath, "go.mod"), []byte("module test\n"), 0600))
	}

	tests := []struct {
		name     string
		files    []string
		expected []modules.ModulePackages
	}{
		{
			name:  "root module package file",
			files: []string{"core/services/app.go"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/services"}},
			},
		},
		{
			name:  "root module root-level file",
			files: []string{"main.go"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"."}},
			},
		},
		{
			name:  "submodule package file",
			files: []string{"deployment/environment/env.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"./environment"}},
			},
		},
		{
			name:  "submodule root-level file",
			files: []string{"deployment/main.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"."}},
			},
		},
		{
			name:  "nested submodule package file",
			files: []string{"core/scripts/cre/workflows/cron/pkg/main.go"},
			expected: []modules.ModulePackages{
				{Module: "core/scripts/cre/workflows/cron", Packages: []string{"./pkg"}},
			},
		},
		{
			name: "multiple files deduplicate packages within module",
			files: []string{
				"core/services/app.go",
				"core/services/db.go",
				"core/logger/log.go",
				"deployment/environment/env.go",
				"deployment/environment/node.go",
				"integration-tests/smoke/test.go",
			},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/logger", "./core/services"}},
				{Module: "deployment", Packages: []string{"./environment"}},
				{Module: "integration-tests", Packages: []string{"./smoke"}},
			},
		},
		{
			name:  "go.mod change triggers all packages in module",
			files: []string{"deployment/go.mod", "deployment/environment/env.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"./..."}},
			},
		},
		{
			name:  "go.sum change triggers all packages in module",
			files: []string{"go.sum"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./..."}},
			},
		},
		{
			name:  "absolute path handling",
			files: []string{filepath.Join(tmpDir, "system-tests/lib/suite/test.go")},
			expected: []modules.ModulePackages{
				{Module: "system-tests/lib", Packages: []string{"./suite"}},
			},
		},
		{
			name:     "empty files slice",
			files:    []string{},
			expected: []modules.ModulePackages{},
		},
		{
			name:     "ignores non-go files",
			files:    []string{"README.md", "docs/architecture.png"},
			expected: []modules.ModulePackages{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mods, err := modules.FindAffectedModules(tmpDir, tc.files)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, mods)
		})
	}
}
