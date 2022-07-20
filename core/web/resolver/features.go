package resolver

import "github.com/smartcontractkit/chainlink/core/config"

type FeaturesResolver struct {
	cfg config.GeneralConfig
}

func NewFeaturesResolver(cfg config.GeneralConfig) *FeaturesResolver {
	return &FeaturesResolver{cfg: cfg}
}

// CSA resolves to whether CSA Keys are enabled
func (r *FeaturesResolver) CSA() bool {
	return r.cfg.FeatureUICSAKeys()
}

// FeedsManager resolves to whether the Feeds Manager is enabled for the UI
func (r *FeaturesResolver) FeedsManager() bool {
	return r.cfg.FeatureFeedsManager()
}

type FeaturesPayloadResolver struct {
	cfg config.GeneralConfig
}

func NewFeaturesPayloadResolver(cfg config.GeneralConfig) *FeaturesPayloadResolver {
	return &FeaturesPayloadResolver{cfg: cfg}
}

func (r *FeaturesPayloadResolver) ToFeatures() (*FeaturesResolver, bool) {
	return NewFeaturesResolver(r.cfg), true
}
