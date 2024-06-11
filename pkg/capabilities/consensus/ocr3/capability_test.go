package ocr3

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const workflowTestID = "consensus-workflow-test-id-1"
const workflowExecutionTestID = "consensus-workflow-execution-test-id-1"
const workflowTestName = "consensus-workflow-test-name-1"
const reportTestId = "rep-id-1"

type mockAggregator struct {
	types.Aggregator
}

func mockAggregatorFactory(_ string, _ values.Map, _ logger.Logger) (types.Aggregator, error) {
	return &mockAggregator{}, nil
}

type encoder struct {
	types.Encoder
}

func mockEncoderFactory(_ *values.Map) (types.Encoder, error) {
	return &encoder{}, nil
}

func TestOCR3Capability_Schema(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Nop()

	s := newStore()
	s.evictedCh = make(chan *request)

	cp := newCapability(s, fc, 1*time.Second, mockAggregatorFactory, mockEncoderFactory, lggr, 10)
	schema, err := cp.Schema()
	require.NoError(t, err)

	var shouldUpdate = false
	if shouldUpdate {
		err = os.WriteFile("./testdata/fixtures/capability/schema.json", []byte(schema), 0600)
		require.NoError(t, err)
	}

	fixture, err := os.ReadFile("./testdata/fixtures/capability/schema.json")
	require.NoError(t, err)

	utils.AssertJSONEqual(t, fixture, []byte(schema))
}

func TestOCR3Capability(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Test(t)

	ctx := tests.Context(t)

	s := newStore()
	s.evictedCh = make(chan *request)

	cp := newCapability(s, fc, 1*time.Second, mockAggregatorFactory, mockEncoderFactory, lggr, 10)
	require.NoError(t, cp.Start(ctx))

	config, err := values.NewMap(
		map[string]any{
			"aggregation_method": "data_feeds",
			"aggregation_config": map[string]any{},
			"encoder_config":     map[string]any{},
			"encoder":            "evm",
			"report_id":          "ffff",
		},
	)
	require.NoError(t, err)

	ethUsdValStr := "1.123456"
	ethUsdValue, err := decimal.NewFromString(ethUsdValStr)
	require.NoError(t, err)
	observationKey := "ETH_USD"
	obs := []any{map[string]any{observationKey: ethUsdValue}}
	inputs, err := values.NewMap(map[string]any{"observations": obs})
	require.NoError(t, err)

	executeReq := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowTestID,
			WorkflowExecutionID: workflowExecutionTestID,
		},
		Config: config,
		Inputs: inputs,
	}
	callback, err := cp.Execute(ctx, executeReq)
	require.NoError(t, err)

	obsv, err := values.NewList(obs)
	require.NoError(t, err)

	// Mock the oracle returning a response
	cp.transmitCh <- &outputs{
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: obsv,
		},
		WorkflowExecutionID: workflowExecutionTestID,
	}
	require.NoError(t, err)

	expectedCapabilityResponse := capabilities.CapabilityResponse{
		Value: obsv,
	}

	gotR := <-s.evictedCh
	assert.Equal(t, workflowExecutionTestID, gotR.WorkflowExecutionID)

	// Test that our unwrapping works
	var actualUnwrappedObs []map[string]decimal.Decimal
	err = gotR.Observations.UnwrapTo(&actualUnwrappedObs)
	assert.NoError(t, err)
	assert.Len(t, actualUnwrappedObs, 1)
	actualObs, ok := actualUnwrappedObs[0][observationKey]
	assert.True(t, ok)
	assert.Equal(t, ethUsdValStr, actualObs.String())

	assert.Equal(t, expectedCapabilityResponse, <-callback)
}

func TestOCR3Capability_Eviction(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Test(t)

	ctx := tests.Context(t)
	rea := time.Second
	s := newStore()
	cp := newCapability(s, fc, rea, mockAggregatorFactory, mockEncoderFactory, lggr, 10)
	require.NoError(t, cp.Start(ctx))

	config, err := values.NewMap(
		map[string]any{
			"aggregation_method": "data_feeds",
			"aggregation_config": map[string]any{},
			"encoder_config":     map[string]any{},
			"encoder":            "evm",
			"report_id":          "aaaa",
		},
	)
	require.NoError(t, err)

	ethUsdValue, err := decimal.NewFromString("1.123456")
	require.NoError(t, err)
	inputs, err := values.NewMap(map[string]any{"observations": []any{map[string]any{"ETH_USD": ethUsdValue}}})
	require.NoError(t, err)

	rid := uuid.New().String()
	executeReq := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowTestID,
			WorkflowExecutionID: rid,
		},
		Config: config,
		Inputs: inputs,
	}

	callback, err := cp.Execute(ctx, executeReq)
	require.NoError(t, err)

	fc.Advance(1 * time.Hour)
	resp := <-callback
	assert.ErrorContains(t, resp.Err, "timeout exceeded: could not process request before expiry")

	_, ok := s.requests[rid]
	assert.False(t, ok)
}

