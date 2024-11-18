package updater

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

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

func parseModuleVersion(modulePath string) string {
    ver := module.Version{Path: modulePath}
    re := regexp.MustCompile(`/v(\d+)$`)
    if match := re.FindStringSubmatch(ver.Path); match != nil {
        return match[1]
    }
    return "0"
}