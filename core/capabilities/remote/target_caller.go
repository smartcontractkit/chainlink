package remote

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

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
	remoteCapabilityInfo    commoncap.CapabilityInfo
	remoteCapabilityDonInfo capabilities.DON
	localDONInfo            capabilities.DON
	dispatcher              types.Dispatcher
	lggr                    logger.Logger
	messageIDToWaitgroup    sync.Map
	messageIDToResponse     sync.Map
}

var _ commoncap.TargetCapability = &remoteTargetCaller{}
var _ types.Receiver = &remoteTargetCaller{}

func NewRemoteTargetCaller(lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, remoteCapabilityDonInfo capabilities.DON, localDonInfo capabilities.DON, dispatcher types.Dispatcher) (*remoteTargetCaller, error) {

	return &remoteTargetCaller{
		remoteCapabilityInfo:    remoteCapabilityInfo,
		remoteCapabilityDonInfo: remoteCapabilityDonInfo,
		localDONInfo:            localDonInfo,
		dispatcher:              dispatcher,
		lggr:                    lggr,
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

	// TODO Assuming here that the capability request is deterministically unique across the nodes, need to confirm this is reasonable assumption
	// TODO also check pb marshalliing is by default deterministic in the version being used

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	messageID := sha256.Sum256(rawRequest)

	responseWaitGroup := &sync.WaitGroup{}
	responseWaitGroup.Add(1)
	c.messageIDToWaitgroup.Store(messageID, responseWaitGroup)

	responseReceived := make(chan struct{})
	go func() {
		responseWaitGroup.Wait()
		close(responseReceived)
	}()

	// Once a response is received from a remote capability further transmission should be cancelled
	ctx, cancelFn := context.WithCancel(parentCtx)
	defer cancelFn()

	if err := c.transmitRequestWithMessageID(ctx, req, messageID); err != nil {
		return nil, fmt.Errorf("failed to transmit request: %w", err)
	}

	select {
	case <-responseReceived:

		response, loaded := c.messageIDToResponse.LoadAndDelete(messageID)
		if !loaded {
			return nil, fmt.Errorf("no response found for message ID %s", messageID)
		}

		msg, ok := response.(*types.MessageBody)
		if !ok {
			return nil, fmt.Errorf("unexpected response type %T for message ID %s", response, messageID)
		}

		if msg.Error != types.Error_OK {
			return nil, fmt.Errorf("remote capability returned error: %s", msg.Error)
		}

		capabilityResponse, err := pb.UnmarshalCapabilityResponse(msg.Payload)
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

// transmitRequestWithMessageID transmits a capability request to remote capabilities according to the transmission configuration
func (c *remoteTargetCaller) transmitRequestWithMessageID(ctx context.Context, req commoncap.CapabilityRequest, messageID [32]byte) error {
	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return fmt.Errorf("failed to marshal capability request: %w", err)
	}

	// TODO should the transmission config be passed into the constructor rather than pulled from the request?
	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	message := &types.MessageBody{
		CapabilityId:    c.remoteCapabilityInfo.ID,
		CapabilityDonId: c.remoteCapabilityDonInfo.ID,
		CallerDonId:     c.localDONInfo.ID,
		Method:          types.MethodExecute,
		Payload:         rawRequest,
		MessageId:       messageID[:],
	}

	peerIDToDelay, err := transmission.GetPeerIDToTransmissionDelay(c.remoteCapabilityDonInfo.Members, c.localDONInfo.Config.SharedSecret, req.Metadata.WorkflowID, req.Metadata.WorkflowExecutionID, tc)
	if err != nil {
		return fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
	}

	for peerID, delay := range peerIDToDelay {
		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				c.lggr.Debugw("executing delayed execution for peer", "peerID", peerID)
				err = c.dispatcher.Send(peerID, message)
				if err != nil {
					c.lggr.Errorw("failed to send message", "peerID", peerID, "err", err)
				}
			}
		}(peerID, delay)
	}

	return nil
}

func (c *remoteTargetCaller) Receive(msg *types.MessageBody) {

	// TODO handle the case where the capability returns a stream of responses
	messageID := getMessageID(msg)

	wg, loaded := c.messageIDToWaitgroup.LoadAndDelete(messageID)
	if loaded {
		wg.(*sync.WaitGroup).Done()
		c.messageIDToResponse.Store(messageID, msg)
		return
	}
}
