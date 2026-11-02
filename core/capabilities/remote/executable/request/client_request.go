package request

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type clientResponse struct {
	Result []byte
	Err    error
}

// messageFields are the fields of the outgoing MessageBody that the caller controls.
type messageFields struct {
	CapabilityID     string
	CapabilityDonID  uint32
	CallerDonID      uint32
	Method           string
	Payload          []byte
	MessageID        []byte
	CapabilityMethod string
}

type ClientRequest struct {
	id               string
	dispatcher       types.Dispatcher
	peers            []p2ptypes.PeerID
	message          messageFields
	responseCh       chan clientResponse
	createdAt        time.Time
	responseReceived map[p2ptypes.PeerID]bool
	lggr             logger.Logger

	requestTimeout time.Duration

	// aggregator decides the outcome of the request from peer replies
	// Default is an identical-response quorum aggregator
	aggregator ResponseAggregator

	respSent bool
	mux      sync.Mutex
}

// NewClientExecuteRequest creates a new client request for a remote executable capability.
// The request is not sent until Send is called.
func NewClientExecuteRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, capMethodName string,
	signers [][]byte, minResponsesToAggregate uint32,
) (*ClientRequest, error) {
	rawRequest, requestID, lggr, workflowExecutionID, err := prepareRequest(lggr, req)
	if err != nil {
		return nil, err
	}
	if remoteCapabilityInfo.DON == nil {
		return nil, errors.New("remote capability info missing DON")
	}
	msg := messageFields{
		CapabilityID:     remoteCapabilityInfo.ID,
		CapabilityDonID:  remoteCapabilityInfo.DON.ID,
		CallerDonID:      localDonInfo.ID,
		Method:           types.MethodExecute,
		Payload:          rawRequest,
		MessageID:        []byte(requestID),
		CapabilityMethod: capMethodName,
	}
	return newClientRequest(lggr, requestID, remoteCapabilityInfo, dispatcher, requestTimeout, msg,
		workflowExecutionID, req.Metadata.ReferenceID, signers, minResponsesToAggregate, nil)
}

// NewClientExecuteRequestWithAggregator is NewClientExecuteRequest with the outcome decided by aggregator.
func NewClientExecuteRequestWithAggregator(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, capMethodName string, aggregator ResponseAggregator,
) (*ClientRequest, error) {
	if aggregator == nil {
		return nil, errors.New("aggregator cannot be nil")
	}
	rawRequest, requestID, lggr, workflowExecutionID, err := prepareRequest(lggr, req)
	if err != nil {
		return nil, err
	}
	if remoteCapabilityInfo.DON == nil {
		return nil, errors.New("remote capability info missing DON")
	}
	msg := messageFields{
		CapabilityID:     remoteCapabilityInfo.ID,
		CapabilityDonID:  remoteCapabilityInfo.DON.ID,
		CallerDonID:      localDonInfo.ID,
		Method:           types.MethodExecute,
		Payload:          rawRequest,
		MessageID:        []byte(requestID),
		CapabilityMethod: capMethodName,
	}
	return newClientRequest(lggr, requestID, remoteCapabilityInfo, dispatcher, requestTimeout, msg,
		workflowExecutionID, req.Metadata.ReferenceID, nil, 0, aggregator)
}

// prepareRequest marshals the request and derives its request ID and enriched logger.
func prepareRequest(lggr logger.Logger, req commoncap.CapabilityRequest) (rawRequest []byte, requestID string, outLggr logger.Logger, workflowExecutionID string, err error) {
	rawRequest, err = proto.MarshalOptions{Deterministic: true}.Marshal(pb.CapabilityRequestToProto(req))
	if err != nil {
		return nil, "", nil, "", fmt.Errorf("failed to marshal capability request: %w", err)
	}

	workflowExecutionID = req.Metadata.WorkflowExecutionID
	if err = validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
		return nil, "", nil, "", fmt.Errorf("workflow execution ID is invalid: %w", err)
	}

	// the requestID must be delineated by the workflow execution ID and the reference ID
	// to ensure that it supports parallel step execution
	requestID = types.MethodExecute + ":" + workflowExecutionID + ":" + req.Metadata.ReferenceID

	outLggr = logger.With(lggr, "requestId", requestID) // cap ID and method name included in the parent logger
	return rawRequest, requestID, outLggr, workflowExecutionID, nil
}

var (
	// Extra time the workflow DON client waits beyond requestTimeout for P2P responses
	// after remote capability nodes hit their execution deadline.
	defaultResponseAggregationGrace = 10 * time.Second
)

