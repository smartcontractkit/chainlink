package cresettings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig/globalconfigtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const tq = "'''" // TOML multi-line literal string delimiter

// TestValidatedCRESettingsSpec_SchemaVersion: v1 specs are unchanged; the top-level config_type
// and offchain_config fields require v2, so nodes predating them (which accept only v1) reject
// such a spec instead of storing it as an empty settings spec.
func TestValidatedCRESettingsSpec_SchemaVersion(t *testing.T) {
	t.Parallel()

	capReg := `config_type = "capabilities_registry"` + "\noffchain_config = " + tq + `{"version":1}` + tq
	settings := "settings = " + tq + `Foo = "bar"` + tq
	for _, tc := range []struct {
		name    string
		version int
		body    string
		wantErr string
	}{
		{"v1 settings", 1, settings, ""},
		{"v2 settings", 2, settings, ""},
		{"v1 embedded shard_assignment", 1, "settings = " + tq + "config_type = \"shard_assignment\"\nstatic_default_assignment = [0]" + tq, ""},
		{"v2 capabilities_registry", 2, capReg, ""},
		{"v1 capabilities_registry", 1, capReg, "require schemaVersion = 2"},
		{"v1 top-level shard_assignment", 1, `config_type = "shard_assignment"` + "\nsettings = " + tq + "static_default_assignment = [0]" + tq, "require schemaVersion = 2"},
		{"v1 offchain_config only", 1, "offchain_config = " + tq + `{"version":1}` + tq, "require schemaVersion = 2"},
		{"v3", 3, settings, "unsupported schemaVersion 3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			spec := fmt.Sprintf("type = \"cresettings\"\nschemaVersion = %d\n%s", tc.version, tc.body)
			_, err := ValidatedCRESettingsSpec(spec)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			// job.ValidateSpec is what the feeds manager, REST and GraphQL call first.
			jt, err := job.ValidateSpec(spec)
			require.NoError(t, err)
			assert.Equal(t, job.CRESettings, jt)
		})
	}

	_, err := job.ValidateSpec("type = \"cresettings\"\nschemaVersion = 3\n" + settings)
	require.ErrorIs(t, err, job.ErrInvalidSchemaVersion)
}

// TestValidatedCRESettingsSpec_DesignDocShape validates a spec in the design doc's format
// (name-keyed DONs, domain/env, schemaVersion 2, fixed externalJobID, name), using real
// CapabilityConfig fields.
func TestValidatedCRESettingsSpec_DesignDocShape(t *testing.T) {
	t.Parallel()

	jb, err := ValidatedCRESettingsSpec(`type = "cresettings"
schemaVersion = 2
config_type = "capabilities_registry"
externalJobID = "a3f8c1d2-4e5b-6a7c-8d9e-0f1a2b3c4d5e"
name = "cre-prod-cap-config"
offchain_config = ` + tq + `
{
  "domain": "cre",
  "env": "prod",
  "version": 201,
  "dons": {
    "workflow_1_zone-a": {
      "capabilities": {
        "cron@1.0.0": {
          "specConfig": { "fields": { "interval": { "stringValue": "30" } } }
        }
      }
    }
  }
}
` + tq)
	require.NoError(t, err)
	assert.Equal(t, uint32(2), jb.SchemaVersion)
	assert.Equal(t, "a3f8c1d2-4e5b-6a7c-8d9e-0f1a2b3c4d5e", jb.ExternalJobID.String())
	reg, err := globalconfig.Parse(jb.CRESettingsSpec.OffchainConfig)
	require.NoError(t, err)
	assert.Equal(t, "cre", reg.GetDomain())
	assert.Equal(t, "prod", reg.GetEnv())
	assert.Equal(t, uint64(201), reg.GetVersion())
	assert.Contains(t, reg.GetDons()["workflow_1_zone-a"].GetCapabilities(), "cron@1.0.0")
}

// TestValidatedCRESettingsSpec_ValidationErrorMetric swaps the package metrics sink, so it does
// not run in parallel.
func TestValidatedCRESettingsSpec_ValidationErrorMetric(t *testing.T) { //nolint:paralleltest // swaps package state
	m, reader := globalconfigtest.NewMetrics(t)
	prev := capRegMetrics
	capRegMetrics = func() *globalconfig.Metrics { return m }
	t.Cleanup(func() { capRegMetrics = prev })

	spec := func(version int, offchain string) string {
		return fmt.Sprintf("type = \"cresettings\"\nschemaVersion = %d\nconfig_type = \"capabilities_registry\"\noffchain_config = %s%s%s", version, tq, offchain, tq)
	}
	_, err := ValidatedCRESettingsSpec(spec(2, `{"domain":"cre","env":"prod","version":1,"dons":{"d":{"capabilityConfigs":{}}}}`))
	require.ErrorContains(t, err, "unknown field")
	_, err = ValidatedCRESettingsSpec(spec(1, `{"domain":"cre","env":"prod","version":2}`))
	require.ErrorContains(t, err, "schemaVersion")
	_, err = ValidatedCRESettingsSpec(spec(2, `{"domain":"cre","env":"staging","version":0}`))
	require.ErrorContains(t, err, "version must be >= 1")
	_, err = ValidatedCRESettingsSpec(spec(2, `{"domain":"cre","env":"prod","version":3}`))
	require.NoError(t, err) // valid: not counted
	_, err = ValidatedCRESettingsSpec("type = \"cresettings\"\nschemaVersion = 1\nsettings = " + tq + "not = [toml" + tq)
	require.Error(t, err) // a settings spec, not a capabilities_registry payload: not counted

	got := globalconfigtest.Collect(t, reader)
	assert.Equal(t, int64(2), got[globalconfig.MetricValidationErrors][globalconfigtest.Series("cre", "prod")])
	assert.Equal(t, int64(1), got[globalconfig.MetricValidationErrors][globalconfigtest.Series("cre", "staging")])
	assert.Len(t, got[globalconfig.MetricValidationErrors], 2)
}
