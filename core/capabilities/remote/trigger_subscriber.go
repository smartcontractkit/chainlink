package remote

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
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
	messageCache  *messagecache.MessageCache[triggerEventKey, p2ptypes.PeerID]
	// In theory we could identify all trigger registrations only by TriggerID (Workflow Engine
	// already includes WorkflowID inside TriggerID). However, keeping the workflowID has some benefits:
	//   - Protection against changes to Engine's logic.
	//   - Better logging and debugging.
	//   - Easier migration.
	// workflowID -> triggerID -> subRegState
	registeredWorkflows map[string]map[string]*subRegState
	// registrationBuckets[i] holds the registrations due to resend during second i of the
	// refresh window. registrationLoop only ever touches the one bucket due "now".
	// Resized (and everything re-bucketed) whenever RegistrationRefresh changes.
	registrationBuckets []map[*subRegState]struct{}
	// Records that the workflow engine already issued ACK fan-out for this logical event
	// (registration TriggerID + TriggerEventId). Used to replay ACK on duplicate Receive without
	// re-aggregating or delivering to the engine again.
	ackReplayCache     map[ackReplayKey]int64
	propagateOrgIDGate limits.GateLimiter
	metrics            triggerSubscriberMetrics
	mu                 sync.RWMutex // protects registeredWorkflows, registrationBuckets, messageCache, and ackReplayCache
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	closeOnce          sync.Once
	lggr               logger.Logger
}

