package drift

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/manifest"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/paths"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/resolve"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func setupRepo(t *testing.T) (paths.Config, *resolve.Resolver) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, ".tool-versions", "golang 1.26.4\nmockery 2.53.0\n")
	writeFile(t, root, "go.mod", "module example.com/test\n\ngo 1.26.4\n")
	require.NoError(t, os.Mkdir(filepath.Join(root, "tools"), 0o755))
	writeFile(t, filepath.Join(root, "tools"), "go-tools.txt", "github.com/jmank88/gomods 0.1.7\n")
	writeFile(t, root, "GNUmakefile", "install:\n\t$(TOOL_VERSION) go-install mockery\n")

	cfg := paths.Config{
		Root:             root,
		ToolVersionsFile: filepath.Join(root, ".tool-versions"),
		GoToolsFile:      filepath.Join(root, "tools", "go-tools.txt"),
		Makefile:         filepath.Join(root, "GNUmakefile"),
		GoMod:            filepath.Join(root, "go.mod"),
	}
	store, err := manifest.New(cfg.ToolVersionsFile, cfg.GoToolsFile)
	require.NoError(t, err)
	return cfg, resolve.New(store)
}

func TestCheckAllClean(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	require.NoError(t, NewChecker(cfg, resolver).Check())
}

func TestCheckMakefilePinViolation(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "GNUmakefile", "install:\n\tgo install github.com/vektra/mockery/v2@v2.99.0\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err)
}

func TestCheckGolangMirror(t *testing.T) {
	t.Parallel()
	cfg, _ := setupRepo(t)
	writeFile(t, cfg.Root, ".tool-versions", "golang 1.26.3\nmockery 2.53.0\n")
	store, err := manifest.New(cfg.ToolVersionsFile, cfg.GoToolsFile)
	require.NoError(t, err)
	resolver := resolve.New(store)

	err = NewChecker(cfg, resolver).Check()
	require.Error(t, err)
}

func TestCheckStrayToolVersions(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	require.NoError(t, os.Mkdir(filepath.Join(cfg.Root, "integration-tests"), 0o755))
	writeFile(t, filepath.Join(cfg.Root, "integration-tests"), ".tool-versions", "golang 1.26.4\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err)
}

func TestCheckAllowedProtocException(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, filepath.Join(cfg.Root, "tools"), "version-exceptions.yaml",
		"- file: GNUmakefile\n  contains: \"protoc-gen-go@\"\n  reason: module-coupled\n")
	writeFile(t, cfg.Root, "GNUmakefile", "install:\n\tgo install google.golang.org/protobuf/cmd/protoc-gen-go@`go list -m -json google.golang.org/protobuf | jq -r .Version`\n")

	exs, err := loadExceptions(cfg.Root)
	require.NoError(t, err)
	err = NewChecker(cfg, resolver).checkMakefilePins(exs)
	require.NoError(t, err)
}

// TestCheckMakefilePinWithException verifies that a managed-module pin in GNUmakefile
// is suppressed when covered by an exception entry.
func TestCheckMakefilePinWithException(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, filepath.Join(cfg.Root, "tools"), "version-exceptions.yaml",
		"- file: GNUmakefile\n  contains: \"github.com/vektra/mockery/v2@\"\n  reason: bootstrap\n")
	writeFile(t, cfg.Root, "GNUmakefile", "install:\n\tgo install github.com/vektra/mockery/v2@v2.99.0\n")

	err := NewChecker(cfg, resolver).Check()
	require.NoError(t, err)
}

func TestCheckRepoWideGoInstall(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "script.sh", "go install gotest.tools/gotestsum@v9.9.9\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err)
}

func TestCheckRepoWideGoInstallClean(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "script.sh", "#!/bin/bash\ngo run ./tools/toolversion go-install somelib\n")

	err := NewChecker(cfg, resolver).Check()
	require.NoError(t, err)
}

func TestCheckRepoWideGoInstallInDockerfile(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "Dockerfile", "FROM golang:1.26\nRUN go install gotest.tools/gotestsum@v9.9.9\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err, "bare Dockerfile with hardcoded go install pin should be flagged")
}

func TestCheckMalformedExceptions(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, filepath.Join(cfg.Root, "tools"), "version-exceptions.yaml", "not: valid: yaml: [\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err, "malformed version-exceptions.yaml must return an error, not silently ignore exceptions")
}

func TestExceptionExactPathNoSubdirLeak(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	// Exception covers only "script.sh" at repo root, NOT "sub/script.sh"
	writeFile(t, filepath.Join(cfg.Root, "tools"), "version-exceptions.yaml",
		"- file: script.sh\n  contains: \"go install gotest.tools\"\n  reason: test\n")
	require.NoError(t, os.Mkdir(filepath.Join(cfg.Root, "sub"), 0o755))
	writeFile(t, filepath.Join(cfg.Root, "sub"), "script.sh", "go install gotest.tools/gotestsum@v9.9.9\n")

	err := NewChecker(cfg, resolver).Check()
	require.Error(t, err, "exception for 'script.sh' must not cover 'sub/script.sh'")
}

func TestScannableFileDockerfile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		base string
		rel  string
		want bool
	}{
		{"Dockerfile", "Dockerfile", true},
		{"Dockerfile", "some/path/Dockerfile", true},
		{"chainlink.Dockerfile", "core/chainlink.Dockerfile", true},
		{"server.dockerfile", "server.dockerfile", true},
		{"Makefile", "Makefile", true},
		{"GNUmakefile", "GNUmakefile", true},
		{"main.go", "main.go", true},
		{"config.yaml", "config.yaml", true},
		{"image.png", "image.png", false},
		{"binary", "binary", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, scannableFile(tt.base, tt.rel), "scannableFile(%q, %q)", tt.base, tt.rel)
	}
}
