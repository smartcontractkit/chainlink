package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
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
		{"CACertFileSet", toml.Telemetry{CACertFile: new("test.pem")}, "test.pem"},
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
		{"EndpointSet", toml.Telemetry{Endpoint: new("localhost:4317")}, "localhost:4317"},
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
		{"TraceSampleRatioSet", toml.Telemetry{TraceSampleRatio: new(0.5)}, 0.5},
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
		{"EmitterBatchProcessorTrue", toml.Telemetry{EmitterBatchProcessor: new(true)}, true},
		{"EmitterBatchProcessorFalse", toml.Telemetry{EmitterBatchProcessor: new(false)}, false},
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
		{"ChipIngressEndpointSet", toml.Telemetry{ChipIngressEndpoint: new("localhost:8080")}, "localhost:8080"},
		{"ChipIngressEndpointNil", toml.Telemetry{ChipIngressEndpoint: nil}, ""},
		{"ChipIngressEndpointEmpty", toml.Telemetry{ChipIngressEndpoint: new("")}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressEndpoint())
		})
	}
}

func TestTelemetryConfig_ChipIngressInsecureConnection(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"ChipIngressInsecureConnectionTrue", toml.Telemetry{ChipIngressInsecureConnection: new(true)}, true},
		{"ChipIngressInsecureConnectionFalse", toml.Telemetry{ChipIngressInsecureConnection: new(false)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressInsecureConnection())
		})
	}
}

func TestTelemetryConfig_ChipIngressBatchEmitterEnabled(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"DefaultNil", toml.Telemetry{ChipIngressBatchEmitterEnabled: nil}, true},
		{"ExplicitTrue", toml.Telemetry{ChipIngressBatchEmitterEnabled: new(true)}, true},
		{"ExplicitFalse", toml.Telemetry{ChipIngressBatchEmitterEnabled: new(false)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressBatchEmitterEnabled())
		})
	}
}

func TestTelemetryConfig_DurableEmitterMaxQueuePayloadBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int64
	}{
		{"Set", toml.Telemetry{DurableEmitterMaxQueuePayloadBytes: new(int64(2048))}, 2048},
		{"Nil", toml.Telemetry{DurableEmitterMaxQueuePayloadBytes: nil}, 1 << 30}, // Default 1 GiB
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.DurableEmitterMaxQueuePayloadBytes())
		})
	}
}

func ptrDuration(d time.Duration) *config.Duration {
	return config.MustNewDuration(d)
}

func TestTelemetryConfig_ChipIngressBufferSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  uint
	}{
		{"ChipIngressBufferSizeSet", toml.Telemetry{ChipIngressBufferSize: new(uint(1000))}, 1000},
		{"ChipIngressBufferSizeNil", toml.Telemetry{ChipIngressBufferSize: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressBufferSize())
		})
	}
}

func TestTelemetryConfig_ChipIngressMaxBatchSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  uint
	}{
		{"ChipIngressMaxBatchSizeSet", toml.Telemetry{ChipIngressMaxBatchSize: new(uint(500))}, 500},
		{"ChipIngressMaxBatchSizeNil", toml.Telemetry{ChipIngressMaxBatchSize: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressMaxBatchSize())
		})
	}
}

func TestTelemetryConfig_ChipIngressMaxConcurrentSends(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int
	}{
		{"ChipIngressMaxConcurrentSendsSet", toml.Telemetry{ChipIngressMaxConcurrentSends: new(10)}, 10},
		{"ChipIngressMaxConcurrentSendsNil", toml.Telemetry{ChipIngressMaxConcurrentSends: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressMaxConcurrentSends())
		})
	}
}

func TestTelemetryConfig_ChipIngressSendInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"ChipIngressSendIntervalSet", toml.Telemetry{ChipIngressSendInterval: ptrDuration(100 * time.Millisecond)}, 100 * time.Millisecond},
		{"ChipIngressSendIntervalNil", toml.Telemetry{ChipIngressSendInterval: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressSendInterval())
		})
	}
}

func TestTelemetryConfig_ChipIngressSendTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"ChipIngressSendTimeoutSet", toml.Telemetry{ChipIngressSendTimeout: ptrDuration(3 * time.Second)}, 3 * time.Second},
		{"ChipIngressSendTimeoutNil", toml.Telemetry{ChipIngressSendTimeout: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressSendTimeout())
		})
	}
}

func TestTelemetryConfig_ChipIngressDrainTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"ChipIngressDrainTimeoutSet", toml.Telemetry{ChipIngressDrainTimeout: ptrDuration(10 * time.Second)}, 10 * time.Second},
		{"ChipIngressDrainTimeoutNil", toml.Telemetry{ChipIngressDrainTimeout: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressDrainTimeout())
		})
	}
}

func TestTelemetryConfig_ChipIngressMaxGRPCRequestSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int
	}{
		{"ChipIngressMaxGRPCRequestSizeSet", toml.Telemetry{ChipIngressMaxGRPCRequestSize: new(10485760)}, 10485760},
		{"ChipIngressMaxGRPCRequestSizeNil", toml.Telemetry{ChipIngressMaxGRPCRequestSize: nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.ChipIngressMaxGRPCRequestSize())
		})
	}
}