type triggerSubscriberMetrics struct {
	registrationOrgIDCounter metric.Int64Counter
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

type ackReplayKey struct {
	triggerID      string
	triggerEventID string
}

type subRegState struct {
	callback chan commoncap.TriggerResponse
	payloads registrationPayloads
	// bucket is the index into triggerSubscriber.registrationBuckets
	// assigned randomly so resends spread across the refresh window
	bucket int
	// sendMu serializes deliveries with closing callback; done unblocks a delivery waiting on a full callback.
	sendMu sync.Mutex
	done   chan struct{}
}

func (r *subRegState) deliver(resp commoncap.TriggerResponse, stopCh <-chan struct{}) bool {
	r.sendMu.Lock()
	defer r.sendMu.Unlock()
	// callback is closed only after done, so checking done first avoids sending on a closed channel.
	select {
	case <-r.done:
		return false
	default:
	}
	select {
	case r.callback <- resp:
		return true
	case <-r.done:
	case <-stopCh:
	}
	return false
}

// close must be called at most once.
func (r *subRegState) close() {
	close(r.done)
	r.sendMu.Lock()
	defer r.sendMu.Unlock()
	close(r.callback)
}

type registrationPayloads struct {
	withOrgID    []byte
	withoutOrgID []byte
	hasOrgID     bool
}

func (p registrationPayloads) forFlag(propagateOrgID bool) []byte {
	if propagateOrgID {
		return p.withOrgID
	}
	return p.withoutOrgID
}

func newRegistrationPayloads(request commoncap.TriggerRegistrationRequest) (registrationPayloads, error) {
	p := registrationPayloads{}
	withOrgID, err := pb.MarshalTriggerRegistrationRequest(request)
	if err != nil {
		return registrationPayloads{}, err
	}
	p.withOrgID = withOrgID
	p.hasOrgID = request.Metadata.OrgID != ""
	request.Metadata.OrgID = ""
	p.withoutOrgID, err = pb.MarshalTriggerRegistrationRequest(request)
	if err != nil {
		return registrationPayloads{}, err
	}

	return p, nil
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

	// registrationBucketDuration is both the registrationLoop poll period and the width of each
	// registrationBuckets slot: registrations are spread across RegistrationRefresh/registrationBucketDuration
	// buckets so each tick only has to look at the one bucket due "now".
	registrationBucketDuration = 1 * time.Second
)

// numRegistrationBuckets is how many registrationBucketDuration-wide buckets span one refresh
// window, floored at 1 so very short refresh values (e.g. in tests) still work.
// Because of the integer division, the effective resend period is RegistrationRefresh rounded
// down to a whole number of seconds, and never less than 1s: e.g. 1.5s resends every 1s, and
// anything below 1s (e.g. 500ms) also resends every 1s, i.e. less often than configured.
func numRegistrationBuckets(refresh time.Duration) int {
	n := int(refresh / registrationBucketDuration)
	if n < 1 {
		return 1
	}
	return n
}

// ensureBucketsLocked resizes registrationBuckets to match refresh, re-bucketing any existing
// registrations if the bucket count changed. Caller must hold s.mu for writing.
func (s *triggerSubscriber) ensureBucketsLocked(refresh time.Duration) {
	want := numRegistrationBuckets(refresh)
	if len(s.registrationBuckets) == want {
		return
	}
	newBuckets := make([]map[*subRegState]struct{}, want)
	for i := range newBuckets {
		newBuckets[i] = make(map[*subRegState]struct{})
	}
	for _, old := range s.registrationBuckets {
		for regState := range old {
			regState.bucket = rand.IntN(want) //nolint:gosec // G404: weak random is fine for send-time jitter
			newBuckets[regState.bucket][regState] = struct{}{}
		}
	}
	s.registrationBuckets = newBuckets
}

func NewTriggerSubscriber(capabilityID string, capMethodName string, dispatcher types.Dispatcher, propagateOrgIDGate limits.GateLimiter, lggr logger.Logger) *triggerSubscriber {
	return &triggerSubscriber{
		capabilityID:        capabilityID,
		capMethodName:       capMethodName,
		dispatcher:          dispatcher,
		messageCache:        messagecache.NewMessageCache[triggerEventKey, p2ptypes.PeerID](),
		ackReplayCache:      make(map[ackReplayKey]int64),
		propagateOrgIDGate:  propagateOrgIDGate,
		registeredWorkflows: make(map[string]map[string]*subRegState),
		stopCh:              make(services.StopChan),
		lggr:                logger.With(logger.Named(lggr, "TriggerSubscriber"), "capabilityID", capabilityID, "capMethodName", capMethodName),
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
	if s.propagateOrgIDGate == nil {
		return errors.New("propagateOrgIDGate set to nil, cannot start triggerSubscriber")
	}
	if err := s.initMetrics(); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
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
	rk := ackReplayKey{triggerID: triggerID, triggerEventID: eventID}
	s.mu.Lock()
	s.ackReplayCache[rk] = time.Now().UnixMilli()
	s.mu.Unlock()
	return nil
}

func (s *triggerSubscriber) RegisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	payloads, err := newRegistrationPayloads(request)
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
	defer s.mu.Unlock()
	s.lggr.Infow("RegisterTrigger called", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID)
	triggerMap, ok := s.registeredWorkflows[request.Metadata.WorkflowID]
	if !ok {
		triggerMap = make(map[string]*subRegState)
		s.registeredWorkflows[request.Metadata.WorkflowID] = triggerMap
	}
	regState, ok := triggerMap[request.TriggerID]
	if !ok {
		// Assign a random bucket so resends spread across the refresh window instead of bursting.
		bucket := rand.IntN(len(s.registrationBuckets)) //nolint:gosec // G404: weak random is fine for send-time jitter
		regState = &subRegState{
			callback: make(chan commoncap.TriggerResponse, sendChannelBufferSize),
			payloads: payloads,
			bucket:   bucket,
			done:     make(chan struct{}),
		}
		triggerMap[request.TriggerID] = regState
		s.registrationBuckets[bucket][regState] = struct{}{}
	} else {
		regState.payloads = payloads
		s.lggr.Warnw("RegisterTrigger re-registering trigger", "donId", cfg.capDonInfo.ID, "workflowID", request.Metadata.WorkflowID, "triggerID", request.TriggerID)
	}

	return regState.callback, nil
}

func (s *triggerSubscriber) registrationLoop() {
	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	defer s.wg.Done()
	ticker := time.NewTicker(registrationBucketDuration)
	defer ticker.Stop()
	lastSentSec := time.Now().Unix()
	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			lastSentSec = s.sendDueBuckets(ctx, lastSentSec, now.Unix())
		}
	}
}

