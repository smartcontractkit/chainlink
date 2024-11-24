package updater

import (
	"context"
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
	gitTimeout          = 30 * time.Second
	gitTimeFormat       = time.RFC3339
	gitRemotePattern    = `^[a-zA-Z0-9][-a-zA-Z0-9_.]*$`
	gitBranchPattern    = `^[a-zA-Z0-9][-a-zA-Z0-9/_]*$`
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

// NewModuleOperator creates a new ModuleOperator
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

// validateGitInput checks if the remote and branch are in the correct format
func validateGitInput(remote, branch string) error {
	remoteRE := regexp.MustCompile(gitRemotePattern)
	if !remoteRE.MatchString(remote) {
		return fmt.Errorf("%w: invalid git remote format", ErrModOperation)
	}

	branchRE := regexp.MustCompile(gitBranchPattern)
	if !branchRE.MatchString(branch) {
		return fmt.Errorf("%w: invalid git branch format", ErrModOperation)
	}
	return nil
}

// GetGitInfo retrieves the latest commit SHA and timestamp from a Git repository
func (m *moduleOperator) GetGitInfo(remote, branch string) (string, time.Time, error) {
	// Validate remote and branch against strict regex patterns before using in exec
	if err := validateGitInput(remote, branch); err != nil {
		return "", time.Time{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	// Get latest SHA
	cmd := exec.CommandContext(ctx, "git", "ls-remote", remote, "refs/heads/"+branch)
	out, err := cmd.Output()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to get SHA: %w", ErrModOperation, err)
	}
	if len(out) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: no output from git ls-remote", ErrModOperation)
	}
	sha := strings.Split(string(out), "\t")[0]
	if len(sha) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: empty SHA from git ls-remote", ErrModOperation)
	}

	// Get commit timestamp
	cmd = exec.CommandContext(ctx, "git", "show", "-s", "--format=%cI", sha)
	out, err = cmd.Output()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get commit time: %w", err)
	}
	if len(out) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: no output from git show", ErrModOperation)
	}

	commitTime, err := time.Parse(gitTimeFormat, strings.TrimSpace(string(out)))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse commit time: %w", err)
	}

	return sha[:gitSHALength], commitTime, nil
}

// GetLatestVersion retrieves the latest pseudo-version for a module
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
