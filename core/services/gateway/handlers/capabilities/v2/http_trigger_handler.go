package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ HTTPTriggerHandler = (*httpTriggerHandler)(nil)

type savedCallback struct {
	callbackCh         chan<- handlers.UserCallbackPayload
	requestStartTime   time.Time
	createdAt          time.Time
	responseAggregator *aggregation.IdenticalNodeResponseAggregator
}

type httpTriggerHandler struct {
	services.StateMachine
	config                  ServiceConfig
	don                     handlers.DON
	donConfig               *config.DONConfig
	lggr                    logger.Logger
	callbacksMu             sync.Mutex
	callbacks               map[string]savedCallback // requestID -> savedCallback
	stopCh                  services.StopChan
	workflowMetadataHandler *WorkflowMetadataHandler
	userRateLimiter         *ratelimit.RateLimiter
	metrics                 *metrics.Metrics
}

type HTTPTriggerHandler interface {
	job.ServiceCtx
	HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload, requestStartTime time.Time) error
	HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error
}

func NewHTTPTriggerHandler(lggr logger.Logger, cfg ServiceConfig, donConfig *config.DONConfig, don handlers.DON, workflowMetadataHandler *WorkflowMetadataHandler, userRateLimiter *ratelimit.RateLimiter, metrics *metrics.Metrics) *httpTriggerHandler {
	return &httpTriggerHandler{
		lggr:                    logger.Named(lggr, "RequestCallbacks"),
		callbacks:               make(map[string]savedCallback),
		config:                  cfg,
		don:                     don,
		donConfig:               donConfig,
		stopCh:                  make(services.StopChan),
		workflowMetadataHandler: workflowMetadataHandler,
		userRateLimiter:         userRateLimiter,
		metrics:                 metrics,
	}
}

func (h *httpTriggerHandler) HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload, requestStartTime time.Time) error {
	triggerReq, err := h.validatedTriggerRequest(ctx, req, callbackCh)
	if err != nil {
		return err
	}

	workflowID, err := h.resolveWorkflowID(ctx, triggerReq, req.ID, callbackCh)
	if err != nil {
		return err
	}

	key, err := h.authorizeRequest(ctx, workflowID, req, callbackCh)
	if err != nil {
		return err
	}

	if err = h.checkRateLimit(ctx, workflowID, req.ID, callbackCh); err != nil {
		return err
	}

	executionID, err := workflows.EncodeExecutionID(workflowID, req.ID)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, internalErrorMessage, callbackCh)
		return errors.New("error generating execution ID: " + err.Error())
	}

	reqWithKey, err := reqWithAuthorizedKey(triggerReq, *key)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInvalidRequest, "Auth failure", callbackCh)
		return errors.Join(errors.New("auth failure"), err)
	}

	if err := h.setupCallback(ctx, req.ID, callbackCh, requestStartTime); err != nil {
		return err
	}

	return h.sendWithRetries(ctx, executionID, reqWithKey)
}

func (h *httpTriggerHandler) validatedTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload) (*jsonrpc.Request[gateway_common.HTTPTriggerRequest], error) {
	if req.Params == nil {
		h.handleUserError(ctx, "", jsonrpc.ErrInvalidRequest, "request params is nil", callbackCh)
		return nil, errors.New("request params is nil")
	}

	triggerReq, err := h.parseTriggerRequest(ctx, req, callbackCh)
	if err != nil {
		return nil, err
	}

	if err := h.validateRequestID(ctx, req.ID, callbackCh); err != nil {
		return nil, err
	}

	if err := h.validateMethod(ctx, req.Method, req.ID, callbackCh); err != nil {
		return nil, err
	}

	if err := h.validateTriggerParams(ctx, triggerReq, req.ID, callbackCh); err != nil {
		return nil, err
	}

	return &jsonrpc.Request[gateway_common.HTTPTriggerRequest]{
		Version: req.Version,
		ID:      req.ID,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  triggerReq,
	}, nil
}

func (h *httpTriggerHandler) parseTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload) (*gateway_common.HTTPTriggerRequest, error) {
	var triggerReq gateway_common.HTTPTriggerRequest
	err := json.Unmarshal(*req.Params, &triggerReq)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrParse, "error decoding payload: "+err.Error(), callbackCh)
		return nil, err
	}
	return &triggerReq, nil
}

