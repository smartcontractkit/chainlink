package coverage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestArtifactsForRunPaginatesFiltersAndUsesDirectLinks(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, time.November, 25, 10, 11, 12, 0, time.UTC)
	expires := created.Add(10 * 365 * 24 * time.Hour)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/repos/acme/project/actions/runs/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("Authorization = %q, want Bearer secret", got)
		}
		if got := r.URL.Query().Get("per_page"); got != "100" {
			t.Errorf("per_page = %q, want 100", got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("page") {
		case "1":
			fmt.Fprintf(w, `{"total_count":3,"artifacts":[{
				"id":101,"name":"coverage-report-main-42","expired":false,
				"created_at":%q,"expires_at":%q,"archive_download_url":"https://download.invalid/101",
				"workflow_run":{"id":42,"name":"Coverage Collector","head_sha":"abc123","conclusion":"success","html_url":"https://web.invalid/acme/project/actions/runs/42"}
			},{
				"id":102,"name":"expired","expired":true,"created_at":%q,"expires_at":%q
			}]}`, created.Format(time.RFC3339), expires.Format(time.RFC3339), created.Format(time.RFC3339), expires.Format(time.RFC3339))
		case "2":
			fmt.Fprintf(w, `{"total_count":3,"artifacts":[{
				"id":103,"name":"coverage-e2e-system-42","expired":false,
				"created_at":%q,"expires_at":%q
			}]}`, created.Format(time.RFC3339), expires.Format(time.RFC3339))
		default:
			http.Error(w, "unexpected page", http.StatusBadRequest)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewGitHubClient(server.URL+"/api", server.URL+"/web", "secret", "acme/project", server.Client())
	artifacts, err := client.ArtifactsForRun(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 {
		t.Fatalf("got %d artifacts, want 2: %+v", len(artifacts), artifacts)
	}
	first := artifacts[0]
	if first.ID != 101 || first.Name != "coverage-report-main-42" || first.RunID != 42 || first.RunName != "Coverage Collector" || first.SHA != "abc123" || first.Conclusion != "success" {
		t.Errorf("unexpected artifact metadata: %+v", first)
	}
	if first.RunURL != "https://web.invalid/acme/project/actions/runs/42" {
		t.Errorf("RunURL = %q", first.RunURL)
	}
	if first.URL != server.URL+"/web/acme/project/actions/runs/42/artifacts/101" {
		t.Errorf("URL = %q, want direct web artifact URL", first.URL)
	}
	if !first.CreatedAt.Equal(created) || !first.ExpiresAt.Equal(expires) {
		t.Errorf("unexpected timestamps: %+v", first)
	}
	if artifacts[1].ID != 103 || artifacts[1].URL != server.URL+"/web/acme/project/actions/runs/42/artifacts/103" {
		t.Errorf("unexpected second-page artifact: %+v", artifacts[1])
	}
}

func TestLatestNightlyArtifactsSelectsNewestByCategory(t *testing.T) {
	t.Parallel()

	type runFixture struct {
		id         int64
		name       string
		event      string
		sha        string
		conclusion string
		artifacts  []map[string]any
	}
	runs := []runFixture{
		{id: 65, name: "Coverage Collector", event: "workflow_dispatch", sha: "manual-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(650), "name": "coverage-report-manual-sha-65"},
		}},
		{id: 60, name: "CI Core", event: "pull_request", sha: "pr-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(600), "name": "coverage-core-go_core_tests-pr-60"},
			{"id": int64(601), "name": "coverage-core-go_core_race_tests-pr-60"},
		}},
		{id: 55, name: "Coverage Collector", event: "workflow_run", sha: "full-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(550), "name": "coverage-report-full-sha-55"},
		}},
		{id: 50, name: "CI Core", event: "schedule", sha: "core-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(500), "name": "coverage-core-go_core_tests_integration-core-sha-50"},
			{"id": int64(501), "name": "coverage-core-go_core_race_tests-core-sha-50"},
			{"id": int64(502), "name": "coverage-core-go_core_tests-core-sha-50"},
		}},
		{id: 45, name: "Deployment Tests", event: "schedule", sha: "integration-sha", conclusion: "failure", artifacts: []map[string]any{
			{"id": int64(450), "name": "coverage-deployment-shard-1-45"},
			{"id": int64(451), "name": "coverage-inmemory-shard-1-45"},
		}},
		{id: 40, name: "System Tests", event: "schedule", sha: "e2e-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(400), "name": "coverage-e2e-system-40"},
		}},
		{id: 35, name: "CI Core", event: "schedule", sha: "old-sha", conclusion: "success", artifacts: []map[string]any{
			{"id": int64(350), "name": "coverage-core-go_core_tests-old-sha-35"},
		}},
	}

	var mu sync.Mutex
	artifactRequests := make(map[int64]int)
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/acme/project/actions/runs", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("branch") != "develop" || query.Get("status") != "completed" || query.Get("per_page") != "100" {
			t.Errorf("unexpected run query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"total_count":7,"workflow_runs":[`)
		for i, run := range runs {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `{"id":%d,"name":%q,"event":%q,"head_sha":%q,"conclusion":%q,"html_url":%q}`,
				run.id, run.name, run.event, run.sha, run.conclusion, fmt.Sprintf("https://github.invalid/acme/project/actions/runs/%d", run.id))
		}
		fmt.Fprint(w, `]}`)
	})
	mux.HandleFunc("/repos/acme/project/actions/runs/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		runID, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
		if err != nil || parts[len(parts)-1] != "artifacts" {
			http.Error(w, "bad artifact path", http.StatusBadRequest)
			return
		}
		mu.Lock()
		artifactRequests[runID]++
		mu.Unlock()
		for _, run := range runs {
			if run.id == runID {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{"total_count": len(run.artifacts), "artifacts": run.artifacts})
				return
			}
		}
		http.NotFound(w, r)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewGitHubClient(server.URL, "https://github.invalid", "token", "acme/project", server.Client())
	artifacts, err := client.LatestNightlyArtifacts(context.Background(), "develop")
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 5 {
		t.Fatalf("got %d artifacts, want full, unit, integration, e2e, race: %+v", len(artifacts), artifacts)
	}
	wantIDs := []int64{550, 502, 450, 400, 501}
	wantSHAs := []string{"full-sha", "core-sha", "integration-sha", "e2e-sha", "core-sha"}
	for i, artifact := range artifacts {
		if artifact.ID != wantIDs[i] || artifact.SHA != wantSHAs[i] {
			t.Errorf("artifact %d = %+v, want ID %d and SHA %q", i, artifact, wantIDs[i], wantSHAs[i])
		}
		if artifact.RunID == 0 || artifact.RunName == "" || artifact.RunURL == "" || artifact.Conclusion == "" {
			t.Errorf("artifact %d is missing run metadata: %+v", i, artifact)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if artifactRequests[65] != 0 {
		t.Error("requested a manual collector report as a nightly baseline")
	}
	if artifactRequests[60] != 0 {
		t.Error("requested raw/race artifacts for a non-scheduled run")
	}
	if artifactRequests[35] != 0 {
		t.Error("continued requesting artifacts after every category was found")
	}
}

func TestUpsertCommentCreatesWhenNoBotMarkerExists(t *testing.T) {
	t.Parallel()

	var posts int
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/acme/project/issues/17/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("per_page") != "100" {
				t.Errorf("unexpected comments query: %s", r.URL.RawQuery)
			}
			fmt.Fprint(w, `[]`)
		case http.MethodPost:
			posts++
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Error(err)
			}
			if payload["body"] != "<!-- coverage -->\nreport" {
				t.Errorf("POST body = %q", payload["body"])
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":1}`)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewGitHubClient(server.URL, server.URL, "token", "acme/project", server.Client())
	if err := client.UpsertComment(context.Background(), 17, "<!-- coverage -->", "<!-- coverage -->\nreport"); err != nil {
		t.Fatal(err)
	}
	if posts != 1 {
		t.Fatalf("POST requests = %d, want 1", posts)
	}
}

func TestUpsertCommentUpdatesBotMarkerOnLaterPage(t *testing.T) {
	t.Parallel()

	var patches, posts int
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/acme/project/issues/17/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			posts++
			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "1" {
			comments := make([]map[string]any, 100)
			for i := range comments {
				comments[i] = map[string]any{"id": i + 1, "body": "unrelated", "user": map[string]string{"login": "someone"}}
			}
			_ = json.NewEncoder(w).Encode(comments)
			return
		}
		fmt.Fprint(w, `[
			{"id":901,"body":"<!-- coverage -->","user":{"login":"someone"}},
			{"id":902,"body":"unrelated","user":{"login":"github-actions[bot]"}},
			{"id":903,"body":"old <!-- coverage --> report","user":{"login":"github-actions[bot]"}}
		]`)
	})
	mux.HandleFunc("/repos/acme/project/issues/comments/903", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		patches++
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		if payload["body"] != "<!-- coverage -->\nnew report" {
			t.Errorf("PATCH body = %q", payload["body"])
		}
		fmt.Fprint(w, `{"id":903}`)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewGitHubClient(server.URL, server.URL, "token", "acme/project", server.Client())
	if err := client.UpsertComment(context.Background(), 17, "<!-- coverage -->", "<!-- coverage -->\nnew report"); err != nil {
		t.Fatal(err)
	}
	if patches != 1 || posts != 0 {
		t.Fatalf("PATCH requests = %d, POST requests = %d; want 1 and 0", patches, posts)
	}
}

