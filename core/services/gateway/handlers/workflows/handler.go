package workflows

import (
	"context"
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

var (
	// TODO: consolidate error handling with functions
	ErrUnsupportedMethod = errors.New("unsupported method")

	promHandlerError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_workflows_handler_error",
		Help: "Metric to track functions handler errors",
	}, []string{"don_id", "error"})

	promCreateSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_workflows_create_success",
		Help: "Metric to track successful create calls",
	}, []string{"don_id"})

	promCreateFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_workflows_create_failure",
		Help: "Metric to track failed create calls",
	}, []string{"don_id"})

	promDeleteSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_workflows_delete_success",
		Help: "Metric to track successful delete calls",
	}, []string{"don_id"})

	promDeleteFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_workflows_delete_failure",
		Help: "Metric to track failed delete calls",
	}, []string{"don_id"})
)

// TODO: deduplicate with functions
type PendingRequest struct {
	request    *api.Message
	responses  map[string]*api.Message
	successful []*api.Message
	errors     []*FailedRequest
}

type FailedRequest struct {
	request *api.Message
	error   error
}

// DummyHandler forwards each request/response without doing any checks.
type workflowSpecHandler struct {
	donConfig       *config.DONConfig
	don             handlers.DON
	savedCallbacks  map[string]*savedCallback
	pendingRequests hc.RequestCache[PendingRequest]

	mu   sync.Mutex
	lggr logger.Logger

	jobORM job.ORM
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
}

var _ handlers.Handler = (*workflowSpecHandler)(nil)

func NewWorkflowSpecHandler(donConfig *config.DONConfig, don handlers.DON, jobORM job.ORM, ds sqlutil.DataSource, lggr logger.Logger) (handlers.Handler, error) {
	return &workflowSpecHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           lggr.Named("DummyHandler." + donConfig.DonId),
		jobORM: 	    jobORM,
	}, nil
}

