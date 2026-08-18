package modules_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func TestFindTestPackages(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		files    []string
		expected []string
	}{
		{
			name:     "single go file in core",
			files:    []string{"core/services/telemetry/ingress.go"},
			expected: []string{"./core/services/telemetry"},
		},
		{
			name:     "single test file in core",
			files:    []string{"core/logger/logger_test.go"},
			expected: []string{"./core/logger"},
		},
		{
			name:     "single go file in deployment",
			files:    []string{"deployment/environment.go"},
			expected: []string{"./deployment"},
		},
		{
			name:     "tools package file",
			files:    []string{"tools/ci-testshard/main.go", "tools/githooks/internal/modules/modules.go"},
			expected: []string{"./tools/ci-testshard", "./tools/githooks/internal/modules"},
		},
		{
			name: "skips full E2E suites system-tests and integration-tests and example workflows",
			files: []string{
				"system-tests/lib/suite.go",
				"system-tests/tests/smoke/cre/solana/main_test.go",
				"integration-tests/smoke/vrf_test.go",
				"integration-tests/load/ocr_test.go",
				"core/scripts/cre/environment/examples/workflows/cron/main.go",
			},
			expected: []string{},
		},
		{
			name: "mix of unit test files and E2E files only runs unit test packages",
			files: []string{
				"core/logger/logger.go",
				"deployment/environment.go",
				"tools/ci-testshard/main.go",
				"system-tests/lib/suite.go",
				"integration-tests/smoke/test.go",
			},
			expected: []string{"./core/logger", "./deployment", "./tools/ci-testshard"},
		},
		{
			name:     "root go.mod change triggers core and tools packages",
			files:    []string{"go.mod"},
			expected: []string{"./core/...", "./tools/..."},
		},
		{
			name:     "deployment go.mod change triggers deployment packages",
			files:    []string{"deployment/go.mod"},
			expected: []string{"./deployment/..."},
		},
		{
			name:     "system-tests or integration-tests go.mod change is skipped",
			files:    []string{"system-tests/lib/go.mod", "integration-tests/go.mod"},
			expected: []string{},
		},
		{
			name:     "absolute path handling",
			files:    []string{filepath.Join(tmpDir, "core/services/job/job.go")},
			expected: []string{"./core/services/job"},
		},
		{
			name:     "empty files",
			files:    []string{},
			expected: []string{},
		},
		{
			name:     "ignores non-go files",
			files:    []string{"README.md", "docs/architecture.png"},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pkgs, err := modules.FindTestPackages(tmpDir, tc.files)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, pkgs)
		})
	}
}
