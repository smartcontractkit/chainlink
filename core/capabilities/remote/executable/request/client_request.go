package request

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/durableemitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type clientResponse struct {
	Result []byte
	Err    error
}

type ClientRequest struct {
	id               string
	cancelFn         context.CancelFunc
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
	wg       *sync.WaitGroup
}

// NewClientExecuteRequest creates a new client request for a remote executable capability.
// V2 capabilities always use the AllAtOnce transmission schedule.
func NewClientExecuteRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, capMethodName string,
	signers [][]byte, minResponsesToAggregate uint32,
) (*ClientRequest, error) {
	rawRequest, err := proto.MarshalOptions{Deterministic: true}.Marshal(pb.CapabilityRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	workflowExecutionID := req.Metadata.WorkflowExecutionID
	if err = validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
		return nil, fmt.Errorf("workflow execution ID is invalid: %w", err)
	}

	// the requestID must be delineated by the workflow execution ID and the reference ID
	// to ensure that it supports parallel step execution
	requestID := types.MethodExecute + ":" + workflowExecutionID + ":" + req.Metadata.ReferenceID

	lggr = logger.With(lggr, "requestId", requestID) // cap ID and method name included in the parent logger
	return newClientRequest(ctx, lggr, requestID, remoteCapabilityInfo, localDonInfo, dispatcher, requestTimeout, types.MethodExecute, rawRequest, workflowExecutionID, req.Metadata.ReferenceID, capMethodName, signers, minResponsesToAggregate)
}

// NewClientExecuteRequestWithAggregator is NewClientExecuteRequest with the outcome decided by aggregator.
func NewClientExecuteRequestWithAggregator(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, capMethodName string, aggregator ResponseAggregator,
) (*ClientRequest, error) {
	if aggregator == nil {
		return nil, errors.New("aggregator cannot be nil")
	}
	cr, err := NewClientExecuteRequest(ctx, lggr, req, remoteCapabilityInfo, localDonInfo, dispatcher, requestTimeout, capMethodName, nil, 0)
	if err != nil {
		return nil, err
	}
	cr.aggregator = aggregator
	return cr, nil
}

var (
	defaultDelayMargin = 10 * time.Second
	// Extra time the workflow DON client waits beyond requestTimeout for P2P responses
	// after remote capability nodes hit their execution deadline.
	defaultResponseAggregationGrace = 10 * time.Second
)

func newClientRequest(ctx context.Context, lggr logger.Logger, requestID string, remoteCapabilityInfo commoncap.CapabilityInfo,
	localDonInfo commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration,
	methodType string, rawRequest []byte, workflowExecutionID string, stepRef string, capMethodName string,
	signers [][]byte, minResponsesToAggregate uint32,
) (*ClientRequest, error) {
	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	peerIDToTransmissionDelay := getPeerIDToTransmissionDelays(remoteCapabilityDonInfo.Members)

	// send schedule through beholder for single execution performance tracking
	if err := emitTransmissionScheduleEvent(ctx,
		workflowExecutionID,
		requestID,
		remoteCapabilityInfo.ID,
		stepRef,
		peerIDToTransmissionDelay,
	); err != nil {
		lggr.Errorw("failed to emit transmission schedule event", "error", err)
	}

	responseReceived := make(map[p2ptypes.PeerID]bool)

	maxDelayDuration := time.Duration(0)
	for _, delay := range peerIDToTransmissionDelay {
		if delay > maxDelayDuration {
			maxDelayDuration = delay
		}
	}

	// Add some margin to allow the last peer to respond
	maxDelayDuration += defaultDelayMargin

	// Instantiate a new context based on the parent, but without its deadline.
	// We set a new deadline instead equal to the original timeout OR the full length
	// of the execution schedule plus some margin, whichever is greater

	// We do this to ensure that we will always execute the entire transmission schedule.
	// This ensures that all capability DON nodes will receive a quorum of requests,
	// and will execute all requests they receive from the workflow DON, preventing
	// quorum errors from lagging members of the workflow DON.
	dl, ok := ctx.Deadline()
	originalTimeout := time.Duration(0)
	if ok {
		originalTimeout = time.Until(dl)
	}
	effectiveTimeout := max(originalTimeout, maxDelayDuration)

	// Now let's create a new context based on the adjusted timeout value.
	// By calling WithoutCancel, we ensure that this context can only be cancelled in
	// one of two ways -- 1) by explicitly calling the cancelFn we create below, or 2)
	// after the adjusted timeout expires.
	ctxWithoutCancel := context.WithoutCancel(ctx)
	ctxWithCancel, cancelFn := context.WithTimeout(ctxWithoutCancel, effectiveTimeout)

	lggr.Debugw("sending request to peers", "schedule", peerIDToTransmissionDelay, "originalTimeout", originalTimeout, "effectiveTimeout", effectiveTimeout)

	var wg sync.WaitGroup
	for peerID, delay := range peerIDToTransmissionDelay {
		responseReceived[peerID] = false

		wg.Add(1)
		go func(innerCtx context.Context, peerID ragep2ptypes.PeerID, delay time.Duration) {
			defer wg.Done()
			message := &types.MessageBody{
				CapabilityId:     remoteCapabilityInfo.ID,
				CapabilityDonId:  remoteCapabilityDonInfo.ID,
				CallerDonId:      localDonInfo.ID,
				Method:           methodType,
				Payload:          rawRequest,
				MessageId:        []byte(requestID),
				CapabilityMethod: capMethodName,
			}

			timer := time.NewTimer(delay)
			defer timer.Stop()

			select {
			case <-innerCtx.Done():
				lggr.Debugw("context done, not sending request to peer", "peerID", peerID)
				return
			case <-timer.C:
				lggr.Debugw("sending request to peer", "peerID", peerID)
				err := dispatcher.Send(peerID, message)
				if err != nil {
					lggr.Errorw("failed to send message", "peerID", peerID, "error", err)
				}
			}
		}(ctxWithCancel, peerID, delay)
	}

	return &ClientRequest{
		id:               requestID,
		cancelFn:         cancelFn,
		createdAt:        time.Now(),
		requestTimeout:   requestTimeout,
		responseReceived: responseReceived,
		responseCh:       make(chan clientResponse, 1),
		wg:               &wg,
		lggr:             lggr,
		aggregator: newQuorumAggregator(lggr, signers, workflowExecutionID, stepRef,
			len(remoteCapabilityDonInfo.Members), remoteCapabilityDonInfo.F, minResponsesToAggregate),
	}, nil
}

