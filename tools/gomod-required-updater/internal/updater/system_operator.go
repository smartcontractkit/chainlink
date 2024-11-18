package updater

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

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

func (so *systemOperator) Walk(root string, fn func(path string, isDir bool) error) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return fn(path, info.IsDir())
	})
}

func (so *systemOperator) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (so *systemOperator) Getwd() (string, error) {
	return os.Getwd()
}

func (so *systemOperator) RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = so.stdout
	cmd.Stderr = so.stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %v", ErrFileOperation, err)
	}
	return nil
}