func (h *httpTriggerHandler) validateRequestID(ctx context.Context, requestID string, callbackCh chan<- handlers.UserCallbackPayload) error {
	if requestID == "" {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "empty request ID", callbackCh)
		return errors.New("empty request ID")
	}
	// Request IDs from users must not contain "/", since this character is reserved
	// for internal node-to-node message routing (e.g., "http_action/{workflowID}/{uuid}").
	if strings.Contains(requestID, "/") {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "request ID must not contain '/'", callbackCh)
		return errors.New("request ID must not contain '/'")
	}
	return nil
}

func (h *httpTriggerHandler) validateMethod(ctx context.Context, method, requestID string, callbackCh chan<- handlers.UserCallbackPayload) error {
	if method != gateway_common.MethodWorkflowExecute {
		h.handleUserError(ctx, requestID, jsonrpc.ErrMethodNotFound, "invalid method: "+method, callbackCh)
		return errors.New("invalid method: " + method)
	}
	return nil
}

func (h *httpTriggerHandler) validateTriggerParams(ctx context.Context, triggerReq *gateway_common.HTTPTriggerRequest, requestID string, callbackCh chan<- handlers.UserCallbackPayload) error {
	if !isValidJSON(triggerReq.Input) {
		h.lggr.Errorw("invalid params JSON", "params", triggerReq.Input)
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "invalid params JSON", callbackCh)
		return errors.New("invalid params JSON")
	}

	return h.validateWorkflowFields(ctx, triggerReq.Workflow, requestID, callbackCh)
}

func (h *httpTriggerHandler) validateWorkflowFields(ctx context.Context, workflow gateway_common.WorkflowSelector, requestID string, callbackCh chan<- handlers.UserCallbackPayload) error {
	if workflow.WorkflowID != "" {
		if !strings.HasPrefix(workflow.WorkflowID, "0x") {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowID must be prefixed with '0x'", callbackCh)
			return errors.New("workflowID must be prefixed with '0x'")
		}
		if workflow.WorkflowID != strings.ToLower(workflow.WorkflowID) {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowID must be lowercase", callbackCh)
			return errors.New("workflowID must be lowercase")
		}
	}
	if workflow.WorkflowOwner != "" {
		if !strings.HasPrefix(workflow.WorkflowOwner, "0x") {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowOwner must be prefixed with '0x'", callbackCh)
			return errors.New("workflowOwner must be prefixed with '0x'")
		}
		if workflow.WorkflowOwner != strings.ToLower(workflow.WorkflowOwner) {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowOwner must be lowercase", callbackCh)
			return errors.New("workflowOwner must be lowercase")
		}
	}
	return nil
}

func (h *httpTriggerHandler) resolveWorkflowID(ctx context.Context, triggerReq *jsonrpc.Request[gateway_common.HTTPTriggerRequest], requestID string, callbackCh chan<- handlers.UserCallbackPayload) (string, error) {
	workflowID := triggerReq.Params.Workflow.WorkflowID
	if workflowID != "" {
		return workflowID, nil
	}

	workflowName := "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName(triggerReq.Params.Workflow.WorkflowName)))
	workflowID, found := h.workflowMetadataHandler.GetWorkflowID(
		triggerReq.Params.Workflow.WorkflowOwner,
		workflowName,
		triggerReq.Params.Workflow.WorkflowTag,
	)
	if !found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflow not found", callbackCh)
		return "", errors.New("workflow not found")
	}
	return workflowID, nil
}

func (h *httpTriggerHandler) authorizeRequest(ctx context.Context, workflowID string, req *jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload) (*gateway_common.AuthorizedKey, error) {
	key, err := h.workflowMetadataHandler.Authorize(workflowID, req.Auth, req)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInvalidRequest, "Auth failure", callbackCh)
		return nil, errors.Join(errors.New("auth failure"), err)
	}
	return key, nil
}

