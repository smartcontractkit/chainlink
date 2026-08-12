package changelog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CommitEntry is one changelog line's worth of commit data.
type CommitEntry struct {
	SHA    string
	Title  string
	PR     int    // extracted from a squash-merge "(#1234)" suffix; 0 if none
	Author string // GitHub @login when known, else the git author name
}

// prSuffix extracts a PR number from a trailing "(#1234)" in a commit title.
var prSuffix = regexp.MustCompile(`\(#[0-9]+\)\s*$`)

var prNumber = regexp.MustCompile(`\(#([0-9]+)\)\s*$`)

// parseTitle splits a commit message's first line into a cleaned title and an
// optional PR number.
func parseTitle(message string) (title string, pr int) {
	title = firstLine(message)
	if m := prNumber.FindStringSubmatch(title); m != nil {
		pr, _ = strconv.Atoi(m[1]) // digits guaranteed by the regex
		title = prSuffix.ReplaceAllString(title, "")
	}
	return strings.TrimSpace(title), pr
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

// compareResult holds the relevant parts of the GitHub compare API response.
type compareResult struct {
	Status       string // "ahead", "behind", "diverged", "identical"
	AheadBy      int
	BehindBy     int
	Commits      []CommitEntry
	TotalCommits int // commits field of the API response (may exceed len(Commits) if capped)
}

// ghClient calls the GitHub REST API.
type ghClient struct {
	httpClient *http.Client
	token      string
	baseURL    string // overridable for tests
}

// newGHClient builds a client using GITHUB_TOKEN/GH_TOKEN, falling back to
// `gh auth token`, and finally to unauthenticated access (public repos only,
// low rate limits).
func newGHClient(ctx context.Context) *ghClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	if token == "" {
		if out, err := exec.CommandContext(ctx, "gh", "auth", "token").Output(); err == nil {
			token = strings.TrimSpace(string(out))
		}
	}
	return &ghClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:    "https://api.github.com",
	}
}

type compareResponse struct {
	Status       string `json:"status"`
	AheadBy      int    `json:"ahead_by"`
	BehindBy     int    `json:"behind_by"`
	TotalCommits int    `json:"total_commits"`
	Commits      []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commit"`
		Author *struct {
			Login string `json:"login"`
		} `json:"author"`
	} `json:"commits"`
}

// Compare runs a GitHub compare of base...head on owner/repo.
func (c *ghClient) Compare(ctx context.Context, owner, repo, base, head string) (compareResult, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/compare/%s...%s", c.baseURL, owner, repo, base, head)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return compareResult{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return compareResult{}, fmt.Errorf("compare %s/%s %s...%s: %w", owner, repo, base, head, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return compareResult{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return compareResult{}, fmt.Errorf("compare %s/%s %s...%s: HTTP %d: %s",
			owner, repo, base, head, resp.StatusCode, truncate(string(body), 200))
	}
	var cr compareResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return compareResult{}, fmt.Errorf("decoding compare response: %w", err)
	}
	res := compareResult{
		Status:       cr.Status,
		AheadBy:      cr.AheadBy,
		BehindBy:     cr.BehindBy,
		TotalCommits: cr.TotalCommits,
	}
	for _, cmt := range cr.Commits {
		title, pr := parseTitle(cmt.Commit.Message)
		author := cmt.Commit.Author.Name
		if cmt.Author != nil && cmt.Author.Login != "" {
			author = "@" + cmt.Author.Login
		}
		res.Commits = append(res.Commits, CommitEntry{
			SHA: cmt.SHA, Title: title, PR: pr, Author: author,
		})
	}
	return res, nil
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
