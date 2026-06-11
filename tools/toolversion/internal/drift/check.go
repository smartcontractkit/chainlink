// Package drift detects version pins outside the canonical manifests.
package drift

import (
	"bufio"
	"fmt"
	"io/fs"
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
	exs, err := loadExceptions(c.cfg.Root)
	if err != nil {
		return fmt.Errorf("load exceptions: %w", err)
	}

	var errs []string
	if err := c.checkMakefilePins(exs); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkGolangMirror(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkStrayToolVersions(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := c.checkRepoWidePins(exs); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "\n"))
}

// isAllowedException reports whether the given line in relFile is covered by an exception.
// relFile must be a slash-separated path relative to the repo root.
func isAllowedException(exs []exception, relFile, line string) bool {
	for _, ex := range exs {
		if relFile != filepath.ToSlash(ex.File) {
			continue
		}
		if strings.Contains(line, ex.Contains) {
			return true
		}
	}
	return false
}

func (c *Checker) checkMakefilePins(exs []exception) error {
	data, err := os.ReadFile(c.cfg.Makefile)
	if err != nil {
		return fmt.Errorf("read GNUmakefile: %w", err)
	}

	type patMod struct {
		pat *regexp.Regexp
		mod string
	}
	mods := c.resolver.ManagedModules()
	patterns := make([]patMod, len(mods))
	for i, mod := range mods {
		patterns[i] = patMod{
			pat: regexp.MustCompile(regexp.QuoteMeta(mod) + `@(v[0-9]|[0-9a-f]{7,})`),
			mod: mod,
		}
	}

	lines := strings.Split(string(data), "\n")
	var violations []string
	for _, pm := range patterns {
		for i, line := range lines {
			if pm.pat.MatchString(line) && !isAllowedException(exs, "GNUmakefile", line) {
				violations = append(violations, fmt.Sprintf("GNUmakefile:%d: hardcoded pin for %s — use $(TOOL_VERSION) go-install instead:\n  %s", i+1, pm.mod, strings.TrimSpace(line)))
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
	root, err := os.OpenRoot(c.cfg.Root)
	if err != nil {
		return err
	}
	defer root.Close()

	var violations []string
	err = fs.WalkDir(root.FS(), ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDir(d.Name(), rel) {
				return fs.SkipDir
			}
			return nil
		}
		if d.Name() != ".tool-versions" {
			return nil
		}
		abs := filepath.Join(c.cfg.Root, rel)
		if filepath.Clean(abs) == rootTV || isFixturePath(rel) {
			return nil
		}
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

func (c *Checker) checkRepoWidePins(exs []exception) error {
	root, err := os.OpenRoot(c.cfg.Root)
	if err != nil {
		return err
	}
	defer root.Close()

	var violations []string
	err = fs.WalkDir(root.FS(), ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDir(d.Name(), rel) {
				return fs.SkipDir
			}
			return nil
		}
		if shouldSkipScan(rel, c.cfg) {
			return nil
		}
		data, err := root.ReadFile(rel)
		if err != nil {
			return nil
		}
		for i, line := range strings.Split(string(data), "\n") {
			if isCommentLine(line) {
				continue
			}
			if !goInstallPin.MatchString(line) {
				continue
			}
			if isAllowedException(exs, filepath.ToSlash(rel), line) {
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

func skipDir(name, rel string) bool {
	switch name {
	case ".git", "vendor", "node_modules", ".cursor", ".bin":
		return true
	}
	return isFixturePath(rel)
}

func isFixturePath(rel string) bool {
	slash := filepath.ToSlash(rel)
	return strings.Contains(slash, "/temp-repo/") || strings.HasPrefix(slash, "temp-repo/")
}

func isCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

func shouldSkipScan(rel string, cfg paths.Config) bool {
	if isFixturePath(rel) {
		return true
	}
	abs := filepath.Join(cfg.Root, rel)
	clean := filepath.Clean(abs)
	if clean == filepath.Clean(cfg.ToolVersionsFile) || clean == filepath.Clean(cfg.GoToolsFile) {
		return true
	}
	if strings.Contains(filepath.ToSlash(rel), "tools/toolversion") {
		return true
	}
	return !scannableFile(filepath.Base(rel), rel)
}

func scannableFile(base, rel string) bool {
	switch base {
	case "GNUmakefile", "Makefile":
		return true
	}
	// Check Dockerfile variants before the extension switch so that bare
	// "Dockerfile" (ext == "") is not short-circuited by the empty-ext false case.
	if strings.HasSuffix(base, "Dockerfile") {
		return true
	}
	ext := filepath.Ext(rel)
	switch ext {
	case ".go", ".sh", ".yml", ".yaml", ".md", ".nix", ".mk", ".Dockerfile", ".dockerfile":
		return true
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".gz", ".zip", ".wasm", ".pb", ".bin", "":
		return false
	}
	return false
}
