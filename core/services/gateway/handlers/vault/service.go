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

	// Track errors from the handler (internal errors)
	promRequestInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_request_internal_error",
		Help: "Metric to track vault handler errors",
	}, []string{"don_id", "error"})

	// Track failed secrets_create calls (user errors)
	promRequestUserError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_request_user_error",
		Help: "Metric to track failed vault requests",
	}, []string{"don_id"})

	// Track successful secrets_create calls
	promRequestSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_request_success",
		Help: "Metric to track successful vault requests",
	}, []string{"don_id"})
)

var _ gw_handlers.Handler = (*service)(nil)

type pendingRequest struct {
	callbackCh chan<- gw_handlers.UserCallbackPayload
	responses  map[string]*jsonrpc.Response[json.RawMessage]
}

type service struct {
	services.StateMachine
	gw_handlers.Handler
	methodConfig Config
	donConfig    *config.DONConfig
	don          gw_handlers.DON
	lggr         logger.Logger
	codec        api.JsonRPCCodec
	mu           sync.RWMutex

	userRateLimiter   *ratelimit.RateLimiter
	nodeRateLimiter   *ratelimit.RateLimiter
	requestTimeoutSec int

	pendingRequests map[string]*pendingRequest
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
		pendingRequests:   make(map[string]*pendingRequest),
		mu:                sync.RWMutex{},
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
	s.lggr.Infof("Received message from node %s: %v", nodeAddr, resp)

	s.mu.Lock()
	pendingRequest, ok := s.pendingRequests[resp.ID]
	if !ok {
		promRequestInternalError.WithLabelValues(s.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return fmt.Errorf("no pending request found for ID: %s", resp.ID)
	}
	s.mu.Unlock()
	defer delete(s.pendingRequests, resp.ID)

	// SENDING DUMMY RESPONSE FOR NOW
	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		promRequestInternalError.WithLabelValues(s.donConfig.DonId, api.NodeReponseEncodingError.String()).Inc()
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				resp.ID,
				api.ToJSONRPCErrorCode(api.NodeReponseEncodingError),
				fmt.Sprintf("Failed to marshal response: %v", err),
				nil,
			),
			ErrorCode: api.NodeReponseEncodingError,
		}, pendingRequest.callbackCh)
	}
	responseObj := gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.NoError,
	}
	// END OF DUMMY RESPONSE

	select {
	case <-pendingRequest.timeoutCtx.Done():
		promRequestInternalError.WithLabelValues(s.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return fmt.Errorf("request %s timed out", resp.ID)
	case pendingRequest.callbackCh <- responseObj:
		promRequestSuccess.WithLabelValues(s.donConfig.DonId).Inc()
		s.lggr.Infof("Processed response for request %s from node %s", resp.ID, nodeAddr)
		return nil
	case <-ctx.Done():
		promRequestInternalError.WithLabelValues(s.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return ctx.Err()
	}
}

	s.lggr.Infof("Processed response for request %s from node %s", resp.ID, nodeAddr)

	return s.sendResponse(ctx, responseObj, pendingRequest.callbackCh)
}

func (s *service) handleSecretsCreate(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {

	var req SecretsCreateRequest
	if err := json.Unmarshal(*jsonRequest.Params, &req); err != nil {
		promRequestUserError.WithLabelValues(s.donConfig.DonId).Inc()
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
	if req.ID == "" || req.Value == "" {
		promRequestUserError.WithLabelValues(s.donConfig.DonId).Inc()
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.InvalidParamsError),
				"Secret ID and value cannot be empty",
				nil,
			),
			ErrorCode: api.InvalidParamsError,
		}, callbackCh)
	}

	s.mu.Lock()
	s.pendingRequests[jsonRequest.ID] = &pendingRequest{
		callbackCh: callbackCh,
		responses:  make(map[string]*jsonrpc.Response[json.RawMessage]),
	}
	s.mu.Unlock()

	// At this point, we know that the request is valid and we can send it to the nodes
	var nodeErrors []error
	for _, node := range s.donConfig.Members {
		err := s.don.SendToNode(ctx, node.Address, &jsonRequest)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			s.lggr.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	// If all nodes failed, return an error to the user
	if len(nodeErrors) == len(s.donConfig.Members) && len(nodeErrors) > 0 {
		promRequestInternalError.WithLabelValues(s.donConfig.DonId, "all_nodes_failed").Inc()
		s.lggr.Errorw("all nodes failed to process request", "request_id", jsonRequest.ID, "error_count", len(nodeErrors))
		return s.sendResponse(ctx, gw_handlers.UserCallbackPayload{
			RawResponse: s.codec.EncodeNewErrorResponse(
				jsonRequest.ID,
				api.ToJSONRPCErrorCode(api.HandlerError),
				"Internal error",
				nil,
			),
			ErrorCode: api.HandlerError,
		}, callbackCh)
	}

	s.lggr.Infof("Processed request: %v", jsonRequest)

	// Block until the channel receives a response
	return nil
}

func (s *service) handleUnsupportedMethod(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	s.lggr.Debugw("unsupported method", "method", jsonRequest.Method)
	promRequestUserError.WithLabelValues(s.donConfig.DonId).Inc()
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
