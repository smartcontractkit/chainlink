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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/trigger/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// TriggerSubscriber is a shim for remote trigger capabilities.
// It translates between capability API calls and network messages.
// Its responsibilities are:
//  1. Periodically refresh all registrations for remote triggers.
//  2. Collect trigger events from remote nodes and aggregate responses via a customizable aggregator.
//  3. Track the registration status for each workflow.
//
// TriggerSubscriber communicates with corresponding TriggerReceivers on remote nodes.
type triggerSubscriber struct {
	capabilityID  string
	capMethodName string
	dispatcher    types.Dispatcher
	cfg           atomic.Pointer[dynamicConfig]
	messageCache  *messagecache.MessageCache[triggerEventKey, p2ptypes.PeerID]
	// In theory we could identify all trigger registrations only by TriggerID (Workflow Engine
	// already includes WorkflowID inside TriggerID). However, keeping the workflowID has some benefits:
	//   - Protection against changes to Engine's logic.
	//   - Better logging and debugging.
	//   - Easier migration.
	// workflowID -> triggerID -> subRegState
	registeredWorkflows map[string]map[string]*registration.SubscriberRegistration
	mu                  sync.RWMutex // protects registeredWorkflows and messageCache
	stopCh              services.StopChan
	wg                  sync.WaitGroup
	lggr                logger.Logger

	// This channel is used to send initial registration requests immediately rather than waiting for the registration refresh cycle
	initialRegistrationRequestCh           chan *registration.SubscriberRegistration
	triggerRegistrationStatusUpdateTimeout limits.TimeLimiter
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
	triggerID      string
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
	maxBatchedWorkflowIDs = 1000
)

