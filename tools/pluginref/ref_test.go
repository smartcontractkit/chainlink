package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

func samplePluginsYAML() string {
	return `plugins:
  solana:
    - moduleURI: "github.com/smartcontractkit/chainlink-solana"
      gitRef: "v1.3.1-0.20260605202330-b5a89c32fdc1"
      installPath: "./pkg/solana/cmd/chainlink-solana"
  starknet:
    - moduleURI: "github.com/smartcontractkit/chainlink-starknet/relayer"
      gitRef: "73cd3d46ad0ce2871160369cb6447aeb9b48513f" # 2026-04-08
      installPath: "./pkg/chainlink/cmd/chainlink-starknet"
  aptos:
    - moduleURI: "github.com/smartcontractkit/chainlink-aptos"
      gitRef: "v0.0.0-20260609211101-71d38bd6a0a9"
      installPath: "./cmd/chainlink-aptos"
  tagged:
    - moduleURI: "github.com/example/repo"
      gitRef: "v1.2.3"
  prefixed:
    - moduleURI: "github.com/example/repo/sub"
      gitRef: "sub/v1.2.3"
`
}

func TestNormalizeVersion(t *testing.T) {
	t.Parallel()

	t.Run("plain tag", func(t *testing.T) {
		t.Parallel()
		mv := normalizeVersion("v1.2.3")
		if mv.Tag != "v1.2.3" || mv.SHA != "" || mv.TagPrefix != "" {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})

	t.Run("prefixed tag", func(t *testing.T) {
		t.Parallel()
		mv := normalizeVersion("sub/dir/v1.2.3")
		if mv.Tag != "v1.2.3" || mv.TagPrefix != "sub/dir" || mv.SHA != "" {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})

	t.Run("pseudo with sha", func(t *testing.T) {
		t.Parallel()
		valids := []struct {
			in   string
			want string
		}{
			{"v0.0.0-20260609211101-71d38bd6a0a9", "71d38bd6a0a9"},
			{"v1.3.1-0.20260605202330-b5a89c32fdc1", "b5a89c32fdc1"},
			{"v1.2.3-0.20250102030405-abcdef123456", "abcdef123456"},
		}
		for _, tc := range valids {
			t.Run(tc.in, func(t *testing.T) {
				t.Parallel()
				mv := normalizeVersion(tc.in)
				if mv.SHA != tc.want {
					t.Fatalf("normalizeVersion(%q).SHA = %q, want %q", tc.in, mv.SHA, tc.want)
				}
			})
		}
	})

	t.Run("raw sha", func(t *testing.T) {
		t.Parallel()
		const sha = "73cd3d46ad0ce2871160369cb6447aeb9b48513f"
		mv := normalizeVersion(sha)
		if mv.SHA != sha {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})
}

func TestCommitRefForGitRef(t *testing.T) {
	t.Parallel()

	cases := []struct {
		gitRef string
		want   string
	}{
		{"v1.3.1-0.20260605202330-b5a89c32fdc1", "b5a89c32fdc1"},
		{"v0.0.0-20260609211101-71d38bd6a0a9", "71d38bd6a0a9"},
		{"73cd3d46ad0ce2871160369cb6447aeb9b48513f", "73cd3d46ad0ce2871160369cb6447aeb9b48513f"},
		{"v1.2.3", "v1.2.3"},
		{"sub/v1.2.3", "sub/v1.2.3"},
	}

	for _, tc := range cases {
		t.Run(tc.gitRef, func(t *testing.T) {
			t.Parallel()
			got, err := commitRefForGitRef(tc.gitRef)
			if err != nil {
				t.Fatalf("commitRefForGitRef(%q): %v", tc.gitRef, err)
			}
			if got != tc.want {
				t.Fatalf("commitRefForGitRef(%q) = %q, want %q", tc.gitRef, got, tc.want)
			}
		})
	}
}

func TestCommitRefForGitRefErrors(t *testing.T) {
	t.Parallel()

	if _, err := commitRefForGitRef(""); err == nil {
		t.Fatal("expected error for empty gitRef")
	}
}

func TestGitRefForPlugin(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	cases := map[string]string{
		"solana":   "v1.3.1-0.20260605202330-b5a89c32fdc1",
		"starknet": "73cd3d46ad0ce2871160369cb6447aeb9b48513f",
		"aptos":    "v0.0.0-20260609211101-71d38bd6a0a9",
		"tagged":   "v1.2.3",
		"prefixed": "sub/v1.2.3",
	}

	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := gitRefForPlugin(path, name)
			if err != nil {
				t.Fatalf("gitRefForPlugin(%q): %v", name, err)
			}
			if got != want {
				t.Fatalf("gitRefForPlugin(%q) = %q, want %q", name, got, want)
			}
		})
	}
}

func TestGitRefForPluginErrors(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	t.Run("missing plugin", func(t *testing.T) {
		t.Parallel()
		if _, err := gitRefForPlugin(path, "missing"); err == nil {
			t.Fatal("expected error for missing plugin")
		}
	})

	t.Run("unknown plugin in empty file", func(t *testing.T) {
		t.Parallel()
		emptyPath := writeFile(t, dir, "empty.yaml", `plugins: {}`)
		if _, err := gitRefForPlugin(emptyPath, "solana"); err == nil {
			t.Fatal("expected error for unknown plugin in empty file")
		}
	})

	t.Run("plugin without gitRef", func(t *testing.T) {
		t.Parallel()
		noRefPath := writeFile(t, dir, "noref.yaml", `plugins:
  solana:
    - moduleURI: "github.com/smartcontractkit/chainlink-solana"
`)
		if _, err := gitRefForPlugin(noRefPath, "solana"); err == nil {
			t.Fatal("expected error for plugin without gitRef")
		}
	})
}
