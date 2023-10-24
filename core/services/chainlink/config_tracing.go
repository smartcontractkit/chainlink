package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type tracingConfig struct {
	s toml.Tracing
}

func (t tracingConfig) Enabled() bool {
	return *t.s.Enabled
}

func (t tracingConfig) CollectorTarget() string {
	return *t.s.CollectorTarget
}

func (t tracingConfig) NodeID() string {
	return *t.s.NodeID
}

func (t tracingConfig) Attributes() map[string]string {
	return t.s.Attributes
}

func (t tracingConfig) SamplingRatio() float64 {
	return *t.s.SamplingRatio
}
