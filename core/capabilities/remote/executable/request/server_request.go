package request

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"slices"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
		metric.WithAttributes(attribute.String("callingDON", s.callingDonID), attribute.String("capabilityID", s.capabilityID), attribute.String("status", status), attribute.String("dispatcherErr", dv)))
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

	mux  sync.Mutex
	lggr logger.Logger

	metrics *srMetrics
}

func NewServerRequest(capability capabilities.ExecutableCapability, method string, capabilityID string, capabilityDonID uint32,
	capabilityPeerID p2ptypes.PeerID,
	callingDon capabilities.DON, requestID string,
	dispatcher types.Dispatcher, requestTimeout time.Duration, lggr logger.Logger) (*ServerRequest, error) {
	lggr = lggr.Named("ServerRequest").With("requestID", requestID, "capabilityID", capabilityID)

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
		lggr:                    lggr,
		metrics:                 m,
	}, nil
}

func (e *ServerRequest) OnMessage(ctx context.Context, msg *types.MessageBody) error {
	e.metrics.countExecutionRequest(ctx)

	e.mux.Lock()
	defer e.mux.Unlock()

	if msg.Sender == nil {
		return errors.New("sender missing from message")
	}

	requester, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return fmt.Errorf("failed to convert message sender to PeerID: %w", err)
	}

	if err := e.addRequester(requester); err != nil {
		return fmt.Errorf("failed to add requester to request: %w", err)
	}

	e.lggr.Debugw("OnMessage called for request", "calls", len(e.requesters),
		"hasResponse", e.response != nil, "requester", requester.String(), "minRequsters", e.callingDon.F+1)

	if e.minimumRequiredRequestsReceived() && !e.hasResponse() {
		switch e.method {
		case types.MethodExecute:
			e.executeRequest(ctx, msg, executeCapabilityRequest)
		default:
			e.setError(types.Error_INTERNAL_ERROR, "unknown method %s"+e.method)
		}
	}

	if err := e.sendResponses(ctx); err != nil {
		return fmt.Errorf("failed to send responses: %w", err)
	}

	return nil
}

func (e *ServerRequest) Expired() bool {
	return time.Since(e.createdTime) > e.requestTimeout
}

func (e *ServerRequest) Cancel(ctx context.Context, err types.Error, msg string) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if !e.hasResponse() {
		e.setError(err, msg)
		if err := e.sendResponses(ctx); err != nil {
			return fmt.Errorf("failed to send responses: %w", err)
		}
	}

	return nil
}

type executeFn func(ctx context.Context, lggr logger.Logger, capability capabilities.ExecutableCapability, payload []byte) ([]byte, error)

func (e *ServerRequest) executeRequest(ctx context.Context, msg *types.MessageBody, method executeFn) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, e.requestTimeout)
	defer cancel()

	success := false
	start := time.Now()
	responsePayload, err := method(ctxWithTimeout, e.lggr, e.capability, msg.Payload)
	if err != nil {
		e.setError(types.Error_INTERNAL_ERROR, err.Error())
	} else {
		success = true
		e.setResult(responsePayload)
	}

	e.metrics.countExecution(ctx, success)
	e.metrics.recordExecutionDuration(ctx, time.Since(start), success)
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
		CapabilityId:    e.capabilityID,
		CapabilityDonId: e.capabilityDonID,
		CallerDonId:     e.callingDon.ID,
		Method:          types.MethodExecute,
		MessageId:       []byte(e.requestMessageID),
		Sender:          e.capabilityPeerID[:],
		Receiver:        requester[:],
	}

	if e.response.error != types.Error_OK {
		responseMsg.Error = e.response.error
		responseMsg.ErrorMsg = e.response.errorMsg
	} else {
		responseMsg.Payload = e.response.response
	}

	e.lggr.Debugw("Sending response", "receiver", requester)
	err := e.dispatcher.Send(requester, &responseMsg)
	e.metrics.countExecutionResponse(ctx, e.response.error.String(), err != nil)
	if err != nil {
		return fmt.Errorf("failed to send response to dispatcher: %w", err)
	}

	e.responseSentToRequester[requester] = true

	return nil
}

func executeCapabilityRequest(ctx context.Context, lggr logger.Logger, capability capabilities.ExecutableCapability, payload []byte) ([]byte, error) {
	capabilityRequest, err := pb.UnmarshalCapabilityRequest(payload)
	if err != nil {
		lggr.Errorw("failed to unmarshal capability request", "err", err)

		// Do not include the unmarshal error in the response as it may contain sensitive information
		return nil, errors.New("failed to unmarshal capability request")
	}

	lggr = lggr.With("metadata", capabilityRequest.Metadata)

	lggr.Debugw("executing capability")
	capResponse, err := capability.Execute(ctx, capabilityRequest)
	if err != nil {
		lggr.Errorw("received execution error", "error", err)

		// If an error is a RemoteReportableError then the information it contains is considered to be safe to report back to the calling nodes
		var reportableError *capabilities.RemoteReportableError
		if errors.As(err, &reportableError) {
			return nil, fmt.Errorf("failed to execute capability: %w", reportableError)
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
