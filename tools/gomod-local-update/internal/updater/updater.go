package updater

import (
	"context"
	"errors"
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

// systemGitExecutor executes git commands on the host system
type systemGitExecutor struct{}

func (g *systemGitExecutor) Command(ctx context.Context, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, "git", args...).Output()
}

const (
	// File and mode constants
	goModFile     = "go.mod"
	goModFileMode = 0644
	gitSHALength  = 12
	gitTimeout    = 30 * time.Second
	gitTimeFormat = time.RFC3339

	// Regex pattern constants
	gitRemotePattern    = `^[a-zA-Z0-9][-a-zA-Z0-9_.]*$`
	gitBranchPattern    = `^[a-zA-Z0-9][-a-zA-Z0-9/_]*$`
	majorVersionPattern = `/v\d+$`
	shaPattern          = `^[a-fA-F0-9]{40}$` // SHA-1 hashes are 40 hexadecimal characters
)

var (
	// Pre-compiled regular expressions
	gitRemoteRE    = regexp.MustCompile(gitRemotePattern)
	gitBranchRE    = regexp.MustCompile(gitBranchPattern)
	gitShaRE       = regexp.MustCompile(shaPattern)
	majorVersionRE = regexp.MustCompile(majorVersionPattern)
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
		git:    &systemGitExecutor{},
	}
}

// validateSHA checks if the SHA consists of exactly 40 hexadecimal digits
func (u *Updater) validateSHA(sha string) error {
	if !gitShaRE.MatchString(sha) {
		return fmt.Errorf("%w: invalid git SHA '%s'", ErrInvalidConfig, sha)
	}
	return nil
}

// getGitInfo retrieves the latest commit SHA and timestamp from a Git repository
func (u *Updater) getGitInfo(remote, branch string) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	// Use u.git.Command for ls-remote
	out, err := u.git.Command(ctx, "ls-remote", remote, "refs/heads/"+branch)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: failed to fetch commit SHA from %s/%s: %w",
			ErrModOperation, remote, branch, err)
	}
	if len(out) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: no output from git ls-remote", ErrModOperation)
	}
	sha := strings.Split(string(out), "\t")[0]
	if len(sha) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: empty SHA from git ls-remote", ErrModOperation)
	}

	// Validate the SHA
	if err := u.validateSHA(sha); err != nil {
		return "", time.Time{}, err
	}

	// Use u.git.Command for show
	showOut, showErr := u.git.Command(ctx, "show", "-s", "--format=%cI", sha)
	if showErr != nil {
		return "", time.Time{}, fmt.Errorf("failed to get commit time: %w", showErr)
	}

	commitTime, parseErr := time.Parse(gitTimeFormat, strings.TrimSpace(string(showOut)))
	if parseErr != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse commit time: %w", parseErr)
	}

	return sha[:gitSHALength], commitTime, nil
}

// Run starts the module update process
func (u *Updater) Run() error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("info: auto-detecting modules with local replace directives")

	f, err := u.readModFile()
	if err != nil {
		return fmt.Errorf("%w: failed to read and parse go.mod file", ErrModOperation)
	}

	// Auto-detect modules to update
	modulesToUpdate, err := u.findLocalReplaceModules()
	if err != nil {
		return fmt.Errorf("%w: failed to detect local replace modules", ErrModOperation)
	}
	if len(modulesToUpdate) == 0 {
		logger.Printf("info: no modules found to update in %s", f.Module.Mod.Path)
		return nil
	}
	logger.Printf("info: found %d modules with local replace directives: %v",
		len(modulesToUpdate), modulesToUpdate)

	// Get commit info once for all modules
	sha, commitTime, err := u.getGitInfo(u.config.RepoRemote, u.config.BranchTrunk)
	if err != nil {
		if errors.Is(err, ErrInvalidConfig) {
			return err
		}
		return fmt.Errorf("%w: failed to get git commit info from remote", ErrModOperation)
	}

	// Update the modules in the same file handle
	if err := u.updateGoMod(f, modulesToUpdate, sha, commitTime); err != nil {
		return fmt.Errorf("%w: failed to update module versions in go.mod", ErrModOperation)
	}

	return u.writeModFile(f)
}

// updateGoMod updates the go.mod file with new pseudo-versions
func (u *Updater) updateGoMod(modFile *modfile.File, modulesToUpdate []string, sha string, commitTime time.Time) error {
	for _, modulePath := range modulesToUpdate {
		majorVersion := getMajorVersion(modulePath)
		pseudoVersion := module.PseudoVersion(majorVersion, "", commitTime, sha[:gitSHALength])

		// Find and update version
		for _, req := range modFile.Require {
			if req.Mod.Path == modulePath {
				if u.config.DryRun {
					log.Printf("[DRY RUN] Would update %s: %s => %s", modulePath, req.Mod.Version, pseudoVersion)
					continue
				}

				if err := modFile.AddRequire(modulePath, pseudoVersion); err != nil {
					return fmt.Errorf("%w: failed to update version for module %s",
						ErrModOperation, modulePath)
				}
				break
			}
		}
	}

	return nil
}

// getMajorVersion extracts the major version number from a module path
// Returns "v2" for /v2, "v0" for no version suffix
func getMajorVersion(modulePath string) string {
	if match := majorVersionRE.FindString(modulePath); match != "" {
		return strings.TrimPrefix(match, "/")
	}
	return "v0"
}

// findLocalReplaceModules finds modules with local replace directives
func (u *Updater) findLocalReplaceModules() ([]string, error) {
	modFile, err := u.readModFile()
	if err != nil {
		return nil, err
	}

	orgPrefix := fmt.Sprintf("github.com/%s/%s", u.config.OrgName, u.config.RepoName)
	localModules := make(map[string]bool)
	var modules []string

	// First find all local replaces for our org
	for _, rep := range modFile.Replace {
		if strings.HasPrefix(rep.Old.Path, orgPrefix) &&
			rep.New.Version == "" &&
			isLocalPath(rep.New.Path) {
			localModules[rep.Old.Path] = true
		}
	}

	// Then check requires that match our replaces
	for _, req := range modFile.Require {
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
		return nil, fmt.Errorf("unable to read go.mod: %w", err)
	}

	modFile, err := modfile.Parse(goModFile, content, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid go.mod format: %w", ErrModOperation, err)
	}

	return modFile, nil
}

// writeModFile writes the go.mod file
func (u *Updater) writeModFile(modFile *modfile.File) error {
	content, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("%w: failed to format go.mod content", ErrModOperation)
	}

	if err := u.system.WriteFile(goModFile, content, goModFileMode); err != nil {
		return fmt.Errorf("%w: failed to write updated go.mod file", ErrModOperation)
	}

	return nil
}