// sendDueBuckets sends registrations in every bucket due in seconds (lastSentSec, nowSec] and
// returns the new lastSentSec. Catching up on all elapsed seconds, rather than only the one
// due "now", ensures no bucket is skipped when ticks are dropped (time.Ticker drops ticks while
// a slow send is in progress) or land on both sides of a second boundary. The catch-up is
// capped at one full refresh window so each registration is sent at most once per call.
func (s *triggerSubscriber) sendDueBuckets(ctx context.Context, lastSentSec, nowSec int64) int64 {
	if nowSec == lastSentSec {
		return lastSentSec // already sent this second's bucket
	}
	if nowSec < lastSentSec {
		// Wall clock moved backwards; resync instead of stalling until it catches up.
		lastSentSec = nowSec - 1
	}

	cfg := s.cfg.Load()

	s.mu.Lock()
	nBuckets := int64(len(s.registrationBuckets))
	fromSec := max(lastSentSec+1, nowSec-nBuckets+1)
	var due []registrationPayloads
	for sec := fromSec; sec <= nowSec; sec++ {
		for regState := range s.registrationBuckets[sec%nBuckets] {
			due = append(due, regState.payloads)
		}
	}
	s.mu.Unlock()

	if len(due) == 0 {
		return nowSec
	}

	propagateOrgID := s.propagateOrgID(ctx)

	var totalP2PSends, totalSendErrors, withOrgID int
	for _, payloads := range due {
		if payloads.hasOrgID {
			withOrgID++
		}
		rawRequest := payloads.forFlag(propagateOrgID)
		for _, peerID := range cfg.capDonInfo.Members {
			m := &types.MessageBody{
				CapabilityId:     cfg.capInfo.ID,
				CapabilityDonId:  cfg.capDonInfo.ID,
				CallerDonId:      cfg.localDonID,
				Method:           types.MethodRegisterTrigger,
				Payload:          rawRequest, // triggerID is in the raw request
				CapabilityMethod: s.capMethodName,
			}
			if err := s.dispatcher.Send(peerID, m); err != nil {
				totalSendErrors++
				s.lggr.Errorw("failed to send message", "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
			} else {
				totalP2PSends++
			}
		}
	}
	s.addRegistrationOrgIDCount(ctx, withOrgID, true, propagateOrgID)
	s.addRegistrationOrgIDCount(ctx, len(due)-withOrgID, false, propagateOrgID)
	s.lggr.Infow("registrationLoop tick: sent due buckets",
		"donId", cfg.capDonInfo.ID,
		"nCapDonMembers", len(cfg.capDonInfo.Members),
		"nBucketsDue", nowSec-fromSec+1,
		"nDue", len(due),
		"p2pSendsSent", totalP2PSends,
		"p2pSendErrors", totalSendErrors)
	return nowSec
}

func (s *triggerSubscriber) initMetrics() error {
	var err error
	s.metrics.registrationOrgIDCounter, err = beholder.GetMeter().Int64Counter("platform_trigger_subscriber_registration_org_id_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_subscriber_registration_org_id_total: %w", err)
	}
	return nil
}

func (s *triggerSubscriber) addRegistrationOrgIDCount(ctx context.Context, n int, hasOrgID, propagated bool) {
	if n == 0 {
		return
	}
	s.metrics.registrationOrgIDCounter.Add(ctx, int64(n), metric.WithAttributes(
		attribute.String("capabilityID", s.capabilityID),
		attribute.String("capMethodName", s.capMethodName),
		attribute.String("hasOrgID", strconv.FormatBool(hasOrgID)),
		attribute.String("orgIDPropagateEnabled", strconv.FormatBool(propagated)),
	))
}

