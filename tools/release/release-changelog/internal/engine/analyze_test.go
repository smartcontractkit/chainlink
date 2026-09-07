package engine

import (
	"strings"
	"testing"
)

func TestKeywordPattern(t *testing.T) {
	t.Parallel()

	hits := []string{
		"fix!: correct CCIP nonce handling",
		"BREAKING: change config format",
		"Revert \"add lane config\"",
		"hotfix for token pool",
		"security patch for verifier",
		"update chain config defaults",
		"Config: bump lanes",
	}
	misses := []string{
		"feat: add new lane",
		"chore: bump deps",
		"reconfigure dashboards", // no word boundary match for "config"
		"fix offramp bug",        // no "!" suffix
		"docs: update README",
	}
	for _, h := range hits {
		if !keywordPattern.MatchString(h) {
			t.Errorf("expected keyword hit for %q", h)
		}
	}
	for _, m := range misses {
		if keywordPattern.MatchString(m) {
			t.Errorf("expected no keyword hit for %q", m)
		}
	}
}

func TestPathMatch(t *testing.T) {
	t.Parallel()

	includes := []string{"core/capabilities/ccip/"}
	cases := []struct {
		name     string
		files    []string
		includes []string
		excludes []string
		want     bool
	}{
		{"include hit", []string{"core/capabilities/ccip/ocr/plugin.go"}, includes, nil, true},
		{"include miss", []string{"core/services/relay/evm/evm.go"}, includes, nil, false},
		{"mixed files", []string{"README.md", "core/capabilities/ccip/x.go"}, includes, nil, true},
		{"no filters", []string{"anything.go"}, nil, nil, true},
		{"exclude only file", []string{"core/capabilities/ccip/gen/x.go"}, includes, []string{"core/capabilities/ccip/gen/"}, false},
		{"exclude partial", []string{"core/capabilities/ccip/gen/x.go", "core/capabilities/ccip/y.go"}, includes, []string{"core/capabilities/ccip/gen/"}, true},
		{"exclude only, no includes", []string{"docs/x.md"}, nil, []string{"docs/"}, false},
		{"empty files", nil, includes, nil, false},
	}
	for _, c := range cases {
		if got := pathMatch(c.files, c.includes, c.excludes); got != c.want {
			t.Errorf("%s: pathMatch = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestNoreplyLogin(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"12345+octocat@users.noreply.github.com": "octocat",
		"octocat@users.noreply.github.com":       "octocat",
		"real@example.com":                       "",
		"":                                       "",
	}
	for email, want := range cases {
		if got := noreplyLogin(email); got != want {
			t.Errorf("noreplyLogin(%q) = %q, want %q", email, got, want)
		}
	}
}

func TestDivergenceNotes(t *testing.T) {
	t.Parallel()

	cfg := RepoConfig{
		Name:       "chainlink-ton",
		GoModules:  []string{"github.com/smartcontractkit/chainlink-ton"},
		PluginKeys: []string{"ton"},
	}
	// Convergent: same SHA everywhere -> no note.
	convergent := repoPin{
		ModuleVersions: map[string]string{"github.com/smartcontractkit/chainlink-ton": "v1.0.5-0.20260629213843-c52e07523035"},
		PluginRefs:     map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260629213843-c52e07523035"}},
	}
	if notes := divergenceNotes(cfg, convergent, "new"); len(notes) != 0 {
		t.Errorf("convergent pins produced notes: %v", notes)
	}
	// Divergent: plugin and module disagree -> one note naming both.
	divergent := repoPin{
		ModuleVersions: map[string]string{"github.com/smartcontractkit/chainlink-ton": "v1.0.5-0.20260629213843-c52e07523035"},
		PluginRefs:     map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260701000000-aaaaaaaaaaaa"}},
	}
	notes := divergenceNotes(cfg, divergent, "new")
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %v", notes)
	}
	if !strings.Contains(notes[0], "c52e07523035") || !strings.Contains(notes[0], "aaaaaaaaaaaa") {
		t.Errorf("note missing SHAs: %s", notes[0])
	}
}

func TestRepoFlags(t *testing.T) {
	t.Parallel()

	cfg := RepoConfig{
		Name:       "chainlink-ton",
		Owner:      "smartcontractkit",
		GoModules:  []string{"github.com/smartcontractkit/chainlink-ton"},
		PluginKeys: []string{"ton"},
	}
	mod := "github.com/smartcontractkit/chainlink-ton"

	t.Run("plugin gitRef changed", func(t *testing.T) {
		t.Parallel()

		oldSnap := DepSnapshot{Plugins: map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260629213843-c52e07523035", ModuleURI: mod}}}
		newSnap := DepSnapshot{Plugins: map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260701000000-aaaaaaaaaaaa", ModuleURI: mod}}}
		rr := RepoReport{Config: cfg}
		flags := repoFlags(rr, oldSnap, newSnap)
		if len(flags) != 1 || !strings.Contains(flags[0], "gitRef changed") {
			t.Errorf("unexpected flags: %v", flags)
		}
	})

	t.Run("plugin removed and added", func(t *testing.T) {
		t.Parallel()

		oldSnap := DepSnapshot{Plugins: map[string]PluginPin{"ton": {GitRef: "v1", ModuleURI: mod}}}
		newSnap := DepSnapshot{Plugins: map[string]PluginPin{}}
		flags := repoFlags(RepoReport{Config: cfg}, oldSnap, newSnap)
		if len(flags) != 1 || !strings.Contains(flags[0], "REMOVED") {
			t.Errorf("unexpected flags: %v", flags)
		}
		flags = repoFlags(RepoReport{Config: cfg}, newSnap, oldSnap)
		if len(flags) != 1 || !strings.Contains(flags[0], "ADDED") {
			t.Errorf("unexpected flags: %v", flags)
		}
	})

	t.Run("dual-source drift", func(t *testing.T) {
		t.Parallel()

		snap := DepSnapshot{
			Modules: map[string]string{mod: "v1.0.5-0.20260629213843-c52e07523035"},
			Plugins: map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260701000000-bbbbbbbbbbbb", ModuleURI: mod}},
		}
		flags := repoFlags(RepoReport{Config: cfg}, snap, snap)
		if len(flags) != 2 { // same drift flagged at both refs
			t.Fatalf("expected 2 drift flags, got %v", flags)
		}
		if !strings.Contains(flags[0], "DRIFT") {
			t.Errorf("unexpected flag text: %s", flags[0])
		}
	})

	t.Run("no drift when SHAs match", func(t *testing.T) {
		t.Parallel()

		snap := DepSnapshot{
			Modules: map[string]string{mod: "v1.0.5-0.20260629213843-c52e07523035"},
			Plugins: map[string]PluginPin{"ton": {GitRef: "v1.0.5-0.20260629213843-c52e07523035", ModuleURI: mod}},
		}
		if flags := repoFlags(RepoReport{Config: cfg}, snap, snap); len(flags) != 0 {
			t.Errorf("unexpected flags: %v", flags)
		}
	})

	t.Run("rollback", func(t *testing.T) {
		t.Parallel()

		rr := RepoReport{
			Config: cfg,
			Status: "behind",
			Old:    repoPin{PrimarySHA: "aaaaaaaaaaaa"},
			New:    repoPin{PrimarySHA: "bbbbbbbbbbbb"},
		}
		flags := repoFlags(rr, DepSnapshot{}, DepSnapshot{})
		if len(flags) != 1 || !strings.Contains(flags[0], "ROLLBACK") {
			t.Errorf("unexpected flags: %v", flags)
		}
	})

	t.Run("keyword hits", func(t *testing.T) {
		t.Parallel()

		rr := RepoReport{
			Config: cfg,
			KeywordHits: []CommitEntry{
				{SHA: "abc123def456", Title: "hotfix for token pool"},
			},
		}
		flags := repoFlags(rr, DepSnapshot{}, DepSnapshot{})
		if len(flags) != 1 || !strings.Contains(flags[0], "keyword match") {
			t.Errorf("unexpected flags: %v", flags)
		}
	})
}
