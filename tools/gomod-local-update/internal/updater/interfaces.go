// Package updater provides functionality to update Go module versions in go.mod files.
// It specializes in handling local replace directives and maintaining consistent
// versioning across related modules.
package updater

import (
	"errors"
	"os"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const (
	gitSHALength = 12
)

// Errors that can be returned by the updater package
var (
	// ErrModOperation indicates a failure in a module-related operation
	ErrModOperation = errors.New("module operation failed")
	// ErrInvalidConfig indicates invalid configuration parameters
	ErrInvalidConfig = errors.New("invalid configuration")
)

// ModuleOperator handles Git repository operations and module version management.
// It provides functionality to retrieve Git information and manage module versions.
type ModuleOperator interface {
	// GetGitInfo retrieves the latest commit SHA and timestamp from a Git repository
	GetGitInfo(remote, branch string) (string, time.Time, error)
	// GetLatestVersion constructs a pseudo-version for a module based on Git information
	GetLatestVersion(modulePath string) (module.Version, error)
	// UpdateRequiredVersions identifies modules that need version updates
	UpdateRequiredVersions(modFile *modfile.File, newVersion string) ([]string, error)
}

// SystemOperator provides an interface for file system operations.
type SystemOperator interface {
	// ReadFile reads the entire contents of a file
	ReadFile(path string) ([]byte, error)
	// WriteFile writes data to a file with specific permissions
	WriteFile(path string, data []byte, perm os.FileMode) error
}
