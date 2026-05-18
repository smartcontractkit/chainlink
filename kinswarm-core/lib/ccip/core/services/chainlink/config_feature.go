package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type featureConfig struct {
	c toml.Feature
}

func (f *featureConfig) FeedsManager() bool {
	return *f.c.FeedsManager
}

func (f *featureConfig) LogPoller() bool {
	return *f.c.LogPoller
}

func (f *featureConfig) UICSAKeys() bool {
	return *f.c.UICSAKeys
}

func (f *featureConfig) CCIP() bool {
	return *f.c.CCIP
}

func (f *featureConfig) MultiFeedsManagers() bool {
	return *f.c.MultiFeedsManagers
}
