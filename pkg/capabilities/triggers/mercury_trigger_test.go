package triggers

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
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
	<-chan capabilities.CapabilityResponse,
	capabilities.CapabilityRequest,
) {
	var unregisterRequest capabilities.CapabilityRequest

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
	triggerEventsCh, err := ts.RegisterTrigger(ctx, registerRequest)
	require.NoError(t, err)

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
	mfr := []datastreams.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(2).Bytes(),
			ObservationTimestamp: 3,
			Signatures:           [][]byte{},
		},
	}
	err = ts.ProcessReport(mfr)
	assert.NoError(t, err)
	msg := <-callback
	triggerEvent, reports := upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, triggerID, triggerEvent.TriggerType)
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
	mfr1 := []datastreams.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(20).Bytes(),
			ObservationTimestamp: 5,
			Signatures:           [][]byte{},
		},
		{
			FeedID:               feedThree,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(25).Bytes(),
			ObservationTimestamp: 8,
			Signatures:           [][]byte{},
		},
		{
			FeedID:               feedTwo,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(30).Bytes(),
			ObservationTimestamp: 10,
			Signatures:           [][]byte{},
		},
		{
			FeedID:               feedFour,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(40).Bytes(),
			ObservationTimestamp: 15,
			Signatures:           [][]byte{},
		},
	}

	err = ts.ProcessReport(mfr1)
	assert.NoError(t, err)

	msg := <-callback1
	triggerEvent, reports := upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, triggerID, triggerEvent.TriggerType)
	assert.Len(t, reports, 3)
	assert.Equal(t, mfr1[0], reports[0])
	assert.Equal(t, mfr1[1], reports[1])
	assert.Equal(t, mfr1[3], reports[2])

	msg = <-callback2
	triggerEvent, reports = upwrapTriggerEvent(t, msg.Value)
	assert.Equal(t, triggerID, triggerEvent.TriggerType)
	assert.Len(t, reports, 2)
	assert.Equal(t, mfr1[2], reports[0])
	assert.Equal(t, mfr1[1], reports[1])

	require.NoError(t, ts.UnregisterTrigger(ctx, cr1))
	mfr2 := []datastreams.FeedReport{
		{
			FeedID:               feedThree,
			FullReport:           []byte("0x1234"),
			BenchmarkPrice:       big.NewInt(50).Bytes(),
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
		require.Equal(t, triggerID, triggerEvent.TriggerType)
		price := big.NewInt(0).SetBytes(reports[1].BenchmarkPrice)
		if price.Cmp(big.NewInt(50)) == 0 {
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
	_, err = ts.RegisterTrigger(ctx, cr)
	require.Error(t, err)

	cm = map[string]interface{}{
		"feedIds":        []string{feedOne},
		"maxFrequencyMs": 0,
	}
	configWrapped, err = values.NewMap(cm)
	require.NoError(t, err)
	cr.Config = configWrapped
	_, err = ts.RegisterTrigger(ctx, cr)
	require.Error(t, err)

	cm = map[string]interface{}{
		"feedIds":        []string{},
		"maxFrequencyMs": 1000,
	}
	configWrapped, err = values.NewMap(cm)
	require.NoError(t, err)
	cr.Config = configWrapped
	_, err = ts.RegisterTrigger(ctx, cr)
	require.Error(t, err)

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

func upwrapTriggerEvent(t *testing.T, wrappedEvent values.Value) (capabilities.TriggerEvent, []datastreams.FeedReport) {
	event := capabilities.TriggerEvent{}
	err := wrappedEvent.UnwrapTo(&event)
	require.NoError(t, err)
	require.NotNil(t, event.Payload)
	mercuryReports, err := testMercuryCodec{}.Unwrap(event.Payload)
	require.NoError(t, err)
	return event, mercuryReports
}

func TestMercuryTrigger_ConfigValidation(t *testing.T) {
	var newConfig = func(t *testing.T, feedIDs []string, maxFrequencyMs int) *values.Map {
		cm := map[string]interface{}{
			"feedIds":        feedIDs,
			"maxFrequencyMs": maxFrequencyMs,
		}
		configWrapped, err := values.NewMap(cm)
		require.NoError(t, err)

		return configWrapped
	}

	var newConfigSingleFeed = func(t *testing.T, feedID string) *values.Map {
		return newConfig(t, []string{feedID}, 1000)
	}

	ts := NewMercuryTriggerService(1000, logger.Nop())
	rawConf := newConfigSingleFeed(t, "012345678901234567890123456789012345678901234567890123456789000000")
	conf, err := ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	rawConf = newConfigSingleFeed(t, "0x1234")
	conf, err = ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	rawConf = newConfigSingleFeed(t, "0x123zzz")
	conf, err = ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	rawConf = newConfigSingleFeed(t, "0x0001013ebd4ed3f5889FB5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	conf, err = ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	passingFeedID := "0x0001013ebd4ed3f5889fb5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292"
	// test maxfreq < 1
	rawConf = newConfig(t, []string{passingFeedID}, 0)
	conf, err = ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	rawConf = newConfig(t, []string{passingFeedID}, -1)
	conf, err = ts.ValidateConfig(rawConf)
	require.Error(t, err)
	require.Empty(t, conf)

	rawConf = newConfigSingleFeed(t, passingFeedID)
	conf, err = ts.ValidateConfig(rawConf)
	require.NoError(t, err)
	require.NotEmpty(t, conf)
}

func TestMercuryTrigger_GenerateSchema(t *testing.T) {
	ts := NewMercuryTriggerService(1000, logger.Nop())
	schema, err := ts.Schema()
	require.NoError(t, err)
	var shouldUpdate = false
	if shouldUpdate {
		err = os.WriteFile("./testdata/fixtures/mercury/schema.json", []byte(schema), 0600)
		require.NoError(t, err)
	}

	fixture, err := os.ReadFile("./testdata/fixtures/mercury/schema.json")
	require.NoError(t, err)

	utils.AssertJSONEqual(t, fixture, []byte(schema))
}
