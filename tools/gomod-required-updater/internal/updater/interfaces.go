package updater

import (
	"fmt"
	"time"
)

// Errors
var (
	ErrInvalidConfig = fmt.Errorf("invalid configuration")
	ErrGitOperation  = fmt.Errorf("git operation failed")
	ErrFileOperation = fmt.Errorf("file operation failed")
)

type GitOperator interface {
	GetSHA(remote, branch string) (string, error)
	GetCommitDate(sha string) (time.Time, error)
	GetRepoInfo(remote string) (org, repo string, err error)
}

type SystemOperator interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(filename string, data []byte, perm uint32) error
	Walk(root string, fn func(path string, isDir bool) error) error
	Chdir(dir string) error
	Getwd() (string, error)
	RunCommand(name string, args ...string) error
}

type Updater struct {
	git    GitOperator
	system SystemOperator
	config *Config
}
