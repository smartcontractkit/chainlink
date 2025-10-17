package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Options struct {
	GoModPath     string
	PluginPaths   []string
	IgnoreModules []string
	Update        bool

	// Test seam; if nil, defaults to getGoModVersion
	GetModVersion func(goModPath, module string) (ModuleVersion, error)
}

type PluginsFile struct {
	Plugins map[string][]Plugin `yaml:"plugins"`
}

type Plugin struct {
	ModuleURI   string   `yaml:"moduleURI"`
	GitRef      string   `yaml:"gitRef"`
	InstallPath string   `yaml:"installPath"`
	Libs        []string `yaml:"libs"`
}

// Normalized version representation.
type ModuleVersion struct {
	Raw       string // the original string (tag/pseudo/SHA)
	SHA       string // extracted commit SHA if available (7..40 lower-hex)
	Tag       string // tag like v0.1.5 (without any subdir prefix)
	TagPrefix string // if raw looked like "sub/dir/vX.Y.Z", TagPrefix is "sub/dir"
}

func main() {
	var (
		flagGoModPath     string
		flagUpdatePlugins bool
		flagPluginPaths   []string
		flagIgnoreModules []string
	)

	rootCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync plugin versions from go.mod to plugins manifest YAML files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(flagPluginPaths) == 0 {
				flagPluginPaths = []string{"./plugins/plugins.public.yaml"}
			}

			opts := Options{
				GoModPath:     flagGoModPath,
				PluginPaths:   flagPluginPaths,
				IgnoreModules: flagIgnoreModules,
				Update:        flagUpdatePlugins,
				GetModVersion: nil, // use default
			}

			hasMismatch, err := runSync(opts)
			if err != nil {
				return err
			}
			if hasMismatch && !opts.Update {
				// Non-zero exit on mismatches in CHECK mode
				os.Exit(1)
			}
			return nil
		},
	}

	rootCmd.Flags().StringVar(&flagGoModPath, "go-mod", "./go.mod", "Path to go.mod file")
	rootCmd.Flags().BoolVar(&flagUpdatePlugins, "update", false, "Write the gitRef using the go.mod version for matching plugins")
	rootCmd.Flags().StringArrayVar(&flagPluginPaths, "plugin-file", nil, "Plugin YAML file to check (can be specified multiple times)")
	rootCmd.Flags().StringArrayVar(&flagIgnoreModules, "ignore-module", nil, "Module URI to ignore (can be specified multiple times)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runSync discovers modules from plugin files (minus ignores), verifies files,
// compares/updates refs, and returns whether any mismatches were found.
func runSync(opts Options) (bool, error) {
	fmt.Println("=== Starting plugin version sync check ===")
	fmt.Printf("Checking go.mod path: %s\n", opts.GoModPath)
	fmt.Printf("Plugin files to check: %s\n", strings.Join(opts.PluginPaths, ", "))

	if opts.Update {
		fmt.Println("Mode: UPDATE (will update plugin gitRef values)")
	} else {
		fmt.Println("Mode: CHECK ONLY (use --update flag to update plugin files)")
	}
	fmt.Println()

	// Validate that go.mod file exists.
	if _, err := os.Stat(opts.GoModPath); os.IsNotExist(err) {
		return false, fmt.Errorf("go.mod file not found at path: %s", opts.GoModPath)
	}

	// Validate that plugin YAML files exist.
	for _, pluginPath := range opts.PluginPaths {
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			return false, fmt.Errorf("plugin YAML file not found: %s", pluginPath)
		}
	}

	get := opts.GetModVersion
	if get == nil {
		get = getGoModVersion
	}

	hasMismatch := false
	for _, p := range opts.PluginPaths {
		fmt.Printf("\n=== Checking plugin file %s ===\n", p)
		modulesToVersions, err := discoverPluginVersions(p)
		if err != nil {
			fmt.Printf("  - ❌ Error discovering plugin versions: %v\n", err)
			hasMismatch = true
			continue
		}

		modulesToCheck := make([]string, 0, len(modulesToVersions))
		for k := range modulesToVersions {
			modulesToCheck = append(modulesToCheck, k)
		}
		sort.Strings(modulesToCheck)

		if len(modulesToCheck) == 0 {
			fmt.Printf("  - No modules found in %s\n", p)
			continue
		}

		for idx, module := range modulesToCheck {
			// Skip ignored modules
			if len(opts.IgnoreModules) > 0 && contains(opts.IgnoreModules, module) {
				fmt.Printf("Ignoring module (%d/%d): %s\n", idx+1, len(modulesToCheck), module)
				continue
			}

			fmt.Printf("Checking Module (%d/%d): %s \n", idx+1, len(modulesToCheck), module)

			normalizedPluginsVersion := normalizeVersion(modulesToVersions[module])
			fmt.Printf("  - Plugins version: %s (%s)\n", modulesToVersions[module], normalizedPluginsVersion.toString())

			goModVersion, err := get(opts.GoModPath, module)
			if err != nil || goModVersion.Raw == "" {
				fmt.Printf("  - No version found in go.mod for %s: %v\n", module, err)
				continue
			}

			matches := versionsMatchForModule(module, goModVersion, normalizedPluginsVersion)
			if matches {
				fmt.Printf("  - ✅ Versions match for %s\n", module)
				continue
			}

			if opts.Update {
				err := updateGitRefInYAML(p, module, goModVersion)
				if err != nil {
					fmt.Printf("  - ❌ Failed to update gitRef in %s: %v\n", p, err)
					hasMismatch = true
					continue
				}
				fmt.Printf("  - ✅ Updated gitRef in %s to %s\n", p, desiredYAMLRefForModule(module, goModVersion))
				continue
			}

			fmt.Printf("  - ❌ MISMATCH for %s: go.mod has %s, plugins file has %s\n", module, goModVersion.toString(), normalizedPluginsVersion.toString())
			hasMismatch = true
		}
	}

	if hasMismatch && !opts.Update {
		fmt.Println("=== Plugin version sync check completed with mismatches ===")
	} else {
		fmt.Println("=== Plugin version sync check completed successfully ===")
	}

	return hasMismatch, nil
}

