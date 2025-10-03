package remote

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// TriggerSubscriber is a shim for remote trigger capabilities.
// It translates between capability API calls and network messages.
// Its responsibilities are:
//  1. Periodically refresh all registrations for remote triggers.
//  2. Collect trigger events from remote nodes and aggregate responses via a customizable aggregator.
//
// TriggerSubscriber communicates with corresponding TriggerReceivers on remote nodes.
type triggerSubscriber struct {
	capabilityID  string
	capMethodName string
	dispatcher    types.Dispatcher
	cfg           atomic.Pointer[dynamicConfig]

	messageCache        *messagecache.MessageCache[triggerEventKey, p2ptypes.PeerID]
	registeredWorkflows map[string]*subRegState
	mu                  sync.RWMutex // protects registeredWorkflows and messageCache
	stopCh              services.StopChan
	wg                  sync.WaitGroup
	lggr                logger.Logger
}

type dynamicConfig struct {
	remoteConfig  *commoncap.RemoteTriggerConfig
	capInfo       commoncap.CapabilityInfo
	capDonInfo    commoncap.DON
	capDonMembers map[p2ptypes.PeerID]struct{}
	localDonID    uint32
	aggregator    types.Aggregator
}

type triggerEventKey struct {
	triggerEventID string
	workflowID     string
}

type subRegState struct {
	callback   chan commoncap.TriggerResponse
	rawRequest []byte
}

type TriggerSubscriber interface {
	commoncap.TriggerCapability
	Receive(ctx context.Context, msg *types.MessageBody)
	SetConfig(config *commoncap.RemoteTriggerConfig, capInfo commoncap.CapabilityInfo, localDONID uint32, remoteDON commoncap.DON, aggregator types.Aggregator) error
}

var _ commoncap.TriggerCapability = &triggerSubscriber{}
var _ types.Receiver = &triggerSubscriber{}
var _ services.Service = &triggerSubscriber{}

const (
	// Engine reads trigger events without blocking and applies its own limits
	sendChannelBufferSize = 1000
	maxBatchedWorkflowIDs = 1000
)

func NewTriggerSubscriber(capabilityID string, capMethodName string, dispatcher types.Dispatcher, lggr logger.Logger) *triggerSubscriber {
	return &triggerSubscriber{
		capabilityID:        capabilityID,
		capMethodName:       capMethodName,
		dispatcher:          dispatcher,
		messageCache:        messagecache.NewMessageCache[triggerEventKey, p2ptypes.PeerID](),
		registeredWorkflows: make(map[string]*subRegState),
		stopCh:              make(services.StopChan),
		lggr:                logger.Named(lggr, "TriggerSubscriber"),
	}
}

func (s *triggerSubscriber) Start(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg := s.cfg.Load()

	// Validate that all required fields are set before starting
	if cfg == nil {
		return errors.New("config not set - call SetConfig() before Start()")
	}
	if cfg.remoteConfig == nil {
		return errors.New("remoteConfig not set - call SetConfig() before Start()")
	}
	if cfg.capInfo.ID == "" {
		return errors.New("capability info not set - call SetConfig() before Start()")
	}
	if cfg.localDonID == 0 {
		return errors.New("local DON ID not set - call SetConfig() before Start()")
	}
	if len(cfg.capDonInfo.Members) == 0 {
		return errors.New("capability DON info not set - call SetConfig() before Start()")
	}
	if cfg.aggregator == nil {
		return errors.New("aggregator not set - call SetAggregator() before Start()")
	}
	if s.dispatcher == nil {
		return errors.New("dispatcher set to nil, cannot start triggerSubscriber")
	}

	s.wg.Add(2)
	go s.registrationLoop()
	go s.eventCleanupLoop()
	s.lggr.Info("TriggerSubscriber started")
	return nil
}

func (s *triggerSubscriber) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	cfg := s.cfg.Load()
	if cfg == nil {
		return commoncap.CapabilityInfo{}, errors.New("config not set - call SetConfig() before Info()")
	}
	return cfg.capInfo, nil
}

