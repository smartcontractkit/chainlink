package remote

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/trigger/registration"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// TriggerPublisher manages all external users of a local trigger capability.
// Its responsibilities are:
//  1. Manage trigger registrations from external nodes (receive, store, aggregate, expire).
//  2. Send out events produced by an underlying, concrete trigger implementation.
//
// TriggerPublisher communicates with corresponding TriggerSubscribers on remote nodes.
type triggerPublisher struct {
	capabilityID  string
	capMethodName string
	dispatcher    types.Dispatcher
	cfg           atomic.Pointer[dynamicPublisherConfig]

	registrations map[registration.ID]*registration.PublisherRegistration
	ackCache      *messagecache.MessageCache[ackKey, p2ptypes.PeerID]
	mu            sync.RWMutex // protects messageCache, ackCache, and registrations
	batchingQueue map[[32]byte]*batchedResponse
	bqMu          sync.Mutex // protects batchingQueue
	stopCh        services.StopChan
	wg            sync.WaitGroup
	lggr          logger.Logger
}

type dynamicPublisherConfig struct {
	remoteConfig    *commoncap.RemoteTriggerConfig
	underlying      commoncap.TriggerCapability
	capDonInfo      commoncap.DON
	workflowDONs    map[uint32]commoncap.DON
	membersCache    map[uint32]map[p2ptypes.PeerID]bool
	batchingEnabled bool
}

type ackKey struct {
	callerDonID    uint32
	triggerEventID string
	triggerID      string // triggerID contains the workflowID
}

type batchedResponse struct {
	rawResponse    []byte
	callerDonID    uint32
	triggerEventID string
	workflowIDs    []string
	triggerIDs     []string
}

type TriggerPublisher interface {
	types.ReceiverService
	SetConfig(config *commoncap.RemoteTriggerConfig, underlying commoncap.TriggerCapability, capDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON) error
}

var _ TriggerPublisher = &triggerPublisher{}
var _ types.ReceiverService = &triggerPublisher{}

const minAllowedBatchCollectionPeriod = 10 * time.Millisecond

func NewTriggerPublisher(capabilityID string, capMethodName string, dispatcher types.Dispatcher, lggr logger.Logger) *triggerPublisher {
	return &triggerPublisher{
		capabilityID:  capabilityID,
		capMethodName: capMethodName,
		dispatcher:    dispatcher,
		ackCache:      messagecache.NewMessageCache[ackKey, p2ptypes.PeerID](),
		registrations: make(map[registration.ID]*registration.PublisherRegistration),
		batchingQueue: make(map[[32]byte]*batchedResponse),
		stopCh:        make(services.StopChan),
		lggr:          logger.With(logger.Named(lggr, "TriggerPublisher"), "capabilityID", capabilityID, "capMethodName", capMethodName),
	}
}

// SetConfig sets the remote trigger configuration, capability info, and DON information dynamically
func (p *triggerPublisher) SetConfig(config *commoncap.RemoteTriggerConfig, underlying commoncap.TriggerCapability, capDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON) error {
	if config == nil {
		p.lggr.Info("SetConfig called with nil config, using defaults")
		config = &commoncap.RemoteTriggerConfig{}
	}
	config.ApplyDefaults()
	if underlying == nil {
		return errors.New("underlying trigger capability cannot be nil")
	}
	if capDonInfo.ID == 0 || len(capDonInfo.Members) == 0 {
		return errors.New("empty capDonInfo provided")
	}
	if workflowDONs == nil {
		workflowDONs = make(map[uint32]commoncap.DON)
	}

	// Build the members cache
	membersCache := make(map[uint32]map[p2ptypes.PeerID]bool)
	for id, don := range workflowDONs {
		cache := make(map[p2ptypes.PeerID]bool)
		for _, member := range don.Members {
			cache[member] = true
		}
		membersCache[id] = cache
	}

	// always replace the whole dynamicPublisherConfig object to avoid inconsistent state
	p.cfg.Store(&dynamicPublisherConfig{
		remoteConfig:    config,
		underlying:      underlying,
		capDonInfo:      capDonInfo,
		workflowDONs:    workflowDONs,
		membersCache:    membersCache,
		batchingEnabled: config.MaxBatchSize > 1 && config.BatchCollectionPeriod >= minAllowedBatchCollectionPeriod,
	})

	return nil
}

