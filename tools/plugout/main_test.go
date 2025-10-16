package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetGoModVersion(t *testing.T) {
	// temp workspace
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// mock go.mod
	mockGoModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module github.com/smartcontractkit/chainlink

go 1.21.0

require (
	github.com/smartcontractkit/chainlink-data-streams v0.1.1-0.20250325191518-036bb568a69d
	github.com/smartcontractkit/chainlink-feeds v0.1.2-0.20250227211209-7cd000095135
	github.com/smartcontractkit/chainlink-solana v1.1.2
	github.com/smartcontractkit/chainlink-tron/relayer v0.1.2-0.20250227211209-7cd111195135 // indirect
)
`
	if err = os.WriteFile(mockGoModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write test go.mod: %v", err)
	}

	versionsMap := map[string]ModuleVersion{
		"github.com/smartcontractkit/chainlink-data-streams": {Raw: "v0.1.1-0.20250325191518-036bb568a69d", SHA: "036bb568a69d"},
		"github.com/smartcontractkit/chainlink-feeds":        {Raw: "v0.1.2-0.20250227211209-7cd000095135", SHA: "7cd000095135"},
		"github.com/smartcontractkit/chainlink-solana":       {Raw: "v1.1.2", Tag: "v1.1.2"},
		"github.com/smartcontractkit/chainlink-tron/relayer": {Raw: "v0.1.2-0.20250227211209-7cd111195135", SHA: "7cd111195135"},
	}

	t.Run("Extract commit hash from pseudoversion", func(t *testing.T) {
		dep := "github.com/smartcontractkit/chainlink-feeds"
		mv, err := getGoModVersion(mockGoModPath, dep)
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		if mv.SHA != versionsMap[dep].SHA {
			t.Errorf("Expected SHA '%s', got '%s' (raw: %s)", versionsMap[dep].SHA, mv.SHA, mv.Raw)
		}
		if mv.Tag != "" {
			t.Errorf("Expected tag empty for pseudoversion, got '%s'", mv.Tag)
		}
	})

	t.Run("Regular tag extraction", func(t *testing.T) {
		dep := "github.com/smartcontractkit/chainlink-solana"
		mv, err := getGoModVersion(mockGoModPath, dep)
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		if mv.Tag != versionsMap[dep].Tag {
			t.Errorf("Expected tag '%s', got '%s' (raw: %s)", versionsMap[dep].Tag, mv.Tag, mv.Raw)
		}
		if mv.SHA != "" {
			t.Errorf("Expected SHA empty for tag, got '%s'", mv.SHA)
		}
	})

	t.Run("Missing module returns error", func(t *testing.T) {
		_, err := getGoModVersion(mockGoModPath, "github.com/smartcontractkit/missing-module")
		if err == nil {
			t.Fatalf("Expected error for missing module, got nil")
		}
	})

	t.Run("Indirect module versions should match", func(t *testing.T) {
		dep := "github.com/smartcontractkit/chainlink-tron/relayer"
		mv, err := getGoModVersion(mockGoModPath, dep)
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		if mv.SHA != versionsMap[dep].SHA {
			t.Errorf("Expected SHA %s, got '%s'", versionsMap[dep].SHA, mv.SHA)
		}
	})
}

func TestVersionMatching_SHAAndPseudo(t *testing.T) {
	// stub getModVersion via Options
	get := func(_ string, module string) (ModuleVersion, error) {
		switch module {
		case "github.com/smartcontractkit/chainlink-data-streams":
			return ModuleVersion{Raw: "v0.1.1-0.20250325191518-036bb568a69d", SHA: "036bb568a69d"}, nil
		case "github.com/smartcontractkit/chainlink-feeds":
			return ModuleVersion{Raw: "v0.1.2-0.20250227211209-7cd000095135", SHA: "7cd000095135"}, nil
		default:
			return ModuleVersion{}, errors.New("module not found")
		}
	}

	// temp workspace
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// (go.mod exists but not read by our stub)
	goModPath := filepath.Join(tempDir, "go.mod")
	if err = os.WriteFile(goModPath, []byte("module x\ngo 1.21\n"), 0600); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// plugin with full & short SHA
	pluginYamlPath1 := filepath.Join(tempDir, "plugins1.yaml")
	yamlContent1 := `plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "036bb568a69d7b7a841017ac72a068f06c50c218" # Full hash
  feeds:
    - moduleURI: "github.com/smartcontractkit/chainlink-feeds"
      gitRef: "7cd000095135" # Short hash
`
	if err = os.WriteFile(pluginYamlPath1, []byte(yamlContent1), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml 1: %v", err)
	}

	pluginYamlPath2 := filepath.Join(tempDir, "plugins2.yaml")
	yamlContent2 := `plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "wrongcommithash123456"
`
	if err = os.WriteFile(pluginYamlPath2, []byte(yamlContent2), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml 2: %v", err)
	}

	t.Run("Full SHA matches pseudo SHA", func(t *testing.T) {
		opts := Options{
			GoModPath:     goModPath,
			PluginPaths:   []string{pluginYamlPath1},
			IgnoreModules: nil,
			Update:        false,
			GetModVersion: get,
		}
		mismatch, err := checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-data-streams", opts)
		if err != nil {
			t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
		}
		if mismatch {
			t.Fatalf("should not mismatch full=036bb... vs pseudo SHA=036bb...")
		}
	})

	t.Run("Short SHA matches pseudo SHA", func(t *testing.T) {
		opts := Options{
			GoModPath:     goModPath,
			PluginPaths:   []string{pluginYamlPath1},
			IgnoreModules: nil,
			Update:        false,
			GetModVersion: get,
		}
		mismatch, err := checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-feeds", opts)
		if err != nil {
			t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
		}
		if mismatch {
			t.Fatalf("should not mismatch short=7cd000095135 vs pseudo SHA=7cd000095135")
		}
	})

	t.Run("Wrong hash yields mismatch", func(t *testing.T) {
		opts := Options{
			GoModPath:     goModPath,
			PluginPaths:   []string{pluginYamlPath2},
			IgnoreModules: nil,
			Update:        false,
			GetModVersion: get,
		}
		mismatch, err := checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-data-streams", opts)
		if err != nil {
			t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
		}
		if !mismatch {
			t.Fatalf("expected mismatch to be true")
		}
	})
}

func TestSubmoduleTagMatchingAndUpdate(t *testing.T) {
	// stub: go.mod tag v1.2.3 for submodule
	get := func(_ string, module string) (ModuleVersion, error) {
		if module == "github.com/smartcontractkit/chainlink-starknet/relayer" {
			return ModuleVersion{Raw: "v1.2.3", Tag: "v1.2.3"}, nil
		}
		return ModuleVersion{}, errors.New("not found")
	}

	// temp workspace
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// minimal go.mod (exists)
	goModPath := filepath.Join(tempDir, "go.mod")
	if err = os.WriteFile(goModPath, []byte("module x\ngo 1.21\n"), 0600); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// YAML has WRONG submodule tag first to trigger mismatch
	pluginYamlPath := filepath.Join(tempDir, "plugins_sub.yaml")
	if err = os.WriteFile(pluginYamlPath, []byte(`plugins:
  starknet:
    - moduleURI: "github.com/smartcontractkit/chainlink-starknet/relayer"
      gitRef: "relayer/v0.9.9"
`), 0600); err != nil {
		t.Fatalf("Failed to write plugin yaml: %v", err)
	}

	// 1) mismatch expected
	opts := Options{
		GoModPath:     goModPath,
		PluginPaths:   []string{pluginYamlPath},
		IgnoreModules: nil,
		Update:        false,
		GetModVersion: get,
	}
	mismatch, err := checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-starknet/relayer", opts)
	if err != nil {
		t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
	}
	if !mismatch {
		t.Fatalf("expected mismatch for YAML relayer/v0.9.9 vs go.mod v1.2.3")
	}

	// 2) update path: expect YAML to be rewritten to relayer/v1.2.3
	opts.Update = true
	mismatch, err = checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-starknet/relayer", opts)
	if err != nil {
		t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
	}
	if !mismatch {
		t.Fatalf("expected mismatch=true even after update")
	}
	// After update, verify file contents
	b, err := os.ReadFile(pluginYamlPath)
	if err != nil {
		t.Fatalf("Failed to read updated yaml: %v", err)
	}
	if !strings.Contains(string(b), `gitRef: "relayer/v1.2.3"`) {
		t.Fatalf("Expected YAML gitRef to be updated to 'relayer/v1.2.3', got:\n%s", string(b))
	}
	// mismatch result can be true or false depending on whether other entries still mismatch; we don't rely on it here
}

func TestDeclaredPluginMissingInGoMod_WarnsAndSkips(t *testing.T) {
	// stub: always "not found in go.mod"
	get := func(_ string, _ string) (ModuleVersion, error) {
		return ModuleVersion{}, errors.New("module not found in go.mod")
	}

	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// minimal go.mod presence
	goModPath := filepath.Join(tempDir, "go.mod")
	if err = os.WriteFile(goModPath, []byte("module x\ngo 1.21\n"), 0600); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// plugin declares a module not present in go.mod
	pluginYamlPath := filepath.Join(tempDir, "plugins_missing.yaml")
	if err = os.WriteFile(pluginYamlPath, []byte(`plugins:
  sui:
    - moduleURI: "github.com/smartcontractkit/chainlink-sui"
      gitRef: "v0.0.1"
`), 0600); err != nil {
		t.Fatalf("Failed to write plugin yaml: %v", err)
	}

	opts := Options{
		GoModPath:     goModPath,
		PluginPaths:   []string{pluginYamlPath},
		IgnoreModules: nil,
		Update:        false,
		GetModVersion: get,
	}

	mismatch, err := checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-sui", opts)
	if err != nil {
		t.Fatalf("checkAndUpdateModuleVersion error: %v", err)
	}
	if mismatch {
		t.Fatalf("Expected warn+skip (no mismatch) when module absent from go.mod")
	}
}

func TestUpdateGitRefInYAML_PreservesFormatting(t *testing.T) {
	// temp workspace
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	pluginYamlPath := filepath.Join(tempDir, "plugins.yaml")
	yamlContent := `# Plugin configuration
plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "oldcommithash123456" # Comment about the version
      installPath: "some/install/path"
  feeds:
    - moduleURI: "github.com/smartcontractkit/chainlink-feeds"
      gitRef: "7cd000095135" # Another comment
`
	if err = os.WriteFile(pluginYamlPath, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml: %v", err)
	}

	mv := ModuleVersion{Raw: "newcommithash789012", SHA: "newcommithash789012"}
	if err = updateGitRefInYAML(pluginYamlPath, "github.com/smartcontractkit/chainlink-data-streams", mv); err != nil {
		t.Fatalf("updateGitRefInYAML failed: %v", err)
	}

	b, err := os.ReadFile(pluginYamlPath)
	if err != nil {
		t.Fatalf("Failed to read updated yaml: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `gitRef: "newcommithash789012" # Comment about the version`) {
		t.Errorf("gitRef not updated correctly or comment not preserved:\n%s", s)
	}
	if !strings.Contains(s, `installPath: "some/install/path"`) {
		t.Errorf("Other content not preserved:\n%s", s)
	}
	if !strings.Contains(s, `gitRef: "7cd000095135" # Another comment`) {
		t.Errorf("Unrelated module affected:\n%s", s)
	}
}

