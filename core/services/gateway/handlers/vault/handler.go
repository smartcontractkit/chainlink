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
	_ gw_handlers.Handler = (*handler)(nil)

	ErrNotAllowlisted = errors.New("sender not allowlisted")
	ErrRateLimited    = errors.New("rate-limited")

	promRequestInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_node_vault_request_internal_error",
		Help: "Gateway node, Vault handler: Metric to track internal errors",
	}, []string{"don_id", "error"})

	promRequestUserError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_node_vault_request_user_error",
		Help: "Gateway node, Vault handler: Metric to track failed requests",
	}, []string{"don_id"})

	promRequestSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_node_vault_request_success",
		Help: "Gateway node, Vault handler: Metric to track successful requests",
	}, []string{"don_id"})
)

type pendingRequest struct {
	req        *jsonrpc.Request[json.RawMessage]
	callbackCh chan<- gw_handlers.UserCallbackPayload
}

type handler struct {
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

func (h *handler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *handler) Name() string {
	return h.lggr.Name()
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

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gw_handlers.DON, lggr logger.Logger) *handler {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		// Return a minimal implementation that will fail gracefully
		return &handler{
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

	return &handler{
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

func (h *handler) Start(ctx context.Context) error {
	return h.StartOnce("VaultService", func() error {
		h.lggr.Info("starting vault service")
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce("VaultMethod", func() error {
		h.lggr.Info("closing vault service")
		return nil
	})
}

func (h *handler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	return errors.New("vault service does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handling vault request", "method", jsonRequest.Method, "id", jsonRequest.ID)
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(h.requestTimeoutSec)*time.Second)
	defer cancel()

	switch jsonRequest.Method {
	case MethodSecretsCreate:
		return h.handleSecretsCreate(timeoutCtx, &jsonRequest, callbackCh)
	default:
		return h.sendResponse(timeoutCtx, h.errorResponse(&jsonRequest, api.UnsupportedMethodError), callbackCh)
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h.lggr.Debugw(fmt.Sprintf("Received response: %v", resp), "nodeAddr", nodeAddr)

	h.mu.Lock()
	defer h.mu.Unlock()
	pendingRequest, ok := h.pendingRequests[resp.ID]
	if !ok {
		promRequestInternalError.WithLabelValues(h.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return fmt.Errorf("no pending request found for ID: %s", resp.ID)
	}
	defer delete(h.pendingRequests, resp.ID)

	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		return h.sendResponse(ctx, h.errorResponse(pendingRequest.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err)), pendingRequest.callbackCh)
	}
	responseObj := gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.NoError,
	}

	select {
	case pendingRequest.callbackCh <- responseObj:
		promRequestSuccess.WithLabelValues(h.donConfig.DonId).Inc()
		h.lggr.Infof("Processed response for request %s from node %s", resp.ID, nodeAddr)
		return nil
	case <-ctx.Done():
		promRequestInternalError.WithLabelValues(h.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return ctx.Err()
	}
}

func (h *handler) handleSecretsCreate(ctx context.Context, jsonRequest *jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	var req SecretsCreateRequest
	if err := json.Unmarshal(*jsonRequest.Params, &req); err != nil {
		return h.sendResponse(ctx, h.errorResponse(jsonRequest, api.UserMessageParseError, err), callbackCh)
	}

	if req.ID == "" || req.Value == "" {
		return h.sendResponse(ctx, h.errorResponse(jsonRequest, api.InvalidParamsError, errors.New("secret id and value cannot be empty")), callbackCh)
	}

	h.mu.Lock()
	h.pendingRequests[jsonRequest.ID] = &pendingRequest{
		callbackCh: callbackCh,
		req:        jsonRequest,
	}
	h.mu.Unlock()

	// At this point, we know that the request is valid and we can send it to the nodes
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, jsonRequest)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			h.lggr.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, h.errorResponse(jsonRequest, api.FatalError, errors.New("failed to forward user request to nodes")), callbackCh)
	}

	h.lggr.Debugf("Forwarded request to Vault nodes: %v", jsonRequest)
	return nil
}

func (h *handler) sendResponse(ctx context.Context, response gw_handlers.UserCallbackPayload, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	select {
	case callbackCh <- response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *handler) errorResponse(req *jsonrpc.Request[json.RawMessage], errorCode api.ErrorCode, errs ...error) gw_handlers.UserCallbackPayload {
	err := errors.New("unknown error")
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}

	switch errorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
		promRequestInternalError.WithLabelValues(h.donConfig.DonId, errorCode.String()).Inc()
		h.lggr.Errorw(err.Error(), "request_id", req.ID)
		// Intentionally hide the error from the user
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		promRequestUserError.WithLabelValues(h.donConfig.DonId).Inc()
		h.lggr.Errorw("invalid params", "request_id", req.ID, "params", string(*req.Params))
		err = errors.New("invalid params error: " + err.Error())
	case api.UnsupportedMethodError:
		promRequestUserError.WithLabelValues(h.donConfig.DonId).Inc()
		h.lggr.Errorw("unsupported method", "request_id", req.ID, "method", req.Method)
		err = errors.New("unsupported method: " + req.Method)
	case api.UserMessageParseError:
		promRequestUserError.WithLabelValues(h.donConfig.DonId).Inc()
		h.lggr.Errorw("user message parse error", "request_id", req.ID, "error", err.Error())
		err = errors.New("user message parse error: " + err.Error())
	case api.NoError:
	case api.UnsupportedDONIdError:
	case api.HandlerError:
	case api.RequestTimeoutError:
		// Unimplemented
	}

	return gw_handlers.UserCallbackPayload{
		RawResponse: h.codec.EncodeNewErrorResponse(
			req.ID,
			api.ToJSONRPCErrorCode(errorCode),
			err.Error(),
			nil,
		),
		ErrorCode: errorCode,
	}
}
