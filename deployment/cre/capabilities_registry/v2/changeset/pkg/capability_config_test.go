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
