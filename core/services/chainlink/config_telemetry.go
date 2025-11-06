package chainlink

import (
	"fmt"
	"maps"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

const defaultHeartbeatInterval = 1 * time.Second

type telemetryConfig struct {
	s toml.Telemetry
}

func (b *telemetryConfig) Enabled() bool { return *b.s.Enabled }

func (b *telemetryConfig) InsecureConnection() bool {
	if b.s.InsecureConnection == nil {
		return false
	}
	return *b.s.InsecureConnection
}

func (b *telemetryConfig) CACertFile() string {
	if b.s.CACertFile == nil {
		return ""
	}
	return *b.s.CACertFile
}

func (b *telemetryConfig) OtelExporterGRPCEndpoint() string {
	if b.s.Endpoint == nil {
		return ""
	}
	return *b.s.Endpoint
}

// ResourceAttributes returns the resource attributes set in the TOML config
// by the user, but first sets OTEL required attributes:
//
//	service.name
//	service.version
//
// These can be overridden by the TOML if the user so chooses
func (b *telemetryConfig) ResourceAttributes() map[string]string {
	sha, ver := static.Short()

	defaults := map[string]string{
		"service.name":         "chainlink",
		"service.version":      static.Version,
		"service.sha":          static.Sha,
		"service.shortversion": fmt.Sprintf("%s@%s", ver, sha),
	}

	maps.Copy(defaults, b.s.ResourceAttributes)

	return defaults
}

func (b *telemetryConfig) TraceSampleRatio() float64 {
	if b.s.TraceSampleRatio == nil {
		return 0.0
	}
	return *b.s.TraceSampleRatio
}

func (b *telemetryConfig) EmitterBatchProcessor() bool {
	if b.s.EmitterBatchProcessor == nil {
		return false
	}
	return *b.s.EmitterBatchProcessor
}

func (b *telemetryConfig) EmitterExportTimeout() time.Duration {
	if b.s.EmitterExportTimeout == nil {
		return 0
	}
	return b.s.EmitterExportTimeout.Duration()
}

func (b *telemetryConfig) ChipIngressEndpoint() string {
	if b.s.ChipIngressEndpoint == nil {
		return ""
	}
	return *b.s.ChipIngressEndpoint
}

func (b *telemetryConfig) ChipIngressInsecureConnection() bool {
	if b.s.ChipIngressInsecureConnection == nil {
		return false
	}
	return *b.s.ChipIngressInsecureConnection
}

func (b *telemetryConfig) HeartbeatInterval() time.Duration {
	if b.s.HeartbeatInterval == nil || b.s.HeartbeatInterval.Duration() <= 0 {
		return defaultHeartbeatInterval
	}
	return b.s.HeartbeatInterval.Duration()
}

func (b *telemetryConfig) LogStreamingEnabled() bool {
	if b.s.LogStreamingEnabled == nil {
		return false
	}
	return *b.s.LogStreamingEnabled
}

func (b *telemetryConfig) LogLevel() zapcore.Level {
	if b.s.LogLevel == nil {
		return zapcore.InfoLevel // Default log level
	}

	var level zapcore.Level
	if err := level.Set(*b.s.LogLevel); err != nil {
		return zapcore.InfoLevel // Fallback to info level on invalid input
	}
	return level
}

func (b *telemetryConfig) LogBatchProcessor() bool {
	if b.s.LogBatchProcessor == nil {
		return true
	}
	return *b.s.LogBatchProcessor
}

func (b *telemetryConfig) LogExportTimeout() time.Duration {
	if b.s.LogExportTimeout == nil {
		return 1 * time.Second
	}
	return b.s.LogExportTimeout.Duration()
}

func (b *telemetryConfig) LogExportMaxBatchSize() int {
	if b.s.LogExportMaxBatchSize == nil {
		return 512
	}
	return *b.s.LogExportMaxBatchSize
}

func (b *telemetryConfig) LogExportInterval() time.Duration {
	if b.s.LogExportInterval == nil {
		return 1 * time.Second
	}
	return b.s.LogExportInterval.Duration()
}

func (b *telemetryConfig) LogMaxQueueSize() int {
	if b.s.LogMaxQueueSize == nil {
		return 2048
	}
	return *b.s.LogMaxQueueSize
}
