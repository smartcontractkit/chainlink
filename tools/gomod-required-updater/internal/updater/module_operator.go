package updater

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const (
	majorVersionPattern = `/v\d+$`
)

// getMajorVersion extracts the major version number from a module path
// Returns "v2" for /v2, "v0" for no version suffix
func getMajorVersion(modulePath string) string {
	re := regexp.MustCompile(majorVersionPattern)
	if match := re.FindString(modulePath); match != "" {
		return "v" + strings.TrimPrefix(match, "/v")
	}
	return "v0"
}

type moduleOperator struct {
	config *Config
}

func NewModuleOperator(config *Config) ModuleOperator {
	if config.RepoRemote == "" {
		config.RepoRemote = "origin"
	}
	if config.BranchTrunk == "" {
		config.BranchTrunk = "develop"
	}
	return &moduleOperator{
		config: config,
	}
}

func (m *moduleOperator) GetGitInfo(remote, branch string) (string, time.Time, error) {
	// Get latest SHA
	cmd := exec.Command("git", "ls-remote", remote, "refs/heads/"+branch)
	out, err := cmd.Output()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to get SHA: %v", ErrModOperation, err)
	}
	sha := strings.Split(string(out), "\t")[0]

	// Get commit timestamp
	cmd = exec.Command("git", "show", "-s", "--format=%cI", sha)
	out, err = cmd.Output()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to get commit time: %v", ErrModOperation, err)
	}
	commitTime, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to parse commit time: %v", ErrModOperation, err)
	}

	return sha[:12], commitTime, nil
}

func (m *moduleOperator) GetLatestVersion(modulePath string) (module.Version, error) {
	sha, commitTime, err := m.GetGitInfo(m.config.RepoRemote, m.config.BranchTrunk)
	if err != nil {
		return module.Version{}, err
	}

	majorVer := strings.TrimPrefix(getMajorVersion(modulePath), "v")
	pseudoVersion := module.PseudoVersion("v"+majorVer, "", commitTime, sha)

	return module.Version{
		Path:    modulePath,
		Version: pseudoVersion,
	}, nil
}

func (m *moduleOperator) validateVersion(modulePath, version string) error {
	expectedMajor := getMajorVersion(modulePath)
	if !strings.HasPrefix(version, expectedMajor) {
		return fmt.Errorf("version %q invalid: should be %s, not v0", version, expectedMajor)
	}
	return nil
}

func (m *moduleOperator) UpdateRequiredVersions(modFile *modfile.File, newVersion string) error {
	for _, req := range modFile.Require {
		if strings.HasPrefix(req.Mod.Path, "github.com/smartcontractkit/chainlink") {
			 // Check if this is a v2+ module path
			modVersion := getMajorVersion(req.Mod.Path)
			if modVersion != "v0" {
				// If module path contains v2+, maintain that version
				newVer := strings.Replace(newVersion, "v0", modVersion, 1)
				req.Mod.Version = newVer
			} else {
				// For modules without version in path, use newVersion as-is
				req.Mod.Version = newVersion
			}
		}
	}
	return nil
}

type moduleInfo struct {
	Path    string    `json:"Path"`
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

// GetModuleInfo gets version info including timestamp
func (m *moduleOperator) GetModuleInfo(modulePath string) (module.Version, time.Time, error) {
	cmd := exec.Command("go", "list", "-m", "-json", modulePath+"@latest")
	out, err := cmd.Output()
	if err != nil {
		return module.Version{}, time.Time{}, fmt.Errorf("%w: %v", ErrModOperation, err)
	}

	var info moduleInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return module.Version{}, time.Time{}, fmt.Errorf("%w: failed to decode response: %v", ErrModOperation, err)
	}

	return module.Version{
		Path:    info.Path,
		Version: info.Version,
	}, info.Time, nil
}

// ParseModulePathParts extracts org and repo from module path
func (m *moduleOperator) ParseModulePathParts(modulePath string) (org, repo string, err error) {
	parts := strings.Split(modulePath, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("%w: invalid module path format: %s", ErrModOperation, modulePath)
	}

	org = parts[1]
	repo = parts[2]

	// Strip version suffix if present (e.g., "repo/v2" -> "repo")
	if i := strings.Index(repo, "/v"); i != -1 {
		repo = repo[:i]
	}

	return org, repo, nil
}
