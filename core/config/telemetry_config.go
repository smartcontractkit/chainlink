package config

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Telemetry interface {
	AuthHeadersTTL() time.Duration
	Enabled() bool
	InsecureConnection() bool
	CACertFile() string
	OtelExporterGRPCEndpoint() string
	ResourceAttributes() map[string]string
	TraceSampleRatio() float64
	EmitterBatchProcessor() bool
	EmitterExportTimeout() time.Duration
	ChipIngressEndpoint() string
	ChipIngressInsecureConnection() bool
	DurableEmitterEnabled() bool
	// DurableEmitterPersistSources lists CloudEvent Source values (beholder_domain) that may be
	// written to the durable Chip queue. See chainlink telemetry config for defaults and wildcards.
	DurableEmitterPersistSources() []string
	HeartbeatInterval() time.Duration
	LogStreamingEnabled() bool
	LogLevel() zapcore.Level
	LogBatchProcessor() bool
	LogExportTimeout() time.Duration
	LogExportMaxBatchSize() int
	LogExportInterval() time.Duration
	LogMaxQueueSize() int
}
