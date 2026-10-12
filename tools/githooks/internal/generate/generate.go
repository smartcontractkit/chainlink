package generate

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

const mockeryConfigName = ".mockery.yaml"

// CommandRunner executes a command in a directory.
type CommandRunner func(ctx context.Context, dir string, args ...string) error

// Config holds options for running generate.
type Config struct {
	Runner CommandRunner
}

func defaultRunner(ctx context.Context, dir string, args ...string) error {
	if len(args) > 0 && args[0] == "modgraph" {
		modgraphScript := filepath.Join(dir, "tools/bin/modgraph")
		cmd := exec.CommandContext(ctx, modgraphScript)
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("modgraph failed: %w", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "go.md"), out, 0o644); err != nil { // #nosec G306 -- must match the tracked go.md permissions (0644)
			return fmt.Errorf("failed to write go.md: %w", err)
		}
		return nil
	}

	if len(args) > 0 && args[0] == "mockery" {
		cmd := exec.CommandContext(ctx, resolveMockeryBin(ctx), args[1:]...) //nolint:gosec // args come from the local hook pipeline, not remote input
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("mockery %v in %s failed: %w (output: %s)", args[1:], dir, err, string(out))
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go %v in %s failed: %w (output: %s)", args, dir, err, string(out))
	}
	return nil
}

// resolveMockeryBin finds the mockery binary the same way GNUmakefile's
// MOCKERY_BIN does, so the hook and `make generate` never run different
// mockery versions. MOCKERY_BIN itself is honored first as an override.
func resolveMockeryBin(ctx context.Context) string {
	if bin := os.Getenv("MOCKERY_BIN"); bin != "" {
		return bin
	}
	if path, err := exec.LookPath("mockery"); err == nil {
		return path
	}
	if out, err := exec.CommandContext(ctx, "go", "env", "GOBIN").Output(); err == nil {
		if gobin := strings.TrimSpace(string(out)); gobin != "" {
			if p := filepath.Join(gobin, "mockery"); isFile(p) {
				return p
			}
		}
	}
	if out, err := exec.CommandContext(ctx, "go", "env", "GOPATH").Output(); err == nil {
		if gopath := strings.TrimSpace(string(out)); gopath != "" {
			if p := filepath.Join(gopath, "bin", "mockery"); isFile(p) {
				return p
			}
		}
	}
	return "mockery"
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// genTarget is a go generate invocation scoped to a module directory.
type genTarget struct {
	modDir  string // absolute module directory "" falls back to repoRoot
	pattern string // package pattern relative to the module root
}

// mockeryConfig is a lazily loaded .mockery.yaml plus the packages it covers.
type mockeryConfig struct {
	raw  []byte
	pkgs map[string]bool
}

// mockeryTarget tracks pending mockery work for a config directory.
type mockeryTarget struct {
	full     bool
	packages map[string]struct{}
}

func (t *mockeryTarget) addPackage(pkg string) {
	if t.packages == nil {
		t.packages = make(map[string]struct{})
	}
	t.packages[pkg] = struct{}{}
}

func (c *mockeryConfig) hasPackage(pkg string) bool {
	return c.pkgs[pkg]
}

func (c *mockeryConfig) hasModulePrefix(modulePath string) bool {
	prefix := modulePath + "/"
	for pkg := range c.pkgs {
		if pkg == modulePath || strings.HasPrefix(pkg, prefix) {
			return true
		}
	}
	return false
}

// loadMockeryConfig reads and parses the .mockery.yaml in dir. Returns (nil, nil)
// when the file does not exist.
func loadMockeryConfig(dir string) (*mockeryConfig, error) {
	path := filepath.Join(dir, mockeryConfigName)
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	pkgs, err := mockeryConfigPackages(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &mockeryConfig{raw: raw, pkgs: pkgs}, nil
}

// mockeryConfigPackages extracts the package import paths from a packages-mode config.
func mockeryConfigPackages(raw []byte) (map[string]bool, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("invalid mockery config: %w", err)
	}

	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, nil
	}

	root := doc.Content[0]
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != "packages" {
			continue
		}

		pkgMap := root.Content[i+1]
		if pkgMap.Kind != yaml.MappingNode {
			return nil, nil
		}

		pkgs := make(map[string]bool, len(pkgMap.Content)/2)
		for j := 0; j+1 < len(pkgMap.Content); j += 2 {
			pkgs[pkgMap.Content[j].Value] = true
		}
		return pkgs, nil
	}

	return nil, nil
}

