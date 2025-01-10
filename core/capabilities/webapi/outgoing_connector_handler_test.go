package webapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

func TestHandleSingleNodeRequest(t *testing.T) {
	t.Run("OK-timeout_is_not_specify_default_timeout_is_expected", func(t *testing.T) {
		ctx := tests.Context(t)
		log := logger.TestLogger(t)
		connector := gcmocks.NewGatewayConnector(t)
		var defaultConfig = ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}
		connectorHandler, err := NewOutgoingConnectorHandler(connector, defaultConfig, ghcapabilities.MethodComputeAction, log)
		require.NoError(t, err)

		msgID := "msgID"
		testURL := "http://localhost:8080"
		connector.EXPECT().DonID().Return("donID")
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})

		// build the expected body with the default timeout
		req := ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: defaultFetchTimeoutMs,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     connector.DonID(),
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the default timeout
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", expectedBody).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResponse(t, msgID))
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(ctx, msgID, ghcapabilities.Request{
			URL: testURL,
		})
		require.NoError(t, err)
	})

	t.Run("OK-timeout_is_specified", func(t *testing.T) {
		ctx := tests.Context(t)
		log := logger.TestLogger(t)
		connector := gcmocks.NewGatewayConnector(t)
		var defaultConfig = ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}
		connectorHandler, err := NewOutgoingConnectorHandler(connector, defaultConfig, ghcapabilities.MethodComputeAction, log)
		require.NoError(t, err)

		msgID := "msgID"
		testURL := "http://localhost:8080"
		connector.EXPECT().DonID().Return("donID")
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})

		// build the expected body with the defined timeout
		req := ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 40000,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     connector.DonID(),
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the defined timeout
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", expectedBody).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResponse(t, msgID))
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(ctx, msgID, ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 40000,
		})
		require.NoError(t, err)
	})
}

func gatewayResponse(t *testing.T, msgID string) *api.Message {
	headers := map[string]string{"Content-Type": "application/json"}
	body := []byte("response body")
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		StatusCode:     200,
		Headers:        headers,
		Body:           body,
		ExecutionError: false,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			Method:    ghcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
}
