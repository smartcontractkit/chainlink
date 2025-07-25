package webapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

type ctxKey string

const (
	privateKey = "6c358b4f16344f03cfce12ebf7b768301bbe6a8977c98a2a2d76699f8bc56161"
)

func TestOutgoingConnectorHandler_AwaitConnection(t *testing.T) {
	gateways := []string{"gateway1", "gateway2"}

	type testCase struct {
		name string

		gatewayConnectorSetup func(*gcmocks.GatewayConnector)
		ctxSetup              func() context.Context
		expectedGateway       string
		expectedError         string
	}

	testCases := []testCase{
		{
			name: "successful connection on first try",
			gatewayConnectorSetup: func(mockConnector *gcmocks.GatewayConnector) {
				mockConnector.EXPECT().AwaitConnection(mock.Anything, "gateway1").Return(nil).Once()
				mockConnector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
			},
			ctxSetup:        context.Background,
			expectedGateway: "gateway1",
		},
		{
			name: "connection timeout then success",
			gatewayConnectorSetup: func(mockConnector *gcmocks.GatewayConnector) {
				mockConnector.EXPECT().AwaitConnection(mock.Anything, "gateway1").Return(errors.New("timeout")).Once()
				mockConnector.EXPECT().AwaitConnection(mock.Anything, "gateway2").Return(nil).Once()
				mockConnector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
			},
			ctxSetup:        context.Background,
			expectedGateway: "gateway2",
		},
		{
			name: "connection timeout then success after backoff",
			gatewayConnectorSetup: func(mockConnector *gcmocks.GatewayConnector) {
				mockConnector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
				mockConnector.EXPECT().AwaitConnection(mock.Anything, "gateway1").Return(errors.New("gateway connection failed: timeout")).Once()
				mockConnector.EXPECT().AwaitConnection(mock.Anything, "gateway2").Return(errors.New("gateway connection failed: timeout")).Once()

				// second call to gateway1 succeeds
				mockConnector.On("AwaitConnection", mock.Anything, "gateway1").Return(nil).Once()
			},
			ctxSetup:        context.Background,
			expectedGateway: "gateway1",
		},
		{
			name: "all gateways fail and context canceled",
			gatewayConnectorSetup: func(mockConnector *gcmocks.GatewayConnector) {
				callCount := 0
				mockConnector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
				mockConnector.EXPECT().AwaitConnection(mock.Anything, mock.Anything).Return(errors.New("gateway connection failed: timeout")).Run(func(ctx context.Context, gatewayID string) {
					callCount++
					if callCount == len(gateways) {
						cancelFunc := ctx.Value(ctxKey("cancelFunc")).(context.CancelFunc)
						cancelFunc()
					}
				})
			},
			ctxSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, ctxKey("cancelFunc"), cancel)
				return ctx
			},
			expectedGateway: "",
			expectedError:   "context canceled",
		},
		{
			name: "context canceled",
			gatewayConnectorSetup: func(mockConnector *gcmocks.GatewayConnector) {
				mockConnector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
			},
			ctxSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel the context immediately
				return ctx
			},
			expectedGateway: "",
			expectedError:   "context canceled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConnector := gcmocks.NewGatewayConnector(t)
			lggr := logger.Test(t)

			if tc.gatewayConnectorSetup != nil {
				tc.gatewayConnectorSetup(mockConnector)
			}

			c := &OutgoingConnectorHandler{
				gc:           mockConnector,
				lggr:         lggr,
				selectorOpts: []func(*gateway.RoundRobinSelector){gateway.WithFixedStart()},
			}

			ctx := tc.ctxSetup()
			gateway, err := c.awaitConnection(ctx, awaitContext{
				workflowID: "test-workflow",
				messageID:  "some-message",
			})

			assert.Equal(t, tc.expectedGateway, gateway)
			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockConnector.AssertExpectations(t)
		})
	}
}

