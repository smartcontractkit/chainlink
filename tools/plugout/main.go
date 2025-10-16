package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

// -----------------------------------------------------------------------------
// Orchestration
// -----------------------------------------------------------------------------

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

	// Discover modules from plugin files and apply ignores.
	modulesToCheck := discoverModulesFromPlugins(opts.PluginPaths)
	if len(modulesToCheck) == 0 {
		return false, errors.New("no modules discovered from plugin files")
	}

	fmt.Println()
	if len(opts.IgnoreModules) > 0 {
		fmt.Printf("Ignoring %d modules as specified:\n  - %s\n", len(opts.IgnoreModules), strings.Join(opts.IgnoreModules, "\n  - "))
		modulesToCheck = without(modulesToCheck, opts.IgnoreModules)
	}
	totalModules := len(modulesToCheck)
	sort.Strings(modulesToCheck)
	fmt.Printf("Modules to check (%d): \n  - %s\n", totalModules, strings.Join(modulesToCheck, "\n  - "))
	fmt.Println()

	hasMismatch := false
	for idx, module := range modulesToCheck {
		fmt.Printf("\n---\n%s - %d/%d\n---\n", module, idx+1, totalModules) // progress
		mismatched, err := checkAndUpdateModuleVersion(module, opts)
		if err != nil {
			return hasMismatch, err
		}
		if mismatched {
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

// discoverModulesFromPlugins returns a unique list of moduleURIs declared in plugin YAMLs.
func discoverModulesFromPlugins(paths []string) []string {
	seen := make(map[string]struct{})
	for _, p := range paths {
		// keep track of modules in this file
		modulesInFile := []string{}
		data, err := os.ReadFile(p)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", p, err)
			continue
		}
		var pf PluginsFile
		if err := yaml.Unmarshal(data, &pf); err != nil {
			fmt.Printf("Warning: failed to parse YAML %s: %v\n", p, err)
			continue
		}
		for _, list := range pf.Plugins {
			for _, pl := range list {
				if pl.ModuleURI != "" {
					seen[pl.ModuleURI] = struct{}{}
					modulesInFile = append(modulesInFile, pl.ModuleURI)
				}
			}
		}
		if len(modulesInFile) == 0 {
			fmt.Printf("Warning: no modules found in %s\n", p)
		} else {
			fmt.Printf("Discovered %d modules in %s\n", len(modulesInFile), p)
			sort.Strings(modulesInFile)
			fmt.Printf("  - %s\n", strings.Join(modulesInFile, ", "))
		}
	}
	out := make([]string, 0, len(seen))
	for m := range seen {
		out = append(out, m)
	}
	return out
}

// without returns a new slice with all elements of input except those in remove.
func without(input, remove []string) []string {
	rm := make(map[string]struct{}, len(remove))
	for _, r := range remove {
		rm[r] = struct{}{}
	}
	out := make([]string, 0, len(input))
	for _, v := range input {
		if _, skip := rm[v]; !skip {
			out = append(out, v)
		}
	}
	return out
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

// -----------------------------------------------------------------------------
// Version parsing & extraction
// -----------------------------------------------------------------------------

var (
	// Examples:
	// v0.0.0-20251013133428-62ab1091a563
	// v1.2.3-0.20250102030405-abcdef123456
	pseudoWithSHARe = regexp.MustCompile(`^v\d+\.\d+\.\d+(?:-[0-9.]+)?-(?:\d{14})-g?([0-9a-f]{7,40})$`)

	// tag like v0.1.5 (allow pre-release/build suffixes)
	plainTagRe = regexp.MustCompile(`^v\d+\.\d+\.\d+([.-].*)?$`)

	// subdir-prefixed tag like sub/dir/v0.1.5
	prefixedTagRe = regexp.MustCompile(`^(.+?)/+(v\d+\.\d+\.\d+(?:[.-].*)?)$`)

	// raw SHA (7..40 hex)
	shaOnlyRe = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
)

// normalizeVersion converts any raw string (tag/pseudo/SHA) into ModuleVersion.
// Pseudoversions never set Tag; subdir-prefixed tags set Tag+TagPrefix.
func normalizeVersion(raw string) ModuleVersion {
	mv := ModuleVersion{Raw: raw}
	low := strings.ToLower(strings.TrimSpace(raw))

	// 1) Pseudoversion? (detect first; return early after extracting SHA)
	if pseudoWithSHARe.MatchString(low) {
		mv.SHA = pseudoWithSHARe.FindStringSubmatch(low)[1]
		return mv
	}
	// Fallback pseudoversion detector (suffix after last dash looks like a SHA)
	if mv.SHA == "" && strings.HasPrefix(low, "v") && strings.Count(low, "-") >= 2 {
		parts := strings.Split(low, "-")
		last := strings.TrimPrefix(parts[len(parts)-1], "g")
		if shaOnlyRe.MatchString(last) {
			mv.SHA = last
			return mv
		}
	}

	// 2) Raw SHA?
	if shaOnlyRe.MatchString(low) {
		mv.SHA = low
		return mv
	}

	// 3) Tag (non-pseudo)
	if plainTagRe.MatchString(low) {
		// subdir-prefixed tag?
		if m := prefixedTagRe.FindStringSubmatch(low); len(m) == 3 && plainTagRe.MatchString(m[2]) {
			orig := strings.TrimSpace(raw)
			if pos := strings.LastIndex(orig, "/"); pos >= 0 && pos+1 < len(orig) {
				mv.Tag = orig[pos+1:]
				mv.TagPrefix = strings.TrimSuffix(orig[:pos], "/")
				return mv
			}
		}
		// plain tag
		mv.Tag = raw
		return mv
	}

	// otherwise leave as-is (unknown form)
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

// getGoModVersion extracts the version for a specific module from go.mod at goModPath.
func getGoModVersion(goModPath, module string) (ModuleVersion, error) {
	fmt.Printf("Extracting module version for %s from %s...\n", module, goModPath)

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ModuleVersion{}, fmt.Errorf("failed to read %s: %w", goModPath, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := scanner.Text()

		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}

		// Look for "<module> <version>" on the same line.
		if strings.Contains(line, module+" ") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == module && i+1 < len(fields) {
					ver := fields[i+1]
					mv := normalizeVersion(ver)
					fmt.Printf("  - Version extracted: %s (SHA:%s Tag:%s Prefix:%s)\n", mv.Raw, mv.SHA, mv.Tag, mv.TagPrefix)
					return mv, nil
				}
			}
		}
	}

	return ModuleVersion{}, errors.New("module " + module + " not found in go.mod")
}

