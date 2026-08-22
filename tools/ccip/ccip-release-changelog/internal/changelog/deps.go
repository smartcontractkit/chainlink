package changelog

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

// PluginPin is one tracked entry of plugins.public.yaml.
type PluginPin struct {
	GitRef    string
	ModuleURI string
}

// DepSnapshot captures the CCIP-relevant dependency pins of the core repo at
// one git ref.
type DepSnapshot struct {
	Ref string
	SHA string
	// Modules maps go.mod module path -> version string (tracked modules only).
	Modules map[string]string
	// Plugins maps plugins.public.yaml key -> plugin pin (tracked plugins only).
	Plugins map[string]PluginPin
}

var (
	// pseudoVersionSHA matches the trailing 12-hex-char commit SHA of a Go
	// pseudo-version, e.g. "v0.1.1-solana.0.20260625091148-e5618f5682ee".
	pseudoVersionSHA = regexp.MustCompile(`-([0-9a-f]{12})$`)
	// rawSHA matches a full 40-char git SHA used directly as a gitRef.
	rawSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// VersionSHA extracts a commit SHA from a Go pseudo-version or raw-SHA
// gitRef. Returns "" if no SHA can be extracted (e.g. a clean release tag).
func VersionSHA(version string) string {
	if rawSHA.MatchString(version) {
		return version
	}
	if m := pseudoVersionSHA.FindStringSubmatch(version); m != nil {
		return m[1]
	}
	return ""
}

// ParseGoMod parses a root go.mod and extracts the versions of all tracked
// modules (per TrackedRepos configuration).
func ParseGoMod(data []byte) (map[string]string, error) {
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing go.mod: %w", err)
	}
	tracked := map[string]bool{}
	for _, repo := range TrackedRepos {
		for _, m := range repo.GoModules {
			tracked[m] = true
		}
	}
	out := map[string]string{}
	for _, req := range f.Require {
		if tracked[req.Mod.Path] {
			out[req.Mod.Path] = req.Mod.Version
		}
	}
	return out, nil
}

// pluginsFile mirrors the relevant structure of plugins/plugins.public.yaml.
type pluginsFile struct {
	Plugins map[string][]struct {
		ModuleURI string `yaml:"moduleURI"`
		GitRef    string `yaml:"gitRef"`
	} `yaml:"plugins"`
}

// ParsePluginsYAML parses plugins.public.yaml and extracts the pins of all
// tracked plugin keys.
func ParsePluginsYAML(data []byte) (map[string]PluginPin, error) {
	var pf pluginsFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parsing plugins.public.yaml: %w", err)
	}
	tracked := map[string]bool{}
	for _, repo := range TrackedRepos {
		for _, k := range repo.PluginKeys {
			tracked[k] = true
		}
	}
	out := map[string]PluginPin{}
	for key, entries := range pf.Plugins {
		if !tracked[key] || len(entries) == 0 {
			continue
		}
		out[key] = PluginPin{GitRef: entries[0].GitRef, ModuleURI: entries[0].ModuleURI}
	}
	return out, nil
}

// LoadSnapshot resolves ref and loads the dependency snapshot at that ref
// from the local core-repo checkout.
func LoadSnapshot(ctx context.Context, g gitRunner, ref string) (DepSnapshot, error) {
	sha, err := g.ResolveRef(ctx, ref)
	if err != nil {
		return DepSnapshot{}, err
	}
	// Read files at the resolved SHA, not the raw ref: the ref may have
	// resolved via the origin/<ref> fallback and not exist locally.
	goMod, err := g.FileAtRef(ctx, sha, "go.mod")
	if err != nil {
		return DepSnapshot{}, err
	}
	plugins, err := g.FileAtRef(ctx, sha, "plugins/plugins.public.yaml")
	if err != nil {
		return DepSnapshot{}, err
	}
	modules, err := ParseGoMod(goMod)
	if err != nil {
		return DepSnapshot{}, err
	}
	pluginRefs, err := ParsePluginsYAML(plugins)
	if err != nil {
		return DepSnapshot{}, err
	}
	return DepSnapshot{
		Ref: ref, SHA: sha,
		Modules: modules, Plugins: pluginRefs,
	}, nil
}

// repoPin holds the resolved pins of one tracked repo at one ref.
type repoPin struct {
	// PrimaryVersion is the version string driving the changelog (plugin
	// gitRef if the repo has plugin entries, else the primary go.mod module).
	PrimaryVersion string
	// PrimarySHA is the extracted commit SHA ("" if unparseable).
	PrimarySHA string
	// ModuleVersions maps module path -> version at this ref.
	ModuleVersions map[string]string
	// PluginRefs maps plugin key -> gitRef at this ref.
	PluginRefs map[string]PluginPin
}

func pinFor(cfg RepoConfig, snap DepSnapshot) repoPin {
	p := repoPin{
		ModuleVersions: map[string]string{},
		PluginRefs:     map[string]PluginPin{},
	}
	for _, m := range cfg.GoModules {
		if v, ok := snap.Modules[m]; ok {
			p.ModuleVersions[m] = v
		}
	}
	for _, k := range cfg.PluginKeys {
		if v, ok := snap.Plugins[k]; ok {
			p.PluginRefs[k] = v
		}
	}
	switch {
	case cfg.primaryIsPlugin() && len(p.PluginRefs) > 0:
		p.PrimaryVersion = p.PluginRefs[cfg.PluginKeys[0]].GitRef
	case len(cfg.GoModules) > 0:
		p.PrimaryVersion = p.ModuleVersions[cfg.GoModules[0]]
	}
	p.PrimarySHA = VersionSHA(p.PrimaryVersion)
	return p
}

// shortSHA trims a SHA for display.
func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}
