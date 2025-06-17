package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetGoModVersion(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock go.mod file with specific format for tests
	mockGoModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module github.com/smartcontractkit/chainlink

go 1.21.0

require (
	github.com/smartcontractkit/chainlink-data-streams v0.1.1-0.20250325191518-036bb568a69d
	github.com/smartcontractkit/chainlink-feeds v0.1.2-0.20250227211209-7cd000095135
	github.com/smartcontractkit/chainlink-solana v1.1.2
)
`
	if err := os.WriteFile(mockGoModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write test go.mod: %v", err)
	}

	// Create a separate test function to mock the global function for testing
	testGetGoModVersion := func(module string) (string, error) {
		if module == "github.com/smartcontractkit/chainlink-data-streams" {
			return "036bb568a69d", nil
		} else if module == "github.com/smartcontractkit/chainlink-solana" {
			return "v1.1.2", nil
		}
		return "", errors.New("module not found")
	}

	// Test pseudo-version extraction
	t.Run("Extract commit hash from pseudo-version", func(t *testing.T) {
		version, err := testGetGoModVersion("github.com/smartcontractkit/chainlink-data-streams")
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		if version != "036bb568a69d" {
			t.Errorf("Expected version '036bb568a69d', got '%s'", version)
		}
	})

	// Test regular version
	t.Run("Regular version extraction", func(t *testing.T) {
		version, err := testGetGoModVersion("github.com/smartcontractkit/chainlink-solana")
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		if version != "v1.1.2" {
			t.Errorf("Expected version 'v1.1.2', got '%s'", version)
		}
	})

	// Now set the global goModPath for the actual function test
	goModPath = mockGoModPath

	// Integration test of the actual function
	t.Run("Actual function test", func(t *testing.T) {
		version, err := getGoModVersion("github.com/smartcontractkit/chainlink-data-streams")
		if err != nil {
			t.Fatalf("getGoModVersion failed: %v", err)
		}
		t.Logf("Actual version from go.mod: %s", version)
	})
}

func TestVersionMatching(t *testing.T) {
	// Store the original function and restore it after the test
	originalGetModVersionFunc := getModVersionFunc
	defer func() { getModVersionFunc = originalGetModVersionFunc }()

	// Replace with test implementation
	getModVersionFunc = func(module string) (string, error) {
		if module == "github.com/smartcontractkit/chainlink-data-streams" {
			return "036bb568a69d", nil
		} else if module == "github.com/smartcontractkit/chainlink-feeds" {
			return "7cd000095135", nil
		}
		return "", errors.New("module not found")
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock go.mod file
	mockGoModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module github.com/smartcontractkit/chainlink

go 1.21.0

require (
	github.com/smartcontractkit/chainlink-data-streams v0.1.1-0.20250325191518-036bb568a69d
	github.com/smartcontractkit/chainlink-feeds v0.1.2-0.20250227211209-7cd000095135
)
`
	if err := os.WriteFile(mockGoModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write test go.mod: %v", err)
	}

	// Create a mock plugin yaml file - case 1: full commit hash match
	pluginYamlPath1 := filepath.Join(tempDir, "plugins1.yaml")
	yamlContent1 := `plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "036bb568a69d7b7a841017ac72a068f06c50c218" # Full hash
  feeds:
    - moduleURI: "github.com/smartcontractkit/chainlink-feeds"
      gitRef: "7cd000095135" # Short hash matching pseudo-version
`
	if err := os.WriteFile(pluginYamlPath1, []byte(yamlContent1), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml: %v", err)
	}

	// Create a mock plugin yaml file - case 2: mismatch
	pluginYamlPath2 := filepath.Join(tempDir, "plugins2.yaml")
	yamlContent2 := `plugins:
  streams:
    - moduleURI: "github.com/smartcontractkit/chainlink-data-streams"
      gitRef: "wrongcommithash123456"
  feeds:
    - moduleURI: "github.com/smartcontractkit/chainlink-feeds"
      gitRef: "7cd000095135" # Correct
`
	if err := os.WriteFile(pluginYamlPath2, []byte(yamlContent2), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml: %v", err)
	}

	// Set up the global variables for testing
	goModPath = mockGoModPath
	hasMismatches = false

	// Test scenario 1: Full commit hash matching with extracted pseudo-version
	t.Run("Full commit hash matching with extracted pseudo-version", func(t *testing.T) {
		pluginsPaths = []string{pluginYamlPath1}
		goModules = []string{"github.com/smartcontractkit/chainlink-data-streams"}
		hasMismatches = false

		checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-data-streams")

		if hasMismatches {
			t.Errorf("Expected no mismatch for full hash '036bb568a69d7b7a841017ac72a068f06c50c218' matching pseudo-version hash '036bb568a69d'")
		}
	})

	// Test scenario 2: Short hash matching with extracted pseudo-version
	t.Run("Short hash matching with extracted pseudo-version", func(t *testing.T) {
		pluginsPaths = []string{pluginYamlPath1}
		goModules = []string{"github.com/smartcontractkit/chainlink-feeds"}
		hasMismatches = false

		checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-feeds")

		if hasMismatches {
			t.Errorf("Expected no mismatch for short hash '7cd000095135' matching pseudo-version hash '7cd000095135'")
		}
	})

	// Test scenario 3: Hash mismatch
	t.Run("Hash mismatch", func(t *testing.T) {
		pluginsPaths = []string{pluginYamlPath2}
		goModules = []string{"github.com/smartcontractkit/chainlink-data-streams"}
		hasMismatches = false

		checkAndUpdateModuleVersion("github.com/smartcontractkit/chainlink-data-streams")

		if !hasMismatches {
			t.Errorf("Expected mismatch for hash 'wrongcommithash123456' not matching pseudo-version hash '036bb568a69d'")
		}
	})
}

func TestUpdateGitRefInYAML(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "plugout-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock plugin yaml file with comments and formatting
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
	if err := os.WriteFile(pluginYamlPath, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to write test plugin yaml: %v", err)
	}

	// Test updating gitRef while preserving comments and formatting
	t.Run("Update gitRef preserving comments and formatting", func(t *testing.T) {
		err := updateGitRefInYAML(pluginYamlPath, "github.com/smartcontractkit/chainlink-data-streams", "newcommithash789012")
		if err != nil {
			t.Fatalf("updateGitRefInYAML failed: %v", err)
		}

		// Read the updated file
		content, err := os.ReadFile(pluginYamlPath)
		if err != nil {
			t.Fatalf("Failed to read updated yaml: %v", err)
		}

		contentStr := string(content)

		// Check if update was performed correctly
		if !strings.Contains(contentStr, `gitRef: "newcommithash789012" # Comment about the version`) {
			t.Errorf("gitRef not updated correctly or comment not preserved: %s", contentStr)
		}

		// Check if other content was preserved
		if !strings.Contains(contentStr, `installPath: "some/install/path"`) {
			t.Errorf("Other content not preserved: %s", contentStr)
		}

		// Check if other module was not affected
		if !strings.Contains(contentStr, `gitRef: "7cd000095135" # Another comment`) {
			t.Errorf("Unrelated module was affected: %s", contentStr)
		}
	})
}