func TestOCR3Capability_EvictionUsingConfig(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Test(t)

	ctx := tests.Context(t)
	// This is the default expired at
	rea := time.Hour
	s := newStore()
	cp := newCapability(s, fc, rea, mockAggregatorFactory, mockEncoderFactory, lggr, 10)
	require.NoError(t, cp.Start(ctx))

	config, err := values.NewMap(
		map[string]any{
			"aggregation_method": "data_feeds",
			"aggregation_config": map[string]any{},
			"encoder_config":     map[string]any{},
			"encoder":            "evm",
			"report_id":          "aaaa",
			"request_timeout_ms": 10000,
		},
	)
	require.NoError(t, err)

	ethUsdValue, err := decimal.NewFromString("1.123456")
	require.NoError(t, err)
	inputs, err := values.NewMap(map[string]any{"observations": []any{map[string]any{"ETH_USD": ethUsdValue}}})
	require.NoError(t, err)

	rid := uuid.New().String()
	executeReq := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowTestID,
			WorkflowExecutionID: rid,
		},
		Config: config,
		Inputs: inputs,
	}

	callback, err := cp.Execute(ctx, executeReq)
	require.NoError(t, err)

	// 1 minute is more than the config timeout we provided, but less than
	// the hardcoded timeout.
	fc.Advance(1 * time.Minute)
	resp := <-callback
	assert.ErrorContains(t, resp.Err, "timeout exceeded: could not process request before expiry")

	_, ok := s.requests[rid]
	assert.False(t, ok)
}

func TestOCR3Capability_Registration(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Test(t)

	ctx := tests.Context(t)
	s := newStore()
	cp := newCapability(s, fc, 1*time.Second, mockAggregatorFactory, mockEncoderFactory, lggr, 10)
	require.NoError(t, cp.Start(ctx))

	config, err := values.NewMap(map[string]any{
		"aggregation_method": "data_feeds",
		"aggregation_config": map[string]any{},
		"encoder":            "",
		"encoder_config":     map[string]any{},
		"report_id":          "000f",
	})
	require.NoError(t, err)

	registerReq := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: workflowTestID,
		},
		Config: config,
	}

	err = cp.RegisterToWorkflow(ctx, registerReq)
	require.NoError(t, err)

	agg, err := cp.getAggregator(workflowTestID)
	require.NoError(t, err)
	assert.NotNil(t, agg)

	unregisterReq := capabilities.UnregisterFromWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: workflowTestID,
		},
	}

	err = cp.UnregisterFromWorkflow(ctx, unregisterReq)
	require.NoError(t, err)

	_, err = cp.getAggregator(workflowTestID)
	assert.ErrorContains(t, err, "no aggregator found for")
}

func TestOCR3Capability_ValidateConfig(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)
	lggr := logger.Test(t)

	s := newStore()
	s.evictedCh = make(chan *request)

	o := newCapability(s, fc, 1*time.Second, mockAggregatorFactory, mockEncoderFactory, lggr, 10)

	t.Run("ValidConfig", func(t *testing.T) {
		config, err := values.NewMap(map[string]any{
			"aggregation_method": "data_feeds",
			"aggregation_config": map[string]any{},
			"encoder":            "",
			"encoder_config":     map[string]any{},
			"report_id":          "aaaa",
		})
		require.NoError(t, err)

		c, err := o.ValidateConfig(config)
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("InvalidConfig null", func(t *testing.T) {
		config, err := values.NewMap(map[string]any{
			"aggregation_method": "data_feeds",
			"report_id":          "aaaa",
		})
		require.NoError(t, err)

		c, err := o.ValidateConfig(config)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected object, but got null") // taken from the error json schema error message
		require.Nil(t, c)
	})

	t.Run("InvalidConfig illegal report_id", func(t *testing.T) {
		config, err := values.NewMap(map[string]any{
			"aggregation_method": "data_feeds",
			"aggregation_config": map[string]any{},
			"encoder":            "",
			"encoder_config":     map[string]any{},
			"report_id":          "aa",
		})
		require.NoError(t, err)

		c, err := o.ValidateConfig(config)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not match pattern") // taken from the error json schema error message
		require.Nil(t, c)
	})
}
