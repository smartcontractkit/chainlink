package config

import "time"

type Telemetry interface {
	Enabled() bool
	InsecureConnection() bool
	CACertFile() string
	OtelExporterGRPCEndpoint() string
	ResourceAttributes() map[string]string
	TraceSampleRatio() float64
	EmitterBatchProcessor() bool
	EmitterExportTimeout() time.Duration
	ChipIngressEndpoint() string
}
