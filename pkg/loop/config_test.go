package loop

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvConfig_parse(t *testing.T) {
	cases := []struct {
		name                           string
		envVars                        map[string]string
		expectError                    bool
		expectedPrometheusPort         int
		expectedTracingEnabled         bool
		expectedTracingCollectorTarget string
		expectedTracingSamplingRatio   float64
	}{
		{
			name: "All variables set correctly",
			envVars: map[string]string{
				envPromPort:                 "8080",
				envTracingEnabled:           "true",
				envTracingCollectorTarget:   "some:target",
				envTracingSamplingRatio:     "1.0",
				envTracingAttribute + "XYZ": "value",
			},
			expectError:                    false,
			expectedPrometheusPort:         8080,
			expectedTracingEnabled:         true,
			expectedTracingCollectorTarget: "some:target",
			expectedTracingSamplingRatio:   1.0,
		},
		{
			name: "CL_PROMETHEUS_PORT parse error",
			envVars: map[string]string{
				envPromPort: "abc",
			},
			expectError: true,
		},
		{
			name: "TRACING_ENABLED parse error",
			envVars: map[string]string{
				envPromPort:       "8080",
				envTracingEnabled: "invalid_bool",
			},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVars {
				t.Setenv(k, v)
			}

			var config EnvConfig
			err := config.parse()

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					if config.PrometheusPort != tc.expectedPrometheusPort {
						t.Errorf("Expected Prometheus port %d, got %d", tc.expectedPrometheusPort, config.PrometheusPort)
					}
					if config.TracingEnabled != tc.expectedTracingEnabled {
						t.Errorf("Expected tracingEnabled %v, got %v", tc.expectedTracingEnabled, config.TracingEnabled)
					}
					if config.TracingCollectorTarget != tc.expectedTracingCollectorTarget {
						t.Errorf("Expected tracingCollectorTarget %s, got %s", tc.expectedTracingCollectorTarget, config.TracingCollectorTarget)
					}
					if config.TracingSamplingRatio != tc.expectedTracingSamplingRatio {
						t.Errorf("Expected tracingSamplingRatio %f, got %f", tc.expectedTracingSamplingRatio, config.TracingSamplingRatio)
					}
				}
			}
		})
	}
}

func TestEnvConfig_AsCmdEnv(t *testing.T) {
	envCfg := EnvConfig{
		PrometheusPort:         9090,
		TracingEnabled:         true,
		TracingCollectorTarget: "http://localhost:9000",
		TracingSamplingRatio:   0.1,
		TracingAttributes:      map[string]string{"key": "value"},
	}
	got := map[string]string{}
	for _, kv := range envCfg.AsCmdEnv() {
		pair := strings.SplitN(kv, "=", 2)
		require.Len(t, pair, 2)
		got[pair[0]] = pair[1]
	}

	assert.Equal(t, strconv.Itoa(9090), got[envPromPort])
	assert.Equal(t, "true", got[envTracingEnabled])
	assert.Equal(t, "http://localhost:9000", got[envTracingCollectorTarget])
	assert.Equal(t, "0.1", got[envTracingSamplingRatio])
	assert.Equal(t, "value", got[envTracingAttribute+"key"])
}
