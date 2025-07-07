package capabilities

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

const (
	defaultSendChannelBufferSize = 1000
	privateKey1                  = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f"
	privateKey2                  = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023e"
	triggerID1                   = "5"
	triggerID2                   = "6"
	workflowID1                  = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1         = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1                       = "0x00000000000000000000000000000000000000aa"
	address1                     = "0x853d51d5d9935964267a5050aC53aa63ECA39bc5"
)

func setupHandler(t *testing.T) (*handler, *mocks.HTTPClient, *handlermocks.DON, []gwcommon.TestNode) {
	lggr := logger.Test(t)
	httpClient := mocks.NewHTTPClient(t)
	don := handlermocks.NewDON(t)
	nodeRateLimiterConfig := ratelimit.RateLimiterConfig{
		GlobalRPS:      100.0,
		GlobalBurst:    100,
		PerSenderRPS:   100.0,
		PerSenderBurst: 100,
	}
	handlerConfig := HandlerConfig{
		NodeRateLimiter:         nodeRateLimiterConfig,
		MaxAllowedMessageAgeSec: 30,
	}

	cfgBytes, err := json.Marshal(handlerConfig)
	require.NoError(t, err)
	donConfig := &config.DONConfig{
		Members: []config.NodeConfig{},
		F:       1,
	}
	nodes := gwcommon.NewTestNodes(t, 2)
	for id, n := range nodes {
		donConfig.Members = append(donConfig.Members, config.NodeConfig{
			Name:    fmt.Sprintf("node_%d", id),
			Address: n.Address,
		})
	}
	handler, err := NewHandler(json.RawMessage(cfgBytes), donConfig, don, httpClient, lggr)
	require.NoError(t, err)
	return handler, httpClient, don, nodes
}

func TestHandler_SendHTTPMessageToClient(t *testing.T) {
	handler, httpClient, don, nodes := setupHandler(t)
	ctx := testutils.Context(t)
	nodeAddr := nodes[0].Address
	payload := Request{
		Method:    "GET",
		URL:       "http://example.com",
		Headers:   map[string]string{},
		Body:      nil,
		TimeoutMs: 2000,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: "123",
			Method:    MethodWebAPITarget,
			DonId:     "testDonId",
			Payload:   json.RawMessage(payloadBytes),
		},
	}
	err = msg.Sign(nodes[0].PrivateKey)
	require.NoError(t, err)
	err = msg.Validate()
	require.NoError(t, err)
	t.Run("happy case", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{},
			Body:       []byte("response body"),
		}, nil).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodWebAPITarget == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				200 == payload.StatusCode &&
				0 == len(payload.Headers) &&
				string(payload.Body) == "response body" &&
				!payload.ExecutionError
		})).Return(nil).Once()
		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("http client non-HTTP error", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 404,
			Headers:    map[string]string{},
			Body:       []byte("access denied"),
		}, nil).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodWebAPITarget == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				404 == payload.StatusCode &&
				string(payload.Body) == "access denied" &&
				0 == len(payload.Headers) &&
				!payload.ExecutionError
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// // ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("http client non-HTTP error", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, errors.New("error while marshalling")).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodWebAPITarget == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				payload.ExecutionError &&
				"error while marshalling" == payload.ErrorMessage
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// // ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})
}

func triggerRequest(t *testing.T, key *ecdsa.PrivateKey, topics []string, methodName string, timestamp string, payload string) *api.Message {
	messageID := "12345"
	if methodName == "" {
		methodName = MethodWebAPITrigger
	}
	if timestamp == "" {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}
	donID := "workflow_don_1"
	var payloadJSON []byte
	if payload == "" {
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		require.NoError(t, err)
		reqPayload := webapicap.TriggerRequestPayload{
			TriggerId:      "web-api-trigger@1.0.0",
			TriggerEventId: "action_1234567890",
			Timestamp:      ts,
			Topics:         topics,
			Params: webapicap.TriggerRequestPayloadParams(map[string]interface{}{
				"bid": "101",
				"ask": "102",
			}),
		}
		payloadJSON, err = json.Marshal(reqPayload)
		require.NoError(t, err)
	} else {
		payloadJSON = []byte(payload)
	}
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: messageID,
			Method:    methodName,
			DonId:     donID,
			Payload:   json.RawMessage(payloadJSON),
		},
	}
	err := msg.Sign(key)
	require.NoError(t, err)
	err = msg.Validate()
	require.NoError(t, err)
	return msg
}

func requireNoChanMsg[T any](t *testing.T, ch <-chan T) {
	timedOut := false
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		timedOut = true
	}
	require.True(t, timedOut)
}

