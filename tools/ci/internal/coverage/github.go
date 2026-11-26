package coverage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const githubPageSize = 100

// GitHubClient is the small subset of the GitHub API needed by coverage jobs.
type GitHubClient struct {
	apiURL     string
	webURL     string
	token      string
	repo       string
	httpClient *http.Client
}

// NewGitHubClient constructs a coverage GitHub API client.
func NewGitHubClient(apiURL, webURL, token, repo string, httpClient *http.Client) *GitHubClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &GitHubClient{
		apiURL:     strings.TrimRight(apiURL, "/"),
		webURL:     strings.TrimRight(webURL, "/"),
		token:      token,
		repo:       strings.Trim(repo, "/"),
		httpClient: httpClient,
	}
}

// Artifact identifies a coverage artifact and the workflow run that produced it.
type Artifact struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	RunID      int64     `json:"run_id"`
	RunName    string    `json:"run_name"`
	SHA        string    `json:"sha"`
	Conclusion string    `json:"conclusion"`
	RunURL     string    `json:"run_url"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type githubRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Event      string `json:"event"`
	HeadSHA    string `json:"head_sha"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL    string `json:"html_url"`
}

type githubArtifact struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Expired     bool       `json:"expired"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	WorkflowRun *githubRun `json:"workflow_run"`
}

// ArtifactsForRun returns all unexpired artifacts for a workflow run.
func (c *GitHubClient) ArtifactsForRun(ctx context.Context, runID int64) ([]Artifact, error) {
	if runID <= 0 {
		return nil, fmt.Errorf("workflow run ID must be positive: %d", runID)
	}

	var artifacts []Artifact
	seen := 0
	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/actions/runs/%d/artifacts", c.apiURL, c.repo, runID)
		var response struct {
			TotalCount int              `json:"total_count"`
			Artifacts  []githubArtifact `json:"artifacts"`
		}
		if err := c.doJSON(ctx, http.MethodGet, endpoint, githubPageQuery(page), nil, &response); err != nil {
			return nil, fmt.Errorf("list artifacts for run %d: %w", runID, err)
		}
		seen += len(response.Artifacts)
		for _, item := range response.Artifacts {
			if item.Expired || (!item.ExpiresAt.IsZero() && !item.ExpiresAt.After(time.Now())) {
				continue
			}
			artifact := Artifact{
				ID:        item.ID,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				ExpiresAt: item.ExpiresAt,
				URL:       fmt.Sprintf("%s/%s/actions/runs/%d/artifacts/%d", c.webURL, c.repo, runID, item.ID),
			}
			if item.WorkflowRun != nil {
				artifact.RunID = item.WorkflowRun.ID
				artifact.RunName = item.WorkflowRun.Name
				artifact.SHA = item.WorkflowRun.HeadSHA
				artifact.Conclusion = item.WorkflowRun.Conclusion
				artifact.RunURL = item.WorkflowRun.HTMLURL
			}
			if artifact.RunID == 0 {
				artifact.RunID = runID
			}
			artifacts = append(artifacts, artifact)
		}
		if githubPageDone(len(response.Artifacts), seen, response.TotalCount) {
			break
		}
	}
	return artifacts, nil
}

// LatestNightlyArtifacts finds the newest artifact in each nightly category.
// The returned order is full, unit, integration, e2e, race.
func (c *GitHubClient) LatestNightlyArtifacts(ctx context.Context, branch string) ([]Artifact, error) {
	const (
		full = iota
		unit
		integration
		e2e
		race
		categoryCount
	)
	selected := make([]Artifact, categoryCount)
	found := make([]bool, categoryCount)
	remaining := categoryCount
	seen := 0

	for page := 1; remaining > 0; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/actions/runs", c.apiURL, c.repo)
		query := githubPageQuery(page)
		query.Set("branch", branch)
		query.Set("status", "completed")
		var response struct {
			TotalCount   int         `json:"total_count"`
			WorkflowRuns []githubRun `json:"workflow_runs"`
		}
		if err := c.doJSON(ctx, http.MethodGet, endpoint, query, nil, &response); err != nil {
			return nil, fmt.Errorf("list completed workflow runs for branch %q: %w", branch, err)
		}
		seen += len(response.WorkflowRuns)
		for _, run := range response.WorkflowRuns {
			// Raw coverage and race data are trusted only from scheduled runs.
			// Full reports are trusted only from the collector's workflow_run hook;
			// manual recalculations are not presented as the nightly baseline.
			if run.Event != "schedule" && !(run.Name == "Coverage Collector" && run.Event == "workflow_run") {
				continue
			}
			artifacts, err := c.ArtifactsForRun(ctx, run.ID)
			if err != nil {
				return nil, err
			}
			for _, artifact := range artifacts {
				fillArtifactRunMetadata(&artifact, run)
				category := nightlyArtifactCategory(artifact.Name, run)
				if category < 0 || found[category] {
					continue
				}
				selected[category] = artifact
				found[category] = true
				remaining--
			}
			if remaining == 0 {
				break
			}
		}
		if githubPageDone(len(response.WorkflowRuns), seen, response.TotalCount) {
			break
		}
	}

	result := make([]Artifact, 0, categoryCount)
	for i, artifact := range selected {
		if found[i] {
			result = append(result, artifact)
		}
	}
	return result, nil
}

func nightlyArtifactCategory(name string, run githubRun) int {
	if run.Name == "Coverage Collector" && run.Event == "workflow_run" && strings.HasPrefix(name, "coverage-report-") {
		return 0
	}
	if run.Event != "schedule" {
		return -1
	}
	lowerName := strings.ToLower(name)
	switch {
	case strings.HasPrefix(name, "coverage-core-go_core_race_tests-"):
		return 4
	case strings.HasPrefix(name, "coverage-core-go_core_tests-") &&
		!strings.Contains(lowerName, "integration") && !strings.Contains(lowerName, "race"):
		return 1
	case strings.HasPrefix(name, "coverage-deployment-"), strings.HasPrefix(name, "coverage-inmemory-"):
		return 2
	case strings.HasPrefix(name, "coverage-e2e-"):
		return 3
	default:
		return -1
	}
}

func fillArtifactRunMetadata(artifact *Artifact, run githubRun) {
	if artifact.RunID == 0 {
		artifact.RunID = run.ID
	}
	if artifact.RunName == "" {
		artifact.RunName = run.Name
	}
	if artifact.SHA == "" {
		artifact.SHA = run.HeadSHA
	}
	if artifact.Conclusion == "" {
		artifact.Conclusion = run.Conclusion
	}
	if artifact.RunURL == "" {
		artifact.RunURL = run.HTMLURL
	}
}

// UpsertComment creates or updates the github-actions bot comment with marker.
func (c *GitHubClient) UpsertComment(ctx context.Context, prNumber int, marker, body string) error {
	if prNumber <= 0 {
		return fmt.Errorf("pull request number must be positive: %d", prNumber)
	}

	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/repos/%s/issues/%d/comments", c.apiURL, c.repo, prNumber)
		var comments []struct {
			ID   int64  `json:"id"`
			Body string `json:"body"`
			User struct {
				Login string `json:"login"`
			} `json:"user"`
		}
		if err := c.doJSON(ctx, http.MethodGet, endpoint, githubPageQuery(page), nil, &comments); err != nil {
			return fmt.Errorf("list comments for PR %d: %w", prNumber, err)
		}
		for _, comment := range comments {
			if comment.User.Login != "github-actions[bot]" || !strings.Contains(comment.Body, marker) {
				continue
			}
			updateEndpoint := fmt.Sprintf("%s/repos/%s/issues/comments/%d", c.apiURL, c.repo, comment.ID)
			if err := c.doJSON(ctx, http.MethodPatch, updateEndpoint, nil, struct {
				Body string `json:"body"`
			}{Body: body}, nil); err != nil {
				return fmt.Errorf("update coverage comment %d: %w", comment.ID, err)
			}
			return nil
		}
		if len(comments) < githubPageSize {
			break
		}
	}

	endpoint := fmt.Sprintf("%s/repos/%s/issues/%d/comments", c.apiURL, c.repo, prNumber)
	if err := c.doJSON(ctx, http.MethodPost, endpoint, nil, struct {
		Body string `json:"body"`
	}{Body: body}, nil); err != nil {
		return fmt.Errorf("create coverage comment for PR %d: %w", prNumber, err)
	}
	return nil
}

func githubPageQuery(page int) url.Values {
	query := make(url.Values)
	query.Set("per_page", strconv.Itoa(githubPageSize))
	query.Set("page", strconv.Itoa(page))
	return query
}

// githubPageDone handles both GitHub's total_count responses and plain arrays.
func githubPageDone(pageLength, seen, total int) bool {
	if pageLength == 0 {
		return true
	}
	if total > 0 {
		return seen >= total
	}
	return pageLength < githubPageSize
}

func (c *GitHubClient) doJSON(ctx context.Context, method, endpoint string, query url.Values, input, output any) error {
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	var requestBody io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		requestBody = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, requestBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
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
		return fmt.Errorf("%s %s returned %s: %s", method, req.URL.Path, response.Status, strings.TrimSpace(string(data)))
	}
	if output == nil {
		_, err = io.Copy(io.Discard, response.Body)
		return err
	}
	if err := json.NewDecoder(response.Body).Decode(output); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