// scopedMockeryConfig rewrites a packages-mode config so it only contains the
// given package import paths, preserving all global defaults and per-package overrides.
func scopedMockeryConfig(raw []byte, pkgs ...string) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("invalid mockery config: %w", err)
	}

	if len(doc.Content) == 0 {
		return nil, errors.New("mockery config is empty")
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, errors.New("mockery config is not a mapping")
	}

	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != "packages" {
			continue
		}

		wanted := make(map[string]bool, len(pkgs))
		for _, pkg := range pkgs {
			wanted[pkg] = true
		}

		pkgMap := root.Content[i+1]
		if pkgMap.Kind == yaml.MappingNode {
			kept := pkgMap.Content[:0]
			for j := 0; j+1 < len(pkgMap.Content); j += 2 {
				if wanted[pkgMap.Content[j].Value] {
					kept = append(kept, pkgMap.Content[j], pkgMap.Content[j+1])
				}
			}
			pkgMap.Content = kept
		}

		out, err := yaml.Marshal(root)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal mockery config: %w", err)
		}
		return out, nil
	}

	return nil, errors.New("mockery config has no packages mapping")
}

// packageImportPath returns the full package import path for a .go file given
// its module directory and module path.
func packageImportPath(modDir, modulePath, fileDir string) string {
	rel, err := filepath.Rel(modDir, fileDir)
	if err != nil || rel == "." {
		return modulePath
	}
	return modulePath + "/" + filepath.ToSlash(rel)
}

