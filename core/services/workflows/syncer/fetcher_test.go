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

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
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
	signature := []byte("signature")

	var (
		url   = "http://example.com"
		msgID = messageID(url)
		donID = "don-id"
	)

	t.Run("OK-valid_request", func(t *testing.T) {
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 200))
		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
			require.NoError(t, err2)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
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
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, inconsistentPayload(t, msgID, donID))
		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
			require.NoError(t, err2)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

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
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayMessage := gatewayResponse(t, msgID, donID, 500) // gateway response that is not signed
		payload, err := json.Marshal(gatewayMessage)
		require.NoError(t, err)
		rawPayload := json.RawMessage(payload)
		gatewayResp := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      gatewayMessage.Body.MessageId,
			Method:  gatewayMessage.Body.Method,
			Params:  &rawPayload,
		}
		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
			require.NoError(t, err2)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)
		ctxwd, cancel := context.WithTimeout(t.Context(), 3*time.Second)
		defer cancel()
		req := ghcapabilities.Request{
			URL:              url,
			Method:           http.MethodGet,
			MaxResponseBytes: 0,
			WorkflowID:       "foo",
		}
		_, err = fetcher.Fetch(ctxwd, msgID, req)
		require.Error(t, err)
		require.ErrorContains(t, err, "context deadline exceeded")
	})

	t.Run("NOK-response_payload_too_large", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "application/json"}
		responsePayload, err := json.Marshal(ghcapabilities.Response{
			StatusCode:   400,
			Headers:      headers,
			ErrorMessage: "http: request body too large",
		})
		require.NoError(t, err)
		gatewayMsg := &api.Message{
			Body: api.MessageBody{
				MessageId: msgID,
				Method:    ghcapabilities.MethodWebAPITarget,
				Payload:   responsePayload,
			},
		}
		payload, err := json.Marshal(gatewayMsg)
		require.NoError(t, err)
		rawPayload := json.RawMessage(payload)
		gatewayResp := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      gatewayMsg.Body.MessageId,
			Method:  gatewayMsg.Body.Method,
			Params:  &rawPayload,
		}
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
			require.NoError(t, err2)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

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
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 500))
		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(mock.Anything, "gateway1", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
			require.NoError(t, err2)
		}).Return(nil).Times(1)
		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
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
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
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
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)
		connector.EXPECT().GatewayIDs(matches.AnyContext).Return([]string{"gateway1", "gateway2"}, nil)

		fetcher := NewFetcherService(lggr, wrapper, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		connector.EXPECT().DonID(matches.AnyContext).Return(donID, nil)
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway1").Return(assert.AnError).Once()
		connector.EXPECT().AwaitConnection(matches.AnyContext, "gateway2").Return(nil).Once()

		gatewayResp := signGatewayResponse(t, gatewayResponse(t, msgID, donID, 200))
		connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return(signature, nil).Once()
		connector.EXPECT().SendToGateway(matches.AnyContext, "gateway2", mock.Anything).Run(func(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) {
			err2 := fetcher.och.HandleGatewayMessage(ctx, "gateway2", gatewayResp)
			require.NoError(t, err2)
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
func signGatewayResponse(t *testing.T, msg *api.Message) *jsonrpc.Request[json.RawMessage] {
	nodeKeys := common.NewTestNodes(t, 1)
	s := &signer{pk: nodeKeys[0].PrivateKey}
	msgToSign := api.GetRawMessageBody(&msg.Body)
	signature, err := s.Sign(msgToSign...)
	require.NoError(t, err)
	msg.Signature = utils.StringToHex(string(signature))

	signerBytes, err := msg.ExtractSigner()
	require.NoError(t, err)

	msg.Body.Receiver = utils.StringToHex(string(signerBytes))
	require.NoError(t, err)
	resp, err := hc.ValidatedRequestFromMessage(msg)
	require.NoError(t, err)
	return resp
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
