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

	defaultRequestTimeoutSec  = 30
	defaultNodeSendTimeoutSec = 10

	// defaultQuorumGraceMillis bounds the extra wait after quorum is reached. It must
	// stay well below the caller's own HTTP deadline, which is what actually cuts the
	// request short when the DON never produces 2F+1 signed responses.
	defaultQuorumGraceMillis = 10_000

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
	completed atomic.Bool

	// graceStarted is set the first time the request holds F+1 signed responses, so
	// the grace timer is armed once per request rather than on every later response.
	graceStarted atomic.Bool
	graceTimer   clockwork.Timer

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

func (ar *activeRequest) setGraceTimer(t clockwork.Timer) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.graceTimer = t
}

func (ar *activeRequest) stopGraceTimer() {
	ar.mu.Lock()
	t := ar.graceTimer
	ar.graceTimer = nil
	ar.mu.Unlock()
	if t != nil {
		t.Stop()
	}
}

type relayBundler interface {
	Bundle(req jsonrpc.Request[json.RawMessage], resps map[string]jsonrpc.Response[json.RawMessage], l logger.Logger) (*BundleSummary, error)
}

type Config struct {
	RequestTimeoutSec int `json:"requestTimeoutSec"`

	// NodeSendTimeoutSec bounds each individual fan-out send to a single relay node, and is
	// clamped to RequestTimeoutSec. It must stay below it so that one node whose connection
	// accepts no writes cannot delay delivery to the rest of the DON.
	NodeSendTimeoutSec int `json:"nodeSendTimeoutSec"`

	// QuorumGraceMillis is how long the handler keeps collecting responses after the
	// first F+1 signed responses arrive, before forwarding whatever it has. It bounds
	// the wait for a DON that answers with quorum but never reaches 2F+1 signed
	// responses, which would otherwise hold the request until RequestTimeoutSec and
	// forward a long-viable bundle after the caller's own deadline has passed.
	// Clamped to RequestTimeoutSec; a negative value disables the grace window and
	// restores waiting until expiry.
	QuorumGraceMillis int `json:"quorumGraceMillis"`
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
	nodeSendTimeout       time.Duration
	quorumGrace           time.Duration

	activeRequests map[string]*activeRequest
	metrics        *metrics

	bundler relayBundler

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
		cfg.RequestTimeoutSec = defaultRequestTimeoutSec
	}

	if cfg.NodeSendTimeoutSec == 0 {
		cfg.NodeSendTimeoutSec = defaultNodeSendTimeoutSec
	}
	if cfg.NodeSendTimeoutSec > cfg.RequestTimeoutSec {
		cfg.NodeSendTimeoutSec = cfg.RequestTimeoutSec
	}

	switch {
	case cfg.QuorumGraceMillis == 0:
		cfg.QuorumGraceMillis = defaultQuorumGraceMillis
	case cfg.QuorumGraceMillis < 0:
		cfg.QuorumGraceMillis = 0
	}
	if maxGraceMillis := cfg.RequestTimeoutSec * 1000; cfg.QuorumGraceMillis > maxGraceMillis {
		cfg.QuorumGraceMillis = maxGraceMillis
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
		nodeSendTimeout:       time.Duration(cfg.NodeSendTimeoutSec) * time.Second,
		quorumGrace:           time.Duration(cfg.QuorumGraceMillis) * time.Millisecond,
		globalNodeRateLimiter: globalNodeRateLimiter,
		perNodeRateLimiters:   perNodeRateLimiters,
		activeRequests:        make(map[string]*activeRequest),
		mu:                    sync.RWMutex{},
		stopCh:                make(services.StopChan),
		metrics:               metrics,
		bundler:               &bundler{},
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
		l := logger.With(h.lggr, "method", er.req.Method, "requestID", er.req.ID)
		l.Debugw("request expired, evaluating collected relay responses",
			"collected", len(responses),
			"nodes", len(h.donConfig.Members),
			"unanswered", len(h.donConfig.Members)-len(responses),
		)
		summary, err := h.bundler.Bundle(er.req, responses, l)
		if err != nil {
			l.Errorw("failed to build relay response bundle", "error", err)
			if sendErr := h.sendResponseAndClearRequest(ctx, er, h.constructErrorResponse(er.req, api.FatalError, err)); sendErr != nil {
				l.Errorw("error returning bundle failure on expiry", "error", sendErr)
			}
			continue
		}
		// Expiry makes further responses unavailable to this request. The common
		// readiness path forwards a viable partial bundle or returns a timeout.
		if err := h.forwardBundleOrTerminateIfReady(ctx, l, er, summary, 0, true); err != nil {
			l.Errorw("error forwarding bundle on expiry", "error", err)
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
	// Bundle readiness counts decodable signed responses, while the enclave remains
	// responsible for signature and quorum verification. A collected response cannot
	// be replaced by a later response from the same node.
	summary, err := h.bundler.Bundle(ar.req, copiedResponses, l)
	if err != nil {
		l.Errorw("failed to build relay response bundle", "error", err)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err))
	}
	return h.forwardBundleOrTerminateIfReady(ctx, l, ar, summary, len(h.donConfig.Members)-summary.Total(), false)
}