func (h *workflowSpecHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {

	mt := msg.Body.Method
	switch mt {
	case Create:
		return h.handleCreate(ctx, msg, callbackCh)
	case Delete:
		return h.handleDelete(ctx, msg, callbackCh)
	case Commit:
		return h.handleCommit(ctx, msg, callbackCh)
	case Abort:
		return h.handleAbort(ctx, msg, callbackCh)
	default:
		// TODO consolidate error handling
		return ErrUnsupportedMethod
	}

}

func (h *workflowSpecHandler) forward(ctx context.Context, msg *api.Message) error {
	h.mu.Lock()
	don := h.don
	h.mu.Unlock()

	var err error
	// Send to all nodes.
	for _, member := range h.donConfig.Members {
		err = errors.Join(err, don.SendToNode(ctx, member.Address, msg))
	}
}

func (h *workflowSpecHandler) handleCreate(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.lggr.Debugw("Handling request", "msg", msg)

	/*
		h.mu.Lock()
		h.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
		don := h.don
		h.mu.Unlock()

		var err error
		// Send to all nodes.
		for _, member := range h.donConfig.Members {
			err = errors.Join(err, don.SendToNode(ctx, member.Address, msg))
		}
	*/
	// prepare state change and forward to all nodes
	j, err := extractJob(msg)
	if err != nil {
		return fmt.Errorf("failed to extract job: %w", err)
	}
	h.lggr.Debugw("handleRequest: processing message", "sender", msg.Body.Sender, "messageId", msg.Body.MessageId)
	err := h.pendingRequests.NewRequest(msg, callbackCh, &PendingRequest{request: msg, responses: make(map[string]*api.Message)})
	if err != nil {
		h.lggr.Warnw("handleRequest: error adding new request", "sender", msg.Body.Sender, "err", err)
		promHandlerError.WithLabelValues(h.donConfig.DonId, err.Error()).Inc()
		return err
	}
	// greedily create the job. this is bit fraught because we don't know if the job will be committed by the other nodes
	// and we need all the nodes to end up in the same state. maybe a two phase commit or a more sophisticated state machine via s4
	err = h.jobORM.CreateJob()
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	h.forward(ctx, msg)
}

func (h *workflowSpecHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {

	h.lggr.Debugw("HandleNodeMessage: processing message", "nodeAddr", nodeAddr, "receiver", msg.Body.Receiver, "id", msg.Body.MessageId)

	switch msg.Body.Method {
	case Create:
		return h.pendingRequests.ProcessResponse(msg, h.processCreate)
	case Delete:
		return h.pendingRequests.ProcessResponse(msg, h.processDelete)
/*
	case Abort:
		return h.pendingRequests.ProcessResponse(msg, h.processAbort)
	case Commit:
		return h.pendingRequests.ProcessResponse(msg, h.processCommit)
*/
	default:
		h.lggr.Debugw("unsupported method", "method", msg.Body.Method)
		return ErrUnsupportedMethod
	}
}

//// Conforms to ResponseProcessor[*PendingRequest]
//func (h *functionsHandler) processSecretsResponse(response *api.Message, responseData *PendingRequest) (*handlers.UserCallbackPayload, *PendingRequest, error) {

func (h *workflowSpecHandler) processCreate(peerMsg *api.Message, pendingState *PendingRequest) (*handlers.UserCallbackPayload, *PendingRequest, error) {
	h.lggr.Debugw("processCreate: processing message", "messageId", msg.Body.MessageId)
	if _, exists := pendingState.responses[peerMsg.Body.Sender]; exists {
		return nil, nil, errors.New("duplicate response")
	}
	if peerMsg.Body.Method != pendingState.request.Body.Method {
		return nil, pendingState, errors.New("invalid method")
	}
	var responsePayload ResponseBase
	err := json.Unmarshal(peerMsg.Body.Payload, &responsePayload)
	if err != nil {
		pendingState.errors = append(pendingState.errors, &FailedRequest{request: peerMsg, error: fmt.Errorf("internal error: failed to unmarshal peer message: %w", err)})
		return nil, pendingState, err
	}

	pr.responses[msg.Body.Sender] = msg
	// if we have any errors, we should abort

	if len(pr.responses) == len(h.donConfig.Members) {
		// all nodes have responded
	}
}
func (h *workflowSpecHandler) Start(context.Context) error {
	return nil
}

func (h *workflowSpecHandler) Close() error {
	return nil
}

type Aborter interface {
	Abort(ctx context.Context, msg *api.Message) error
}

type CreateAborter interface {
	Create(ctx context.Context, msg *api.Message) error
	Aborter
}

type DeleteAborter interface {
	Delete(ctx context.Context, msg *api.Message) error
	Aborter
}

type TwoPhaseCommit[T any] interface {
//	kv[T]
	Prepare(ctx context.Context, msg *api.Message) error
	Commit(ctx context.Context, msg *api.Message) error
	Abort(ctx context.Context, msg *api.Message) error
}

type kv[T any] interface {
	Get(key string) (T, error)
	Put(key string, value T) error
	Delete(key string) error
}

type job struct {
	Name string
	TOML string
}

// almost behaves like a 2PC store, but doesn't doesn't have durable storage for write ahead log
type almostTwoPhaseCommitStore struct {
	mu    sync.Mutex
	store kv[*job]
	wal   kv[*job]
}

func NewAlmostTwoPhaseCommitStore() TwoPhaseCommit {
	return &almostTwoPhaseCommitStore]{
		store: kv[*job],
		wal:   kv[*job],
	}
}

func (s *almostTwoPhaseCommitStore) Prepare(ctx context.Context, msg *api.Message) error {
	j, err := extractJob(msg)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wal.Put(msg.Body.MessageId, j)
}

func (s *almostTwoPhaseCommitStore[T]) Commit(ctx context.Context, msg *api.Message) error {
	j, err := extractJob(msg)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	j, err = s.wal.Get(j.Name)
	if err != nil {
		return err
	}
	s.store[msg.Body.MessageId] = &job{
		Name: msg.Body.Method,
		TOML: string(msg.Body.Payload),
	}
	delete(s.wal, msg.Body.MessageId)
	return nil
}

func (s *almostTwoPhaseCommitStore) Abort(ctx context.Context, msg *api.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.wal, msg.Body.MessageId)
	return nil
}

func extractJob(msg *api.Message) (*job.Job, error) {
	var j job.Job
	if err := json.Unmarshal(msg.Body.Payload, &j); err != nil {
		return nil, err
	}
	return &j, nil
}
