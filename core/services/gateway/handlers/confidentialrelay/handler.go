package confidentialrelay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gwhandlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

const (
	defaultCleanUpPeriod = 5 * time.Second

	// Re-exported from chainlink-common for local use and test convenience.
	MethodSecretsGet     = relaytypes.MethodSecretsGet
	MethodCapabilityExec = relaytypes.MethodCapabilityExec
)

var _ gwhandlers.Handler = (*handler)(nil)

type metrics struct {
	requestInternalError metric.Int64Counter
	requestUserError     metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("confidential_relay_gateway_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}

	requestUserError, err := beholder.GetMeter().Int64Counter("confidential_relay_gateway_request_user_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register user error counter: %w", err)
	}

	requestSuccess, err := beholder.GetMeter().Int64Counter("confidential_relay_gateway_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}

	return &metrics{
		requestInternalError: requestInternalError,
		requestUserError:     requestUserError,
		requestSuccess:       requestSuccess,
	}, nil
}

type activeRequest struct {
	req       jsonrpc.Request[json.RawMessage]
	responses map[string]*jsonrpc.Response[json.RawMessage]
	mu        sync.Mutex

	createdAt time.Time
	gwhandlers.Callback
}

func (ar *activeRequest) addResponseForNode(nodeAddr string, resp *jsonrpc.Response[json.RawMessage]) bool {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	_, exists := ar.responses[nodeAddr]
	if exists {
		return false
	}

	ar.responses[nodeAddr] = resp
	return true
}

func (ar *activeRequest) copiedResponses() map[string]jsonrpc.Response[json.RawMessage] {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	copied := make(map[string]jsonrpc.Response[json.RawMessage], len(ar.responses))
	for k, response := range ar.responses {
		var copiedResponse jsonrpc.Response[json.RawMessage]
		if response != nil {
			copiedResponse = *response
			if response.Result != nil {
				copiedResult := *response.Result
				copiedResponse.Result = &copiedResult
			}
			if response.Error != nil {
				copiedError := *response.Error
				copiedResponse.Error = &copiedError
			}
		}
		copied[k] = copiedResponse
	}
	return copied
}

type relayAggregator interface {
	Aggregate(resps map[string]jsonrpc.Response[json.RawMessage], donF int, donMembersCount int, l logger.Logger) (*jsonrpc.Response[json.RawMessage], error)
}

type Config struct {
	RequestTimeoutSec int `json:"requestTimeoutSec"`
}

type handler struct {
	services.StateMachine
	donConfig *config.DONConfig
	don       gwhandlers.DON
	codec     api.JsonRPCCodec
	lggr      logger.Logger
	mu        sync.RWMutex
	stopCh    services.StopChan

	globalNodeRateLimiter limits.RateLimiter
	perNodeRateLimiters   map[string]limits.RateLimiter
	requestTimeout        time.Duration

	activeRequests map[string]*activeRequest
	metrics        *metrics

	aggregator relayAggregator

	clock clockwork.Clock
}

func (h *handler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *handler) Name() string {
	return h.lggr.Name()
}

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gwhandlers.DON, lggr logger.Logger, clock clockwork.Clock, limitsFactory limits.Factory) (*handler, error) {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal method config: %w", err)
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	globalNodeRateLimiter, err := limitsFactory.MakeRateLimiter(cresettings.Default.GatewayConfidentialRelayGlobalRate)
	if err != nil {
		return nil, fmt.Errorf("failed to create global node rate limiter: %w", err)
	}

	perNodeRateLimiters := make(map[string]limits.RateLimiter, len(donConfig.Members))
	for _, member := range donConfig.Members {
		rl, makeErr := limitsFactory.MakeRateLimiter(cresettings.Default.GatewayConfidentialRelayPerNodeRate)
		if makeErr != nil {
			return nil, fmt.Errorf("failed to create per-node rate limiter for %s: %w", member.Address, makeErr)
		}
		perNodeRateLimiters[member.Address] = rl
	}

	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &handler{
		donConfig:             donConfig,
		don:                   don,
		lggr:                  logger.Named(lggr, "ConfidentialRelayHandler:"+donConfig.DonId),
		requestTimeout:        time.Duration(cfg.RequestTimeoutSec) * time.Second,
		globalNodeRateLimiter: globalNodeRateLimiter,
		perNodeRateLimiters:   perNodeRateLimiters,
		activeRequests:        make(map[string]*activeRequest),
		mu:                    sync.RWMutex{},
		stopCh:                make(services.StopChan),
		metrics:               metrics,
		aggregator:            &aggregator{},
		clock:                 clock,
	}, nil
}

