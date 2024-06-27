package config

type Tracing interface {
	Enabled() bool
	CollectorTarget() string
	NodeID() string
	Attributes() map[string]string
	SamplingRatio() float64
}
