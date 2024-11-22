package updater

import (
	"os"
)

type systemOperator struct {
	stdout *os.File
	stderr *os.File
}

func NewSystemOperator() SystemOperator {
	return &systemOperator{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (s *systemOperator) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s *systemOperator) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}
