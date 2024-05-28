package remote

import (
	"context"
	sync "sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// TriggerPublisher manages all external users of a local trigger capability.
// Its responsibilities are:
//  1. Manage trigger registrations from external nodes (receive, store, aggregate, expire).
//  2. Send out events produced by an underlying, concrete trigger implementation.
//
// TriggerPublisher communicates with corresponding TriggerSubscribers on remote nodes.
type triggerPublisher struct {
	config        types.RemoteTriggerConfig
	underlying    commoncap.TriggerCapability
	capInfo       commoncap.CapabilityInfo
	capDonInfo    commoncap.DON
	workflowDONs  map[string]commoncap.DON
	dispatcher    types.Dispatcher
	messageCache  *messageCache[registrationKey, p2ptypes.PeerID]
	registrations map[registrationKey]*pubRegState
	mu            sync.RWMutex // protects messageCache and registrations
	stopCh        services.StopChan
	wg            sync.WaitGroup
	lggr          logger.Logger
}

type registrationKey struct {
	callerDonId string
	workflowId  string
}

type pubRegState struct {
	callback <-chan commoncap.CapabilityResponse
	request  commoncap.CapabilityRequest
}

var _ types.Receiver = &triggerPublisher{}
var _ services.Service = &triggerPublisher{}

func NewTriggerPublisher(config types.RemoteTriggerConfig, underlying commoncap.TriggerCapability, capInfo commoncap.CapabilityInfo, capDonInfo commoncap.DON, workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, lggr logger.Logger) *triggerPublisher {
	config.ApplyDefaults()
	return &triggerPublisher{
		config:        config,
		underlying:    underlying,
		capInfo:       capInfo,
		capDonInfo:    capDonInfo,
		workflowDONs:  workflowDONs,
		dispatcher:    dispatcher,
		messageCache:  NewMessageCache[registrationKey, p2ptypes.PeerID](),
		registrations: make(map[registrationKey]*pubRegState),
		stopCh:        make(services.StopChan),
		lggr:          lggr,
	}
}

func (p *triggerPublisher) Start(ctx context.Context) error {
	p.wg.Add(1)
	go p.registrationCleanupLoop()
	p.lggr.Info("TriggerPublisher started")
	return nil
}

func (p *triggerPublisher) Receive(msg *types.MessageBody) {
	sender := ToPeerID(msg.Sender)
	if msg.Method == types.MethodRegisterTrigger {
		req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
		if err != nil {
			p.lggr.Errorw("failed to unmarshal capability request", "capabilityId", p.capInfo.ID, "err", err)
			return
		}
		callerDon, ok := p.workflowDONs[msg.CallerDonId]
		if !ok {
			p.lggr.Errorw("received a message from unsupported workflow DON", "capabilityId", p.capInfo.ID, "callerDonId", msg.CallerDonId)
			return
		}
		p.lggr.Debugw("received trigger registration", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID, "sender", sender)
		key := registrationKey{msg.CallerDonId, req.Metadata.WorkflowID}
		nowMs := time.Now().UnixMilli()
		p.mu.Lock()
		defer p.mu.Unlock()
		p.messageCache.Insert(key, sender, nowMs, msg.Payload)
		_, exists := p.registrations[key]
		if exists {
			p.lggr.Debugw("trigger registration already exists", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID)
			return
		}
		// NOTE: require 2F+1 by default, introduce different strategies later (KS-76)
		minRequired := uint32(2*callerDon.F + 1)
		ready, payloads := p.messageCache.Ready(key, minRequired, nowMs-int64(p.config.RegistrationExpiryMs), false)
		if !ready {
			p.lggr.Debugw("not ready to aggregate yet", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID, "minRequired", minRequired)
			return
		}
		aggregated, err := AggregateModeRaw(payloads, uint32(callerDon.F+1))
		if err != nil {
			p.lggr.Errorw("failed to aggregate trigger registrations", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID, "err", err)
			return
		}
		unmarshaled, err := pb.UnmarshalCapabilityRequest(aggregated)
		if err != nil {
			p.lggr.Errorw("failed to unmarshal request", "capabilityId", p.capInfo.ID, "err", err)
			return
		}
		ctx, cancel := p.stopCh.NewCtx()
		callbackCh, err := p.underlying.RegisterTrigger(ctx, unmarshaled)
		cancel()
		if err == nil {
			p.registrations[key] = &pubRegState{
				callback: callbackCh,
				request:  unmarshaled,
			}
			p.wg.Add(1)
			go p.triggerEventLoop(callbackCh, key)
			p.lggr.Debugw("updated trigger registration", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID)
		} else {
			p.lggr.Errorw("failed to register trigger", "capabilityId", p.capInfo.ID, "workflowId", req.Metadata.WorkflowID, "err", err)
		}
	} else {
		p.lggr.Errorw("received trigger request with unknown method", "method", msg.Method, "sender", sender)
	}
}

func (p *triggerPublisher) registrationCleanupLoop() {
	defer p.wg.Done()
	ticker := time.NewTicker(time.Duration(p.config.RegistrationExpiryMs) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			now := time.Now().UnixMilli()
			p.mu.RLock()
			for key, req := range p.registrations {
				callerDon := p.workflowDONs[key.callerDonId]
				ready, _ := p.messageCache.Ready(key, uint32(2*callerDon.F+1), now-int64(p.config.RegistrationExpiryMs), false)
				if !ready {
					p.lggr.Infow("trigger registration expired", "capabilityId", p.capInfo.ID, "callerDonID", key.callerDonId, "workflowId", key.workflowId)
					ctx, cancel := p.stopCh.NewCtx()
					err := p.underlying.UnregisterTrigger(ctx, req.request)
					cancel()
					p.lggr.Infow("unregistered trigger", "capabilityId", p.capInfo.ID, "callerDonID", key.callerDonId, "workflowId", key.workflowId, "err", err)
					// after calling UnregisterTrigger, the underlying trigger will not send any more events to the channel
					delete(p.registrations, key)
					p.messageCache.Delete(key)
				}
			}
			p.mu.RUnlock()
		}
	}
}