func TestTelemetryConfig_HeartbeatInterval(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"HeartbeatIntervalSet", toml.Telemetry{HeartbeatInterval: ptrDuration(5 * time.Second)}, 5 * time.Second},
		{"HeartbeatIntervalNil", toml.Telemetry{HeartbeatInterval: nil}, 1 * time.Second},             // Default value
		{"HeartbeatIntervalZero", toml.Telemetry{HeartbeatInterval: ptrDuration(0)}, 1 * time.Second}, // Zero value results in default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.HeartbeatInterval())
		})
	}
}

func TestTelemetryConfig_LogStreamingEnabled(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"LogStreamingEnabledTrue", toml.Telemetry{LogStreamingEnabled: new(true)}, true},
		{"LogStreamingEnabledFalse", toml.Telemetry{LogStreamingEnabled: new(false)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogStreamingEnabled())
		})
	}
}

func TestTelemetryConfig_LogLevel(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  zapcore.Level
	}{
		{"LogLevelSet", toml.Telemetry{LogLevel: new("debug")}, zapcore.DebugLevel},
		{"LogLevelInfo", toml.Telemetry{LogLevel: new("info")}, zapcore.InfoLevel},
		{"LogLevelWarn", toml.Telemetry{LogLevel: new("warn")}, zapcore.WarnLevel},
		{"LogLevelError", toml.Telemetry{LogLevel: new("error")}, zapcore.ErrorLevel},
		{"LogLevelNil", toml.Telemetry{LogLevel: nil}, zapcore.InfoLevel},
		{"LogLevelInvalid", toml.Telemetry{LogLevel: new("invalid")}, zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogLevel())
		})
	}
}

func TestTelemetryConfig_LogBatchProcessor(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  bool
	}{
		{"LogBatchProcessorTrue", toml.Telemetry{LogBatchProcessor: new(true)}, true},
		{"LogBatchProcessorFalse", toml.Telemetry{LogBatchProcessor: new(false)}, false},
		{"LogBatchProcessorNil", toml.Telemetry{LogBatchProcessor: nil}, true}, // Default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogBatchProcessor())
		})
	}
}

func TestTelemetryConfig_LogExportTimeout(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"LogExportTimeoutSet", toml.Telemetry{LogExportTimeout: ptrDuration(5 * time.Second)}, 5 * time.Second},
		{"LogExportTimeoutNil", toml.Telemetry{LogExportTimeout: nil}, 1 * time.Second}, // Default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogExportTimeout())
		})
	}
}
func TestTelemetryConfig_LogExportMaxBatchSize(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int
	}{
		{"LogExportMaxBatchSizeSet", toml.Telemetry{LogExportMaxBatchSize: new(512)}, 512},
		{"LogExportMaxBatchSizeNil", toml.Telemetry{LogExportMaxBatchSize: nil}, 512}, // Default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogExportMaxBatchSize())
		})
	}
}

func TestTelemetryConfig_LogExportInterval(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  time.Duration
	}{
		{"LogExportIntervalSet", toml.Telemetry{LogExportInterval: ptrDuration(5 * time.Second)}, 5 * time.Second},
		{"LogExportIntervalNil", toml.Telemetry{LogExportInterval: nil}, 1 * time.Second}, // Default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogExportInterval())
		})
	}
}

func TestTelemetryConfig_LogMaxQueueSize(t *testing.T) {
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int
	}{
		{"LogMaxQueueSizeSet", toml.Telemetry{LogMaxQueueSize: new(2048)}, 2048},
		{"LogMaxQueueSizeNil", toml.Telemetry{LogMaxQueueSize: nil}, 2048}, // Default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.LogMaxQueueSize())
		})
	}
}

func TestTelemetryConfig_MetricViewsDenyAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  []string
	}{
		{"DenylistSet", toml.Telemetry{MetricViewsDenyAttributes: []string{"event_id"}}, []string{"event_id"}},
		{"DenylistNil", toml.Telemetry{MetricViewsDenyAttributes: nil}, nil},
		{"DenylistEmpty", toml.Telemetry{MetricViewsDenyAttributes: []string{}}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.MetricViewsDenyAttributes())
		})
	}
}

func TestTelemetryConfig_MetricCardinalityLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		telemetry toml.Telemetry
		expected  int
	}{
		{"MetricCardinalityLimitSet", toml.Telemetry{MetricCardinalityLimit: new(500)}, 500},
		{"MetricCardinalityLimitZero", toml.Telemetry{MetricCardinalityLimit: new(0)}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tc := telemetryConfig{s: tt.telemetry}
			assert.Equal(t, tt.expected, tc.MetricCardinalityLimit())
		})
	}

	t.Run("MetricCardinalityLimitDefaultFromCore", func(t *testing.T) {
		t.Parallel()
		defaults := docs.CoreDefaults()
		tc := telemetryConfig{s: defaults.Telemetry}
		assert.Equal(t, 100000, tc.MetricCardinalityLimit())
	})
}