func (h *httpTriggerHandler) checkRateLimit(ctx context.Context, workflowID, requestID string, callbackCh chan<- handlers.UserCallbackPayload) error {
	workflowRef, found := h.workflowMetadataHandler.GetWorkflowReference(workflowID)
	if !found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflow reference not found", callbackCh)
		return errors.New("workflow reference not found")
	}
	workflowOwnerAllow, globalAllow := h.userRateLimiter.AllowVerbose(workflowRef.workflowOwner)
	if !workflowOwnerAllow {
		h.metrics.Trigger.IncrementWorkflowOwnerThrottled(ctx, h.lggr)
		h.handleUserError(ctx, requestID, jsonrpc.ErrLimitExceeded, "rate limit exceeded", callbackCh)
		return errors.New("workflow owner rate limit exceeded")
	}
	if !globalAllow {
		h.metrics.Trigger.IncrementGlobalThrottled(ctx, h.lggr)
		h.handleUserError(ctx, requestID, jsonrpc.ErrLimitExceeded, "rate limit exceeded", callbackCh)
		return errors.New("global rate limit exceeded")
	}
	return nil
}

func (h *httpTriggerHandler) setupCallback(ctx context.Context, requestID string, callbackCh chan<- handlers.UserCallbackPayload, requestStartTime time.Time) error {
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()

	if _, found := h.callbacks[requestID]; found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrConflict, "in-flight request", callbackCh)
		return errors.New("in-flight request ID: " + requestID)
	}

	// 2f + 1 is chosen to ensure that majority of honest nodes are executing the request
	agg, err := aggregation.NewIdenticalNodeResponseAggregator(2*h.donConfig.F + 1)
	if err != nil {
		return errors.New("failed to create response aggregator: " + err.Error())
	}

	h.callbacks[requestID] = savedCallback{
		callbackCh:         callbackCh,
		requestStartTime:   requestStartTime,
		createdAt:          time.Now(),
		responseAggregator: agg,
	}
	return nil
}

func (h *httpTriggerHandler) HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h.lggr.Debugw("handling trigger response", "requestID", resp.ID, "nodeAddr", nodeAddr, "error", resp.Error, "result", resp.Result)
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()
	saved, exists := h.callbacks[resp.ID]
	if !exists {
		return errors.New("callback not found for request ID: " + resp.ID)
	}
	aggResp, err := saved.responseAggregator.CollectAndAggregate(resp, nodeAddr)
	if err != nil {
		return err
	}
	if aggResp == nil {
		h.lggr.Debugw("Not enough responses to aggregate", "requestID", resp.ID, "nodeAddress", nodeAddr)
		return nil
	}
	rawResp, err := json.Marshal(aggResp)
	if err != nil {
		return errors.New("failed to marshal response: " + err.Error())
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case saved.callbackCh <- handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   api.NoError,
	}:
		delete(h.callbacks, resp.ID)
		latencyMs := time.Since(saved.requestStartTime).Milliseconds()
		h.metrics.Trigger.RecordRequestHandlerLatency(ctx, latencyMs, h.lggr)
	}
	return nil
}

func (h *httpTriggerHandler) Start(ctx context.Context) error {
	return h.StartOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Starting HTTPTriggerHandler")
		go func() {
			ticker := time.NewTicker(time.Duration(h.config.CleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.reapExpiredCallbacks(ctx)
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *httpTriggerHandler) Close() error {
	return h.StopOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Closing HTTPTriggerHandler")
		close(h.stopCh)
		return nil
	})
}

// reapExpiredCallbacks removes callbacks that are older than the maximum age
func (h *httpTriggerHandler) reapExpiredCallbacks(ctx context.Context) {
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()
	now := time.Now()
	var expiredCount int
	for reqID, callback := range h.callbacks {
		if now.Sub(callback.createdAt) > time.Duration(h.config.MaxTriggerRequestDurationMs)*time.Millisecond {
			delete(h.callbacks, reqID)
			expiredCount++
		}
	}
	if expiredCount > 0 {
		h.metrics.Trigger.IncrementPendingRequestsCleanUpCount(ctx, int64(expiredCount), h.lggr)
		h.lggr.Infow("Removed expired callbacks", "count", expiredCount, "remaining", len(h.callbacks))
	}
	h.metrics.Trigger.RecordPendingRequestsCount(ctx, int64(len(h.callbacks)), h.lggr)
}

func isValidJSON(data []byte) bool {
	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return false
	}

	switch val.(type) {
	case map[string]any, []any:
		return true
	default:
		return false
	}
}

