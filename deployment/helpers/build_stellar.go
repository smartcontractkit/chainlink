package helpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// EnvStellarSourceDir names a writable chainlink-stellar checkout to compile contract WASM from. When unset, the pinned Go module's own directory is used.
	EnvStellarSourceDir = "STELLAR_CONTRACTS_SOURCE_DIR"
	stellarGoModule     = "github.com/smartcontractkit/chainlink-stellar"
)

// LocalStellarBuildConfig controls building a contract WASM from source with the Soroban toolchain.
type LocalStellarBuildConfig struct {
	BuildLocally bool
	// SourceDir is the chainlink-stellar workspace root that BuildCmd runs in.
	SourceDir string
	// BuildCmd defaults to `stellar contract build` when empty.
	BuildCmd []string
	// BuiltWASMRelPath is the produced artifact path relative to SourceDir; defaults
	// to the standard cargo/soroban release output for the artifact.
	BuiltWASMRelPath string
	// CargoTargetDir, when set, is exported as CARGO_TARGET_DIR so builds work against
	// a read-only Go module cache checkout (compile at runtime).
	CargoTargetDir      string
	CleanDestinationDir bool
}

// BuildStellarConfig sources a compiled contract WASM at deploy time by compiling it
// from chainlink-stellar sources with the Soroban toolchain runtime build.
type BuildStellarConfig struct {
	// DestinationDir is where the resolved WASM is written and read back from.
	DestinationDir string
	// ArtifactFile is the WASM filename to source (e.g. chainlink-stellar deployment/cre.ReadFixtureWasm)
	ArtifactFile string
	LocalBuild   LocalStellarBuildConfig
}

func (c BuildStellarConfig) wasmFile() (string, error) {
	if c.ArtifactFile == "" {
		return "", errors.New("missing artifact file path")
	}
	return c.ArtifactFile, nil
}

// WASMPath is the path the resolved WASM lands at under DestinationDir.
func (c BuildStellarConfig) WASMPath() (string, error) {
	f, err := c.wasmFile()
	if err != nil {
		return "", err
	}
	return filepath.Join(c.DestinationDir, f), nil
}

// BuildStellarConfigFor resolves how to compile the named contract WASM the artifact
// filename, e.g. chainlink-stellar deployment/cre.ForwarderWasm, re-exported by deployment/cre/stellar.
func BuildStellarConfigFor(ctx context.Context, artifactFile string) (BuildStellarConfig, error) {
	destDir, err := os.MkdirTemp("", "stellar-wasm-")
	if err != nil {
		return BuildStellarConfig{}, fmt.Errorf("create stellar wasm dest dir: %w", err)
	}
	cfg := BuildStellarConfig{DestinationDir: destDir, ArtifactFile: artifactFile}

	// An explicit checkout wins. It is writable, so cargo can use its default target dir.
	if src := os.Getenv(EnvStellarSourceDir); src != "" {
		cfg.LocalBuild = LocalStellarBuildConfig{BuildLocally: true, SourceDir: src}
		return cfg, nil
	}

	srcDir, err := stellarModuleDir(ctx)
	if err != nil {
		return BuildStellarConfig{}, fmt.Errorf(
			"no Stellar contract sources available: set %s to a chainlink-stellar checkout, or ensure module %s is on disk: %w",
			EnvStellarSourceDir, stellarGoModule, err)
	}
	cfg.LocalBuild = LocalStellarBuildConfig{
		BuildLocally: true,
		SourceDir:    srcDir,
		// The module cache is read-only, so cargo cannot write target/ inside it.
		CargoTargetDir: filepath.Join(destDir, "cargo-target"),
	}
	return cfg, nil
}

// stellarModuleDir returns the on-disk directory of the pinned chainlink-stellar module.
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
		return "", fmt.Errorf("go list %s returned an empty Dir", stellarGoModule)
	}
	// A module present in the cache but not extracted has no Cargo.toml, and
	// `stellar contract build` would fail deep inside cargo with a far less obvious error.
	if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err != nil {
		return "", fmt.Errorf("chainlink-stellar module dir %q has no Cargo.toml: %w", dir, err)
	}
	return dir, nil
}

// BuildStellar builds the contract WASM from chainlink-stellar sources and returns its bytes.
func BuildStellar(ctx context.Context, config BuildStellarConfig) ([]byte, error) {
	if config.DestinationDir == "" {
		return nil, errors.New("buildStellar: DestinationDir is required")
	}
	if !config.LocalBuild.BuildLocally {
		return nil, errors.New("buildStellar: LocalBuild.BuildLocally must be set — the contract WASM is always built from source (prebuilt-artifact download is not supported)")
	}
	if config.LocalBuild.CleanDestinationDir {
		if err := os.RemoveAll(config.DestinationDir); err != nil {
			return nil, fmt.Errorf("buildStellar: clean destination dir: %w", err)
		}
	}
	if err := os.MkdirAll(config.DestinationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("buildStellar: create destination dir: %w", err)
	}

	if err := buildStellarLocally(ctx, config); err != nil {
		return nil, err
	}

	wasmPath, err := config.WASMPath()
	if err != nil {
		return nil, fmt.Errorf("buildStellar: %w", err)
	}
	wasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("buildStellar: read built WASM %s: %w", wasmPath, err)
	}
	return wasm, nil
}

func buildStellarLocally(ctx context.Context, config BuildStellarConfig) error {
	if config.LocalBuild.SourceDir == "" {
		return errors.New("buildStellar: LocalBuild.SourceDir is required for a local build")
	}

	cmdArgs := config.LocalBuild.BuildCmd
	if len(cmdArgs) == 0 {
		cmdArgs = []string{"stellar", "contract", "build"}
	}
	//nolint:gosec // command comes from deploy config, not untrusted input
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = config.LocalBuild.SourceDir
	if config.LocalBuild.CargoTargetDir != "" {
		if err := os.MkdirAll(config.LocalBuild.CargoTargetDir, 0o755); err != nil {
			return fmt.Errorf("buildStellar: create CARGO_TARGET_DIR: %w", err)
		}
		cmd.Env = append(os.Environ(), "CARGO_TARGET_DIR="+config.LocalBuild.CargoTargetDir)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("buildStellar: %v in %s failed: %w\n%s", cmdArgs, config.LocalBuild.SourceDir, err, out)
	}

	wasmFileName, err := config.wasmFile()
	if err != nil {
		return fmt.Errorf("buildStellar: %w", err)
	}
	var built string
	if config.LocalBuild.CargoTargetDir != "" {
		built = filepath.Join(config.LocalBuild.CargoTargetDir, "wasm32v1-none", "release", wasmFileName)
	} else {
		relPath := config.LocalBuild.BuiltWASMRelPath
		if relPath == "" {
			relPath = filepath.Join("target", "wasm32v1-none", "release", wasmFileName)
		}
		built = filepath.Join(config.LocalBuild.SourceDir, relPath)
	}
	if err = copyFile(built, config.DestinationDir); err != nil {
		return fmt.Errorf("buildStellar: copy built WASM %s: %w", built, err)
	}
	return nil
}
