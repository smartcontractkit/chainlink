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
	requestStartTime time.Time
	createdAt        time.Time
	// processed is set once the aggregated response has been sent to the user.
	// The entry stays in the callbacks map so late node responses can be told
	// apart from unknown request IDs, and is removed later by the async reaper.
	processed bool
	// responseAggregators holds one IdenticalNodeResponseAggregator per shard the
	// workflow is assigned to, keyed by shard donID. The first shard to reach its
	// quorum produces the user response.
	responseAggregators map[string]*aggregation.IdenticalNodeResponseAggregator
	doneCh              chan struct{} // closed when callback is responded to (processed) or reaped unanswered. signals sendWithRetries to stop retrying
}

type httpTriggerHandler struct {
	services.StateMachine
	config                  ServiceConfig
	shards                  []*shardEndpoint
	nodeAddrToShard         map[string]*shardEndpoint
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

func NewHTTPTriggerHandler(lggr logger.Logger, cfg ServiceConfig, shards []*shardEndpoint, nodeAddrToShard map[string]*shardEndpoint, workflowMetadataHandler *WorkflowMetadataHandler, userRateLimiter limits.RateLimiter, metrics *metrics.Metrics) *httpTriggerHandler {
	return &httpTriggerHandler{
		lggr:                    logger.Named(lggr, "RequestCallbacks"),
		callbacks:               make(map[string]savedCallback),
		config:                  cfg,
		shards:                  shards,
		nodeAddrToShard:         nodeAddrToShard,
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

	strippedWorkflowID := strings.TrimPrefix(workflowID, "0x")
	legacyExecutionID, err := workflows.EncodeExecutionID(strippedWorkflowID, req.ID) //nolint:staticcheck // legacy ID kept for observability comparison
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, internalErrorMessage, callback)
		return errors.New("error generating execution ID: " + err.Error())
	}
	// Workflows shouldn't use more than one HTTP trigger. If we ever need to support multiple triggers, we'd need to pass
	// trigger index to the Gateway handler and somehow allow senders to pick. For now, we use trigger index 0.
	// Execution IDs here are used only for logging.
	executionIDWithTriggerIndex, err := workflows.GenerateExecutionIDWithTriggerIndex(strippedWorkflowID, req.ID, 0)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, internalErrorMessage, callback)
		return errors.New("error generating execution ID with trigger index: " + err.Error())
	}
	h.lggr.Debugw("processing request",
		"legacyExecutionID", legacyExecutionID,
		"executionIDWithTriggerIndex", executionIDWithTriggerIndex,
		"requestID", req.ID,
		"workflowID", workflowID)

	reqWithKey, err := reqWithAuthorizedKey(triggerReq, *key)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, internalErrorMessage, callback)
		return errors.New("error marshaling trigger request: " + err.Error())
	}

	doneCh, err := h.setupCallback(ctx, req.ID, callback, requestStartTime, workflowID)
	if err != nil {
		return err
	}

	return h.sendWithRetries(ctx, legacyExecutionID, executionIDWithTriggerIndex, reqWithKey, workflowID, doneCh)
}

func (h *httpTriggerHandler) validatedTriggerRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback) (*jsonrpc.Request[gateway_common.HTTPTriggerRequest], error) {
	if req.Params == nil {
		h.handleUserError(ctx, "", jsonrpc.ErrInvalidRequest, "'params' field is missing. Include a valid 'params' object", callback)
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
		h.handleUserError(ctx, req.ID, jsonrpc.ErrParse, "payload is not a valid JSON. Ensure that the request body is a well-formed JSON", callback)
		return nil, err
	}
	return &triggerReq, nil
}

func (h *httpTriggerHandler) validateRequestID(ctx context.Context, requestID string, callback handlers.Callback) error {
	if requestID == "" {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "'id' field is required and cannot be empty. Use a new unique request 'id' for each request", callback)
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
		h.handleUserError(ctx, requestID, jsonrpc.ErrMethodNotFound, fmt.Sprintf("'%s' is not a valid method. Ensure that method is set to 'workflows.execute'", method), callback)
		return errors.New("invalid method: " + method)
	}
	return nil
}

