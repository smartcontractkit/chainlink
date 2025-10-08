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

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
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

	messageCache  *messagecache.MessageCache[registrationKey, p2ptypes.PeerID]
	registrations map[registrationKey]*pubRegState
	mu            sync.RWMutex // protects messageCache and registrations
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

type registrationKey struct {
	callerDonID uint32
	workflowID  string
}

type pubRegState struct {
	callback <-chan commoncap.TriggerResponse
	request  commoncap.TriggerRegistrationRequest
	cancel   context.CancelFunc
}

type batchedResponse struct {
	rawResponse    []byte
	callerDonID    uint32
	triggerEventID string
	workflowIDs    []string
}

type TriggerPublisher interface {
	types.ReceiverService
	SetConfig(config *commoncap.RemoteTriggerConfig, underlying commoncap.TriggerCapability, capDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON) error
}

var _ TriggerPublisher = &triggerPublisher{}
var _ types.ReceiverService = &triggerPublisher{}

const minAllowedBatchCollectionPeriod = 10 * time.Millisecond

func NewTriggerPublisher(capID, method string, dispatcher types.Dispatcher, lggr logger.Logger) *triggerPublisher {
	return &triggerPublisher{
		capabilityID:  capID,
		capMethodName: method,
		dispatcher:    dispatcher,
		messageCache:  messagecache.NewMessageCache[registrationKey, p2ptypes.PeerID](),
		registrations: make(map[registrationKey]*pubRegState),
		batchingQueue: make(map[[32]byte]*batchedResponse),
		stopCh:        make(services.StopChan),
		lggr:          logger.Named(lggr, "TriggerPublisher"),
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
	go p.registrationCleanupLoop()
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
			p.lggr.Errorw("failed to unmarshal trigger registration request", "capabilityId", p.capabilityID, "err", err)
			return
		}
		callerDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			p.lggr.Errorw("received a message from unsupported workflow DON", "capabilityId", p.capabilityID, "callerDonId", msg.CallerDonId)
			return
		}
		if !cfg.membersCache[msg.CallerDonId][sender] {
			p.lggr.Errorw("sender not a member of its workflow DON", "capabilityId", p.capabilityID, "callerDonId", msg.CallerDonId, "sender", sender)
			return
		}
		if err = validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
			p.lggr.Errorw("received trigger request with invalid workflow ID", "capabilityId", p.capabilityID, "workflowId", SanitizeLogString(req.Metadata.WorkflowID), "err", err)
			return
		}
		p.lggr.Debugw("received trigger registration", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID, "sender", sender)
		key := registrationKey{msg.CallerDonId, req.Metadata.WorkflowID}
		nowMs := time.Now().UnixMilli()
		p.mu.Lock()
		defer p.mu.Unlock()
		p.messageCache.Insert(key, sender, nowMs, msg.Payload)
		_, exists := p.registrations[key]
		if exists {
			p.lggr.Debugw("trigger registration already exists", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID)
			return
		}
		// NOTE: require 2F+1 by default, introduce different strategies later (KS-76)
		minRequired := uint32(2*callerDon.F + 1)
		ready, payloads := p.messageCache.Ready(key, minRequired, nowMs-cfg.remoteConfig.RegistrationExpiry.Milliseconds(), false)
		if !ready {
			p.lggr.Debugw("not ready to aggregate yet", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID, "minRequired", minRequired)
			return
		}
		aggregated, err := aggregation.AggregateModeRaw(payloads, uint32(callerDon.F+1))
		if err != nil {
			p.lggr.Errorw("failed to aggregate trigger registrations", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID, "err", err)
			return
		}
		unmarshaled, err := pb.UnmarshalTriggerRegistrationRequest(aggregated)
		if err != nil {
			p.lggr.Errorw("failed to unmarshal request", "capabilityId", p.capabilityID, "err", err)
			return
		}
		ctx, cancel := p.stopCh.NewCtx()
		callbackCh, err := cfg.underlying.RegisterTrigger(ctx, unmarshaled)
		if err == nil {
			p.registrations[key] = &pubRegState{
				callback: callbackCh,
				request:  unmarshaled,
				cancel:   cancel,
			}
			p.wg.Add(1)
			go p.triggerEventLoop(callbackCh, key)
			p.lggr.Debugw("updated trigger registration", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID)
		} else {
			cancel()
			p.lggr.Errorw("failed to register trigger", "capabilityId", p.capabilityID, "workflowId", req.Metadata.WorkflowID, "err", err)
		}
	case types.MethodTriggerEvent:
		p.lggr.Errorw("trigger request failed with error",
			"method", SanitizeLogString(msg.Method), "sender", sender, "errorMsg", SanitizeLogString(msg.ErrorMsg))
	default:
		p.lggr.Errorw("received message with unknown method",
			"method", SanitizeLogString(msg.Method), "sender", sender)
	}
}

