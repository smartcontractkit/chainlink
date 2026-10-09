package executable

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// server manages all external users of a local executable capability.
// Its responsibilities are:
//  1. Manage requests from external nodes executing the executable capability once sufficient requests are received.
//  2. Send out responses produced by an underlying capability to all requesters.
//
// server communicates with corresponding client on remote nodes.
type server struct {
	services.StateMachine
	capabilityID  string
	capMethodName string
	peerID        p2ptypes.PeerID
	dispatcher    types.Dispatcher
	cfg           atomic.Pointer[dynamicServerConfig]
	lggr          logger.Logger

	requestIDToRequest map[string]requestAndMsgID

	// Used to detect messages with the same message id but different payloads
	messageIDToRequestIDsCount map[string]map[string]int

	receiveLock sync.Mutex
	stopCh      services.StopChan
	wg          sync.WaitGroup

	parallelExecutor        *remote.ParallelExecutor
	rejectedRequestsCounter metric.Int64Counter

	// Messages that could not get an executor slot immediately wait for one in their own goroutine,
	// never in the dispatcher's receive goroutine. The number of such waiting messages is capped.
	waitingForSlot    atomic.Int64
	maxWaitingForSlot int64

	// workflowDONBindingGate, when open, makes each ServerRequest require the
	// request's Metadata.WorkflowDonID to match the authenticated calling DON.
	workflowDONBindingGate limits.GateLimiter
}

const (
	// errServerAtCapacity is sent to requesters whose request was rejected because no executor slot became available.
	errServerAtCapacity = "executable capability server at capacity, request rejected"

	// maxWaitingForSlotPerSlot bounds the number of messages waiting for an executor slot, relative to the number of slots.
	maxWaitingForSlotPerSlot = 10

	rejectReasonWaitQueueFull = "wait_queue_full"
	rejectReasonWaitTimeout   = "wait_timeout"
)

type dynamicServerConfig struct {
	remoteExecutableConfig *commoncap.RemoteExecutableConfig
	hasher                 types.MessageHasher
	underlying             commoncap.ExecutableCapability
	capInfo                commoncap.CapabilityInfo
	localDonInfo           commoncap.DON
	workflowDONs           map[uint32]commoncap.DON
}

type Server interface {
	types.Receiver
	services.Service
	SetConfig(remoteExecutableConfig *commoncap.RemoteExecutableConfig, underlying commoncap.ExecutableCapability,
		capInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON,
		messageHasher types.MessageHasher) error
}

var _ Server = &server{}
var _ types.Receiver = &server{}
var _ services.Service = &server{}

type requestAndMsgID struct {
	request   *request.ServerRequest
	messageID string
}

func NewServer(capabilityID, methodName string, peerID p2ptypes.PeerID, dispatcher types.Dispatcher, workflowDONBindingGate limits.GateLimiter, lggr logger.Logger) *server {
	return &server{
		capabilityID:               capabilityID,
		capMethodName:              methodName,
		peerID:                     peerID,
		dispatcher:                 dispatcher,
		lggr:                       logger.With(logger.Named(lggr, "ExecutableCapabilityServer"), "capabilityID", capabilityID, "capMethodName", methodName),
		requestIDToRequest:         map[string]requestAndMsgID{},
		messageIDToRequestIDsCount: map[string]map[string]int{},
		stopCh:                     make(services.StopChan),
		workflowDONBindingGate:     workflowDONBindingGate,
	}
}

