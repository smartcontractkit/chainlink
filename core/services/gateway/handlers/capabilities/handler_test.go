package capabilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
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
		NodeRateLimiter: nodeRateLimiterConfig,
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
	ctx := t.Context()
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
			MessageID: "123",
			Method:    MethodWorkflowSyncer,
			DonID:     "testDonId",
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
			return m.Body.MessageID == "123" &&
				MethodWorkflowSyncer == m.Body.Method &&
				m.Body.DonID == "testDonId" &&
				payload.StatusCode == http.StatusOK &&
				len(payload.Headers) == 0 &&
				string(payload.Body) == "response body" &&
				!payload.ExecutionError
		})).Return(nil).Once()
		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		handler.wg.Wait()
		require.True(t, httpClient.AssertExpectations(t))
		require.True(t, don.AssertExpectations(t))
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
			return m.Body.MessageID == "123" &&
				MethodWorkflowSyncer == m.Body.Method &&
				m.Body.DonID == "testDonId" &&
				payload.StatusCode == http.StatusNotFound &&
				string(payload.Body) == "access denied" &&
				len(payload.Headers) == 0 &&
				!payload.ExecutionError
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		handler.wg.Wait()
		require.True(t, httpClient.AssertExpectations(t))
		require.True(t, don.AssertExpectations(t))
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
			return m.Body.MessageID == "123" &&
				MethodWorkflowSyncer == m.Body.Method &&
				m.Body.DonID == "testDonId" &&
				payload.ExecutionError &&
				payload.ErrorMessage == "error while marshalling"
		})).Return(nil).Once()

		resp, err := hc.ValidatedResponseFromMessage(msg)
		require.NoError(t, err)
		err = handler.HandleNodeMessage(ctx, resp, nodeAddr)
		require.NoError(t, err)

		handler.wg.Wait()
		require.True(t, httpClient.AssertExpectations(t))
		require.True(t, don.AssertExpectations(t))
	})
}

func TestHandlerStartClose(t *testing.T) {
	handler, _, _, _ := setupHandler(t)
	ctx := t.Context()

	require.NoError(t, handler.Start(ctx))
	require.NoError(t, handler.Close())

	require.ErrorContains(t, handler.Start(ctx), "has already been started")
	require.ErrorContains(t, handler.Close(), "has already been stopped")
}
