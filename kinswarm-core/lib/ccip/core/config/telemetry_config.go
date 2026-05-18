package config

type Telemetry interface {
	Enabled() bool
	InsecureConnection() bool
	CACertFile() string
	OtelExporterGRPCEndpoint() string
	ResourceAttributes() map[string]string
	TraceSampleRatio() float64
}