func TestHandleSingleNodeRequest(t *testing.T) {
	t.Run("uses default timeout if no timeout is provided", func(t *testing.T) {
		msgID := "msgID"
		testURL := "http://localhost:8080"
		connector, connectorHandler := newFunctionWithDefaultConfig(
			t,
			func(gc *gcmocks.GatewayConnector) {
				gc.EXPECT().DonID(matches.AnyContext).Return("donID", nil)
				gc.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
				gc.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1"}, nil)
			},
		)

		// build the expected body with the default timeout
		req := ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: defaultFetchTimeoutMs,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)
		donID, err := connector.DonID(t.Context())
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the default timeout
		connector.EXPECT().SignMessage(mock.Anything, common.Flatten(api.GetRawMessageBody(expectedBody)...)).Return([]byte("signature"), nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResponse(t, msgID, privateKey))
			require.NoError(t, err2)
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL: testURL,
		})
		require.NoError(t, err)
	})

	t.Run("uses timeout", func(t *testing.T) {
		msgID := "msgID"
		testURL := "http://localhost:8080"
		connector, connectorHandler := newFunctionWithDefaultConfig(
			t,
			func(gc *gcmocks.GatewayConnector) {
				gc.EXPECT().DonID(matches.AnyContext).Return("donID", nil)
				gc.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
				gc.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1"}, nil)
			},
		)

		// build the expected body with the defined timeout
		req := ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 40000,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)
		donID, err := connector.DonID(t.Context())
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the defined timeout
		connector.EXPECT().SignMessage(mock.Anything, common.Flatten(api.GetRawMessageBody(expectedBody)...)).Return([]byte("signature"), nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResponse(t, msgID, privateKey))
			require.NoError(t, err2)
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 40000,
		})
		_, found := connectorHandler.responses.get(msgID)
		assert.False(t, found)
		require.NoError(t, err)
	})

	t.Run("cleans up in event of a timeout", func(t *testing.T) {
		msgID := "msgID"
		testURL := "http://localhost:8080"
		connector, connectorHandler := newFunctionWithDefaultConfig(
			t,
			func(gc *gcmocks.GatewayConnector) {
				gc.EXPECT().DonID(matches.AnyContext).Return("donID", nil)
				gc.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
				gc.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1"}, nil)
			},
		)

		// build the expected body with the defined timeout
		req := ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 10,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)
		donID, err := connector.DonID(t.Context())
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the defined timeout
		connector.EXPECT().SignMessage(mock.Anything, common.Flatten(api.GetRawMessageBody(expectedBody)...)).Return([]byte("signature"), nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			// don't call HandleGatewayMessage here; i.e. simulate a failure to receive a response
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL:       testURL,
			TimeoutMs: 10,
		})
		_, found := connectorHandler.responses.get(msgID)
		assert.False(t, found)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("rate limits outgoing traffic", func(t *testing.T) {
		msgID := "msgID"
		testURL := "http://localhost:8080"
		var config = ServiceConfig{
			OutgoingRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      2.0,
				GlobalBurst:    2,
				PerSenderRPS:   1.0,
				PerSenderBurst: 1,
			},
			RateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}
		connector, connectorHandler := newFunction(
			t,
			func(gc *gcmocks.GatewayConnector) {
				gc.EXPECT().DonID(matches.AnyContext).Return("donID", nil)
				gc.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
				gc.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1"}, nil)
			},
			config,
		)

		// build the expected body with the default timeout
		req := ghcapabilities.Request{
			URL:        testURL,
			WorkflowID: "1",
			TimeoutMs:  defaultFetchTimeoutMs,
		}
		payload, err := json.Marshal(req)
		require.NoError(t, err)
		donID, err := connector.DonID(t.Context())
		require.NoError(t, err)

		expectedBody := &api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodComputeAction,
			Payload:   payload,
		}

		// expect the request body to contain the default timeout
		connector.EXPECT().SignMessage(mock.Anything, common.Flatten(api.GetRawMessageBody(expectedBody)...)).Return([]byte("signature"), nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResponse(t, msgID, privateKey))
			require.NoError(t, err2)
		}).Return(nil).Times(1)

		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL:        testURL,
			WorkflowID: "1",
		})
		require.NoError(t, err)

		// Second request should error from workflow ratelimit
		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL:        testURL,
			WorkflowID: "1",
		})
		require.Error(t, err)
		require.ErrorContains(t, err, errorOutgoingRatelimitWorkflow)

		// Third request should error from global ratelimit
		_, err = connectorHandler.HandleSingleNodeRequest(t.Context(), msgID, ghcapabilities.Request{
			URL:        testURL,
			WorkflowID: "2",
		})
		require.Error(t, err)
		require.ErrorContains(t, err, errorOutgoingRatelimitGlobal)
	})
}

