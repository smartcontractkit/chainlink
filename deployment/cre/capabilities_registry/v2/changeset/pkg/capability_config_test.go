package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

func TestCapabilityConfig_MarshalUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("matching keys to proto", func(t *testing.T) {
		t.Parallel()

		rawMap := map[string]any{
			"restrictedConfig": map[string]any{
				"fields": map[string]any{
					"spendRatios": map[string]any{
						"mapValue": map[string]any{
							"fields": map[string]any{
								"RESOURCE_TYPE_COMPUTE": map[string]any{
									"stringValue": "1.0",
								},
							},
						},
					},
				},
			},
			"methodConfigs": map[string]any{
				"BalanceAt": map[string]any{
					"remoteExecutableConfig": map[string]any{
						"requestTimeout":            "30s",
						"serverMaxParallelRequests": float64(10),
					},
					"aggregatorConfig": map[string]any{
						"aggregatorType": "SignedReport",
					},
				},
			},
		}

		cfg := pkg.CapabilityConfig(rawMap)
		protoBts, err := cfg.MarshalProto()

		require.NoError(t, err)

		result := pkg.CapabilityConfig(map[string]any{})
		require.NoError(t, result.UnmarshalProto(protoBts))

		assert.Equal(t, rawMap, map[string]any(result))
	})
}

func TestCapabilityConfig_UnmarshalWithValidation(t *testing.T) {
	t.Parallel()

	t.Run("unknown fields are detected and reported", func(t *testing.T) {
		t.Parallel()

		// Create a config with an intentional typo - "methodConfigsTypo" instead of "methodConfigs"
		rawMap := map[string]any{
			"restrictedConfig": map[string]any{
				"fields": map[string]any{
					"spendRatios": map[string]any{
						"mapValue": map[string]any{
							"fields": map[string]any{
								"RESOURCE_TYPE_COMPUTE": map[string]any{
									"stringValue": "1.0",
								},
							},
						},
					},
				},
			},
			"methodConfigsTypo": map[string]any{ // intentional typo: should be "methodConfigs"
				"BalanceAt": map[string]any{
					"remoteExecutableConfig": map[string]any{
						"requestTimeout":            "30s",
						"serverMaxParallelRequests": float64(10),
					},
				},
			},
		}

		cfg := pkg.CapabilityConfig(rawMap)
		_, err := cfg.MarshalProto()

		// Should return an error with the specific field name
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config validation failed")
		assert.Contains(t, err.Error(), "methodConfigsTypo")
	})

	t.Run("deltaStageNanos vs deltaStage field name typo", func(t *testing.T) {
		t.Parallel()

		// Using incorrect field name "deltaStageNanos" instead of correct "deltaStage"
		rawMap := map[string]any{
			"methodConfigs": map[string]any{
				"BalanceAt": map[string]any{
					"remoteExecutableConfig": map[string]any{
						"requestTimeout":            "30s",
						"serverMaxParallelRequests": float64(10),
						"deltaStageNanos":           "1000000", // Typo: should be "deltaStage"
					},
				},
			},
		}

		cfg := pkg.CapabilityConfig(rawMap)
		_, err := cfg.MarshalProto()

		// Should return an error with the specific field name
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config validation failed")
		assert.Contains(t, err.Error(), "deltaStageNanos")
	})

	t.Run("valid config passes validation", func(t *testing.T) {
		t.Parallel()

		// Valid config matching proto schema
		rawMap := map[string]any{
			"restrictedConfig": map[string]any{
				"fields": map[string]any{
					"spendRatios": map[string]any{
						"mapValue": map[string]any{
							"fields": map[string]any{
								"RESOURCE_TYPE_COMPUTE": map[string]any{
									"stringValue": "1.0",
								},
							},
						},
					},
				},
			},
			"methodConfigs": map[string]any{
				"BalanceAt": map[string]any{
					"remoteExecutableConfig": map[string]any{
						"requestTimeout":            "30s",
						"serverMaxParallelRequests": float64(10),
					},
					"aggregatorConfig": map[string]any{
						"aggregatorType": "SignedReport",
					},
				},
			},
		}

		cfg := pkg.CapabilityConfig(rawMap)
		protoBts, err := cfg.MarshalProto()

		// Should succeed with no errors
		require.NoError(t, err)
		require.NotNil(t, protoBts)
	})
}