// forwardBundleOrTerminateIfReady applies the same completion rules to live and
// expired requests. Remaining counts responses that can still arrive; expiry
// passes zero because waiting can no longer improve the bundle.
func (h *handler) forwardBundleOrTerminateIfReady(ctx context.Context, l logger.Logger, ar *activeRequest, summary *BundleSummary, remaining int, expired bool) error {
	earlyNeed := 2*h.donConfig.F + 1
	minQuorum := h.donConfig.F + 1
	nodes := len(h.donConfig.Members)
	unanswered := nodes - summary.Total()
	maxPossibleSigned := summary.Signed() + remaining
	if summary.Signed() >= earlyNeed {
		return h.forwardBundle(ctx, l, ar, summary)
	}
	if remaining == 0 && summary.Signed() >= minQuorum {
		l.Infow("no more responses expected; forwarding partial signed bundle",
			"signed", summary.Signed(),
			"minQuorum", minQuorum,
			"earlyNeed", earlyNeed,
			"collected", summary.Total(),
			"nodes", nodes,
			"unanswered", unanswered,
			"expired", expired,
			"errors", summary.Error(),
			"undecodable", summary.Undecodable(),
			"nodeErrors", summary.NodeErrorsFormatted(),
		)
		return h.forwardBundle(ctx, l, ar, summary)
	}
	if maxPossibleSigned < minQuorum {
		if expired {
			l.Warnw("request expired before relay quorum was reached",
				"signed", summary.Signed(),
				"minQuorum", minQuorum,
				"earlyNeed", earlyNeed,
				"collected", summary.Total(),
				"nodes", nodes,
				"unanswered", unanswered,
				"errors", summary.Error(),
				"undecodable", summary.Undecodable(),
				"nodeErrors", summary.NodeErrorsFormatted(),
			)
			return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.RequestTimeoutError,
				fmt.Errorf("request expired: %d signed relay responses, need at least %d to reach quorum (collected=%d nodes=%d unanswered=%d errors=%d undecodable=%d)",
					summary.Signed(), minQuorum, summary.Total(), nodes, unanswered, summary.Error(), summary.Undecodable())))
		}
		l.Warnw("relay quorum unreachable; returning error without waiting",
			"signed", summary.Signed(),
			"minQuorum", minQuorum,
			"earlyNeed", earlyNeed,
			"collected", summary.Total(),
			"nodes", nodes,
			"remaining", remaining,
			"maxPossibleSigned", maxPossibleSigned,
			"errors", summary.Error(),
			"undecodable", summary.Undecodable(),
			"nodeErrors", summary.NodeErrorsFormatted(),
		)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError,
			fmt.Errorf("relay quorum unreachable: %d signed responses, at most %d possible, need %d (collected=%d nodes=%d remaining=%d errors=%d undecodable=%d)",
				summary.Signed(), maxPossibleSigned, minQuorum, summary.Total(), nodes, remaining, summary.Error(), summary.Undecodable())))
	}
	if !expired && summary.Signed() >= minQuorum {
		h.startQuorumGrace(l, ar, summary)
	}
	l.Debugw("waiting for more signed relay responses before forwarding bundle",
		"signed", summary.Signed(),
		"earlyNeed", earlyNeed,
		"minQuorum", minQuorum,
		"collected", summary.Total(),
		"nodes", nodes,
		"remaining", remaining,
		"maxPossibleSigned", maxPossibleSigned,
		"errors", summary.Error(),
		"undecodable", summary.Undecodable(),
		"nodeErrors", summary.NodeErrorsFormatted(),
	)
	return nil
}

// startQuorumGrace arms the grace window the first time a request holds minQuorum
// signed responses. It bounds how long a request that will never reach earlyNeed
// keeps waiting for the rest of the DON: without it such a request is held until
// requestTimeout, long after the caller's own deadline has elapsed, so a bundle
// that was viable within milliseconds is forwarded to nobody.
func (h *handler) startQuorumGrace(l logger.Logger, ar *activeRequest, summary *BundleSummary) {
	if h.quorumGrace <= 0 {
		return
	}
	if !ar.graceStarted.CompareAndSwap(false, true) {
		return
	}
	l.Infow("relay quorum reached below earlyNeed; starting grace window",
		"signed", summary.Signed(),
		"minQuorum", h.donConfig.F+1,
		"earlyNeed", 2*h.donConfig.F+1,
		"collected", summary.Total(),
		"nodes", len(h.donConfig.Members),
		"grace", h.quorumGrace,
	)
	ar.setGraceTimer(h.clock.AfterFunc(h.quorumGrace, func() { h.onQuorumGraceElapsed(ar) }))
}

