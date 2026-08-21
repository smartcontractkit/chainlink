package changelog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// syntheticReport builds a representative report covering: plugin change,
// drift, rollback, keyword hits, divergence note, filtered repo, errors.
func syntheticReport() *Report {
	oldSnap := DepSnapshot{
		Ref: "v2.55.0", SHA: "oldoldoldoldoldoldoldoldoldoldoldold0000",
		Modules: map[string]string{
			"github.com/smartcontractkit/chainlink-ccip":            "v0.1.1-solana.0.20260101000000-aaaaaaaaaaaa",
			"github.com/smartcontractkit/chainlink-ccip/chains/evm": "v0.0.0-20260101000000-bbbbbbbbbbbb",
			"github.com/smartcontractkit/chainlink-ton":             "v1.0.5-0.20260101000000-111111111111",
			"github.com/smartcontractkit/chainlink-evm":             "v0.3.4-0.20260101000000-222222222222",
		},
		Plugins: map[string]PluginPin{
			"ton":    {GitRef: "v1.0.5-0.20260101000000-111111111111", ModuleURI: "github.com/smartcontractkit/chainlink-ton"},
			"evm":    {GitRef: "v0.3.4-0.20260101000000-222222222222", ModuleURI: "github.com/smartcontractkit/chainlink-evm"},
			"solana": {GitRef: "v1.3.1-0.20260101000000-333333333333", ModuleURI: "github.com/smartcontractkit/chainlink-solana"},
		},
	}
	newSnap := DepSnapshot{
		Ref: "v2.56.0", SHA: "newnewnewnewnewnewnewnewnewnewnewnew0000",
		Modules: map[string]string{
			"github.com/smartcontractkit/chainlink-ccip":            "v0.1.1-solana.0.20260202000000-cccccccccccc",
			"github.com/smartcontractkit/chainlink-ccip/chains/evm": "v0.0.0-20260202000000-dddddddddddd",
			"github.com/smartcontractkit/chainlink-ton":             "v1.0.5-0.20260202000000-444444444444",
			"github.com/smartcontractkit/chainlink-evm":             "v0.3.4-0.20260101000000-222222222222",
		},
		Plugins: map[string]PluginPin{
			"ton":    {GitRef: "v1.0.5-0.20260202000000-555555555555", ModuleURI: "github.com/smartcontractkit/chainlink-ton"},
			"evm":    {GitRef: "v0.3.4-0.20260101000000-222222222222", ModuleURI: "github.com/smartcontractkit/chainlink-evm"},
			"solana": {GitRef: "v1.3.1-0.20260101000000-333333333333", ModuleURI: "github.com/smartcontractkit/chainlink-solana"},
		},
	}

	ccipCfg := TrackedRepos[0]
	tonCfg := TrackedRepos[4]
	evmCfg := TrackedRepos[5]
	solCfg := TrackedRepos[3]
	coreCfg := TrackedRepos[6]

	rep := &Report{
		Old: oldSnap, New: newSnap,
		Generated: time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC),
	}

	ccip := RepoReport{
		Config: ccipCfg,
		Old:    repoPin{PrimaryVersion: oldSnap.Modules[ccipCfg.GoModules[0]], PrimarySHA: "aaaaaaaaaaaa"},
		New:    repoPin{PrimaryVersion: newSnap.Modules[ccipCfg.GoModules[0]], PrimarySHA: "cccccccccccc"},
		Status: "ahead", TotalInRange: 2,
		Commits: []CommitEntry{
			{SHA: "deadbeef0000", Title: "feat: add fast lane config", PR: 1234, Author: "@octocat"},
			{SHA: "feedface0000", Title: "chore: bump deps", PR: 1235, Author: "Jane Doe"},
		},
	}
	ccip.KeywordHits = []CommitEntry{ccip.Commits[0]} // "config" in title
	ccip.Notes = []string{"divergent pins (new ref): chainlink-ccip at `cccccccccccc`; chainlink-ccip/chains/evm at `dddddddddddd`"}

	ton := RepoReport{
		Config: tonCfg,
		Old:    repoPin{PrimaryVersion: "v1.0.5-0.20260101000000-111111111111", PrimarySHA: "111111111111"},
		New:    repoPin{PrimaryVersion: "v1.0.5-0.20260202000000-555555555555", PrimarySHA: "555555555555"},
		Status: "ahead", TotalInRange: 1,
		Commits: []CommitEntry{
			{SHA: "cafe0000cafe", Title: "hotfix for gas estimator", PR: 99, Author: "@tondev"},
		},
	}
	ton.KeywordHits = ton.Commits

	evm := RepoReport{
		Config: evmCfg, Status: "identical",
		Old: repoPin{PrimaryVersion: "v0.3.4-0.20260101000000-222222222222", PrimarySHA: "222222222222"},
		New: repoPin{PrimaryVersion: "v0.3.4-0.20260101000000-222222222222", PrimarySHA: "222222222222"},
	}

	sol := RepoReport{
		Config: solCfg, Status: "behind",
		Old: repoPin{PrimaryVersion: "v1.3.1-0.20260101000000-333333333333", PrimarySHA: "333333333333"},
		New: repoPin{PrimaryVersion: "v1.3.1-0.20260101000000-333333333333", PrimarySHA: "333333333333"},
		Err: "",
	}
	sol.Status = "behind"
	sol.Commits = nil

	core := RepoReport{
		Config: coreCfg, Status: "ahead",
		TotalInRange: 87,
		Commits: []CommitEntry{
			{SHA: "010101010101", Title: "feat(ccip): new capability", PR: 2000, Author: "@coredev"},
		},
	}

	broken := RepoReport{
		Config: TrackedRepos[1], // aptos
		Err:    "compare smartcontractkit/chainlink-aptos abc...def: HTTP 404: Not Found",
	}

	rep.Repos = []RepoReport{ccip, broken, {Config: TrackedRepos[2], Status: "identical"}, sol, ton, evm, core}

	for _, rr := range rep.Repos {
		rep.Flags = append(rep.Flags, repoFlags(rr, oldSnap, newSnap)...)
	}
	return rep
}

func TestRenderMarkdown_Golden(t *testing.T) {
	t.Parallel()

	got := RenderMarkdown(syntheticReport())
	goldenPath := filepath.Join("testdata", "report.golden.md")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file (run with UPDATE_GOLDEN=1 to create): %v", err)
	}
	if got != string(want) {
		t.Errorf("markdown mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRenderSlackSummary_Golden(t *testing.T) {
	t.Parallel()

	got := RenderSlackSummary(syntheticReport())
	goldenPath := filepath.Join("testdata", "slack-summary.golden.txt")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file (run with UPDATE_GOLDEN=1 to create): %v", err)
	}
	if got != string(want) {
		t.Errorf("slack summary mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
