package executable

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

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

	parallelExecutor *parallelExecutor
}

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

func NewServer(capabilityID, methodName string, peerID p2ptypes.PeerID, dispatcher types.Dispatcher, lggr logger.Logger) *server {
	return &server{
		capabilityID:               capabilityID,
		capMethodName:              methodName,
		peerID:                     peerID,
		dispatcher:                 dispatcher,
		lggr:                       logger.Named(lggr, "ExecutableCapabilityServer"),
		requestIDToRequest:         map[string]requestAndMsgID{},
		messageIDToRequestIDsCount: map[string]map[string]int{},
		stopCh:                     make(services.StopChan),
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
		r.lggr.Warn("no message hasher provided, using default V1 hasher")
		messageHasher = NewV1Hasher(remoteExecutableConfig.RequestHashExcludedAttributes)
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
		r.parallelExecutor = newParallelExecutor(int(cfg.remoteExecutableConfig.ServerMaxParallelRequests))

		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			ticker := time.NewTicker(getServerTickerInterval(cfg))
			defer ticker.Stop()

			r.lggr.Info("executable capability server started")
			for {
				select {
				case <-r.stopCh:
					return
				case <-ticker.C:
					ticker.Reset(getServerTickerInterval(cfg))
					r.expireRequests()
				}
			}
		}()

		err := r.parallelExecutor.Start(ctx)
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

		r.lggr.Info("executable capability server closed")
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

	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

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

	msgHash, err := cfg.hasher.Hash(msg)
	if err != nil {
		r.lggr.Errorw("failed to get message hash", "err", err)
		return
	}

	// A request is uniquely identified by the message id and the hash of the payload to prevent a malicious
	// actor from sending a different payload with the same message id
	requestID := messageID + hex.EncodeToString(msgHash[:])

	r.lggr.Debugw("received request", "msgId", msg.MessageId, "requestID", requestID)

	if requestIDs, ok := r.messageIDToRequestIDsCount[messageID]; ok {
		requestIDs[requestID]++
	} else {
		r.messageIDToRequestIDsCount[messageID] = map[string]int{requestID: 1}
	}

	requestIDs := r.messageIDToRequestIDsCount[messageID]
	if len(requestIDs) > 1 {
		// This is a potential attack vector as well as a situation that will occur if the client is sending non-deterministic payloads
		// so a warning is logged
		r.lggr.Warnw("received messages with the same id and different payloads", "messageID", messageID, "lenRequestIDs", len(requestIDs))
	}

	if _, ok := r.requestIDToRequest[requestID]; !ok {
		callingDon, ok := cfg.workflowDONs[msg.CallerDonId]
		if !ok {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return
		}

		sr, ierr := request.NewServerRequest(cfg.underlying, msg.Method, cfg.capInfo.ID, cfg.localDonInfo.ID, r.peerID,
			callingDon, messageID, r.dispatcher, cfg.remoteExecutableConfig.RequestTimeout, r.capMethodName, r.lggr)
		if ierr != nil {
			r.lggr.Errorw("failed to instantiate server request", "err", ierr)
			return
		}

		r.requestIDToRequest[requestID] = requestAndMsgID{
			request:   sr,
			messageID: messageID,
		}
	}

	reqAndMsgID := r.requestIDToRequest[requestID]
	if executeTaskErr := r.parallelExecutor.ExecuteTask(ctx,
		func(ctx context.Context) {
			err = reqAndMsgID.request.OnMessage(ctx, msg)
			if err != nil {
				r.lggr.Errorw("failed to execute on message", "messageID", reqAndMsgID.messageID, "err", err)
			}
		}); executeTaskErr != nil {
		r.lggr.Errorw("failed to execute on message task", "messageID", messageID, "err", executeTaskErr)
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
