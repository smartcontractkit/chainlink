package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var _ handlers.Handler = (*gatewayHandler)(nil)

const (
	handlerName                          = "HTTPCapabilityHandler"
	defaultCleanUpPeriodMs               = 1000 * 60 * 10 // 10 minutes
	defaultMaxTriggerRequestDurationMs   = 1000 * 60      // 1 minute
	defaultInitialIntervalMs             = 100
	defaultMaxIntervalTimeMs             = 1000 * 30 // 30 seconds
	defaultMultiplier                    = 2.0
	defaultMetadataPullIntervalMs        = 1000 * 60 // 1 minute
	defaultMetadataAggregationIntervalMs = 1000 * 60 // 1 minute
	defaultMetadataPullRequestTimeoutMs  = 1000 * 30 // 30 seconds
	internalErrorMessage                 = "Internal server error occurred while processing the request"
	defaultOutboundRequestCacheTTLMs     = 1000 * 60 * 10 // 10 minutes
)

type gatewayHandler struct {
	services.StateMachine
	config          ServiceConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *ratelimit.RateLimiter // Rate limiter for node requests (e.g. outgoing HTTP requests, HTTP trigger response, auth metadata exchange)
	userRateLimiter limits.RateLimiter     // Rate limiter for user requests that trigger workflow executions
	wg              sync.WaitGroup
	stopCh          services.StopChan
	responseCache   ResponseCache // Caches HTTP responses to avoid redundant requests for outbound HTTP actions
	triggerHandler  HTTPTriggerHandler
	metadataHandler *WorkflowMetadataHandler // Handles authorization for HTTP trigger requests
	metrics         *metrics.Metrics
}

type ResponseCache interface {
	// Set caches a response if it is cacheable (2xx or 4xx status codes) and the cache is empty or expired for the given request.
	Set(workflowID string, req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse)

	// Fetch retrieves a response from the cache if it exists and the age of cached response is less than the max age of the request.
	// If the cached response is expired or not cached, it fetches a new response from the fetchFn.
	// The response is cached if it is cacheable and storeOnFetch is true.
	Fetch(ctx context.Context, workflowID string, req gateway_common.OutboundHTTPRequest, fetchFn func() gateway_common.OutboundHTTPResponse, storeOnFetch bool) gateway_common.OutboundHTTPResponse

	// DeleteExpired removes all cached responses that have exceeded their TTL (Time To Live).
	DeleteExpired(ctx context.Context) int
}

type ServiceConfig struct {
	NodeRateLimiter               ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	MaxTriggerRequestDurationMs   int                         `json:"maxTriggerRequestDurationMs"`
	RetryConfig                   RetryConfig                 `json:"retryConfig"`
	CleanUpPeriodMs               int                         `json:"cleanUpPeriodMs"`
	MetadataPullIntervalMs        int                         `json:"metadataPullIntervalMs"`
	MetadataAggregationIntervalMs int                         `json:"metadataAggregationIntervalMs"`
	MetadataPullRequestTimeoutMs  int                         `json:"metadataPullRequestTimeoutMs"`
	OutboundRequestCacheTTLMs     int                         `json:"outboundRequestCacheTTLMs"`
}

type RetryConfig struct {
	InitialIntervalMs int     `json:"initialIntervalMs"`
	MaxIntervalTimeMs int     `json:"maxIntervalTimeMs"`
	Multiplier        float64 `json:"multiplier"`
}

func NewGatewayHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger, lf limits.Factory) (*gatewayHandler, error) {
	var cfg ServiceConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	cfg = WithDefaults(cfg)
	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to create node rate limiter: %w", err)
	}
	userRateLimiter, err := lf.MakeRateLimiter(cresettings.Default.PerWorkflow.HTTPTrigger.RateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create user rate limiter: %w", err)
	}

	metrics, err := metrics.NewMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	metadataHandler := NewWorkflowMetadataHandler(lggr, cfg, don, donConfig, metrics)
	triggerHandler := NewHTTPTriggerHandler(lggr, cfg, donConfig, don, metadataHandler, userRateLimiter, metrics)
	return &gatewayHandler{
		config:          cfg,
		don:             don,
		donConfig:       donConfig,
		lggr:            logger.With(logger.Named(lggr, handlerName), "donId", donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		userRateLimiter: userRateLimiter,
		stopCh:          make(services.StopChan),
		responseCache:   newResponseCache(lggr, cfg.OutboundRequestCacheTTLMs, metrics),
		triggerHandler:  triggerHandler,
		metadataHandler: metadataHandler,
		metrics:         metrics,
	}, nil
}

