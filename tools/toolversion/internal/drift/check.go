// Package drift detects version pins outside the canonical manifests.
package drift

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/paths"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/resolve"
)

// Checker validates that consumers do not hardcode managed tool versions.
type Checker struct {
	cfg      paths.Config
	resolver *resolve.Resolver
}

func NewChecker(cfg paths.Config, resolver *resolve.Resolver) *Checker {
	return &Checker{cfg: cfg, resolver: resolver}
}

// Check runs all drift rules and returns a combined error if any fail.
func (c *Checker) Check() error {
	var errs []string
	if err := c.checkMakefilePins(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkGolangMirror(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkStrayToolVersions(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkRepoWidePins(); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "\n"))
}

func (c *Checker) checkMakefilePins() error {
	data, err := os.ReadFile(c.cfg.Makefile)
	if err != nil {
		return fmt.Errorf("read GNUmakefile: %w", err)
	}
	content := string(data)
	var violations []string
	for _, mod := range c.resolver.ManagedModules() {
		pat := regexp.MustCompile(regexp.QuoteMeta(mod) + `@(v[0-9]|[0-9a-f]{7,})`)
		for i, line := range strings.Split(content, "\n") {
			if pat.MatchString(line) && !isAllowedException(c.cfg.Root, "GNUmakefile", line) {
				violations = append(violations, fmt.Sprintf("GNUmakefile:%d: hardcoded pin for %s — use $(TOOL_VERSION) go-install instead:\n  %s", i+1, mod, strings.TrimSpace(line)))
			}
		}
	}
	if len(violations) > 0 {
		return fmt.Errorf("check-tool-versions:\n%s", strings.Join(violations, "\n"))
	}
	return nil
}

func (c *Checker) checkGolangMirror() error {
	tvGo, err := c.resolver.Get("golang")
	if err != nil {
		return err
	}
	modGo, err := parseGoModDirective(c.cfg.GoMod)
	if err != nil {
		return err
	}
	if tvGo != modGo {
		return fmt.Errorf("check-tool-versions: golang in .tool-versions (%s) != go directive in go.mod (%s); align .tool-versions to go.mod", tvGo, modGo)
	}
	return nil
}

func parseGoModDirective(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("no go directive in %s", path)
}

func (c *Checker) checkStrayToolVersions() error {
	rootTV := filepath.Clean(c.cfg.ToolVersionsFile)
	var violations []string
	err := filepath.WalkDir(c.cfg.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDir(d.Name(), path) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != ".tool-versions" {
			return nil
		}
		if filepath.Clean(path) == rootTV || isFixturePath(path) {
			return nil
		}
		rel, _ := filepath.Rel(c.cfg.Root, path)
		violations = append(violations, fmt.Sprintf("check-tool-versions: stray .tool-versions at %s — only repo root .tool-versions is allowed", rel))
		return nil
	})
	if err != nil {
		return err
	}
	if len(violations) > 0 {
		return fmt.Errorf("%s", strings.Join(violations, "\n"))
	}
	return nil
}

var goInstallPin = regexp.MustCompile(`go install [^[:space:]]+@(v[0-9]+(\.[0-9]+)*|latest|[0-9a-f]{7,})`)

func (c *Checker) checkRepoWidePins() error {
	var violations []string

	err := filepath.WalkDir(c.cfg.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDir(d.Name(), path) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldSkipScan(path, c.cfg) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(c.cfg.Root, path)
		for i, line := range strings.Split(string(data), "\n") {
			if isCommentLine(line) {
				continue
			}
			if !goInstallPin.MatchString(line) {
				continue
			}
			if isAllowedException(c.cfg.Root, rel, line) {
				continue
			}
			violations = append(violations, fmt.Sprintf("%s:%d: hardcoded go install pin — use go run ./tools/toolversion go-install <key>:\n  %s", rel, i+1, strings.TrimSpace(line)))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(violations) > 0 {
		return fmt.Errorf("check-tool-versions:\n%s", strings.Join(violations, "\n"))
	}
	return nil
}

func skipDir(name, path string) bool {
	switch name {
	case ".git", "vendor", "node_modules", ".cursor", ".bin":
		return true
	}
	return isFixturePath(path)
}

func isFixturePath(path string) bool {
	return strings.Contains(filepath.ToSlash(path), "/temp-repo/")
}

func isCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

func shouldSkipScan(path string, cfg paths.Config) bool {
	if isFixturePath(path) {
		return true
	}
	clean := filepath.Clean(path)
	if clean == filepath.Clean(cfg.ToolVersionsFile) || clean == filepath.Clean(cfg.GoToolsFile) {
		return true
	}
	if strings.Contains(path, string(filepath.Join("tools", "toolversion"))) {
		return true
	}
	base := filepath.Base(path)
	if !scannableFile(base, path) {
		return true
	}
	return false
}

func scannableFile(base, path string) bool {
	switch base {
	case "GNUmakefile", "Makefile":
		return true
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".sh", ".yml", ".yaml", ".md", ".nix", ".mk", ".Dockerfile", ".dockerfile":
		return true
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".gz", ".zip", ".wasm", ".pb", ".bin", "":
		return false
	}
	if strings.HasSuffix(base, "Dockerfile") {
		return true
	}
	return false
}