func (p *triggerPublisher) Start(ctx context.Context) error {
	cfg := p.cfg.Load()

	// Validate that all required fields are set before starting
	if cfg == nil {
		return errors.New("config not set - call SetConfig() before Start()")
	}
	if cfg.remoteConfig == nil {
		return errors.New("remoteConfig not set - call SetConfig() before Start()")
	}
	if cfg.underlying == nil {
		return errors.New("underlying trigger capability not set - call SetConfig() before Start()")
	}
	if len(cfg.capDonInfo.Members) == 0 {
		return errors.New("capability DON info not set - call SetConfig() before Start()")
	}
	if p.dispatcher == nil {
		return errors.New("dispatcher set to nil, cannot start triggerPublisher")
	}

	p.wg.Add(1)
	go p.cacheCleanupLoop()
	p.wg.Add(1)
	go p.batchingLoop()
	p.lggr.Info("TriggerPublisher started")
	return nil
}

func (p *triggerPublisher) Receive(_ context.Context, msg *types.MessageBody) {
	cfg := p.cfg.Load()
	if cfg == nil {
		p.lggr.Errorw("received message but config is not set")
		return
	}

	sender, err := ToPeerID(msg.Sender)
	if err != nil {
		p.lggr.Errorw("failed to convert message sender to PeerID", "err", err)
		return
	}

	if msg.ErrorMsg != "" {
		p.lggr.Errorw("received a message with error",
			"method", SanitizeLogString(msg.Method), "sender", sender, "errorMsg", SanitizeLogString(msg.ErrorMsg))
	}

	switch msg.Method {
	case types.MethodRegisterTrigger:
		req, err := pb.UnmarshalTriggerRegistrationRequest(msg.Payload)
		if err != nil {
			p.lggr.Errorw("failed to unmarshal trigger registration request", "err", err)
			return
		}
		callerDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			p.lggr.Errorw("received a message from unsupported workflow DON", "callerDonId", msg.CallerDonId)
			return
		}
		if !cfg.membersCache[msg.CallerDonId][sender] {
			p.lggr.Errorw("sender not a member of its workflow DON", "callerDonId", msg.CallerDonId, "sender", sender)
			return
		}
		if err = validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
			p.lggr.Errorw("received trigger request with invalid workflow ID", "workflowId", SanitizeLogString(req.Metadata.WorkflowID), "err", err)
			return
		}
		p.lggr.Debugw("received trigger registration", "workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "sender", sender)
		regID := registration.NewID(msg.CallerDonId, req.Metadata.WorkflowID, req.TriggerID)
		p.mu.Lock()
		defer p.mu.Unlock()
		reg, exists := p.registrations[regID]
		if !exists {
			p.lggr.Debugw("creating new trigger registration", "registrationID", regID)
			reg = registration.NewPublisherRegistration(p.lggr, regID, p.publishResponse,
				cfg.underlying,
				p.capabilityID, cfg.capDonInfo.ID, p.capMethodName, p.dispatcher)

			p.registrations[regID] = reg
		} else {
			p.lggr.Debugw("using existing trigger registration already exists", "registrationID", regID)
		}

		ctx, _ := p.stopCh.NewCtx()
		reg.AddRegistrationRequest(ctx, sender, msg.Payload, callerDon, cfg.remoteConfig.RegistrationExpiry)
	case types.MethodTriggerEvent:
		p.lggr.Errorw("trigger request failed with error",
			"method", SanitizeLogString(msg.Method), "sender", sender, "errorMsg", SanitizeLogString(msg.ErrorMsg))
	case types.MethodTriggerEventAck:
		triggerMetadata := msg.GetTriggerEventMetadata()
		if triggerMetadata == nil {
			p.lggr.Errorw("received empty trigger event ack metadata", "sender", sender)
			break
		}
		triggerEventID := triggerMetadata.TriggerEventId
		p.lggr.Debugw("received trigger event ACK", "sender", sender, "trigger event ID", triggerEventID)

		p.mu.Lock()
		defer p.mu.Unlock()
		callerDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			p.lggr.Errorw("received a message from unsupported workflow DON", "callerDonId", msg.CallerDonId)
			return
		}
		if !cfg.membersCache[msg.CallerDonId][sender] {
			p.lggr.Errorw("sender not a member of its workflow DON", "callerDonId", msg.CallerDonId, "sender", sender)
			return
		}

		if len(triggerMetadata.TriggerIds) != 1 {
			p.lggr.Errorw("did not receive single triggerID in ACK request", "callerDonId", msg.CallerDonId, "sender", sender, "triggerIDs", triggerMetadata.TriggerIds)
			return
		}
		triggerID := triggerMetadata.TriggerIds[0]

		key := ackKey{msg.CallerDonId, triggerEventID, triggerID}
		nowMs := time.Now().UnixMilli()
		p.ackCache.Insert(key, sender, nowMs, msg.Payload)
		minRequired := uint32(2*callerDon.F + 1)
		ready, _, _ := p.ackCache.Ready(key, minRequired, nowMs-cfg.remoteConfig.MessageExpiry.Milliseconds(), false)
		if !ready {
			p.lggr.Debugw("not ready to ACK trigger event yet", "triggerEventId", triggerEventID, "minRequired", minRequired)
			return
		}

		ctx, cancel := p.stopCh.NewCtx()
		defer cancel()
		p.lggr.Debugw("ACKing trigger event", "triggerEventId", triggerEventID)
		err = cfg.underlying.AckEvent(ctx, triggerID, triggerEventID, p.capMethodName)
		if err != nil {
			p.lggr.Errorw("failed to AckEvent on underlying trigger capability",
				"eventID", triggerEventID, "capabilityID", p.capabilityID, "err", err)
		}
	default:
		p.lggr.Errorw("received message with unknown method",
			"method", SanitizeLogString(msg.Method), "sender", sender)
	}
}

