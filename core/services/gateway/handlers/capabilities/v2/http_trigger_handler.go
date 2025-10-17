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

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ HTTPTriggerHandler = (*httpTriggerHandler)(nil)

const (
	// Reference: https://github.com/smartcontractkit/chainlink-evm/blob/develop/contracts/src/v0.8/workflow/dev/v2/WorkflowRegistry.sol
	workflowIDLength       = 66 // 0x + 64 hex characters = 32 bytes
	workflowOwnerLength    = 42 // 0x + 40 hex characters = 20 bytes
	maxWorkflowNameLength  = 64 // Maximum workflow name length
	WorkflowNameHashLength = 22 // 0x + 20 hex characters = 10 bytes
	maxWorkflowTagLength   = 32 // Maximum workflow tag length
)

type savedCallback struct {
	handlers.Callback
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
	userRateLimiter         limits.RateLimiter
	metrics                 *metrics.Metrics
	wg                      sync.WaitGroup
}

type HTTPTriggerHandler interface {
	job.ServiceCtx
	HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback, requestStartTime time.Time) error
	HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error
}

func NewHTTPTriggerHandler(lggr logger.Logger, cfg ServiceConfig, donConfig *config.DONConfig, don handlers.DON, workflowMetadataHandler *WorkflowMetadataHandler, userRateLimiter limits.RateLimiter, metrics *metrics.Metrics) *httpTriggerHandler {
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

func (h *httpTriggerHandler) HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback, requestStartTime time.Time) error {
	triggerReq, err := h.validatedTriggerRequest(ctx, req, callback)
	if err != nil {
		return err
	}

	workflowID, err := h.resolveWorkflowID(ctx, triggerReq, req.ID, callback)
	if err != nil {
		return err
	}

	key, err := h.authorizeRequest(ctx, workflowID, req, callback)
	if err != nil {
		return err
	}

	if err = h.checkRateLimit(ctx, workflowID, req.ID, callback); err != nil {
		return err
	}

	executionID, err := workflows.EncodeExecutionID(strings.TrimPrefix(workflowID, "0x"), req.ID)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, internalErrorMessage, callback)
		return errors.New("error generating execution ID: " + err.Error())
	}

	reqWithKey, err := reqWithAuthorizedKey(triggerReq, *key)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInvalidRequest, "Auth failure", callback)
		return errors.Join(errors.New("auth failure"), err)
	}

	if err := h.setupCallback(ctx, req.ID, callback, requestStartTime); err != nil {
		return err
	}

	return h.sendWithRetries(ctx, executionID, reqWithKey)
}

func (h *httpTriggerHandler) validatedTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback) (*jsonrpc.Request[gateway_common.HTTPTriggerRequest], error) {
	if req.Params == nil {
		h.handleUserError(ctx, "", jsonrpc.ErrInvalidRequest, "request params is nil", callback)
		return nil, errors.New("request params is nil")
	}

	triggerReq, err := h.parseTriggerRequest(ctx, req, callback)
	if err != nil {
		return nil, err
	}

	if err := h.validateRequestID(ctx, req.ID, callback); err != nil {
		return nil, err
	}

	if err := h.validateMethod(ctx, req.Method, req.ID, callback); err != nil {
		return nil, err
	}

	if err := h.validateTriggerParams(ctx, triggerReq, req.ID, callback); err != nil {
		return nil, err
	}

	return &jsonrpc.Request[gateway_common.HTTPTriggerRequest]{
		Version: req.Version,
		ID:      req.ID,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  triggerReq,
	}, nil
}

func (h *httpTriggerHandler) parseTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback) (*gateway_common.HTTPTriggerRequest, error) {
	var triggerReq gateway_common.HTTPTriggerRequest
	err := json.Unmarshal(*req.Params, &triggerReq)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrParse, "error decoding payload: "+err.Error(), callback)
		return nil, err
	}
	return &triggerReq, nil
}

func (h *httpTriggerHandler) validateRequestID(ctx context.Context, requestID string, callback handlers.Callback) error {
	if requestID == "" {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "empty request ID", callback)
		return errors.New("empty request ID")
	}
	// Request IDs from users must not contain "/", since this character is reserved
	// for internal node-to-node message routing (e.g., "http_action/{workflowID}/{uuid}").
	if strings.Contains(requestID, "/") {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "request ID must not contain '/'", callback)
		return errors.New("request ID must not contain '/'")
	}
	return nil
}

func (h *httpTriggerHandler) validateMethod(ctx context.Context, method, requestID string, callback handlers.Callback) error {
	if method != gateway_common.MethodWorkflowExecute {
		h.handleUserError(ctx, requestID, jsonrpc.ErrMethodNotFound, "invalid method: "+method, callback)
		return errors.New("invalid method: " + method)
	}
	return nil
}