// -----------------------------------------------------------------------------
// Discovery & helpers
// -----------------------------------------------------------------------------

// discoverPluginVersions returns a unique list of module URIs and their versions from a plugin YAML file
func discoverPluginVersions(path string) (map[string]string, error) {
	versions := make(map[string]string) // moduleURI -> gitRef
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	var pf PluginsFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("failed to parse YAML %s: %w", path, err)
	}
	for _, list := range pf.Plugins {
		for _, plugin := range list {
			if versions[plugin.ModuleURI] != "" {
				fmt.Printf("  - Warning: duplicate moduleURI %s in %s\n", plugin.ModuleURI, path)
			}

			if plugin.ModuleURI != "" && plugin.GitRef != "" {
				versions[plugin.ModuleURI] = plugin.GitRef
			}
		}
	}

	return versions, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// moduleSubdir returns the path within the repo after the "host/org/repo" prefix.
// e.g.  github.com/smartcontractkit/chainlink-starknet/relayer -> "relayer"
func moduleSubdir(module string) string {
	parts := strings.Split(module, "/")
	if len(parts) <= 3 {
		return ""
	}
	return strings.Join(parts[3:], "/")
}

var (
	pseudoWithSHARe = regexp.MustCompile(
		`^v\d+\.\d+\.\d+(?:-\d{14}|-(?:0|[0-9A-Za-z-]+\.0)\.\d{14})-([0-9a-f]{7,40})$`,
	)
	plainTagRe    = regexp.MustCompile(`^v\d+\.\d+\.\d+([.-].*)?$`)
	prefixedTagRe = regexp.MustCompile(`^(.+?)/+(v\d+\.\d+\.\d+(?:[.-].*)?)$`)
	shaOnlyRe     = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	// NEW: matches the middle part of a valid pseudoversion
	pseudoMiddleRe = regexp.MustCompile(`^(?:\d{14}|(?:0|[0-9A-Za-z-]+\.0)\.\d{14})$`)
)

