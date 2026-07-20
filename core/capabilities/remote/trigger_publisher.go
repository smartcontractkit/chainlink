package remote

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
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

	messageCache     *messagecache.MessageCache[registrationKey, p2ptypes.PeerID]
	registrations    map[registrationKey]*pubRegState
	unregisterCache  *messagecache.MessageCache[registrationKey, p2ptypes.PeerID]
	ackCache         *messagecache.MessageCache[ackKey, p2ptypes.PeerID]
	mu               sync.RWMutex // protects messageCache, ackCache, unregisterCache, and registrations
	batchingQueue    map[[32]byte]*batchedResponse
	bqMu             sync.Mutex // protects batchingQueue
	stopCh           services.StopChan
	wg               sync.WaitGroup
	ackExecutor      *ParallelExecutor
	registerExecutor *ParallelExecutor
	lggr             logger.Logger
	metrics          triggerPublisherMetrics
}

type triggerPublisherMetrics struct {
	registerTriggerCounter      metric.Int64Counter
	unregisterTriggerCounter    metric.Int64Counter
	ackEventCounter             metric.Int64Counter
	registerTriggerDurationMs   metric.Int64Histogram
	unregisterTriggerDurationMs metric.Int64Histogram
	ackEventDurationMs          metric.Int64Histogram
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
	triggerID   string
}

type ackKey struct {
	callerDonID    uint32
	triggerEventID string
	triggerID      string // triggerID contains the workflowID
}

type pubRegState struct {
	callback        <-chan commoncap.TriggerResponse
	request         commoncap.TriggerRegistrationRequest
	cancel          context.CancelFunc
	registrationErr error // non-nil if RegisterTrigger returned an error; used to suppress retries
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

const (
	minAllowedBatchCollectionPeriod = 10 * time.Millisecond
	defaultMaxParallelAcks          = 100 // TODO: make this configurable https://smartcontract-it.atlassian.net/browse/PLEX-3266
	defaultMaxParallelRegisters     = 100 // TODO: make this configurable https://smartcontract-it.atlassian.net/browse/PLEX-3266
)

func NewTriggerPublisher(capabilityID string, capMethodName string, dispatcher types.Dispatcher, lggr logger.Logger) *triggerPublisher {
	slotUsageAttrs := []attribute.KeyValue{
		attribute.String("capabilityID", capabilityID),
		attribute.String("capMethodName", capMethodName),
	}

	return &triggerPublisher{
		capabilityID:    capabilityID,
		capMethodName:   capMethodName,
		dispatcher:      dispatcher,
		messageCache:    messagecache.NewMessageCache[registrationKey, p2ptypes.PeerID](),
		ackCache:        messagecache.NewMessageCache[ackKey, p2ptypes.PeerID](),
		unregisterCache: messagecache.NewMessageCache[registrationKey, p2ptypes.PeerID](),
		registrations:   make(map[registrationKey]*pubRegState),
		batchingQueue:   make(map[[32]byte]*batchedResponse),
		stopCh:          make(services.StopChan),
		lggr:            logger.With(logger.Named(lggr, "TriggerPublisher"), "capabilityID", capabilityID, "capMethodName", capMethodName),
		ackExecutor: NewParallelExecutor(defaultMaxParallelAcks,
			"trigger_publisher_ack_executor",
			slotUsageAttrs...,
		),
		registerExecutor: NewParallelExecutor(defaultMaxParallelRegisters,
			"trigger_publisher_register_executor",
			slotUsageAttrs...,
		),
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

func (p *triggerPublisher) initMetrics() error {
	var err error
	p.metrics.registerTriggerCounter, err = beholder.GetMeter().Int64Counter("platform_trigger_publisher_register_trigger_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_register_trigger_total: %w", err)
	}
	p.metrics.unregisterTriggerCounter, err = beholder.GetMeter().Int64Counter("platform_trigger_publisher_unregister_trigger_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_unregister_trigger_total: %w", err)
	}
	p.metrics.ackEventCounter, err = beholder.GetMeter().Int64Counter("platform_trigger_publisher_ack_event_total")
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_ack_event_total: %w", err)
	}
	durationBuckets := metric.WithExplicitBucketBoundaries(1, 50, 250, 1_000, 5_000, 30_000)
	p.metrics.registerTriggerDurationMs, err = beholder.GetMeter().Int64Histogram("platform_trigger_publisher_register_trigger_duration_ms", durationBuckets)
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_register_trigger_duration_ms: %w", err)
	}
	p.metrics.unregisterTriggerDurationMs, err = beholder.GetMeter().Int64Histogram("platform_trigger_publisher_unregister_trigger_duration_ms", durationBuckets)
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_unregister_trigger_duration_ms: %w", err)
	}
	p.metrics.ackEventDurationMs, err = beholder.GetMeter().Int64Histogram("platform_trigger_publisher_ack_event_duration_ms", durationBuckets)
	if err != nil {
		return fmt.Errorf("failed to register platform_trigger_publisher_ack_event_duration_ms: %w", err)
	}
	return nil
}

