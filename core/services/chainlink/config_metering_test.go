package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestMeteringConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		mc := meteringConfig{s: toml.Metering{}}
		assert.False(t, mc.MeterRecordsEnabled())
		assert.False(t, mc.MeterSnapshotsEnabled())
		assert.Empty(t, mc.Product())
		assert.Empty(t, mc.Tenant())
		assert.Empty(t, mc.Environment())
		assert.Empty(t, mc.Zone())
		assert.Empty(t, mc.NodeID())
	})

	t.Run("explicit values", func(t *testing.T) {
		mc := meteringConfig{s: toml.Metering{
			MeterRecordsEnabled:   ptr(true),
			MeterSnapshotsEnabled: ptr(true),
			Product:               ptr("cre"),
			Tenant:                ptr("mainline"),
			Environment:           ptr("production"),
			Zone:                  ptr("wf-zone-a"),
			NodeID:                ptr("csa-pubkey-1"),
		}}
		assert.True(t, mc.MeterRecordsEnabled())
		assert.True(t, mc.MeterSnapshotsEnabled())
		assert.Equal(t, "cre", mc.Product())
		assert.Equal(t, "mainline", mc.Tenant())
		assert.Equal(t, "production", mc.Environment())
		assert.Equal(t, "wf-zone-a", mc.Zone())
		assert.Equal(t, "csa-pubkey-1", mc.NodeID())
	})
}
