package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
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
)

type metrics struct {
	requestInternalError metric.Int64Counter
	requestUserError     metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}

	requestUserError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_user_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register user error counter: %w", err)
	}

	requestSuccess, err := beholder.GetMeter().Int64Counter("gateway_vault_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}

	return &metrics{
		requestInternalError: requestInternalError,
		requestUserError:     requestUserError,
		requestSuccess:       requestSuccess,
	}, nil
}

type activeRequest struct {
	req        jsonrpc.Request[json.RawMessage]
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

	activeRequests map[string]activeRequest
	metrics        *metrics
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

	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &handler{
		methodConfig:    cfg,
		donConfig:       donConfig,
		don:             don,
		lggr:            logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
		requestTimeout:  time.Duration(cfg.RequestTimeoutSec) * time.Second,
		nodeRateLimiter: nodeRateLimiter,
		activeRequests:  make(map[string]activeRequest),
		mu:              sync.RWMutex{},
		stopCh:          make(services.StopChan),
		metrics:         metrics,
	}, nil
}

func (h *handler) Start(ctx context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		go func() {
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
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
	var expiredRequests []activeRequest
	now := time.Now()
	for _, userRequest := range h.activeRequests {
		if now.Sub(userRequest.createdAt) > h.requestTimeout {
			expiredRequests = append(expiredRequests, userRequest)
		}
	}
	h.mu.RUnlock()

	for _, er := range expiredRequests {
		err := h.sendResponse(ctx, er, h.errorResponse(er.req, api.RequestTimeoutError))
		if err != nil {
			h.lggr.Errorw("error sending response to user", "request_id", er.req.ID, "error", err)
		}
	}
}

func (h *handler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	return errors.New("vault handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handling vault request", "method", req.Method, "id", req.ID)
	ar := activeRequest{
		callbackCh: callbackCh,
		req:        req,
		createdAt:  time.Now(),
	}

	h.mu.Lock()
	h.activeRequests[req.ID] = ar
	h.mu.Unlock()
	switch req.Method {
	case MethodSecretsCreate:
		return h.handleSecretsCreate(ctx, ar)
	default:
		return h.sendResponse(ctx, ar, h.errorResponse(req, api.UnsupportedMethodError))
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h.lggr.Debugw(fmt.Sprintf("Received response: %v", resp), "nodeAddr", nodeAddr)

	if !h.nodeRateLimiter.Allow(nodeAddr) {
		h.lggr.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}

	h.mu.RLock()
	ur, ok := h.activeRequests[resp.ID]
	h.mu.RUnlock()
	if !ok {
		// Request is not found, so we don't need to send a response to the user
		// This might happen if the response is stale
		h.lggr.Errorw("no pending request found for ID", "request_id", resp.ID)
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", api.StaleNodeResponseError.String()),
		))
		return nil
	}

	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		return h.sendResponse(ctx, ur, h.errorResponse(ur.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err)))
	}

	var errorCode api.ErrorCode
	if resp.Error != nil {
		errorCode = api.FromJSONRPCErrorCode(resp.Error.Code)
	} else {
		errorCode = api.NoError
	}

	successResp := gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   errorCode,
	}
	return h.sendResponse(ctx, ur, successResp)
}

func (h *handler) handleSecretsCreate(ctx context.Context, ar activeRequest) error {
	var secretsCreateRequest SecretsCreateRequest
	if err := json.Unmarshal(*ar.req.Params, &secretsCreateRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	if secretsCreateRequest.ID == "" || secretsCreateRequest.Value == "" {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id and value cannot be empty")))
	}

	// At this point, we know that the request is valid and we can send it to the nodes
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, &ar.req)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			h.lggr.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	h.lggr.Debugf("Forwarded request to Vault nodes: %v", ar.req)
	return nil
}

func (h *handler) errorResponse(
	req jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	errs ...error,
) gw_handlers.UserCallbackPayload {
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
	case api.StaleNodeResponseError:
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

func (h *handler) sendResponse(ctx context.Context, userRequest activeRequest, resp gw_handlers.UserCallbackPayload) error {
	switch resp.ErrorCode {
	case api.StaleNodeResponseError:
	case api.FatalError:
	case api.NodeReponseEncodingError:
	case api.RequestTimeoutError:
	case api.HandlerError:
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", resp.ErrorCode.String()),
		))
	case api.InvalidParamsError:
	case api.UnsupportedMethodError:
	case api.UserMessageParseError:
	case api.UnsupportedDONIdError:
		h.metrics.requestUserError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	case api.NoError:
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	}

	select {
	case userRequest.callbackCh <- resp:
		h.lggr.Debugw("sent response", "request_id", userRequest.req.ID)
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return nil
	case <-ctx.Done():
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return ctx.Err()
	}
}
