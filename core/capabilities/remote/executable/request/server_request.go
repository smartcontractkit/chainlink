package request

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type srMetrics struct {
	capabilityID         string
	callingDonID         string
	executeDuration      metric.Int64Histogram
	executeCount         metric.Int64Counter
	executeRequestCount  metric.Int64Counter
	executeResponseCount metric.Int64Counter
}

func (s *srMetrics) recordExecutionDuration(ctx context.Context, d time.Duration, success bool) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	s.executeDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", successStr), attribute.String("callingDON", s.callingDonID), attribute.String("capabilityID", s.capabilityID),
	))
}

func (s *srMetrics) countExecution(ctx context.Context, success bool) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	s.executeCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String("success", successStr), attribute.String("callingDON", s.callingDonID), attribute.String("capabilityID", s.capabilityID),
	))
}

func (s *srMetrics) countExecutionRequest(ctx context.Context) {
	s.executeRequestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String("callingDON", s.callingDonID), attribute.String("capabilityID", s.capabilityID),
	))
}

func (s *srMetrics) countExecutionResponse(ctx context.Context, status string, dispatcherErr bool) {
	// Beholder doesn't support non-string attributes
	dv := "false"
	if dispatcherErr {
		dv = "true"
	}
	s.executeResponseCount.Add(
		ctx, 1,
		metric.WithAttributes(attribute.String("callingDON", s.callingDonID), attribute.String("capabilityID", s.capabilityID), attribute.String("status", status), attribute.String("dispatcherErr", dv)),
	)
}

func newSrMetrics(capabilityID string, callingDonID uint32) (*srMetrics, error) {
	h, err := beholder.GetMeter().Int64Histogram("platform_executable_capability_server_execute_duration_ms")
	if err != nil {
		return nil, err
	}

	ec, err := beholder.GetMeter().Int64Counter("platform_executable_capability_server_execute_count")
	if err != nil {
		return nil, err
	}

	erc, err := beholder.GetMeter().Int64Counter("platform_executable_capability_server_execute_request_count")
	if err != nil {
		return nil, err
	}

	erspc, err := beholder.GetMeter().Int64Counter("platform_executable_capability_server_execute_response_count")
	if err != nil {
		return nil, err
	}

	return &srMetrics{
		capabilityID:         capabilityID,
		callingDonID:         strconv.FormatUint(uint64(callingDonID), 10),
		executeDuration:      h,
		executeCount:         ec,
		executeRequestCount:  erc,
		executeResponseCount: erspc,
	}, nil
}

type response struct {
	response []byte
	error    types.Error
	errorMsg string
}

type ServerRequest struct {
	capability capabilities.ExecutableCapability

	capabilityPeerID p2ptypes.PeerID
	capabilityID     string
	capabilityDonID  uint32

	dispatcher types.Dispatcher

	requesters              map[p2ptypes.PeerID]bool
	responseSentToRequester map[p2ptypes.PeerID]bool

	createdTime time.Time

	response *response

	callingDon capabilities.DON

	requestMessageID string
	method           string
	requestTimeout   time.Duration
	capMethodName    string

	// workflowDONBindingGate, when open, requires the request's
	// Metadata.WorkflowDonID to match the authenticated calling DON.
	workflowDONBindingGate limits.GateLimiter

	// stateMux guards requesters, responseSentToRequester, response, and executionCancel.
	// It is held only for short map/field operations, never during capability execution.
	stateMux sync.Mutex

	// executionClaimed is set to true exactly once (via CAS) by the message that
	// wins the right to execute the capability. All other messages skip execution.
	executionClaimed atomic.Bool

	// executionCancel cancels the in-flight capability execution context.
	// Set under stateMux when execution starts; called by Cancel to stop
	// execution early (e.g. on request expiry) instead of waiting for completion.
	executionCancel context.CancelFunc

	lggr logger.Logger

	metrics *srMetrics
}

func NewServerRequest(capability capabilities.ExecutableCapability, method, capabilityID string, capabilityDonID uint32,
	capabilityPeerID p2ptypes.PeerID,
	callingDon capabilities.DON, requestID string,
	dispatcher types.Dispatcher, requestTimeout time.Duration, capMethodName string,
	workflowDONBindingGate limits.GateLimiter, lggr logger.Logger,
) (*ServerRequest, error) {
	lggr = logger.With(logger.Named(lggr, "ServerRequest"), "requestID", requestID) // cap ID and method name included in the parent logger

	if workflowDONBindingGate == nil {
		return nil, errors.New("workflowDONBindingGate must not be nil")
	}

	m, err := newSrMetrics(capabilityID, callingDon.ID)
	if err != nil {
		return nil, err
	}

	return &ServerRequest{
		capability:              capability,
		createdTime:             time.Now(),
		capabilityID:            capabilityID,
		capabilityDonID:         capabilityDonID,
		capabilityPeerID:        capabilityPeerID,
		dispatcher:              dispatcher,
		requesters:              map[p2ptypes.PeerID]bool{},
		responseSentToRequester: map[p2ptypes.PeerID]bool{},
		callingDon:              callingDon,
		requestMessageID:        requestID,
		method:                  method,
		requestTimeout:          requestTimeout,
		capMethodName:           capMethodName,
		workflowDONBindingGate:  workflowDONBindingGate,
		lggr:                    lggr,
		metrics:                 m,
	}, nil
}

