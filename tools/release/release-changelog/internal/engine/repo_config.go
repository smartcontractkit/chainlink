package engine

// RepoConfig describes one repository tracked by a product's changelog.
type RepoConfig struct {
	// Name is the repository name (e.g. "chainlink-ccip").
	Name string
	// Owner is the GitHub org/user (e.g. "smartcontractkit").
	Owner string

	// GoModules lists module paths in the root go.mod of the core repo that
	// come from this repository, in decreasing importance order. The first
	// entry is the "primary module" whose pin drives the commit changelog
	// when the repo has no plugin entry.
	GoModules []string

	// PluginKeys lists keys in plugins/plugins.public.yaml that install from
	// this repository. When set, the plugin gitRef is what gets built into
	// the release image, so it is the primary pin for the commit changelog.
	PluginKeys []string

	// IncludePaths, when non-empty, restricts the commit changelog to commits
	// touching at least one of these path prefixes.
	IncludePaths []string
	// ExcludePaths drops commits that only touch these path prefixes.
	ExcludePaths []string

	// Local indicates the repository is the one this tool runs inside
	// (the core chainlink repo). Its commit log is read from the local git
	// checkout instead of the GitHub compare API.
	Local bool
}

// primaryIsPlugin reports whether the plugin gitRef (rather than a go.mod
// module pin) is the authoritative pin for this repo.
func (c RepoConfig) primaryIsPlugin() bool { return len(c.PluginKeys) > 0 }
