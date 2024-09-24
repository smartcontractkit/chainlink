package webapicapabilities

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

func setupHandler(t *testing.T) (*handler, *mocks.HTTPClient, *handlermocks.DON, []gwcommon.TestNode) {
	lggr := logger.TestLogger(t)
	httpClient := mocks.NewHTTPClient(t)
	don := handlermocks.NewDON(t)
	nodeRateLimiterConfig := common.RateLimiterConfig{
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
	ctx := testutils.Context(t)
	nodeAddr := nodes[0].Address
	payload := TargetRequestPayload{
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

	t.Run("happy case", func(t *testing.T) {
		httpClient.EXPECT().Send(ctx, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{},
			Body:       []byte("response body"),
		}, nil).Once()

		don.EXPECT().SendToNode(ctx, nodes[0].Address, mock.MatchedBy(func(m *api.Message) bool {
			var payload TargetResponsePayload
			err = json.Unmarshal(m.Body.Payload, &payload)
			if err != nil {
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

		err = handler.HandleNodeMessage(ctx, msg, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("http client non-HTTP error", func(t *testing.T) {
		httpClient.EXPECT().Send(ctx, mock.Anything).Return(&network.HTTPResponse{
			StatusCode: 404,
			Headers:    map[string]string{},
			Body:       []byte("access denied"),
		}, nil).Once()

		don.EXPECT().SendToNode(ctx, nodes[0].Address, mock.MatchedBy(func(m *api.Message) bool {
			var payload TargetResponsePayload
			err = json.Unmarshal(m.Body.Payload, &payload)
			if err != nil {
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

		err = handler.HandleNodeMessage(ctx, msg, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})

	t.Run("http client non-HTTP error", func(t *testing.T) {
		httpClient.EXPECT().Send(ctx, mock.Anything).Return(nil, fmt.Errorf("error while marshalling")).Once()

		don.EXPECT().SendToNode(ctx, nodes[0].Address, mock.MatchedBy(func(m *api.Message) bool {
			var payload TargetResponsePayload
			err = json.Unmarshal(m.Body.Payload, &payload)
			if err != nil {
				return false
			}
			return "123" == m.Body.MessageId &&
				MethodWebAPITarget == m.Body.Method &&
				"testDonId" == m.Body.DonId &&
				payload.ExecutionError &&
				"error while marshalling" == payload.ErrorMessage
		})).Return(nil).Once()

		err = handler.HandleNodeMessage(ctx, msg, nodeAddr)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return httpClient.AssertExpectations(t) && don.AssertExpectations(t)
		}, tests.WaitTimeout(t), 100*time.Millisecond)
	})
}