func (e *ServerRequest) OnMessage(ctx context.Context, msg *types.MessageBody) error {
	e.metrics.countExecutionRequest(ctx)

	if msg.Sender == nil {
		return errors.New("sender missing from message")
	}

	requester, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return fmt.Errorf("failed to convert message sender to PeerID: %w", err)
	}

	e.stateMux.Lock()
	if err := e.addRequester(requester); err != nil {
		e.stateMux.Unlock()
		return fmt.Errorf("failed to add requester to request: %w", err)
	}

	quorumReached := e.minimumRequiredRequestsReceived()
	calls := len(e.requesters)
	hasResponse := e.hasResponse()
	e.stateMux.Unlock()

	e.lggr.Debugw("OnMessage called for request", "requester", requester.String(),
		"quorumReached", quorumReached, "hasResponse", hasResponse, "minRequesters", e.callingDon.F+1, "calls", calls)

	// Only one message wins the right to execute. All others skip execution
	// and either wait for the executor's fan-out or self-send if the response
	// is already available.
	if quorumReached && !hasResponse && e.executionClaimed.CompareAndSwap(false, true) {
		switch e.method {
		case types.MethodExecute:
			ctxWithTimeout, cancel := context.WithTimeout(ctx, e.requestTimeout)
			defer cancel()

			// Expose the cancel func so Cancel can stop the in-flight execution early.
			e.stateMux.Lock()
			e.executionCancel = cancel
			e.stateMux.Unlock()
			success := false
			start := time.Now()
			responsePayload, responseErr := executeCapabilityRequest(ctxWithTimeout, e.lggr, e.capability, msg.Payload, e.callingDon.ID, e.workflowDONBindingGate)

			e.stateMux.Lock()
			// Cancel may have already set a timeout error response; never overwrite an
			// existing response.
			if !e.hasResponse() {
				if responseErr != nil {
					e.setError(types.Error_INTERNAL_ERROR, responseErr.Error())
				} else {
					success = true
					e.setResult(responsePayload)
				}
			}
			e.stateMux.Unlock()

			e.metrics.countExecution(ctxWithTimeout, success)
			e.metrics.recordExecutionDuration(ctxWithTimeout, time.Since(start), success)
		default:
			e.stateMux.Lock()
			e.setError(types.Error_INTERNAL_ERROR, "unknown method %s"+e.method)
			e.stateMux.Unlock()
		}
	}

	e.stateMux.Lock()
	defer e.stateMux.Unlock()
	if err := e.sendResponses(ctx); err != nil {
		return fmt.Errorf("failed to send responses: %w", err)
	}

	return nil
}

func (e *ServerRequest) Expired() bool {
	return time.Since(e.createdTime) > e.requestTimeout
}

func (e *ServerRequest) Evictable(minRetention time.Duration) bool {
	age := time.Since(e.createdTime)
	return age > e.requestTimeout && age > minRetention
}

// Cancel stops any in-flight execution by cancelling its context and, if no
// response has been produced yet, records err as the response and fans it out
// to all requesters.
func (e *ServerRequest) Cancel(ctx context.Context, err types.Error, msg string) error {
	e.stateMux.Lock()
	defer e.stateMux.Unlock()

	// Cancel the in-flight execution, if any. The executor goroutine returns
	// early on the cancelled context and skips overwriting the response set
	// below (guarded by hasResponse).
	if e.executionCancel != nil {
		e.executionCancel()
	}

	// Only set cancellation error if no response exists (matches original behavior)
	if !e.hasResponse() {
		e.setError(err, msg)
		if err := e.sendResponses(ctx); err != nil {
			return fmt.Errorf("failed to send responses: %w", err)
		}
	}

	return nil
}

func (e *ServerRequest) addRequester(from p2ptypes.PeerID) error {
	fromPeerInCallingDon := slices.Contains(e.callingDon.Members, from)

	if !fromPeerInCallingDon {
		return fmt.Errorf("request received from peer %s not in calling don", from)
	}

	if e.requesters[from] {
		return fmt.Errorf("request already received from peer %s", from)
	}

	e.requesters[from] = true

	return nil
}