func newFunctionWithDefaultConfig(t *testing.T, mockFn func(*gcmocks.GatewayConnector)) (*gcmocks.GatewayConnector, *OutgoingConnectorHandler) {
	var defaultConfig = ServiceConfig{
		OutgoingRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
		RateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	}
	return newFunction(t, mockFn, defaultConfig)
}

func newFunction(t *testing.T, mockFn func(*gcmocks.GatewayConnector), serviceConfig ServiceConfig) (*gcmocks.GatewayConnector, *OutgoingConnectorHandler) {
	log := logger.Test(t)
	connector := gcmocks.NewGatewayConnector(t)

	mockFn(connector)

	connectorHandler, err := NewOutgoingConnectorHandler(connector, serviceConfig, ghcapabilities.MethodComputeAction, log, gateway.WithFixedStart())
	require.NoError(t, err)
	return connector, connectorHandler
}

func gatewayResponse(t *testing.T, msgID string, privateKey string) *jsonrpc.Request[json.RawMessage] {
	headers := map[string]string{"Content-Type": "application/json"}
	body := []byte("response body")
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		StatusCode:     200,
		Headers:        headers,
		Body:           body,
		ExecutionError: false,
	})
	require.NoError(t, err)
	m := &api.Message{
		Body: api.MessageBody{
			DonId:     "donID",
			MessageId: msgID,
			Method:    ghcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)
	err = m.Sign(key)
	require.NoError(t, err)
	req, err := hc.ValidatedRequestFromMessage(m)
	require.NoError(t, err)
	return req
}

func TestServiceConfigDefaults(t *testing.T) {
	t.Run("fills default RateLimiterConfigs", func(t *testing.T) {
		var cfg ServiceConfig

		tomlErr := toml.Unmarshal([]byte{}, &cfg)
		require.NoError(t, tomlErr)

		iRLConf := incomingRateLimiterConfigDefaults(cfg.RateLimiter)
		require.Equal(t, DefaultGlobalBurst, iRLConf.GlobalBurst)
		require.InDelta(t, DefaultGlobalRPS, iRLConf.GlobalRPS, 0.001)
		require.Equal(t, DefaultPerSenderBurst, iRLConf.PerSenderBurst)
		require.InDelta(t, DefaultPerSenderRPS, iRLConf.PerSenderRPS, 0.001)

		oRLConf := outgoingRateLimiterConfigDefaults(cfg.RateLimiter)
		require.Equal(t, DefaultGlobalBurst, oRLConf.GlobalBurst)
		require.InDelta(t, DefaultGlobalRPS, oRLConf.GlobalRPS, 0.001)
		require.Equal(t, DefaultWorkflowBurst, oRLConf.PerSenderBurst)
		require.InDelta(t, DefaultWorkflowRPS, oRLConf.PerSenderRPS, 0.001)
	})
}
func TestOutgoingConnectorHandler_HandleGatewayMessage_InvalidMessage(t *testing.T) {
	_, handler := newFunctionWithDefaultConfig(
		t,
		func(gc *gcmocks.GatewayConnector) {},
	)
	invalidMsg := api.Message{
		Body: api.MessageBody{
			// MessageId is empty, which should fail Validate()
			Method: "some-method",
		},
	}
	params, err := json.Marshal(invalidMsg)
	require.NoError(t, err)
	rawParams := json.RawMessage(params)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		Method:  invalidMsg.Body.Method,
		Params:  &rawParams,
	}
	err = handler.HandleGatewayMessage(context.Background(), "gateway1", req)
	require.NoError(t, err)
}
