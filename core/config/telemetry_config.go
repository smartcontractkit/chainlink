package config

import (
	"time"

	"go.uber.org/zap/zapcore"
)

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
	ChipIngressInsecureConnection() bool
	HeartbeatInterval() time.Duration
	LogStreamingEnabled() bool
	LogLevel() zapcore.Level
	LogBatchProcessor() bool
	LogExportTimeout() time.Duration
	LogExportMaxBatchSize() int
	LogExportInterval() time.Duration
	LogMaxQueueSize() int
}