func (p *triggerPublisher) recordMethodDuration(ctx context.Context, hist metric.Int64Histogram, outcome string, d time.Duration) {
	hist.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("capabilityID", p.capabilityID),
		attribute.String("capMethodName", p.capMethodName),
		attribute.String("outcome", outcome),
	))
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
	if err := p.initMetrics(); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	p.wg.Add(1)
	go p.cacheCleanupLoop()
	p.wg.Add(1)
	go p.registrationCheckLoop()
	p.wg.Add(1)
	go p.batchingLoop()

	if err := p.ackExecutor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start ack executor: %w", err)
	}
	if err := p.registerExecutor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start register executor: %w", err)
	}

	p.lggr.Info("TriggerPublisher started")
	return nil
}

// TODO: Reduce complexity -> https://smartcontract-it.atlassian.net/browse/PLEX-3258
func (p *triggerPublisher) Receive(ctx context.Context, msg *types.MessageBody) {
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
		registerStart := time.Now()
		registerOutcome := "invalid"
		// Duration is recorded either by this defer (synchronous Receive exits) or by the
		// register executor goroutine (success/error). recordRegisterDuration=false on
		// successful schedule prevents double-counting.
		//
		// Interpreting outcomes:
		//   - success/error: quorum reached, RegisterTrigger ran (includes executor queue
		//     wait + underlying work such as DB). Use for end-to-end register latency.
		//   - anything else (not_ready, quorum_already_met, skipped_exists, schedule_failed, ...):
		//     returned without calling RegisterTrigger. These did early-exit as intended.
		//     High p99 here means Receive itself was slow—typically p.mu wait and/or
		//     time blocked acquiring an executor slot (schedule_failed), not DB writes.
		recordRegisterDuration := true
		defer func() {
			if recordRegisterDuration {
				p.recordMethodDuration(ctx, p.metrics.registerTriggerDurationMs, registerOutcome, time.Since(registerStart))
			}
		}()

		req, unmarshalErr := pb.UnmarshalTriggerRegistrationRequest(msg.Payload)
		if unmarshalErr != nil {
			p.lggr.Errorw("failed to unmarshal trigger registration request", "err", unmarshalErr)
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
		key := registrationKey{msg.CallerDonId, req.Metadata.WorkflowID, req.TriggerID}
		nowMs := time.Now().UnixMilli()
		p.mu.Lock()
		p.messageCache.Insert(key, sender, nowMs, msg.Payload)
		if existing, exists := p.registrations[key]; exists {
			p.mu.Unlock()
			if existing.registrationErr != nil {
				p.lggr.Debugw("skipping trigger registration; previous attempt failed with user error",
					"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "err", existing.registrationErr)
			} else {
				p.lggr.Debugw("trigger registration already exists", "workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID)
			}
			return
		}
		// NOTE: require 2F+1 by default, introduce different strategies later (KS-76)
		minRequired := uint32(2*callerDon.F + 1)
		// once=true: only one register attempt per registration cycle. Whichever workflow
		// node message satisfies quorum schedules RegisterTrigger; the following messages
		// from other nodes do not (Ready returns false, WasReady returns true). If that attempt fails, registration
		// must wait for the next cycle (e.g. re-quorum after messageCache.Delete on system
		// error, or a fresh registration flow).
		//
		// This did not matter when registration was synchronous—the "already exists" check
		// above returned early without causing unncessary RegisterTrigger calls. With async registration,
		// that check may pass in time, so once=true prevents duplicate RegisterTrigger calls for the same trigger.
		ready, payloads := p.messageCache.Ready(key, minRequired, nowMs-cfg.remoteConfig.RegistrationExpiry.Milliseconds(), true)
		if !ready {
			// Read cache state before unlocking; the register executor goroutine may
			// mutate messageCache (Delete on system error) concurrently under p.mu.
			peersReceived := len(p.messageCache.Peers(key))
			wasReady := p.messageCache.WasReady(key)
			p.mu.Unlock()
			if wasReady {
				p.lggr.Debugw("quorum already met; ignoring duplicate registration message",
					"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID,
					"peersReceived", peersReceived, "minRequired", minRequired)
			} else {
				p.lggr.Debugw("quorum not reached yet",
					"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID,
					"peersReceived", peersReceived, "minRequired", minRequired)
			}

			return
		}
		aggregated, aggregateErr := aggregation.AggregateModeRaw(payloads, uint32(callerDon.F+1))
		if aggregateErr != nil {
			p.mu.Unlock()
			p.lggr.Errorw("failed to aggregate trigger registrations", "workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "err", aggregateErr)
			return
		}
		unmarshaled, unmarshalErr := pb.UnmarshalTriggerRegistrationRequest(aggregated)
		if unmarshalErr != nil {
			p.mu.Unlock()
			p.lggr.Errorw("failed to unmarshal request", "err", unmarshalErr)
			return
		}
		p.mu.Unlock()

		capAttrs := metric.WithAttributes(
			attribute.String("capabilityID", p.capabilityID),
			attribute.String("callerDonID", strconv.FormatUint(uint64(key.callerDonID), 10)),
		)
		scheduleErr := p.registerExecutor.ExecuteTask(ctx, func(taskCtx context.Context) {
			callbackCh, registerErr := cfg.underlying.RegisterTrigger(taskCtx, unmarshaled)
			taskOutcome := "success"
			if registerErr != nil {
				taskOutcome = "error"
			}
			p.recordMethodDuration(taskCtx, p.metrics.registerTriggerDurationMs, taskOutcome, time.Since(registerStart))

			p.mu.Lock()
			defer p.mu.Unlock()
			if registerErr == nil {
				p.metrics.registerTriggerCounter.Add(taskCtx, 1, capAttrs, metric.WithAttributes(attribute.String("outcome", "success")))
				p.registrations[key] = &pubRegState{
					callback: callbackCh,
					request:  unmarshaled,
				}
				p.wg.Add(1)
				go p.triggerEventLoop(callbackCh, key)
				p.lggr.Debugw("updated trigger registration", "workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID)
				return
			}

			p.metrics.registerTriggerCounter.Add(taskCtx, 1, capAttrs, metric.WithAttributes(attribute.String("outcome", "error")))
			var capErr caperrors.Error
			if errors.As(registerErr, &capErr) && capErr.Origin() == caperrors.OriginUser {
				p.registrations[key] = &pubRegState{registrationErr: registerErr}
				p.lggr.Errorw("trigger registration failed with user error; will not retry",
					"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "err", registerErr)
				return
			}
			// System errors should retry. Ready(..., once=true) sets wasReady on quorum and never
			// resets it; without Delete, later messages keep getting not_ready even though we
			// never registered. Delete clears wasReady so a fresh quorum can schedule again.
			p.lggr.Errorw("trigger registration failed with system error; will retry",
				"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "err", registerErr)
			p.messageCache.Delete(key)
		})
		if scheduleErr != nil {
			p.lggr.Errorw("failed to schedule RegisterTrigger task",
				"workflowId", req.Metadata.WorkflowID, "triggerID", req.TriggerID, "err", scheduleErr)
			return
		}
		// Hand off to executor; disable defer recording—goroutine records success/error.
		recordRegisterDuration = false
	case types.MethodUnregisterTrigger:
		unregisterStart := time.Now()
		unregisterOutcome := "invalid"
		defer func() {
			p.recordMethodDuration(ctx, p.metrics.unregisterTriggerDurationMs, unregisterOutcome, time.Since(unregisterStart))
		}()

		meta := msg.GetTriggerEventMetadata()
		if meta == nil {
			p.lggr.Errorw("received unregister with nil metadata", "sender", sender)
			return
		}
		if len(meta.WorkflowIds) != 1 || len(meta.TriggerIds) != 1 {
			p.lggr.Errorw("received unregister with unexpected metadata sizes",
				"sender", sender, "workflowIdsLen", len(meta.WorkflowIds), "triggerIdsLen", len(meta.TriggerIds))
			return
		}
		callerDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			p.lggr.Errorw("received unregister from unsupported workflow DON", "callerDonId", msg.CallerDonId)
			return
		}
		if !cfg.membersCache[msg.CallerDonId][sender] {
			p.lggr.Errorw("unregister sender not a member of its workflow DON", "callerDonId", msg.CallerDonId, "sender", sender)
			return
		}

		key := registrationKey{
			callerDonID: msg.CallerDonId,
			workflowID:  meta.WorkflowIds[0],
			triggerID:   meta.TriggerIds[0],
		}

		p.mu.Lock()
		reg, exists := p.registrations[key]
		if !exists {
			p.mu.Unlock()
			return
		}
		nowMs := time.Now().UnixMilli()
		minRequired := uint32(2*callerDon.F + 1)
		p.unregisterCache.Insert(key, sender, nowMs, nil)
		ready, _ := p.unregisterCache.Ready(key, minRequired, nowMs-cfg.remoteConfig.RegistrationExpiry.Milliseconds(), true)
		if !ready {
			peersReceived := len(p.unregisterCache.Peers(key))
			if p.unregisterCache.WasReady(key) {
				p.mu.Unlock()
				p.lggr.Debugw("unregister quorum already met; ignoring duplicate message",
					"workflowID", key.workflowID, "triggerID", key.triggerID, "sender", sender,
					"peersReceived", peersReceived, "minRequired", minRequired)
			} else {
				p.mu.Unlock()
				p.lggr.Debugw("unregister quorum not reached yet",
					"workflowID", key.workflowID, "triggerID", key.triggerID, "sender", sender,
					"peersReceived", peersReceived, "minRequired", minRequired)
			}
			return
		}
		p.lggr.Debugw("unregister trigger", "workflowID", key.workflowID, "triggerID", key.triggerID,
			"callerDonId", msg.CallerDonId, "minRequired", minRequired, "sender", sender)
		delete(p.registrations, key)
		p.unregisterCache.Delete(key)
		p.messageCache.Delete(key)
		p.mu.Unlock()

		ctx, cancel := p.stopCh.NewCtx()
		err = p.cfg.Load().underlying.UnregisterTrigger(ctx, reg.request)
		if err != nil {
			unregisterOutcome = "error"
			p.lggr.Errorw("failed to unregister trigger on underlying", "workflowID", key.workflowID, "triggerID", key.triggerID, "err", err)
		} else {
			unregisterOutcome = "success"
		}
		cancel()
		p.lggr.Infow("unregistered trigger", "workflowID", key.workflowID, "triggerID", key.triggerID)
	case types.MethodTriggerEvent:
		p.lggr.Errorw("trigger request failed with error",
			"method", SanitizeLogString(msg.Method), "sender", sender, "errorMsg", SanitizeLogString(msg.ErrorMsg))
	case types.MethodTriggerEventAck:
		ackHandlerStart := time.Now()
		ackOutcome := "invalid"
		// Same two-path duration model as register: defer for synchronous Receive exits,
		// executor goroutine for success/error after quorum. See register case comment.
		recordAckDuration := true
		defer func() {
			if recordAckDuration {
				p.recordMethodDuration(ctx, p.metrics.ackEventDurationMs, ackOutcome, time.Since(ackHandlerStart))
			}
		}()

		triggerMetadata := msg.GetTriggerEventMetadata()
		if triggerMetadata == nil {
			p.lggr.Errorw("received empty trigger event ack metadata", "sender", sender)
			return
		}
		triggerEventID := triggerMetadata.TriggerEventId
		p.lggr.Debugw("received trigger event ACK", "sender", sender, "trigger event ID", triggerEventID)

		p.mu.Lock()
		callerDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			p.mu.Unlock()
			p.lggr.Errorw("received a message from unsupported workflow DON", "callerDonId", msg.CallerDonId)
			return
		}
		if !cfg.membersCache[msg.CallerDonId][sender] {
			p.mu.Unlock()
			p.lggr.Errorw("sender not a member of its workflow DON", "callerDonId", msg.CallerDonId, "sender", sender)
			return
		}
		if len(triggerMetadata.TriggerIds) != 1 {
			p.mu.Unlock()
			p.lggr.Errorw("did not receive single triggerID in ACK request", "callerDonId", msg.CallerDonId, "sender", sender, "triggerIDs", triggerMetadata.TriggerIds)
			return
		}
		triggerID := triggerMetadata.TriggerIds[0]

		key := ackKey{msg.CallerDonId, triggerEventID, triggerID}
		nowMs := time.Now().UnixMilli()
		p.ackCache.Insert(key, sender, nowMs, msg.Payload)
		minRequired := uint32(2*callerDon.F + 1)
		// false is set to <once> to return ready even after quorum
		// to add redundancy in case AckEvent fails below.
		ready, _ := p.ackCache.Ready(key, minRequired, 0, false)
		ackCount := len(p.ackCache.Peers(key))
		if !ready {
			p.mu.Unlock()
			p.lggr.Debugw("ACK quorum not reached yet",
				"triggerEventId", triggerEventID,
				"triggerID", triggerID,
				"acksReceived", ackCount,
				"minRequired", minRequired)
			return
		}

		p.mu.Unlock()

		callerDonID := msg.CallerDonId
		ackAttrs := metric.WithAttributes(
			attribute.String("capabilityID", p.capabilityID),
			attribute.String("callerDonID", strconv.FormatUint(uint64(callerDonID), 10)),
		)

		p.lggr.Infow("ACK quorum reached, forwarding to underlying trigger",
			"triggerEventId", triggerEventID,
			"triggerID", triggerID,
			"minRequired", minRequired)

		scheduleErr := p.ackExecutor.ExecuteTask(ctx, func(taskCtx context.Context) {
			p.lggr.Infow("executing AckEvent task",
				"triggerEventId", triggerEventID,
				"triggerID", triggerID,
				"callerDonId", callerDonID,
				"minRequired", minRequired)
			ackErr := cfg.underlying.AckEvent(taskCtx, triggerID, triggerEventID, p.capMethodName)
			taskOutcome := "success"
			if ackErr != nil {
				taskOutcome = "error"
			}
			p.recordMethodDuration(taskCtx, p.metrics.ackEventDurationMs, taskOutcome, time.Since(ackHandlerStart))
			if ackErr != nil {
				p.metrics.ackEventCounter.Add(taskCtx, 1, ackAttrs, metric.WithAttributes(attribute.String("outcome", "error")))
				p.lggr.Errorw("failed to AckEvent on underlying trigger capability",
					"eventID", triggerEventID, "capabilityID", p.capabilityID, "err", ackErr)
				return
			}
			p.lggr.Infow("AckEvent completed successfully",
				"triggerEventId", triggerEventID,
				"triggerID", triggerID,
				"callerDonId", callerDonID,
				"minRequired", minRequired)
			p.metrics.ackEventCounter.Add(taskCtx, 1, ackAttrs, metric.WithAttributes(attribute.String("outcome", "success")))
		})
		if scheduleErr != nil {
			p.lggr.Errorw("failed to schedule AckEvent task",
				"triggerEventId", triggerEventID, "triggerID", triggerID, "err", scheduleErr)
			return
		}
		recordAckDuration = false
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
			if cfg.remoteConfig.MessageExpiry != cleanupInterval {
				cleanupInterval = cfg.remoteConfig.MessageExpiry
				ticker.Reset(cleanupInterval)
			}
			now := time.Now().UnixMilli()

			// Registrations are no longer expired by timeout. Subscribers that
			// no longer want a registration respond to MethodTriggerRegistrationCheck
			// with MethodUnregisterTrigger, which is handled in Receive.
			p.mu.Lock()
			deleted := p.ackCache.DeleteOlderThan(now - cfg.remoteConfig.MessageExpiry.Milliseconds())
			p.mu.Unlock()

			if deleted > 0 {
				p.lggr.Debugw("cleaned expired AckCache entries", "deleted", deleted)
			}
		}
	}
}

func (p *triggerPublisher) registrationCheckLoop() {
	defer p.wg.Done()
	firstCfg := p.cfg.Load()
	if firstCfg == nil {
		p.lggr.Errorw("registrationCheckLoop started but config not set")
		return
	}
	checkInterval := firstCfg.remoteConfig.RegistrationRefresh
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			cfg := p.cfg.Load()
			if cfg.remoteConfig.RegistrationRefresh != checkInterval {
				checkInterval = cfg.remoteConfig.RegistrationRefresh
				ticker.Reset(checkInterval)
			}
			p.sendRegistrationChecks()
		}
	}
}

func (p *triggerPublisher) sendRegistrationChecks() {
	cfg := p.cfg.Load()
	if cfg == nil {
		return
	}

	now := time.Now().UnixMilli()

	p.mu.Lock()
	deletedUnreg := p.unregisterCache.DeleteOlderThan(now - cfg.remoteConfig.RegistrationExpiry.Milliseconds())
	p.mu.Unlock()

	if deletedUnreg > 0 {
		p.lggr.Debugw("cleaned expired unregister cache entries", "deleted", deletedUnreg)
	}

	p.mu.RLock()
	totalRegistrations := len(p.registrations)
	grouped := make(map[uint32][]registrationKey, len(cfg.workflowDONs))
	for key := range p.registrations {
		grouped[key.callerDonID] = append(grouped[key.callerDonID], key)
	}
	p.mu.RUnlock()

	p.lggr.Infow("sendRegistrationChecks: tick",
		"totalRegistrations", totalRegistrations,
		"nWorkflowDONs", len(grouped))

	for callerDonID, keys := range grouped {
		if len(keys) == 0 {
			continue
		}

		chunkSize := int(cfg.remoteConfig.MaxBatchSize)
		if chunkSize < 1 {
			chunkSize = commoncap.DefaultBatchSize
		}

		don, ok := cfg.workflowDONs[callerDonID]
		if !ok {
			continue
		}

		for offset := 0; offset < len(keys); offset += chunkSize {
			end := min(offset+chunkSize, len(keys))
			chunk := keys[offset:end]
			workflowIDs := make([]string, len(chunk))
			triggerIDs := make([]string, len(chunk))
			for i, key := range chunk {
				workflowIDs[i] = key.workflowID
				triggerIDs[i] = key.triggerID
			}

			msg := &types.MessageBody{
				CapabilityId:     p.capabilityID,
				CapabilityDonId:  cfg.capDonInfo.ID,
				CallerDonId:      callerDonID,
				Method:           types.MethodTriggerRegistrationCheck,
				CapabilityMethod: p.capMethodName,
				Metadata: &types.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &types.TriggerEventMetadata{
						WorkflowIds: workflowIDs,
						TriggerIds:  triggerIDs,
					},
				},
			}

			for _, peerID := range don.Members {
				err := p.dispatcher.Send(peerID, msg)
				if err != nil {
					p.lggr.Errorw("failed to send message", "donId", cfg.capDonInfo.ID, "peerId", peerID, "err", err)
				}
			}
			p.lggr.Debugw("sent batched registration check",
				"callerDonId", callerDonID,
				"chunkOffset", offset,
				"chunkLen", len(chunk),
				"totalRegistrations", len(keys),
				"nPeers", len(don.Members))
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
				p.lggr.Infow("triggerEventLoop channel closed", "workflowId", key.workflowID, "triggerID", key.triggerID)
				return
			}

			triggerEvent := response.Event
			p.lggr.Debugw("received trigger event", "workflowId", key.workflowID, "triggerID", key.triggerID, "triggerEventID", triggerEvent.ID)
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
					triggerIDs:     []string{key.triggerID},
				})
			}
		}
	}
}

