package ccip_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ccip"
)

// fakeProber publishes exactly the tags in its set; everything else is
// definitively absent (false, nil). Optional failTags force a transient error
// for the first failCount attempts, to exercise retry logic. Resolve probes
// candidates sequentially, so the call counter needs no synchronization.
type fakeProber struct {
	published  map[string]bool
	failTags   map[string]bool
	failCount  int
	callsByTag map[string]int
}

func newFakeProber(tags ...string) *fakeProber {
	m := make(map[string]bool, len(tags))
	for _, t := range tags {
		m[t] = true
	}
	return &fakeProber{published: m, callsByTag: map[string]int{}}
}

func (p *fakeProber) IsPublished(_ context.Context, tag string) (bool, error) {
	p.callsByTag[tag]++
	n := p.callsByTag[tag]
	if p.failTags != nil && p.failTags[tag] && n <= p.failCount {
		return false, fmt.Errorf("transient error for %s (attempt %d)", tag, n)
	}
	return p.published[tag], nil
}

// --- DerivePRTag ---

func TestDerivePRTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"rcN", "v2.63.1-rc.4", "2.63.1-ccip-rc.4"},
		{"rc0 patch", "v2.63.1-rc.0", "2.63.1-ccip-rc.0"},
		{"rc0 minor", "v2.63.0-rc.0", "2.63.0-ccip-rc.0"},
		{"stable no ccip suffix", "v2.63.0", "2.63.0"},
		{"beta", "v2.63.0-beta.0", "2.63.0-ccip-beta.0"},
		{"double-digit rc", "v2.63.0-rc.12", "2.63.0-ccip-rc.12"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ccip.DerivePRTag(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDerivePRTag_Errors(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		_, err := ccip.DerivePRTag("")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be set")
	})
	t.Run("non-v", func(t *testing.T) {
		t.Parallel()
		_, err := ccip.DerivePRTag("2.63.0-rc.0")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must begin with 'v'")
	})
}

// --- Candidates ---

