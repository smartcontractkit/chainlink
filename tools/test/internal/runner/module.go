package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// resolveModuleDir returns the Go module root and adjusted goTestArgs for the
// package patterns in goTestArgs. When relative patterns (./foo/...) point into
// a subdirectory that owns its own go.mod, the working directory is moved to
// that submodule root and the patterns are rewritten relative to it.
// An error is returned when patterns span more than one module.
func resolveModuleDir(repoRoot string, goTestArgs []string) (string, []string, error) {
	patterns := packagePatternsFromEnd(goTestArgs)

	var relPatterns []string
	for _, p := range patterns {
		if isRelativePackagePattern(p) {
			relPatterns = append(relPatterns, p)
		}
	}

	if len(relPatterns) == 0 {
		return repoRoot, goTestArgs, nil
	}

	var moduleRoot string
	for _, p := range relPatterns {
		dir := patternBaseDir(p)
		abs := filepath.Join(repoRoot, dir)
		mod := nearestModuleRoot(abs, repoRoot)
		if moduleRoot == "" {
			moduleRoot = mod
		} else if moduleRoot != mod {
			return "", nil, fmt.Errorf("package patterns span multiple Go modules (%s and %s): run go test for each module separately", moduleRoot, mod)
		}
	}

	if moduleRoot == repoRoot {
		return repoRoot, goTestArgs, nil
	}

	adjusted := rewriteRelativePatterns(repoRoot, moduleRoot, goTestArgs)
	return moduleRoot, adjusted, nil
}

func isRelativePackagePattern(p string) bool {
	return strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") || p == "." || p == ".."
}

// patternBaseDir strips trailing /... and /. wildcards to get the directory part.
func patternBaseDir(p string) string {
	if s, ok := strings.CutSuffix(p, "/..."); ok {
		return s
	}
	if s, ok := strings.CutSuffix(p, "/."); ok {
		return s
	}
	return p
}

// nearestModuleRoot walks up from dir toward stopAt looking for a go.mod.
// Returns stopAt if no intermediate go.mod is found.
func nearestModuleRoot(dir, stopAt string) string {
	d := filepath.Clean(dir)
	stop := filepath.Clean(stopAt)
	for d != stop {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return stop
		}
		d = parent
	}
	return stop
}

// rewriteRelativePatterns rewrites each relative package pattern in goTestArgs
// so it is expressed relative to moduleRoot instead of repoRoot.
func rewriteRelativePatterns(repoRoot, moduleRoot string, goTestArgs []string) []string {
	result := make([]string, len(goTestArgs))
	for i, arg := range goTestArgs {
		if !isRelativePackagePattern(arg) || !looksLikeGoPackagePattern(arg) {
			result[i] = arg
			continue
		}
		base := patternBaseDir(arg)
		suffix := arg[len(base):]
		abs := filepath.Join(repoRoot, base)
		rel, err := filepath.Rel(moduleRoot, abs)
		if err != nil {
			result[i] = arg
			continue
		}
		if rel == "." {
			result[i] = "." + suffix
		} else {
			result[i] = "./" + rel + suffix
		}
	}
	return result
}
