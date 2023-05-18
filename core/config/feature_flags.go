package config

// FeatureFlags contains bools that toggle various features or chains
// TODO: document the new ones
type FeatureFlags interface {
	FeatureExternalInitiators() bool
	FeatureFeedsManager() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	FeatureUICSAKeys() bool
	FeatureLogPoller() bool

	AutoPprofEnabled() bool
	EVMEnabled() bool
	EVMRPCEnabled() bool
	P2PEnabled() bool
	CosmosEnabled() bool
	SolanaEnabled() bool
	StarkNetEnabled() bool
}
