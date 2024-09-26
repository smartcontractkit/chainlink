package transmission

import (
	"fmt"
	"testing"
	"time"

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
				"one":   300 * time.Millisecond,
				"two":   0 * time.Millisecond,
				"three": 100 * time.Millisecond,
				"four":  200 * time.Millisecond,
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
				"one":   300 * time.Millisecond,
				"two":   100 * time.Millisecond,
				"three": 200 * time.Millisecond,
				"four":  0 * time.Millisecond,
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