func (p *triggerPublisher) cacheCleanupLoop() {
	defer p.wg.Done()

	// Get initial config for ticker setup
	firstCfg := p.cfg.Load()
	if firstCfg == nil {
		p.lggr.Errorw("cacheCleanupLoop started but config not set")
		return
	}
	cleanupInterval := firstCfg.remoteConfig.MessageExpiry
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			cfg := p.cfg.Load()
			// Update cleanup interval if config has changed
			if cfg.remoteConfig.MessageExpiry != cleanupInterval {
				cleanupInterval = cfg.remoteConfig.MessageExpiry
				ticker.Reset(cleanupInterval)
			}
			now := time.Now().UnixMilli()

			p.mu.Lock()
			for regID, reg := range p.registrations {
				callerDon := cfg.workflowDONs[regID.CallerDonID()]
				if !reg.IsLive(cfg.remoteConfig.RegistrationExpiry, callerDon) {
					p.lggr.Infow("trigger registration expired", "ID", regID)
					ctx, cancel := p.stopCh.NewCtx()
					err := reg.Close(ctx)
					cancel()
					p.lggr.Infow("unregistered trigger", "ID", regID, "err", err)
					// after calling UnregisterTrigger, the underlying trigger will not send any more events to the channel
					delete(p.registrations, regID)
				}
			}

			deleted := p.ackCache.DeleteOlderThan(now - cfg.remoteConfig.MessageExpiry.Milliseconds())
			p.mu.Unlock()

			if deleted > 0 {
				p.lggr.Debugw("cleaned expired AckCache entries", "deleted", deleted)
			}
		}
	}
}

func (p *triggerPublisher) publishResponse(registrationID registration.ID, response commoncap.TriggerResponse) {
	triggerEvent := response.Event
	p.lggr.Debugw("received trigger event", "registrationID", registrationID, "triggerEventID", triggerEvent.ID)
	marshaledResponse, err := pb.MarshalTriggerResponse(response)
	if err != nil {
		p.lggr.Debugw("can't marshal trigger event", "err", err)
		return
	}

	cfg := p.cfg.Load()
	if cfg.batchingEnabled {
		p.enqueueForBatching(marshaledResponse, registrationID, triggerEvent.ID)
	} else {
		// a single-element "batch"
		p.sendBatch(&batchedResponse{
			rawResponse:    marshaledResponse,
			callerDonID:    registrationID.CallerDonID(),
			triggerEventID: triggerEvent.ID,
			workflowIDs:    []string{registrationID.WorkflowID()},
			triggerIDs:     []string{registrationID.TriggerID()},
		})
	}
}

func (p *triggerPublisher) enqueueForBatching(rawResponse []byte, registrationID registration.ID, triggerEventID string) {
	// put in batching queue, group by hash(callerDonId, triggerEventID, response)
	combined := make([]byte, 4)
	binary.LittleEndian.PutUint32(combined, registrationID.CallerDonID())
	combined = append(combined, []byte(triggerEventID)...)
	combined = append(combined, rawResponse...)
	sha := sha256.Sum256(combined)
	p.bqMu.Lock()
	elem, exists := p.batchingQueue[sha]
	if !exists {
		elem = &batchedResponse{
			rawResponse:    rawResponse,
			callerDonID:    registrationID.CallerDonID(),
			triggerEventID: triggerEventID,
			workflowIDs:    []string{registrationID.WorkflowID()},
			triggerIDs:     []string{registrationID.TriggerID()},
		}
		p.batchingQueue[sha] = elem
	} else {
		elem.workflowIDs = append(elem.workflowIDs, registrationID.WorkflowID())
		elem.triggerIDs = append(elem.triggerIDs, registrationID.TriggerID())
	}
	p.bqMu.Unlock()
}

