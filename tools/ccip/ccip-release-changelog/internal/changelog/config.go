package changelog

// RepoConfig describes one repository tracked by the changelog tool.
//
// EDIT THIS LIST to add/remove repositories or tune which paths are tracked.
// Everything below is configuration: no CLI flags or Slack inputs are needed
// to change tracking behavior.
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

// TrackedRepos is the editable configuration block for the tool.
//
// To track an additional repo, add an entry. To narrow a repo's changelog to
// specific paths (like the core repo entry), set IncludePaths/ExcludePaths.
var TrackedRepos = []RepoConfig{
	{
		Name:  "chainlink-ccip",
		Owner: "smartcontractkit",
		GoModules: []string{
			"github.com/smartcontractkit/chainlink-ccip", // primary
			"github.com/smartcontractkit/chainlink-ccip/chains/evm",
			"github.com/smartcontractkit/chainlink-ccip/chains/solana",
			"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings",
		},
	},
	{
		Name:       "chainlink-aptos",
		Owner:      "smartcontractkit",
		GoModules:  []string{"github.com/smartcontractkit/chainlink-aptos/codec"},
		PluginKeys: []string{"aptos"}, // primary pin
	},
	{
		Name:       "chainlink-sui",
		Owner:      "smartcontractkit",
		GoModules:  []string{"github.com/smartcontractkit/chainlink-sui/codec"},
		PluginKeys: []string{"sui"}, // primary pin
	},
	{
		Name:       "chainlink-solana",
		Owner:      "smartcontractkit",
		PluginKeys: []string{"solana"}, // primary pin (not in root go.mod)
	},
	{
		Name:       "chainlink-ton",
		Owner:      "smartcontractkit",
		GoModules:  []string{"github.com/smartcontractkit/chainlink-ton"},
		PluginKeys: []string{"ton"}, // dual-source: checked against go.mod
	},
	{
		Name:  "chainlink-evm",
		Owner: "smartcontractkit",
		GoModules: []string{
			"github.com/smartcontractkit/chainlink-evm",
			"github.com/smartcontractkit/chainlink-evm/gethwrappers",
			// contracts/cre/gobindings intentionally not tracked (CRE-owned)
		},
		PluginKeys: []string{"evm"}, // dual-source: checked against go.mod
	},
	{
		// The core repo itself, restricted to the CCIP capability tree.
		Name:         "chainlink",
		Owner:        "smartcontractkit",
		IncludePaths: []string{"core/capabilities/ccip/"},
		ExcludePaths: []string{},
		Local:        true,
	},
}
