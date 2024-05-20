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
	capInfo              commoncap.CapabilityInfo
	donInfo              capabilities.DON
	dispatcher           types.Dispatcher
	lggr                 logger.Logger
	messageIDToWaitgroup sync.Map
	messageIDToResponse  sync.Map
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(lggr logger.Logger, capInfo commoncap.CapabilityInfo, donInfo capabilities.DON, dispatcher types.Dispatcher) *remoteTargetCaller {
	return &remoteTargetCaller{
		capInfo:    capInfo,
		donInfo:    donInfo,
		dispatcher: dispatcher,
		lggr:       lggr,
	}
}

func (c *remoteTargetCaller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.capInfo, nil
}

func (c *remoteTargetCaller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *remoteTargetCaller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {

	if c.capInfo.DON == nil {
		return nil, errors.New("missing remote capability DON info")
	}

	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, err
	}

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	peerIDToDelay, err := transmission.GetPeerIDToTransmissionDelay(c.capInfo.DON.Members, c.donInfo.Config.SharedSecret, req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID, tc)
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

	for peerID, delay := range peerIDToDelay {

		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				c.lggr.Debugw("executing delayed execution for peer", "peerID", peerID)
				m := &types.MessageBody{
					CapabilityId:    c.capInfo.ID,
					CapabilityDonId: c.capInfo.DON.ID,
					CallerDonId:     c.donInfo.ID,
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

		capabilityResponse, ok := response.(commoncap.CapabilityResponse)
		if !ok {
			return nil, fmt.Errorf("failed to cast response to CapabilityResponse: %v", response)
		}

		// TODO going to need to handle the case where the capability returns a stream of responses
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
	wg, loaded := c.messageIDToWaitgroup.LoadAndDelete(msg.MessageId)
	if loaded {
		wg.(*sync.WaitGroup).Done()
		c.messageIDToResponse.Store(msg.MessageId, msg.Payload)
		return
	}
}
