package sequences

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

func TestDeepMergeMaps(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		result := deepMergeMaps(nil, nil)
		assert.Nil(t, result)
	})

	t.Run("nil base", func(t *testing.T) {
		override := map[string]any{"key": "value"}
		result := deepMergeMaps(nil, override)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("nil override", func(t *testing.T) {
		base := map[string]any{"key": "value"}
		result := deepMergeMaps(base, nil)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("scalar override", func(t *testing.T) {
		base := map[string]any{"a": 1, "b": 2}
		override := map[string]any{"b": 3}
		result := deepMergeMaps(base, override)
		assert.Equal(t, map[string]any{"a": 1, "b": 3}, result)
	})

	t.Run("add new key", func(t *testing.T) {
		base := map[string]any{"a": 1}
		override := map[string]any{"b": 2}
		result := deepMergeMaps(base, override)
		assert.Equal(t, map[string]any{"a": 1, "b": 2}, result)
	})

	t.Run("nested map merge", func(t *testing.T) {
		base := map[string]any{
			"outer": map[string]any{
				"keep": "yes",
				"inner": map[string]any{
					"x": 1,
					"y": 2,
				},
			},
		}
		override := map[string]any{
			"outer": map[string]any{
				"inner": map[string]any{
					"y": 99,
				},
			},
		}
		result := deepMergeMaps(base, override)
		expected := map[string]any{
			"outer": map[string]any{
				"keep": "yes",
				"inner": map[string]any{
					"x": 1,
					"y": 99,
				},
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("does not mutate inputs", func(t *testing.T) {
		base := map[string]any{
			"nested": map[string]any{"a": 1},
		}
		override := map[string]any{
			"nested": map[string]any{"b": 2},
		}
		_ = deepMergeMaps(base, override)

		assert.Equal(t, map[string]any{"nested": map[string]any{"a": 1}}, base)
		assert.Equal(t, map[string]any{"nested": map[string]any{"b": 2}}, override)
	})

	t.Run("realistic minResponsesToAggregate override", func(t *testing.T) {
		base := map[string]any{
			"methodConfigs": map[string]any{
				"LogTrigger": map[string]any{
					"remoteTriggerConfig": map[string]any{
						"registrationRefresh":     "20s",
						"registrationExpiry":      "60s",
						"minResponsesToAggregate": 4,
						"messageExpiry":           "120s",
					},
				},
				"BalanceAt": map[string]any{
					"remoteExecutableConfig": map[string]any{
						"requestTimeout": "30s",
					},
				},
			},
		}
		override := map[string]any{
			"methodConfigs": map[string]any{
				"LogTrigger": map[string]any{
					"remoteTriggerConfig": map[string]any{
						"minResponsesToAggregate": 2,
					},
				},
			},
		}
		result := deepMergeMaps(base, override)

		logTrigger := result["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)
		assert.Equal(t, 2, logTrigger["minResponsesToAggregate"])
		assert.Equal(t, "20s", logTrigger["registrationRefresh"])
		assert.Equal(t, "60s", logTrigger["registrationExpiry"])
		assert.Equal(t, "120s", logTrigger["messageExpiry"])

		balanceAt := result["methodConfigs"].(map[string]any)["BalanceAt"].(map[string]any)["remoteExecutableConfig"].(map[string]any)
		assert.Equal(t, "30s", balanceAt["requestTimeout"])
	})
}

func TestResolveCapabilityConfigsForDON(t *testing.T) {
	baseConfigs := []contracts.CapabilityConfig{
		{
			Capability: contracts.Capability{CapabilityID: "cap-a@1.0.0"},
			Config: map[string]any{
				"methodConfigs": map[string]any{
					"LogTrigger": map[string]any{
						"remoteTriggerConfig": map[string]any{
							"minResponsesToAggregate": 4,
							"registrationRefresh":     "20s",
						},
					},
				},
			},
		},
		{
			Capability: contracts.Capability{CapabilityID: "cap-b@1.0.0"},
			Config: map[string]any{
				"methodConfigs": map[string]any{
					"LogTrigger": map[string]any{
						"remoteTriggerConfig": map[string]any{
							"minResponsesToAggregate": 4,
							"registrationRefresh":     "20s",
						},
					},
				},
			},
		},
	}

	t.Run("no overrides returns base", func(t *testing.T) {
		result, err := resolveCapabilityConfigsForDON(baseConfigs, nil)
		require.NoError(t, err)
		assert.Equal(t, baseConfigs, result)
	})

	t.Run("empty overrides returns base", func(t *testing.T) {
		result, err := resolveCapabilityConfigsForDON(baseConfigs, []CapabilityConfigOverride{})
		require.NoError(t, err)
		assert.Equal(t, baseConfigs, result)
	})

	t.Run("override specific capability", func(t *testing.T) {
		overrides := []CapabilityConfigOverride{
			{
				CapabilityID: "cap-a@1.0.0",
				Config: map[string]any{
					"methodConfigs": map[string]any{
						"LogTrigger": map[string]any{
							"remoteTriggerConfig": map[string]any{
								"minResponsesToAggregate": 2,
							},
						},
					},
				},
			},
		}

		result, err := resolveCapabilityConfigsForDON(baseConfigs, overrides)
		require.NoError(t, err)
		require.Len(t, result, 2)

		capAConfig := result[0].Config["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)
		assert.Equal(t, 2, capAConfig["minResponsesToAggregate"])
		assert.Equal(t, "20s", capAConfig["registrationRefresh"])

		capBConfig := result[1].Config["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)
		assert.Equal(t, 4, capBConfig["minResponsesToAggregate"])
	})

	t.Run("override all capabilities with empty ID", func(t *testing.T) {
		overrides := []CapabilityConfigOverride{
			{
				CapabilityID: "",
				Config: map[string]any{
					"methodConfigs": map[string]any{
						"LogTrigger": map[string]any{
							"remoteTriggerConfig": map[string]any{
								"minResponsesToAggregate": 2,
							},
						},
					},
				},
			},
		}

		result, err := resolveCapabilityConfigsForDON(baseConfigs, overrides)
		require.NoError(t, err)
		require.Len(t, result, 2)

		for _, cfg := range result {
			logTrigger := cfg.Config["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)
			assert.Equal(t, 2, logTrigger["minResponsesToAggregate"])
			assert.Equal(t, "20s", logTrigger["registrationRefresh"])
		}
	})

	t.Run("error on non-existent capability ID", func(t *testing.T) {
		overrides := []CapabilityConfigOverride{
			{
				CapabilityID: "non-existent@1.0.0",
				Config:       map[string]any{"foo": "bar"},
			},
		}

		_, err := resolveCapabilityConfigsForDON(baseConfigs, overrides)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-existent@1.0.0")
	})

	t.Run("nil config override is skipped", func(t *testing.T) {
		overrides := []CapabilityConfigOverride{
			{CapabilityID: "", Config: nil},
		}

		result, err := resolveCapabilityConfigsForDON(baseConfigs, overrides)
		require.NoError(t, err)
		assert.Equal(t, 4, result[0].Config["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)["minResponsesToAggregate"])
	})

	t.Run("does not mutate base configs", func(t *testing.T) {
		overrides := []CapabilityConfigOverride{
			{
				CapabilityID: "",
				Config: map[string]any{
					"methodConfigs": map[string]any{
						"LogTrigger": map[string]any{
							"remoteTriggerConfig": map[string]any{
								"minResponsesToAggregate": 99,
							},
						},
					},
				},
			},
		}

		_, err := resolveCapabilityConfigsForDON(baseConfigs, overrides)
		require.NoError(t, err)

		originalVal := baseConfigs[0].Config["methodConfigs"].(map[string]any)["LogTrigger"].(map[string]any)["remoteTriggerConfig"].(map[string]any)["minResponsesToAggregate"]
		assert.Equal(t, 4, originalVal)
	})
}