func normalizeVersion(raw string) ModuleVersion {
	mv := ModuleVersion{Raw: raw}
	low := strings.ToLower(strings.TrimSpace(raw))

	// 1) Pseudoversion? (strict)
	if m := pseudoWithSHARe.FindStringSubmatch(low); m != nil {
		mv.SHA = m[1]
		return mv
	}

	// 1b) Fallback pseudoversion (guarded)
	if strings.HasPrefix(low, "v") && strings.Count(low, "-") == 2 {
		parts := strings.Split(low, "-") // ["vX.Y.Z", middle, sha-ish]
		middle := parts[1]
		last := strings.TrimPrefix(parts[2], "g") // tolerate optional 'g' prefix
		if pseudoMiddleRe.MatchString(middle) && shaOnlyRe.MatchString(last) {
			mv.SHA = last
			return mv
		}
	}

	// 2) Raw SHA?
	if shaOnlyRe.MatchString(low) {
		mv.SHA = low
		return mv
	}

	// 3) Prefixed tag?
	if m := prefixedTagRe.FindStringSubmatch(low); len(m) == 3 && plainTagRe.MatchString(m[2]) {
		orig := strings.TrimSpace(raw)
		if pos := strings.LastIndex(orig, "/"); pos >= 0 && pos+1 < len(orig) {
			return ModuleVersion{
				Raw:       raw,
				Tag:       orig[pos+1:],
				TagPrefix: strings.TrimSuffix(orig[:pos], "/"),
			}
		}
	}

	// 4) Plain tag?
	if plainTagRe.MatchString(low) {
		mv.Tag = raw
		return mv
	}

	return mv
}

func (m *ModuleVersion) toString() string {
	if m.Tag != "" && m.TagPrefix != "" {
		return fmt.Sprintf("Tag: %s/%s", m.TagPrefix, m.Tag)
	}
	if m.Tag != "" {
		return "Tag: " + m.Tag
	}
	if m.SHA != "" {
		return "SHA: " + m.SHA
	}
	return "Raw: " + m.Raw
}

func shaEqual(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	// Allow prefix matches (e.g., 12-char vs 40-char).
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}

// getGoModVersion extracts the version for a specific module by invoking `go list -m`.
func getGoModVersion(goModPath, module string) (ModuleVersion, error) {
	// The working directory should be the module root containing go.mod
	modDir := filepath.Dir(goModPath)

	cmd := exec.CommandContext(context.Background(), "go", "list", "-m", "-json", "-mod=readonly", module)
	cmd.Dir = modDir

	out, err := cmd.Output()
	if err != nil {
		return ModuleVersion{}, fmt.Errorf("failed to run go list -m: %w", err)
	}

	var m struct {
		Path    string
		Version string
		Replace *struct {
			Path    string
			Version string
		}
	}
	if err := json.Unmarshal(out, &m); err != nil {
		return ModuleVersion{}, fmt.Errorf("failed to parse go list output: %w", err)
	}

	version := m.Version
	if m.Replace != nil && m.Replace.Version != "" {
		version = m.Replace.Version
	}

	mv := normalizeVersion(version)
	fmt.Printf("  - Version in go.mod: %s (%s)\n", version, mv.toString())
	return mv, nil
}

