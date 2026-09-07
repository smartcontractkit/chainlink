package engine

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseTitle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		message   string
		wantTitle string
		wantPR    int
	}{
		{"feat: add lane (#1234)", "feat: add lane", 1234},
		{"feat: add lane (#1234)\n\nSome body text\n", "feat: add lane", 1234},
		{"direct push commit", "direct push commit", 0},
		{"fix: parens (not a PR) in middle (#42)", "fix: parens (not a PR) in middle", 42},
	}
	for _, c := range cases {
		title, pr := parseTitle(c.message)
		if title != c.wantTitle || pr != c.wantPR {
			t.Errorf("parseTitle(%q) = (%q, %d), want (%q, %d)", c.message, title, pr, c.wantTitle, c.wantPR)
		}
	}
}

const compareFixture = `{
  "status": "ahead",
  "ahead_by": 2,
  "behind_by": 0,
  "total_commits": 2,
  "commits": [
    {
      "sha": "1111111111111111111111111111111111111111",
      "commit": {
        "message": "feat: add lane (#1234)",
        "author": { "name": "Octo Cat" }
      },
      "author": { "login": "octocat" }
    },
    {
      "sha": "2222222222222222222222222222222222222222",
      "commit": {
        "message": "direct push commit",
        "author": { "name": "No Account" }
      },
      "author": null
    }
  ]
}`

func TestCompare(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/repos/smartcontractkit/chainlink-ccip/compare/aaa...bbb"
		if r.URL.Path != wantPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing/incorrect auth header")
		}
		fmt.Fprint(w, compareFixture)
	}))
	defer srv.Close()

	c := &ghClient{httpClient: srv.Client(), token: "test-token", baseURL: srv.URL}
	res, err := c.Compare(context.Background(), "smartcontractkit", "chainlink-ccip", "aaa", "bbb")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "ahead" || res.AheadBy != 2 {
		t.Errorf("unexpected result: %+v", res)
	}
	if len(res.Commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(res.Commits))
	}
	first := res.Commits[0]
	if first.Title != "feat: add lane" || first.PR != 1234 || first.Author != "@octocat" {
		t.Errorf("unexpected first commit: %+v", first)
	}
	second := res.Commits[1]
	if second.PR != 0 || second.Author != "No Account" {
		t.Errorf("unexpected second commit: %+v", second)
	}
}

func TestCompareHTTPError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"Not Found"}`)
	}))
	defer srv.Close()

	c := &ghClient{httpClient: srv.Client(), token: "", baseURL: srv.URL}
	_, err := c.Compare(context.Background(), "smartcontractkit", "chainlink-ccip", "aaa", "bbb")
	if err == nil {
		t.Fatal("expected error")
	}
}