// getPeerIDToTransmissionDelays returns a map of PeerID to the time.Duration that the node
// with that PeerID should wait before transmitting the capability request. If a node is not
// in the map, it should not transmit. With the AllAtOnce schedule every DON member transmits
// immediately.
func getPeerIDToTransmissionDelays(donPeerIDs []p2ptypes.PeerID) map[p2ptypes.PeerID]time.Duration {
	peerIDToTransmissionDelay := make(map[p2ptypes.PeerID]time.Duration, len(donPeerIDs))
	for _, peerID := range donPeerIDs {
		peerIDToTransmissionDelay[peerID] = 0
	}
	return peerIDToTransmissionDelay
}

func emitTransmissionScheduleEvent(ctx context.Context, workflowExecutionID, transmissionID, capabilityID, stepRef string, peerIDToTransmissionDelay map[p2ptypes.PeerID]time.Duration) error {
	// map of peer IDs to their transmission delays in milliseconds
	peerDelaysMap := make(map[string]int64, len(peerIDToTransmissionDelay))
	for peerID, delay := range peerIDToTransmissionDelay {
		peerDelaysMap[peerID.String()] = delay.Milliseconds()
	}

	msg := &events.TransmissionsScheduledEvent{
		Timestamp:              time.Now().Format(time.RFC3339),
		ScheduleType:           ScheduleAllAtOnce,
		WorkflowExecutionID:    workflowExecutionID,
		TransmissionID:         transmissionID,
		CapabilityID:           capabilityID,
		StepRef:                stepRef,
		PeerTransmissionDelays: peerDelaysMap,
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal TransmissionScheduleEvent: %w", err)
	}

	// emit transmission schedule event to track which nodes are successful when called to emit
	entity := fmt.Sprintf("%s.%s", TransmissionEventProtoPkg, TransmissionEventEntity)
	if err = beholder.GetEmitter().Emit(ctx, b,
		"beholder_data_schema", TransmissionEventSchema, // required
		"beholder_domain", "platform", // required
		"beholder_entity", entity); err != nil { // required
		return err
	}

	err = durableemitter.GlobalEmit(ctx, b, "source", "platform", "type", entity)
	if err != nil && !errors.Is(err, durableemitter.ErrNotInitialized) {
		return err
	}
	return nil
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
	c.cancelFn()
	c.wg.Wait()
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
	var obs observationError
	if errors.As(err, &obs) {
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
