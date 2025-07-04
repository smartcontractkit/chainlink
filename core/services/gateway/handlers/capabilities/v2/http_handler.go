package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var _ handlers.Handler = (*gatewayHandler)(nil)

const (
	handlerName            = "HTTPCapabilityHandler"
	defaultCleanUpPeriodMs = 1000 * 60 * 10 // 10 minutes
)

type gatewayHandler struct {
	services.StateMachine
	config          ServiceConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *ratelimit.RateLimiter // Rate limiter for node requests (e.g. outgoing HTTP requests, HTTP trigger response, auth metadata exchange)
	userRateLimiter *ratelimit.RateLimiter // Rate limiter for user requests that trigger workflow executions
	wg              sync.WaitGroup
	stopCh          services.StopChan
	responseCache   ResponseCache
}

type ResponseCache interface {
	Set(req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse, ttl time.Duration)
	Get(req gateway_common.OutboundHTTPRequest) *gateway_common.OutboundHTTPResponse
	DeleteExpired() int
}

type ServiceConfig struct {
	NodeRateLimiter ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	UserRateLimiter ratelimit.RateLimiterConfig `json:"userRateLimiter"`
	CleanUpPeriodMs int                         `json:"cleanUpPeriodMs"`
}

func NewGatewayHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*gatewayHandler, error) {
	var cfg ServiceConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	cfg = WithDefaults(cfg)
	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, err
	}
	userRateLimiter, err := ratelimit.NewRateLimiter(cfg.UserRateLimiter)
	if err != nil {
		return nil, err
	}
	return &gatewayHandler{
		config:          cfg,
		don:             don,
		donConfig:       donConfig,
		lggr:            logger.With(logger.Named(lggr, handlerName), "donId", donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		userRateLimiter: userRateLimiter,
		wg:              sync.WaitGroup{},
		stopCh:          make(services.StopChan),
		responseCache:   newResponseCache(lggr),
	}, nil
}

func WithDefaults(cfg ServiceConfig) ServiceConfig {
	if cfg.CleanUpPeriodMs == 0 {
		cfg.CleanUpPeriodMs = defaultCleanUpPeriodMs
	}
	return cfg
}

func (h *gatewayHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	if resp.ID == "" {
		return fmt.Errorf("received response with empty request ID from node %s", nodeAddr)
	}
	if resp.Result == nil {
		return fmt.Errorf("received response with nil result from node %s", nodeAddr)
	}
	h.lggr.Debugw("handling incoming node message", "requestID", resp.ID, "nodeAddr", nodeAddr)
	var outboundReq gateway_common.OutboundHTTPRequest
	err := json.Unmarshal(*resp.Result, &outboundReq)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP request from node %s: %w", nodeAddr, err)
	}
	return h.handleOutgoingRequest(ctx, resp.ID, outboundReq, nodeAddr)
}

func (h *gatewayHandler) HandleLegacyUserMessage(context.Context, *api.Message, chan<- handlers.UserCallbackPayload) error {
	return errors.New("HTTP capability gateway handler does not support legacy messages")
}

func (h *gatewayHandler) HandleJSONRPCUserMessage(context.Context, jsonrpc.Request[json.RawMessage], chan<- handlers.UserCallbackPayload) error {
	// TODO: Implement trigger request handling
	return nil
}

func (h *gatewayHandler) handleOutgoingRequest(ctx context.Context, requestID string, req gateway_common.OutboundHTTPRequest, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "requestID", requestID, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}

	if req.CacheSettings.ReadFromCache {
		cached := h.responseCache.Get(req)
		if cached != nil {
			h.lggr.Debugw("Using cached HTTP response", "requestID", requestID, "nodeAddr", nodeAddr)
			return h.sendResponseToNode(ctx, requestID, *cached, nodeAddr)
		}
	}

	timeout := time.Duration(req.TimeoutMs) * time.Millisecond
	httpReq := network.HTTPRequest{
		Method:           req.Method,
		URL:              req.URL,
		Headers:          req.Headers,
		Body:             req.Body,
		MaxResponseBytes: req.MaxResponseBytes,
		Timeout:          timeout,
	}

	// send response to node async
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := logger.With(h.lggr, "requestID", requestID, "method", req.Method, "timeout", req.TimeoutMs)
		l.Debug("Sending request to client")
		var outboundResp gateway_common.OutboundHTTPResponse
		resp, err := h.httpClient.Send(ctx, httpReq)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			outboundResp = gateway_common.OutboundHTTPResponse{
				ErrorMessage: err.Error(),
			}
		} else {
			outboundResp = gateway_common.OutboundHTTPResponse{
				StatusCode: resp.StatusCode,
				Headers:    resp.Headers,
				Body:       resp.Body,
			}

			if req.CacheSettings.StoreInCache && isCacheableStatusCode(resp.StatusCode) {
				cacheTTLMs := req.CacheSettings.TTLMs
				if cacheTTLMs > 0 {
					h.responseCache.Set(req, outboundResp, time.Duration(cacheTTLMs)*time.Millisecond)
					l.Debugw("Cached HTTP response", "ttlMs", cacheTTLMs)
				}
			}
		}
		err = h.sendResponseToNode(newCtx, requestID, outboundResp, nodeAddr)
		if err != nil {
			l.Errorw("error sending response to node", "err", err, "nodeAddr", nodeAddr)
		}
	}()
	return nil
}

func (h *gatewayHandler) HealthReport() map[string]error {
	return map[string]error{handlerName: h.Healthy()}
}

func (h *gatewayHandler) Name() string {
	return handlerName
}

func (h *gatewayHandler) Start(context.Context) error {
	return h.StartOnce(handlerName, func() error {
		h.lggr.Info("Starting " + handlerName)
		go func() {
			ticker := time.NewTicker(time.Duration(h.config.CleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.responseCache.DeleteExpired()
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *gatewayHandler) Close() error {
	return h.StopOnce(handlerName, func() error {
		h.lggr.Info("Closing " + handlerName)
		close(h.stopCh)
		h.wg.Wait()
		return nil
	})
}

func (h *gatewayHandler) sendResponseToNode(ctx context.Context, requestID string, resp gateway_common.OutboundHTTPResponse, nodeAddr string) error {
	params, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	rawParams := json.RawMessage(params)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      requestID,
		Method:  gateway_common.MethodHTTPAction,
		Params:  &rawParams,
	}

	err = h.don.SendToNode(ctx, nodeAddr, req)
	if err != nil {
		return err
	}

	h.lggr.Debugw("sent response to node", "to", nodeAddr)
	return nil
}

// isCacheableStatusCode returns true if the HTTP status code indicates a cacheable response.
// This includes successful responses (2xx) and client errors (4xx)
func isCacheableStatusCode(statusCode int) bool {
	return (statusCode >= 200 && statusCode < 300) || (statusCode >= 400 && statusCode < 500)
}
