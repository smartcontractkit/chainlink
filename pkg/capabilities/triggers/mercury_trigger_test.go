package triggers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// registerTrigger will do the following:
//
//  1. Register a trigger with the given feedIDs and triggerID
//  2. Return the trigger events channel, registerUnregisterRequest, and test context
func registerTrigger(
	ctx context.Context,
	t *testing.T,
	ts *MercuryTriggerService,
	feedIDs []string,
	triggerID string,
) (
	triggerEventsCh chan capabilities.CapabilityResponse,
	unregisterRequest capabilities.CapabilityRequest,
) {

	inputs, err := values.NewMap(map[string]interface{}{
		"triggerId": triggerID,
	})
	require.NoError(t, err)

	config, err := values.NewMap(map[string]interface{}{
		"feedIds":        feedIDs,
		"maxFrequencyMs": 100,
	})
	require.NoError(t, err)

	requestMetadata := capabilities.RequestMetadata{
		WorkflowID: "workflow-id-1",
	}
	registerRequest := capabilities.CapabilityRequest{
		Metadata: requestMetadata,
		Inputs:   inputs,
		Config:   config,
	}
	triggerEventsCh = make(chan capabilities.CapabilityResponse, 1000)
	require.NoError(t, ts.RegisterTrigger(ctx, triggerEventsCh, registerRequest))

	unregisterRequest = capabilities.CapabilityRequest{
		Metadata: requestMetadata,
		Inputs:   inputs,
	}

	return triggerEventsCh, unregisterRequest
}

var (
	feedOne   = "0x1111111111111111111100000000000000000000000000000000000000000000"
	feedTwo   = "0x2222222222222222222200000000000000000000000000000000000000000000"
	feedThree = "0x3333333333333333333300000000000000000000000000000000000000000000"
	feedFour  = "0x4444444444444444444400000000000000000000000000000000000000000000"
	feedFive  = "0x5555555555555555555500000000000000000000000000000000000000000000"
)

func TestMercuryTrigger(t *testing.T) {
	ts := NewMercuryTriggerService(100, logger.Nop())
	ctx := tests.Context(t)
	err := ts.Start(ctx)
	require.NoError(t, err)
	// use registerTriggerHelper to register a trigger
	callback, registerUnregisterRequest := registerTrigger(
		ctx,
		t,
		ts,
		[]string{feedOne},
		"test-id-1",
	)

	// Send events to trigger and check for them in the callback
	mfr := []mercury.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       2,
			ObservationTimestamp: 3,
		},
	}
	err = ts.ProcessReport(mfr)
	assert.NoError(t, err)
	msg := <-callback
	triggerEvent, reports := upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, "mercury", triggerEvent.TriggerType)
	assert.Len(t, reports, 1)
	assert.Equal(t, mfr[0], reports[0])

	// Unregister the trigger and check that events no longer go on the callback
	require.NoError(t, ts.UnregisterTrigger(ctx, registerUnregisterRequest))
	err = ts.ProcessReport(mfr)
	require.NoError(t, err)
	require.Len(t, callback, 0)
	require.NoError(t, ts.Close())
}

