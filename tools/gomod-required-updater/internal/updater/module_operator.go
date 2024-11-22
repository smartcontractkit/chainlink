package updater

import (
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
// Returns "v2" for /v2, "v1" for no version suffix
func getMajorVersion(modulePath string) string {
	re := regexp.MustCompile(majorVersionPattern)
	if match := re.FindString(modulePath); match != "" {
		return strings.TrimPrefix(match, "/")
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

// extractVersionParts splits a pseudo-version into its components
func extractVersionParts(version string) (base, timestamp, sha string) {
	parts := strings.Split(version, "-")
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

// UpdateRequiredVersions identifies modules that need version updates
func (m *moduleOperator) UpdateRequiredVersions(modFile *modfile.File, newVersion string) ([]string, error) {
	_, _, sha := extractVersionParts(newVersion)
	if sha == "" {
		return nil, fmt.Errorf("%w: invalid version format: %s", ErrModOperation, newVersion)
	}

	var modulesToUpdate []string
	localModules := make(map[string]bool)
	for _, rep := range modFile.Replace {
		if rep.New.Version == "" {
			localModules[rep.Old.Path] = true
		}
	}

	for _, req := range modFile.Require {
		if localModules[req.Mod.Path] {
			modulesToUpdate = append(modulesToUpdate, req.Mod.Path)
		}
	}
	return modulesToUpdate, nil
}
