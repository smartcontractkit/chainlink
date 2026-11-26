package coverage

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunsForCommitAndDownloadArtifact(t *testing.T) {
	t.Parallel()
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	entry, err := zw.Create("nested/covmeta.test")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("coverage")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/org/repo/actions/runs":
			if r.URL.Query().Get("head_sha") != "abc" {
				t.Errorf("missing head_sha query: %s", r.URL.RawQuery)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"total_count":1,"workflow_runs":[{"id":42,"name":"CI Core","status":"completed","conclusion":"success"}]}`)
		case "/repos/org/repo/actions/artifacts/10/zip":
			w.Header().Set("Content-Type", "application/zip")
			_, _ = w.Write(archive.Bytes())
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewGitHubClient(server.URL, "https://github.com", "token", "org/repo", server.Client())
	runs, err := client.RunsForCommit(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 || runs[0].ID != 42 || runs[0].Name != "CI Core" || runs[0].Status != "completed" {
		t.Fatalf("unexpected runs: %+v", runs)
	}
	destination := filepath.Join(t.TempDir(), "artifact")
	if err := client.DownloadArtifact(context.Background(), Artifact{ID: 10, Name: "coverage-unit"}, destination); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(destination, "nested", "covmeta.test"))
	if err != nil || string(data) != "coverage" {
		t.Fatalf("downloaded data=%q err=%v", data, err)
	}
}

func TestDownloadArtifactRejectsZipSlip(t *testing.T) {
	t.Parallel()
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	entry, err := zw.Create("../escape")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte("bad"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive.Bytes())
	}))
	defer server.Close()
	client := NewGitHubClient(server.URL, "https://github.com", "", "org/repo", server.Client())
	root := t.TempDir()
	if err := client.DownloadArtifact(context.Background(), Artifact{ID: 10}, filepath.Join(root, "artifact")); err == nil {
		t.Fatal("expected zip-slip path to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, "escape")); !os.IsNotExist(err) {
		t.Fatal("zip-slip file was written")
	}
}
