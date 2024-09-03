package target

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// server manages all external users of a local target capability.
// Its responsibilities are:
//  1. Manage requests from external nodes executing the target capability once sufficient requests are received.
//  2. Send out responses produced by an underlying capability to all requesters.
//
// server communicates with corresponding client on remote nodes.
type server struct {
	services.StateMachine
	lggr logger.Logger

	config       *commoncap.RemoteTargetConfig
	peerID       p2ptypes.PeerID
	underlying   commoncap.TargetCapability
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
}

var _ types.Receiver = &server{}
var _ services.Service = &server{}

type requestAndMsgID struct {
	request   *request.ServerRequest
	messageID string
}

func NewServer(config *commoncap.RemoteTargetConfig, peerID p2ptypes.PeerID, underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON,
	workflowDONs map[uint32]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration, lggr logger.Logger) *server {
	if config == nil {
		lggr.Info("no config provided, using default values")
		config = &commoncap.RemoteTargetConfig{}
	}
	return &server{
		config:       config,
		underlying:   underlying,
		peerID:       peerID,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		requestIDToRequest:         map[string]requestAndMsgID{},
		messageIDToRequestIDsCount: map[string]map[string]int{},
		requestTimeout:             requestTimeout,

		lggr:   lggr.Named("TargetServer"),
		stopCh: make(services.StopChan),
	}
}

func (r *server) Start(ctx context.Context) error {
	return r.StartOnce(r.Name(), func() error {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			ticker := time.NewTicker(r.requestTimeout)
			defer ticker.Stop()
			r.lggr.Info("TargetServer started")
			for {
				select {
				case <-r.stopCh:
					return
				case <-ticker.C:
					r.expireRequests()
				}
			}
		}()
		return nil
	})
}

func (r *server) Close() error {
	return r.StopOnce(r.Name(), func() error {
		close(r.stopCh)
		r.wg.Wait()
		r.lggr.Info("TargetServer closed")
		return nil
	})
}

func (r *server) expireRequests() {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()

	for requestID, executeReq := range r.requestIDToRequest {
		if executeReq.request.Expired() {
			err := executeReq.request.Cancel(types.Error_TIMEOUT, "request expired")
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

	if msg.Method != types.MethodExecute {
		r.lggr.Errorw("received request for unsupported method type", "method", remote.SanitizeLogString(msg.Method))
		return
	}

	messageId, err := GetMessageID(msg)
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
	requestID := messageId + hex.EncodeToString(msgHash[:])

	r.lggr.Debugw("received request", "msgId", msg.MessageId, "requestID", requestID)

	if requestIDs, ok := r.messageIDToRequestIDsCount[messageId]; ok {
		requestIDs[requestID] = requestIDs[requestID] + 1
	} else {
		r.messageIDToRequestIDsCount[messageId] = map[string]int{requestID: 1}
	}

	requestIDs := r.messageIDToRequestIDsCount[messageId]
	if len(requestIDs) > 1 {
		// This is a potential attack vector as well as a situation that will occur if the client is sending non-deterministic payloads
		// so a warning is logged
		r.lggr.Warnw("received messages with the same id and different payloads", "messageID", messageId, "lenRequestIDs", len(requestIDs))
	}

	if _, ok := r.requestIDToRequest[requestID]; !ok {
		callingDon, ok := r.workflowDONs[msg.CallerDonId]
		if !ok {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return
		}

		r.requestIDToRequest[requestID] = requestAndMsgID{
			request: request.NewServerRequest(r.underlying, r.capInfo.ID, r.localDonInfo.ID, r.peerID,
				callingDon, messageId, r.dispatcher, r.requestTimeout, r.lggr),
			messageID: messageId,
		}
	}

	reqAndMsgID := r.requestIDToRequest[requestID]

	err = reqAndMsgID.request.OnMessage(ctx, msg)
	if err != nil {
		r.lggr.Errorw("request failed to OnMessage new message", "messageID", reqAndMsgID.messageID, "err", err)
	}
}

func (r *server) getMessageHash(msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	for _, path := range r.config.RequestHashExcludedAttributes {
		if !req.Inputs.DeleteAtPath(path) {
			return [32]byte{}, fmt.Errorf("failed to delete attribute from map at path: %s", path)
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
		return "", fmt.Errorf("invalid message id")
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
	return "TargetServer"
}