// getYAMLVersion extracts the gitRef for a specific module from a plugin YAML file.
func getYAMLVersion(pluginPath, module string) (ModuleVersion, error) {
	fmt.Printf("Extracting plugins version for %s from %s...\n", module, pluginPath)

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return ModuleVersion{}, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var pluginsFile PluginsFile
	if err := yaml.Unmarshal(data, &pluginsFile); err != nil {
		return ModuleVersion{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, plugins := range pluginsFile.Plugins {
		for _, plugin := range plugins {
			if plugin.ModuleURI == module {
				normalizedVersion := normalizeVersion(plugin.GitRef)
				fmt.Printf("  - Version extracted: %s (from %s)\n", normalizedVersion.toString(), plugin.GitRef)
				return normalizedVersion, nil
			}
		}
	}

	return ModuleVersion{}, errors.New("module " + module + " not found in " + pluginPath)
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

// -----------------------------------------------------------------------------
// Core compare/update
// -----------------------------------------------------------------------------

// checkAndUpdateModuleVersion compares YAML refs against go.mod for `module`.
// Returns whether any mismatches were found for this module.
func checkAndUpdateModuleVersion(module string, opts Options) (bool, error) {
	get := opts.GetModVersion
	if get == nil {
		get = getGoModVersion
	}

	goModMV, err := get(opts.GoModPath, module)
	if err != nil || goModMV.Raw == "" {
		fmt.Printf("  - ⚠️  %v\n", err)
		return false, nil // warn & skip, no mismatch
	}

	mismatchFound := false
	for _, pluginPath := range opts.PluginPaths {
		yamlMV, err := getYAMLVersion(pluginPath, module)
		if err != nil {
			fmt.Printf("  - ⚠️  %v\n", err)
			continue
		}

		if !versionsMatchForModule(module, goModMV, yamlMV) {
			mismatchFound = true
			fmt.Printf("  - ❌ MISMATCH: %s\n", module)
			if opts.Update {
				if err := updateGitRefInYAML(pluginPath, module, goModMV); err != nil {
					fmt.Printf("  - ❌ %v\n", err)
				} else {
					fmt.Printf("  - ✅ Updated gitRef in %s\n", pluginPath)
				}
			}
		} else {
			fmt.Printf("  - ✅ %s versions match in %s\n", module, pluginPath)
		}
	}
	return mismatchFound, nil
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
