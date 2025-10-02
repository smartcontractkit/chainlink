package remote

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var (
	ErrReceiverExists = errors.New("receiver already exists")
)

// dispatcher en/decodes messages and routes traffic between peers and capabilities
type dispatcher struct {
	cfg               config.Dispatcher
	peerWrapper       p2ptypes.PeerWrapper
	peer              p2ptypes.Peer
	peerID            p2ptypes.PeerID
	signer            p2ptypes.Signer
	don2donSharedPeer p2ptypes.SharedPeer
	registry          core.CapabilitiesRegistry
	rateLimiter       *ratelimit.RateLimiter
	receivers         map[key]*receiver
	mu                sync.RWMutex
	stopCh            services.StopChan
	wg                sync.WaitGroup
	lggr              logger.Logger

	metrics dispatcherMetrics
}

type dispatcherMetrics struct {
	externalPeerMsgsRcvdCounter metric.Int64Counter
	sharedPeerMsgsRcvdCounter   metric.Int64Counter
}

var _ types.Dispatcher = &dispatcher{}

type key struct {
	capID      string
	donID      uint32
	methodName string
}

var _ services.Service = &dispatcher{}

func NewDispatcher(cfg config.Dispatcher, peerWrapper p2ptypes.PeerWrapper, don2donSharedPeer p2ptypes.SharedPeer, signer p2ptypes.Signer, registry core.CapabilitiesRegistry, lggr logger.Logger) (*dispatcher, error) {
	rl, err := ratelimit.NewRateLimiter(ratelimit.RateLimiterConfig{
		GlobalRPS:      cfg.RateLimit().GlobalRPS(),
		GlobalBurst:    cfg.RateLimit().GlobalBurst(),
		PerSenderRPS:   cfg.RateLimit().PerSenderRPS(),
		PerSenderBurst: cfg.RateLimit().PerSenderBurst(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rate limiter")
	}
	return &dispatcher{
		cfg:               cfg,
		peerWrapper:       peerWrapper,
		signer:            signer,
		registry:          registry,
		rateLimiter:       rl,
		receivers:         make(map[key]*receiver),
		stopCh:            make(services.StopChan),
		lggr:              logger.Named(lggr, "Dispatcher"),
		don2donSharedPeer: don2donSharedPeer,
	}, nil
}

func (d *dispatcher) initMetrics() error {
	var err error
	d.metrics.externalPeerMsgsRcvdCounter, err = beholder.GetMeter().Int64Counter("platform_don2don_dispatcher_external_peer_msgs_rcvd_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_don2don_dispatcher_external_peer_msgs_rcvd_total): %w", err)
	}
	d.metrics.sharedPeerMsgsRcvdCounter, err = beholder.GetMeter().Int64Counter("platform_don2don_dispatcher_shared_peer_msgs_rcvd_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_don2don_dispatcher_shared_peer_msgs_rcvd_total): %w", err)
	}
	return nil
}