func (s *triggerSubscriber) propagateOrgID(ctx context.Context) bool {
	open, err := s.propagateOrgIDGate.IsOpen(ctx)
	if err != nil {
		s.lggr.Warnw("failed to evaluate PropagateOrgIDInRequestMetadata; omitting OrgID from trigger registrations", "err", err)
		return false
	}
	return open
}

func (s *triggerSubscriber) UnregisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	triggerMap, ok := s.registeredWorkflows[request.Metadata.WorkflowID]
	if !ok {
		return nil
	}
	state := triggerMap[request.TriggerID]
	if state != nil {
		state.close()
		delete(s.registrationBuckets[state.bucket], state)
	}
	delete(triggerMap, request.TriggerID)
	if len(triggerMap) == 0 {
		delete(s.registeredWorkflows, request.Metadata.WorkflowID)
	}
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

	switch msg.Method {
	case types.MethodTriggerEvent:
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
			s.mu.RLock()
			_, found := s.registeredWorkflows[workflowID]
			s.mu.RUnlock()
			if !found {
				s.lggr.Errorw("received message for unregistered workflow/trigger", "workflowID", SanitizeLogString(workflowID), "sender", sender)
				continue
			}
			s.mu.Lock()
			var triggerID string
			if idx < len(meta.TriggerIds) {
				triggerID = meta.TriggerIds[idx]
			}
			var registration *subRegState
			triggerMap, found := s.registeredWorkflows[workflowID]
			if found {
				if triggerID != "" {
					// received a message from updated publisher, which provided a triggerID
					registration = triggerMap[triggerID]
				} else {
					// legacy flow, expect there to be only a single trigger of each type per workflow
					for tid, reg := range triggerMap {
						registration = reg
						triggerID = tid // canonical registration id for caches and AckEvent replay
						break
					}
					if len(triggerMap) > 1 {
						s.lggr.Errorw("received message without triggerID but workflow has multiple trigger - picking a random one", "workflowID", SanitizeLogString(workflowID), "sender", sender)
					}
				}
			}
			if registration == nil {
				s.lggr.Errorw("received message for unregistered workflow/trigger", "workflowID", SanitizeLogString(workflowID), "triggerID", triggerID, "sender", sender)
				continue // was unregistered in the meantime
			}

			key := triggerEventKey{
				triggerEventID: meta.TriggerEventId,
				workflowID:     workflowID,
				triggerID:      triggerID,
			}
			rk := ackReplayKey{triggerID: triggerID, triggerEventID: meta.TriggerEventId}
			if _, ok := s.ackReplayCache[rk]; ok {
				// Event has already been ACKd by engine, so we don't need to re-deliver
				s.mu.Unlock()
				ctx, cancel := s.stopCh.NewCtx()
				err := s.AckEvent(ctx, triggerID, meta.TriggerEventId, s.capMethodName)
				cancel()
				if err != nil {
					s.lggr.Errorw("replay AckEvent failed", "triggerID", triggerID, "triggerEventID", meta.TriggerEventId, "err", err)
				} else {
					s.lggr.Infow("replayed ACK fan-out for duplicate trigger event after prior engine ACK",
						"triggerID", triggerID, "triggerEventID", meta.TriggerEventId, "sender", sender)
				}
				continue
			}

			nowMs := time.Now().UnixMilli()
			creationTs := s.messageCache.Insert(key, sender, nowMs, msg.Payload)
			ready, payloads := s.messageCache.Ready(key, cfg.remoteConfig.MinResponsesToAggregate, nowMs-cfg.remoteConfig.MessageExpiry.Milliseconds(), true)
			if !ready {
				s.mu.Unlock()
				s.lggr.Debugw("trigger event received", "triggerEventId", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "sender", sender, "ready", ready, "nowTs", nowMs, "creationTs", creationTs, "minResponsesToAggregate", cfg.remoteConfig.MinResponsesToAggregate)
				continue
			}

			aggregatedResponse, err := cfg.aggregator.Aggregate(meta.TriggerEventId, payloads)
			if err != nil {
				s.mu.Unlock()
				s.lggr.Debugw("trigger event received", "triggerEventId", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "sender", sender, "ready", ready, "nowTs", nowMs, "creationTs", creationTs, "minResponsesToAggregate", cfg.remoteConfig.MinResponsesToAggregate)
				s.lggr.Errorw("failed to aggregate responses", "triggerEventID", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "err", err)
				continue
			}
			s.messageCache.MarkDelivered(key)
			s.mu.Unlock()
			s.lggr.Debugw("trigger event received", "triggerEventId", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID, "sender", sender, "ready", ready, "nowTs", nowMs, "creationTs", creationTs, "minResponsesToAggregate", cfg.remoteConfig.MinResponsesToAggregate)
			s.lggr.Infow("remote trigger event aggregated", "triggerEventID", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID)
			if !registration.deliver(aggregatedResponse, s.stopCh) {
				s.lggr.Warnw("dropped aggregated trigger event: trigger unregistered or subscriber stopping", "triggerEventID", meta.TriggerEventId, "workflowId", workflowID, "triggerID", triggerID)
			}
		}
	case types.MethodTriggerRegistrationCheck:
		meta := msg.GetTriggerEventMetadata()
		if meta == nil {
			s.lggr.Errorw("received registration check with nil metadata", "sender", sender)
			return
		}

		if len(meta.WorkflowIds) != len(meta.TriggerIds) {
			s.lggr.Errorw("received registration check with mismatched workflow and trigger IDs",
				"sender", sender,
				"workflowIdsLen", len(meta.WorkflowIds),
				"triggerIdsLen", len(meta.TriggerIds))
			return
		}

		s.lggr.Infow("received registration check", "sender", sender)
		for i, workflowID := range meta.WorkflowIds {
			triggerID := meta.TriggerIds[i]

			s.mu.RLock()
			triggerMap, ok := s.registeredWorkflows[workflowID]
			reg := ok && triggerMap[triggerID] != nil
			s.mu.RUnlock()

			if !reg {
				s.lggr.Infow("sending unregister in response to registration check", "workflowID", workflowID, "triggerID", triggerID)
				// Registration was removed locally — tell the publisher to clean up.
				s.sendUnregister(workflowID, triggerID)
			}
			// For existing registrations we intentionally do NOT resend.
			// The periodic registrationLoop already refreshes registrations,
			// and the publisher ignores duplicate MethodRegisterTrigger for
			// registrations it already has.
		}
	default:
		s.lggr.Errorw("received trigger event with unknown method", "method", SanitizeLogString(msg.Method), "sender", sender, "err", SanitizeLogString(msg.ErrorMsg))
	}
}

