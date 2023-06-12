package resolver

import "github.com/smartcontractkit/chainlink/v2/core/config"

type FeaturesResolver struct {
	cfg config.Feature
}

func NewFeaturesResolver(cfg config.Feature) *FeaturesResolver {
	return &FeaturesResolver{cfg: cfg}
}

// CSA resolves to whether CSA Keys are enabled
func (r *FeaturesResolver) CSA() bool {
	return r.cfg.UICSAKeys()
}

// FeedsManager resolves to whether the Feeds Manager is enabled for the UI
func (r *FeaturesResolver) FeedsManager() bool {
	return r.cfg.FeedsManager()
}

type FeaturesPayloadResolver struct {
	cfg config.Feature
}

func NewFeaturesPayloadResolver(cfg config.Feature) *FeaturesPayloadResolver {
	return &FeaturesPayloadResolver{cfg: cfg}
}

func (r *FeaturesPayloadResolver) ToFeatures() (*FeaturesResolver, bool) {
	return NewFeaturesResolver(r.cfg), true
}
