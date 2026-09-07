package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestMeteringConfig(t *testing.T) {
	t.Parallel()
	t.Run("defaults", func(t *testing.T) {
		t.Parallel()
		mc := meteringConfig{s: toml.Metering{}}
		assert.False(t, mc.MeterRecordsEnabled())
		assert.False(t, mc.MeterSnapshotsEnabled())
		// A zero-value toml.Metering (not run through setDefaults) has a nil
		// Product pointer, so Product() returns "unset". The parsed config
		// applies a "cre" default via docs.CoreDefaults (covered by the
		// LogConfiguration effective-TOML test), so metering is never enabled
		// with an empty product dimension.
		assert.Equal(t, "unset", mc.Product())
		assert.Empty(t, mc.Tenant())
		assert.Empty(t, mc.NumericTenantID())
		assert.Empty(t, mc.Environment())
		assert.Empty(t, mc.Zone())
		assert.Empty(t, mc.NodeID())
	})

	t.Run("explicit values", func(t *testing.T) {
		t.Parallel()
		mc := meteringConfig{s: toml.Metering{
			MeterRecordsEnabled:   new(true),
			MeterSnapshotsEnabled: new(true),
			Product:               new("cre"),
			Tenant:                new("mainline"),
			NumericTenantID:       new("42"),
			Environment:           new("production"),
			Zone:                  new("wf-zone-a"),
			NodeID:                new("clp-cre-wf-zone-a-1"),
		}}
		assert.True(t, mc.MeterRecordsEnabled())
		assert.True(t, mc.MeterSnapshotsEnabled())
		assert.Equal(t, "cre", mc.Product())
		assert.Equal(t, "mainline", mc.Tenant())
		assert.Equal(t, "42", mc.NumericTenantID())
		assert.Equal(t, "production", mc.Environment())
		assert.Equal(t, "wf-zone-a", mc.Zone())
		assert.Equal(t, "clp-cre-wf-zone-a-1", mc.NodeID())
	})
}
