package drift

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/manifest"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/paths"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/resolve"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func setupRepo(t *testing.T) (paths.Config, *resolve.Resolver) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, ".tool-versions", "golang 1.26.4\nmockery 2.53.0\n")
	writeFile(t, root, "go.mod", "module example.com/test\n\ngo 1.26.4\n")
	if err := os.Mkdir(filepath.Join(root, "tools"), 0o755); err != nil {
		t.Fatal(err)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	return cfg, resolve.New(store)
}

func TestCheckMakefilePinViolation(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "GNUmakefile", "install:\n\tgo install github.com/vektra/mockery/v2@v2.99.0\n")

	err := NewChecker(cfg, resolver).Check()
	if err == nil {
		t.Fatal("expected violation")
	}
}

func TestCheckGolangMirror(t *testing.T) {
	t.Parallel()
	cfg, _ := setupRepo(t)
	writeFile(t, cfg.Root, ".tool-versions", "golang 1.26.3\nmockery 2.53.0\n")
	store, err := manifest.New(cfg.ToolVersionsFile, cfg.GoToolsFile)
	if err != nil {
		t.Fatal(err)
	}
	resolver := resolve.New(store)

	err = NewChecker(cfg, resolver).Check()
	if err == nil {
		t.Fatal("expected golang mirror violation")
	}
}

func TestCheckStrayToolVersions(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	if err := os.Mkdir(filepath.Join(cfg.Root, "integration-tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(cfg.Root, "integration-tests"), ".tool-versions", "golang 1.26.4\n")

	err := NewChecker(cfg, resolver).Check()
	if err == nil {
		t.Fatal("expected stray .tool-versions violation")
	}
}

func TestCheckAllowedProtocException(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, filepath.Join(cfg.Root, "tools"), "version-exceptions.yaml", "- file: GNUmakefile\n  contains: \"protoc-gen-go@\"\n  reason: module-coupled\n")
	writeFile(t, cfg.Root, "GNUmakefile", "install:\n\tgo install google.golang.org/protobuf/cmd/protoc-gen-go@`go list -m -json google.golang.org/protobuf | jq -r .Version`\n")

	err := NewChecker(cfg, resolver).checkMakefilePins()
	if err != nil {
		t.Fatalf("expected protoc exception to pass makefile check: %v", err)
	}
}

func TestCheckRepoWideGoInstall(t *testing.T) {
	t.Parallel()
	cfg, resolver := setupRepo(t)
	writeFile(t, cfg.Root, "script.sh", "go install gotest.tools/gotestsum@v9.9.9\n")

	err := NewChecker(cfg, resolver).Check()
	if err == nil {
		t.Fatal("expected repo-wide go install violation")
	}
}