func (p *triggerPublisher) registrationCleanupLoop() {
	defer p.wg.Done()

	// Get initial config for ticker setup
	firstCfg := p.cfg.Load()
	if firstCfg == nil {
		p.lggr.Errorw("registrationCleanupLoop started but config not set")
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
			for key, req := range p.registrations {
				callerDon := cfg.workflowDONs[key.callerDonID]
				ready, _ := p.messageCache.Ready(key, uint32(2*callerDon.F+1), now-cfg.remoteConfig.RegistrationExpiry.Milliseconds(), false)
				if !ready {
					p.lggr.Infow("trigger registration expired", "capabilityId", p.capabilityID, "callerDonID", key.callerDonID, "workflowId", key.workflowID)
					ctx, cancel := p.stopCh.NewCtx()
					err := cfg.underlying.UnregisterTrigger(ctx, req.request)
					cancel()
					p.registrations[key].cancel() // Cancel context on register trigger
					p.lggr.Infow("unregistered trigger", "capabilityId", p.capabilityID, "callerDonID", key.callerDonID, "workflowId", key.workflowID, "err", err)
					// after calling UnregisterTrigger, the underlying trigger will not send any more events to the channel
					delete(p.registrations, key)
					p.messageCache.Delete(key)
				}
			}
			p.mu.Unlock()
		}
	}
}

func (p *triggerPublisher) triggerEventLoop(callbackCh <-chan commoncap.TriggerResponse, key registrationKey) {
	defer p.wg.Done()
	for {
		select {
		case <-p.stopCh:
			return
		case response, ok := <-callbackCh:
			if !ok {
				p.lggr.Infow("triggerEventLoop channel closed", "capabilityId", p.capabilityID, "workflowId", key.workflowID)
				return
			}

			triggerEvent := response.Event
			p.lggr.Debugw("received trigger event", "capabilityId", p.capabilityID, "workflowId", key.workflowID, "triggerEventID", triggerEvent.ID)
			marshaledResponse, err := pb.MarshalTriggerResponse(response)
			if err != nil {
				p.lggr.Debugw("can't marshal trigger event", "err", err)
				break
			}

			cfg := p.cfg.Load()
			if cfg.batchingEnabled {
				p.enqueueForBatching(marshaledResponse, key, triggerEvent.ID)
			} else {
				// a single-element "batch"
				p.sendBatch(&batchedResponse{
					rawResponse:    marshaledResponse,
					callerDonID:    key.callerDonID,
					triggerEventID: triggerEvent.ID,
					workflowIDs:    []string{key.workflowID},
				})
			}
		}
	}
}

func (p *triggerPublisher) enqueueForBatching(rawResponse []byte, key registrationKey, triggerEventID string) {
	// put in batching queue, group by hash(callerDonId, triggerEventID, response)
	combined := make([]byte, 4)
	binary.LittleEndian.PutUint32(combined, key.callerDonID)
	combined = append(combined, []byte(triggerEventID)...)
	combined = append(combined, rawResponse...)
	sha := sha256.Sum256(combined)
	p.bqMu.Lock()
	elem, exists := p.batchingQueue[sha]
	if !exists {
		elem = &batchedResponse{
			rawResponse:    rawResponse,
			callerDonID:    key.callerDonID,
			triggerEventID: triggerEventID,
			workflowIDs:    []string{key.workflowID},
		}
		p.batchingQueue[sha] = elem
	} else {
		elem.workflowIDs = append(elem.workflowIDs, key.workflowID)
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
		idBatch := resp.workflowIDs
		if cfg.batchingEnabled && int64(len(idBatch)) > int64(cfg.remoteConfig.MaxBatchSize) {
			idBatch = idBatch[:cfg.remoteConfig.MaxBatchSize]
			resp.workflowIDs = resp.workflowIDs[cfg.remoteConfig.MaxBatchSize:]
		} else {
			resp.workflowIDs = nil
		}
		msg := &types.MessageBody{
			CapabilityId:    p.capabilityID,
			CapabilityDonId: cfg.capDonInfo.ID,
			CallerDonId:     resp.callerDonID,
			Method:          types.MethodTriggerEvent,
			Payload:         resp.rawResponse,
			Metadata: &types.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &types.TriggerEventMetadata{
					WorkflowIds:    idBatch,
					TriggerEventId: resp.triggerEventID,
				},
			},
			CapabilityMethod: p.capMethodName,
		}
		// NOTE: send to all nodes by default, introduce different strategies later (KS-76)
		for _, peerID := range cfg.workflowDONs[resp.callerDonID].Members {
			err := p.dispatcher.Send(peerID, msg)
			if err != nil {
				p.lggr.Errorw("failed to send trigger event", "capabilityId", p.capabilityID, "peerID", peerID, "err", err)
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
