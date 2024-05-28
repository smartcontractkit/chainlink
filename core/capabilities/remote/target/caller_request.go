package target

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type callerRequest struct {
	transmissionCtx      context.Context
	responseCh           chan commoncap.CapabilityResponse
	transmissionCancelFn context.CancelFunc
	createdAt            time.Time
	responseIDCount      map[[32]byte]int
	responseReceived     map[p2ptypes.PeerID]bool

	requiredIdenticalResponses int

	respSent bool
}

func newCallerRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest, messageID string,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher) (*callerRequest, error) {

	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	tc, err := transmission.ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request config: %w", err)
	}

	peerIDToTransmissionDelay, err := transmission.GetPeerIDToTransmissionDelay(remoteCapabilityDonInfo.Members, localDonInfo.Config.SharedSecret,
		messageID, tc)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer ID to transmission delay: %w", err)
	}

	transmissionCtx, transmissionCancelFn := context.WithCancel(ctx)
	responseReceived := make(map[p2ptypes.PeerID]bool)
	for peerID, delay := range peerIDToTransmissionDelay {
		responseReceived[peerID] = false
		go func(peerID ragep2ptypes.PeerID, delay time.Duration) {
			message := &types.MessageBody{
				CapabilityId:    remoteCapabilityInfo.ID,
				CapabilityDonId: remoteCapabilityDonInfo.ID,
				CallerDonId:     localDonInfo.ID,
				Method:          types.MethodExecute,
				Payload:         rawRequest,
				MessageId:       []byte(messageID),
			}

			select {
			case <-transmissionCtx.Done():
				return
			case <-time.After(delay):
				err = dispatcher.Send(peerID, message)
				if err != nil {
					lggr.Errorw("failed to send message", "peerID", peerID, "err", err)
				}
			}
		}(peerID, delay)
	}

	return &callerRequest{
		createdAt:                  time.Now(),
		transmissionCancelFn:       transmissionCancelFn,
		requiredIdenticalResponses: int(remoteCapabilityDonInfo.F + 1),
		responseIDCount:            make(map[[32]byte]int),
		responseReceived:           responseReceived,
		responseCh:                 make(chan commoncap.CapabilityResponse, 1),
	}, nil
}

func (c *callerRequest) responseSent() bool {
	return c.respSent
}

// TODO addResponse assumes that only one response is received from each peer, if streaming responses need to be supported this will need to be updated
func (c *callerRequest) addResponse(sender p2ptypes.PeerID, response []byte) error {
	if _, ok := c.responseReceived[sender]; !ok {
		return fmt.Errorf("response from peer %s not expected", sender)
	}

	if c.responseReceived[sender] {
		return fmt.Errorf("response from peer %s already received", sender)
	}

	c.responseReceived[sender] = true

	responseID := sha256.Sum256(response)
	c.responseIDCount[responseID]++

	if c.responseIDCount[responseID] == c.requiredIdenticalResponses {
		capabilityResponse, err := pb.UnmarshalCapabilityResponse(response)
		if err != nil {
			c.sendResponse(commoncap.CapabilityResponse{Err: fmt.Errorf("failed to unmarshal capability response: %w", err)})
		} else {
			c.sendResponse(commoncap.CapabilityResponse{Value: capabilityResponse.Value})
		}
	}

	return nil
}

func (c *callerRequest) sendResponse(response commoncap.CapabilityResponse) {
	c.responseCh <- response
	close(c.responseCh)
	c.transmissionCancelFn()
	c.respSent = true
}

func (c *callerRequest) cancelRequest(reason string) {
	c.transmissionCancelFn()
	if !c.responseSent() {
		c.sendResponse(commoncap.CapabilityResponse{Err: errors.New(reason)})
	}
}
