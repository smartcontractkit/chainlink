package common_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

type requestState struct {
	counter int
}

func TestRequestCache_Simple(t *testing.T) {
	t.Parallel()

	cache := common.NewRequestCache[requestState](time.Hour, 1000)
	callback := common.NewCallback()

	req := &api.Message{Body: api.MessageBody{MessageId: "aa", Sender: "0x1234"}}
	initialState := &requestState{}
	lggr := logger.Test(t)
	require.NoError(t, cache.NewRequest(lggr, req, callback, initialState))

	nodeResp := &api.Message{Body: api.MessageBody{MessageId: "aa", Receiver: "0x1234"}}
	go func() {
		assert.NoError(t, cache.ProcessResponse(nodeResp, func(response *api.Message, responseData *requestState) (aggregated *handlers.UserCallbackPayload, newResponseData *requestState, err error) {
			// ready after first response
			var rawResponse json.RawMessage
			rawResponse, err = json.Marshal(response)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
			}
			return &handlers.UserCallbackPayload{RawResponse: rawResponse}, nil, nil
		}))
	}()
	finalResp, err := callback.Wait(t.Context())
	require.NoError(t, err)
	var msg api.Message
	require.NoError(t, json.Unmarshal(finalResp.RawResponse, &msg))
	require.Equal(t, "aa", msg.Body.MessageId)
}

func TestRequestCache_MultiResponse(t *testing.T) {
	t.Parallel()

	nRequests := 10
	nResponsesPerRequest := 100
	maxDelayMillis := 100

	lggr := logger.Test(t)
	cache := common.NewRequestCache[requestState](time.Hour, 1000)
	cbs := make([]*common.Callback, nRequests)
	reqs := make([]*api.Message, nRequests)
	for i := range nRequests {
		cb := common.NewCallback()
		cbs[i] = cb
		reqs[i] = &api.Message{Body: api.MessageBody{MessageId: "abcd", Sender: fmt.Sprintf("sender_%d", i)}}
		initialState := &requestState{counter: 0}
		require.NoError(t, cache.NewRequest(lggr, reqs[i], cbs[i], initialState))
	}

	for i := range nRequests {
		resp := &api.Message{Body: api.MessageBody{MessageId: "abcd"}}
		resp.Body.Receiver = reqs[i].Body.Sender
		for range nResponsesPerRequest {
			go func() {
				n := rand.Intn(maxDelayMillis) + 1
				time.Sleep(time.Duration(n) * time.Millisecond)
				assert.NoError(t, cache.ProcessResponse(resp, func(response *api.Message, responseData *requestState) (aggregated *handlers.UserCallbackPayload, newResponseData *requestState, err error) {
					responseData.counter++
					if responseData.counter == nResponsesPerRequest {
						var rawResponse json.RawMessage
						rawResponse, err = json.Marshal(response)
						if err != nil {
							return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
						}
						return &handlers.UserCallbackPayload{RawResponse: rawResponse}, nil, nil
					}
					return nil, responseData, nil
				}))
			}()
		}
	}

	for i := range nRequests {
		resp, err := cbs[i].Wait(t.Context())
		require.NoError(t, err)
		var msg api.Message
		require.NoError(t, json.Unmarshal(resp.RawResponse, &msg))
		require.Equal(t, "abcd", msg.Body.MessageId)
		require.Equal(t, reqs[i].Body.Sender, msg.Body.Receiver)
	}
}

func TestRequestCache_Timeout(t *testing.T) {
	t.Parallel()

	cache := common.NewRequestCache[requestState](time.Millisecond*10, 1000)
	callback := common.NewCallback()
	lggr := logger.Test(t)

	req := &api.Message{Body: api.MessageBody{MessageId: "aa", Sender: "0x1234"}}
	initialState := &requestState{}
	require.NoError(t, cache.NewRequest(lggr, req, callback, initialState))

	finalResp, err := callback.Wait(t.Context())
	require.NoError(t, err)
	codec := api.JsonRPCCodec{}
	rawResp, err := codec.DecodeLegacyResponse(finalResp.RawResponse)
	require.NoError(t, err)
	require.Equal(t, "aa", rawResp.Body.MessageId)
	require.Equal(t, api.RequestTimeoutError, finalResp.ErrorCode)
}

func TestRequestCache_MaxSize(t *testing.T) {
	t.Parallel()

	cache := common.NewRequestCache[requestState](time.Hour, 2)
	callback := common.NewCallback()
	lggr := logger.Test(t)
	initialState := &requestState{}

	req := &api.Message{Body: api.MessageBody{MessageId: "aa", Sender: "0x1234"}}
	require.NoError(t, cache.NewRequest(lggr, req, callback, initialState))

	req.Body.MessageId = "bb"
	require.NoError(t, cache.NewRequest(lggr, req, callback, initialState))

	req.Body.MessageId = "cc"
	require.Error(t, cache.NewRequest(lggr, req, callback, initialState))
}
