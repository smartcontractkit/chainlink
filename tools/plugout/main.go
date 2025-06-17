package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Configuration options
var (
	goModPath     string
	updatePlugins bool
	hasMismatches bool
	goModules     []string
	pluginsPaths  []string
)

// Default values for modules and plugin paths
var (
	defaultModules = []string{
		"github.com/smartcontractkit/chainlink-data-streams",
		"github.com/smartcontractkit/chainlink-feeds",
		"github.com/smartcontractkit/chainlink-solana",
	}
	defaultPluginsPaths = []string{
		"./plugins/plugins.public.yaml",
	}
)

// Plugin schema for YAML parsing
type PluginsFile struct {
	Plugins map[string][]Plugin `yaml:"plugins"`
}

type Plugin struct {
	ModuleURI string `yaml:"moduleURI"`
	GitRef    string `yaml:"gitRef"`
}

// For testing purposes
var getModVersionFunc = getGoModVersion

func main() {
	rootCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync plugin versions from go.mod to plugins manifest YAML files",
		PreRun: func(cmd *cobra.Command, args []string) {
			// If no modules provided, use defaults
			if len(goModules) == 0 {
				goModules = defaultModules
			}

			// If no plugin paths provided, use defaults
			if len(pluginsPaths) == 0 {
				pluginsPaths = defaultPluginsPaths
			}
		},
		Run: runSync,
	}

	rootCmd.Flags().StringVar(&goModPath, "go-mod", "./go.mod", "Path to go.mod file")
	rootCmd.Flags().BoolVar(&updatePlugins, "update", false, "Write the gitRef using the go.mod version for matching plugins")
	rootCmd.Flags().StringArrayVar(&goModules, "module", nil, "Go module URI to check (can be specified multiple times)")
	rootCmd.Flags().StringArrayVar(&pluginsPaths, "plugin-file", nil, "Plugin YAML file to check (can be specified multiple times)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if hasMismatches && !updatePlugins {
		os.Exit(1)
	}
}

func runSync(cmd *cobra.Command, args []string) {
	fmt.Println("=== Starting plugin version sync check ===")
	fmt.Printf("Checking go.mod path: %s\n", goModPath)
	fmt.Printf("Modules to check: %s\n", strings.Join(goModules, ", "))
	fmt.Printf("Plugin files to check: %s\n", strings.Join(pluginsPaths, ", "))

	if updatePlugins {
		fmt.Println("Mode: UPDATE (will update plugin gitRef values)")
	} else {
		fmt.Println("Mode: CHECK ONLY (use --update flag to update plugin files)")
	}
	fmt.Println()

	// Validate that go.mod file exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		fmt.Printf("Error: go.mod file not found at path: %s\n", goModPath)
		os.Exit(1)
	}

	// Validate all plugin YAML files exist before continuing
	for _, pluginPath := range pluginsPaths {
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			fmt.Printf("Error: Plugin YAML file not found: %s\n", pluginPath)
			os.Exit(1)
		}
	}

	for _, module := range goModules {
		checkAndUpdateModuleVersion(module)
	}

	if hasMismatches && !updatePlugins {
		fmt.Println("=== Plugin version sync check completed with mismatches ===")
	} else {
		fmt.Println("=== Plugin version sync check completed successfully ===")
	}
}

// getGoModVersion extracts the version from go.mod for a specific module
func getGoModVersion(module string) (string, error) {
	fmt.Printf("Extracting version for %s from go.mod...\n", module)

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Parse go.mod to find the module version
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	foundModule := false
	pseudoVersionPattern := regexp.MustCompile(`v[\d]+\.[\d]+\.[\d]+-[\d]+-g([0-9a-f]+)`)

	for scanner.Scan() {
		line := scanner.Text()
		// Check if this line contains the module we're looking for
		if strings.Contains(line, module+" ") {
			foundModule = true
			// If the version is on the same line
			fields := strings.Fields(line)
			if len(fields) > 1 {
				version := fields[len(fields)-1]

				// Check if it's a pseudo-version and extract the commit hash
				if pseudoVersionPattern.MatchString(version) {
					matches := pseudoVersionPattern.FindStringSubmatch(version)
					if len(matches) >= 2 {
						version = matches[1]
					}
				}

				fmt.Printf("Version extracted: %s\n", version)
				return version, nil
			}
			continue
		}

		// If we previously found the module, and this line might contain the version
		if foundModule {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				version := fields[0]

				// Check if it's a pseudo-version and extract the commit hash
				if pseudoVersionPattern.MatchString(version) {
					matches := pseudoVersionPattern.FindStringSubmatch(version)
					if len(matches) >= 2 {
						version = matches[1]
					}
				}

				fmt.Printf("Version extracted: %s\n", version)
				return version, nil
			}
		}
	}

	return "", errors.New("module " + module + " not found in go.mod")
}

