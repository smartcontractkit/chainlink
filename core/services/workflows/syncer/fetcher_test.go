package syncer

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

type wrapper struct {
	c connector.GatewayConnector
}

func newConnectorWrapper(c connector.GatewayConnector) *wrapper {
	return &wrapper{
		c: c,
	}
}

func (w *wrapper) GetGatewayConnector() connector.GatewayConnector {
	return w.c
}

func TestNewFetcherService(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)

	connector := gcmocks.NewGatewayConnector(t)
	wrapper := &wrapper{c: connector}

	var (
		url   = "http://example.com"
		msgID = messageID(url)
		donID = "don-id"
	)

	t.Run("OK-valid_request", func(t *testing.T) {
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 200))
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: 0,
			WorkflowID:       "foo",
		}
		payload, err := fetcher.Fetch(ctx, msgID, req)
		require.NoError(t, err)

		expectedPayload := []byte("response body")
		require.Equal(t, expectedPayload, payload)
	})

	t.Run("fails with invalid payload response", func(t *testing.T) {
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, inconsistentPayload(t, msgID, donID))
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: 0,
			WorkflowID:       "foo",
		}
		_, err := fetcher.Fetch(ctx, msgID, req)
		require.Error(t, err)
	})

	t.Run("fails due to invalid gateway response", func(t *testing.T) {
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := gatewayResponse(t, msgID, donID, 500) // gateway response that is not signed
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: 0,
			WorkflowID:       "foo",
		}
		_, err := fetcher.Fetch(ctx, msgID, req)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid response from gateway")
	})

	t.Run("NOK-response_payload_too_large", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "application/json"}
		responsePayload, err := json.Marshal(ghcapabilities.Response{
			StatusCode:   400,
			Headers:      headers,
			ErrorMessage: "http: request body too large",
		})
		require.NoError(t, err)
		gatewayResponse := &api.Message{
			Body: api.MessageBody{
				MessageId: msgID,
				Method:    ghcapabilities.MethodWebAPITarget,
				Payload:   responsePayload,
			},
		}

		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResponse)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: math.MaxUint32,
			WorkflowID:       "foo",
		}
		_, err = fetcher.Fetch(ctx, msgID, req)
		require.Error(t, err, "execution error from gateway: http: request body too large")
	})

	t.Run("NOK-bad_request", func(t *testing.T) {
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 500))
		connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: math.MaxUint32,
			WorkflowID:       "foo",
		}
		payload, err := fetcher.Fetch(ctx, msgID, req)
		require.ErrorContains(t, err, "request failed with status code")

		expectedPayload := []byte("response body")
		require.Equal(t, expectedPayload, payload)
	})

	// Connector handler never makes a connection to a gateway and the context expires.
	t.Run("NOK-request_context_deadline_exceeded", func(t *testing.T) {
		connector := gcmocks.NewGatewayConnector(t)
		wrapper := newConnectorWrapper(connector)
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(assert.AnError).Maybe()
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway2").Return(assert.AnError).Maybe()

		ctxwd, cancel := context.WithTimeout(t.Context(), 3*time.Second)
		defer cancel()
		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: math.MaxUint32,
			WorkflowID:       "foo",
		}
		_, err := fetcher.Fetch(ctxwd, url, req)
		require.Error(t, err)
		require.ErrorContains(t, err, "context deadline exceeded")
	})

	// Connector handler cycles to next available gateway after first connection fails.
	t.Run("OK-connector_handler_awaits_working_gateway", func(t *testing.T) {
		connector := gcmocks.NewGatewayConnector(t)
		wrapper := newConnectorWrapper(connector)
		connector.EXPECT().AddHandler([]string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs().Return([]string{"gateway1", "gateway2"})

		fetcher := NewFetcherService(lggr, wrapper, webapi.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().DonID().Return(donID)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(assert.AnError).Once()
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway2").Return(nil).Once()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 200))
		connector.EXPECT().SignAndSendToGateway(matches.AnyContext, "gateway2", mock.Anything).Run(func(ctx context.Context, gatewayID string, msg *api.MessageBody) {
			fetcher.och.HandleGatewayMessage(ctx, "gateway2", gatewayResp)
		}).Return(nil).Times(1)

		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: 0,
			WorkflowID:       "foo",
		}
		payload, err := fetcher.Fetch(ctx, msgID, req)
		require.NoError(t, err)

		expectedPayload := []byte("response body")
		require.Equal(t, expectedPayload, payload)
	})
}

// gatewayResponse creates an unsigned gateway response with a response body.
func gatewayResponse(t *testing.T, msgID string, donID string, statusCode int) *api.Message {
	headers := map[string]string{"Content-Type": "application/json"}
	body := []byte("response body")
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
}

// inconsistentPayload creates an unsigned gateway response with an inconsistent payload.  The
// ExecutionError is true, but there is no ErrorMessage, so it is invalid.
func inconsistentPayload(t *testing.T, msgID string, donID string) *api.Message {
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		ExecutionError: true,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			DonId:     donID,
			Method:    ghcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
}

// signGatewayResponse signs the gateway response with a private key and arbitrarily sets the receiver
// to the signer's address.  A signature and receiver are required for a valid gateway response.
func signGatewayResponse(t *testing.T, msg *api.Message) *api.Message {
	nodeKeys := common.NewTestNodes(t, 1)
	s := &signer{pk: nodeKeys[0].PrivateKey}
	signature, err := s.Sign(api.GetRawMessageBody(&msg.Body)...)
	require.NoError(t, err)
	msg.Signature = utils.StringToHex(string(signature))

	signerBytes, err := msg.ExtractSigner()
	require.NoError(t, err)

	msg.Body.Receiver = utils.StringToHex(string(signerBytes))
	return msg
}

type signer struct {
	pk *ecdsa.PrivateKey
}

func (s *signer) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(s.pk, data...)
}

func messageID(url string, parts ...string) string {
	h := sha256.New()
	h.Write([]byte(url))
	for _, p := range parts {
		h.Write([]byte(p))
	}
	hash := hex.EncodeToString(h.Sum(nil))
	p := []string{ghcapabilities.MethodWorkflowSyncer, hash}
	return strings.Join(p, "/")
}
