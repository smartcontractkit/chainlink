package syncer

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

func TestNewFetcherFunc(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)

	config := webapi.ServiceConfig{
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	}

	connector := gcmocks.NewGatewayConnector(t)
	och, err := webapi.NewOutgoingConnectorHandler(connector, config, ghcapabilities.MethodComputeAction, lggr)
	require.NoError(t, err)

	workflowID := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID := "25c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1"
	idGenerator := func() string { return "uniqueID" }
	url := "http://example.com"

	msgID := strings.Join([]string{
		workflowID,
		workflowExecutionID,
		ghcapabilities.MethodWorkflowSyncer,
		idGenerator(),
	}, "/")

	t.Run("OK-valid_request", func(t *testing.T) {
		gatewayResp := gatewayResponse(t, msgID)
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return("don-id")
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		fetcher := NewFetcherFunc(ctx, lggr, och, workflowID, workflowExecutionID, idGenerator)

		payload, err := fetcher(ctx, url)
		require.NoError(t, err)

		expectedPayload := []byte("response body")
		require.Equal(t, expectedPayload, payload)
	})

	t.Run("NOK-invalid_workflow_id", func(t *testing.T) {
		invalidWorkflowID := "invalidWorkflowID"
		fetcher := NewFetcherFunc(ctx, lggr, och, invalidWorkflowID, workflowExecutionID, idGenerator)

		_, err := fetcher(ctx, url)
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflow ID")
	})

	t.Run("NOK-invalid_workflow_execution_id", func(t *testing.T) {
		invalidWorkflowExecutionID := "invalidWorkflowExecutionID"
		fetcher := NewFetcherFunc(ctx, lggr, och, workflowID, invalidWorkflowExecutionID, idGenerator)

		_, err := fetcher(ctx, url)
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflow execution ID")
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