func TestHandlerReceiveHTTPMessageFromClient(t *testing.T) {
	handler, _, don, nodes := setupHandler(t)
	ctx := testutils.Context(t)
	msg := triggerRequest(t, nodes[0].PrivateKey, []string{"daily_price_update"}, "", "", "")
	codec := api.JsonRPCCodec{}

	t.Run("happy case", func(t *testing.T) {
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)

		// sends to 2 dons
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			nodeReq := nodeRequest(msg)
			require.Equal(t, nodeReq, args.Get(2))
		}).Return(nil).Once()
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			nodeReq := nodeRequest(msg)
			require.Equal(t, nodeReq, args.Get(2))
		}).Return(nil).Once()

		err := handler.HandleLegacyUserMessage(ctx, msg, ch)
		require.NoError(t, err)
		requireNoChanMsg(t, ch)

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodes[0].Address)
		require.NoError(t, err)

		userPayload := <-ch
		require.Equal(t, handlers.UserCallbackPayload{RawResponse: codec.EncodeLegacyResponse(msg), ErrorCode: api.NoError}, userPayload)
		_, open := <-ch
		require.False(t, open)
	})

	t.Run("sad case invalid method", func(t *testing.T) {
		invalidMsg := triggerRequest(t, nodes[0].PrivateKey, []string{"daily_price_update"}, "foo", "", "")
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleLegacyUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch

		require.Equal(t, handlers.UserCallbackPayload{
			RawResponse: codec.EncodeNewErrorResponse(
				invalidMsg.Body.MessageId,
				api.ToJSONRPCErrorCode(api.UnsupportedMethodError),
				"invalid method foo",
				nil,
			),
			ErrorCode: api.UnsupportedMethodError,
		}, resp)
		_, open := <-ch
		require.False(t, open)
	})

	t.Run("sad case stale message", func(t *testing.T) {
		invalidMsg := triggerRequest(t, nodes[0].PrivateKey, []string{"daily_price_update"}, "", "123456", "")
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleLegacyUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{
			RawResponse: codec.EncodeNewErrorResponse(
				invalidMsg.Body.MessageId,
				api.ToJSONRPCErrorCode(api.HandlerError),
				"stale message",
				nil,
			),
			ErrorCode: api.HandlerError,
		}, resp)
		_, open := <-ch
		require.False(t, open)
	})

	t.Run("sad case empty payload", func(t *testing.T) {
		invalidMsg := triggerRequest(t, nodes[0].PrivateKey, []string{"daily_price_update"}, "", "123456", "{}")
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleLegacyUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{
			RawResponse: codec.EncodeNewErrorResponse(
				invalidMsg.Body.MessageId,
				api.ToJSONRPCErrorCode(api.UserMessageParseError),
				"error decoding payload field params in TriggerRequestPayload: required",
				nil,
			),
			ErrorCode: api.UserMessageParseError,
		}, resp)
		_, open := <-ch
		require.False(t, open)
	})

	t.Run("sad case invalid payload", func(t *testing.T) {
		invalidMsg := triggerRequest(t, nodes[0].PrivateKey, []string{"daily_price_update"}, "", "123456", `{"foo":"bar"}`)
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleLegacyUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{
			RawResponse: codec.EncodeNewErrorResponse(
				invalidMsg.Body.MessageId,
				api.ToJSONRPCErrorCode(api.UserMessageParseError),
				"error decoding payload field params in TriggerRequestPayload: required",
				nil,
			),
			ErrorCode: api.UserMessageParseError,
		}, resp)
		_, open := <-ch
		require.False(t, open)
	})
	// TODO: Validate Senders and rate limit chck, pending question in trigger about where senders and rate limits are validated
}

func TestHandleComputeActionMessage(t *testing.T) {
	handler, httpClient, don, nodes := setupHandler(t)
	ctx := testutils.Context(t)
	nodeAddr := nodes[0].Address
	payload := Request{
		Method:    "GET",
		URL:       "http://example.com",
		Headers:   map[string]string{},
		Body:      nil,
		TimeoutMs: 2000,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: "123",
			Method:    MethodComputeAction,
			DonId:     "testDonId",
			Payload:   json.RawMessage(payloadBytes),
		},
	}
	err = msg.Sign(nodes[0].PrivateKey)
	require.NoError(t, err)

	t.Run("OK-compute_with_fetch", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{},
			Body:       []byte("response body"),
		}, nil).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodComputeAction == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				200 == payload.StatusCode &&
				0 == len(payload.Headers) &&
				string(payload.Body) == "response body" &&
				!payload.ExecutionError
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("NOK-payload_error_making_external_request", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 404,
			Headers:    map[string]string{},
			Body:       []byte("access denied"),
		}, nil).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodComputeAction == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				404 == payload.StatusCode &&
				string(payload.Body) == "access denied" &&
				0 == len(payload.Headers) &&
				!payload.ExecutionError
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// // ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("NOK-error_outside_payload", func(t *testing.T) {
		httpClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, errors.New("error while marshalling")).Once()

		don.EXPECT().SendToNode(mock.Anything, nodes[0].Address, mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var m api.Message
			err2 := json.Unmarshal(*req.Params, &m)
			if err2 != nil {
				return false
			}
			var payload Response
			err2 = json.Unmarshal(m.Body.Payload, &payload)
			if err2 != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodComputeAction == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				payload.ExecutionError &&
				"error while marshalling" == payload.ErrorMessage
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			// // ensure all goroutines close
			err2 := handler.Close()
			require.NoError(t, err2)
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})
}

func nodeRequest(msg *api.Message) *jsonrpc.Request[json.RawMessage] {
	req, err := hc.ValidatedRequestFromMessage(msg)
	if err != nil {
		panic(fmt.Sprintf("failed to create node request: %v", err))
	}
	return req
}