func NewTriggerSubscriber(capabilityID string, capMethodName string, dispatcher types.Dispatcher, lggr logger.Logger,
	triggerRegistrationStatusUpdateTimeout limits.TimeLimiter) *triggerSubscriber {
	return &triggerSubscriber{
		capabilityID:                           capabilityID,
		capMethodName:                          capMethodName,
		dispatcher:                             dispatcher,
		messageCache:                           messagecache.NewMessageCache[triggerEventKey, p2ptypes.PeerID](),
		registeredWorkflows:                    make(map[string]map[string]*registration.SubscriberRegistration),
		stopCh:                                 make(services.StopChan),
		lggr:                                   logger.With(logger.Named(lggr, "TriggerSubscriber"), "capabilityID", capabilityID, "capMethodName", capMethodName),
		initialRegistrationRequestCh:           make(chan *registration.SubscriberRegistration),
		triggerRegistrationStatusUpdateTimeout: triggerRegistrationStatusUpdateTimeout,
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

func (s *triggerSubscriber) AckEvent(ctx context.Context, triggerID string, eventID string, method string) error {
	s.lggr.Debugw("AckEvent called on subscriber", "triggerID", triggerID, "eventID", eventID)
	cfg := s.cfg.Load()
	for _, peerID := range cfg.capDonInfo.Members {
		m := &types.MessageBody{
			CapabilityId:     cfg.capInfo.ID,
			CapabilityDonId:  cfg.capDonInfo.ID,
			CallerDonId:      cfg.localDonID,
			Method:           types.MethodTriggerEventAck,
			CapabilityMethod: s.capMethodName,
			Metadata: &types.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &types.TriggerEventMetadata{
					TriggerEventId: eventID,
					TriggerIds:     []string{triggerID}, // triggerID contains workflowID
				},
			},
		}
		err := s.dispatcher.Send(peerID, m)
		if err != nil {
			s.lggr.Errorw("failed to send message", "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
		}
	}
	return nil
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

	s.mu.Lock()
	s.lggr.Infow("RegisterTrigger called", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID)
	triggerMap, ok := s.registeredWorkflows[request.Metadata.WorkflowID]
	if !ok {
		triggerMap = make(map[string]*registration.SubscriberRegistration)
		s.registeredWorkflows[request.Metadata.WorkflowID] = triggerMap
	}
	reg, existingRegistration := triggerMap[request.TriggerID]
	if !existingRegistration {
		reg = registration.NewSubscriberRegistration(s.lggr, rawRequest, request.Metadata.WorkflowID, request.TriggerID)
		triggerMap[request.TriggerID] = reg
	} else {
		reg.UpdateRegistrationRequest(rawRequest)
		s.lggr.Warnw("RegisterTrigger re-registering trigger", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID)
	}
	s.mu.Unlock()

	if existingRegistration {
		return reg.GetTriggerResponseChannel(), nil
	}

	// Send the first registration request immediately, do not wait for the registration refresh cycle
	s.initialRegistrationRequestCh <- reg

	// Wait for the registration status with a timeout.
	subCtx, subCancel, err := s.triggerRegistrationStatusUpdateTimeout.WithTimeout(ctx)
	if err != nil {
		s.lggr.Errorw("failed to create timeout context for trigger registration status update", "err", err)
		return reg.GetTriggerResponseChannel(), commoncap.ErrUnableToDetermineRegistrationStatus
	}
	defer subCancel()

	// Await registration with timeout - it is not guaranteed that registration status will be able to be determined
	// as may be running against a legacy remote capability that does not support registration status updates.
	start := time.Now()
	s.lggr.Debug("Starting registration wait", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID)
	regErr := reg.AwaitRegistration(subCtx)
	s.lggr.Debugw("Finished registration wait", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID, "waitTime", time.Since(start).String(), "regErr", regErr)
	if regErr == nil {
		return reg.GetTriggerResponseChannel(), nil
	}

	// In the case that the error is ErrUnableToDetermineRegistrationStatus, the caller may choose to ignore the error
	// and continue to wait for trigger events, this replicates the legacy behaviour to support remote nodes that do
	// not provide registration status updates.
	if errors.Is(regErr, commoncap.ErrUnableToDetermineRegistrationStatus) {
		s.lggr.Warnw("unable to determine registration status", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID)
		return reg.GetTriggerResponseChannel(), regErr
	}

	s.lggr.Errorw("registration error occurred", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "error", regErr)
	return nil, regErr
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
		case reg := <-s.initialRegistrationRequestCh:
			cfg := s.cfg.Load()
			s.sendRegistrationRequestToCapabilityDon(cfg, reg)
		case <-ticker.C:
			cfg := s.cfg.Load()
			if cfg.remoteConfig.RegistrationRefresh != tickerDuration {
				tickerDuration = cfg.remoteConfig.RegistrationRefresh
				ticker.Reset(tickerDuration)
			}

			s.mu.RLock()
			s.lggr.Infow("register trigger for remote capability", "donId", cfg.capDonInfo.ID, "nMembers", len(cfg.capDonInfo.Members), "nWorkflows", len(s.registeredWorkflows))
			if len(s.registeredWorkflows) == 0 {
				s.lggr.Infow("no workflows to register")
			}

			for _, regMap := range s.registeredWorkflows {
				for _, reg := range regMap {
					s.sendRegistrationRequestToCapabilityDon(cfg, reg)
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *triggerSubscriber) sendRegistrationRequestToCapabilityDon(cfg *dynamicConfig, reg *registration.SubscriberRegistration) {
	for _, peerID := range cfg.capDonInfo.Members {
		m := &types.MessageBody{
			CapabilityId:     cfg.capInfo.ID,
			CapabilityDonId:  cfg.capDonInfo.ID,
			CallerDonId:      cfg.localDonID,
			Method:           types.MethodRegisterTrigger,
			Payload:          reg.GetRawRequest(),
			CapabilityMethod: s.capMethodName,
		}
		err := s.dispatcher.Send(peerID, m)
		if err != nil {
			s.lggr.Errorw("failed to send message", "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
		}
	}
}

func (s *triggerSubscriber) UnregisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	triggerMap, ok := s.registeredWorkflows[request.Metadata.WorkflowID]
	if !ok {
		return nil
	}
	reg := triggerMap[request.TriggerID]
	if reg != nil {
		reg.Close()
	}
	delete(triggerMap, request.TriggerID)
	if len(triggerMap) == 0 {
		delete(s.registeredWorkflows, request.Metadata.WorkflowID)
	}
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
		s.lggr.Errorw("received message from unexpected node", "sender", sender)
		return
	}

	if msg.Method == types.MethodTriggerEvent {
		meta := msg.GetTriggerEventMetadata()
		if meta == nil {
			s.lggr.Errorw("received message with invalid trigger metadata", "sender", sender)
			return
		}
		if len(meta.WorkflowIds) > maxBatchedWorkflowIDs {
			s.lggr.Errorw("received message with too many workflow IDs - truncating", "nWorkflows", len(meta.WorkflowIds), "sender", sender)
			meta.WorkflowIds = meta.WorkflowIds[:maxBatchedWorkflowIDs]
		}
		for idx, workflowID := range meta.WorkflowIds {
			var triggerID string
			if idx < len(meta.TriggerIds) {
				triggerID = meta.TriggerIds[idx]
			}
			s.mu.RLock()
			triggerMap, found := s.registeredWorkflows[workflowID]
			var registration *registration.SubscriberRegistration
			if found {
				if triggerID != "" {
					// received a message from updated publisher, which provided a triggerID
					registration = triggerMap[triggerID]
				} else {
					// legacy flow, expect there to be only a single trigger of each type per workflow
					for _, reg := range triggerMap {
						registration = reg
						break
					}
					if len(triggerMap) > 1 {
						s.lggr.Errorw("received message without triggerID but workflow has multiple trigger - picking a random one", "workflowID", SanitizeLogString(workflowID), "sender", sender)
					}
				}
			}
			s.mu.RUnlock()
			if registration == nil {
				s.lggr.Errorw("received message for unregistered workflow/trigger", "workflowID", SanitizeLogString(workflowID), "triggerID", triggerID, "sender", sender)
				continue
			}
			key := triggerEventKey{
				triggerEventID: meta.TriggerEventId,
				workflowID:     workflowID,
				triggerID:      triggerID,
			}
			nowMs := time.Now().UnixMilli()
			s.mu.Lock()
			creationTs := s.messageCache.Insert(key, sender, nowMs, msg.Payload)
			ready, _, payloads := s.messageCache.Ready(key, cfg.remoteConfig.MinResponsesToAggregate, nowMs-cfg.remoteConfig.MessageExpiry.Milliseconds(), true)
			s.mu.Unlock()
			s.lggr.Debugw("trigger event received", "triggerEventId", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "sender", sender, "ready", ready, "nowTs", nowMs, "creationTs", creationTs, "minResponsesToAggregate", cfg.remoteConfig.MinResponsesToAggregate)
			if ready {
				aggregatedResponse, err := cfg.aggregator.Aggregate(meta.TriggerEventId, payloads)
				if err != nil {
					s.lggr.Errorw("failed to aggregate responses", "triggerEventID", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "err", err)
					continue
				}
				s.lggr.Infow("remote trigger event aggregated", "triggerEventID", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID)
				registration.SendAggregatedEvent(aggregatedResponse)
			}
		}
	} else if msg.Method == types.TriggerRegistrationStatus {
		meta := msg.GetTriggerRegistrationMetadata()
		if meta == nil {
			s.lggr.Errorw("received trigger registration status message with nil trigger registration metadata", "sender", sender)
			return
		}
		s.mu.Lock()

		triggerMap, ok := s.registeredWorkflows[meta.WorkflowId]
		if !ok {
			triggerMap = make(map[string]*registration.SubscriberRegistration)
			s.registeredWorkflows[meta.WorkflowId] = triggerMap
		}
		reg, found := triggerMap[meta.TriggerId]
		s.mu.Unlock()

		if !found {
			s.lggr.Warnw("received trigger registration status message for unregistered workflow", "workflowID", SanitizeLogString(meta.WorkflowId), "triggerID", SanitizeLogString(meta.TriggerId), "sender", sender)
			return
		}

		reg.HandleTriggerRegistrationStatusUpdate(sender, msg, cfg.remoteConfig.MinResponsesToAggregate, cfg.remoteConfig.MessageExpiry,
			len(cfg.capDonInfo.Members))
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
