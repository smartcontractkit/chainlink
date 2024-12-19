package syncer

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

type wrapper struct {
	c connector.GatewayConnector
}

func (w *wrapper) GetGatewayConnector() connector.GatewayConnector {
	return w.c
}

func TestNewFetcherService(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)

	connector := gcmocks.NewGatewayConnector(t)
	wrapper := &wrapper{c: connector}

	url := "http://example.com"

	msgID := strings.Join([]string{ghcapabilities.MethodWorkflowSyncer, hash(url)}, "/")

	t.Run("OK-valid_request", func(t *testing.T) {
		connector.EXPECT().AddHandler([]string{capabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper)
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := gatewayResponse(t, msgID)
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return("don-id")
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		payload, err := fetcher.Fetch(ctx, url)
		require.NoError(t, err)

		expectedPayload := []byte("response body")
		require.Equal(t, expectedPayload, payload)
	})

	t.Run("NOK-response_payload_too_large", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "application/json"}
		responsePayload, err := json.Marshal(ghcapabilities.Response{
			StatusCode:     400,
			Headers:        headers,
			ErrorMessage:   "http: request body too large",
			ExecutionError: true,
		})
		require.NoError(t, err)
		gatewayResponse := &api.Message{
			Body: api.MessageBody{
				MessageId: msgID,
				Method:    ghcapabilities.MethodWebAPITarget,
				Payload:   responsePayload,
			},
		}

		connector.EXPECT().AddHandler([]string{capabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper)
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResponse)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return("don-id")
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		_, err = fetcher.Fetch(ctx, url)
		require.Error(t, err, "execution error from gateway: http: request body too large")
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
