package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

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

func (b *telemetryConfig) ResourceAttributes() map[string]string {
	return b.s.ResourceAttributes
}

func (b *telemetryConfig) TraceSampleRatio() float64 {
	if b.s.TraceSampleRatio == nil {
		return 0.0
	}
	return *b.s.TraceSampleRatio
}