func (p *triggerPublisher) triggerEventLoop(callbackCh <-chan commoncap.CapabilityResponse, key registrationKey) {
	defer p.wg.Done()
	for {
		select {
		case <-p.stopCh:
			return
		case response, ok := <-callbackCh:
			if !ok {
				p.lggr.Infow("triggerEventLoop channel closed", "capabilityId", p.capInfo.ID, "workflowId", key.workflowId)
				return
			}
			triggerEvent := capabilities.TriggerEvent{}
			err := response.Value.UnwrapTo(&triggerEvent)
			if err != nil {
				p.lggr.Errorw("can't unwrap trigger event", "capabilityId", p.capInfo.ID, "workflowId", key.workflowId, "err", err)
				break
			}
			p.lggr.Debugw("received trigger event", "capabilityId", p.capInfo.ID, "workflowId", key.workflowId, "triggerEventID", triggerEvent.ID)
			marshaled, err := pb.MarshalCapabilityResponse(response)
			if err != nil {
				p.lggr.Debugw("can't marshal trigger event", "err", err)
				break
			}
			msg := &types.MessageBody{
				CapabilityId:    p.capInfo.ID,
				CapabilityDonId: p.capDonInfo.ID,
				CallerDonId:     key.callerDonId,
				Method:          types.MethodTriggerEvent,
				Payload:         marshaled,
				Metadata: &types.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &types.TriggerEventMetadata{
						// NOTE: optionally introduce batching across workflows as an optimization
						WorkflowIds:    []string{key.workflowId},
						TriggerEventId: triggerEvent.ID,
					},
				},
			}
			// NOTE: send to all nodes by default, introduce different strategies later (KS-76)
			for _, peerID := range p.workflowDONs[key.callerDonId].Members {
				err = p.dispatcher.Send(peerID, msg)
				if err != nil {
					p.lggr.Errorw("failed to send trigger event", "capabilityId", p.capInfo.ID, "peerID", peerID, "err", err)
				}
			}
		}
	}
}

func (p *triggerPublisher) Close() error {
	close(p.stopCh)
	p.wg.Wait()
	p.lggr.Info("TriggerPublisher closed")
	return nil
}

func (p *triggerPublisher) Ready() error {
	return nil
}

func (p *triggerPublisher) HealthReport() map[string]error {
	return nil
}

func (p *triggerPublisher) Name() string {
	return "TriggerPublisher"
}
