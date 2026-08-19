package generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

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
		cmd := exec.CommandContext(ctx, "mockery", args[1:]...) // #nosec G204 -- args come from the local hook pipeline, not remote input
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

		baseName := filepath.Base(cleanRel)
		fileDir := filepath.Dir(absFile)

		// Check if go.mod / go.sum changed -> trigger go.md modgraph generation
		if baseName == "go.mod" || baseName == "go.sum" {
			runModGraph = true
		}

		// Check if config docs generator should run
		if strings.HasPrefix(cleanRel, "core/config") {
			runConfigDocs = true
		}

		// Check proto or generate files
		if strings.HasSuffix(cleanRel, ".proto") ||
			baseName == "generate.go" ||
			baseName == "gen.go" {
			modDir, targetModErr := modules.FindModuleDir(absRoot, absFile)
			if targetModErr != nil {
				return fmt.Errorf("failed to find Go module for %s: %w", file, targetModErr)
			}
			if modDir == "" {
				modDir = absRoot
			}

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
		if filepath.Ext(file) != ".go" && baseName != "go.mod" && baseName != "go.sum" {
			continue
		}

		modDir, modErr := modules.FindModuleDir(absRoot, absFile)
		if modErr != nil {
			return fmt.Errorf("failed to find Go module for %s: %w", file, modErr)
		}
		if modDir == "" {
			continue
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

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, t := range targets {
		wg.Go(func() {
			if err := runner(ctx, t.modDir, "generate", t.pattern); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("go generate %s: %w", t.pattern, err))
				mu.Unlock()
			}
		})
	}

	if runConfigDocs {
		wg.Go(func() {
			if err := runner(ctx, absRoot, "run", "./core/config/docs/cmd/generate", "-o", "./docs/"); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("generate config docs: %w", err))
				mu.Unlock()
			}
		})
	}

	if runModGraph {
		wg.Go(func() {
			if err := runner(ctx, absRoot, "modgraph"); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("generate go.md: %w", err))
				mu.Unlock()
			}
		})
	}

	// Mockery runs: scoped per-package configs for changed packages, full runs
	// for config-file or dependency changes matching CI's `make generate`.
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
				dir := dir
				wg.Go(func() {
					if err := runner(ctx, dir, "mockery"); err != nil {
						mu.Lock()
						errs = append(errs, fmt.Errorf("mockery in %s: %w", dir, err))
						mu.Unlock()
					}
				})
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

			runDir := dir
			wg.Go(func() {
				if err := runner(ctx, runDir, "mockery", "--config", tmpConfig); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("mockery in %s: %w", runDir, err))
					mu.Unlock()
				}
			})
		}
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
