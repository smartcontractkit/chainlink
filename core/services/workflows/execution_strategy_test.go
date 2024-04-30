package workflows

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func assertBetween(t *testing.T, got time.Duration, low time.Duration, high time.Duration) {
	assert.GreaterOrEqual(t, got, low)
	assert.LessOrEqual(t, got, high)
}

func TestScheduledExecutionStrategy_LocalDON(t *testing.T) {
	var gotTime time.Time
	var called bool

	// Our capability has DONInfo == nil, so we'll treat it as a local
	// capability and use the local DON Info to determine the transmission
	// schedule.
	mt := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_polygon-testnet-mumbai",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mumbai testnet",
			"v1.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			gotTime = time.Now()
			called = true
			return capabilities.CapabilityResponse{}, nil
		},
	)

	l := logger.TestLogger(t)

	// The combination of this key and the metadata above
	// will yield the permutation [3, 2, 0, 1]
	key, err := hex.DecodeString("fb13ca015a9ec60089c7141e9522de79")
	require.NoError(t, err)

	testCases := []struct {
		name     string
		position int
		schedule string
		low      time.Duration
		high     time.Duration
	}{
		{
			name:     "position 0; oneAtATime",
			position: 0,
			schedule: "oneAtATime",
			low:      300 * time.Millisecond,
			high:     400 * time.Millisecond,
		},
		{
			name:     "position 1; oneAtATime",
			position: 1,
			schedule: "oneAtATime",
			low:      200 * time.Millisecond,
			high:     300 * time.Millisecond,
		},
		{
			name:     "position 2; oneAtATime",
			position: 2,
			schedule: "oneAtATime",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
		},
		{
			name:     "position 3; oneAtATime",
			position: 3,
			schedule: "oneAtATime",
			low:      100 * time.Millisecond,
			high:     200 * time.Millisecond,
		},
		{
			name:     "position 0; allAtOnce",
			position: 0,
			schedule: "allAtOnce",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
		},
		{
			name:     "position 1; allAtOnce",
			position: 1,
			schedule: "allAtOnce",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
		},
		{
			name:     "position 2; allAtOnce",
			position: 2,
			schedule: "allAtOnce",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
		},
		{
			name:     "position 3; allAtOnce",
			position: 3,
			schedule: "allAtOnce",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()

			m, err := values.NewMap(map[string]any{
				"schedule":   tc.schedule,
				"deltaStage": "100ms",
			})
			require.NoError(t, err)

			req := capabilities.CapabilityRequest{
				Config: m,
				Metadata: capabilities.RequestMetadata{
					WorkflowID:          "mock-workflow-id",
					WorkflowExecutionID: "mock-execution-id",
				},
			}

			ids := []p2ptypes.PeerID{
				randKey(),
				randKey(),
				randKey(),
				randKey(),
			}
			don := &capabilities.DON{
				Members: ids,
				Config: capabilities.DONConfig{
					SharedSecret: [16]byte(key),
				},
			}
			peerID := ids[tc.position]
			de := scheduledExecution{
				DON:      don,
				PeerID:   &peerID,
				Position: tc.position,
			}
			_, err = de.Apply(tests.Context(t), l, mt, req)
			require.NoError(t, err)
			require.True(t, called)

			assertBetween(t, gotTime.Sub(startTime), tc.low, tc.high)
		})
	}
}

func randKey() [32]byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return [32]byte(key)
}
