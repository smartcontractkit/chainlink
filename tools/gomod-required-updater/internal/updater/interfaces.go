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

type Updater struct {
	git    GitOperator
	system SystemOperator
	config *Config
}
