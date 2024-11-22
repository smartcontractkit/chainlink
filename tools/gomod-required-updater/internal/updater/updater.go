package updater

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// New creates a new Updater
func New(mod ModuleOperator, system SystemOperator, config *Config) *Updater {
	return &Updater{
		mod:    mod,
		system: system,
		config: config,
	}
}

// Run the update process
func (u *Updater) Run() error {
	if len(u.config.ModulesToUpdate) == 0 {
		log.Printf("No modules specified, will auto-detect modules with local replace directives")
	} else {
		log.Printf("Starting update process for modules: %v", u.config.ModulesToUpdate)
	}

	// Get org and repo info from current module first
	content, err := u.system.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// Get org/repo from the current module path
	org, repo, err := u.mod.ParseModulePathParts(f.Module.Mod.Path)
	if err != nil {
		return fmt.Errorf("failed to get repo info from current module: %w", err)
	}
	u.config.OrgName = org
	u.config.RepoName = repo

	// Find modules to update first if none specified
	if len(u.config.ModulesToUpdate) == 0 {
		modulesToAdd, err := u.findLocalReplaceModules()
		if err != nil {
			return fmt.Errorf("failed to find local replace modules: %w", err)
		}
		if len(modulesToAdd) == 0 {
			log.Printf("No modules found to update in %s", f.Module.Mod.Path)
			return nil // This is now a non-error case
		}
		u.config.ModulesToUpdate = modulesToAdd
		log.Printf("Found %d modules with local replace directives: %v", len(modulesToAdd), modulesToAdd)
	}

	// Get commit info once for all modules
	sha, commitTime, err := u.mod.GetGitInfo(u.config.RepoRemote, u.config.BranchTrunk)
	if err != nil {
		return fmt.Errorf("failed to get git info: %w", err)
	}

	if err := u.updateGoMod("go.mod", sha, commitTime); err != nil {
		return fmt.Errorf("error updating go.mod: %w", err)
	}

	return nil
}

// updateGoMod updates the go.mod file with new pseudo-versions
func (u *Updater) updateGoMod(path string, sha string, commitTime time.Time) error {
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

		// Get major version from module path suffix (/v2, /v3, etc)
		majorVersion := "v0"
		if idx := strings.LastIndex(modulePath, "/v"); idx != -1 {
			versionSuffix := modulePath[idx+1:] // get everything after the /v
			if _, err := fmt.Sscanf(versionSuffix, "v%d", new(int)); err == nil {
				majorVersion = versionSuffix
			}
		}

		// Create proper pseudo-version using x/mod/module
		// Note: older version parameter is empty since we're creating a new pseudo-version
		pseudoVersion := module.PseudoVersion(majorVersion, "", commitTime, sha[:12])

		// Find and update version
		for _, req := range f.Require {
			if req.Mod.Path == modulePath {
				moduleExists = true
				if u.config.DryRun {
					log.Printf("[DRY RUN] Would update %s: %s => %s", modulePath, req.Mod.Version, pseudoVersion)
					continue
				}

				if err := f.AddRequire(modulePath, pseudoVersion); err != nil {
					return fmt.Errorf("failed to add requirement: %w", err)
				}
				break
			}
		}

		if !moduleExists {
			continue
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

// findLocalReplaceModules finds modules with local replace directives
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

// isLocalPath checks if the path is a local path
func isLocalPath(path string) bool {
	return path == "." || path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../")
}