func WithDefaults(cfg ServiceConfig) ServiceConfig {
	if cfg.CleanUpPeriodMs == 0 {
		cfg.CleanUpPeriodMs = defaultCleanUpPeriodMs
	}
	if cfg.MaxTriggerRequestDurationMs == 0 {
		cfg.MaxTriggerRequestDurationMs = defaultMaxTriggerRequestDurationMs
	}
	if cfg.MetadataPullIntervalMs == 0 {
		cfg.MetadataPullIntervalMs = defaultMetadataPullIntervalMs
	}
	if cfg.MetadataAggregationIntervalMs == 0 {
		cfg.MetadataAggregationIntervalMs = defaultMetadataPullIntervalMs
	}
	if cfg.MetadataPullRequestTimeoutMs == 0 {
		cfg.MetadataPullRequestTimeoutMs = defaultMetadataPullRequestTimeoutMs
	}
	if cfg.RetryConfig.InitialIntervalMs == 0 {
		cfg.RetryConfig.InitialIntervalMs = defaultInitialIntervalMs
	}
	if cfg.RetryConfig.MaxIntervalTimeMs == 0 {
		cfg.RetryConfig.MaxIntervalTimeMs = defaultMaxIntervalTimeMs
	}
	if cfg.RetryConfig.Multiplier == 0 {
		cfg.RetryConfig.Multiplier = defaultMultiplier
	}
	if cfg.OutboundRequestCacheTTLMs == 0 {
		cfg.OutboundRequestCacheTTLMs = defaultOutboundRequestCacheTTLMs
	}
	return cfg
}

func (h *gatewayHandler) Methods() []string {
	return []string{
		gateway_common.MethodWorkflowExecute,
		gateway_common.MethodHTTPAction,
		gateway_common.MethodPushWorkflowMetadata,
		gateway_common.MethodPullWorkflowMetadata,
	}
}

func (h *gatewayHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	if resp.ID == "" {
		return fmt.Errorf("received response with empty request ID from node %s", nodeAddr)
	}
	h.lggr.Debugw("handling incoming node message", "requestID", resp.ID, "nodeAddr", nodeAddr)
	nodeAllow, globalAllow := h.nodeRateLimiter.AllowVerbose(nodeAddr)
	if !nodeAllow {
		h.metrics.Common.IncrementCapabilityNodeThrottled(ctx, nodeAddr, h.lggr)
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	if !globalAllow {
		h.metrics.Common.IncrementGlobalThrottled(ctx, h.lggr)
		return errors.New("global rate limit exceeded")
	}
	// Node messages follow the format "<methodName>/<workflowID>/<uuid>" or
	// "<methodName>/<workflowID>/<workflowExecutionID>/<uuid>". Messages are routed
	// based on the method in the ID.
	// Any messages without "/" is assumed to be a trigger response to a prior user request.
	if strings.Contains(resp.ID, "/") {
		if resp.Result == nil {
			h.lggr.Errorw("received response with empty result from node", "nodeAddr", nodeAddr, "error", resp.Error)
			return fmt.Errorf("received response with empty result from node %s", nodeAddr)
		}
		parts := strings.Split(resp.ID, "/")
		methodName := parts[0]
		switch methodName {
		case gateway_common.MethodHTTPAction:
			start := time.Now()
			h.metrics.Action.IncrementRequestCount(ctx, nodeAddr, h.lggr)
			err := h.makeOutgoingRequest(ctx, resp, nodeAddr)
			if err != nil {
				h.metrics.Action.IncrementRequestFailures(ctx, nodeAddr, h.lggr)
			}
			h.metrics.Action.RecordRequestLatency(ctx, time.Since(start).Milliseconds(), h.lggr)
			return err
		case gateway_common.MethodPushWorkflowMetadata:
			h.metrics.Trigger.IncrementMetadataRequestCount(ctx, nodeAddr, gateway_common.MethodPushWorkflowMetadata, h.lggr)
			err := h.metadataHandler.OnMetadataPush(ctx, resp, nodeAddr)
			if err != nil {
				h.metrics.Trigger.IncrementMetadataProcessingFailures(ctx, nodeAddr, gateway_common.MethodPushWorkflowMetadata, h.lggr)
			}
			return err
		case gateway_common.MethodPullWorkflowMetadata:
			h.metrics.Trigger.IncrementMetadataRequestCount(ctx, nodeAddr, gateway_common.MethodPullWorkflowMetadata, h.lggr)
			err := h.metadataHandler.OnMetadataPullResponse(ctx, resp, nodeAddr)
			if err != nil {
				h.metrics.Trigger.IncrementMetadataProcessingFailures(ctx, nodeAddr, gateway_common.MethodPullWorkflowMetadata, h.lggr)
			}
			return err
		default:
			return fmt.Errorf("unsupported method %s in node message ID %s", methodName, resp.ID)
		}
	}
	return h.triggerHandler.HandleNodeTriggerResponse(ctx, resp, nodeAddr)
}

// createHTTPRequestCallback creates a callback function that makes the actual HTTP request
func (h *gatewayHandler) createHTTPRequestCallback(ctx context.Context, requestID string, httpReq network.HTTPRequest, req gateway_common.OutboundHTTPRequest) func() gateway_common.OutboundHTTPResponse {
	return func() gateway_common.OutboundHTTPResponse {
		l := logger.With(h.lggr, "requestID", requestID, "method", req.Method, "timeout", req.TimeoutMs)
		l.Debug("Sending request to client")
		start := time.Now()
		resp, err := h.httpClient.Send(ctx, httpReq)
		externalEndpointLatency := time.Since(start)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			isExternalEndpointError := errors.Is(err, network.ErrHTTPSend) || errors.Is(err, network.ErrHTTPRead)
			return gateway_common.OutboundHTTPResponse{
				ErrorMessage:            err.Error(),
				IsExternalEndpointError: isExternalEndpointError,
				ExternalEndpointLatency: externalEndpointLatency,
			}
		}
		h.metrics.Action.IncrementCustomerEndpointResponseCount(ctx, strconv.Itoa(resp.StatusCode), h.lggr)
		h.metrics.Action.RecordCustomerEndpointRequestLatency(ctx, time.Since(start).Milliseconds(), h.lggr)
		return gateway_common.OutboundHTTPResponse{
			StatusCode:              resp.StatusCode,
			Headers:                 resp.Headers,
			Body:                    resp.Body,
			ExternalEndpointLatency: externalEndpointLatency,
		}
	}
}