func TestUpsertCommentValidatesPRNumber(t *testing.T) {
	t.Parallel()

	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)
	client := NewGitHubClient(server.URL, server.URL, "token", "acme/project", server.Client())
	for _, prNumber := range []int{0, -1} {
		if err := client.UpsertComment(context.Background(), prNumber, "marker", "body"); err == nil {
			t.Errorf("UpsertComment(%d) returned nil error", prNumber)
		}
	}
	if requests != 0 {
		t.Fatalf("made %d requests for invalid PR numbers", requests)
	}
}

func TestArtifactsForRunRejectsInvalidRunID(t *testing.T) {
	t.Parallel()
	client := NewGitHubClient("https://api.invalid", "https://github.invalid", "", "acme/project", &http.Client{})
	if _, err := client.ArtifactsForRun(context.Background(), 0); err == nil {
		t.Fatal("ArtifactsForRun(0) returned nil error")
	}
}

func TestLatestNightlyArtifactsEscapesBranchQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("branch"); got != "release/next" {
			t.Errorf("branch = %q", got)
		}
		if _, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			t.Errorf("invalid query: %v", err)
		}
		fmt.Fprint(w, `{"total_count":0,"workflow_runs":[]}`)
	}))
	t.Cleanup(server.Close)
	client := NewGitHubClient(server.URL, server.URL, "", "acme/project", server.Client())
	artifacts, err := client.LatestNightlyArtifacts(context.Background(), "release/next")
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 0 {
		t.Fatalf("artifacts = %+v, want none", artifacts)
	}
}