func (h *httpTriggerHandler) validateTriggerParams(ctx context.Context, triggerReq *gateway_common.HTTPTriggerRequest, requestID string, callback handlers.Callback) error {
	if !isValidJSON(triggerReq.Input) {
		h.lggr.Errorw("invalid params JSON", "params", triggerReq.Input)
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "invalid params JSON", callback)
		return errors.New("invalid params JSON")
	}

	return h.validateWorkflowFields(ctx, triggerReq.Workflow, requestID, callback)
}

func (h *httpTriggerHandler) validateWorkflowFields(ctx context.Context, workflow gateway_common.WorkflowSelector, requestID string, callback handlers.Callback) error {
	hasWorkflowID := workflow.WorkflowID != ""
	hasWorkflowName := workflow.WorkflowName != ""
	hasWorkflowOwner := workflow.WorkflowOwner != ""
	hasWorkflowTag := workflow.WorkflowTag != ""

	if !hasWorkflowID {
		if !hasWorkflowName {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowName is required when workflowID is not provided", callback)
			return errors.New("workflowName is required when workflowID is not provided")
		}
		if !hasWorkflowOwner {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowOwner is required when workflowID is not provided", callback)
			return errors.New("workflowOwner is required when workflowID is not provided")
		}
		if !hasWorkflowTag {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowTag is required when workflowID is not provided", callback)
			return errors.New("workflowTag is required when workflowID is not provided")
		}
	}

	if hasWorkflowID {
		if err := h.validateWorkflowID(ctx, workflow.WorkflowID, requestID, callback); err != nil {
			return err
		}
	}
	if hasWorkflowOwner {
		if err := h.validateWorkflowOwner(ctx, workflow.WorkflowOwner, requestID, callback); err != nil {
			return err
		}
	}
	if hasWorkflowName {
		if err := h.validateWorkflowName(ctx, workflow.WorkflowName, requestID, callback); err != nil {
			return err
		}
	}

	if workflow.WorkflowTag != "" {
		if err := h.validateWorkflowTag(ctx, workflow.WorkflowTag, requestID, callback); err != nil {
			return err
		}
	}

	return nil
}

func validateHexInput(input string, expectedLength int) error {
	if input != strings.ToLower(input) {
		return errors.New("must be lowercase")
	}

	if len(input) > expectedLength {
		return fmt.Errorf("hex string too long: expected at most %d characters, got %d", expectedLength, len(input))
	}

	hexStr := strings.TrimPrefix(input, "0x")
	_, err := hex.DecodeString(hexStr)
	if err != nil {
		return errors.New("must be a valid hex string")
	}

	return nil
}

func (h *httpTriggerHandler) validateWorkflowID(ctx context.Context, workflowID string, requestID string, callback handlers.Callback) error {
	if err := validateHexInput(workflowID, workflowIDLength); err != nil {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowID "+err.Error(), callback)
		return errors.New("workflowID " + err.Error())
	}

	return nil
}

func (h *httpTriggerHandler) validateWorkflowOwner(ctx context.Context, workflowOwner string, requestID string, callback handlers.Callback) error {
	if err := validateHexInput(workflowOwner, workflowOwnerLength); err != nil {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowOwner "+err.Error(), callback)
		return errors.New("workflowOwner " + err.Error())
	}

	return nil
}

// validateWorkflowName validates the workflowName length and format
func (h *httpTriggerHandler) validateWorkflowName(ctx context.Context, workflowName string, requestID string, callback handlers.Callback) error {
	if len(workflowName) == 0 {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowName cannot be empty", callback)
		return errors.New("workflowName cannot be empty")
	}

	if len(workflowName) > maxWorkflowNameLength {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, fmt.Sprintf("workflowName cannot exceed %d characters, got %d", maxWorkflowNameLength, len(workflowName)), callback)
		return fmt.Errorf("workflowName cannot exceed %d characters, got %d", maxWorkflowNameLength, len(workflowName))
	}

	return nil
}

// validateWorkflowTag validates the workflowTag length and format
func (h *httpTriggerHandler) validateWorkflowTag(ctx context.Context, workflowTag string, requestID string, callback handlers.Callback) error {
	if len(workflowTag) == 0 {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowTag cannot be empty", callback)
		return errors.New("workflowTag cannot be empty")
	}

	if len(workflowTag) > maxWorkflowTagLength {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, fmt.Sprintf("workflowTag cannot exceed %d characters, got %d", maxWorkflowTagLength, len(workflowTag)), callback)
		return fmt.Errorf("workflowTag cannot exceed %d characters, got %d", maxWorkflowTagLength, len(workflowTag))
	}

	return nil
}

