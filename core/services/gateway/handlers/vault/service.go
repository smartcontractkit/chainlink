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

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_handlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
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

var _ gw_handlers.Handler = (*service)(nil)

type service struct {
	services.StateMachine
	gw_handlers.Handler
	methodConfig Config
	donConfig    *config.DONConfig
	don          gw_handlers.DON
	lggr         logger.Logger
	codec        api.JsonRPCCodec

	userRateLimiter   *ratelimit.RateLimiter
	nodeRateLimiter   *ratelimit.RateLimiter
	requestTimeoutSec int

	// In-memory secret storage for demo purposes
	// In production, this would be a proper storage backend
	secretsStore map[string]map[string]SecretEntry
	storeMu      sync.RWMutex
}

func (s *service) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *service) Name() string {
	return s.lggr.Name()
}

type SecretEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type Config struct {
	UserRateLimiterConfig ratelimit.RateLimiterConfig `json:"user_rate_limiter"`
	NodeRateLimiterConfig ratelimit.RateLimiterConfig `json:"node_rate_limiter"`
	RequestTimeoutSec     int                         `json:"request_timeout_sec"`
}

func NewService(methodConfig json.RawMessage, donConfig *config.DONConfig, don gw_handlers.DON, lggr logger.Logger) *service {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		// Return a minimal implementation that will fail gracefully
		return &service{
			donConfig:         donConfig,
			don:               don,
			lggr:              logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
			codec:             api.JsonRPCCodec{},
			requestTimeoutSec: 30,
			secretsStore:      make(map[string]map[string]SecretEntry),
		}
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	userRateLimiter, _ := ratelimit.NewRateLimiter(cfg.UserRateLimiterConfig)
	nodeRateLimiter, _ := ratelimit.NewRateLimiter(cfg.NodeRateLimiterConfig)

	return &service{
		methodConfig:      cfg,
		donConfig:         donConfig,
		don:               don,
		lggr:              logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
		requestTimeoutSec: cfg.RequestTimeoutSec,
		userRateLimiter:   userRateLimiter,
		nodeRateLimiter:   nodeRateLimiter,
		secretsStore:      make(map[string]map[string]SecretEntry),
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.StartOnce("VaultService", func() error {
		s.lggr.Info("starting vault service")
		return nil
	})
}

func (s *service) Close() error {
	return s.StopOnce("VaultMethod", func() error {
		s.lggr.Info("closing vault service")
		return nil
	})
}

func (s *service) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	return errors.New("vault service does not support legacy messages")
}

func (s *service) HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	s.lggr.Debugw("handling vault request", "method", jsonRequest.Method, "id", jsonRequest.ID)

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.requestTimeoutSec)*time.Second)
	defer cancel()

	// Process request based on method
	switch jsonRequest.Method {
	case MethodSecretsCreate:
		return s.handleSecretsCreate(timeoutCtx, jsonRequest, callbackCh)
	default:
		return s.handleUnsupportedMethod(timeoutCtx, jsonRequest, callbackCh)
	}
}

func (s *service) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	return errors.New("node message support not implemented")
}

func (s *service) handleSecretsCreate(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	var req SecretsCreateRequest
	if err := json.Unmarshal(*jsonRequest.Params, &req); err != nil {
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.InvalidParamsError),
				fmt.Sprintf("Failed to parse request: %v", err),
				nil,
			),
			ErrorCode: api.InvalidParamsError,
		}, callbackCh)
	}

	// Validate request
	if req.ID == "" {
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.InvalidParamsError),
				"Secret ID cannot be empty",
				nil,
			),
			ErrorCode: api.InvalidParamsError,
		}, callbackCh)
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
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.InvalidParamsError),
				"Secret with this ID already exists",
				nil,
			),
			ErrorCode: api.InvalidParamsError,
		}, callbackCh)
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

	var resultBytes json.RawMessage
	resultBytes, err := json.Marshal(responseData)
	if err != nil {
		promSecretsCreateFailure.WithLabelValues(s.donConfig.DonId).Inc()
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.NodeReponseEncodingError),
				fmt.Sprintf("Failed to marshal response: %v", err),
				nil,
			),
			ErrorCode: api.NodeReponseEncodingError,
		}, callbackCh)
	}
	jsonResponse := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      jsonRequest.ID,
		Result:  &resultBytes,
	}
	rawResponse, err := json.Marshal(jsonResponse)
	if err != nil {
		promSecretsCreateFailure.WithLabelValues(s.donConfig.DonId).Inc()
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.NodeReponseEncodingError),
				fmt.Sprintf("Failed to marshal response: %v", err),
				nil,
			),
			ErrorCode: api.NodeReponseEncodingError,
		}, callbackCh)
	}
	responseObj := gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.NoError,
	}

	promSecretsCreateSuccess.WithLabelValues(s.donConfig.DonId).Inc()
	return s.sendResponse(ctx, responseObj, callbackCh)
}

func (s *service) handleUnsupportedMethod(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	s.lggr.Debugw("unsupported method", "method", jsonRequest.Method)
	promHandlerError.WithLabelValues(s.donConfig.DonId, ErrUnsupportedMethod.Error()).Inc()

	return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
		RawResponse: s.codec.EncodeNewErrorResponse(
			jsonRequest.ID,
			api.ToJSONRPCErrorCode(api.UnsupportedMethodError),
			"Unsupported method: "+jsonRequest.Method,
			nil,
		),
		ErrorCode: api.UnsupportedMethodError,
	}, callbackCh)
}

func (s *service) sendResponse(ctx context.Context, response gw_handlers.UserCallbackPayload, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	select {
	case callbackCh <- response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
