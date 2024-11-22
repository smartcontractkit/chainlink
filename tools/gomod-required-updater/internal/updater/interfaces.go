package updater

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// Errors
var (
	ErrModOperation = fmt.Errorf("module operation failed")
)

// ModuleOperator handles module-related operations
type ModuleOperator interface {
	GetGitInfo(remote, branch string) (string, time.Time, error)
	GetLatestVersion(modulePath string) (module.Version, error)
	UpdateRequiredVersions(modFile *modfile.File, newVersion string) ([]string, error)
}

// SystemOperator handles file system operations
type SystemOperator interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
}
