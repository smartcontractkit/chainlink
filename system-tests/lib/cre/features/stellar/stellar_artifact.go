package stellar

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

const (
	envStellarSourceDir = "STELLAR_CONTRACTS_SOURCE_DIR"
	stellarGoModule     = "github.com/smartcontractkit/chainlink-stellar"
)

// autoSourceWarnOnce ensures the auto-fallback warning is logged at most once per process.
var autoSourceWarnOnce sync.Once

// stellarBuildConfig resolves how to build the given contract WASM, identified by its artifact filename.
func stellarBuildConfig(ctx context.Context, artifactFile string) (helpers.BuildStellarConfig, error) {
	destDir, err := os.MkdirTemp("", "stellar-wasm-")
	if err != nil {
		return helpers.BuildStellarConfig{}, fmt.Errorf("create stellar wasm dest dir: %w", err)
	}
	cfg := helpers.BuildStellarConfig{DestinationDir: destDir, ArtifactFile: artifactFile}

	if src := os.Getenv(envStellarSourceDir); src != "" {
		cfg.LocalBuild = helpers.LocalStellarBuildConfig{BuildLocally: true, SourceDir: src}
		return cfg, nil
	}

	moduleDir, modErr := stellarModuleDir(ctx)
	if modErr != nil {
		return helpers.BuildStellarConfig{}, fmt.Errorf(
			"no Stellar WASM source configured: set %s to a chainlink-stellar checkout (local soroban build), or ensure module %s is available on disk: %w",
			envStellarSourceDir, stellarGoModule, modErr)
	}
	autoSourceWarnOnce.Do(func() {
		framework.L.Warn().
			Str("module", stellarGoModule).
			Str("moduleDir", moduleDir).
			Msgf("%s unset; building Stellar contract WASM from the go-list module Dir — this is fragile on a read-only module cache. Prefer exporting %s to a writable chainlink-stellar checkout.",
				envStellarSourceDir, envStellarSourceDir)
	})
	cfg.LocalBuild = helpers.LocalStellarBuildConfig{
		BuildLocally:   true,
		SourceDir:      moduleDir,
		CargoTargetDir: filepath.Join(destDir, "cargo-target"),
	}
	return cfg, nil
}

// stellarModuleDir resolves the on-disk path of the pinned chainlink-stellar module
func stellarModuleDir(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-f", "{{.Dir}}", stellarGoModule)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go list %s: %w (%s)", stellarGoModule, err, strings.TrimSpace(stderr.String()))
	}
	dir := strings.TrimSpace(stdout.String())
	if dir == "" || dir == "|" {
		return "", fmt.Errorf("go list %s returned empty Dir", stellarGoModule)
	}
	if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err != nil {
		return "", fmt.Errorf("chainlink-stellar module Dir %q missing Cargo.toml: %w", dir, err)
	}
	return dir, nil
}
