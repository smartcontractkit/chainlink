package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/cmd"
)

// fakeRegistry serves a minimal Docker Registry v2 API: /token returns a bearer
// token and /v2/{repo}/manifests/{tag} returns 200 for known tags, 404 otherwise.
func fakeRegistry(t *testing.T, knownTags []string) string {
	t.Helper()
	known := make(map[string]bool, len(knownTags))
	for _, tg := range knownTags {
		known[tg] = true
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"token":"fake-bearer-token"}`)
	})
	mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/v2/")
		idx := strings.LastIndex(rest, "/manifests/")
		if idx < 0 {
			http.NotFound(w, r)
			return
		}
		tag := rest[idx+len("/manifests/"):]
		if known[tag] {
			fmt.Fprint(w, `{"schemaVersion":2}`)
			return
		}
		http.NotFound(w, r)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestCcipResolveBaseline_CLI_GitHubOutput(t *testing.T) {
	base := fakeRegistry(t, []string{"2.63.0-ccip-rc.0"})

	tmpDir := t.TempDir()
	ghOutput := filepath.Join(tmpDir, "gh_output")
	require.NoError(t, os.WriteFile(ghOutput, []byte{}, 0o600))
	t.Setenv("GITHUB_OUTPUT", ghOutput)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--chainlink-version", "v2.63.1-rc.4",
		"--registry-url", base,
		"--probe-retry-delay", "0s",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	ghContent, err := os.ReadFile(ghOutput)
	require.NoError(t, err)
	assert.Contains(t, string(ghContent), "ccip_pr_tag")
	assert.Contains(t, string(ghContent), "2.63.1-ccip-rc.4")
	assert.Contains(t, string(ghContent), "baseline_image_tag")
	assert.Contains(t, string(ghContent), "2.63.0-ccip-rc.0")
	assert.Contains(t, string(ghContent), "skip<<_GitHubActionsFileCommandDelimeter_\nfalse\n_GitHubActionsFileCommandDelimeter_\n") //typos:ignore `Delimeter` should be `Delimiter` // From GitHub

	// Human-readable summary on stdout.
	assert.Contains(t, out.String(), "Baseline for v2.63.1-rc.4: public.ecr.aws/chainlink/ccip:2.63.0-ccip-rc.0")
}

func TestCcipResolveBaseline_CLI_SkipWarning(t *testing.T) {
	// Empty registry: every candidate 404s.
	base := fakeRegistry(t, nil)

	tmpDir := t.TempDir()
	ghOutput := filepath.Join(tmpDir, "gh_output")
	require.NoError(t, os.WriteFile(ghOutput, []byte{}, 0o600))
	t.Setenv("GITHUB_OUTPUT", ghOutput)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--chainlink-version", "v2.63.0-rc.0",
		"--registry-url", base,
		"--max-fallback", "2",
		"--probe-retry-delay", "0s",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	ghContent, err := os.ReadFile(ghOutput)
	require.NoError(t, err)
	assert.Contains(t, string(ghContent), "baseline_image_tag<<")                                                                   // empty multiline value
	assert.Contains(t, string(ghContent), "skip<<_GitHubActionsFileCommandDelimeter_\ntrue\n_GitHubActionsFileCommandDelimeter_\n") //typos:ignore `Delimeter` should be `Delimiter` // From GitHub
	// Warning annotation emitted.
	assert.Contains(t, out.String(), "::warning::No published baseline CCIP rc.0 image found")
}

func TestCcipResolveBaseline_CLI_JSON(t *testing.T) {
	t.Parallel()
	base := fakeRegistry(t, []string{"2.63.0-ccip-rc.0"})

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--chainlink-version", "v2.63.1-rc.4",
		"--registry-url", base,
		"--probe-retry-delay", "0s",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &payload))
	assert.Equal(t, "2.63.1-ccip-rc.4", payload["ccip_pr_tag"])
	assert.Equal(t, "2.63.0-ccip-rc.0", payload["baseline_image_tag"])
	assert.Equal(t, false, payload["skip"])
}

func TestCcipResolveBaseline_CLI_LocalKV(t *testing.T) {
	t.Parallel()
	base := fakeRegistry(t, []string{"2.62.0-ccip-rc.0"})

	// No GITHUB_OUTPUT -> k=v lines on stdout for local execution.
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--chainlink-version", "v2.63.0",
		"--registry-url", base,
		"--probe-retry-delay", "0s",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	assert.Contains(t, out.String(), "ccip_pr_tag=2.63.0")
	assert.Contains(t, out.String(), "baseline_image_tag=2.62.0-ccip-rc.0")
	assert.Contains(t, out.String(), "skip=false")
}

func TestCcipResolveBaseline_CLI_EnvFallback(t *testing.T) {
	base := fakeRegistry(t, []string{"2.63.0-ccip-rc.0"})
	t.Setenv("CHAINLINK_VERSION", "v2.63.1-rc.4")
	t.Setenv("CCIP_IMAGE_REPO", "public.ecr.aws/chainlink/ccip")

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--registry-url", base,
		"--probe-retry-delay", "0s",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out.String(), "ccip_pr_tag=2.63.1-ccip-rc.4")
	assert.Contains(t, out.String(), "baseline_image_tag=2.63.0-ccip-rc.0")
}

func TestCcipResolveBaseline_CLI_RejectsNonVTag(t *testing.T) {
	t.Parallel()
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"ccip", "resolve-baseline",
		"--chainlink-version", "2.63.0-rc.0",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must begin with 'v'")
}
