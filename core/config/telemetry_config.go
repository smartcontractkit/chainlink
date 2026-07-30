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
	ChipIngressBatchEmitterEnabled() bool
	ChipIngressBufferSize() uint
	ChipIngressMaxBatchSize() uint
	ChipIngressMaxConcurrentSends() int
	ChipIngressSendInterval() time.Duration
	ChipIngressSendTimeout() time.Duration
	ChipIngressDrainTimeout() time.Duration
	ChipIngressMaxGRPCRequestSize() int
	DurableEmitterEnabled() bool
	DurableEmitterRetransmitBatchSize() int
	DurableEmitterEventTTL() time.Duration
	DurableEmitterMaxQueuePayloadBytes() int64
	DurableEmitterInsertBatchFlushInterval() time.Duration
	HeartbeatInterval() time.Duration
	LogStreamingEnabled() bool
	LogLevel() zapcore.Level
	LogBatchProcessor() bool
	LogExportTimeout() time.Duration
	LogExportMaxBatchSize() int
	LogExportInterval() time.Duration
	LogMaxQueueSize() int
	MetricViewsDenyAttributes() []string
	MetricCardinalityLimit() int
	PrometheusBridge() PrometheusBridge
}

type PrometheusBridge interface {
	Enabled() bool
	Prefixes() []string
}