// SetConfig sets the remote server configuration dynamically
func (r *server) SetConfig(remoteExecutableConfig *commoncap.RemoteExecutableConfig, underlying commoncap.ExecutableCapability,
	capInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON, messageHasher types.MessageHasher) error {
	currCfg := r.cfg.Load()
	if remoteExecutableConfig == nil {
		r.lggr.Info("no remote config provided, using default values")
		remoteExecutableConfig = &commoncap.RemoteExecutableConfig{}
	}
	if messageHasher == nil {
		return errors.New("message hasher must be provided")
	}
	if capInfo.ID == "" || capInfo.ID != r.capabilityID {
		return fmt.Errorf("capability info provided does not match the server's capabilityID: %s != %s", capInfo.ID, r.capabilityID)
	}
	if underlying == nil {
		return errors.New("underlying capability cannot be nil")
	}
	if len(localDonInfo.Members) == 0 {
		return errors.New("empty localDonInfo provided")
	}
	if len(workflowDONs) == 0 {
		return errors.New("empty workflowDONs provided")
	}
	if remoteExecutableConfig.RequestTimeout <= 0 {
		return errors.New("cfg.RequestTimeout must be positive")
	}
	if remoteExecutableConfig.ServerMaxParallelRequests <= 0 {
		return errors.New("cfg.ServerMaxParallelRequests must be positive")
	}

	if currCfg != nil && currCfg.remoteExecutableConfig != nil &&
		currCfg.remoteExecutableConfig.ServerMaxParallelRequests > 0 &&
		remoteExecutableConfig.ServerMaxParallelRequests != currCfg.remoteExecutableConfig.ServerMaxParallelRequests {
		r.lggr.Warn("ServerMaxParallelRequests changed but it won't be applied until node restart")
	}

	// always replace the whole dynamicServerConfig object to avoid inconsistent state
	r.cfg.Store(&dynamicServerConfig{
		remoteExecutableConfig: remoteExecutableConfig,
		hasher:                 messageHasher,
		underlying:             underlying,
		capInfo:                capInfo,
		localDonInfo:           localDonInfo,
		workflowDONs:           workflowDONs,
	})
	return nil
}

func (r *server) Start(ctx context.Context) error {
	return r.StartOnce(r.Name(), func() error {
		cfg := r.cfg.Load()

		// Validate that all required fields are set before starting
		if cfg == nil {
			return errors.New("config not set - call SetConfig() before Start()")
		}
		if cfg.remoteExecutableConfig == nil {
			return errors.New("remote executable config not set - call SetConfig() before Start()")
		}
		if cfg.underlying == nil {
			return errors.New("underlying capability not set - call SetConfig() before Start()")
		}
		if cfg.capInfo.ID == "" {
			return errors.New("capability info not set - call SetConfig() before Start()")
		}
		if len(cfg.localDonInfo.Members) == 0 {
			return errors.New("local DON info not set - call SetConfig() before Start()")
		}
		if cfg.remoteExecutableConfig.RequestTimeout <= 0 {
			return errors.New("cfg.RequestTimeout not set - call SetConfig() before Start()")
		}
		if cfg.remoteExecutableConfig.ServerMaxParallelRequests <= 0 {
			return errors.New("cfg.ServerMaxParallelRequests not set - call SetConfig() before Start()")
		}
		if r.dispatcher == nil {
			return errors.New("dispatcher set to nil, cannot start server")
		}

		// Initialize parallel executor with the configured max parallel requests
		slotUsageAttrs := []attribute.KeyValue{
			attribute.String("capabilityID", r.capabilityID),
			attribute.String("capMethodName", r.capMethodName),
		}
		r.parallelExecutor = remote.NewParallelExecutor(int(cfg.remoteExecutableConfig.ServerMaxParallelRequests), "executable_server", slotUsageAttrs...)
		r.maxWaitingForSlot = int64(cfg.remoteExecutableConfig.ServerMaxParallelRequests) * maxWaitingForSlotPerSlot

		r.wg.Go(func() {
			ticker := time.NewTicker(getServerTickerInterval(cfg))
			defer ticker.Stop()

			r.lggr.Info("ExecutableCapabilityServer started")
			for {
				select {
				case <-r.stopCh:
					return
				case <-ticker.C:
					ticker.Reset(getServerTickerInterval(cfg))
					r.expireRequests()
				}
			}
		})

		var err error
		r.rejectedRequestsCounter, err = beholder.GetMeter().Int64Counter("platform_executable_capability_server_rejected_request_count")
		if err != nil {
			return fmt.Errorf("failed to register platform_executable_capability_server_rejected_request_count: %w", err)
		}

		err = r.parallelExecutor.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start parallel executor: %w", err)
		}
		return nil
	})
}

func getServerTickerInterval(cfg *dynamicServerConfig) time.Duration {
	if cfg.remoteExecutableConfig.RequestTimeout > 0 {
		return cfg.remoteExecutableConfig.RequestTimeout
	}
	return defaultExpiryCheckInterval
}

func (r *server) Close() error {
	return r.StopOnce(r.Name(), func() error {
		close(r.stopCh)
		r.wg.Wait()
		if r.parallelExecutor != nil {
			err := r.parallelExecutor.Close()
			if err != nil {
				return fmt.Errorf("failed to close parallel executor: %w", err)
			}
		}

		r.lggr.Info("ExecutableCapabilityServer closed")
		return nil
	})
}