func (h *httpTriggerHandler) validateTriggerParams(ctx context.Context, triggerReq *gateway_common.HTTPTriggerRequest, requestID string, callback handlers.Callback) error {
	if !isValidJSON(triggerReq.Input) {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "'params' must be {} or [] (JSON object or array). Primitives (null, '', numbers, booleans) are not allowed. Use {} if none.", callback)
		return errors.New("invalid params JSON: " + string(triggerReq.Input))
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

func (h *httpTriggerHandler) validateWorkflowID(ctx context.Context, workflowID, requestID string, callback handlers.Callback) error {
	if err := validateHexInput(workflowID, workflowIDLength); err != nil {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowID "+err.Error(), callback)
		return errors.New("workflowID " + err.Error())
	}

	return nil
}

func (h *httpTriggerHandler) validateWorkflowOwner(ctx context.Context, workflowOwner, requestID string, callback handlers.Callback) error {
	if err := validateHexInput(workflowOwner, workflowOwnerLength); err != nil {
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "workflowOwner "+err.Error(), callback)
		return errors.New("workflowOwner " + err.Error())
	}

	return nil
}

// validateWorkflowName validates the workflowName length and format
func (h *httpTriggerHandler) validateWorkflowName(ctx context.Context, workflowName, requestID string, callback handlers.Callback) error {
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
func (h *httpTriggerHandler) validateWorkflowTag(ctx context.Context, workflowTag, requestID string, callback handlers.Callback) error {
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
	h.lggr.Debugw("resolving workflow ID", "workflowID", triggerReq.Params.Workflow.WorkflowID, "workflowOwner", triggerReq.Params.Workflow.WorkflowOwner, "workflowName", triggerReq.Params.Workflow.WorkflowName, "workflowTag", triggerReq.Params.Workflow.WorkflowTag, "requestID", requestID)
	workflowID := triggerReq.Params.Workflow.WorkflowID
	if workflowID != "" {
		workflowID = normalizeHex(workflowID, workflowIDLength)
		_, found := h.workflowMetadataHandler.GetWorkflowReference(workflowID)
		if !found {
			h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, fmt.Sprintf("Workflow not found. 'workflowID' %s is not a valid workflow ID", workflowID), callback)
			return "", errors.New("workflow not found")
		}
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
		h.handleUserError(ctx, requestID, jsonrpc.ErrInvalidRequest, "Workflow not found. Provide either a valid 'workflowID' or a valid combination of 'workflowOwner', 'workflowName', and 'workflowTag'", callback)
		return "", errors.New("workflow not found")
	}
	return workflowID, nil
}

func (h *httpTriggerHandler) authorizeRequest(ctx context.Context, workflowID string, req *jsonrpc.Request[json.RawMessage], callback handlers.Callback) (*gateway_common.AuthorizedKey, error) {
	h.lggr.Debugw("authorizing request", "workflowID", workflowID, "requestID", req.ID)
	key, err := h.workflowMetadataHandler.Authorize(workflowID, req.Auth, req)
	if err != nil {
		h.handleUserError(ctx, req.ID, jsonrpc.ErrInvalidRequest, "Auth failure: "+err.Error(), callback)
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

	// TODO orgID https://smartcontract-it.atlassian.net/browse/CRE-1707
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: workflowRef.workflowOwner, Workflow: workflowID})
	if err := h.userRateLimiter.AllowErr(ctx); err != nil {
		lggr := logger.With(h.lggr, platform.KeyWorkflowID, workflowID, platform.KeyWorkflowOwner, workflowRef.workflowOwner, "requestID", requestID, "err", err)
		if errLimited, ok := errors.AsType[limits.ErrorRateLimited](err); ok {
			switch errLimited.Scope {
			case settings.ScopeWorkflow:
				lggr.Errorf("failed to start execution: per workflow rate limit exceeded")
				h.metrics.IncrementWorkflowThrottled(ctx, h.lggr)
			default:
				lggr.Errorf("failed to start execution: unexpected rate limit for scope %s", errLimited.Scope)
			}
			h.handleUserError(ctx, requestID, jsonrpc.ErrLimitExceeded, "rate limit exceeded", callback)
			return err
		}
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	return nil
}

