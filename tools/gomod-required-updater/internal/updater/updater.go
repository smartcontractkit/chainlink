package updater

import (
	"fmt"
	"log"
	"path/filepath"
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
	
	absRoot, err := filepath.Abs(u.config.RootPath)
	if err != nil {
		return fmt.Errorf("failed to resolve root path: %w", err)
	}

	// If org modules update is enabled, collect all modules with local replaces
	if u.config.UpdateOrgModules {
		modulesToAdd, err := u.findLocalReplaceModules(absRoot)
		if err != nil {
			return fmt.Errorf("failed to find local replace modules: %w", err)
		}
		if len(modulesToAdd) > 0 {
			log.Printf("Found %d modules with local replace directives", len(modulesToAdd))
			u.config.ModulesToUpdate = append(u.config.ModulesToUpdate, modulesToAdd...)
		}

		rootMod, err := u.findRootModule()
		if err != nil {
			return fmt.Errorf("failed to find root module: %w", err)
		}
		u.config.ModulesToUpdate = append(u.config.ModulesToUpdate, rootMod)
	}

	// Get SHA after collecting modules
	log.Printf("Fetching latest SHA from %s/%s", u.config.RepoRemote, u.config.BranchTrunk)
	sha, err := u.git.GetSHA(u.config.RepoRemote, u.config.BranchTrunk)
	if err != nil {
		return fmt.Errorf("failed to get SHA: %w", err)
	}
	log.Printf("Using SHA: %s", sha)

	return u.system.Walk(absRoot, func(path string, isDir bool) error {
		if filepath.Base(path) == "go.mod" {
			for _, module := range u.config.ModulesToUpdate {
				if err := u.updateGoMod(path, module, sha); err != nil {
					return fmt.Errorf("error updating %s: %w", path, err)
				}
			}
		}
		return nil
	})
}

func (u *Updater) findRootModule() (string, error) {
	content, err := u.system.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to read root go.mod: %w", err)
	}

	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse root go.mod: %w", err)
	}

	return f.Module.Mod.Path, nil
}

func (u *Updater) updateGoMod(path, modulePath, sha string) error {
	log.Printf("Processing %s for module %s", path, modulePath)
    
    content, err := u.system.ReadFile(path)
    if (err != nil) {
        return err
    }

    f, err := modfile.Parse(path, content, nil)
    if err != nil {
        return err
    }

    moduleExists := false
    var currentVersion string
    for _, req := range f.Require {
        if req.Mod.Path == modulePath {
            moduleExists = true
            currentVersion = req.Mod.Version
            break
        }
    }

	// Check for replace directive if updating org modules
	if u.config.UpdateOrgModules {
		for _, rep := range f.Replace {
			if rep.Old.Path == modulePath && isLocalPath(rep.New.Path) {
				log.Printf("Found local replace for %s => %s", modulePath, rep.New.Path)
				moduleExists = true
				break
			}
		}
	}

    if !moduleExists {
        log.Printf("Skipping %s: module %s not found", path, modulePath)
        return nil
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
        return nil
    }

    if err := f.AddRequire(modulePath, pseudoVersion); err != nil {
        return fmt.Errorf("failed to add requirement: %w", err)
    }

    newContent, err := f.Format()
    if err != nil {
        return fmt.Errorf("failed to format go.mod: %w", err)
    }

    if err := u.system.WriteFile(path, newContent, 0644); err != nil {
        return fmt.Errorf("failed to write go.mod: %w", err)
    }

    dir := filepath.Dir(path)
    origDir, err := u.system.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get current directory: %w", err)
    }
    
    if err := u.system.Chdir(dir); err != nil {
        return fmt.Errorf("failed to change to directory %s: %w", dir, err)
    }
	// Use defer to ensure we always return to the initial directory.
    defer func() {
        if err := u.system.Chdir(origDir); err != nil {
            log.Printf("Warning: failed to return to original directory: %v", err)
        }
    }()

    log.Printf("Running go mod tidy in %s", dir)
    if err := u.system.RunCommand("go", "mod", "tidy"); err != nil {
        return fmt.Errorf("go mod tidy failed: %w", err)
    }

    return nil
}

func (u *Updater) findLocalReplaceModules(root string) ([]string, error) {
    var modules []string
    seen := make(map[string]bool)
    orgPrefix := fmt.Sprintf("github.com/%s/%s", u.config.OrgName, u.config.RepoName)

    err := u.system.Walk(root, func(path string, isDir bool) error {
        if filepath.Base(path) != "go.mod" {
            return nil
        }

        content, err := u.system.ReadFile(path)
        if err != nil {
            return err
        }

        f, err := modfile.Parse(path, content, nil)
        if err != nil {
            return err
        }

        for _, rep := range f.Replace {
            // Only process modules from our org/repo that have local replaces
            if strings.HasPrefix(rep.Old.Path, orgPrefix) && 
               isLocalPath(rep.New.Path) && 
               !seen[rep.Old.Path] {
                log.Printf("Found local replace in %s: %s => %s", path, rep.Old.Path, rep.New.Path)
                modules = append(modules, rep.Old.Path)
                seen[rep.Old.Path] = true
            }
        }
        return nil
    })

    return modules, err
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