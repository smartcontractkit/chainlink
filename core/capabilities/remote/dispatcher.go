package remote

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var (
	ErrReceiverExists = errors.New("receiver already exists")
)

// dispatcher en/decodes messages and routes traffic between peers and capabilities
type dispatcher struct {
	cfg         config.Dispatcher
	peerWrapper p2ptypes.PeerWrapper
	peer        p2ptypes.Peer
	peerID      p2ptypes.PeerID
	signer      p2ptypes.Signer
	registry    core.CapabilitiesRegistry
	rateLimiter *common.RateLimiter
	receivers   map[key]*receiver
	mu          sync.RWMutex
	stopCh      services.StopChan
	wg          sync.WaitGroup
	lggr        logger.Logger
}

type key struct {
	capId string
	donId uint32
}

var _ services.Service = &dispatcher{}

func NewDispatcher(cfg config.Dispatcher, peerWrapper p2ptypes.PeerWrapper, signer p2ptypes.Signer, registry core.CapabilitiesRegistry, lggr logger.Logger) (*dispatcher, error) {
	rl, err := common.NewRateLimiter(common.RateLimiterConfig{
		GlobalRPS:      cfg.RateLimit().GlobalRPS(),
		GlobalBurst:    cfg.RateLimit().GlobalBurst(),
		PerSenderRPS:   cfg.RateLimit().PerSenderRPS(),
		PerSenderBurst: cfg.RateLimit().PerSenderBurst(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rate limiter")
	}
	return &dispatcher{
		cfg:         cfg,
		peerWrapper: peerWrapper,
		signer:      signer,
		registry:    registry,
		rateLimiter: rl,
		receivers:   make(map[key]*receiver),
		stopCh:      make(services.StopChan),
		lggr:        lggr.Named("Dispatcher"),
	}, nil
}

func (d *dispatcher) Start(ctx context.Context) error {
	d.peer = d.peerWrapper.GetPeer()
	d.peerID = d.peer.ID()
	if d.peer == nil {
		return fmt.Errorf("peer is not initialized")
	}
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.receive()
	}()

	d.lggr.Info("dispatcher started")
	return nil
}

func (d *dispatcher) Close() error {
	close(d.stopCh)
	d.wg.Wait()
	d.lggr.Info("dispatcher closed")
	return nil
}

var capReceiveChannelUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "capability_receive_channel_usage",
	Help: "The usage of the receive channel for each capability, 0 indicates empty, 1 indicates full.",
}, []string{"capabilityId", "donId"})

type receiver struct {
	cancel context.CancelFunc
	ch     chan *types.MessageBody
}

func (d *dispatcher) SetReceiver(capabilityId string, donId uint32, rec types.Receiver) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	k := key{capabilityId, donId}
	_, ok := d.receivers[k]
	if ok {
		return fmt.Errorf("%w: receiver already exists for capability %s and don %d", ErrReceiverExists, capabilityId, donId)
	}

	receiverCh := make(chan *types.MessageBody, d.cfg.ReceiverBufferSize())

	ctx, cancelCtx := d.stopCh.NewCtx()
	d.wg.Add(1)
	go func() {
		defer cancelCtx()
		defer d.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-receiverCh:
				rec.Receive(ctx, msg)
			}
		}
	}()

	d.receivers[k] = &receiver{
		cancel: cancelCtx,
		ch:     receiverCh,
	}

	d.lggr.Debugw("receiver set", "capabilityId", capabilityId, "donId", donId)
	return nil
}

func (d *dispatcher) RemoveReceiver(capabilityId string, donId uint32) {
	d.mu.Lock()
	defer d.mu.Unlock()

	receiverKey := key{capabilityId, donId}
	if receiver, ok := d.receivers[receiverKey]; ok {
		receiver.cancel()
		delete(d.receivers, receiverKey)
		d.lggr.Debugw("receiver removed", "capabilityId", capabilityId, "donId", donId)
	}
}

func (d *dispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	msgBody.Version = uint32(d.cfg.SupportedVersion())
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
	msg := &types.Message{Signature: signature, Body: rawBody}
	rawMsg, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return d.peer.Send(peerID, rawMsg)
}

func (d *dispatcher) receive() {
	recvCh := d.peer.Receive()
	for {
		select {
		case <-d.stopCh:
			d.lggr.Info("stopped - exiting receive")
			return
		case msg := <-recvCh:
			if !d.rateLimiter.Allow(msg.Sender.String()) {
				d.lggr.Debugw("rate limit exceeded, dropping message", "sender", msg.Sender)
				continue
			}
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
				d.lggr.Debugw("received message for unregistered capability", "capabilityId", SanitizeLogString(k.capId), "donId", k.donId)
				d.tryRespondWithError(msg.Sender, body, types.Error_CAPABILITY_NOT_FOUND)
				continue
			}

			receiverQueueUsage := float64(len(receiver.ch)) / float64(d.cfg.ReceiverBufferSize())
			capReceiveChannelUsage.WithLabelValues(k.capId, fmt.Sprint(k.donId)).Set(receiverQueueUsage)
			select {
			case receiver.ch <- body:
			default:
				d.lggr.Warnw("receiver channel full, dropping message", "capabilityId", k.capId, "donId", k.donId)
			}
		}
	}
}

func (d *dispatcher) tryRespondWithError(peerID p2ptypes.PeerID, body *types.MessageBody, errType types.Error) {
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

func (d *dispatcher) Ready() error {
	return nil
}

func (d *dispatcher) HealthReport() map[string]error {
	return nil
}

func (d *dispatcher) Name() string {
	return "Dispatcher"
}
