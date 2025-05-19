package transmission

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
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

	log := logger.TestLogger(t)

	// Our capability has DONInfo == nil, so we'll treat it as a local
	// capability and use the local DON Info to determine the transmission
	// schedule.
	mt := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_polygon-testnet-mumbai@1.0.0",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mumbai testnet",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			gotTime = time.Now()
			called = true
			return capabilities.CapabilityResponse{}, nil
		},
	)

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
			low:      200 * time.Millisecond,
			high:     300 * time.Millisecond,
		},
		{
			name:     "position 1; oneAtATime",
			position: 1,
			schedule: "oneAtATime",
			low:      300 * time.Millisecond,
			high:     400 * time.Millisecond,
		},
		{
			name:     "position 2; oneAtATime",
			position: 2,
			schedule: "oneAtATime",
			low:      100 * time.Millisecond,
			high:     200 * time.Millisecond,
		},
		{
			name:     "position 3; oneAtATime",
			position: 3,
			schedule: "oneAtATime",
			low:      0 * time.Millisecond,
			high:     100 * time.Millisecond,
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
					WorkflowID:          "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0",
					WorkflowExecutionID: "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
				},
			}

			ids := []p2ptypes.PeerID{
				randKey(),
				randKey(),
				randKey(),
				randKey(),
			}
			localDON := capabilities.Node{
				WorkflowDON: capabilities.DON{
					ID:      1,
					Members: ids,
				},
				PeerID: &ids[tc.position],
			}
			localTargetCapability := NewLocalTargetCapability(log, "capabilityID", localDON, mt)

			_, err = localTargetCapability.Execute(t.Context(), req)

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

type mockCapability struct {
	capabilities.CapabilityInfo
	capabilities.Executable
	response  chan capabilities.CapabilityResponse
	transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)
}

func newMockCapability(info capabilities.CapabilityInfo, transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)) *mockCapability {
	return &mockCapability{
		transform:      transform,
		CapabilityInfo: info,
		response:       make(chan capabilities.CapabilityResponse, 10),
	}
}

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	cr, err := m.transform(req)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	m.response <- cr
	return cr, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
