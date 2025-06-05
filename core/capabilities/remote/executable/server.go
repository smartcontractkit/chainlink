package executable

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
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
	lggr logger.Logger

	config       *commoncap.RemoteExecutableConfig
	peerID       p2ptypes.PeerID
	underlying   commoncap.ExecutableCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo commoncap.DON
	workflowDONs map[uint32]commoncap.DON
	dispatcher   types.Dispatcher

	requestIDToRequest map[string]requestAndMsgID
	requestTimeout     time.Duration

	// Used to detect messages with the same message id but different payloads
	messageIDToRequestIDsCount map[string]map[string]int

	receiveLock sync.Mutex
	stopCh      services.StopChan
	wg          sync.WaitGroup

	parallelExecutor *parallelExecutor
}

var _ types.Receiver = &server{}
var _ services.Service = &server{}

type requestAndMsgID struct {
	request   *request.ServerRequest
	messageID string
}

func NewServer(remoteExecutableConfig *commoncap.RemoteExecutableConfig, peerID p2ptypes.PeerID, underlying commoncap.ExecutableCapability,
	capInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON,
	workflowDONs map[uint32]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration,
	maxParallelRequests int,
	lggr logger.Logger) *server {
	if remoteExecutableConfig == nil {
		lggr.Info("no remote config provided, using default values")
		remoteExecutableConfig = &commoncap.RemoteExecutableConfig{}
	}
	return &server{
		config:       remoteExecutableConfig,
		underlying:   underlying,
		peerID:       peerID,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		requestIDToRequest:         map[string]requestAndMsgID{},
		messageIDToRequestIDsCount: map[string]map[string]int{},
		requestTimeout:             requestTimeout,

		lggr:   logger.Named(lggr, "ExecutableCapabilityServer"),
		stopCh: make(services.StopChan),

		parallelExecutor: newParallelExecutor(maxParallelRequests),
	}
}

func (r *server) Start(ctx context.Context) error {
	return r.StartOnce(r.Name(), func() error {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			tickerInterval := expiryCheckInterval
			if r.requestTimeout < tickerInterval {
				tickerInterval = r.requestTimeout
			}
			ticker := time.NewTicker(tickerInterval)
			defer ticker.Stop()

			r.lggr.Info("executable capability server started")
			for {
				select {
				case <-r.stopCh:
					return
				case <-ticker.C:
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

func (r *server) Close() error {
	return r.StopOnce(r.Name(), func() error {
		close(r.stopCh)
		r.wg.Wait()
		err := r.parallelExecutor.Close()
		if err != nil {
			return fmt.Errorf("failed to close parallel executor: %w", err)
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

	msgHash, err := r.getMessageHash(msg)
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
		callingDon, ok := r.workflowDONs[msg.CallerDonId]
		if !ok {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return
		}

		sr, ierr := request.NewServerRequest(r.underlying, msg.Method, r.capInfo.ID, r.localDonInfo.ID, r.peerID,
			callingDon, messageID, r.dispatcher, r.requestTimeout, r.lggr)
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

func (r *server) getMessageHash(msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	// An attribute called StepDependency is used to define a data dependency between steps,
	// and not to provide input values; we should therefore disregard it when hashing the request
	if len(r.config.RequestHashExcludedAttributes) == 0 {
		r.config.RequestHashExcludedAttributes = []string{"StepDependency"}
	}

	for _, path := range r.config.RequestHashExcludedAttributes {
		if req.Inputs != nil {
			req.Inputs.DeleteAtPath(path)
		}
	}

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	hash := sha256.Sum256(reqBytes)
	return hash, nil
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