func (s *triggerSubscriber) sendUnregister(workflowID, triggerID string) {
	if s.dispatcher == nil {
		return
	}

	cfg := s.cfg.Load()

	for _, peerID := range cfg.capDonInfo.Members {
		m := &types.MessageBody{
			CapabilityId:     cfg.capInfo.ID,
			CapabilityDonId:  cfg.capDonInfo.ID,
			CallerDonId:      cfg.localDonID,
			Method:           types.MethodUnregisterTrigger,
			CapabilityMethod: s.capMethodName,
			Metadata: &types.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &types.TriggerEventMetadata{
					WorkflowIds: []string{workflowID},
					TriggerIds:  []string{triggerID},
				},
			},
		}
		err := s.dispatcher.Send(peerID, m)
		if err != nil {
			s.lggr.Errorw("failed to send message", "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
		}
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
			cutoff := time.Now().UnixMilli() - remoteConfig.MessageExpiry.Milliseconds()
			for k, ts := range s.ackReplayCache {
				if ts < cutoff {
					delete(s.ackReplayCache, k)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *triggerSubscriber) Close() error {
	s.closeOnce.Do(func() {
		close(s.stopCh)
		s.wg.Wait()
		s.lggr.Info("TriggerSubscriber closed")
	})
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

	s.mu.Lock()
	s.ensureBucketsLocked(config.RegistrationRefresh)
	s.mu.Unlock()

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
