package requests_test

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func Test_Handler_SendsResponse(t *testing.T) {
	lggr := logger.Test(t)
	ctx := tests.Context(t)

	h := requests.NewHandler(lggr, requests.NewStore(), clockwork.NewFakeClockAt(time.Now()), 1*time.Second)
	servicetest.Run(t, h)

	responseCh := make(chan capabilities.CapabilityResponse, 10)
	h.SendRequest(ctx, &requests.Request{
		WorkflowExecutionID: "test",
		CallbackCh:          responseCh,
		ExpiresAt:           time.Now().Add(1 * time.Hour),
	})

	testVal, err := values.Wrap("testval")
	require.NoError(t, err)

	h.SendResponse(ctx, &requests.Response{
		WorkflowExecutionID: "test",
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: testVal,
			Err:   nil,
		},
	})

	resp := <-responseCh
	require.Equal(t, testVal, resp.Value)

}

func Test_Handler_SendsResponseToLateRequest(t *testing.T) {
	lggr := logger.Test(t)
	ctx := tests.Context(t)

	h := requests.NewHandler(lggr, requests.NewStore(), clockwork.NewFakeClockAt(time.Now()), 1*time.Second)
	servicetest.Run(t, h)

	testVal, err := values.Wrap("testval")
	require.NoError(t, err)

	h.SendResponse(ctx, &requests.Response{
		WorkflowExecutionID: "test",
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: testVal,
			Err:   nil,
		},
	})

	responseCh := make(chan capabilities.CapabilityResponse, 10)
	h.SendRequest(ctx, &requests.Request{
		WorkflowExecutionID: "test",
		CallbackCh:          responseCh,
		ExpiresAt:           time.Now().Add(1 * time.Hour),
	})

	resp := <-responseCh
	require.Equal(t, testVal, resp.Value)

}

func Test_Handler_SendsResponseToLateRequestOnlyOnce(t *testing.T) {
	lggr := logger.Test(t)
	ctx := tests.Context(t)

	h := requests.NewHandler(lggr, requests.NewStore(), clockwork.NewFakeClockAt(time.Now()), 1*time.Second)
	servicetest.Run(t, h)

	testVal, err := values.Wrap("testval")
	require.NoError(t, err)

	h.SendResponse(ctx, &requests.Response{
		WorkflowExecutionID: "test",
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: testVal,
			Err:   nil,
		},
	})

	responseCh := make(chan capabilities.CapabilityResponse, 10)
	h.SendRequest(ctx, &requests.Request{
		WorkflowExecutionID: "test",
		CallbackCh:          responseCh,
		ExpiresAt:           time.Now().Add(1 * time.Hour),
	})

	require.NoError(t, err)

	resp := <-responseCh
	require.Equal(t, testVal, resp.Value)

	responseCh = make(chan capabilities.CapabilityResponse, 10)
	h.SendRequest(ctx, &requests.Request{
		WorkflowExecutionID: "test",
		CallbackCh:          responseCh,
		ExpiresAt:           time.Now().Add(1 * time.Hour),
	})

	select {
	case <-responseCh:
		t.Fatal("Should not have received a response")
	default:
	}

}

func Test_Handler_PendingRequestsExpiry(t *testing.T) {
	ctx := tests.Context(t)

	lggr := logger.Test(t)
	clock := clockwork.NewFakeClockAt(time.Now())
	h := requests.NewHandler(lggr, requests.NewStore(), clock, 1*time.Second)
	servicetest.Run(t, h)

	responseCh := make(chan capabilities.CapabilityResponse, 10)
	h.SendRequest(ctx, &requests.Request{
		WorkflowExecutionID: "test",
		CallbackCh:          responseCh,
		ExpiresAt:           time.Now().Add(1 * time.Second),
	})

	clock.Advance(2 * time.Second)

	resp := <-responseCh

	assert.ErrorContains(t, resp.Err, "timeout exceeded: could not process request before expiry")
}