// normalizeHex normalizes a hex string by stripping 0x prefix, padding with leading zeros, and adding 0x prefix back
func normalizeHex(input string, length int) string {
	hexStr := strings.TrimPrefix(input, "0x")
	// length-2 because we'll add "0x" prefix
	expectedHexLength := length - 2
	paddedHex := strings.Repeat("0", expectedHexLength-len(hexStr)) + hexStr
	return "0x" + paddedHex
}

func (h *httpTriggerHandler) resolveWorkflowID(ctx context.Context, triggerReq *jsonrpc.Request[gateway_common.HTTPTriggerRequest], requestID string, callback handlers.Callback) (string, error) {
	workflowID := triggerReq.Params.Workflow.WorkflowID
	if workflowID != "" {
		workflowID = normalizeHex(workflowID, workflowIDLength)
		return workflowID, nil
	}
	workflowOwner := normalizeHex(triggerReq.Params.Workflow.WorkflowOwner, workflowOwnerLength)
	workflowName := "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName(triggerReq.Params.Workflow.WorkflowName)))
	workflowID, found := h.workflowMetadataHandler.GetWorkflowID(
		workflowOwner,
		workflowName,
		triggerReq.Params.Workflow.WorkflowTag,
	)
	if !found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflow not found", callback)
		return "", errors.New("workflow not found")
	}
	return workflowID, nil
}

func (h *httpTriggerHandler) authorizeRequest(ctx context.Context, workflowID string, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback) (*gateway_common.AuthorizedKey, error) {
	key, err := h.workflowMetadataHandler.Authorize(workflowID, req.Auth, req)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInvalidRequest, "Auth failure", callback)
		return nil, errors.Join(errors.New("auth failure"), err)
	}
	return key, nil
}

func (h *httpTriggerHandler) checkRateLimit(ctx context.Context, workflowID, requestID string, callback handlers.Callback) error {
	workflowRef, found := h.workflowMetadataHandler.GetWorkflowReference(workflowID)
	if !found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflow reference not found", callback)
		return errors.New("workflow reference not found")
	}
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: workflowRef.workflowOwner, Workflow: workflowID})
	if err := h.userRateLimiter.AllowErr(ctx); err != nil {
		lggr := logger.With(h.lggr, platform.KeyWorkflowID, workflowID, platform.KeyWorkflowOwner, workflowRef.workflowOwner, "requestID", requestID, "err", err)
		var errLimited limits.ErrorRateLimited
		if errors.As(err, &errLimited) {
			switch errLimited.Scope {
			case settings.ScopeWorkflow:
				lggr.Errorf("failed to start execution: per workflow rate limit exceeded")
				h.metrics.Trigger.IncrementWorkflowThrottled(ctx, h.lggr)
			default:
				lggr.Errorf("failed to start execution: unexpected rate limit for scope %s", errLimited.Scope)
			}
			h.handleUserError(ctx, requestID, jsonrpc.ErrLimitExceeded, "rate limit exceeded", callback)
			return err
		}
	}
	return nil
}

func (h *httpTriggerHandler) setupCallback(ctx context.Context, requestID string, callback handlers.Callback, requestStartTime time.Time) error {
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()

	if _, found := h.callbacks[requestID]; found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrConflict, "in-flight request", callback)
		return errors.New("in-flight request ID: " + requestID)
	}

	// (N+F)//2 + 1 threshold where N = number of nodes, F = number of faulty nodes
	threshold := (len(h.donConfig.Members)+h.donConfig.F)/2 + 1
	agg, err := aggregation.NewIdenticalNodeResponseAggregator(threshold)
	if err != nil {
		return errors.New("failed to create response aggregator: " + err.Error())
	}

	h.callbacks[requestID] = savedCallback{
		Callback:           callback,
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

	err = saved.SendResponse(handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   api.NoError,
	})
	if err != nil {
		return err
	}
	delete(h.callbacks, resp.ID)
	latencyMs := time.Since(saved.requestStartTime).Milliseconds()
	h.metrics.Trigger.RecordRequestHandlerLatency(ctx, latencyMs, h.lggr)
	return nil
}

func (h *httpTriggerHandler) Start(ctx context.Context) error {
	return h.StartOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Starting HTTPTriggerHandler")
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
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
		h.wg.Wait()
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
			h.metrics.Trigger.IncrementRequestErrors(ctx, jsonrpc.ErrInternal, h.lggr)
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

func (h *httpTriggerHandler) handleUserError(ctx context.Context, requestID string, code int64, message string, callback handlers.Callback) {
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
	errorCode := api.FromJSONRPCErrorCode(code)
	h.metrics.Trigger.IncrementRequestErrors(ctx, code, h.lggr)
	err = callback.SendResponse(handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   errorCode,
	})
	if err != nil {
		h.lggr.Errorw("failed to send user callback", "err", err, "requestID", requestID)
		return
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