func (p *triggerPublisher) sendBatch(resp *batchedResponse) {
	cfg := p.cfg.Load()
	if cfg == nil {
		p.lggr.Errorw("config not set during sendBatch")
		return
	}

	for len(resp.workflowIDs) > 0 {
		workflowBatch := resp.workflowIDs
		triggerBatch := resp.triggerIDs
		if cfg.batchingEnabled && int64(len(workflowBatch)) > int64(cfg.remoteConfig.MaxBatchSize) {
			workflowBatch = workflowBatch[:cfg.remoteConfig.MaxBatchSize]
			triggerBatch = triggerBatch[:cfg.remoteConfig.MaxBatchSize]
			resp.workflowIDs = resp.workflowIDs[cfg.remoteConfig.MaxBatchSize:]
			resp.triggerIDs = resp.triggerIDs[cfg.remoteConfig.MaxBatchSize:]
		} else {
			resp.workflowIDs = nil
			resp.triggerIDs = nil
		}

		ackSnapshot := make(map[string]map[p2ptypes.PeerID]bool)
		p.mu.RLock()
		for _, triggerID := range triggerBatch {
			key := ackKey{
				callerDonID:    resp.callerDonID,
				triggerEventID: resp.triggerEventID,
				triggerID:      triggerID,
			}
			ackSnapshot[triggerID] = p.ackCache.Peers(key)
		}
		p.mu.RUnlock()

		for _, peerID := range cfg.workflowDONs[resp.callerDonID].Members {
			var missingTriggerIDs []string
			var missingWorkflowIDs []string

			// determine which triggerIDs / workflowIDs have not yet ACKd this trigger event
			for i, triggerID := range triggerBatch {
				peers := ackSnapshot[triggerID]
				if peers == nil || !peers[peerID] {
					missingTriggerIDs = append(missingTriggerIDs, triggerID)
					missingWorkflowIDs = append(missingWorkflowIDs, workflowBatch[i])
				}
			}

			if len(missingTriggerIDs) == 0 {
				p.lggr.Debugw("skipping trigger event send; all triggerIDs already ACKed by peer",
					"peerID", peerID,
					"callerDonID", resp.callerDonID,
					"triggerEventID", resp.triggerEventID,
					"triggerIDs", triggerBatch,
				)
				continue
			}

			p.lggr.Debugw("sending trigger event to peer",
				"peerID", peerID,
				"callerDonID", resp.callerDonID,
				"triggerEventID", resp.triggerEventID,
				"workflowIDs", missingWorkflowIDs,
				"triggerIDs", missingTriggerIDs,
			)

			msg := &types.MessageBody{
				CapabilityId:     p.capabilityID,
				CapabilityDonId:  cfg.capDonInfo.ID,
				CallerDonId:      resp.callerDonID,
				Method:           types.MethodTriggerEvent,
				Payload:          resp.rawResponse,
				CapabilityMethod: p.capMethodName,
				Metadata: &types.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &types.TriggerEventMetadata{
						WorkflowIds:    missingWorkflowIDs,
						TriggerIds:     missingTriggerIDs,
						TriggerEventId: resp.triggerEventID,
					},
				},
			}

			err := p.dispatcher.Send(peerID, msg)
			if err != nil {
				p.lggr.Errorw("failed to send trigger event", "peerID", peerID, "err", err)
			}
		}
	}
}

func (p *triggerPublisher) batchingLoop() {
	defer p.wg.Done()

	// Get initial config for ticker setup
	firstCfg := p.cfg.Load()
	if firstCfg == nil {
		p.lggr.Errorw("batchingLoop started but config not set")
		return
	}
	interval := firstCfg.remoteConfig.BatchCollectionPeriod
	ticker := time.NewTicker(interval)

	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			cfg := p.cfg.Load()
			// Update cleanup interval if config has changed
			if cfg.remoteConfig.MessageExpiry != interval {
				interval = cfg.remoteConfig.BatchCollectionPeriod
				ticker.Reset(interval)
			}

			p.bqMu.Lock()
			queue := p.batchingQueue
			p.batchingQueue = make(map[[32]byte]*batchedResponse)
			p.bqMu.Unlock()

			for _, elem := range queue {
				p.sendBatch(elem)
			}
		}
	}
}

func (p *triggerPublisher) Close() error {
	close(p.stopCh)
	p.wg.Wait()

	// Close all registrations
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, reg := range p.registrations {
		err := reg.Close(context.Background())
		if err != nil {
			p.lggr.Errorw("error closing registration", "err", err)
		}
	}

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
	return p.lggr.Name()
}
