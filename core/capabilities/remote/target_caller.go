package remote

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// remoteTargetCaller/Receiver are shims translating between capability API calls and network messages
type remoteTargetCaller struct {
	remoteCapabilityInfo commoncap.CapabilityInfo
	localDONInfo         capabilities.DON
	dispatcher           types.Dispatcher
	lggr                 logger.Logger
	messageIDToWaitgroup sync.Map
	messageIDToResponse  sync.Map
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher) (*remoteTargetCaller, error) {

	if remoteCapabilityInfo.DON == nil {
		return nil, errors.New("missing remote capability DON info")
	}

	return &remoteTargetCaller{
		remoteCapabilityInfo: remoteCapabilityInfo,
		localDONInfo:         localDonInfo,
		dispatcher:           dispatcher,
		lggr:                 lggr,
	}, nil
}

func (c *remoteTargetCaller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.remoteCapabilityInfo, nil
}

func (c *remoteTargetCaller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) Execute(parentCtx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {

	// TODO should the transmission config be passed into the constructor rather than pulled from the request?
	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	peerIDToDelay, err := transmission.GetPeerIDToTransmissionDelay(c.remoteCapabilityInfo.DON.Members, c.localDONInfo.Config.SharedSecret, req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID, tc)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
	}

	messageID := uuid.New().String()

	responseWaitGroup := &sync.WaitGroup{}
	responseWaitGroup.Add(1)
	c.messageIDToWaitgroup.Store(messageID, responseWaitGroup)

	responseReceived := make(chan struct{})
	go func() {
		responseWaitGroup.Wait()
		close(responseReceived)
	}()

	// Once a response is returned from a remote capability any pending scheduled calls can be cancelled
	ctx, cancelFn := context.WithCancel(parentCtx)
	defer cancelFn()

	for peerID, delay := range peerIDToDelay {
		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				c.lggr.Debugw("executing delayed execution for peer", "peerID", peerID)
				m := &types.MessageBody{
					CapabilityId:    c.remoteCapabilityInfo.ID,
					CapabilityDonId: c.remoteCapabilityInfo.DON.ID,
					CallerDonId:     c.localDONInfo.ID,
					Method:          types.MethodExecute,
					Payload:         rawRequest,
					MessageId:       []byte(messageID),
				}
				err = c.dispatcher.Send(peerID, m)
				if err != nil {
					c.lggr.Errorw("failed to send message", "peerID", peerID, "err", err)
				}
			}
		}(peerID, delay)
	}

	select {
	case <-responseReceived:
		response, loaded := c.messageIDToResponse.LoadAndDelete(messageID)
		if !loaded {
			return nil, fmt.Errorf("no response found for message ID %s", messageID)
		}

		payload, ok := response.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected response type %T for message ID %s", response, messageID)
		}

		capabilityResponse, err := pb.UnmarshalCapabilityResponse(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal capability response: %w", err)
		}

		// TODO handle the case where the capability returns a stream of responses
		resultCh := make(chan commoncap.CapabilityResponse, 1)
		resultCh <- capabilityResponse
		close(resultCh)

		return resultCh, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {

	// TODO handle the case where the capability returns a stream of responses
	wg, loaded := c.messageIDToWaitgroup.LoadAndDelete(string(msg.MessageId))
	if loaded {
		wg.(*sync.WaitGroup).Done()
		c.messageIDToResponse.Store(string(msg.MessageId), msg.Payload)
		return
	}
}