// extractWorkflowIDFromRequestPath extracts the workflowID from an outgoing request path string.
// The workflowID is expected to be the first element after splitting the string by "/".
func extractWorkflowIDFromRequestPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func (h *gatewayHandler) HandleLegacyUserMessage(context.Context, *api.Message, handlers.Callback) error {
	return errors.New("HTTP capability gateway handler does not support legacy messages")
}

func (h *gatewayHandler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callback handlers.Callback) error {
	h.metrics.Trigger.IncrementRequestCount(ctx, h.lggr)
	err := h.triggerHandler.HandleUserTriggerRequest(ctx, &req, callback, time.Now())
	if err != nil {
		h.lggr.Errorw("failed to handle user trigger request", "requestID",
			req.ID, "err", err)
		// error response is sent to the response channel by the trigger handler
		// so return nil after logging
	}
	return nil
}

func (h *gatewayHandler) makeOutgoingRequest(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	requestID := resp.ID
	h.lggr.Debugw("handling webAPI outgoing message", "requestID", requestID, "nodeAddr", nodeAddr)
	var req gateway_common.OutboundHTTPRequest
	err := json.Unmarshal(*resp.Result, &req)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP request from node %s: %w", nodeAddr, err)
	}
	workflowID := extractWorkflowIDFromRequestPath(requestID)
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
		var outboundResp gateway_common.OutboundHTTPResponse
		callback := h.createHTTPRequestCallback(newCtx, requestID, httpReq, req)
		if req.CacheSettings.MaxAgeMs > 0 {
			h.metrics.Action.IncrementCacheReadCount(ctx, h.lggr)
			outboundResp = h.responseCache.Fetch(ctx, workflowID, req, callback, req.CacheSettings.Store)
		} else {
			outboundResp = callback()
			if req.CacheSettings.Store {
				h.responseCache.Set(workflowID, req, outboundResp)
			}
		}
		h.metrics.Action.IncrementCapabilityRequestCount(ctx, nodeAddr, h.lggr)
		err := h.sendResponseToNode(newCtx, requestID, outboundResp, nodeAddr)
		if err != nil {
			l.Errorw("error sending response to node", "err", err, "nodeAddr", nodeAddr, "requestID", requestID)
			h.metrics.Action.IncrementCapabilityFailures(ctx, nodeAddr, h.lggr)
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

func (h *gatewayHandler) Start(ctx context.Context) error {
	return h.StartOnce(handlerName, func() error {
		h.lggr.Info("Starting " + handlerName)
		err := h.triggerHandler.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start HTTP trigger handler: %w", err)
		}
		err = h.metadataHandler.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start HTTP auth handler: %w", err)
		}
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
			ticker := time.NewTicker(time.Duration(h.config.CleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.responseCache.DeleteExpired(ctx)
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
		err := h.triggerHandler.Close()
		if err != nil {
			h.lggr.Errorw("failed to close HTTP trigger handler", "err", err)
		}
		err = h.metadataHandler.Close()
		if err != nil {
			h.lggr.Errorw("failed to close HTTP auth handler", "err", err)
		}
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
		Version: jsonrpc.JsonRpcVersion,
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
