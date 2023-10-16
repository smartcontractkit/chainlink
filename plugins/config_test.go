package plugins

import (
	"testing"
)

func TestGetEnvConfig(t *testing.T) {
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
				"CL_PROMETHEUS_PORT":       "8080",
				"TRACING_ENABLED":          "true",
				"TRACING_COLLECTOR_TARGET": "some:target",
				"TRACING_SAMPLING_RATIO":   "1.0",
				"TRACING_ATTRIBUTE_XYZ":    "value",
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
				"CL_PROMETHEUS_PORT": "abc",
			},
			expectError: true,
		},
		{
			name: "TRACING_ENABLED parse error",
			envVars: map[string]string{
				"CL_PROMETHEUS_PORT": "8080",
				"TRACING_ENABLED":    "invalid_bool",
			},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVars {
				t.Setenv(k, v)
			}

			config, err := GetEnvConfig()

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					if config.PrometheusPort() != tc.expectedPrometheusPort {
						t.Errorf("Expected Prometheus port %d, got %d", tc.expectedPrometheusPort, config.PrometheusPort())
					}
					if config.TracingEnabled() != tc.expectedTracingEnabled {
						t.Errorf("Expected tracingEnabled %v, got %v", tc.expectedTracingEnabled, config.TracingEnabled())
					}
					if config.TracingCollectorTarget() != tc.expectedTracingCollectorTarget {
						t.Errorf("Expected tracingCollectorTarget %s, got %s", tc.expectedTracingCollectorTarget, config.TracingCollectorTarget())
					}
					if config.TracingSamplingRatio() != tc.expectedTracingSamplingRatio {
						t.Errorf("Expected tracingSamplingRatio %f, got %f", tc.expectedTracingSamplingRatio, config.TracingSamplingRatio())
					}
				}
			}
		})
	}
}