// getYAMLVersion extracts the gitRef from a plugin YAML file for a specific module
func getYAMLVersion(pluginPath, module string) (string, error) {
	fmt.Printf("Extracting version for %s from %s...\n", module, pluginPath)

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file: %w", err)
	}

	var pluginsFile PluginsFile
	if err := yaml.Unmarshal(data, &pluginsFile); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, plugins := range pluginsFile.Plugins {
		for _, plugin := range plugins {
			if plugin.ModuleURI == module {
				return plugin.GitRef, nil
			}
		}
	}

	return "", errors.New("module " + module + " not found in " + pluginPath)
}

// updateGitRefInYAML updates the gitRef in a plugin YAML file for a specific module
func updateGitRefInYAML(pluginPath, module, newGitRef string) error {
	fmt.Printf("  Updating gitRef for %s to %s in %s\n", module, newGitRef, pluginPath)

	// Read the original file content
	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// First unmarshal to verify the module exists
	var pluginsFile PluginsFile
	if err := yaml.Unmarshal(data, &pluginsFile); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Check if module exists in the file
	moduleExists := false
	for _, plugins := range pluginsFile.Plugins {
		for _, plugin := range plugins {
			if plugin.ModuleURI == module {
				moduleExists = true
				break
			}
		}
		if moduleExists {
			break
		}
	}

	if !moduleExists {
		return errors.New("module " + module + " not found in " + pluginPath)
	}

	// Convert to string for line-by-line processing
	content := string(data)
	lines := strings.Split(content, "\n")

	// Track if we found the module and if we've updated the gitRef
	foundModule := false
	updated := false

	// Process line by line
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if this line contains the moduleURI
		if strings.Contains(trimmedLine, "moduleURI:") && strings.Contains(trimmedLine, module) {
			foundModule = true
			continue
		}

		// If we found the module, look for the gitRef line to update
		if foundModule && strings.Contains(trimmedLine, "gitRef:") {
			// Extract indentation from the original line
			indentation := ""
			for _, char := range line {
				if char != ' ' && char != '\t' {
					break
				}
				indentation += string(char)
			}

			// Keep any comment that might be on the same line
			comment := ""
			if idx := strings.Index(line, "#"); idx >= 0 {
				comment = " " + strings.TrimSpace(line[idx:])
			}

			// Replace the line with updated gitRef, preserving indentation and comment
			lines[i] = fmt.Sprintf("%sgitRef: %q%s", indentation, newGitRef, comment)
			updated = true
			foundModule = false // Reset to find next occurrence if any
			continue
		}
	}

	if !updated {
		return errors.New("failed to update gitRef for module " + module)
	}

	// Join lines back into content
	newContent := strings.Join(lines, "\n")

	// Write the updated content back to file
	if err := os.WriteFile(pluginPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// checkAndUpdateModuleVersion checks the version of a module in go.mod and all plugin files
func checkAndUpdateModuleVersion(module string) {
	goModVersion, err := getModVersionFunc(module)
	if err != nil {
		fmt.Printf("  ⚠️  %v\n", err)
		return
	}

	for _, pluginPath := range pluginsPaths {
		yamlVersion, err := getYAMLVersion(pluginPath, module)
		if err != nil {
			fmt.Printf("  ⚠️  %v\n", err)
			continue
		}

		// Check if versions match
		versionMatch := false

		// Case 1: Direct exact match
		if goModVersion == yamlVersion {
			versionMatch = true
		} else if !strings.HasPrefix(goModVersion, "v") && strings.HasPrefix(yamlVersion, goModVersion) {
			// Case 2: the yaml gitRef/version contains sha sum from the go.mod pseudo-version
			versionMatch = true
		}

		if !versionMatch {
			fmt.Printf("  ❌ MISMATCH: %s in %s\n", module, pluginPath)
			fmt.Printf("    go.mod: %s\n", goModVersion)
			fmt.Printf("    yaml  : %s\n", yamlVersion)

			if updatePlugins {
				if err := updateGitRefInYAML(pluginPath, module, goModVersion); err != nil {
					fmt.Printf("    ❌ %v\n", err)
				} else {
					fmt.Printf("    ✅ Updated gitRef in %s\n", pluginPath)
				}
			} else {
				hasMismatches = true
			}
		} else {
			fmt.Printf("  ✅ %s versions match in %s\n", module, pluginPath)
		}
	}
}
