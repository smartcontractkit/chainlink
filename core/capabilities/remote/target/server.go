package target

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
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
	lggr         logger.Logger
	peerID       p2ptypes.PeerID
	underlying   commoncap.TargetCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo commoncap.DON
	workflowDONs map[string]commoncap.DON
	dispatcher   types.Dispatcher

	requestIDToRequest map[string]*request.ServerRequest
	requestTimeout     time.Duration

	receiveLock sync.Mutex
	stopCh      services.StopChan
	wg          sync.WaitGroup
}

var _ types.Receiver = &server{}
var _ services.Service = &server{}

func NewServer(peerID p2ptypes.PeerID, underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON,
	workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration, lggr logger.Logger) *server {
	return &server{
		underlying:   underlying,
		peerID:       peerID,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		requestIDToRequest: map[string]*request.ServerRequest{},
		requestTimeout:     requestTimeout,

		lggr:   lggr,
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
		if executeReq.Expired() {
			err := executeReq.Cancel(types.Error_TIMEOUT, "request expired")
			if err != nil {
				r.lggr.Errorw("failed to cancel request", "request", executeReq, "err", err)
			}
			delete(r.requestIDToRequest, requestID)
		}
	}
}

// Receive handles incoming messages from remote nodes and dispatches them to the corresponding request without blocking
// the client.
func (r *server) Receive(msg *types.MessageBody) {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()
	ctx, _ := r.stopCh.NewCtx()

	if msg.Method != types.MethodExecute {
		r.lggr.Errorw("received request for unsupported method type", "method", msg.Method)
		return
	}

	// A request is uniquely identified by the message id and the hash of the payload to prevent a malicious
	// actor from sending a different payload with the same message id
	messageId := GetMessageID(msg)
	hash := sha256.Sum256(msg.Payload)
	requestID := messageId + hex.EncodeToString(hash[:])

	if _, ok := r.requestIDToRequest[requestID]; !ok {
		callingDon, ok := r.workflowDONs[msg.CallerDonId]
		if !ok {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return
		}

		r.requestIDToRequest[requestID] = request.NewServerRequest(r.underlying, r.capInfo.ID, r.localDonInfo.ID, r.peerID,
			callingDon, messageId, r.dispatcher, r.requestTimeout, r.lggr)
	}

	req := r.requestIDToRequest[requestID]

	go func() {
		err := req.OnMessage(ctx, msg)
		if err != nil {
			r.lggr.Errorw("request failed to OnMessage new message", "request", req, "err", err)
		}
	}()
}

func GetMessageID(msg *types.MessageBody) string {
	return string(msg.MessageId)
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
