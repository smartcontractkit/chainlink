package updater

import (
	"os"
	"path/filepath"
)

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
