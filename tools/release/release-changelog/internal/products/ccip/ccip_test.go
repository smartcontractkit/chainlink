package ccip

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/tools/release/release-changelog/internal/engine"
)

func TestNormalizeRef(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		// CCIP image tags invert to their source git tags.
		"2.56.1-ccip-rc.2":       "v2.56.1-rc.2",
		"2.34.0-ccip-beta.0":     "v2.34.0-beta.0",
		"2.56.1-ccip-rc.2-arm64": "v2.56.1-rc.2",
		"2.56.1-ccip-rc.0-amd64": "v2.56.1-rc.0",
		// Hybrid: image tag with a "v" prefix (git-tag-shaped mistake).
		"v2.56.1-ccip-rc.1":       "v2.56.1-rc.1",
		"v2.56.1-ccip-rc.1-arm64": "v2.56.1-rc.1",
		// Full image URIs.
		"public.ecr.aws/chainlink/ccip:2.56.1-ccip-rc.2": "v2.56.1-rc.2",
		// Plain git refs pass through untouched.
		"v2.55.0":        "v2.55.0",
		"v2.56.1-rc.2":   "v2.56.1-rc.2",
		"release/2.57.2": "release/2.57.2",
		"eee47b86a6826870b0d8e92e8f979f3b8cce9134": "eee47b86a6826870b0d8e92e8f979f3b8cce9134",
	}
	for in, want := range cases {
		if got := normalizeRef(in); got != want {
			t.Errorf("normalizeRef(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestProductSanity guards the real CCIP config against accidental breakage:
// the primary pin of every tracked repo must be resolvable from a realistic
// go.mod + plugins.yaml pair, and the core-repo entry must stay path-scoped.
func TestProductSanity(t *testing.T) {
	t.Parallel()

	if Product.DisplayName != "CCIP" {
		t.Errorf("DisplayName = %q, want CCIP", Product.DisplayName)
	}
	if Product.NormalizeRef == nil {
		t.Error("NormalizeRef must be set")
	}

	const goMod = `module github.com/smartcontractkit/chainlink/v2

go 1.26.4

require (
	github.com/smartcontractkit/chainlink-ccip v0.1.1-solana.0.20260625091148-e5618f5682ee
	github.com/smartcontractkit/chainlink-ton v1.0.5-0.20260629213843-c52e07523035
)
`
	const pluginsYAML = `plugins:
  aptos:
    - moduleURI: "github.com/smartcontractkit/chainlink-aptos"
      gitRef: "v0.0.0-20260708114855-e953eeb028a7"
  evm:
    - moduleURI: "github.com/smartcontractkit/chainlink-evm"
      gitRef: "v0.3.4-0.20260715161014-611d8ac32364"
`
	mods, err := engine.ParseGoMod([]byte(goMod), Product.Repos)
	if err != nil {
		t.Fatal(err)
	}
	if mods["github.com/smartcontractkit/chainlink-ccip"] == "" {
		t.Error("chainlink-ccip module not tracked by Product.Repos")
	}
	if mods["github.com/smartcontractkit/chainlink-ton"] == "" {
		t.Error("chainlink-ton module not tracked by Product.Repos")
	}

	plugins, err := engine.ParsePluginsYAML([]byte(pluginsYAML), Product.Repos)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"aptos", "evm"} {
		if plugins[key].GitRef == "" {
			t.Errorf("plugin %q not tracked by Product.Repos", key)
		}
	}

	// The core repo entry must remain path-scoped (an unfiltered core
	// changelog duplicates the changesets CHANGELOG).
	var core *engine.RepoConfig
	for i := range Product.Repos {
		if Product.Repos[i].Local {
			core = &Product.Repos[i]
		}
	}
	if core == nil {
		t.Fatal("no Local (core repo) entry in Product.Repos")
	}
	if len(core.IncludePaths) == 0 {
		t.Error("core repo entry must keep IncludePaths set (e.g. core/capabilities/ccip/)")
	}
}