func (h *httpTriggerHandler) setupCallback(ctx context.Context, requestID string, callback handlers.Callback, requestStartTime time.Time, workflowID string) (<-chan struct{}, error) {
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()

	if _, found := h.callbacks[requestID]; found {
		h.handleUserError(ctx, requestID, jsonrpc.ErrConflict, fmt.Sprintf("requestID: %s has already been used. Ensure the requestID is unique for each request.", requestID), callback)
		return nil, fmt.Errorf("in-flight request ID: %s", requestID)
	}

	// Build one response aggregator per shard the workflow is assigned to.
	assigned := h.workflowMetadataHandler.WorkflowShards(workflowID)
	if len(assigned) == 0 {
		// this shouldn't happen because we checked it in authorizeRequest()
		h.handleUserError(ctx, requestID, jsonrpc.ErrInternal, fmt.Sprintf("Workflow %s is not assigned to any DONs", workflowID), callback)
		return nil, errors.New("workflow is not assigned to any shards")
	}

	aggregators := make(map[string]*aggregation.IdenticalNodeResponseAggregator, len(assigned))
	for _, shard := range assigned {
		// (N+F)//2 + 1 threshold where N = number of nodes, F = number of faulty nodes
		threshold := (len(shard.members)+shard.f)/2 + 1
		agg, err := aggregation.NewIdenticalNodeResponseAggregator(threshold)
		if err != nil {
			return nil, errors.New("failed to create response aggregator: " + err.Error())
		}
		aggregators[shard.donID] = agg
	}

	doneCh := make(chan struct{})
	h.callbacks[requestID] = savedCallback{
		Callback:            callback,
		requestStartTime:    requestStartTime,
		createdAt:           time.Now(),
		responseAggregators: aggregators,
		doneCh:              doneCh,
	}
	return doneCh, nil
}

// cleanupCallback removes a callback and, if it was never processed, closes
// doneCh to signal sendWithRetries to stop. When the entry was already
// processed, doneCh was closed by markCallbackProcessed, so it is only
// deleted here.
// Must be called while holding callbacksMu lock.
func (h *httpTriggerHandler) cleanupCallback(requestID string) {
	saved, exists := h.callbacks[requestID]
	if !exists {
		return
	}
	if !saved.processed {
		close(saved.doneCh)
	}
	delete(h.callbacks, requestID)
}

func (h *httpTriggerHandler) HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h.lggr.Debugw("handling trigger response", "requestID", resp.ID, "nodeAddr", nodeAddr, "error", resp.Error, "result", resp.Result)
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()
	saved, exists := h.callbacks[resp.ID]
	if !exists {
		return errors.New("callback not found for request ID: " + resp.ID)
	}
	if saved.processed {
		h.lggr.Debugw("request already processed, ignoring late response", "requestID", resp.ID, "nodeAddr", nodeAddr)
		return nil
	}

	// Route the response into the aggregator for the shard that owns this node.
	shard, ok := h.nodeAddrToShard[nodeAddr]
	if !ok {
		return fmt.Errorf("received trigger response from unknown node %s (no owning shard)", nodeAddr)
	}
	agg, ok := saved.responseAggregators[shard.donID]
	if !ok {
		// The node belongs to a shard this workflow isn't assigned to (or the
		// callback was captured before the workflow was assigned there).
		return fmt.Errorf("node %s (shard %s) is not assigned to workflow for request ID %s", nodeAddr, shard.donID, resp.ID)
	}
	aggResp, err := agg.CollectAndAggregate(resp, nodeAddr)
	if err != nil {
		return err
	}
	if aggResp == nil {
		h.lggr.Debugw("Not enough responses to aggregate", "requestID", resp.ID, "nodeAddress", nodeAddr, "shard", shard.donID)
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

	// First shard to reach quorum wins: after successfully sending the response,
	// mark the callback as processed and close doneCh, stopping all shard sends.
	// The entry is kept in the map so late responses are recognized as such and
	// is removed later by the periodic reaper.
	saved.processed = true
	close(saved.doneCh)
	h.callbacks[resp.ID] = saved
	latencyMs := time.Since(saved.requestStartTime).Milliseconds()
	h.metrics.RecordRequestHandlerLatency(ctx, latencyMs, h.lggr)
	h.metrics.IncrementRequestSuccess(ctx, h.lggr)
	h.lggr.Debugw("Sent response to user", "requestID", resp.ID, "latencyMs", latencyMs)
	return nil
}

