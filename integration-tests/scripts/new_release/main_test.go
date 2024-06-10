package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMainFunction(t *testing.T) {
	t.Run("MissingArguments", func(t *testing.T) {
		os.Args = []string{"main"}
		require.PanicsWithError(t, "usage: go run main.go <repository_name> <days>", func() {
			main()
		})
	})

	t.Run("InvalidDaysArgument", func(t *testing.T) {
		os.Args = []string{"main", "some/repo", "invalid"}
		require.PanicsWithError(t, "error: days must be an integer, but 'invalid' is not an integer", func() {
			main()
		})
	})

	t.Run("ValidArgumentsNoRelease", func(t *testing.T) {
		server := httptest.NewServer(http.NotFoundHandler())
		defer server.Close()

		oldRepoURL := repoURL
		repoURL = server.URL + "/%s/releases/latest"
		defer func() { repoURL = oldRepoURL }()

		os.Args = []string{"main", "some/repo", "30"}
		require.PanicsWithError(t, "error fetching release: unexpected status code: 404\n", func() {
			main()
		})
	})

	t.Run("ValidArgumentsWithRelease", func(t *testing.T) {
		release := Release{
			TagName:     "v1.0.0",
			PublishedAt: time.Now().AddDate(0, 0, -1), // 1 day ago
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(release)
			require.NoError(t, err)
		}))
		defer server.Close()

		oldRepoURL := repoURL
		repoURL = server.URL + "/%s/releases/latest"
		defer func() { repoURL = oldRepoURL }()

		os.Args = []string{"main", "some/repo", "30"}
		output := captureOutput(func() {
			main()
		})
		require.Contains(t, "v1.0.0", output)
	})

	t.Run("ValidArgumentsWithOldRelease", func(t *testing.T) {
		release := Release{
			TagName:     "v1.0.0",
			PublishedAt: time.Now().AddDate(0, 0, -10), // 1 day ago
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(release)
			require.NoError(t, err)
		}))
		defer server.Close()

		oldRepoURL := repoURL
		repoURL = server.URL + "/%s/releases/latest"
		defer func() { repoURL = oldRepoURL }()

		os.Args = []string{"main", "some/repo", "9"}
		output := captureOutput(func() {
			main()
		})
		require.Equal(t, "none\n", output)
	})
}

func TestGetLatestRelease(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		release := Release{
			TagName:     "v1.0.0",
			PublishedAt: time.Now(),
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(release)
			require.NoError(t, err)
		}))
		defer server.Close()

		result, err := getLatestRelease(server.URL, http.DefaultClient)
		require.NoError(t, err)
		require.Equal(t, release.TagName, result.TagName)
		require.True(t, release.PublishedAt.Equal(result.PublishedAt))
	})

	t.Run("Non200StatusCode", func(t *testing.T) {
		server := httptest.NewServer(http.NotFoundHandler())
		defer server.Close()

		_, err := getLatestRelease(server.URL, http.DefaultClient)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected status code: 404")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprintln(w, "invalid json")
			require.NoError(t, err)
		}))
		defer server.Close()

		_, err := getLatestRelease(server.URL, http.DefaultClient)
		require.Error(t, err)
	})
}

func TestIsReleaseRecent(t *testing.T) {
	t.Run("RecentRelease", func(t *testing.T) {
		publishedAt := time.Now().AddDate(0, 0, -5) // 5 days ago
		days := 10
		require.True(t, isReleaseRecent(publishedAt, days))
	})

	t.Run("OldRelease", func(t *testing.T) {
		publishedAt := time.Now().AddDate(0, 0, -15) // 15 days ago
		days := 10
		require.False(t, isReleaseRecent(publishedAt, days))
	})
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
