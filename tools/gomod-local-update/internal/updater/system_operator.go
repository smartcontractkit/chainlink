package updater

import (
	"os"
	"path/filepath"
)

// SystemOperator provides an interface for file system operations.
type SystemOperator interface {
	// ReadFile reads the entire contents of a file
	ReadFile(path string) ([]byte, error)
	// WriteFile writes data to a file with specific permissions
	WriteFile(path string, data []byte, perm os.FileMode) error
}

type systemOperator struct{}

func NewSystemOperator() SystemOperator {
	return &systemOperator{}
}

func (s *systemOperator) ReadFile(path string) ([]byte, error) {
	path = filepath.Clean(path)
	return os.ReadFile(path)
}

func (s *systemOperator) WriteFile(path string, data []byte, perm os.FileMode) error {
	path = filepath.Clean(path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