func (s *triggerSubscriber) RegisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	if err != nil {
		return nil, err
	}
	if request.Metadata.WorkflowID == "" {
		return nil, errors.New("empty workflowID")
	}

	cfg := s.cfg.Load()
	if cfg == nil {
		return nil, errors.New("config not set - call SetConfig() first")
	}
	capID, capDonID := cfg.capInfo.ID, cfg.capDonInfo.ID

	s.mu.Lock()
	defer s.mu.Unlock()
	s.lggr.Infow("RegisterTrigger called", "capabilityId", capID, "donId", capDonID, "workflowID", request.Metadata.WorkflowID)
	regState, ok := s.registeredWorkflows[request.Metadata.WorkflowID]
	if !ok {
		regState = &subRegState{
			callback:   make(chan commoncap.TriggerResponse, sendChannelBufferSize),
			rawRequest: rawRequest,
		}
		s.registeredWorkflows[request.Metadata.WorkflowID] = regState
	} else {
		regState.rawRequest = rawRequest
		s.lggr.Warnw("RegisterTrigger re-registering trigger", "capabilityId", capID, "donId", capDonID, "workflowID", request.Metadata.WorkflowID)
	}

	return regState.callback, nil
}

func (s *triggerSubscriber) registrationLoop() {
	defer s.wg.Done()
	cfg := s.cfg.Load()
	tickerDuration := cfg.remoteConfig.RegistrationRefresh
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			cfg := s.cfg.Load()
			if cfg.remoteConfig.RegistrationRefresh != tickerDuration {
				tickerDuration = cfg.remoteConfig.RegistrationRefresh
				ticker.Reset(tickerDuration)
			}

			s.mu.RLock()
			s.lggr.Infow("register trigger for remote capability", "capabilityId", cfg.capInfo.ID, "donId", cfg.capDonInfo.ID, "nMembers", len(cfg.capDonInfo.Members), "nWorkflows", len(s.registeredWorkflows))
			if len(s.registeredWorkflows) == 0 {
				s.lggr.Infow("no workflows to register")
			}

			for _, registration := range s.registeredWorkflows {
				for _, peerID := range cfg.capDonInfo.Members {
					m := &types.MessageBody{
						CapabilityId:     cfg.capInfo.ID,
						CapabilityDonId:  cfg.capDonInfo.ID,
						CallerDonId:      cfg.localDonID,
						Method:           types.MethodRegisterTrigger,
						Payload:          registration.rawRequest,
						CapabilityMethod: s.capMethodName,
					}
					err := s.dispatcher.Send(peerID, m)
					if err != nil {
						s.lggr.Errorw("failed to send message", "capabilityId", cfg.capInfo.ID, "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
					}
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *triggerSubscriber) UnregisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.registeredWorkflows[request.Metadata.WorkflowID]
	if state != nil && state.callback != nil {
		close(state.callback)
	}
	delete(s.registeredWorkflows, request.Metadata.WorkflowID)
	// Registrations will quickly expire on all remote nodes.
	// Alternatively, we could send UnregisterTrigger messages right away.
	return nil
}

func (s *triggerSubscriber) Receive(_ context.Context, msg *types.MessageBody) {
	sender, err := ToPeerID(msg.Sender)
	if err != nil {
		s.lggr.Errorw("failed to convert message sender to PeerID", "err", err)
		return
	}
	cfg := s.cfg.Load()
	if cfg == nil {
		s.lggr.Errorw("config not set - call SetConfig() first")
		return
	}
	if _, found := cfg.capDonMembers[sender]; !found {
		s.lggr.Errorw("received message from unexpected node", "capabilityId", cfg.capInfo.ID, "sender", sender)
		return
	}

	if msg.Method == types.MethodTriggerEvent {
		meta := msg.GetTriggerEventMetadata()
		if meta == nil {
			s.lggr.Errorw("received message with invalid trigger metadata", "capabilityId", cfg.capInfo.ID, "sender", sender)
			return
		}
		if len(meta.WorkflowIds) > maxBatchedWorkflowIDs {
			s.lggr.Errorw("received message with too many workflow IDs - truncating", "capabilityId", cfg.capInfo.ID, "nWorkflows", len(meta.WorkflowIds), "sender", sender)
			meta.WorkflowIds = meta.WorkflowIds[:maxBatchedWorkflowIDs]
		}
		for _, workflowID := range meta.WorkflowIds {
			s.mu.RLock()
			registration, found := s.registeredWorkflows[workflowID]
			s.mu.RUnlock()
			if !found {
				s.lggr.Errorw("received message for unregistered workflow", "capabilityId", cfg.capInfo.ID, "workflowID", SanitizeLogString(workflowID), "sender", sender)
				continue
			}
			key := triggerEventKey{
				triggerEventID: meta.TriggerEventId,
				workflowID:     workflowID,
			}
			nowMs := time.Now().UnixMilli()
			s.mu.Lock()
			creationTs := s.messageCache.Insert(key, sender, nowMs, msg.Payload)
			ready, payloads := s.messageCache.Ready(key, cfg.remoteConfig.MinResponsesToAggregate, nowMs-cfg.remoteConfig.MessageExpiry.Milliseconds(), true)
			s.mu.Unlock()
			s.lggr.Debugw("trigger event received", "triggerEventId", meta.TriggerEventId, "capabilityId", cfg.capInfo.ID, "workflowId", workflowID, "sender", sender, "ready", ready, "nowTs", nowMs, "creationTs", creationTs, "minResponsesToAggregate", cfg.remoteConfig.MinResponsesToAggregate)
			if ready {
				aggregatedResponse, err := cfg.aggregator.Aggregate(meta.TriggerEventId, payloads)
				if err != nil {
					s.lggr.Errorw("failed to aggregate responses", "triggerEventID", meta.TriggerEventId, "capabilityId", cfg.capInfo.ID, "workflowId", workflowID, "err", err)
					continue
				}
				s.lggr.Infow("remote trigger event aggregated", "triggerEventID", meta.TriggerEventId, "capabilityId", cfg.capInfo.ID, "workflowId", workflowID)
				registration.callback <- aggregatedResponse
			}
		}
	} else {
		s.lggr.Errorw("received trigger event with unknown method", "method", SanitizeLogString(msg.Method), "sender", sender, "err", SanitizeLogString(msg.ErrorMsg))
	}
}

func (s *triggerSubscriber) eventCleanupLoop() {
	defer s.wg.Done()
	cfg := s.cfg.Load()
	cleanupInterval := cfg.remoteConfig.MessageExpiry
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			freshCfg := s.cfg.Load()
			remoteConfig := freshCfg.remoteConfig
			// Update cleanup interval if config has changed
			if remoteConfig.MessageExpiry != cleanupInterval {
				cleanupInterval = remoteConfig.MessageExpiry
				ticker.Reset(cleanupInterval)
			}
			s.mu.Lock()
			s.messageCache.DeleteOlderThan(time.Now().UnixMilli() - remoteConfig.MessageExpiry.Milliseconds())
			s.mu.Unlock()
		}
	}
}

