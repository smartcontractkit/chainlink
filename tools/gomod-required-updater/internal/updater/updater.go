package updater

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const (
    goModFile = "go.mod"
)

type Updater struct {
	mod    ModuleOperator
	system SystemOperator
	config *Config
}


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

	// Use helper method instead of direct file read
	f, err := u.readModFile()
	if err != nil {
		return err
	}

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
	// Use helper method instead of direct file read
	f, err := u.readModFile()
	if err != nil {
		return err
	}

	for _, modulePath := range u.config.ModulesToUpdate {
		moduleExists := false
		majorVersion := getMajorVersion(modulePath)
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

	return u.writeModFile(f)
}

// findLocalReplaceModules finds modules with local replace directives
func (u *Updater) findLocalReplaceModules() ([]string, error) {
	var modules []string
	seen := make(map[string]bool)
	orgPrefix := fmt.Sprintf("github.com/%s/%s", u.config.OrgName, u.config.RepoName)

	// Use helper method instead of direct file read
	f, err := u.readModFile()
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

func (u *Updater) readModFile() (*modfile.File, error) {
    content, err := u.system.ReadFile("go.mod")
    if err != nil {
        return nil, fmt.Errorf("failed to read go.mod: %w", err)
    }

    f, err := modfile.Parse("go.mod", content, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to parse go.mod: %w", err)
    }

    return f, nil
}

func (u *Updater) writeModFile(f *modfile.File) error {
    content, err := f.Format()
    if err != nil {
        return fmt.Errorf("failed to format go.mod: %w", err)
    }

    if err := u.system.WriteFile("go.mod", content, 0644); err != nil {
        return fmt.Errorf("failed to write go.mod: %w", err)
    }

    return nil
}