func (d *dispatcher) Start(ctx context.Context) error {
	if d.peerWrapper == nil && d.don2donSharedPeer == nil {
		return errors.New("either peerWrapper or don2donSharedPeer must be set")
	}
	if d.peerWrapper != nil {
		d.peer = d.peerWrapper.GetPeer()
		d.peerID = d.peer.ID()
		if d.peer == nil {
			return errors.New("peer is not initialized")
		}
	}
	if d.don2donSharedPeer != nil {
		if (d.peerID != p2ptypes.PeerID{}) && d.peerID != d.don2donSharedPeer.ID() {
			return errors.New("peer ID from peerWrapper and don2donSharedPeer do not match")
		}
		d.peerID = d.don2donSharedPeer.ID()
	}
	err := d.signer.Initialize()
	if err != nil {
		return errors.Wrap(err, "failed to initialize signer")
	}
	if err = d.initMetrics(); err != nil {
		return errors.Wrap(err, "failed to initialize metrics")
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

func (d *dispatcher) SetReceiverForMethod(capabilityID string, donID uint32, method string, rec types.Receiver) error {
	return d.setReceiver(key{capabilityID, donID, method}, rec)
}

func (d *dispatcher) SetReceiver(capabilityID string, donID uint32, rec types.Receiver) error {
	return d.setReceiver(key{capabilityID, donID, ""}, rec) // empty method name
}

func (d *dispatcher) setReceiver(k key, rec types.Receiver) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.receivers[k]
	if ok {
		return fmt.Errorf("%w: receiver already exists for capability %s, donID %d, method %s", ErrReceiverExists, k.capID, k.donID, k.methodName)
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

	d.lggr.Debugw("receiver set", "capabilityId", k.capID, "donId", k.donID, "methodName", k.methodName)
	return nil
}

func (d *dispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, method string) {
	d.removeReceiver(key{capabilityID, donID, method})
}

func (d *dispatcher) RemoveReceiver(capabilityID string, donID uint32) {
	d.removeReceiver(key{capabilityID, donID, ""}) // empty method name
}

func (d *dispatcher) removeReceiver(k key) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if receiver, ok := d.receivers[k]; ok {
		// NOTE: receiver.ch is not drained or closed - handle it if receivers are ever dynamically removed/re-added.
		receiver.cancel()
		delete(d.receivers, k)
		d.lggr.Debugw("receiver removed", "capabilityId", k.capID, "donId", k.donID, "methodName", k.methodName)
	}
}

func (d *dispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	//nolint:gosec // disable G115
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
	if d.cfg.SendToSharedPeer() {
		return d.don2donSharedPeer.Send(peerID, rawMsg)
	}
	if d.peer != nil {
		return d.peer.Send(peerID, rawMsg)
	}
	return errors.New("no peer available to send message")
}

func (d *dispatcher) receive() {
	externalPeerRecvCh := make(<-chan p2ptypes.Message)
	if d.peer != nil {
		externalPeerRecvCh = d.peer.Receive()
	}
	sharedPeerRecvCh := make(<-chan p2ptypes.Message)
	if d.don2donSharedPeer != nil {
		sharedPeerRecvCh = d.don2donSharedPeer.Receive()
	}
	ctx, cancel := d.stopCh.NewCtx()
	defer cancel()
	for {
		select {
		case <-d.stopCh:
			d.lggr.Info("stopped - exiting receive")
			return
		case msg := <-externalPeerRecvCh: // deprecated, will be removed in favor of SharedPeer (CRE-707)
			d.metrics.externalPeerMsgsRcvdCounter.Add(ctx, 1)
			d.handleMessage(&msg)
		case msg, ok := <-sharedPeerRecvCh:
			if !ok {
				d.lggr.Info("shared peer channel closed - exiting receive")
				return
			}
			d.metrics.sharedPeerMsgsRcvdCounter.Add(ctx, 1)
			d.handleMessage(&msg)
		}
	}
}

func (d *dispatcher) handleMessage(msg *p2ptypes.Message) {
	if !d.rateLimiter.Allow(msg.Sender.String()) {
		d.lggr.Errorw("rate limit exceeded, dropping message", "sender", msg.Sender)
		return
	}
	body, err := ValidateMessage(msg, d.peerID)
	if err != nil {
		d.lggr.Debugw("received invalid message", "error", err)
		d.tryRespondWithError(msg.Sender, body, types.Error_VALIDATION_FAILED)
		return
	}
	// CapabilityMethod will be empty for legacy "v1" messages
	k := key{body.CapabilityId, body.CapabilityDonId, body.CapabilityMethod}
	d.mu.RLock()
	receiver, ok := d.receivers[k]
	d.mu.RUnlock()
	if !ok {
		d.lggr.Debugw("received message for unregistered capability or method", "capabilityId", SanitizeLogString(k.capID), "donId", k.donID, "method", k.methodName)
		d.tryRespondWithError(msg.Sender, body, types.Error_CAPABILITY_NOT_FOUND)
		return
	}

	receiverQueueUsage := float64(0)
	if d.cfg.ReceiverBufferSize() > 0 {
		receiverQueueUsage = float64(len(receiver.ch)) / float64(d.cfg.ReceiverBufferSize())
	}
	capReceiveChannelUsage.WithLabelValues(k.capID, strconv.FormatUint(uint64(k.donID), 10)).Set(receiverQueueUsage)
	select {
	case receiver.ch <- body:
	default:
		d.lggr.Warnw("receiver channel full, dropping message", "capabilityId", k.capID, "donId", k.donID)
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
	return d.lggr.Name()
}
