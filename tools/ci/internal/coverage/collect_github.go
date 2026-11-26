package coverage

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RunsForCommit lists workflow runs for a commit, newest first.
func (c *GitHubClient) RunsForCommit(ctx context.Context, sha string) ([]WorkflowRun, error) {
	if sha == "" {
		return nil, fmt.Errorf("commit SHA is required")
	}
	var runs []WorkflowRun
	seen := 0
	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/actions/runs", c.apiURL, c.repo)
		query := githubPageQuery(page)
		query.Set("head_sha", sha)
		var response struct {
			TotalCount   int         `json:"total_count"`
			WorkflowRuns []githubRun `json:"workflow_runs"`
		}
		if err := c.doJSON(ctx, http.MethodGet, endpoint, query, nil, &response); err != nil {
			return nil, fmt.Errorf("list runs for commit %s: %w", sha, err)
		}
		seen += len(response.WorkflowRuns)
		for _, run := range response.WorkflowRuns {
			runs = append(runs, WorkflowRun{ID: run.ID, Name: run.Name, Status: run.Status, Conclusion: run.Conclusion})
		}
		if githubPageDone(len(response.WorkflowRuns), seen, response.TotalCount) {
			break
		}
	}
	return runs, nil
}

// DownloadArtifact downloads and safely extracts an artifact ZIP.
func (c *GitHubClient) DownloadArtifact(ctx context.Context, artifact Artifact, destination string) error {
	if artifact.ID <= 0 {
		return fmt.Errorf("artifact ID must be positive: %d", artifact.ID)
	}
	endpoint := fmt.Sprintf("%s/repos/%s/actions/artifacts/%d/zip", c.apiURL, c.repo, artifact.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return fmt.Errorf("download artifact %s returned %s: %s", artifact.Name, response.Status, strings.TrimSpace(string(data)))
	}

	tmp, err := os.CreateTemp("", "coverage-artifact-*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := io.Copy(tmp, response.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	archive, err := zip.OpenReader(tmpPath)
	if err != nil {
		return fmt.Errorf("open artifact %s: %w", artifact.Name, err)
	}
	defer archive.Close()
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	root, err := filepath.Abs(destination)
	if err != nil {
		return err
	}
	for _, entry := range archive.File {
		path := filepath.Join(root, filepath.Clean(entry.Name))
		if path != root && !strings.HasPrefix(path, root+string(filepath.Separator)) {
			return fmt.Errorf("artifact %s contains unsafe path %q", artifact.Name, entry.Name)
		}
		if entry.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		input, err := entry.Open()
		if err != nil {
			return err
		}
		output, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			input.Close()
			return err
		}
		_, copyErr := io.Copy(output, input)
		closeErr := output.Close()
		inputErr := input.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if inputErr != nil {
			return inputErr
		}
	}
	return nil
}
