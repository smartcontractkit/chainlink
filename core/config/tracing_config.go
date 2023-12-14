package config

type Tracing interface {
	Enabled() bool
	CollectorTarget() string
	NodeID() string
	SamplingRatio() float64
	TLSCertPath() string
	Mode() string
	Attributes() map[string]string
}