// onQuorumGraceElapsed forwards the bundle collected during the grace window.
// Collected responses are never replaced, so the signed count only grows: a request
// that armed the timer still holds at least minQuorum signed responses here. If the
// earlyNeed path already answered, the send is a no-op.
func (h *handler) onQuorumGraceElapsed(ar *activeRequest) {
	ctx, cancel := h.stopCh.NewCtx()
	defer cancel()

	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)
	summary, err := h.bundler.Bundle(ar.req, ar.copiedResponses(), l)
	if err != nil {
		l.Errorw("failed to build relay response bundle after quorum grace", "error", err)
		if sendErr := h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err)); sendErr != nil {
			l.Errorw("error returning bundle failure after quorum grace", "error", sendErr)
		}
		return
	}

	minQuorum := h.donConfig.F + 1
	nodes := len(h.donConfig.Members)
	if summary.Signed() < minQuorum {
		l.Warnw("quorum grace elapsed below quorum; leaving request to expiry",
			"signed", summary.Signed(),
			"minQuorum", minQuorum,
			"collected", summary.Total(),
			"nodes", nodes,
		)
		return
	}

	l.Infow("quorum grace elapsed; forwarding partial signed bundle",
		"signed", summary.Signed(),
		"minQuorum", minQuorum,
		"earlyNeed", 2*h.donConfig.F+1,
		"collected", summary.Total(),
		"nodes", nodes,
		"unanswered", nodes-summary.Total(),
		"grace", h.quorumGrace,
		"errors", summary.Error(),
		"undecodable", summary.Undecodable(),
		"nodeErrors", summary.NodeErrorsFormatted(),
	)
	if err := h.forwardBundle(ctx, l, ar, summary); err != nil {
		l.Errorw("error forwarding bundle after quorum grace", "error", err)
	}
}

// forwardBundle sends a previously-built bundle to the enclave. The gateway makes
// no trust decision; the enclave verifies signatures and reaches quorum.
func (h *handler) forwardBundle(ctx context.Context, l logger.Logger, ar *activeRequest, summary *BundleSummary) error {
	nodes := len(h.donConfig.Members)
	rawResponse, err := jsonrpc.EncodeResponse(summary.Response())
	if err != nil {
		l.Errorw("failed to encode response", "error", err)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.NodeReponseEncodingError, err))
	}
	l.Infow("forwarding relay response bundle to enclave",
		"signedResponses", summary.Signed(),
		"earlyNeed", 2*h.donConfig.F+1,
		"collected", summary.Total(),
		"nodes", nodes,
		"remaining", nodes-summary.Total(),
		"errorResponses", summary.Error(),
		"undecodableResponses", summary.Undecodable(),
		"nodeErrors", summary.NodeErrorsFormatted(),
	)
	return h.sendResponseAndClearRequest(ctx, ar, gwhandlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   api.NoError,
	})
}

func (h *handler) fanOutToNodes(ctx context.Context, l logger.Logger, ar *activeRequest) error {
	var (
		group      errgroup.Group
		nodeErrors atomic.Uint32
	)

	// Each send is bounded independently. A node whose websocket accepts no writes blocks
	// until its context is cancelled, and because the caller only reads the response callback
	// after this function returns, an unbounded send would hold the request open until the
	// client gives up, discarding a bundle that already reached quorum.
	sendCtx, cancel := context.WithTimeout(ctx, h.nodeSendTimeout)
	defer cancel()

	for _, node := range h.donConfig.Members {
		group.Go(func() error {
			err := h.don.SendToNode(sendCtx, node.Address, &ar.req)
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
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to relay nodes")
	return nil
}

// sendResponseAndClearRequest claims the request, sends payload, and removes it from
// activeRequests. Concurrent completion paths (node-message forward,
// terminal-state forward, quorum grace, expiry) may all race here; only the first
// claimer sends. Metrics are recorded only after a successful send so losers do not
// double-count.
func (h *handler) sendResponseAndClearRequest(ctx context.Context, ar *activeRequest, payload gwhandlers.UserCallbackPayload) error {
	if !ar.completed.CompareAndSwap(false, true) {
		// Another path already answered this request.
		return nil
	}
	ar.stopGraceTimer()

	sendErr := ar.SendResponse(payload)

	h.mu.Lock()
	delete(h.activeRequests, ar.req.ID)
	h.mu.Unlock()

	if sendErr != nil {
		h.lggr.Errorw("error sending response to user", "requestID", ar.req.ID, "error", sendErr)
		return sendErr
	}

	h.recordMetrics(ctx, payload.ErrorCode)
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
