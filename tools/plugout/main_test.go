package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// -----------------------------
// helpers
// -----------------------------

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// Basic plugin YAML with indentation and comments that should be preserved.
func samplePluginsYAML() string {
	return `plugins:
  example:
    - moduleURI: github.com/example/repo
      installPath: ""
      gitRef: "v1.2.3" # keep-comment-root
  exampleTwo:
    - moduleURI: github.com/example/repo/sub
      libs: ["x","y"]
      gitRef: "sub/v1.2.3"
  misc:
    - moduleURI: github.com/foo/bar
      gitRef: "v0.0.0-20250102030405-abcdef123456"
`
}

// Same structure but without a gitRef under a module (to hit an error path).
func yamlMissingGitRefFor(module string) string {
	return `plugins:
  data:
    - moduleURI: ` + module + `
      installPath: ""
      # no gitRef here, update should fail
`
}

// -----------------------------
// unit tests: small helpers
// -----------------------------

func TestContains(t *testing.T) {
	if !contains([]string{"a", "b"}, "a") {
		t.Fatal("expected contains to find element")
	}
	if contains([]string{"a", "b"}, "c") {
		t.Fatal("expected contains to NOT find element")
	}
}

func TestModuleSubdir(t *testing.T) {
	cases := map[string]string{
		"github.com/org/repo":         "",
		"github.com/org/repo/relayer": "relayer",
		"github.com/org/repo/sub/dir": "sub/dir",
		"github.com/org/repo/v2":      "v2",
		"github.com/org/repo/v2/sub":  "v2/sub",
		"not/a/valid/module/at/all":   "module/at/all", // function drops first 3 segments
		"g.com/o/r/s":                 "s",
		"g.com/o/r":                   "",
	}
	for in, want := range cases {
		if got := moduleSubdir(in); got != want {
			t.Errorf("moduleSubdir(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestShaEqual(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"abcdef1", "abcdef1", true},
		{"abcdef1234", "abcdef1", true}, // prefix ok
		{"abcdef1", "abcdef1234", true}, // prefix ok
		{"abcdef1", "", false},
		{"", "abcdef1", false},
		{"abc", "def", false},
	}
	for _, c := range cases {
		if got := shaEqual(c.a, c.b); got != c.want {
			t.Errorf("shaEqual(%q,%q)=%v want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestNormalizeVersion(t *testing.T) {
	t.Run("plain tag", func(t *testing.T) {
		mv := normalizeVersion("v1.2.3")
		if mv.Tag != "v1.2.3" || mv.SHA != "" || mv.TagPrefix != "" || mv.Raw != "v1.2.3" {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})

	t.Run("prefixed tag", func(t *testing.T) {
		mv := normalizeVersion("sub/dir/v1.2.3")
		if mv.Raw != "sub/dir/v1.2.3" || mv.Tag != "v1.2.3" || mv.TagPrefix != "sub/dir" || mv.SHA != "" {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})

	t.Run("pseudo with sha (go style)", func(t *testing.T) {
		valids := []string{
			"v0.0.0-20251013133428-62ab1091a563",
			"v1.2.3-0.20250102030405-abcdef123456",
			"v1.2.3-rc.0.20240102030405-deadbeefcafebabe",
		}
		for _, v := range valids {
			mv := normalizeVersion(v)
			if mv.SHA == "" || mv.Tag != "" || mv.Raw != v {
				t.Fatalf("unexpected mv for valid pseudo %q: %+v", v, mv)
			}
		}
	})

	t.Run("bad pseudos should not match", func(t *testing.T) {
		invalids := []string{
			"v1.2.3--20240102030405-deadbeef",   // extra hyphen
			"v1.2.3-20240102030405g-deadbeef",   // junk in timestamp
			"v1.2.3-0-20240102030405-deadbeef",  // missing dot after 0
			"v1.2.3-rc.20240102030405-deadbeef", // missing .0 before timestamp for pre-release form
		}
		for _, inv := range invalids {
			mv := normalizeVersion(inv)
			if mv.SHA != "" || mv.Raw != inv {
				t.Fatalf("expected no match for invalid pseudo version %q, got %+v", inv, mv)
			}
		}
	})

	t.Run("raw sha", func(t *testing.T) {
		mv := normalizeVersion("abcdef1234567890")
		if mv.SHA != "abcdef1234567890" || mv.Tag != "" || mv.Raw != "abcdef1234567890" {
			t.Fatalf("unexpected mv: %+v", mv)
		}
	})
}

func TestDesiredYAMLRefForModule(t *testing.T) {
	cases := []struct {
		module string
		mv     ModuleVersion
		want   string
	}{
		// Root module: function returns mv.Raw when Tag present but no subdir.
		{"github.com/example/repo", ModuleVersion{Tag: "v1.2.4", Raw: "v1.2.4"}, "v1.2.4"},
		// Submodule: write "sub/.../vX.Y.Z"
		{"github.com/example/repo/sub", ModuleVersion{Tag: "v1.2.4", TagPrefix: "sub", Raw: "sub/v1.2.4"}, "sub/v1.2.4"},
		{"github.com/example/repo/sub/dir", ModuleVersion{Tag: "v1.2.4", TagPrefix: "sub/dir", Raw: "sub/dir/v1.2.4"}, "sub/dir/v1.2.4"},
		// Pseudo/other raw forms fall back to Raw
		{"github.com/example/repo", ModuleVersion{Raw: "v0.0.0-2025...-deadbeef"}, "v0.0.0-2025...-deadbeef"},
	}
	for _, c := range cases {
		got := desiredYAMLRefForModule(c.module, c.mv)
		if got != c.want {
			t.Errorf("desiredYAMLRefForModule(%q,%+v)=%q want %q", c.module, c.mv, got, c.want)
		}
	}
}

func TestTagsMatchWithSubdir(t *testing.T) {
	// root modules: tag equality only
	if !tagsMatchWithSubdir("github.com/example/repo",
		ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"},
		ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"}) {
		t.Fatal("expected root tag equality")
	}

	// sub modules: tag with prefix equivalence against raw "sub/vX.Y.Z"
	if !tagsMatchWithSubdir("github.com/example/repo/sub",
		ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"},
		ModuleVersion{Raw: "sub/v1.2.3"}) {
		t.Fatal("expected subdir tag equivalence against raw")
	}

	// one side prefixed TagPrefix
	if !tagsMatchWithSubdir("github.com/example/repo/sub",
		ModuleVersion{Tag: "v1.2.3", TagPrefix: "sub", Raw: "sub/v1.2.3"},
		ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"}) {
		t.Fatal("expected tag equality when one side prefixed")
	}

	// v2 path with deeper subdir
	if !tagsMatchWithSubdir("github.com/example/repo/v2/relayer",
		ModuleVersion{Tag: "v2.0.1", Raw: "v2.0.1"},
		ModuleVersion{Raw: "v2/relayer/v2.0.1"}) {
		t.Fatal("expected v2 subdir tag equivalence")
	}

	// mismatch
	if tagsMatchWithSubdir("github.com/example/repo/sub",
		ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"},
		ModuleVersion{Tag: "v1.2.4", Raw: "v1.2.4"}) {
		t.Fatal("did not expect tag equivalence")
	}
}

func TestVersionsMatchForModule(t *testing.T) {
	module := "github.com/example/repo/sub"

	t.Run("SHA equality (prefix)", func(t *testing.T) {
		a := ModuleVersion{SHA: "abcdef1234", Raw: "v0.0.0-...-abcdef1234"}
		b := ModuleVersion{SHA: "abcdef", Raw: "abcdef"}
		if !versionsMatchForModule(module, a, b) {
			t.Fatal("expected SHA prefix equality to match")
		}
	})

	t.Run("YAML raw contains go.mod SHA", func(t *testing.T) {
		a := ModuleVersion{SHA: "deadbeef", Raw: "v0.0.0-...-deadbeef"}
		b := ModuleVersion{Raw: "v0.0.0-20250102030405-deadbeef"}
		if !versionsMatchForModule(module, a, b) {
			t.Fatal("expected match when YAML raw contains SHA")
		}
	})

	t.Run("Tag with subdir equivalence", func(t *testing.T) {
		a := ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"}
		b := ModuleVersion{Raw: "sub/v1.2.3"}
		if !versionsMatchForModule(module, a, b) {
			t.Fatal("expected tag/subdir equivalence")
		}
	})

	t.Run("Raw equality fallback", func(t *testing.T) {
		a := ModuleVersion{Raw: "weird-form-1"}
		b := ModuleVersion{Raw: "weird-form-1"}
		if !versionsMatchForModule(module, a, b) {
			t.Fatal("expected raw fallback equality")
		}
	})
}

// -----------------------------
// YAML discovery/update
// -----------------------------

func TestDiscoverPluginVersions(t *testing.T) {
	dir := t.TempDir()
	yamlPath := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	got, err := discoverPluginVersions(yamlPath)
	if err != nil {
		t.Fatalf("discoverPluginVersions error: %v", err)
	}

	want := map[string]string{
		"github.com/example/repo":     "v1.2.3",
		"github.com/example/repo/sub": "sub/v1.2.3",
		"github.com/foo/bar":          "v0.0.0-20250102030405-abcdef123456",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("discoverPluginVersions got=%v want=%v", got, want)
	}
}

func TestUpdateGitRefInYAML_SuccessAndPreserveFormatting(t *testing.T) {
	dir := t.TempDir()
	yamlPath := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	// Update root module to new tag (Raw must be populated)
	err := updateGitRefInYAML(yamlPath, "github.com/example/repo", ModuleVersion{Tag: "v1.2.4", Raw: "v1.2.4"})
	if err != nil {
		t.Fatalf("updateGitRefInYAML root: %v", err)
	}

	// Update submodule; ensure "sub/vX.Y.Z"
	err = updateGitRefInYAML(yamlPath, "github.com/example/repo/sub", ModuleVersion{Tag: "v1.2.5", Raw: "v1.2.5"})
	if err != nil {
		t.Fatalf("updateGitRefInYAML sub: %v", err)
	}

	content := readFile(t, yamlPath)

	// Quotes preserved and comment kept
	if !strings.Contains(content, `gitRef: "v1.2.4" # keep-comment-root`) {
		t.Fatalf("expected updated quoted gitRef with preserved comment, got:\n%s", content)
	}
	if !strings.Contains(content, `gitRef: "sub/v1.2.5"`) {
		t.Fatalf("expected updated subdir tag, got:\n%s", content)
	}
}

func TestUpdateGitRefInYAML_ModuleNotFound(t *testing.T) {
	dir := t.TempDir()
	yamlPath := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())
	err := updateGitRefInYAML(yamlPath, "github.com/does/not/exist", ModuleVersion{Tag: "v0.1.0", Raw: "exist/v0.1.0"})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got: %v", err)
	}
}

func TestUpdateGitRefInYAML_NoGitRefLineToReplace(t *testing.T) {
	dir := t.TempDir()
	mod := "github.com/example/repo"
	yamlPath := writeFile(t, dir, "plugins.yaml", yamlMissingGitRefFor(mod))
	err := updateGitRefInYAML(yamlPath, mod, ModuleVersion{Tag: "v1.2.3", Raw: "v1.2.3"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "failed to update gitref") {
		t.Fatalf("expected failure to update gitRef, got: %v", err)
	}
}

// -----------------------------
// runSync end-to-end with seam
// -----------------------------

func TestRunSync_CheckMode_WithMismatch(t *testing.T) {
	dir := t.TempDir()
	// go.mod only needs to exist
	goMod := writeFile(t, dir, "go.mod", "module github.com/example/repo\n")
	plugins := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	opts := Options{
		GoModPath:     goMod,
		PluginPaths:   []string{plugins},
		IgnoreModules: nil,
		Update:        false,
		GetModVersion: func(goModPath, module string) (ModuleVersion, error) {
			// Force mismatch for two modules
			switch module {
			case "github.com/example/repo":
				return ModuleVersion{Tag: "v1.2.4", Raw: "v1.2.4"}, nil // YAML has v1.2.3 -> mismatch
			case "github.com/example/repo/sub":
				return ModuleVersion{Tag: "v1.2.3", TagPrefix: "sub", Raw: "sub/v1.2.3"}, nil // YAML has sub/v1.2.3 -> match via subdir rule
			case "github.com/foo/bar":
				return ModuleVersion{SHA: "deadbeef", Raw: "deadbeef"}, nil // YAML has pseudo with different SHA -> mismatch
			default:
				return ModuleVersion{}, nil
			}
		},
	}

	hasMismatch, err := runSync(opts)
	if err != nil {
		t.Fatalf("runSync error: %v", err)
	}
	if !hasMismatch {
		t.Fatalf("expected mismatches in CHECK mode")
	}
}

func TestRunSync_UpdateMode_AppliesChanges(t *testing.T) {
	dir := t.TempDir()
	goMod := writeFile(t, dir, "go.mod", "module github.com/example/repo\n")
	plugins := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	opts := Options{
		GoModPath:   goMod,
		PluginPaths: []string{plugins},
		Update:      true,
		GetModVersion: func(goModPath, module string) (ModuleVersion, error) {
			switch module {
			case "github.com/example/repo":
				return ModuleVersion{Tag: "v1.2.4", Raw: "v1.2.4"}, nil
			case "github.com/example/repo/sub":
				return ModuleVersion{Tag: "v1.2.5", TagPrefix: "sub", Raw: "sub/v1.2.5"}, nil
			case "github.com/foo/bar":
				// Simulate same pseudo version -> no change
				return normalizeVersion("v0.0.0-20250102030405-abcdef123456"), nil
			default:
				return ModuleVersion{}, nil
			}
		},
	}

	hasMismatch, err := runSync(opts)
	if err != nil {
		t.Fatalf("runSync error: %v", err)
	}
	if hasMismatch {
		t.Fatalf("did not expect mismatches after UPDATE")
	}

	updated := readFile(t, plugins)
	if !strings.Contains(updated, `gitRef: "v1.2.4"`) {
		t.Fatalf("root module not updated:\n%s", updated)
	}
	if !strings.Contains(updated, `gitRef: "sub/v1.2.5"`) {
		t.Fatalf("sub module not updated:\n%s", updated)
	}
}

func TestRunSync_IgnoreModules_SkipsChecks(t *testing.T) {
	dir := t.TempDir()
	goMod := writeFile(t, dir, "go.mod", "module github.com/example/repo\n")
	plugins := writeFile(t, dir, "plugins.yaml", samplePluginsYAML())

	opts := Options{
		GoModPath:     goMod,
		PluginPaths:   []string{plugins},
		IgnoreModules: []string{"github.com/example/repo", "github.com/foo/bar"},
		Update:        false,
		GetModVersion: func(goModPath, module string) (ModuleVersion, error) {
			// Would mismatch the 2 ignored modules
			return ModuleVersion{Tag: "v1.2.3", TagPrefix: "sub", Raw: "sub/v1.2.3"}, nil
		},
	}

	hasMismatch, err := runSync(opts)
	if err != nil {
		t.Fatalf("runSync error: %v", err)
	}
	// Only non-ignored module is github.com/example/repo/sub which matches (sub/v1.2.3 vs v1.2.3)
	if hasMismatch {
		t.Fatalf("did not expect mismatch when ignored modules are skipped")
	}
}

func TestRunSync_FileValidation(t *testing.T) {
	dir := t.TempDir()
	// Missing go.mod
	_, err := runSync(Options{
		GoModPath:   filepath.Join(dir, "missing.go.mod"),
		PluginPaths: []string{filepath.Join(dir, "plugins.yaml")},
	})
	if err == nil || !strings.Contains(err.Error(), "go.mod file not found") {
		t.Fatalf("expected go.mod not found error, got: %v", err)
	}

	// go.mod exists, missing plugins
	goMod := writeFile(t, dir, "go.mod", "module x\n")
	_, err = runSync(Options{
		GoModPath:   goMod,
		PluginPaths: []string{filepath.Join(dir, "plugins.yaml")},
	})
	if err == nil || !strings.Contains(err.Error(), "plugin YAML file not found") {
		t.Fatalf("expected plugin file not found error, got: %v", err)
	}
}
