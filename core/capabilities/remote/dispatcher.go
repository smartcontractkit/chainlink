package remote

import (
	"context"
	"fmt"
	sync "sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// dispatcher en/decodes messages and routes traffic between peers and capabilities
type dispatcher struct {
	peerWrapper p2ptypes.PeerWrapper
	peer        p2ptypes.Peer
	peerID      p2ptypes.PeerID
	signer      p2ptypes.Signer
	registry    core.CapabilitiesRegistry
	receivers   map[key]remotetypes.Receiver
	mu          sync.RWMutex
	stopCh      services.StopChan
	wg          sync.WaitGroup
	lggr        logger.Logger
}

type key struct {
	capId string
	donId string
}

var _ services.Service = &dispatcher{}

const supportedVersion = 1

func NewDispatcher(peerWrapper p2ptypes.PeerWrapper, signer p2ptypes.Signer, registry core.CapabilitiesRegistry, lggr logger.Logger) *dispatcher {
	return &dispatcher{
		peerWrapper: peerWrapper,
		signer:      signer,
		registry:    registry,
		receivers:   make(map[key]remotetypes.Receiver),
		stopCh:      make(services.StopChan),
		lggr:        lggr.Named("Dispatcher"),
	}
}

func (d *dispatcher) Start(ctx context.Context) error {
	d.peer = d.peerWrapper.GetPeer()
	d.peerID = d.peer.ID()
	if d.peer == nil {
		return fmt.Errorf("peer is not initialized")
	}
	d.wg.Add(1)
	go d.receive()
	d.lggr.Info("dispatcher started")
	return nil
}

func (d *dispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	k := key{capabilityId, donId}
	_, ok := d.receivers[k]
	if ok {
		return fmt.Errorf("receiver already exists for capability %s and don %s", capabilityId, donId)
	}
	d.receivers[k] = receiver
	return nil
}

func (d *dispatcher) RemoveReceiver(capabilityId string, donId string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.receivers, key{capabilityId, donId})
}

func (d *dispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = supportedVersion
	msgBody.Sender = d.peerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	rawBody, err := proto.Marshal(msgBody)
	if err != nil {
		return err
	}
	signature, err := d.signer.Sign(rawBody)
	if err != nil {
		return err
	}
	msg := &remotetypes.Message{Signature: signature, Body: rawBody}
	rawMsg, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return d.peer.Send(peerID, rawMsg)
}

func (d *dispatcher) receive() {
	defer d.wg.Done()
	recvCh := d.peer.Receive()
	for {
		select {
		case <-d.stopCh:
			d.lggr.Info("stopped - exiting receive")
			return
		case msg := <-recvCh:
			body, err := ValidateMessage(msg, d.peerID)
			if err != nil {
				d.lggr.Debugw("received invalid message", "error", err)
				d.tryRespondWithError(msg.Sender, body, types.Error_VALIDATION_FAILED)
				continue
			}
			k := key{body.CapabilityId, body.CapabilityDonId}
			d.mu.RLock()
			receiver, ok := d.receivers[k]
			d.mu.RUnlock()
			if !ok {
				d.lggr.Debugw("received message for unregistered capability", "capabilityId", k.capId, "donId", k.donId)
				d.tryRespondWithError(msg.Sender, body, types.Error_CAPABILITY_NOT_FOUND)
				continue
			}
			receiver.Receive(body)
		}
	}
}

func (d *dispatcher) tryRespondWithError(peerID p2ptypes.PeerID, body *remotetypes.MessageBody, errType types.Error) {
	if body == nil {
		return
	}
	if body.Error != types.Error_OK {
		d.lggr.Debug("received an invalid message with error field set - not responding to avoid an infinite loop")
		return
	}
	body.Error = errType
	// clear payload to reduce message size
	body.Payload = nil
	err := d.Send(peerID, body)
	if err != nil {
		d.lggr.Debugw("failed to send error response", "error", err)
	}
}

func (d *dispatcher) Close() error {
	close(d.stopCh)
	d.wg.Wait()
	d.lggr.Info("dispatcher closed")
	return nil
}

func (d *dispatcher) Ready() error {
	return nil
}

func (d *dispatcher) HealthReport() map[string]error {
	return nil
}

func (d *dispatcher) Name() string {
	return "Dispatcher"
}
