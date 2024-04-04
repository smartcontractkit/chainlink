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

var (
	feedOne   = "0x1111111111111111111100000000000000000000000000000000000000000000"
	feedTwo   = "0x2222222222222222222200000000000000000000000000000000000000000000"
	feedThree = "0x3333333333333333333300000000000000000000000000000000000000000000"
	feedFour  = "0x4444444444444444444400000000000000000000000000000000000000000000"
	feedFive  = "0x5555555555555555555500000000000000000000000000000000000000000000"
)

func TestMercuryTrigger(t *testing.T) {
	ts := NewMercuryTriggerService()
	ctx := tests.Context(t)
	require.NotNil(t, ts)

	cm := map[string]interface{}{
		"feedIds": []string{feedOne},
	}

	configWrapped, err := values.NewMap(cm)
	require.NoError(t, err)

	im := map[string]interface{}{
		"triggerId": "test-id-1",
	}

	inputsWrapped, err := values.NewMap(im)
	require.NoError(t, err)

	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Config: configWrapped,
		Inputs: inputsWrapped,
	}
	callback := make(chan capabilities.CapabilityResponse, 10)
	require.NoError(t, ts.RegisterTrigger(ctx, callback, cr))

	// Send events to trigger and check for them in the callback
	fr := []FeedReport{
		{
			FeedID:               mercury.FeedID(feedOne).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       2,
			ObservationTimestamp: 3,
		},
	}

	mfr := []mercury.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       2,
			ObservationTimestamp: 3,
		},
	}
	err = ts.ProcessReport(fr)
	assert.NoError(t, err)
	assert.Len(t, callback, 1)
	msg := <-callback
	unwrapped, err := mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	require.NoError(t, err)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	assert.Equal(t, GenerateTriggerEventID(mfr), unwrapped.ID)
	assert.Len(t, unwrapped.BatchedPayload, 1)
	assert.Equal(t, mfr[0], unwrapped.BatchedPayload[feedOne])

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

	im1 := map[string]interface{}{
		"triggerId": "test-id-1",
	}

	iwrapped1, err := values.NewMap(im1)
	require.NoError(t, err)

	cm1 := map[string]interface{}{
		"feedIds": []string{
			feedOne,
			feedThree,
			feedFour,
		},
	}

	cwrapped1, err := values.NewMap(cm1)
	require.NoError(t, err)

	cr1 := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Inputs: iwrapped1,
		Config: cwrapped1,
	}

	im2 := map[string]interface{}{
		"triggerId": "test-id-2",
	}

	iwrapped2, err := values.NewMap(im2)
	require.NoError(t, err)

	cm2 := map[string]interface{}{
		"feedIds": []string{
			feedTwo,
			feedThree,
			feedFive,
		},
	}

	cwrapped2, err := values.NewMap(cm2)
	require.NoError(t, err)

	cr2 := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Inputs: iwrapped2,
		Config: cwrapped2,
	}

	callback1 := make(chan capabilities.CapabilityResponse, 10)
	callback2 := make(chan capabilities.CapabilityResponse, 10)

	require.NoError(t, ts.RegisterTrigger(ctx, callback1, cr1))
	require.NoError(t, ts.RegisterTrigger(ctx, callback2, cr2))

	// Send events to trigger and check for them in the callback
	fr1 := []FeedReport{
		{
			FeedID:               mercury.FeedID(feedOne).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       20,
			ObservationTimestamp: 5,
		},
		{
			FeedID:               mercury.FeedID(feedThree).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       25,
			ObservationTimestamp: 8,
		},
		{
			FeedID:               mercury.FeedID(feedTwo).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       30,
			ObservationTimestamp: 10,
		},
		{
			FeedID:               mercury.FeedID(feedFour).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       40,
			ObservationTimestamp: 15,
		},
	}

	mfr1 := []mercury.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       20,
			ObservationTimestamp: 5,
		},
		{
			FeedID:               feedThree,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       25,
			ObservationTimestamp: 8,
		},
		{
			FeedID:               feedTwo,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       30,
			ObservationTimestamp: 10,
		},
		{
			FeedID:               feedFour,
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
	payload = append(payload, mfr1[0], mfr1[1], mfr1[3])
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.BatchedPayload, 3)
	assert.Equal(t, mfr1[0], unwrapped.BatchedPayload[feedOne])
	assert.Equal(t, mfr1[1], unwrapped.BatchedPayload[feedThree])
	assert.Equal(t, mfr1[3], unwrapped.BatchedPayload[feedFour])

	msg = <-callback2
	unwrapped, _ = mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	payload = make([]mercury.FeedReport, 0)
	payload = append(payload, mfr1[1], mfr1[2]) // Because GenerateTriggerEventID sorts the reports by feedID, this works
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.BatchedPayload, 2)
	assert.Equal(t, mfr1[2], unwrapped.BatchedPayload[feedTwo])
	assert.Equal(t, mfr1[1], unwrapped.BatchedPayload[feedThree])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr1))
	fr2 := []FeedReport{
		{
			FeedID:               mercury.FeedID(feedThree).Bytes(),
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       50,
			ObservationTimestamp: 20,
		},
	}
	mfr2 := []mercury.FeedReport{
		{
			FeedID:               feedThree,
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
	unwrapped, err = mercury.Codec{}.UnwrapMercuryTriggerEvent(msg.Value)
	require.NoError(t, err)
	assert.Equal(t, "mercury", unwrapped.TriggerType)
	payload = make([]mercury.FeedReport, 0)
	payload = append(payload, mfr2[0])
	assert.Equal(t, GenerateTriggerEventID(payload), unwrapped.ID)
	assert.Len(t, unwrapped.BatchedPayload, 1)
	assert.Equal(t, mfr2[0], unwrapped.BatchedPayload[feedThree])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr2))
	err = ts.ProcessReport(fr1)
	assert.NoError(t, err)
	assert.Len(t, callback1, 0)
	assert.Len(t, callback2, 0)
}