func (r *server) expireRequests() {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	for requestID, executeReq := range r.requestIDToRequest {
		if executeReq.request.Expired() {
			ctx, cancelFn := r.stopCh.NewCtx()
			err := executeReq.request.Cancel(ctx, types.Error_TIMEOUT, "request expired by executable server")
			cancelFn()
			if err != nil {
				r.lggr.Errorw("failed to cancel request", "request", executeReq, "err", err)
			}
		}
		if executeReq.request.Evictable(commoncap.DefaultExecutableRequestTimeout) {
			delete(r.requestIDToRequest, requestID)
			delete(r.messageIDToRequestIDsCount, executeReq.messageID)
		}
	}
}

func (r *server) Receive(ctx context.Context, msg *types.MessageBody) {
	cfg := r.cfg.Load()
	if cfg == nil {
		r.lggr.Errorw("config not set, cannot process request")
		return
	}

	switch msg.Method {
	case types.MethodExecute:
	default:
		r.lggr.Errorw("received request for unsupported method type", "method", remote.SanitizeLogString(msg.Method))
		return
	}

	messageID, err := GetMessageID(msg)
	if err != nil {
		r.lggr.Errorw("invalid message id", "err", err, "id", remote.SanitizeLogString(string(msg.MessageId)))
		return
	}

	msgHash, err := cfg.hasher.Hash(ctx, msg)
	if err != nil {
		r.lggr.Errorw("failed to get message hash", "err", err)
		return
	}

	// A request is uniquely identified by the message id and the hash of the payload to prevent a malicious
	// actor from sending a different payload with the same message id
	requestID := messageID + hex.EncodeToString(msgHash[:])

	r.lggr.Debugw("received request", "msgId", msg.MessageId, "requestID", requestID)

	reqAndMsgID, ok := r.getOrCreateRequest(cfg, msg, messageID, requestID)
	if !ok {
		return
	}

	onMessage := func(ctx context.Context) {
		if err := reqAndMsgID.request.OnMessage(ctx, msg); err != nil {
			r.lggr.Errorw("failed to execute on message", "messageID", reqAndMsgID.messageID, "err", err)
		}
	}

	// Once another message has claimed the execution, OnMessage only records the requester and fans out
	// the response (if any), so it is handled inline without taking an executor slot.
	if reqAndMsgID.request.ExecutionClaimed() {
		onMessage(ctx)
		return
	}

	// Never block waiting for an executor slot: Receive is called from the dispatcher's single receive
	// goroutine for this capability, so blocking here stalls all inbound messages and causes drops.
	executeTaskErr := r.parallelExecutor.TryExecuteTask(ctx, onMessage)
	if errors.Is(executeTaskErr, remote.ErrNoSlotAvailable) {
		r.waitForSlot(ctx, cfg, msg, reqAndMsgID, onMessage)
		return
	}
	if executeTaskErr != nil {
		r.lggr.Errorw("failed to execute on message task", "messageID", messageID, "err", executeTaskErr)
	}
}

// waitForSlot waits for an executor slot in a separate goroutine until the request deadline and then
// runs onMessage. If too many messages are already waiting, or no slot frees up before the deadline,
// the message is rejected so that the requester gets an error instead of a silent timeout.
func (r *server) waitForSlot(ctx context.Context, cfg *dynamicServerConfig, msg *types.MessageBody, reqAndMsgID requestAndMsgID, onMessage func(context.Context)) {
	if r.waitingForSlot.Add(1) > r.maxWaitingForSlot {
		r.waitingForSlot.Add(-1)
		r.rejectMessage(ctx, cfg, msg, reqAndMsgID.messageID, rejectReasonWaitQueueFull)
		return
	}

	started := r.IfNotStopped(func() {
		r.wg.Go(func() {
			defer r.waitingForSlot.Add(-1)

			waitCtx, cancel := r.stopCh.Ctx(ctx)
			defer cancel()
			waitCtx, cancelDeadline := context.WithDeadline(waitCtx, reqAndMsgID.request.Deadline())
			defer cancelDeadline()

			err := r.parallelExecutor.ExecuteTaskWithWaitContext(waitCtx, ctx, onMessage)
			switch {
			case err == nil:
			case errors.Is(err, context.DeadlineExceeded):
				r.rejectMessage(ctx, cfg, msg, reqAndMsgID.messageID, rejectReasonWaitTimeout)
			default:
				r.lggr.Debugw("stopped waiting for executor slot", "messageID", reqAndMsgID.messageID, "err", err)
			}
		})
	})
	if !started {
		r.waitingForSlot.Add(-1)
	}
}

