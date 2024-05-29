package target

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type receiverRequest interface {
	Receive(ctx context.Context, msg *types.MessageBody) error
	Expired() bool
	Cancel(err types.Error, msg string) error
}

type remoteTargetReceiver struct {
	lggr         logger.Logger
	peerID       p2ptypes.PeerID
	underlying   commoncap.TargetCapability
	capInfo      commoncap.CapabilityInfo
	localDonInfo capabilities.DON
	workflowDONs map[string]commoncap.DON
	dispatcher   types.Dispatcher

	requestIDToRequest map[string]receiverRequest
	requestTimeout     time.Duration

	receiveLock sync.Mutex
}

var _ types.Receiver = &remoteTargetReceiver{}

func NewRemoteTargetReceiver(ctx context.Context, lggr logger.Logger, peerID p2ptypes.PeerID, underlying commoncap.TargetCapability, capInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON,
	workflowDONs map[string]commoncap.DON, dispatcher types.Dispatcher, requestTimeout time.Duration) *remoteTargetReceiver {

	receiver := &remoteTargetReceiver{
		underlying:   underlying,
		peerID:       peerID,
		capInfo:      capInfo,
		localDonInfo: localDonInfo,
		workflowDONs: workflowDONs,
		dispatcher:   dispatcher,

		requestIDToRequest: map[string]receiverRequest{},
		requestTimeout:     requestTimeout,

		lggr: lggr,
	}

	go func() {
		ticker := time.NewTicker(requestTimeout)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				receiver.ExpireRequests()
			}
		}
	}()

	return receiver
}

func (r *remoteTargetReceiver) ExpireRequests() {
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

func (r *remoteTargetReceiver) Receive(msg *types.MessageBody) {
	r.receiveLock.Lock()
	defer r.receiveLock.Unlock()
	// TODO should the dispatcher be passing in a context?
	ctx := context.Background()

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
		if callingDon, ok := r.workflowDONs[msg.CallerDonId]; ok {
			r.requestIDToRequest[requestID] = request.NewReceiverRequest(r.underlying, r.capInfo.ID, r.localDonInfo.ID, r.peerID,
				callingDon, messageId, r.dispatcher, r.requestTimeout)
		} else {
			r.lggr.Errorw("received request from unregistered don", "donId", msg.CallerDonId)
			return
		}
	}

	req := r.requestIDToRequest[requestID]

	go func() {
		err := req.Receive(ctx, msg)
		if err != nil {
			r.lggr.Errorw("request failed to Receive new message", "request", req, "err", err)
		}
	}()
}

func GetMessageID(msg *types.MessageBody) string {
	return string(msg.MessageId)
}
