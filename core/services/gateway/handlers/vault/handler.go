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

const (
	handlerName          = "VaultHandler"
	defaultCleanUpPeriod = 5 * time.Second
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
	createdAt  time.Time
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
	stopCh       services.StopChan

	nodeRateLimiter *ratelimit.RateLimiter
	requestTimeout  time.Duration

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
	NodeRateLimiterConfig ratelimit.RateLimiterConfig `json:"node_rate_limiter"`
	RequestTimeoutSec     int                         `json:"request_timeout_sec"`
}

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gw_handlers.DON, lggr logger.Logger) (*handler, error) {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal method config: %w", err)
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create node rate limiter: %w", err)
	}

	return &handler{
		methodConfig:    cfg,
		donConfig:       donConfig,
		don:             don,
		lggr:            logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
		requestTimeout:  time.Duration(cfg.RequestTimeoutSec) * time.Second,
		nodeRateLimiter: nodeRateLimiter,
		pendingRequests: make(map[string]*pendingRequest),
		mu:              sync.RWMutex{},
		stopCh:          make(services.StopChan),
	}, nil
}

func (h *handler) Start(ctx context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		ctx, _ := h.stopCh.NewCtx()
		go func() {
			ticker := time.NewTicker(defaultCleanUpPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.removeExpiredRequests(ctx)
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce("VaultHandler", func() error {
		h.lggr.Info("closing vault handler")
		close(h.stopCh)
		return nil
	})
}

// removeExpiredRequests removes expired requests from the pending requests map
func (h *handler) removeExpiredRequests(ctx context.Context) {
	h.mu.RLock()
	var expiredRequests []*pendingRequest
	now := time.Now()
	for _, pendingRequest := range h.pendingRequests {
		if now.Sub(pendingRequest.createdAt) > h.requestTimeout {
			expiredRequests = append(expiredRequests, pendingRequest)
		}
	}
	h.mu.RUnlock()

	for _, pendingRequest := range expiredRequests {
		err := h.sendResponse(ctx, pendingRequest.req.ID, h.errorResponse(pendingRequest.req, api.RequestTimeoutError), pendingRequest.callbackCh)
		if err != nil {
			h.lggr.Errorw("error sending response to user", "request_id", pendingRequest.req.ID, "error", err)
		}
	}
}

func (h *handler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	return errors.New("vault handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handling vault request", "method", req.Method, "id", req.ID)

	h.mu.Lock()
	h.pendingRequests[req.ID] = &pendingRequest{
		callbackCh: callbackCh,
		req:        &req,
		createdAt:  time.Now(),
	}
	h.mu.Unlock()
	switch req.Method {
	case MethodSecretsCreate:
		return h.handleSecretsCreate(ctx, &req, callbackCh)
	default:
		return h.sendResponse(ctx, req.ID, h.errorResponse(&req, api.UnsupportedMethodError), callbackCh)
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h.lggr.Debugw(fmt.Sprintf("Received response: %v", resp), "nodeAddr", nodeAddr)

	if !h.nodeRateLimiter.Allow(nodeAddr) {
		h.lggr.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	pendingRequest, ok := h.pendingRequests[resp.ID]
	if !ok {
		promRequestInternalError.WithLabelValues(h.donConfig.DonId, api.RequestTimeoutError.String()).Inc()
		return fmt.Errorf("no pending request found for ID: %s", resp.ID)
	}

	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		return h.sendResponse(ctx, resp.ID, h.errorResponse(pendingRequest.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err)), pendingRequest.callbackCh)
	}

	return h.sendResponse(ctx, resp.ID, gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.FromJSONRPCErrorCode(resp.Error.Code),
	}, pendingRequest.callbackCh)
}

func (h *handler) handleSecretsCreate(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	var secretsCreateRequest SecretsCreateRequest
	if err := json.Unmarshal(*req.Params, &secretsCreateRequest); err != nil {
		return h.sendResponse(ctx, req.ID, h.errorResponse(req, api.UserMessageParseError, err), callbackCh)
	}

	if secretsCreateRequest.ID == "" || secretsCreateRequest.Value == "" {
		return h.sendResponse(ctx, req.ID, h.errorResponse(req, api.InvalidParamsError, errors.New("secret id and value cannot be empty")), callbackCh)
	}

	// At this point, we know that the request is valid and we can send it to the nodes
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, req)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			h.lggr.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, req.ID, h.errorResponse(req, api.FatalError, errors.New("failed to forward user request to nodes")), callbackCh)
	}

	h.lggr.Debugf("Forwarded request to Vault nodes: %v", req)
	return nil
}

func (h *handler) errorResponse(req *jsonrpc.Request[json.RawMessage], errorCode api.ErrorCode, errs ...error) gw_handlers.UserCallbackPayload {
	err := errors.New("unknown error")
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}

	switch errorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
		h.lggr.Errorw(err.Error(), "request_id", req.ID)
		// Intentionally hide the error from the user
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		h.lggr.Errorw("invalid params", "request_id", req.ID, "params", string(*req.Params))
		err = errors.New("invalid params error: " + err.Error())
	case api.UnsupportedMethodError:
		h.lggr.Errorw("unsupported method", "request_id", req.ID, "method", req.Method)
		err = errors.New("unsupported method: " + req.Method)
	case api.UserMessageParseError:
		h.lggr.Errorw("user message parse error", "request_id", req.ID, "error", err.Error())
		err = errors.New("user message parse error: " + err.Error())
	case api.NoError:
	case api.UnsupportedDONIdError:
	case api.HandlerError:
	case api.RequestTimeoutError:
		// Unused in this handler
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

func (h *handler) sendResponse(ctx context.Context, reqID string, resp gw_handlers.UserCallbackPayload, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	h.mu.Lock()
	delete(h.pendingRequests, reqID)
	h.mu.Unlock()

	switch resp.ErrorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
	case api.RequestTimeoutError:
	case api.HandlerError:
		promRequestInternalError.WithLabelValues(h.donConfig.DonId, resp.ErrorCode.String()).Inc()
	case api.InvalidParamsError:
	case api.UnsupportedMethodError:
	case api.UserMessageParseError:
	case api.UnsupportedDONIdError:
		promRequestUserError.WithLabelValues(h.donConfig.DonId).Inc()
	case api.NoError:
		promRequestSuccess.WithLabelValues(h.donConfig.DonId).Inc()
	}

	select {
	case callbackCh <- resp:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
