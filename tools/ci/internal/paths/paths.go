package paths

import (
	"os"
	"path/filepath"
)

// ResolveFromRepoRoot returns p, or the repo-root-relative variant ("../../"+p) when p does
// not exist relative to the current working directory. It exists because the CLI is invoked
// both from the repository root and from tools/ci (e.g. go -C tools/ci run .). Existing paths
// are returned as absolute paths; the input is returned unchanged when nothing matches.
func ResolveFromRepoRoot(p string) string {
	if p == "" || filepath.IsAbs(p) {
		return p
	}
	if abs := absIfExists(p); abs != "" {
		return abs
	}
	if alt := filepath.Join("../..", p); absIfExists(alt) != "" {
		return absIfExists(alt)
	}
	return p
}

func absIfExists(p string) string {
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}