func (s *triggerSubscriber) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	s.lggr.Info("TriggerSubscriber closed")
	return nil
}

func (s *triggerSubscriber) Ready() error {
	return nil
}

func (s *triggerSubscriber) HealthReport() map[string]error {
	return nil
}

func (s *triggerSubscriber) Name() string {
	return s.lggr.Name()
}

// SetConfig sets the remote trigger configuration, capability info, and DON information dynamically
func (s *triggerSubscriber) SetConfig(config *commoncap.RemoteTriggerConfig, capInfo commoncap.CapabilityInfo, localDONID uint32, remoteDON commoncap.DON, aggregator types.Aggregator) error {
	if config == nil {
		s.lggr.Info("SetConfig called with nil config, using defaults")
		config = &commoncap.RemoteTriggerConfig{}
	}
	config.ApplyDefaults()
	if capInfo.ID == "" || capInfo.ID != s.capabilityID {
		return fmt.Errorf("capability info provided does not match the subscriber's capabilityID: %s != %s", capInfo.ID, s.capabilityID)
	}
	if localDONID == 0 {
		return errors.New("localDONID=0 provided")
	}
	if remoteDON.ID == 0 || len(remoteDON.Members) == 0 {
		return errors.New("empty remoteDON provided")
	}
	if aggregator == nil {
		return errors.New("aggregator not set - call SetAggregator() before SetConfig()")
	}
	// Rebuild the capDonMembers map
	capDonMembers := make(map[p2ptypes.PeerID]struct{})
	for _, member := range remoteDON.Members {
		capDonMembers[member] = struct{}{}
	}

	// always replace the whole dynamicConfig object to avoid inconsistent state
	s.cfg.Store(&dynamicConfig{
		remoteConfig:  config,
		capInfo:       capInfo,
		capDonInfo:    remoteDON,
		capDonMembers: capDonMembers,
		localDonID:    localDONID,
		aggregator:    aggregator,
	})
	return nil
}
