package request

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	responseIDCount  map[[32]byte]int
	errorCount       map[string]int
	responseReceived map[p2ptypes.PeerID]bool
	lggr             logger.Logger

	requiredIdenticalResponses int

	requestTimeout time.Duration

	respSent bool
	mux      sync.Mutex
	wg       *sync.WaitGroup
}

func NewClientExecuteRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration) (*ClientRequest, error) {
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

	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request: %w", err)
	}

	lggr = lggr.With("requestId", requestID, "capabilityID", remoteCapabilityInfo.ID)
	return newClientRequest(ctx, lggr, requestID, remoteCapabilityInfo, localDonInfo, dispatcher, requestTimeout, tc, types.MethodExecute, rawRequest)
}

var (
	defaultDelayMargin = 10 * time.Second
)

func newClientRequest(ctx context.Context, lggr logger.Logger, requestID string, remoteCapabilityInfo commoncap.CapabilityInfo,
	localDonInfo commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration,
	tc transmission.TransmissionConfig, methodType string, rawRequest []byte) (*ClientRequest, error) {
	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	peerIDToTransmissionDelay, err := transmission.GetPeerIDToTransmissionDelaysForConfig(remoteCapabilityDonInfo.Members, requestID, tc)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
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
	effectiveTimeout := originalTimeout
	if originalTimeout < maxDelayDuration {
		effectiveTimeout = maxDelayDuration
	}

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
				CapabilityId:    remoteCapabilityInfo.ID,
				CapabilityDonId: remoteCapabilityDonInfo.ID,
				CallerDonId:     localDonInfo.ID,
				Method:          methodType,
				Payload:         rawRequest,
				MessageId:       []byte(requestID),
			}

			select {
			case <-innerCtx.Done():
				lggr.Debugw("context done, not sending request to peer", "peerID", peerID)
				return
			case <-time.After(delay):
				lggr.Debugw("sending request to peer", "peerID", peerID)
				err := dispatcher.Send(peerID, message)
				if err != nil {
					lggr.Errorw("failed to send message", "peerID", peerID, "error", err)
				}
			}
		}(ctxWithCancel, peerID, delay)
	}

	return &ClientRequest{
		id:                         requestID,
		cancelFn:                   cancelFn,
		createdAt:                  time.Now(),
		requestTimeout:             requestTimeout,
		requiredIdenticalResponses: int(remoteCapabilityDonInfo.F + 1),
		responseIDCount:            make(map[[32]byte]int),
		errorCount:                 make(map[string]int),
		responseReceived:           responseReceived,
		responseCh:                 make(chan clientResponse, 1),
		wg:                         &wg,
		lggr:                       lggr,
	}, nil
}

func (c *ClientRequest) ID() string {
	return c.id
}

func (c *ClientRequest) ResponseChan() <-chan clientResponse {
	return c.responseCh
}

func (c *ClientRequest) Expired() bool {
	return time.Since(c.createdAt) > c.requestTimeout
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

	if msg.Error == types.Error_OK {
		responseID := sha256.Sum256(msg.Payload)
		c.responseIDCount[responseID]++

		lggr := c.lggr.With("responseID", hex.EncodeToString(responseID[:]), "count", c.responseIDCount[responseID], "requiredCount", c.requiredIdenticalResponses, "peer", sender)
		lggr.Debug("received response from peer")

		if len(c.responseIDCount) > 1 {
			lggr.Warn("received multiple different responses for the same request, number of different responses received: %d", len(c.responseIDCount))
		}

		if c.responseIDCount[responseID] == c.requiredIdenticalResponses {
			c.sendResponse(clientResponse{Result: msg.Payload})
		}
	} else {
		c.lggr.Debug("received error from peer", "error", msg.Error, "errorMsg", msg.ErrorMsg, "peer", sender)
		c.errorCount[msg.ErrorMsg]++

		if len(c.errorCount) > 1 {
			c.lggr.Warn("received multiple different errors for the same request, number of different errors received", "errorCount", len(c.errorCount))
		}

		if c.errorCount[msg.ErrorMsg] == c.requiredIdenticalResponses {
			c.sendResponse(clientResponse{Err: fmt.Errorf("%s : %s", msg.Error, msg.ErrorMsg)})
		}
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
	c.lggr.Debugw("received OK response", "count", c.requiredIdenticalResponses)
}