func (h *httpTriggerHandler) Start(ctx context.Context) error {
	return h.StartOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Starting HTTPTriggerHandler")
		h.wg.Go(func() {
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
		})
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
		if now.Sub(callback.createdAt) > time.Duration(h.config.CleanUpPeriodMs)*time.Millisecond {
			if !callback.processed {
				h.metrics.IncrementRequestErrors(ctx, jsonrpc.ErrInternal, h.lggr)
			}
			h.cleanupCallback(reqID)
			expiredCount++
		}
	}
	if expiredCount > 0 {
		h.metrics.IncrementPendingRequestsCleanUpCount(ctx, int64(expiredCount), h.lggr)
		h.lggr.Infow("Removed expired callbacks", "count", expiredCount, "remaining", len(h.callbacks))
	}
	h.metrics.RecordPendingRequestsCount(ctx, int64(len(h.callbacks)), h.lggr)
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
	switch code {
	case jsonrpc.ErrInternal, jsonrpc.ErrServerOverloaded, jsonrpc.ErrUnknown, jsonrpc.ErrLimitExceeded, jsonrpc.ErrConflict:
		h.lggr.Errorw("returning error to user", "code", code, "message", message, "requestID", requestID)
	default:
		h.lggr.Warnw("returning error to user", "code", code, "message", message, "requestID", requestID)
	}
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
	h.metrics.IncrementRequestErrors(ctx, code, h.lggr)
	err = callback.SendResponse(handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   errorCode,
	})
	if err != nil {
		h.lggr.Errorw("failed to send user callback", "err", err, "requestID", requestID)
		return
	}
}

type nodeSendResult struct {
	nodeAddress string
	err         error
}

// sendWithRetries fans the request out to every shard the workflow is assigned
// to, running an independent per-shard retry loop in parallel. Each per-shard
// loop sends to that shard's members via the shard's own connection manager,
// retrying failed nodes until all succeed or the max trigger request duration
// is reached. Each send attempt is bounded by a per-node timeout, smaller than
// the overall request duration, so that a single slow or unresponsive node can't
// delay delivery to the rest of the DON.
// doneCh is closed when the callback has been responded to (first shard reaches
// quorum), allowing immediate termination of all shard loops.
func (h *httpTriggerHandler) sendWithRetries(ctx context.Context, legacyExecutionID, executionIDWithTriggerIndex string, req *jsonrpc.Request[json.RawMessage], workflowID string, doneCh <-chan struct{}) error {
	if doneCh == nil {
		return errors.New("doneCh cannot be nil")
	}

	assigned := h.workflowMetadataHandler.WorkflowShards(workflowID)
	if len(assigned) == 0 {
		// this shouldn't happen because we checked it in authorizeRequest()
		h.callbacksMu.Lock()
		saved, exists := h.callbacks[req.ID]
		if exists {
			h.handleUserError(ctx, req.ID, jsonrpc.ErrInternal, fmt.Sprintf("Workflow %s is not assigned to any DONs", workflowID), saved.Callback)
			h.cleanupCallback(req.ID)
		}
		h.callbacksMu.Unlock()
		return fmt.Errorf("workflow %s not assigned to any shard", workflowID)
	}

	// Create a context that will be cancelled when the max request duration is reached
	maxDuration := time.Duration(h.config.MaxTriggerRequestDurationMs) * time.Millisecond
	ctxWithTimeout, cancel := context.WithTimeout(ctx, maxDuration)
	defer cancel()

	// Run one send loop per assigned shard in parallel.
	errCh := make(chan error, len(assigned))
	for _, shard := range assigned {
		h.wg.Go(func() {
			errCh <- h.sendToShard(ctxWithTimeout, shard, legacyExecutionID, executionIDWithTriggerIndex, req, doneCh)
		})
	}

	var combinedErr error
	for range assigned {
		if err := <-errCh; err != nil {
			combinedErr = errors.Join(combinedErr, err)
		}
	}
	return combinedErr
}