func newClientRequest(lggr logger.Logger, requestID string, remoteCapabilityInfo commoncap.CapabilityInfo,
	dispatcher types.Dispatcher, requestTimeout time.Duration, message messageFields,
	workflowExecutionID, stepRef string, signers [][]byte, minResponsesToAggregate uint32,
	aggregator ResponseAggregator,
) (*ClientRequest, error) {
	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	responseReceived := make(map[p2ptypes.PeerID]bool, len(remoteCapabilityDonInfo.Members))
	for _, peerID := range remoteCapabilityDonInfo.Members {
		responseReceived[peerID] = false
	}

	if aggregator == nil {
		aggregator = newQuorumAggregator(lggr, signers, workflowExecutionID, stepRef,
			len(remoteCapabilityDonInfo.Members), remoteCapabilityDonInfo.F, minResponsesToAggregate)
	}

	return &ClientRequest{
		id:               requestID,
		dispatcher:       dispatcher,
		peers:            remoteCapabilityDonInfo.Members,
		message:          message,
		createdAt:        time.Now(),
		requestTimeout:   requestTimeout,
		responseReceived: responseReceived,
		responseCh:       make(chan clientResponse, 1),
		lggr:             lggr,
		aggregator:       aggregator,
	}, nil
}

// Send dispatches the request to all remote capability peers. Callers must register the
// request to receive responses before calling Send, since peers may reply immediately.
func (c *ClientRequest) Send() {
	for _, peerID := range c.peers {
		// the dispatcher mutates the message, so each peer gets its own
		message := &types.MessageBody{
			CapabilityId:     c.message.CapabilityID,
			CapabilityDonId:  c.message.CapabilityDonID,
			CallerDonId:      c.message.CallerDonID,
			Method:           c.message.Method,
			Payload:          c.message.Payload,
			MessageId:        c.message.MessageID,
			CapabilityMethod: c.message.CapabilityMethod,
		}
		if err := c.dispatcher.Send(peerID, message); err != nil {
			c.lggr.Errorw("failed to send message", "peerID", peerID, "error", err)
		}
	}
}

func (c *ClientRequest) ID() string {
	return c.id
}

func (c *ClientRequest) ResponseChan() <-chan clientResponse {
	return c.responseCh
}

func (c *ClientRequest) Expired() bool {
	return time.Since(c.createdAt) > c.requestTimeout+defaultResponseAggregationGrace
}

func (c *ClientRequest) Cancel(err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if !c.respSent {
		c.sendResponse(clientResponse{Err: err})
	}
}

func (c *ClientRequest) OnMessage(_ context.Context, msg *types.MessageBody) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.respSent {
		return nil
	}

	if msg.Sender == nil {
		return errors.New("sender missing from message")
	}

	c.lggr.Debugw("OnMessage called for client request")

	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return fmt.Errorf("failed to convert message sender to PeerID: %w", err)
	}

	received, expected := c.responseReceived[sender]
	if !expected {
		return fmt.Errorf("response from peer %s not expected", sender)
	}

	if received {
		return fmt.Errorf("response from peer %s already received", sender)
	}

	c.responseReceived[sender] = true

	return c.onMessageWithAggregator(sender, msg)
}

// Must be called with c.mux held.
func (c *ClientRequest) onMessageWithAggregator(sender p2ptypes.PeerID, msg *types.MessageBody) error {
	var (
		final *commoncap.CapabilityResponse
		err   error
	)
	if msg.Error == types.Error_OK {
		resp, uerr := pb.UnmarshalCapabilityResponse(msg.Payload)
		if uerr != nil {
			final, err = c.aggregator.OnUnmarshalable(sender, uerr)
		} else {
			final, err = c.aggregator.OnResponse(sender, resp)
		}
	} else {
		c.lggr.Debugw("received error from peer", "error", msg.Error, "errorMsg", msg.ErrorMsg, "peer", sender)
		final, err = c.aggregator.OnError(sender, msg.Error, msg.ErrorMsg)
	}

	// observation errors are returned to the OnMessage caller without
	// deciding the request
	if obs, ok := errors.AsType[observationError](err); ok {
		return obs.err
	}

	switch {
	case err != nil:
		c.sendResponse(clientResponse{Err: err})
	case final != nil:
		payload, merr := pb.MarshalCapabilityResponse(*final)
		if merr != nil {
			c.sendResponse(clientResponse{Err: fmt.Errorf("failed to marshal aggregated response: %w", merr)})
			return nil
		}
		c.sendResponse(clientResponse{Result: payload})
	}
	return nil
}

func (c *ClientRequest) sendResponse(response clientResponse) {
	c.responseCh <- response
	close(c.responseCh)
	c.respSent = true
	if response.Err != nil {
		c.lggr.Warnw("received error response", "error", remote.SanitizeLogString(response.Err.Error()))
		return
	}
	c.lggr.Debugw("received OK response")
}