func (e *ServerRequest) minimumRequiredRequestsReceived() bool {
	return len(e.requesters) >= int(e.callingDon.F+1)
}

func (e *ServerRequest) setResult(result []byte) {
	e.lggr.Debug("setting result on request")
	e.response = &response{
		response: result,
	}
}

func (e *ServerRequest) setError(err types.Error, errMsg string) {
	e.lggr.Debugw("setting error on request", "type", err, "error", errMsg)
	e.response = &response{
		error:    err,
		errorMsg: errMsg,
	}
}

func (e *ServerRequest) hasResponse() bool {
	return e.response != nil
}

func (e *ServerRequest) sendResponses(ctx context.Context) error {
	if e.hasResponse() {
		for requester := range e.requesters {
			if !e.responseSentToRequester[requester] {
				e.responseSentToRequester[requester] = true
				if err := e.sendResponse(ctx, requester); err != nil {
					return fmt.Errorf("failed to send response to requester %s: %w", requester, err)
				}
			}
		}
	}

	return nil
}

func (e *ServerRequest) sendResponse(ctx context.Context, requester p2ptypes.PeerID) error {
	responseMsg := types.MessageBody{
		CapabilityId:     e.capabilityID,
		CapabilityDonId:  e.capabilityDonID,
		CallerDonId:      e.callingDon.ID,
		Method:           types.MethodExecute,
		MessageId:        []byte(e.requestMessageID),
		Sender:           e.capabilityPeerID[:],
		Receiver:         requester[:],
		CapabilityMethod: e.capMethodName,
	}

	if e.response.error != types.Error_OK {
		responseMsg.Error = e.response.error
		responseMsg.ErrorMsg = e.response.errorMsg
	} else {
		responseMsg.Payload = e.response.response
	}

	e.lggr.Debugw("Sending response", "receiver", requester, "capabilityId", e.capabilityID, "donId", e.capabilityDonID, "method", e.capMethodName)
	err := e.dispatcher.Send(requester, &responseMsg)
	e.metrics.countExecutionResponse(ctx, e.response.error.String(), err != nil)
	if err != nil {
		return fmt.Errorf("failed to send response to dispatcher: %w", err)
	}

	e.responseSentToRequester[requester] = true

	return nil
}

func executeCapabilityRequest(ctx context.Context, lggr logger.Logger, capability capabilities.ExecutableCapability, payload []byte, callingDonID uint32, workflowDONBindingGate limits.GateLimiter) ([]byte, error) {
	capabilityRequest, err := pb.UnmarshalCapabilityRequest(payload)
	if err != nil {
		lggr.Errorw("failed to unmarshal capability request", "err", err)

		// Do not include the unmarshal error in the response as it may contain sensitive information
		return nil, errors.New("failed to unmarshal capability request")
	}

	// When enabled, bind the caller-supplied WorkflowDonID to the authenticated
	// calling DON so it cannot be spoofed. All F+1 aggregated requests share this
	// payload (WorkflowDonID is part of the request hash), so a single check here
	// covers the quorum. The gate is guaranteed non-nil by NewServerRequest.
	enabled, gerr := workflowDONBindingGate.Limit(ctx)
	if gerr != nil {
		lggr.Errorw("failed to evaluate workflow DON binding gate", "err", gerr)
		return nil, errors.New("failed to evaluate workflow DON binding gate")
	}
	if enabled && capabilityRequest.Metadata.WorkflowDonID != callingDonID {
		lggr.Errorw("workflow DON ID in request metadata does not match calling DON",
			"metadataWorkflowDonID", capabilityRequest.Metadata.WorkflowDonID, "callingDonID", callingDonID)
		return nil, fmt.Errorf("workflow DON ID %d in request metadata does not match calling DON ID %d",
			capabilityRequest.Metadata.WorkflowDonID, callingDonID)
	}

	lggr = logger.With(lggr, "metadata", capabilityRequest.Metadata)

	lggr.Debugw("executing capability")
	capResponse, err := capability.Execute(ctx, capabilityRequest)
	if err != nil {
		lggr.Errorw("received execution error", "error", err)

		if capError, ok := errors.AsType[caperrors.Error](err); ok {
			return nil, errors.New(capError.SerializeToRemoteString())
		}
		return nil, errors.New("failed to execute capability")
	}

	responsePayload, err := pb.MarshalCapabilityResponse(capResponse)
	if err != nil {
		lggr.Errorw("failed to marshal capability request", "error", err)

		// Do not include the marshal error in the response as it may contain sensitive information
		return nil, errors.New("failed to marshal capability request")
	}

	lggr.Debug("received execution results")
	return responsePayload, nil
}
