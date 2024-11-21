package updater

import (
	"io"
	"os"
)

type SystemOperator interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(filename string, data []byte, perm uint32) error
}

type systemOperator struct {
	stdout io.Writer
	stderr io.Writer
}

func NewSystemOperator() SystemOperator {
	return &systemOperator{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// For testing
func NewSystemOperatorWithIO(stdout, stderr io.Writer) SystemOperator {
	return &systemOperator{
		stdout: stdout,
		stderr: stderr,
	}
}

func (so *systemOperator) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (so *systemOperator) WriteFile(filename string, data []byte, perm uint32) error {
	return os.WriteFile(filename, data, os.FileMode(perm))
}