// updateGitRefInYAML updates the gitRef in a plugin YAML file for a specific module.
// If go.mod provided a TAG and the module is a submodule, we write "sub/dir/vX.Y.Z".
func updateGitRefInYAML(pluginPath, module string, goModMV ModuleVersion) error {
	newGitRef := desiredYAMLRefForModule(module, goModMV)
	fmt.Printf("  Updating gitRef for %s to %s in %s\n", module, newGitRef, pluginPath)

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Verify the module exists in YAML.
	var pluginsFile PluginsFile
	if err := yaml.Unmarshal(data, &pluginsFile); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
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

	// Line-wise replace to preserve formatting/comments.
	content := string(data)
	lines := strings.Split(content, "\n")

	foundModule := false
	updated := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "moduleURI:") && strings.Contains(trimmed, module) {
			foundModule = true
			continue
		}

		if foundModule && strings.Contains(trimmed, "gitRef:") {
			// preserve indentation
			indent := ""
			for _, ch := range line {
				if ch != ' ' && ch != '\t' {
					break
				}
				indent += string(ch)
			}
			// preserve comment
			comment := ""
			if idx := strings.Index(line, "#"); idx >= 0 {
				comment = " " + strings.TrimSpace(line[idx:])
			}
			lines[i] = fmt.Sprintf("%sgitRef: %q%s", indent, newGitRef, comment)
			updated = true
			foundModule = false
			continue
		}
	}

	if !updated {
		return errors.New("failed to update gitRef for module " + module)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(pluginPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}
	return nil
}

// desiredYAMLRefForModule returns the string to write to YAML for a module/ref.
// If ref is a tag and module has a subdir, return "subdir/tag". Otherwise return raw.
func desiredYAMLRefForModule(module string, mv ModuleVersion) string {
	if mv.Tag != "" {
		if sub := moduleSubdir(module); sub != "" {
			return sub + "/" + mv.Tag
		}
	}
	return mv.Raw
}

// versionsMatchForModule extends versionsMatch with submodule tag-prefix logic.
func versionsMatchForModule(module string, a, b ModuleVersion) bool {
	// 1) SHA equality (prefix-friendly).
	if shaEqual(a.SHA, b.SHA) {
		return true
	}

	// 2) If YAML raw contains the go.mod SHA (e.g., YAML has full pseudo, go.mod SHA was normalized).
	if a.SHA != "" && strings.Contains(strings.ToLower(b.Raw), a.SHA) {
		return true
	}

	// 3) Tag equality, accounting for submodule tag prefixes.
	if tagsMatchWithSubdir(module, a, b) {
		return true
	}

	// 4) Raw equality fallback.
	return a.Raw != "" && a.Raw == b.Raw
}

// tagsMatchWithSubdir considers these equivalent for module "repo/sub":
//
//	go.mod: v1.2.3  <=>  YAML: sub/v1.2.3
//
// Works with multi-segment subdirs and v2+ paths like "v2/sub".
func tagsMatchWithSubdir(module string, a, b ModuleVersion) bool {
	if a.Tag == "" && b.Tag == "" {
		return false
	}
	if a.Tag != "" && b.Tag != "" && a.Tag == b.Tag {
		return true
	}

	sub := moduleSubdir(module)
	if sub == "" {
		// root module: just compare plain tags if present
		return a.Tag != "" && b.Tag != "" && a.Tag == b.Tag
	}

	// If either side has a prefix already, normalize both possibilities.
	// Accept:
	//   a.Tag == b.Tag && (a.TagPrefix == sub || b.TagPrefix == sub)
	// Or if prefixes absent on one side, accept sub+"/"+plainTag equality against other's raw.
	if a.Tag != "" && b.Tag != "" && a.Tag == b.Tag && (strings.EqualFold(a.TagPrefix, sub) || strings.EqualFold(b.TagPrefix, sub)) {
		return true
	}
	// Compare raw strings for cases where YAML holds "sub/vX" while go.mod has "vX".
	if a.Tag != "" && strings.EqualFold(b.Raw, sub+"/"+a.Tag) {
		return true
	}
	if b.Tag != "" && strings.EqualFold(a.Raw, sub+"/"+b.Tag) {
		return true
	}
	return false
}
