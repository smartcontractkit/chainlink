package updater

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func New(git GitOperator, system SystemOperator, config *Config) *Updater {
	return &Updater{
		git:    git,
		system: system,
		config: config,
	}
}

func (u *Updater) Run() error {
	log.Printf("Starting update process with remote '%s' and branch '%s'", u.config.RepoRemote, u.config.BranchTrunk)

	// Get org and repo info for finding local modules
	org, repo, err := u.git.GetRepoInfo(u.config.RepoRemote)
	if err != nil {
		return fmt.Errorf("failed to get repo info: %w", err)
	}
	u.config.OrgName = org
	u.config.RepoName = repo

	// Find modules to update
	modulesToAdd, err := u.findLocalReplaceModules()
	if err != nil {
		return fmt.Errorf("failed to find local replace modules: %w", err)
	}
	if len(modulesToAdd) > 0 {
		log.Printf("Found %d modules with local replace directives", len(modulesToAdd))
		u.config.ModulesToUpdate = append(u.config.ModulesToUpdate, modulesToAdd...)
	}

	// Get SHA after collecting modules
	log.Printf("Fetching latest SHA from %s/%s", u.config.RepoRemote, u.config.BranchTrunk)
	sha, err := u.git.GetSHA(u.config.RepoRemote, u.config.BranchTrunk)
	if err != nil {
		return fmt.Errorf("failed to get SHA: %w", err)
	}
	log.Printf("Using SHA: %s", sha)

	// Update go.mod in current directory
	if err := u.updateGoMod("go.mod", sha); err != nil {
		return fmt.Errorf("error updating go.mod: %w", err)
	}

	return nil
}

func (u *Updater) updateGoMod(path string, sha string) error {
	content, err := u.system.ReadFile(path)
	if err != nil {
		return err
	}

	f, err := modfile.Parse(path, content, nil)
	if err != nil {
		return err
	}

	for _, modulePath := range u.config.ModulesToUpdate {
		moduleExists := false
		var currentVersion string
		for _, req := range f.Require {
			if req.Mod.Path == modulePath {
				moduleExists = true
				currentVersion = req.Mod.Version
				break
			}
		}

		// Check for replace directive
		for _, rep := range f.Replace {
			if rep.Old.Path == modulePath && isLocalPath(rep.New.Path) {
				log.Printf("Found local replace for %s => %s", modulePath, rep.New.Path)
				moduleExists = true
				break
			}
		}

		if !moduleExists {
			continue
		}

		log.Printf("Current version: %s", currentVersion)

		commitDate, err := u.git.GetCommitDate(sha)
		if err != nil {
			return fmt.Errorf("failed to get commit date: %w", err)
		}

		shortSHA := sha[:12]
		versionPrefix := parseModuleVersion(modulePath)
		pseudoVersion := module.PseudoVersion("v"+versionPrefix, "", commitDate, shortSHA)

		log.Printf("Updating to version: %s", pseudoVersion)

		if u.config.DryRun {
			log.Printf("[DRY RUN] Would update %s: %s => %s", modulePath, currentVersion, pseudoVersion)
			continue
		}

		if err := f.AddRequire(modulePath, pseudoVersion); err != nil {
			return fmt.Errorf("failed to add requirement: %w", err)
		}
	}

	newContent, err := f.Format()
	if err != nil {
		return fmt.Errorf("failed to format go.mod: %w", err)
	}

	if err := u.system.WriteFile(path, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	return nil
}

func (u *Updater) findLocalReplaceModules() ([]string, error) {
	var modules []string
	seen := make(map[string]bool)
	orgPrefix := fmt.Sprintf("github.com/%s/%s", u.config.OrgName, u.config.RepoName)

	content, err := u.system.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	for _, rep := range f.Replace {
		// Only process modules from our org/repo that have local replaces
		if strings.HasPrefix(rep.Old.Path, orgPrefix) &&
			isLocalPath(rep.New.Path) &&
			!seen[rep.Old.Path] {
			log.Printf("Found local replace: %s => %s", rep.Old.Path, rep.New.Path)
			modules = append(modules, rep.Old.Path)
			seen[rep.Old.Path] = true
		}
	}

	return modules, nil
}

func parseModuleVersion(modulePath string) string {
	ver := module.Version{Path: modulePath}
	re := regexp.MustCompile(`/v(\d+)$`)
	if match := re.FindStringSubmatch(ver.Path); match != nil {
		return match[1]
	}
	return "0"
}

func isLocalPath(path string) bool {
	return path == "." || path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../")
}
