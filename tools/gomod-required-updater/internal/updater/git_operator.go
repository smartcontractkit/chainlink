package updater

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type gitOperator struct{}

// NewGitOperator creates a new instance of GitOperator
func NewGitOperator() GitOperator {
	return &gitOperator{}
}

// GetSHA returns the SHA of a given branch in a remote repository
func (g *gitOperator) GetSHA(remote, branch string) (string, error) {
	// Use the full ref path to get exact match
	cmd := exec.Command("git", "ls-remote", remote, "refs/heads/"+branch)
	log.Printf("Running git command: git ls-remote %s refs/heads/%s", remote, branch)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%w: failed to get remote SHA: %v", ErrGitOperation, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", fmt.Errorf("%w: no SHA found for branch %s in remote %s", ErrGitOperation, branch, remote)
	}

	// Parse the first line that exactly matches our ref
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		sha := strings.TrimSpace(parts[0])
		ref := strings.TrimSpace(parts[1])

		if ref == "refs/heads/"+branch {
			log.Printf("Found remote SHA: %s for ref: %s", sha, ref)
			return sha, nil
		}
	}

	return "", fmt.Errorf("%w: no exact match found for refs/heads/%s", ErrGitOperation, branch)
}

// GetCommitDate returns the commit date of a given SHA used for in the go.mod
// psuedo-version for the module.
func (g *gitOperator) GetCommitDate(sha string) (time.Time, error) {
	// Get commit date in ISO 8601 time format
	cmd := exec.Command("git", "show", "-s", "--format=%cI", sha)
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get commit date: %w", err)
	}

	dateStr := strings.TrimSpace(string(output))
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return t, nil
}

func (g *gitOperator) GetRepoInfo(remote string) (org, repo string, err error) {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", remote))
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get repo info for remote %s: %w", remote, err)
	}

	// Handle different URL formats:
	// https://github.com/org/repo.git
	// git@github.com:org/repo.git
	url := strings.TrimSpace(string(output))
	parts := strings.Split(strings.TrimSuffix(url, ".git"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("unexpected git URL format from remote %s: %s", remote, url)
	}

	repo = parts[len(parts)-1]
	org = parts[len(parts)-2]
	return org, repo, nil
}
