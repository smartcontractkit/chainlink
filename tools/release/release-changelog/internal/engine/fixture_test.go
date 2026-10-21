package engine

// testProduct is an inline fixture mirroring the shape of the CCIP product.
// Engine tests deliberately do NOT import a real product definition, so that
// editing a product config never breaks engine tests (or golden files).
var testProduct = Product{
	DisplayName: "CCIP",
	Repos: []RepoConfig{
		{
			Name:  "chainlink-ccip",
			Owner: "smartcontractkit",
			GoModules: []string{
				"github.com/smartcontractkit/chainlink-ccip",
				"github.com/smartcontractkit/chainlink-ccip/chains/evm",
				"github.com/smartcontractkit/chainlink-ccip/chains/solana",
				"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings",
			},
		},
		{
			Name:       "chainlink-aptos",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-aptos/codec"},
			PluginKeys: []string{"aptos"},
		},
		{
			Name:       "chainlink-sui",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-sui/codec"},
			PluginKeys: []string{"sui"},
		},
		{
			Name:       "chainlink-solana",
			Owner:      "smartcontractkit",
			PluginKeys: []string{"solana"},
		},
		{
			Name:       "chainlink-ton",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-ton"},
			PluginKeys: []string{"ton"},
		},
		{
			Name:  "chainlink-evm",
			Owner: "smartcontractkit",
			GoModules: []string{
				"github.com/smartcontractkit/chainlink-evm",
				"github.com/smartcontractkit/chainlink-evm/gethwrappers",
			},
			PluginKeys: []string{"evm"},
		},
		{
			Name:         "chainlink",
			Owner:        "smartcontractkit",
			IncludePaths: []string{"core/capabilities/ccip/"},
			Local:        true,
		},
	},
}
