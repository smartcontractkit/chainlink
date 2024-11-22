package updater

import (
	"fmt"
	"time"

	"golang.org/x/mod/module"
)

// Errors
var (
	ErrInvalidConfig = fmt.Errorf("invalid configuration")
	ErrModOperation  = fmt.Errorf("module operation failed")
	ErrFileOperation = fmt.Errorf("file operation failed")
)

type ModuleOperator interface {
	// GetLatestVersion gets the latest pseudo-version based on current git state
	GetLatestVersion(modulePath string) (module.Version, error)
	// GetModuleInfo gets version info including timestamp
	GetModuleInfo(modulePath string) (module.Version, time.Time, error)
	// ParseModulePathParts extracts org and repo from module path
	ParseModulePathParts(modulePath string) (org, repo string, err error)
	// GetGitInfo gets the latest git SHA and commit time
	GetGitInfo(remote, branch string) (sha string, commitTime time.Time, err error)
}

type Updater struct {
	mod    ModuleOperator
	system SystemOperator
	config *Config
}
