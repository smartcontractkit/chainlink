package updater

import (
	"fmt"
	"log"
	"strings"

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

	// Get latest version once and apply to all modules
	version, err := u.mod.GetLatestVersion(u.config.ModulesToUpdate[0])
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}
	log.Printf("Using version: %s for all modules", version.Version)

	if err := u.updateGoMod("go.mod", version); err != nil {
		return fmt.Errorf("error updating go.mod: %w", err)
	}

	return nil
}

// updateGoMod updates the go.mod file with the new module version
func (u *Updater) updateGoMod(path string, newVersion module.Version) error {
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

        // Get timestamp and commit hash from version string
        parts := strings.Split(newVersion.Version, "-")
        var timestamp, commitHash string
        if len(parts) >= 3 {
            timestamp = parts[1]
            commitHash = parts[2]
        } else {
            timestamp = "00000000000000"
            commitHash = newVersion.Version // use full version as commit hash if can't parse
        }

        // Format the version based on module's major version from path
        targetVersion := fmt.Sprintf("%s.0.0-%s-%s", majorVersion, timestamp, commitHash)
        if majorVersion == "v0" {
            targetVersion = fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHash)
        }

		// Find current version
		for _, req := range f.Require {
			if req.Mod.Path == modulePath {
				moduleExists = true
				if u.config.DryRun {
					log.Printf("[DRY RUN] Would update %s: %s => %s", modulePath, req.Mod.Version, targetVersion)
					continue
				}

				if err := f.AddRequire(modulePath, targetVersion); err != nil {
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