func (p *triggerPublisher) enqueueForBatching(rawResponse []byte, key registrationKey, triggerEventID string) {
	// put in batching queue, group by hash(callerDonId, triggerEventID, response)
	combined := make([]byte, 0, 4+len(triggerEventID)+len(rawResponse))
	combined = binary.LittleEndian.AppendUint32(combined, key.callerDonID)
	combined = append(combined, triggerEventID...)
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
			triggerIDs:     []string{key.triggerID},
		}
		p.batchingQueue[sha] = elem
	} else {
		elem.workflowIDs = append(elem.workflowIDs, key.workflowID)
		elem.triggerIDs = append(elem.triggerIDs, key.triggerID)
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

		peersSent := 0
		for _, peerID := range cfg.workflowDONs[resp.callerDonID].Members {
			peersSent++

			p.lggr.Debugw("sending trigger event to peer",
				"peerID", peerID,
				"callerDonID", resp.callerDonID,
				"triggerEventID", resp.triggerEventID,
				"workflowIDs", workflowBatch,
				"triggerIDs", triggerBatch,
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
						WorkflowIds:    workflowBatch,
						TriggerIds:     triggerBatch,
						TriggerEventId: resp.triggerEventID,
					},
				},
			}

			err := p.dispatcher.Send(peerID, msg)
			if err != nil {
				p.lggr.Errorw("failed to send trigger event", "peerID", peerID, "err", err)
			}
		}
		p.lggr.Infow("sendBatch: event dispatched",
			"triggerEventID", resp.triggerEventID,
			"callerDonID", resp.callerDonID,
			"batchSize", len(triggerBatch),
			"peersSent", peersSent)
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

	if p.ackExecutor != nil {
		if err := p.ackExecutor.Close(); err != nil {
			return fmt.Errorf("failed to close ack executor: %w", err)
		}
	}
	if p.registerExecutor != nil {
		if err := p.registerExecutor.Close(); err != nil {
			return fmt.Errorf("failed to close register executor: %w", err)
		}
	}

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