func (h *handler) Start(_ context.Context) error {
	return h.StartOnce("ConfidentialRelayHandler", func() error {
		h.lggr.Info("starting confidential relay handler")
		go func() {
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
			ticker := h.clock.NewTicker(defaultCleanUpPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.Chan():
					h.removeExpiredRequests(ctx)
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce("ConfidentialRelayHandler", func() error {
		h.lggr.Info("closing confidential relay handler")
		close(h.stopCh)
		var err error
		if h.globalNodeRateLimiter != nil {
			err = errors.Join(err, h.globalNodeRateLimiter.Close())
		}
		for _, rl := range h.perNodeRateLimiters {
			err = errors.Join(err, rl.Close())
		}
		return err
	})
}

func (h *handler) removeExpiredRequests(ctx context.Context) {
	h.mu.RLock()
	var expiredRequests []*activeRequest
	now := h.clock.Now()
	for _, userRequest := range h.activeRequests {
		if now.Sub(userRequest.createdAt) > h.requestTimeout {
			expiredRequests = append(expiredRequests, userRequest)
		}
	}
	h.mu.RUnlock()

	for _, er := range expiredRequests {
		responses := er.copiedResponses()
		h.lggr.Debugw("request expired without quorum", "requestID", er.req.ID, "responseCount", len(responses), "required", h.donConfig.F+1)
		err := h.sendResponseAndCleanup(ctx, er, h.constructErrorResponse(er.req, api.RequestTimeoutError, fmt.Errorf("request expired: got %d/%d responses", len(responses), h.donConfig.F+1)))
		if err != nil {
			h.lggr.Errorw("error sending response to user", "requestID", er.req.ID, "error", err)
		}
	}
}

func (h *handler) Methods() []string {
	return []string{MethodSecretsGet, MethodCapabilityExec}
}

func (h *handler) HandleLegacyUserMessage(_ context.Context, _ *api.Message, _ gwhandlers.Callback) error {
	return errors.New("confidential relay handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callback gwhandlers.Callback) error {
	if req.ID == "" {
		return errors.New("request ID cannot be empty")
	}
	if len(req.ID) > 200 {
		return errors.New("request ID is too long: " + strconv.Itoa(len(req.ID)) + ". max is 200 characters")
	}

	l := logger.With(h.lggr, "method", req.Method, "requestID", req.ID)
	l.Debugw("handling confidential relay request")

	ar, err := h.newActiveRequest(req, callback)
	if err != nil {
		return err
	}

	return h.fanOutToNodes(ctx, l, ar)
}

func (h *handler) newActiveRequest(req jsonrpc.Request[json.RawMessage], callback gwhandlers.Callback) (*activeRequest, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.activeRequests[req.ID] != nil {
		h.lggr.Errorw("request id already exists", "requestID", req.ID)
		return nil, errors.New("request ID already exists: " + req.ID)
	}
	ar := &activeRequest{
		Callback:  callback,
		req:       req,
		createdAt: h.clock.Now(),
		responses: map[string]*jsonrpc.Response[json.RawMessage]{},
	}
	h.activeRequests[req.ID] = ar
	return ar, nil
}

func (h *handler) getActiveRequest(requestID string) *activeRequest {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.activeRequests[requestID]
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	l := logger.With(h.lggr, "method", resp.Method, "requestID", resp.ID, "nodeAddr", nodeAddr)
	l.Debugw("handling node response")

	nodeRateLimiter, ok := h.perNodeRateLimiters[nodeAddr]
	if !ok {
		return fmt.Errorf("received message from unexpected node %s", nodeAddr)
	}
	if !nodeRateLimiter.Allow(ctx) {
		l.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}
	if !h.globalNodeRateLimiter.Allow(ctx) {
		l.Debug("global relay rate limit exceeded")
		return nil
	}

	ar := h.getActiveRequest(resp.ID)
	if ar == nil {
		l.Debugw("no pending request found for ID")
		return nil
	}

	added := ar.addResponseForNode(nodeAddr, resp)
	if !added {
		l.Errorw("duplicate response from node, ignoring", "nodeAddr", nodeAddr)
		return nil
	}

	copiedResponses := ar.copiedResponses()
	aggregatedResp, err := h.aggregator.Aggregate(copiedResponses, h.donConfig.F, len(h.donConfig.Members), l)
	switch {
	case errors.Is(err, errInsufficientResponsesForQuorum):
		l.Debugw("aggregating responses, waiting for other nodes...", "error", err)
		return nil
	case errors.Is(err, errQuorumUnobtainable):
		l.Errorw("quorum unobtainable, returning error to user", "error", err)
		return h.sendResponseAndCleanup(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err))
	case err != nil:
		l.Errorw("unexpected aggregation error", "error", err)
		return h.sendResponseAndCleanup(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err))
	}

	rawResponse, err := jsonrpc.EncodeResponse(aggregatedResp)
	if err != nil {
		h.lggr.Errorw("failed to encode response", "requestID", ar.req.ID, "error", err)
		return h.sendResponseAndCleanup(ctx, ar, h.constructErrorResponse(ar.req, api.NodeReponseEncodingError, err))
	}
	return h.sendResponseAndCleanup(ctx, ar, gwhandlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.NoError,
	})
}

func (h *handler) fanOutToNodes(ctx context.Context, l logger.Logger, ar *activeRequest) error {
	var (
		group      errgroup.Group
		nodeErrors atomic.Uint32
	)

	for _, node := range h.donConfig.Members {
		group.Go(func() error {
			err := h.don.SendToNode(ctx, node.Address, &ar.req)
			if err != nil {
				nodeErrors.Add(1)
				l.Errorw("error sending request to node", "node", node.Address, "error", err)
			}
			return nil
		})
	}

	_ = group.Wait()

	numNodeErrors := nodeErrors.Load()
	remainingPossibleResponses := len(h.donConfig.Members) - int(numNodeErrors)
	if remainingPossibleResponses < h.donConfig.F+1 && numNodeErrors > 0 {
		return h.sendResponseAndCleanup(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to relay nodes")
	return nil
}

// sendResponseAndCleanup sends payload.
// The request is always removed from activeRequests
// regardless of whether the send succeeds, since a failed callback cannot
// be retried.
func (h *handler) sendResponseAndCleanup(ctx context.Context, ar *activeRequest, payload gwhandlers.UserCallbackPayload) error {
	h.recordMetrics(ctx, payload.ErrorCode)
	sendErr := ar.SendResponse(payload)

	h.mu.Lock()
	delete(h.activeRequests, ar.req.ID)
	h.mu.Unlock()

	if sendErr != nil {
		h.lggr.Errorw("error sending response to user", "requestID", ar.req.ID, "error", sendErr)
		return sendErr
	}

	h.lggr.Debugw("response sent to user", "requestID", ar.req.ID, "errorCode", payload.ErrorCode)
	return nil
}

func (h *handler) recordMetrics(ctx context.Context, errorCode api.ErrorCode) {
	//nolint:exhaustive // do not record other errors
	switch errorCode {
	case api.HandlerError:
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", errorCode.String()),
		))
	case api.UnsupportedDONIdError:
		h.metrics.requestUserError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	case api.NoError:
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	}
}

func (h *handler) constructErrorResponse(req jsonrpc.Request[json.RawMessage], errorCode api.ErrorCode, err error) gwhandlers.UserCallbackPayload {
	//nolint:exhaustive // do not modify other error codes
	switch errorCode {
	case api.NodeReponseEncodingError:
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		err = fmt.Errorf("invalid params error: %w", err)
	case api.UnsupportedMethodError:
		err = fmt.Errorf("unsupported method(%s): %w", req.Method, err)
	case api.UserMessageParseError:
		err = fmt.Errorf("user message parse error: %w", err)
	}
	return gwhandlers.UserCallbackPayload{
		RawResponse: h.codec.EncodeNewErrorResponse(
			req.ID,
			api.ToJSONRPCErrorCode(errorCode),
			err.Error(),
			nil,
		),
		ErrorCode: errorCode,
	}
}