func TestMultipleMercuryTriggers(t *testing.T) {
	ts := NewMercuryTriggerService(100, logger.Nop())
	ctx := tests.Context(t)
	err := ts.Start(ctx)
	require.NoError(t, err)
	callback1, cr1 := registerTrigger(
		ctx,
		t,
		ts,
		[]string{
			feedOne,
			feedThree,
			feedFour,
		},
		"test-id-1",
	)

	callback2, cr2 := registerTrigger(
		ctx,
		t,
		ts,
		[]string{
			feedTwo,
			feedThree,
			feedFive,
		},
		"test-id-2",
	)

	// Send events to trigger and check for them in the callback
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

	err = ts.ProcessReport(mfr1)
	assert.NoError(t, err)

	msg := <-callback1
	triggerEvent, reports := upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, "mercury", triggerEvent.TriggerType)
	assert.Len(t, reports, 3)
	assert.Equal(t, mfr1[0], reports[0])
	assert.Equal(t, mfr1[1], reports[1])
	assert.Equal(t, mfr1[3], reports[2])

	msg = <-callback2
	triggerEvent, reports = upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, "mercury", triggerEvent.TriggerType)
	assert.Len(t, reports, 2)
	assert.Equal(t, mfr1[2], reports[0])
	assert.Equal(t, mfr1[1], reports[1])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr1))
	mfr2 := []mercury.FeedReport{
		{
			FeedID:               feedThree,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       50,
			ObservationTimestamp: 20,
		},
	}
	err = ts.ProcessReport(mfr2)
	assert.NoError(t, err)

	retryCount := 0
	for rMsg := range callback2 {
		triggerEvent, reports = upwrapTriggerEvent(t, rMsg.Value)
		require.NoError(t, err)
		require.Len(t, reports, 2)
		require.Equal(t, "mercury", triggerEvent.TriggerType)
		if reports[1].BenchmarkPrice == 50 {
			// expect to eventually get updated feed value
			break
		}
		require.True(t, retryCount < 100)
		retryCount++
	}

	require.NoError(t, ts.UnregisterTrigger(ctx, cr2))
	err = ts.ProcessReport(mfr1)
	assert.NoError(t, err)
	assert.Len(t, callback1, 0)
	assert.Len(t, callback2, 0)
	require.NoError(t, ts.Close())
}

func TestMercuryTrigger_RegisterTriggerErrors(t *testing.T) {
	ts := NewMercuryTriggerService(100, logger.Nop())
	ctx := tests.Context(t)
	require.NoError(t, ts.Start(ctx))

	im := map[string]interface{}{
		"triggerId": "test-id-1",
	}
	inputsWrapped, err := values.NewMap(im)
	require.NoError(t, err)

	cm := map[string]interface{}{
		"feedIds":        []string{feedOne},
		"maxFrequencyMs": 90,
	}
	configWrapped, err := values.NewMap(cm)
	require.NoError(t, err)

	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "workflow-id-1",
		},
		Config: configWrapped,
		Inputs: inputsWrapped,
	}
	callback := make(chan capabilities.CapabilityResponse)
	require.Error(t, ts.RegisterTrigger(ctx, callback, cr))

	cm = map[string]interface{}{
		"feedIds":        []string{feedOne},
		"maxFrequencyMs": 0,
	}
	configWrapped, err = values.NewMap(cm)
	require.NoError(t, err)
	cr.Config = configWrapped
	require.Error(t, ts.RegisterTrigger(ctx, callback, cr))

	cm = map[string]interface{}{
		"feedIds":        []string{},
		"maxFrequencyMs": 1000,
	}
	configWrapped, err = values.NewMap(cm)
	require.NoError(t, err)
	cr.Config = configWrapped
	require.Error(t, ts.RegisterTrigger(ctx, callback, cr))

	require.NoError(t, ts.Close())
}

func TestGetNextWaitIntervalMs(t *testing.T) {
	// getNextWaitIntervalMs args = (lastTs, tickerResolutionMs, currentTs)

	// expected cases
	assert.Equal(t, int64(900), getNextWaitIntervalMs(12000, 1000, 12100))
	assert.Equal(t, int64(200), getNextWaitIntervalMs(12000, 1000, 12800))

	// slow processing
	assert.Equal(t, int64(0), getNextWaitIntervalMs(12000, 1000, 13000))
	assert.Equal(t, int64(0), getNextWaitIntervalMs(12000, 1000, 14600))
}

func upwrapTriggerEvent(t *testing.T, wrappedEvent values.Value) (capabilities.TriggerEvent, []mercury.FeedReport) {
	event := capabilities.TriggerEvent{}
	err := wrappedEvent.UnwrapTo(&event)
	require.NoError(t, err)
	require.NotNil(t, event.Payload)
	mercuryReports, err := mercury.Codec{}.Unwrap(event.Payload)
	require.NoError(t, err)
	return event, mercuryReports
}
