package updater

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// gitExecutor allows mocking git commands in tests
type gitExecutor interface {
	Command(ctx context.Context, args ...string) ([]byte, error)
}

// realGitExecutor implements actual git command execution
type realGitExecutor struct{}

func (g *realGitExecutor) Command(ctx context.Context, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, "git", args...).Output()
}

const (
	goModFile           = "go.mod"
	goModFileMode       = 0644
	gitSHALength        = 12
	gitTimeout          = 30 * time.Second
	gitTimeFormat       = time.RFC3339
	gitRemotePattern    = `^[a-zA-Z0-9][-a-zA-Z0-9_.]*$`
	gitBranchPattern    = `^[a-zA-Z0-9][-a-zA-Z0-9/_]*$`
	majorVersionPattern = `/v\d+$`
)

type Updater struct {
	config *Config
	system SystemOperator
	git    gitExecutor
}

// New creates a new Updater
func New(config *Config, system SystemOperator) *Updater {
	return &Updater{
		config: config,
		system: system,
		git:    &realGitExecutor{},
	}
}

// validateGitInput checks if the remote and branch are in the correct format
func (u *Updater) validateGitInput(remote, branch string) error {
	remoteRE := regexp.MustCompile(gitRemotePattern)
	if !remoteRE.MatchString(remote) {
		return fmt.Errorf("%w: git remote '%s' contains invalid characters", ErrInvalidConfig, remote)
	}

	branchRE := regexp.MustCompile(gitBranchPattern)
	if !branchRE.MatchString(branch) {
		return fmt.Errorf("%w: git branch '%s' contains invalid characters", ErrInvalidConfig, branch)
	}
	return nil
}

// getGitInfo retrieves the latest commit SHA and timestamp from a Git repository
func (u *Updater) getGitInfo(remote, branch string) (string, time.Time, error) {
	if err := u.validateGitInput(remote, branch); err != nil {
		return "", time.Time{}, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	out, err := u.git.Command(ctx, "ls-remote", remote, "refs/heads/"+branch)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to fetch commit SHA from %s/%s: %v",
			ErrModOperation, remote, branch, err)
	}
	if len(out) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: no output from git ls-remote", ErrModOperation)
	}
	sha := strings.Split(string(out), "\t")[0]
	if len(sha) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: empty SHA from git ls-remote", ErrModOperation)
	}

	cmd := exec.CommandContext(ctx, "git", "show", "-s", "--format=%cI", sha)
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

// Run starts the module update process
func (u *Updater) Run() error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if len(u.config.ModulesToUpdate) == 0 {
		logger.Printf("info: auto-detecting modules with local replace directives")
	} else {
		logger.Printf("info: updating modules: %v", u.config.ModulesToUpdate)
	}

	f, err := u.readModFile()
	if err != nil {
		return fmt.Errorf("%w: failed to read and parse go.mod file: %v", ErrModOperation, err)
	}

	// Find modules to update first if none specified
	if len(u.config.ModulesToUpdate) == 0 {
		u.config.ModulesToUpdate, err = u.findLocalReplaceModules()
		if err != nil {
			return fmt.Errorf("%w: failed to detect local replace modules: %v", ErrModOperation, err)
		}
		if len(u.config.ModulesToUpdate) == 0 {
			logger.Printf("info: no modules found to update in %s", f.Module.Mod.Path)
			return nil
		}
		logger.Printf("info: found %d modules with local replace directives: %v",
			len(u.config.ModulesToUpdate), u.config.ModulesToUpdate)
	}

	// Get commit info once for all modules
	sha, commitTime, err := u.getGitInfo(u.config.RepoRemote, u.config.BranchTrunk)
	if err != nil {
		return fmt.Errorf("%w: failed to get git commit info from remote: %v", ErrModOperation, err)
	}

	// Update the modules in the same file handle
	if err := u.updateGoMod(f, sha, commitTime); err != nil {
		return fmt.Errorf("%w: failed to update module versions in go.mod: %v", ErrModOperation, err)
	}

	return u.writeModFile(f)
}

// updateGoMod updates the go.mod file with new pseudo-versions
func (u *Updater) updateGoMod(f *modfile.File, sha string, commitTime time.Time) error {
	for _, modulePath := range u.config.ModulesToUpdate {
		moduleExists := false
		majorVersion := getMajorVersion(modulePath)
		pseudoVersion := module.PseudoVersion(majorVersion, "", commitTime, sha[:gitSHALength])

		// Find and update version
		for _, req := range f.Require {
			if req.Mod.Path == modulePath {
				moduleExists = true
				if u.config.DryRun {
					log.Printf("[DRY RUN] Would update %s: %s => %s", modulePath, req.Mod.Version, pseudoVersion)
					continue
				}

				if err := f.AddRequire(modulePath, pseudoVersion); err != nil {
					return fmt.Errorf("%w: failed to update version for module %s: %v",
						ErrModOperation, modulePath, err)
				}
				break
			}
		}

		if !moduleExists {
			continue
		}
	}

	return nil
}

// getMajorVersion extracts the major version number from a module path
// Returns "v2" for /v2, "v1" for no version suffix
func getMajorVersion(modulePath string) string {
	re := regexp.MustCompile(majorVersionPattern)
	if match := re.FindString(modulePath); match != "" {
		return strings.TrimPrefix(match, "/")
	}
	return "v0"
}

// findLocalReplaceModules finds modules with local replace directives
func (u *Updater) findLocalReplaceModules() ([]string, error) {
	f, err := u.readModFile()
	if err != nil {
		return nil, err
	}

	orgPrefix := fmt.Sprintf("github.com/%s/%s", u.config.OrgName, u.config.RepoName)
	localModules := make(map[string]bool)
	var modules []string

	// First find all local replaces for our org
	for _, rep := range f.Replace {
		if strings.HasPrefix(rep.Old.Path, orgPrefix) &&
			rep.New.Version == "" &&
			isLocalPath(rep.New.Path) {
			localModules[rep.Old.Path] = true
		}
	}

	// Then check requires that match our replaces
	for _, req := range f.Require {
		if localModules[req.Mod.Path] {
			modules = append(modules, req.Mod.Path)
		}
	}

	return modules, nil
}

// isLocalPath checks if the path is a local path
func isLocalPath(path string) bool {
	return path == "." || path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../")
}

// readModFile reads the go.mod file
func (u *Updater) readModFile() (*modfile.File, error) {
	content, err := u.system.ReadFile(goModFile)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read go.mod: %v", ErrModOperation, err)
	}

	f, err := modfile.Parse(goModFile, content, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid go.mod format: %v", ErrModOperation, err)
	}

	return f, nil
}

// writeModFile writes the go.mod file
func (u *Updater) writeModFile(f *modfile.File) error {
	content, err := f.Format()
	if err != nil {
		return fmt.Errorf("%w: failed to format go.mod content: %v", ErrModOperation, err)
	}

	if err := u.system.WriteFile(goModFile, content, goModFileMode); err != nil {
		return fmt.Errorf("%w: failed to write updated go.mod file: %v", ErrModOperation, err)
	}

	return nil
}