func (h *httpTriggerHandler) handleUserError(ctx context.Context, requestID string, code int64, message string, callbackCh chan<- handlers.UserCallbackPayload) {
	resp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      requestID,
		Error: &jsonrpc.WireError{
			Code:    code,
			Message: message,
		},
	}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		h.lggr.Errorw("failed to marshal error response", "err", err, "requestID", requestID)
		return
	}
	errorCode := api.ErrorCode(code)
	h.metrics.Trigger.IncrementRequestErrors(ctx, errorCode.String(), h.lggr)
	callbackCh <- handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   errorCode,
	}
}

// sendWithRetries attempts to send the request to all DON members,
// retrying failed nodes until either all succeed or the max trigger request duration is reached.
func (h *httpTriggerHandler) sendWithRetries(ctx context.Context, executionID string, req *jsonrpc.Request[json.RawMessage]) error {
	// Create a context that will be cancelled when the max request duration is reached
	maxDuration := time.Duration(h.config.MaxTriggerRequestDurationMs) * time.Millisecond
	ctxWithTimeout, cancel := context.WithTimeout(ctx, maxDuration)
	defer cancel()

	successfulNodes := make(map[string]bool)
	b := backoff.Backoff{
		Min:    time.Duration(h.config.RetryConfig.InitialIntervalMs) * time.Millisecond,
		Max:    time.Duration(h.config.RetryConfig.MaxIntervalTimeMs) * time.Millisecond,
		Factor: h.config.RetryConfig.Multiplier,
		Jitter: true,
	}

	for {
		// Retry sending to nodes that haven't received the message
		allNodesSucceeded := true
		var combinedErr error

		for _, member := range h.donConfig.Members {
			if successfulNodes[member.Address] {
				continue
			}
			h.metrics.Trigger.IncrementCapabilityRequestCount(ctx, member.Address, gateway_common.MethodWorkflowExecute, h.lggr)
			err := h.don.SendToNode(ctxWithTimeout, member.Address, req)
			if err != nil {
				allNodesSucceeded = false
				h.metrics.Trigger.IncrementCapabilityRequestFailures(ctx, member.Address, gateway_common.MethodWorkflowExecute, h.lggr)
				err = errors.Join(combinedErr, err)
				h.lggr.Debugw("Failed to send trigger request to node, will retry",
					"node", member.Address,
					"executionID", executionID,
					"error", err)
			} else {
				// Mark this node as successful
				successfulNodes[member.Address] = true
			}
		}

		if allNodesSucceeded {
			h.lggr.Infow("Successfully sent trigger request to all nodes",
				"executionID", executionID,
				"nodeCount", len(h.donConfig.Members))
			return nil
		}

		// Not all nodes succeeded, wait and retry
		h.lggr.Debugw("Retrying failed nodes for trigger request",
			"executionID", executionID,
			"failedCount", len(h.donConfig.Members)-len(successfulNodes),
			"errors", combinedErr)

		select {
		case <-time.After(b.Duration()):
			continue
		case <-ctxWithTimeout.Done():
			return fmt.Errorf("request retry time exceeded, some nodes may not have received the request: executionID=%s, successNodes=%d, totalNodes=%d",
				executionID, len(successfulNodes), len(h.donConfig.Members))
		}
	}
}

func reqWithAuthorizedKey(req *jsonrpc.Request[gateway_common.HTTPTriggerRequest], key gateway_common.AuthorizedKey) (*jsonrpc.Request[json.RawMessage], error) {
	params := *req.Params
	params.Key = key
	msg, err := json.Marshal(params)
	if err != nil {
		return nil, errors.New("error marshaling trigger request")
	}
	rawMsg := json.RawMessage(msg)
	r := &jsonrpc.Request[json.RawMessage]{
		Version: req.Version,
		ID:      req.ID,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawMsg,
	}
	return r, err
}
