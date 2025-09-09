package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/storage"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var (
	ErrNoGatewayConnector  = errors.New("failed to start fetcher service: gateway connector is not configured")
	ErrNoStorageClient     = errors.New("failed to start fetcher service: storage client is not configured")
	ErrEmptyStorageRequest = errors.New("storage service request must not be empty")
)

// FetcherService is a service that fetches data from the gateway using the OutgoingConnectorHandler.
type FetcherService struct {
	services.StateMachine
	lggr          logger.Logger
	och           *webapi.OutgoingConnectorHandler
	wrapper       gatewayConnector
	selectorOpts  []func(*gateway.RoundRobinSelector)
	storageClient storage.WorkflowClient
}

type WorkflowClient interface {
	DownloadArtifact(ctx context.Context, req *storage_service.DownloadArtifactRequest) (*storage_service.DownloadArtifactResponse, error)
	Close() error
}

type gatewayConnector interface {
	GetGatewayConnector() connector.GatewayConnector
}

func NewFetcherService(lggr logger.Logger, wrapper gatewayConnector, storageClient storage.WorkflowClient, selectorOpts ...func(*gateway.RoundRobinSelector)) *FetcherService {
	return &FetcherService{
		lggr:          lggr.Named("FetcherService"),
		wrapper:       wrapper,
		storageClient: storageClient,
		selectorOpts:  selectorOpts,
	}
}

func (s *FetcherService) Start(ctx context.Context) error {
	return s.StartOnce("FetcherService", func() error {
		if s.wrapper == nil {
			return ErrNoGatewayConnector
		}
		if s.storageClient == nil {
			return ErrNoStorageClient
		}

		connector := s.wrapper.GetGatewayConnector()

		outgoingConnectorLggr := s.lggr.Named("OutgoingConnectorHandler")

		webAPIConfig := webapi.ServiceConfig{
			OutgoingRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      webapi.DefaultGlobalRPS,
				GlobalBurst:    webapi.DefaultGlobalBurst,
				PerSenderRPS:   webapi.DefaultWorkflowRPS,
				PerSenderBurst: webapi.DefaultWorkflowBurst,
			},
			RateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}

		och, err := webapi.NewOutgoingConnectorHandler(connector,
			webAPIConfig,
			ghcapabilities.MethodWorkflowSyncer, outgoingConnectorLggr, s.selectorOpts...)
		if err != nil {
			return fmt.Errorf("could not create outgoing connector handler: %w", err)
		}

		s.och = och
		return och.Start(ctx)
	})
}

func (s *FetcherService) Close() error {
	return s.StopOnce("FetcherService", func() error {
		return s.och.Close()
	})
}

func (s *FetcherService) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *FetcherService) Name() string {
	return s.lggr.Name()
}

// RetrieveURL gets an ephemeral endpoint to download the given artifact from the storage service.
func (s *FetcherService) RetrieveURL(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error) {
	if req == nil {
		return "", ErrEmptyStorageRequest
	}

	storageResp, err := s.storageClient.DownloadArtifact(ctx, req)
	if err != nil {
		return "", err
	}

	s.lggr.Debugw("received response from storage service", "ID", req.Id, "Type", req.Type, "Expiry", storageResp.Expiry)

	return storageResp.GetUrl(), nil
}

// Fetch fetches the given URL and returns the response body.  n is the maximum number of bytes to
// read from the response body.  Set n to zero to use the default size limit specified by the
// configured gateway's http client, if any.
func (s *FetcherService) Fetch(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error) {
	if req.WorkflowID == "" {
		return nil, errors.New("invalid call to fetch, must provide workflow ID")
	}

	resp, err := s.och.HandleSingleNodeRequest(ctx, messageID, req)
	if err != nil {
		return nil, err
	}

	if err = resp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid response from gateway: %w", err)
	}

	s.lggr.Debugw("received gateway response", "donID", resp.Body.DonId, "msgID", resp.Body.MessageId, "receiver", resp.Body.Receiver, "sender", resp.Body.Sender)

	var payload ghcapabilities.Response
	if err = json.Unmarshal(resp.Body.Payload, &payload); err != nil {
		return nil, err
	}

	if err = payload.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload received from gateway message: %w", err)
	}

	if payload.ExecutionError {
		return nil, fmt.Errorf("execution error from gateway: %s", payload.ErrorMessage)
	}

	if payload.StatusCode < 200 || payload.StatusCode >= 300 {
		// NOTE: redirects are currently not supported
		return payload.Body, fmt.Errorf("request failed with status code: %d", payload.StatusCode)
	}

	return payload.Body, nil
}

// NewFetcher creates a new FetcherFunc based on the provided URL configuration
// The implementation supports both file and HTTP(S) URLs and bypasses the gateway
func NewFetcherFunc(baseURL string, lggr logger.Logger) (types.FetcherFunc, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL cannot be empty")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	switch u.Scheme {
	case "file":
		// Ensure the basePath is absolute
		if !filepath.IsAbs(u.Path) {
			return nil, fmt.Errorf("basePath must be an absolute path, got: %s", u.Path)
		}
		return newFileFetcher(u.Path, lggr), nil
	case "http", "https":
		return newHTTPFetcher(baseURL, lggr), nil
	default:
		return nil, fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
	}
}

func newFileFetcher(basePath string, lggr logger.Logger) types.FetcherFunc {
	return func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// the incoming request URL is expected to be a relative path or a path within the basePath
		if req.URL == "" {
			return nil, errors.New("request URL cannot be empty")
		}
		u, err := url.Parse(req.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		fullPath := filepath.Clean(u.Path)

		// ensure that the incoming request URL is either relative or absolute but within the basePath
		if !filepath.IsAbs(fullPath) {
			// If it's not absolute, we assume it's relative to the basePath
			fullPath = filepath.Join(basePath, fullPath)
		}
		if !strings.HasPrefix(fullPath, basePath) {
			return nil, fmt.Errorf("request URL %s is not within the basePath %s", fullPath, basePath)
		}

		lggr.Debugw("Fetching file", "messageID", messageID, "path", fullPath)

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return data, nil
	}
}

func newHTTPFetcher(baseURL string, lggr logger.Logger) types.FetcherFunc {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error) {
		// Clean the path to prevent directory traversal
		cleanPath := strings.TrimPrefix(filepath.Clean(req.URL), "/")

		// Join base URL with path
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}

		u.Path = filepath.Join(u.Path, cleanPath)
		fetchURL := u.String()

		lggr.Debugw("Fetching HTTP resource", "url", fetchURL)

		req2, err := http.NewRequestWithContext(ctx, http.MethodGet, fetchURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req2)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return data, nil
	}
}
