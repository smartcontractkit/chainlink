package updater

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const (
    goModFile     = "go.mod"
    goModFileMode = 0644
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
        return err
    }

    // Find modules to update first if none specified
    if len(u.config.ModulesToUpdate) == 0 {
        modulesToAdd, err := u.findLocalReplaceModules()
        if err != nil {
            return fmt.Errorf("failed to find local replace modules: %w", err)
        }
        if len(modulesToAdd) == 0 {
            logger.Printf("info: no modules found to update in %s", f.Module.Mod.Path)
            return nil // This is now a non-error case
        }
        u.config.ModulesToUpdate = modulesToAdd
        logger.Printf("info: found %d modules with local replace directives: %v", len(modulesToAdd), modulesToAdd)
    }

    // Get commit info once for all modules
    sha, commitTime, err := u.mod.GetGitInfo(u.config.RepoRemote, u.config.BranchTrunk)
    if err != nil {
        return fmt.Errorf("failed to get git info: %w", err)
    }

    // Update the modules in the same file handle
    if err := u.updateGoMod(f, sha, commitTime); err != nil {
        return fmt.Errorf("error updating %s: %w", goModFile, err)
    }

    // Write the changes
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
                    return fmt.Errorf("failed to add requirement: %w", err)
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

// readModFile reads the go.mod file
func (u *Updater) readModFile() (*modfile.File, error) {
    content, err := u.system.ReadFile(goModFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read %s: %w", goModFile, err)
    }

    f, err := modfile.Parse(goModFile, content, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to parse %s: %w", goModFile, err)
    }

    return f, nil
}

// writeModFile writes the go.mod file
func (u *Updater) writeModFile(f *modfile.File) error {
    content, err := f.Format()
    if err != nil {
        return fmt.Errorf("failed to format %s: %w", goModFile, err)
    }

    if err := u.system.WriteFile(goModFile, content, goModFileMode); err != nil {
        return fmt.Errorf("failed to write %s: %w", goModFile, err)
    }

    return nil
}