func TestCandidates(t *testing.T) {
	t.Parallel()
	// maxFallback capped small to keep assertions short.
	tests := []struct {
		name  string
		in    string
		maxFB int
		want  []string
	}{
		{
			name:  "rcN prepends same-version rc0 then minor base rc0s",
			in:    "v2.63.1-rc.4",
			maxFB: 2,
			want: []string{
				"2.63.1-ccip-rc.0",
				"2.63.0-ccip-rc.0",
				"2.62.0-ccip-rc.0",
				"2.61.0-ccip-rc.0",
			},
		},
		{
			name:  "rc0 patch uses this minor .0 rc0 then prior",
			in:    "v2.63.1-rc.0",
			maxFB: 2,
			want: []string{
				"2.63.0-ccip-rc.0",
				"2.62.0-ccip-rc.0",
				"2.61.0-ccip-rc.0",
			},
		},
		{
			name:  "rc0 minor uses previous minor .0 rc0 then prior",
			in:    "v2.63.0-rc.0",
			maxFB: 2,
			want: []string{
				"2.62.0-ccip-rc.0",
				"2.61.0-ccip-rc.0",
				"2.60.0-ccip-rc.0",
			},
		},
		{
			name:  "stable uses previous minor .0 rc0",
			in:    "v2.63.0",
			maxFB: 1,
			want: []string{
				"2.62.0-ccip-rc.0",
				"2.61.0-ccip-rc.0",
			},
		},
		{
			name:  "beta uses previous minor .0 rc0",
			in:    "v2.63.0-beta.0",
			maxFB: 1,
			want: []string{
				"2.62.0-ccip-rc.0",
				"2.61.0-ccip-rc.0",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ccip.Candidates(tc.in, tc.maxFB)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCandidates_StopsAtZeroMinor(t *testing.T) {
	t.Parallel()
	// v2.0.0-rc.0 -> baseMinor = -1, so the walk adds nothing.
	got, err := ccip.Candidates("v2.0.0-rc.0", 12)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestCandidates_ParseErrors(t *testing.T) {
	t.Parallel()
	t.Run("non-v", func(t *testing.T) {
		t.Parallel()
		_, err := ccip.Candidates("2.63.0-rc.0", 12)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must begin with 'v'")
	})
	t.Run("bad core", func(t *testing.T) {
		t.Parallel()
		_, err := ccip.Candidates("v2.63.x-rc.0", 12)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse version core")
	})
	t.Run("bad rc number", func(t *testing.T) {
		t.Parallel()
		_, err := ccip.Candidates("v2.63.0-rc.xx", 12)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse rc number")
	})
}

// --- Resolve (matches every bash stubbed test case) ---

func TestResolve(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		version   string
		published []string // tags the fake registry reports as published
		maxFB     int
		wantPR    string
		wantBase  string
		wantSkip  bool
	}{
		{
			name:      "rcN uses same-version rc0",
			version:   "v2.63.1-rc.4",
			published: []string{"2.63.1-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.1-ccip-rc.4",
			wantBase:  "2.63.1-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "rcN falls back to minor base rc0",
			version:   "v2.63.1-rc.4",
			published: []string{"2.63.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.1-ccip-rc.4",
			wantBase:  "2.63.0-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "rc0 patch uses minor base rc0",
			version:   "v2.63.1-rc.0",
			published: []string{"2.63.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.1-ccip-rc.0",
			wantBase:  "2.63.0-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "rc0 minor uses previous minor rc0",
			version:   "v2.63.0-rc.0",
			published: []string{"2.62.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.0-ccip-rc.0",
			wantBase:  "2.62.0-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "walks past missing minor rc0s",
			version:   "v2.63.0-rc.0",
			published: []string{"2.60.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.0-ccip-rc.0",
			wantBase:  "2.60.0-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "skips when nothing published",
			version:   "v2.63.0-rc.0",
			published: nil,
			maxFB:     12,
			wantPR:    "2.63.0-ccip-rc.0",
			wantBase:  "",
			wantSkip:  true,
		},
		{
			name:      "stable derives no ccip suffix",
			version:   "v2.63.0",
			published: []string{"2.62.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.0",
			wantBase:  "2.62.0-ccip-rc.0",
			wantSkip:  false,
		},
		{
			name:      "beta derives ccip beta tag",
			version:   "v2.63.0-beta.0",
			published: []string{"2.62.0-ccip-rc.0"},
			maxFB:     12,
			wantPR:    "2.63.0-ccip-beta.0",
			wantBase:  "2.62.0-ccip-rc.0",
			wantSkip:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prober := newFakeProber(tc.published...)
			res, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
				ChainlinkVersion: tc.version,
				MaxFallback:      tc.maxFB,
				ProbeAttempts:    1,
				ProbeRetryDelay:  0,
			}, prober)
			require.NoError(t, err)
			assert.Equal(t, tc.wantPR, res.PRTag)
			assert.Equal(t, tc.wantBase, res.BaselineTag)
			assert.Equal(t, tc.wantSkip, res.Skip)
		})
	}
}

func TestResolve_RetriesTransientThenSucceeds(t *testing.T) {
	t.Parallel()
	// First probe of the winning tag errors; the retry succeeds.
	prober := newFakeProber("2.63.0-ccip-rc.0")
	prober.failTags = map[string]bool{"2.63.0-ccip-rc.0": true}
	prober.failCount = 1 // fail the 1st attempt only

	res, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
		ChainlinkVersion: "v2.63.1-rc.0",
		MaxFallback:      12,
		ProbeAttempts:    3,
		ProbeRetryDelay:  1 * time.Millisecond,
	}, prober)
	require.NoError(t, err)
	assert.False(t, res.Skip)
	assert.Equal(t, "2.63.0-ccip-rc.0", res.BaselineTag)
	// Two calls to the winning tag: one failed, one succeeded.
	assert.Equal(t, 2, prober.callsByTag["2.63.0-ccip-rc.0"])
}

func TestResolve_RejectsNonVTag(t *testing.T) {
	t.Parallel()
	_, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
		ChainlinkVersion: "2.63.0-rc.0",
		ProbeAttempts:    1,
	}, newFakeProber())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must begin with 'v'")
}

// --- RegistryProbe against an httptest fake registry ---

// fakeRegistry serves a Docker Registry v2-ish API: /token returns a bearer
// token and /v2/{repo}/manifests/{tag} returns 200 for known tags, 404
// otherwise (and 500 for tags in serverErrorTags to exercise the transient path).
func fakeRegistry(t *testing.T, knownTags []string, serverErrorTags []string) (*httptest.Server, string) {
	t.Helper()
	known := make(map[string]bool, len(knownTags))
	for _, tg := range knownTags {
		known[tg] = true
	}
	bad := make(map[string]bool, len(serverErrorTags))
	for _, tg := range serverErrorTags {
		bad[tg] = true
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"token":"fake-bearer-token"}`)
	})
	mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		// path: /v2/{repo}/manifests/{tag}
		rest := strings.TrimPrefix(r.URL.Path, "/v2/")
		idx := strings.LastIndex(rest, "/manifests/")
		if idx < 0 {
			http.NotFound(w, r)
			return
		}
		tag := rest[idx+len("/manifests/"):]
		switch {
		case bad[tag]:
			http.Error(w, "server error", http.StatusInternalServerError)
		case known[tag]:
			w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
			fmt.Fprint(w, `{"schemaVersion":2}`)
		default:
			http.NotFound(w, r)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, srv.URL
}

func TestRegistryProbe_IsPublished(t *testing.T) {
	t.Parallel()
	_, base := fakeRegistry(t, []string{"2.63.0-ccip-rc.0"}, nil)
	probe := ccip.NewRegistryProbe(base, "chainlink/ccip", nil)

	got, err := probe.IsPublished(context.Background(), "2.63.0-ccip-rc.0")
	require.NoError(t, err)
	assert.True(t, got)

	got, err = probe.IsPublished(context.Background(), "9.9.9-ccip-rc.0")
	require.NoError(t, err) // definitively absent: 404, no error
	assert.False(t, got)
}

func TestRegistryProbe_ServerErrorIsTransient(t *testing.T) {
	t.Parallel()
	_, base := fakeRegistry(t, nil, []string{"bad.0-ccip-rc.0"})
	probe := ccip.NewRegistryProbe(base, "chainlink/ccip", nil)

	got, err := probe.IsPublished(context.Background(), "bad.0-ccip-rc.0")
	require.Error(t, err)
	assert.False(t, got)
	assert.Contains(t, err.Error(), "500")
}

func TestResolve_AgainstHTTPFakeRegistry(t *testing.T) {
	t.Parallel()
	_, base := fakeRegistry(t, []string{"2.63.0-ccip-rc.0"}, nil)
	probe := ccip.NewRegistryProbe(base, "chainlink/ccip", nil)

	res, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
		ChainlinkVersion: "v2.63.1-rc.4",
		MaxFallback:      12,
		ProbeAttempts:    1,
	}, probe)
	require.NoError(t, err)
	assert.Equal(t, "2.63.1-ccip-rc.4", res.PRTag)
	// same-version rc.0 is absent in the fake -> falls back to 2.63.0-ccip-rc.0.
	assert.Equal(t, "2.63.0-ccip-rc.0", res.BaselineTag)
	assert.False(t, res.Skip)
}

func TestResolve_AgainstHTTPFakeRegistry_SkipWhenAbsent(t *testing.T) {
	t.Parallel()
	// Empty registry: every candidate 404s.
	_, base := fakeRegistry(t, nil, nil)
	probe := ccip.NewRegistryProbe(base, "chainlink/ccip", nil)

	res, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
		ChainlinkVersion: "v2.63.0-rc.0",
		MaxFallback:      3,
		ProbeAttempts:    1,
	}, probe)
	require.NoError(t, err)
	assert.True(t, res.Skip)
	assert.Empty(t, res.BaselineTag)
}

// TestResolve_LivePublicECR validates the real public-ECR token + manifest
// endpoints end-to-end. Skipped unless CCIP_LIVE_ECR=1 (run manually / in a
// dedicated job); mirrors the former bash harness's LIVE_ECR=1 cases.
func TestResolve_LivePublicECR(t *testing.T) {
	t.Parallel()
	if os.Getenv("CCIP_LIVE_ECR") != "1" {
		t.Skip("skipping live public-ECR test; set CCIP_LIVE_ECR=1 to run")
	}
	probe := ccip.NewRegistryProbe("https://public.ecr.aws", "chainlink/ccip", nil)

	// A definitively-absent tag must return (false, nil), not error.
	got, err := probe.IsPublished(context.Background(), "0.0.0-ccip-rc.0")
	require.NoError(t, err)
	assert.False(t, got)

	// Across recent rc.0 versions at least one baseline should resolve against
	// the live registry (and every probe must complete without error).
	resolvedAny := false
	for _, v := range []string{"v2.63.0-rc.0", "v2.62.0-rc.0", "v2.61.0-rc.0"} {
		res, err := ccip.Resolve(context.Background(), ccip.ResolveOptions{
			ChainlinkVersion: v,
			MaxFallback:      12,
			ProbeAttempts:    3,
			ProbeRetryDelay:  2 * time.Second,
		}, probe)
		require.NoError(t, err)
		assert.Regexp(t, `^[0-9].*-ccip-rc\.0$`, res.PRTag)
		if !res.Skip {
			assert.Regexp(t, `^[0-9].*-ccip-rc\.0$`, res.BaselineTag)
			resolvedAny = true
		}
	}
	assert.True(t, resolvedAny, "expected at least one recent rc.0 to resolve a live baseline")
}
