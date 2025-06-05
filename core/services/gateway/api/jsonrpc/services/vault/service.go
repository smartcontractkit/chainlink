package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api/jsonrpc"
	gw_jsonrpc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/api/jsonrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_handlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

var (
	ErrNotAllowlisted    = errors.New("sender not allowlisted")
	ErrRateLimited       = errors.New("rate-limited")
	ErrUnsupportedMethod = errors.New("unsupported method")

	promHandlerError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_handler_error",
		Help: "Metric to track vault handler errors",
	}, []string{"don_id", "error"})

	promSecretsCreateSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_secrets_create_success",
		Help: "Metric to track successful vault secrets_create calls",
	}, []string{"don_id"})

	promSecretsCreateFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_secrets_create_failure",
		Help: "Metric to track failed vault secrets_create calls",
	}, []string{"don_id"})
)

type service struct {
	services.StateMachine

	methodConfig VaultConfig
	donConfig    *config.DONConfig
	don          gw_handlers.DON
	lggr         logger.Logger

	mu                sync.RWMutex
	userRateLimiter   *hc.RateLimiter
	nodeRateLimiter   *hc.RateLimiter
	requestTimeoutSec int

	// In-memory secret storage for demo purposes
	// In production, this would be a proper storage backend
	secretsStore map[string]map[string]SecretEntry
	storeMu      sync.RWMutex
}

type SecretEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type VaultConfig struct {
	UserRateLimiterConfig hc.RateLimiterConfig `json:"user_rate_limiter"`
	NodeRateLimiterConfig hc.RateLimiterConfig `json:"node_rate_limiter"`
	RequestTimeoutSec     int                  `json:"request_timeout_sec"`
}

var _ gw_jsonrpc.Service = (*service)(nil)

func New(methodConfig json.RawMessage, donConfig *config.DONConfig, don gw_handlers.DON, lggr logger.Logger) gw_jsonrpc.Service {
	var cfg VaultConfig
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		// Return a minimal implementation that will fail gracefully
		return &service{
			donConfig:         donConfig,
			don:               don,
			lggr:              lggr.Named("VaultMethod"),
			requestTimeoutSec: 30,
			secretsStore:      make(map[string]map[string]SecretEntry),
		}
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	userRateLimiter, _ := hc.NewRateLimiter(cfg.UserRateLimiterConfig)
	nodeRateLimiter, _ := hc.NewRateLimiter(cfg.NodeRateLimiterConfig)

	return &service{
		methodConfig:      cfg,
		donConfig:         donConfig,
		don:               don,
		lggr:              lggr.Named("VaultMethod"),
		requestTimeoutSec: cfg.RequestTimeoutSec,
		userRateLimiter:   userRateLimiter,
		nodeRateLimiter:   nodeRateLimiter,
		secretsStore:      make(map[string]map[string]SecretEntry),
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.StartOnce("VaultMethod", func() error {
		s.lggr.Info("starting vault method")
		return nil
	})
}

func (s *service) Close() error {
	return s.StopOnce("VaultMethod", func() error {
		s.lggr.Info("closing vault method")
		return nil
	})
}

func (s *service) HandleUserRequest(ctx context.Context, request *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) error {
	s.lggr.Debugw("handling vault request", "method", request.Method, "id", request.ID)

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.requestTimeoutSec)*time.Second)
	defer cancel()

	// Process request based on method
	switch request.Method {
	case MethodSecretsCreate:
		return s.handleSecretsCreate(timeoutCtx, request, callbackCh)
	default:
		s.lggr.Debugw("unsupported method", "method", request.Method)
		promHandlerError.WithLabelValues(s.donConfig.DonId, ErrUnsupportedMethod.Error()).Inc()

		response := &jsonrpc.Response{
			Version: "2.0",
			ID:      request.ID,
			Error: &jsonrpc.Error{
				Code:    -32601, // Method not found
				Message: "Method not found",
				Data:    json.RawMessage(fmt.Sprintf("Unsupported method: %s", request.Method)),
			},
		}

		select {
		case callbackCh <- response:
			return nil
		case <-timeoutCtx.Done():
			return timeoutCtx.Err()
		}
	}
}

func (s *service) handleSecretsCreate(ctx context.Context, request *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) error {
	var req SecretsCreateRequest
	if err := json.Unmarshal(request.Params, &req); err != nil {
		response := &jsonrpc.Response{
			Version: "2.0",
			ID:      request.ID,
			Error: &jsonrpc.Error{
				Code:    -32602, // Invalid params
				Message: "Invalid parameters",
				Data:    json.RawMessage(fmt.Sprintf("Failed to parse request: %v", err)),
			},
		}
		return s.sendResponse(ctx, response, callbackCh)
	}

	// Validate request
	if req.ID == "" {
		response := &jsonrpc.Response{
			Version: "2.0",
			ID:      request.ID,
			Error: &jsonrpc.Error{
				Code:    -32602, // Invalid params
				Message: "Invalid parameters",
				Data:    json.RawMessage("Secret ID cannot be empty"),
			},
		}
		return s.sendResponse(ctx, response, callbackCh)
	}

	// Store secret
	s.storeMu.Lock()
	defer s.storeMu.Unlock()

	// Extract sender from request metadata (this would come from JWT or other auth)
	senderAddr := "default" // In a real implementation, extract from authenticated context
	if s.secretsStore[senderAddr] == nil {
		s.secretsStore[senderAddr] = make(map[string]SecretEntry)
	}

	// Check if secret already exists
	if _, exists := s.secretsStore[senderAddr][req.ID]; exists {
		response := &jsonrpc.Response{
			Version: "2.0",
			ID:      request.ID,
			Error: &jsonrpc.Error{
				Code:    -32000, // Server error
				Message: "Secret already exists",
				Data:    json.RawMessage("Secret with this ID already exists"),
			},
		}
		return s.sendResponse(ctx, response, callbackCh)
	}

	// Create new secret
	secret := SecretEntry{
		ID:        req.ID,
		Value:     req.Value,
		CreatedAt: time.Now().Unix(),
	}

	s.secretsStore[senderAddr][req.ID] = secret

	// Create success response
	responseData := SecretsCreateResponse{
		ResponseBase: ResponseBase{
			Success: true,
		},
		ID: req.ID,
	}

	resultBytes, err := json.Marshal(responseData)
	if err != nil {
		promSecretsCreateFailure.WithLabelValues(s.donConfig.DonId).Inc()
		response := &jsonrpc.Response{
			Version: "2.0",
			ID:      request.ID,
			Error: &jsonrpc.Error{
				Code:    -32603, // Internal error
				Message: "Internal error",
				Data:    json.RawMessage(fmt.Sprintf("Failed to marshal response: %v", err)),
			},
		}
		return s.sendResponse(ctx, response, callbackCh)
	}

	response := &jsonrpc.Response{
		Version: "2.0",
		ID:      request.ID,
		Result:  resultBytes,
	}

	promSecretsCreateSuccess.WithLabelValues(s.donConfig.DonId).Inc()
	return s.sendResponse(ctx, response, callbackCh)
}

func (s *service) sendResponse(ctx context.Context, response *jsonrpc.Response, callbackCh chan<- *jsonrpc.Response) error {
	select {
	case callbackCh <- response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
