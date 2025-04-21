package transmission

import (
	"fmt"
	"testing"
	"time"

	types2 "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_GetPeerIDToTransmissionDelay(t *testing.T) {
	peer1 := [32]byte([]byte(fmt.Sprintf("%-32s", "one")))
	peer2 := [32]byte([]byte(fmt.Sprintf("%-32s", "two")))
	peer3 := [32]byte([]byte(fmt.Sprintf("%-32s", "three")))
	peer4 := [32]byte([]byte(fmt.Sprintf("%-32s", "four")))

	ids := []p2ptypes.PeerID{
		peer1, peer2, peer3, peer4,
	}

	testCases := []struct {
		name                string
		peerName            string
		schedule            string
		deltaStage          string
		workflowExecutionID string
		expectedDelays      map[string]time.Duration
	}{
		{
			"TestOneAtATime",
			"one",
			"oneAtATime",
			"100ms",
			"15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
			map[string]time.Duration{
				"one":   100 * time.Millisecond,
				"two":   200 * time.Millisecond,
				"three": 300 * time.Millisecond,
				"four":  0 * time.Millisecond,
			},
		},

		{
			"TestAllAtOnce",
			"one",
			"allAtOnce",
			"100ms",
			"15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
			map[string]time.Duration{
				"one":   0 * time.Millisecond,
				"two":   0 * time.Millisecond,
				"three": 0 * time.Millisecond,
				"four":  0 * time.Millisecond,
			},
		},

		{
			"TestOneAtATimeWithDifferentExecutionID",
			"one",
			"oneAtATime",
			"100ms",
			"16c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			map[string]time.Duration{
				"one":   200 * time.Millisecond,
				"two":   100 * time.Millisecond,
				"three": 0 * time.Millisecond,
				"four":  300 * time.Millisecond,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transmissionCfg, err := values.NewMap(map[string]any{
				"schedule":   tc.schedule,
				"deltaStage": tc.deltaStage,
			})
			require.NoError(t, err)

			capabilityRequest := capabilities.CapabilityRequest{
				Config: transmissionCfg,
				Metadata: capabilities.RequestMetadata{
					WorkflowID:          "17c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
					WorkflowExecutionID: tc.workflowExecutionID,
				},
			}

			peerIdToDelay, err := GetPeerIDToTransmissionDelay(ids, capabilityRequest)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedDelays["one"], peerIdToDelay[peer1])
			assert.Equal(t, tc.expectedDelays["two"], peerIdToDelay[peer2])
			assert.Equal(t, tc.expectedDelays["three"], peerIdToDelay[peer3])
			assert.Equal(t, tc.expectedDelays["four"], peerIdToDelay[peer4])
		})
	}
}

func TestGetPeerIDToTransmissionDelaysForConfig_ConsistentSchedule(t *testing.T) {
	// Create test peer IDs
	peerIDs := make([]types2.PeerID, 5)
	reversedPeerIDs := make([]types2.PeerID, 5)
	for i := range peerIDs {
		peerIDs[i] = [32]byte{byte(i)}
		reversedPeerIDs[len(peerIDs)-1-i] = peerIDs[i]
	}

	// Test configuration
	transmissionID := "test-transmission"
	config := TransmissionConfig{
		Schedule:   Schedule_OneAtATime,
		DeltaStage: 1 * time.Second,
	}

	// Get first schedule
	schedule1, err := GetPeerIDToTransmissionDelaysForConfig(peerIDs, transmissionID, config)
	require.NoError(t, err)

	// Get second schedule with same inputs
	schedule2, err := GetPeerIDToTransmissionDelaysForConfig(peerIDs, transmissionID, config)
	require.NoError(t, err)

	// Verify each peer has the same delay in both schedules
	for peerID, delay1 := range schedule1 {
		delay2, exists := schedule2[peerID]
		require.True(t, exists, "peer %v should exist in both schedules", peerID)
		require.Equal(t, delay1, delay2, "peer %v should have same delay in both schedules", peerID)
	}

	// Verify the order of delays is consistent
	delays1 := make([]time.Duration, 0, len(schedule1))
	delays2 := make([]time.Duration, 0, len(schedule2))
	for _, peerID := range peerIDs {
		delays1 = append(delays1, schedule1[peerID])
		delays2 = append(delays2, schedule2[peerID])
	}
	require.Equal(t, delays1, delays2, "delays should be in same order in both schedules")

	// Verify with different transmission ID
	differentSchedule, err := GetPeerIDToTransmissionDelaysForConfig(peerIDs, "different-transmission", config)
	require.NoError(t, err)
	require.NotEqual(t, schedule1, differentSchedule, "different transmission ID should produce different schedule")

	// Verify with different peer order
	reversedSchedule, err := GetPeerIDToTransmissionDelaysForConfig(reversedPeerIDs, transmissionID, config)
	require.NoError(t, err)
	for i := range peerIDs {
		require.Equal(t, schedule1[[32]byte{byte(i)}], reversedSchedule[[32]byte{byte(i)}], "peer order should not affect schedule")
	}
}
