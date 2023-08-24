package config

// CommitPluginJobSpecConfig contains the plugin specific variables for the ccip.CCIPCommit plugin.
// We use ID here to keep it as general as possible, e.g. abstracting for chains which don't have an address concept.
type CommitPluginJobSpecConfig struct {
	SourceStartBlock, DestStartBlock int64  // Only for first time job add.
	OffRamp                          string `json:"offRamp"`
	// TokenPricesUSDPipeline should contain a token price pipeline for the following tokens:
	//		The SOURCE chain wrapped native
	// 		The DESTINATION supported tokens (including fee tokens) as defined in destination OffRamp and PriceRegistry.
	TokenPricesUSDPipeline string `json:"tokenPricesUSDPipeline"`
}

// ExecutionPluginJobSpecConfig contains the plugin specific variables for the ccip.CCIPExecution plugin.
type ExecutionPluginJobSpecConfig struct {
	SourceStartBlock, DestStartBlock int64 // Only for first time job add.
}
