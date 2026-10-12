package tools_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/tools"
)

func TestDiscoverTargets_SyntheticDir(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Submodule 1 with go.mod and tests
	submod1 := filepath.Join(tmpDir, "tools", "submod1")
	require.NoError(t, os.MkdirAll(submod1, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(submod1, "go.mod"), []byte("module submod1\n\ngo 1.26\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(submod1, "foo_test.go"), []byte("package submod1\n"), 0o600))

	// Submodule 2 with go.mod
	submod2 := filepath.Join(tmpDir, "tools", "submod2")
	require.NoError(t, os.MkdirAll(submod2, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(submod2, "go.mod"), []byte("module submod2\n\ngo 1.26\n"), 0o600))

	// Root module tool package with test
	rootPkg := filepath.Join(tmpDir, "tools", "rootpkg")
	require.NoError(t, os.MkdirAll(rootPkg, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(rootPkg, "bar_test.go"), []byte("package rootpkg\n"), 0o600))

	// Non-test directory
	docDir := filepath.Join(tmpDir, "tools", "docs")
	require.NoError(t, os.MkdirAll(docDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(docDir, "README.md"), []byte("# Docs\n"), 0o600))

	targets, err := tools.DiscoverTargets(tmpDir, filepath.Join(tmpDir, "tools"))
	require.NoError(t, err)

	expected := []tools.Target{
		{
			Name:        "tools/root",
			Dir:         ".",
			Packages:    "./tools/...",
			IsSubmodule: false,
		},
		{
			Name:        "tools/submod1",
			Dir:         "tools/submod1",
			Packages:    "./...",
			IsSubmodule: true,
		},
		{
			Name:        "tools/submod2",
			Dir:         "tools/submod2",
			Packages:    "./...",
			IsSubmodule: true,
		},
	}

	assert.Equal(t, expected, targets)
}

func TestComputeMatrix_AllFlag(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
		{Name: "tools/githooks", Dir: "tools/githooks", Packages: "./...", IsSubmodule: true},
	}

	matrix := tools.ComputeMatrix(targets, tools.MatrixOptions{
		All: true,
	})

	assert.Equal(t, targets, matrix)
}

func TestComputeMatrix_ScheduleAndWorkflowDispatch(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
	}

	matrixSchedule := tools.ComputeMatrix(targets, tools.MatrixOptions{
		EventName: "schedule",
	})
	assert.Equal(t, targets, matrixSchedule)

	matrixDispatch := tools.ComputeMatrix(targets, tools.MatrixOptions{
		EventName: "workflow_dispatch",
	})
	assert.Equal(t, targets, matrixDispatch)
}

func TestComputeMatrix_WorkflowFilesChanged(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
	}

	matrix := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{".github/workflows/ci-core.yml"},
	})
	assert.Equal(t, targets, matrix)

	matrixOtherWorkflow := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{".github/workflows/cre-system-tests.yaml"},
	})
	assert.Equal(t, targets, matrixOtherWorkflow)

	matrixAction := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{".github/actions/setup-go/action.yml"},
	})
	assert.Equal(t, targets, matrixAction)
}

func TestComputeMatrix_SubmoduleChanged(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
		{Name: "tools/githooks", Dir: "tools/githooks", Packages: "./...", IsSubmodule: true},
	}

	matrix := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{"tools/ci/cmd/version.go"},
	})
	require.Len(t, matrix, 1)
	assert.Equal(t, "tools/ci", matrix[0].Name)

	matrixGithooks := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{"tools/githooks/internal/eof/eof.go"},
	})
	require.Len(t, matrixGithooks, 1)
	assert.Equal(t, "tools/githooks", matrixGithooks[0].Name)
}

func TestComputeMatrix_RootToolsChanged(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
	}

	matrixTxtar := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{"tools/txtar/visitor.go"},
	})
	require.Len(t, matrixTxtar, 1)
	assert.Equal(t, "tools/root", matrixTxtar[0].Name)

	matrixGoMod := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{"go.mod"},
	})
	require.Len(t, matrixGoMod, 1)
	assert.Equal(t, "tools/root", matrixGoMod[0].Name)
}

func TestComputeMatrix_UnrelatedChanges(t *testing.T) {
	t.Parallel()
	targets := []tools.Target{
		{Name: "tools/root", Dir: ".", Packages: "./tools/..."},
		{Name: "tools/ci", Dir: "tools/ci", Packages: "./...", IsSubmodule: true},
	}

	matrix := tools.ComputeMatrix(targets, tools.MatrixOptions{
		ChangedFiles: []string{"core/services/chainlink.go", "deployment/environment.go"},
	})
	assert.Empty(t, matrix)
	assert.NotNil(t, matrix)
	assert.Equal(t, []tools.Target{}, matrix)

	jsonData, err := json.Marshal(matrix)
	require.NoError(t, err)
	assert.JSONEq(t, "[]", string(jsonData))
}
