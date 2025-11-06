package v2

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2/mocks"
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
	storageService := mocks.NewWorkflowClient(t)
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

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

	t.Run("OK-retrieve-url", func(t *testing.T) {
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		expectedURL := "some-url"
		response := storage_service.DownloadArtifactResponse{
			Url: expectedURL,
		}
		storageService.EXPECT().DownloadArtifact(mock.Anything, mock.Anything).Return(&response, nil).Once()

		req := storage_service.DownloadArtifactRequest{
			Id:   "artifact-id",
			Type: storage_service.ArtifactType_ARTIFACT_TYPE_UNSPECIFIED,
		}
		payload, err := fetcher.RetrieveURL(ctx, &req)
		require.NoError(t, err)

		require.Equal(t, expectedURL, payload)
	})

	t.Run("NOK-retrieve-url-empty-req", func(t *testing.T) {
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
		require.NoError(t, fetcher.Start(ctx))
		defer fetcher.Close()

		_, err := fetcher.RetrieveURL(ctx, nil)
		require.ErrorIs(t, err, ErrEmptyStorageRequest)
	})

	t.Run("fails with invalid payload response", func(t *testing.T) {
		connector.EXPECT().AddHandler(matches.AnyContext, []string{ghcapabilities.MethodWorkflowSyncer}, mock.Anything).Return(nil)

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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
		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

		fetcher := NewFetcherService(lggr, wrapper, storageService, gateway.WithFixedStart())
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

	t.Run("NOK-no-gateway-connector", func(t *testing.T) {
		fetcher := NewFetcherService(lggr, nil, storageService, gateway.WithFixedStart())
		require.ErrorIs(t, fetcher.Start(ctx), ErrNoGatewayConnector)
		defer fetcher.Close()
	})

	t.Run("NOK-no-storage-client", func(t *testing.T) {
		fetcher := NewFetcherService(lggr, wrapper, nil, gateway.WithFixedStart())
		require.ErrorIs(t, fetcher.Start(ctx), ErrNoStorageClient)
		defer fetcher.Close()
	})
}

func TestNewFetcherFunc(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()
	testContent := []byte("test content")

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name    string
			baseURL string
			errMsg  string
		}{
			{
				name:    "empty url",
				baseURL: "",
				errMsg:  "baseURL cannot be empty",
			},
			{
				name:    "invalid url",
				baseURL: "://invalid-url",
				errMsg:  "invalid URL",
			},
			{
				name:    "unsupported scheme",
				baseURL: "ftp://example.com",
				errMsg:  "unsupported URL scheme: ftp",
			},
			{
				name:    "relative file path",
				baseURL: "file:relative/path",
				errMsg:  "basePath must be an absolute path",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewFetcherFunc(tc.baseURL, lggr)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			})
		}
	})

	t.Run("file fetcher", func(t *testing.T) {
		// Create temp dir for test files
		tempDir := t.TempDir()
		testFilePath := filepath.Join(tempDir, "test.txt")

		// Write test content to file
		err := os.WriteFile(testFilePath, testContent, 0600)
		require.NoError(t, err)

		baseURL := "file://" + tempDir
		fetcher, err := NewFetcherFunc(baseURL, lggr)
		require.NoError(t, err)
		require.NotNil(t, fetcher)

		// Test fetching valid file
		resp, err := fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "test.txt",
		})
		require.NoError(t, err)
		assert.Equal(t, testContent, resp)

		// Test fetching non-existent file
		_, err = fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "nonexistent.txt",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")

		// Test path traversal attempt
		_, err = fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "../../../etc/passwd",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not within the basePath")

		// Test fetching full path
		resp, err = fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: testFilePath,
		})
		require.NoError(t, err)
		assert.Equal(t, testContent, resp)

		// Test full path with file:// prefix
		resp, err = fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "file://" + testFilePath,
		})
		require.NoError(t, err)
		assert.Equal(t, testContent, resp)
	})

	t.Run("http fetcher", func(t *testing.T) {
		// Create test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/test.json" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(testContent)
				assert.NoError(t, err)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		baseURL := server.URL + "/workflows"
		fetcher, err := NewFetcherFunc(baseURL, lggr)
		require.NoError(t, err)
		require.NotNil(t, fetcher)

		// Test fetching valid URL
		resp, err := fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "test.json",
		})
		require.NoError(t, err)
		assert.Equal(t, testContent, resp)

		// Test fetching non-existent resource
		_, err = fetcher(ctx, "test-msg-id", ghcapabilities.Request{
			URL: "nonexistent.json",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP request failed with status code: 404")
	})

	t.Run("context cancellation", func(t *testing.T) {
		tempDir := t.TempDir()
		baseURL := "file://" + tempDir

		fetcher, err := NewFetcherFunc(baseURL, lggr)
		require.NoError(t, err)

		// Create a canceled context
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		// Try to fetch with canceled context
		_, err = fetcher(canceledCtx, "test-msg-id", ghcapabilities.Request{
			URL: "anything.txt",
		})
		require.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("timeout handling", func(t *testing.T) {
		// Create a slow HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond) // Delay response
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(testContent)
			assert.NoError(t, err)
		}))
		defer server.Close()

		baseURL := server.URL
		fetcher, err := NewFetcherFunc(baseURL, lggr)
		require.NoError(t, err)

		// Create a context with short timeout
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Try to fetch with timeout context - should fail with deadline exceeded
		_, err = fetcher(timeoutCtx, "test-msg-id", ghcapabilities.Request{
			URL: "anything",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
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