// Run maps changed files to specific code generator targets and runs them.
func Run(ctx context.Context, repoRoot string, files []string, cfg ...Config) error {
	if len(files) == 0 {
		return nil
	}

	runner := defaultRunner
	if len(cfg) > 0 && cfg[0].Runner != nil {
		runner = cfg[0].Runner
	}

	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get absolute repo root: %w", err)
	}

	genTargets := make(map[genTarget]struct{})
	runConfigDocs := false
	runModGraph := false

	mockeryTargets := make(map[string]*mockeryTarget)
	configCache := make(map[string]*mockeryConfig)

	loadConfig := func(dir string) (*mockeryConfig, error) {
		if cached, ok := configCache[dir]; ok {
			return cached, nil
		}
		loaded, loadErr := loadMockeryConfig(dir)
		if loadErr != nil {
			return nil, loadErr
		}
		configCache[dir] = loaded
		return loaded, nil
	}

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		baseName := filepath.Base(file)
		ext := filepath.Ext(file)

		// Fast pre-filter: skip files that cannot trigger generators or mockery
		if ext != ".go" && ext != ".proto" && ext != ".yaml" && ext != ".yml" &&
			baseName != "go.mod" && baseName != "go.sum" && baseName != "TAG" &&
			!strings.Contains(filepath.ToSlash(file), "core/config") {
			continue
		}

		absFile := file
		if !filepath.IsAbs(absFile) {
			absFile = filepath.Join(absRoot, file)
		}

		relFromRoot, relErr := filepath.Rel(absRoot, absFile)
		if relErr != nil {
			return fmt.Errorf("failed to get relative path from repo root: %w", relErr)
		}
		cleanRel := filepath.ToSlash(filepath.Clean(relFromRoot))
		cleanRel = strings.TrimPrefix(cleanRel, "./")

		baseName = filepath.Base(cleanRel)
		fileDir := filepath.Dir(absFile)

		// Check if go.mod / go.sum changed -> trigger go.md modgraph generation
		if baseName == "go.mod" || baseName == "go.sum" {
			runModGraph = true
		}

		// Check if config docs generator should run
		if strings.HasPrefix(cleanRel, "core/config") {
			runConfigDocs = true
		}

		// Check if operator UI tag changed -> trigger core/web asset generation
		if cleanRel == "operator_ui/TAG" {
			genTargets[genTarget{modDir: absRoot, pattern: "./core/web"}] = struct{}{}
		}

		// Check proto or generate files
		var (
			modDir       string
			modResolved  bool
			targetModErr error
		)

		isGoFile := filepath.Ext(absFile) == ".go"
		isProto := strings.HasSuffix(cleanRel, ".proto")
		isGenFile := baseName == "generate.go" || baseName == "gen.go"
		isDirectiveGo := isGoFile && hasGoGenerateDirective(absFile)

		if isProto || isGenFile || isDirectiveGo {
			modDir, targetModErr = modules.FindModuleDir(absRoot, absFile)
			if targetModErr != nil {
				return fmt.Errorf("failed to find Go module for %s: %w", file, targetModErr)
			}
			if modDir == "" {
				modDir = absRoot
			}
			modResolved = true

			relPkg, relPkgErr := filepath.Rel(modDir, filepath.Dir(absFile))
			if relPkgErr != nil {
				return fmt.Errorf("failed to get package path relative to module: %w", relPkgErr)
			}

			pkgPattern := "."
			if relPkg != "." {
				pkgPattern = "./" + filepath.ToSlash(filepath.Clean(relPkg))
			}
			genTargets[genTarget{modDir: modDir, pattern: pkgPattern}] = struct{}{}
		}

		// Check mockery config change -> full mockery run for that config
		if baseName == mockeryConfigName {
			target := mockeryTargets[fileDir]
			if target == nil {
				target = &mockeryTarget{}
				mockeryTargets[fileDir] = target
			}
			target.full = true
			continue
		}

		// Check Go file or dependency changes for mockery coverage
		if !isGoFile && baseName != "go.mod" && baseName != "go.sum" {
			continue
		}

		if !modResolved {
			modDir, targetModErr = modules.FindModuleDir(absRoot, absFile)
			if targetModErr != nil {
				return fmt.Errorf("failed to find Go module for %s: %w", file, targetModErr)
			}
			if modDir == "" {
				continue
			}
		}

		modulePath, pathErr := modules.ReadModulePath(modDir)
		if pathErr != nil {
			return fmt.Errorf("failed to resolve module path for %s: %w", file, pathErr)
		}

		isDepChange := baseName == "go.mod" || baseName == "go.sum"
		pkgImportPath := packageImportPath(modDir, modulePath, fileDir)

		for dir := fileDir; ; dir = filepath.Dir(dir) {
			cfgPtr, cfgErr := loadConfig(dir)
			if cfgErr != nil {
				return cfgErr
			}
			if cfgPtr == nil {
				if dir == absRoot || dir == filepath.Dir(dir) {
					break
				}
				continue
			}

			var covered bool
			if isDepChange {
				covered = cfgPtr.hasModulePrefix(modulePath)
			} else {
				covered = cfgPtr.hasPackage(pkgImportPath)
			}
			if covered {
				target := mockeryTargets[dir]
				if target == nil {
					target = &mockeryTarget{}
					mockeryTargets[dir] = target
				}
				if isDepChange {
					target.full = true
				} else {
					target.addPackage(pkgImportPath)
				}
			}

			if dir == absRoot || dir == filepath.Dir(dir) {
				break
			}
		}
	}

	targets := make([]genTarget, 0, len(genTargets))
	for t := range genTargets {
		targets = append(targets, t)
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].modDir != targets[j].modDir {
			return targets[i].modDir < targets[j].modDir
		}
		return targets[i].pattern < targets[j].pattern
	})

	if len(targets) == 0 && !runConfigDocs && !runModGraph && len(mockeryTargets) == 0 {
		return nil
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.GOMAXPROCS(0))

	for _, t := range targets {
		g.Go(func() error {
			if err := runner(gCtx, t.modDir, "generate", t.pattern); err != nil {
				return fmt.Errorf("go generate %s: %w", t.pattern, err)
			}
			return nil
		})
	}

	if runConfigDocs {
		g.Go(func() error {
			if err := runner(gCtx, absRoot, "run", "./core/config/docs/cmd/generate", "-o", "./docs/"); err != nil {
				return fmt.Errorf("generate config docs: %w", err)
			}
			return nil
		})
	}

	if runModGraph {
		g.Go(func() error {
			if err := runner(gCtx, absRoot, "modgraph"); err != nil {
				return fmt.Errorf("generate go.md: %w", err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// Mockery runs: scoped per-package configs for changed packages, full runs
	// for config-file or dependency changes matching CI's `make generate`.
	// Mockery internally uses all CPU cores, so configs are run sequentially to
	// avoid thread oversubscription and memory contention.
	if len(mockeryTargets) > 0 {
		mockeryDirs := make([]string, 0, len(mockeryTargets))
		for dir := range mockeryTargets {
			mockeryDirs = append(mockeryDirs, dir)
		}
		sort.Strings(mockeryDirs)

		needsScoped := false
		for _, dir := range mockeryDirs {
			target := mockeryTargets[dir]
			if target != nil && !target.full && len(target.packages) > 0 {
				needsScoped = true
				break
			}
		}

		var tmpDir string
		if needsScoped {
			var tmpErr error
			if tmpDir, tmpErr = os.MkdirTemp("", "githooks-mockery-"); tmpErr != nil {
				return fmt.Errorf("failed to create mockery config temp dir: %w", tmpErr)
			}
			defer os.RemoveAll(tmpDir)
		}

		for _, dir := range mockeryDirs {
			target := mockeryTargets[dir]
			if target == nil || (!target.full && len(target.packages) == 0) {
				continue
			}

			if target.full {
				if err := runner(ctx, dir, "mockery"); err != nil {
					return fmt.Errorf("mockery in %s: %w", dir, err)
				}
				continue
			}

			pkgs := make([]string, 0, len(target.packages))
			for pkg := range target.packages {
				pkgs = append(pkgs, pkg)
			}
			sort.Strings(pkgs)

			baseCfg := configCache[dir]
			if baseCfg == nil {
				var loadErr error
				if baseCfg, loadErr = loadMockeryConfig(dir); loadErr != nil {
					return loadErr
				}
			}
			if baseCfg == nil {
				continue
			}

			scoped, scopedErr := scopedMockeryConfig(baseCfg.raw, pkgs...)
			if scopedErr != nil {
				return scopedErr
			}

			tmpFile, tmpFileErr := os.CreateTemp(tmpDir, filepath.Base(dir)+"-mockery-*.yaml")
			if tmpFileErr != nil {
				return fmt.Errorf("failed to create scoped mockery config temp file: %w", tmpFileErr)
			}
			tmpConfig := tmpFile.Name()
			if _, writeErr := tmpFile.Write(scoped); writeErr != nil {
				_ = tmpFile.Close()
				return fmt.Errorf("failed to write scoped mockery config: %w", writeErr)
			}
			if closeErr := tmpFile.Close(); closeErr != nil {
				return fmt.Errorf("failed to close scoped mockery config: %w", closeErr)
			}

			if err := runner(ctx, dir, "mockery", "--config", tmpConfig); err != nil {
				return fmt.Errorf("mockery in %s: %w", dir, err)
			}
		}
	}

	return nil
}

// hasGoGenerateDirective checks if a Go source file contains a //go:generate directive.
func hasGoGenerateDirective(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), []byte("//go:generate")) {
			return true
		}
	}
	if scanner.Err() != nil {
		// Fail open: assume a directive may be present rather than silently
		// skipping generation for a file we couldn't fully scan.
		return true
	}
	return false
}
