package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

func TestBuildCREMetadataV2_IncludesZone(t *testing.T) {
	t.Parallel()

	labels := map[string]string{
		platform.KeyDonID:  "5",
		platform.KeyP2PID:  "12D3KooWRp2Gdj1JPot5425NshLBeyJJPno95GAIrzTNUm6cgx3B",
		platform.KeyZone:   "zone-a",
		platform.KeyDonF:   "1",
		platform.KeyDonN:   "4",
	}

	creInfo := buildCREMetadataV2(labels)
	require.NotNil(t, creInfo)
	assert.Equal(t, int32(5), creInfo.DonID)
	assert.Equal(t, "12D3KooWRp2Gdj1JPot5425NshLBeyJJPno95GAIrzTNUm6cgx3B", creInfo.P2PID)
	assert.Equal(t, "zone-a", creInfo.Zone)
}

func TestBuildCREMetadataV2_EmptyZoneWhenAbsent(t *testing.T) {
	t.Parallel()

	creInfo := buildCREMetadataV2(map[string]string{
		platform.KeyDonID: "2",
	})
	require.NotNil(t, creInfo)
	assert.Empty(t, creInfo.Zone)
}
