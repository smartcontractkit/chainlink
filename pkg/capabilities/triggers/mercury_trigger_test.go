package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func TestMercuryTrigger(t *testing.T) {
	ts := NewMercuryTriggerService()
	ctx := tests.Context(t)
	require.NotNil(t, ts)

	m := map[string]interface{}{
		"feedIds":   []int64{1},
		"triggerId": "test-id-1",
	}

	wrapped, err := values.NewMap(m)
	require.NoError(t, err)

	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Inputs: wrapped,
	}
	callback := make(chan capabilities.CapabilityResponse, 10)
	require.NoError(t, ts.RegisterTrigger(ctx, callback, cr))

	// Send events to trigger and check for them in the callback
	fr := []mercury.FeedReport{
		{
			FeedID:               1,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       2,
			ObservationTimestamp: 3,
		},
	}
	err = ts.ProcessReport(fr)
	assert.NoError(t, err)
	assert.Len(t, callback, 1)
	msg := <-callback
	unwrapped, _ := mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	assert.Equal(t, GenerateTriggerEventID(fr), unwrapped.ID)
	assert.Len(t, unwrapped.Payload, 1)
	assert.Equal(t, fr[0], unwrapped.Payload[0])

	// Unregister the trigger and check that events no longer go on the callback
	require.NoError(t, ts.UnregisterTrigger(ctx, cr))
	err = ts.ProcessReport(fr)
	assert.NoError(t, err)
	assert.Len(t, callback, 0)
}

func TestMultipleMercuryTriggers(t *testing.T) {
	ts := NewMercuryTriggerService()
	ctx := tests.Context(t)
	require.NotNil(t, ts)

	m1 := map[string]interface{}{
		"feedIds":   []int64{1, 3, 4},
		"triggerId": "test-id-1",
	}

	m2 := map[string]interface{}{
		"feedIds":   []int64{2, 3, 5},
		"triggerId": "test-id-2",
	}

	wrapped1, err := values.NewMap(m1)
	require.NoError(t, err)

	wrapped2, err := values.NewMap(m2)
	require.NoError(t, err)

	cr1 := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Inputs: wrapped1,
	}
	cr2 := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Inputs: wrapped2,
	}

	callback1 := make(chan capabilities.CapabilityResponse, 10)
	callback2 := make(chan capabilities.CapabilityResponse, 10)

	require.NoError(t, ts.RegisterTrigger(ctx, callback1, cr1))
	require.NoError(t, ts.RegisterTrigger(ctx, callback2, cr2))

	// Send events to trigger and check for them in the callback
	fr1 := []mercury.FeedReport{
		{
			FeedID:               1,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       20,
			ObservationTimestamp: 5,
		},
		{
			FeedID:               3,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       25,
			ObservationTimestamp: 8,
		},
		{
			FeedID:               2,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       30,
			ObservationTimestamp: 10,
		},
		{
			FeedID:               4,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       40,
			ObservationTimestamp: 15,
		},
	}

	err = ts.ProcessReport(fr1)
	assert.NoError(t, err)
	assert.Len(t, callback1, 1)
	assert.Len(t, callback2, 1)

	msg := <-callback1
	unwrapped, _ := mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	payload := make([]mercury.FeedReport, 0)
	payload = append(payload, fr1[0], fr1[1], fr1[3])
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.Payload, 3)
	assert.Equal(t, fr1[0], unwrapped.Payload[0])
	assert.Equal(t, fr1[1], unwrapped.Payload[1])
	assert.Equal(t, fr1[3], unwrapped.Payload[2])

	msg = <-callback2
	unwrapped, _ = mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	payload = make([]mercury.FeedReport, 0)
	payload = append(payload, fr1[1], fr1[2]) // Because GenerateTriggerEventID sorts the reports by feedID, this works
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.Payload, 2)
	assert.Equal(t, fr1[2], unwrapped.Payload[0])
	assert.Equal(t, fr1[1], unwrapped.Payload[1])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr1))
	fr2 := []mercury.FeedReport{
		{
			FeedID:               3,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       50,
			ObservationTimestamp: 20,
		},
	}
	err = ts.ProcessReport(fr2)
	assert.NoError(t, err)
	assert.Len(t, callback1, 0)
	assert.Len(t, callback2, 1)

	msg = <-callback2
	unwrapped, _ = mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	payload = make([]mercury.FeedReport, 0)
	payload = append(payload, fr2[0])
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.Payload, 1)
	assert.Equal(t, fr2[0], unwrapped.Payload[0])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr2))
	err = ts.ProcessReport(fr1)
	assert.NoError(t, err)
	assert.Len(t, callback1, 0)
	assert.Len(t, callback2, 0)
}
