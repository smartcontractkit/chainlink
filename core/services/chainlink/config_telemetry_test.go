package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func TestTelemetryConfig_Enabled(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"EnabledTrue", toml.Telemetry{Enabled: &trueVal}, true},
		{"EnabledFalse", toml.Telemetry{Enabled: &falseVal}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.Enabled())
		})
	}
}

func TestTelemetryConfig_InsecureConnection(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"InsecureConnectionTrue", toml.Telemetry{InsecureConnection: &trueVal}, true},
		{"InsecureConnectionFalse", toml.Telemetry{InsecureConnection: &falseVal}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.InsecureConnection())
		})
	}
}

func TestTelemetryConfig_CACertFile(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  string
	}{
		{"CACertFileSet", toml.Telemetry{CACertFile: ptr("test.pem")}, "test.pem"},
		{"CACertFileNil", toml.Telemetry{CACertFile: nil}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.CACertFile())
		})
	}
}

func TestTelemetryConfig_OtelExporterGRPCEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  string
	}{
		{"EndpointSet", toml.Telemetry{Endpoint: ptr("localhost:4317")}, "localhost:4317"},
		{"EndpointNil", toml.Telemetry{Endpoint: nil}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.OtelExporterGRPCEndpoint())
		})
	}
}

func TestTelemetryConfig_ResourceAttributes(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  map[string]string
	}{
		{
			"DefaultAttributes",
			toml.Telemetry{ResourceAttributes: nil},
			map[string]string{
				"service.name":         "chainlink",
				"service.sha":          "unset",
				"service.shortversion": "unset@unset",
				"service.version":      static.Version,
			},
		},
		{
			"CustomAttributes",
			toml.Telemetry{ResourceAttributes: map[string]string{"custom.key": "custom.value"}},
			map[string]string{
				"service.name":         "chainlink",
				"service.sha":          "unset",
				"service.shortversion": "unset@unset",
				"service.version":      static.Version,
				"custom.key":           "custom.value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ResourceAttributes())
		})
	}
}

func TestTelemetryConfig_TraceSampleRatio(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  float64
	}{
		{"TraceSampleRatioSet", toml.Telemetry{TraceSampleRatio: ptrFloat(0.5)}, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.InEpsilon(t, tt.expected, tc.TraceSampleRatio(), 0.0001)
		})
	}
}

func TestTelemetryConfig_EmitterBatchProcessor(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"EmitterBatchProcessorTrue", toml.Telemetry{EmitterBatchProcessor: ptr(true)}, true},
		{"EmitterBatchProcessorFalse", toml.Telemetry{EmitterBatchProcessor: ptr(false)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.EmitterBatchProcessor())
		})
	}
}

func TestTelemetryConfig_EmitterExportTimeout(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"EmitterExportTimeoutSet", toml.Telemetry{EmitterExportTimeout: ptrDuration(5 * time.Second)}, 5 * time.Second},
		{"EmitterExportTimeoutNil", toml.Telemetry{EmitterExportTimeout: nil}, 0},
		{"EmitterExportTimeoutZero", toml.Telemetry{EmitterExportTimeout: ptrDuration(0)}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.EmitterExportTimeout())
		})
	}
}

func TestTelemetryConfig_ChipIngressEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  string
	}{
		{"ChipIngressEndpointSet", toml.Telemetry{ChipIngressEndpoint: ptr("localhost:8080")}, "localhost:8080"},
		{"ChipIngressEndpointNil", toml.Telemetry{ChipIngressEndpoint: nil}, ""},
		{"ChipIngressEndpointEmpty", toml.Telemetry{ChipIngressEndpoint: ptr("")}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressEndpoint())
		})
	}
}

func ptrDuration(d time.Duration) *config.Duration {
	return config.MustNewDuration(d)
}

func ptrFloat(f float64) *float64 {
	return &f
}