// getOrCreateRequest records the message in the request bookkeeping maps and returns the request it
// belongs to, creating it if needed. receiveLock is held only for this bookkeeping and never across
// blocking calls, so that expireRequests can always make progress.
func (r *server) getOrCreateRequest(cfg *dynamicServerConfig, msg *types.MessageBody, messageID, requestID string) (requestAndMsgID, bool) {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	requestIDs, requestIDsOK := r.messageIDToRequestIDsCount[messageID]
	if requestIDsOK {
		requestIDs[requestID]++
		if len(requestIDs) > 1 {
			// This is a potential attack vector as well as a situation that will occur if the client is sending non-deterministic payloads
			// so a warning is logged
			r.lggr.Warnw("received messages with the same id and different payloads", "messageID", messageID, "lenRequestIDs", len(requestIDs))
		}
	}

	if _, ok := r.requestIDToRequest[requestID]; !ok {
		callingDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return requestAndMsgID{}, false
		}

		sr, ierr := request.NewServerRequest(cfg.underlying, msg.Method, cfg.capInfo.ID, cfg.localDonInfo.ID, r.peerID,
			callingDon, messageID, r.dispatcher, cfg.remoteExecutableConfig.RequestTimeout, r.capMethodName, r.workflowDONBindingGate, r.lggr)
		if ierr != nil {
			r.lggr.Errorw("failed to instantiate server request", "err", ierr)
			return requestAndMsgID{}, false
		}

		r.requestIDToRequest[requestID] = requestAndMsgID{
			request:   sr,
			messageID: messageID,
		}
	}

	if !requestIDsOK {
		// This is separate from inverse case because we want to wait until after the early returns in between.
		r.messageIDToRequestIDsCount[messageID] = map[string]int{requestID: 1}
	}

	return r.requestIDToRequest[requestID], true
}

// rejectMessage responds to the sender with a TIMEOUT error without executing the request. It does not
// touch the shared ServerRequest, so execution on behalf of other requesters is unaffected. The error
// message is constant so that clients can aggregate identical errors from multiple overloaded nodes.
func (r *server) rejectMessage(ctx context.Context, cfg *dynamicServerConfig, msg *types.MessageBody, messageID string, reason string) {
	r.lggr.Warnw("no executor slot available, rejecting request", "messageID", messageID, "reason", reason,
		"maxParallelRequests", cfg.remoteExecutableConfig.ServerMaxParallelRequests)
	r.rejectedRequestsCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("capabilityID", r.capabilityID),
		attribute.String("capMethodName", r.capMethodName),
		attribute.String("reason", reason),
	))

	requester, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		r.lggr.Errorw("failed to convert message sender to PeerID", "messageID", messageID, "err", err)
		return
	}
	callingDon, ok := cfg.workflowDONs[msg.CallerDonId]
	if !ok || !slices.Contains(callingDon.Members, requester) {
		r.lggr.Errorw("rejected request from peer not in calling don", "messageID", messageID, "donId", msg.CallerDonId, "peer", requester)
		return
	}

	err = r.dispatcher.Send(requester, &types.MessageBody{
		CapabilityId:     cfg.capInfo.ID,
		CapabilityDonId:  cfg.localDonInfo.ID,
		CallerDonId:      msg.CallerDonId,
		Method:           types.MethodExecute,
		MessageId:        []byte(messageID),
		Sender:           r.peerID[:],
		Receiver:         requester[:],
		CapabilityMethod: r.capMethodName,
		Error:            types.Error_TIMEOUT,
		ErrorMsg:         errServerAtCapacity,
	})
	if err != nil {
		r.lggr.Errorw("failed to send rejection response", "messageID", messageID, "err", err)
	}
}

func GetMessageID(msg *types.MessageBody) (string, error) {
	idStr := string(msg.MessageId)
	if !validation.IsValidID(idStr) {
		return "", errors.New("invalid message id")
	}
	return idStr, nil
}

func (r *server) Ready() error {
	return nil
}

func (r *server) HealthReport() map[string]error {
	return nil
}

func (r *server) Name() string {
	return r.lggr.Name()
}