// sendToShard sends the request to all members of a single shard, retrying
// failures until all succeed, the callback is responded to (doneCh), or ctx is
// cancelled (overall max duration).
func (h *httpTriggerHandler) sendToShard(ctx context.Context, shard *shardEndpoint, legacyExecutionID, executionIDWithTriggerIndex string, req *jsonrpc.Request[json.RawMessage], doneCh <-chan struct{}) error {
	nodeTimeout := time.Duration(h.config.NodeSendTimeoutMs) * time.Millisecond

	successfulNodes := make(map[string]bool)
	b := backoff.Backoff{
		Min:    time.Duration(h.config.RetryConfig.InitialIntervalMs) * time.Millisecond,
		Max:    time.Duration(h.config.RetryConfig.MaxIntervalTimeMs) * time.Millisecond,
		Factor: h.config.RetryConfig.Multiplier,
		Jitter: true,
	}

	for {
		var pending []string
		for _, member := range shard.members {
			if !successfulNodes[member.Address] {
				pending = append(pending, member.Address)
			}
		}

		// Buffered so every goroutine can send its result and exit without waiting on a reader.
		results := make(chan nodeSendResult, len(pending))
		var wg sync.WaitGroup
		for _, nodeAddress := range pending {
			wg.Add(1)
			go func(nodeAddress string) {
				defer wg.Done()

				nodeCtx, nodeCancel := context.WithTimeout(ctx, nodeTimeout)
				defer nodeCancel()

				h.metrics.IncrementTriggerCapabilityRequestCount(ctx, nodeAddress, gateway_common.MethodWorkflowExecute, h.lggr)
				sendStart := time.Now()
				err := shard.connMgr.SendToNode(nodeCtx, nodeAddress, req)
				h.metrics.RecordGatewayToNodeLatency(ctx, time.Since(sendStart).Milliseconds(), nodeAddress, gateway_common.MethodWorkflowExecute, h.lggr)
				if err != nil {
					h.metrics.IncrementTriggerCapabilityRequestFailures(ctx, nodeAddress, gateway_common.MethodWorkflowExecute, h.lggr)
				}
				results <- nodeSendResult{nodeAddress: nodeAddress, err: err}
			}(nodeAddress)
		}
		wg.Wait()
		close(results)

		// Only this goroutine reads from results, so successfulNodes and combinedErr need no locking.
		var combinedErr error
		for res := range results {
			if res.err != nil {
				combinedErr = errors.Join(combinedErr, fmt.Errorf("node %s: %w", res.nodeAddress, res.err))
				h.lggr.Debugw("Failed to send trigger request to node, will retry",
					"node", res.nodeAddress,
					"shard", shard.donID,
					"legacyExecutionID", legacyExecutionID,
					"executionIDWithTriggerIndex", executionIDWithTriggerIndex,
					"error", res.err)
			} else {
				successfulNodes[res.nodeAddress] = true
			}
		}

		if len(successfulNodes) == len(shard.members) {
			h.lggr.Infow("Successfully sent trigger request to all nodes in shard",
				"shard", shard.donID,
				"legacyExecutionID", legacyExecutionID,
				"executionIDWithTriggerIndex", executionIDWithTriggerIndex,
				"nodeCount", len(shard.members))
			return nil
		}

		// Not all nodes succeeded, wait and retry
		h.lggr.Debugw("Retrying failed nodes for trigger request",
			"shard", shard.donID,
			"legacyExecutionID", legacyExecutionID,
			"executionIDWithTriggerIndex", executionIDWithTriggerIndex,
			"failedCount", len(shard.members)-len(successfulNodes),
			"errors", combinedErr)

		select {
		case <-doneCh:
			h.lggr.Infow("Callback already responded to, stopping retries",
				"shard", shard.donID,
				"legacyExecutionID", legacyExecutionID,
				"executionIDWithTriggerIndex", executionIDWithTriggerIndex,
				"requestID", req.ID,
				"successNodes", len(successfulNodes),
				"totalNodes", len(shard.members))
			return nil
		case <-time.After(b.Duration()):
			continue
		case <-ctx.Done():
			return fmt.Errorf("shard %s: request retry time exceeded, some nodes may not have received the request: legacyExecutionID=%s, executionIDWithTriggerIndex=%s, successNodes=%d, totalNodes=%d",
				shard.donID, legacyExecutionID, executionIDWithTriggerIndex, len(successfulNodes), len(shard.members))
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