func TestRunSync_Integration_UsesDiscoveryAndIgnore(t *testing.T) {
	// stub: two modules with different SHAs
	get := func(_ string, module string) (ModuleVersion, error) {
		switch module {
		case "github.com/smartcontractkit/chainlink-data-streams":
			return ModuleVersion{Raw: "v0.0.0-20250101000000-aaaaaaaaaaaa", SHA: "aaaaaaaaaaaa"}, nil
		case "github.com/smartcontractkit/chainlink-feeds":
			return ModuleVersion{Raw: "v0.0.0-20250101000000-bbbbbbbbbbbb", SHA: "bbbbbbbbbbbb"}, nil
		default:
			return ModuleVersion{}, errors.New("module not found")
		}
	}

	// temp workspace
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	if err = os.WriteFile(goModPath, []byte("module x\ngo 1.21\n"), 0600); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	pluginsPath := filepath.Join(tempDir, "plugins.yaml")
	if err = os.WriteFile(pluginsPath, []byte(`plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "aaaaaaaaaaaa"   # matches SHA
  feeds:
    - moduleURI: "github.com/smartcontractkit/chainlink-feeds"
      gitRef: "WRONGHASH"      # mismatch
`), 0600); err != nil {
		t.Fatalf("Failed to write plugins.yaml: %v", err)
	}

	// CHECK mode: expect mismatch due to 'feeds'
	opts := Options{
		GoModPath:     goModPath,
		PluginPaths:   []string{pluginsPath},
		IgnoreModules: nil,
		Update:        false,
		GetModVersion: get,
	}
	hasMismatch, err := runSync(opts)
	if err != nil {
		t.Fatalf("runSync error: %v", err)
	}
	if !hasMismatch {
		t.Fatalf("Expected hasMismatch=true due to the 'feeds' mismatch")
	}

	// Now ignore the mismatching module: expect no mismatches
	opts.IgnoreModules = []string{"github.com/smartcontractkit/chainlink-feeds"}
	hasMismatch, err = runSync(opts)
	if err != nil {
		t.Fatalf("runSync error: %v", err)
	}
	if hasMismatch {
		t.Fatalf("Expected hasMismatch=false when ignoring mismatching module")
	}
}
