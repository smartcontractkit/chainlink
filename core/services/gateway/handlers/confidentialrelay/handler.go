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
	gwhandlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

const (
	// defaultCleanUpPeriod is how often expired requests are swept and closed grace
	// windows are forwarded, so it also bounds how far past its deadline a grace
	// window can run.
	defaultCleanUpPeriod = time.Second

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

// requestLabels are the identifiers the gateway pulls out of a relay request's
// params purely for logging, so a gateway line can be correlated with the
// relay-DON's and the enclave's lines for the same workflow execution. The
// gateway stays a dumb relay: it does not otherwise interpret params.
type requestLabels struct {
	WorkflowID  string `json:"workflow_id"`
	ExecutionID string `json:"execution_id"`
}

// extractRequestLabels best-effort decodes the logging identifiers from a
// request's params. Both relay methods' params carry these fields. A decode
// failure leaves them empty and is only logged: these labels are for
// correlation, and the params themselves are validated by the relay nodes,
// not here, so a request whose params do not decode is still fanned out and
// rejected there. ProcessRequest has already parsed the envelope as valid
// JSON by this point, so a failure here means params is not an object or
// carries non-string identifiers — malformed input rather than a gateway bug,
// hence debug level to avoid handing a caller a log-spam lever.
func (h *handler) extractRequestLabels(req jsonrpc.Request[json.RawMessage]) requestLabels {
	var labels requestLabels
	if req.Params == nil {
		return labels
	}
	if err := json.Unmarshal(*req.Params, &labels); err != nil {
		h.lggr.Debugw("could not decode relay request params for logging labels",
			"method", req.Method, "requestID", req.ID, "err", err)
	}
	return labels
}

type activeRequest struct {
	req       jsonrpc.Request[json.RawMessage]
	labels    requestLabels
	responses map[string]*jsonrpc.Response[json.RawMessage]
	mu        sync.Mutex
	completed atomic.Bool

	// graceStarted is set the first time the request holds F+1 signed responses, so
	// the grace deadline is armed once per request rather than moved forward by every
	// later response. graceDeadline is guarded by mu and is only meaningful once
	// graceStarted is set.
	graceStarted  atomic.Bool
	graceDeadline time.Time

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

func (ar *activeRequest) armGraceDeadline(deadline time.Time) bool {
	if !ar.graceStarted.CompareAndSwap(false, true) {
		return false
	}
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.graceDeadline = deadline
	return true
}

// graceElapsed reports whether the request reached quorum and its grace window has
// since closed, meaning the collected bundle should be forwarded now.
func (ar *activeRequest) graceElapsed(now time.Time) bool {
	if !ar.graceStarted.Load() {
		return false
	}
	ar.mu.Lock()
	defer ar.mu.Unlock()
	return !ar.graceDeadline.IsZero() && !now.Before(ar.graceDeadline)
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
	// shards is the DON×shard matrix this handler relays for. Quorum is evaluated
	// independently per shard (each with its own F and members); the first shard to
	// reach quorum produces the user response.
	shards []*gwhandlers.ShardEndpoint
	// donID names the handler for logging/metrics; it is the first shard's DON name.
	donID  string
	codec  api.JSONRPCCodec
	lggr   logger.Logger
	mu     sync.RWMutex
	stopCh services.StopChan

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

func NewHandler(methodConfig json.RawMessage, dons *gwhandlers.ShardedDONs, lggr logger.Logger, clock clockwork.Clock, limitsFactory limits.Factory) (*handler, error) {
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

	shards, _, err := dons.BuildShardEndpoints()
	if err != nil {
		return nil, err
	}
	if len(shards) == 0 {
		return nil, errors.New("confidential relay handler requires at least one shard")
	}
	donID := shards[0].DonName

	globalNodeRateLimiter, err := limitsFactory.MakeRateLimiter(cresettings.Default.GatewayConfidentialRelayGlobalRate)
	if err != nil {
		return nil, fmt.Errorf("failed to create global node rate limiter: %w", err)
	}

	members := gwhandlers.AllMembers(shards)
	perNodeRateLimiters := make(map[string]limits.RateLimiter, len(members))
	for _, member := range members {
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
		shards:                shards,
		donID:                 donID,
		lggr:                  logger.Named(lggr, "ConfidentialRelayHandler:"+donID),
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

// totalNodes is the number of relay nodes across every shard, used only for logging.
func (h *handler) totalNodes() int {
	n := 0
	for _, shard := range h.shards {
		n += len(shard.Members)
	}
	return n
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
					h.forwardGracedRequests(ctx)
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
		l := h.requestLogger(er.req, er.labels)
		collected := len(er.copiedResponses())
		l.Debugw("request expired, evaluating collected relay responses",
			"collected", collected,
			"nodes", h.totalNodes(),
			"unanswered", h.totalNodes()-collected,
		)
		// Expiry makes further responses unavailable to this request. The common
		// readiness path forwards a viable partial bundle from whichever shard
		// qualifies, or returns a timeout.
		if err := h.evaluateShards(ctx, l, er, true); err != nil {
			l.Errorw("error evaluating relay response bundles on expiry", "error", err)
		}
	}
}

func (h *handler) Methods() []string {
	return []string{MethodSecretsGet, MethodCapabilityExec}
}

func (h *handler) HandleLegacyUserMessage(_ context.Context, _ *api.Message, _ gwhandlers.Callback) error {
	return errors.New("confidential relay handler does not support legacy messages")
}

// requestLogger returns the logger for one relay request: the gateway request
// id plus the workflow/execution identity it carries, so a line can be
// correlated with the relay DON's and the enclave's logs for the same
// execution. The request id changes per enclave retry; the execution identity
// does not.
func (h *handler) requestLogger(req jsonrpc.Request[json.RawMessage], labels requestLabels) logger.Logger {
	return logger.With(h.lggr,
		"method", req.Method,
		"requestID", req.ID,
		"workflowID", labels.WorkflowID,
		"executionID", labels.ExecutionID,
	)
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callback gwhandlers.Callback) error {
	if req.ID == "" {
		return errors.New("request ID cannot be empty")
	}
	if len(req.ID) > 200 {
		return errors.New("request ID is too long: " + strconv.Itoa(len(req.ID)) + ". max is 200 characters")
	}

	labels := h.extractRequestLabels(req)
	l := h.requestLogger(req, labels)
	l.Debugw("handling confidential relay request", "shards", len(h.shards), "nodes", h.totalNodes())

	ar, err := h.newActiveRequest(req, labels, callback)
	if err != nil {
		return err
	}

	return h.fanOutToNodes(ctx, l, ar)
}

func (h *handler) newActiveRequest(req jsonrpc.Request[json.RawMessage], labels requestLabels, callback gwhandlers.Callback) (*activeRequest, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.activeRequests[req.ID] != nil {
		h.lggr.Errorw("request id already exists", "requestID", req.ID, "executionID", labels.ExecutionID)
		return nil, errors.New("request ID already exists: " + req.ID)
	}
	ar := &activeRequest{
		Callback:  callback,
		req:       req,
		labels:    labels,
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
	// A node response carries no params, so the execution identity comes from
	// the request it answers.
	l = logger.With(l, "workflowID", ar.labels.WorkflowID, "executionID", ar.labels.ExecutionID)

	added := ar.addResponseForNode(nodeAddr, resp)
	if !added {
		l.Errorw("duplicate response from node, ignoring", "nodeAddr", nodeAddr)
		return nil
	}

	return h.evaluateShards(ctx, l, ar, false)
}

// filterResponses restricts a request's collected responses to one shard's members.
// Node addresses are disjoint across shards (enforced by BuildShardEndpoints), so this
// never pulls in a response from another shard.
func filterResponses(all map[string]jsonrpc.Response[json.RawMessage], shard *gwhandlers.ShardEndpoint) map[string]jsonrpc.Response[json.RawMessage] {
	out := make(map[string]jsonrpc.Response[json.RawMessage], len(shard.Members))
	for _, m := range shard.Members {
		if r, ok := all[m.Address]; ok {
			out[m.Address] = r
		}
	}
	return out
}

// evaluateShards applies the completion rules to live and expired requests, once per
// shard: each shard's quorum is evaluated against its own F and members, independently
// of every other shard, and the first shard whose bundle is ready wins the request.
// Only when every shard can no longer reach its own quorum does the request fail.
// Expiry makes further responses unavailable, so remaining is forced to zero for every
// shard rather than counted.
func (h *handler) evaluateShards(ctx context.Context, l logger.Logger, ar *activeRequest, expired bool) error {
	all := ar.copiedResponses()

	var (
		allDead        = true
		firstUserErr   *jsonrpc.WireError
		totalSigned    int
		totalCollected int
	)

	for _, shard := range h.shards {
		summary, err := h.bundler.Bundle(ar.req, filterResponses(all, shard), l)
		if err != nil {
			l.Errorw("failed to build relay response bundle", "shard", shard.DonID, "error", err)
			return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err))
		}

		minQuorum := shard.F + 1
		earlyNeed := 2*shard.F + 1
		nodes := len(shard.Members)
		remaining := nodes - summary.Total()
		if expired {
			remaining = 0
		}
		maxPossibleSigned := summary.Signed() + remaining

		totalSigned += summary.Signed()
		totalCollected += summary.Total()

		if summary.Signed() >= earlyNeed {
			return h.forwardBundle(ctx, l, ar, shard, summary)
		}
		if remaining == 0 && summary.Signed() >= minQuorum {
			l.Infow("no more responses expected in shard; forwarding partial signed bundle",
				"shard", shard.DonID,
				"signed", summary.Signed(),
				"minQuorum", minQuorum,
				"earlyNeed", earlyNeed,
				"collected", summary.Total(),
				"nodes", nodes,
				"unanswered", nodes-summary.Total(),
				"expired", expired,
				"errors", summary.Error(),
				"undecodable", summary.Undecodable(),
				"nodeErrors", summary.NodeErrorsFormatted(),
			)
			return h.forwardBundle(ctx, l, ar, shard, summary)
		}
		if maxPossibleSigned < minQuorum {
			// A user-level node error means the request is invalid for every node,
			// so report it rather than the generic quorum failure.
			if userErr := summary.UserError(); userErr != nil && firstUserErr == nil {
				firstUserErr = userErr
			}
			continue
		}

		allDead = false
		if !expired && summary.Signed() >= minQuorum {
			h.startQuorumGrace(l, ar, summary, shard)
		}
	}

	if allDead {
		return h.failUnreachableQuorum(ctx, l, ar, firstUserErr, expired, totalSigned, totalCollected)
	}

	l.Debugw("waiting for more signed relay responses before forwarding bundle",
		"signed", totalSigned,
		"collected", totalCollected,
		"nodes", h.totalNodes(),
		"shards", len(h.shards),
	)
	return nil
}

// failUnreachableQuorum is reached only once every shard's quorum has become
// unreachable: a single-shard "quorum unreachable" verdict no longer fails the
// request outright, since another shard may still complete it.
func (h *handler) failUnreachableQuorum(ctx context.Context, l logger.Logger, ar *activeRequest, userErr *jsonrpc.WireError, expired bool, signed, collected int) error {
	nodes := h.totalNodes()
	if userErr != nil {
		l.Warnw("relay quorum unreachable in every shard due to a user error; propagating node error",
			"signed", signed,
			"collected", collected,
			"nodes", nodes,
			"shards", len(h.shards),
			"expired", expired,
			"userErrorCode", userErr.Code,
		)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.InvalidParamsError,
			errors.New(sanitizeNodeErrorMessage(userErr.Message))))
	}
	if expired {
		l.Warnw("request expired before relay quorum was reached in any shard",
			"signed", signed,
			"collected", collected,
			"nodes", nodes,
			"shards", len(h.shards),
		)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.RequestTimeoutError,
			fmt.Errorf("request expired: %d signed relay responses across %d shard(s), none reached quorum (collected=%d nodes=%d)",
				signed, len(h.shards), collected, nodes)))
	}
	l.Warnw("relay quorum unreachable in every shard; returning error without waiting",
		"signed", signed,
		"collected", collected,
		"nodes", nodes,
		"shards", len(h.shards),
	)
	return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError,
		fmt.Errorf("relay quorum unreachable: %d signed responses across %d shard(s), none can reach quorum (collected=%d nodes=%d)",
			signed, len(h.shards), collected, nodes)))
}

// startQuorumGrace arms the grace window the first time a request holds minQuorum
// signed responses in some shard. It bounds how long a request that will never reach
// earlyNeed keeps waiting for the rest of the DON: without it such a request is held
// until requestTimeout, long after the caller's own deadline has elapsed, so a bundle
// that was viable within milliseconds is forwarded to nobody. The deadline is swept
// by the cleanup goroutine rather than a timer, so no completion path runs on an
// untracked goroutine. It is armed once per request (not per shard): the first shard
// to reach minQuorum starts the clock for whichever shard ends up winning.
func (h *handler) startQuorumGrace(l logger.Logger, ar *activeRequest, summary *BundleSummary, shard *gwhandlers.ShardEndpoint) {
	if h.quorumGrace <= 0 {
		return
	}
	if !ar.armGraceDeadline(h.clock.Now().Add(h.quorumGrace)) {
		return
	}
	l.Infow("relay quorum reached below earlyNeed; starting grace window",
		"shard", shard.DonID,
		"signed", summary.Signed(),
		"minQuorum", shard.F+1,
		"earlyNeed", 2*shard.F+1,
		"collected", summary.Total(),
		"nodes", len(shard.Members),
		"grace", h.quorumGrace,
	)
}

// forwardGracedRequests forwards the bundle for every request whose grace window has
// closed. Collected responses are never replaced, so the signed count only grows: a
// request that armed the window still holds at least minQuorum signed responses here.
func (h *handler) forwardGracedRequests(ctx context.Context) {
	h.mu.RLock()
	var graced []*activeRequest
	now := h.clock.Now()
	for _, ar := range h.activeRequests {
		if ar.graceElapsed(now) {
			graced = append(graced, ar)
		}
	}
	h.mu.RUnlock()

	for _, ar := range graced {
		h.forwardAfterGrace(ctx, ar)
	}
}

func (h *handler) forwardAfterGrace(ctx context.Context, ar *activeRequest) {
	l := h.requestLogger(ar.req, ar.labels)
	all := ar.copiedResponses()

	for _, shard := range h.shards {
		summary, err := h.bundler.Bundle(ar.req, filterResponses(all, shard), l)
		if err != nil {
			l.Errorw("failed to build relay response bundle after quorum grace", "shard", shard.DonID, "error", err)
			if sendErr := h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, err)); sendErr != nil {
				l.Errorw("error returning bundle failure after quorum grace", "error", sendErr)
			}
			return
		}

		minQuorum := shard.F + 1
		if summary.Signed() < minQuorum {
			continue
		}

		nodes := len(shard.Members)
		l.Infow("quorum grace elapsed; forwarding partial signed bundle",
			"shard", shard.DonID,
			"signed", summary.Signed(),
			"minQuorum", minQuorum,
			"earlyNeed", 2*shard.F+1,
			"collected", summary.Total(),
			"nodes", nodes,
			"unanswered", nodes-summary.Total(),
			"grace", h.quorumGrace,
			"errors", summary.Error(),
			"undecodable", summary.Undecodable(),
			"nodeErrors", summary.NodeErrorsFormatted(),
		)
		if err := h.forwardBundle(ctx, l, ar, shard, summary); err != nil {
			l.Errorw("error forwarding bundle after quorum grace", "error", err)
		}
		return
	}

	l.Warnw("quorum grace elapsed without any shard reaching quorum; leaving request to expiry")
}

// forwardBundle sends a previously-built shard bundle to the enclave. The gateway makes
// no trust decision; the enclave verifies signatures and reaches quorum.
func (h *handler) forwardBundle(ctx context.Context, l logger.Logger, ar *activeRequest, shard *gwhandlers.ShardEndpoint, summary *BundleSummary) error {
	nodes := len(shard.Members)
	rawResponse, err := jsonrpc.EncodeResponse(summary.Response())
	if err != nil {
		l.Errorw("failed to encode response", "error", err)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.NodeReponseEncodingError, err))
	}
	l.Infow("forwarding relay response bundle to enclave",
		"shard", shard.DonID,
		"signedResponses", summary.Signed(),
		"earlyNeed", 2*shard.F+1,
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

// fanOutToNodes sends the request to every member of every shard, each via that
// shard's own connection manager. A shard whose own quorum becomes unreachable here
// does not fail the request by itself: another shard may still complete it. The
// request only fails if every shard is left unreachable.
func (h *handler) fanOutToNodes(ctx context.Context, l logger.Logger, ar *activeRequest) error {
	var group errgroup.Group
	shardNodeErrors := make([]atomic.Uint32, len(h.shards))

	// Each send is bounded independently. A node whose websocket accepts no writes blocks
	// until its context is cancelled, and because the caller only reads the response callback
	// after this function returns, an unbounded send would hold the request open until the
	// client gives up, discarding a bundle that already reached quorum.
	sendCtx, cancel := context.WithTimeout(ctx, h.nodeSendTimeout)
	defer cancel()

	for i, shard := range h.shards {
		for _, node := range shard.Members {
			group.Go(func() error {
				err := shard.ConnMgr.SendToNode(sendCtx, node.Address, &ar.req)
				if err != nil {
					shardNodeErrors[i].Add(1)
					l.Errorw("error sending request to node", "node", node.Address, "shard", shard.DonID, "error", err)
				}
				return nil
			})
		}
	}

	_ = group.Wait()

	var anyShardReachable bool
	var totalNodeErrors uint32
	for i, shard := range h.shards {
		numNodeErrors := shardNodeErrors[i].Load()
		totalNodeErrors += numNodeErrors
		reachable := len(shard.Members) - int(numNodeErrors)
		if reachable >= shard.F+1 {
			anyShardReachable = true
		} else if numNodeErrors > 0 {
			l.Errorw("too few relay nodes reachable in shard to reach quorum",
				"shard", shard.DonID,
				"nodeErrors", numNodeErrors,
				"reachable", reachable,
				"minQuorum", shard.F+1,
				"nodes", len(shard.Members),
			)
		}
	}
	if !anyShardReachable {
		// The per-shard failures are logged above; this is the aggregate decision
		// that ends the request, so it gets its own line.
		l.Errorw("too few relay nodes reachable in any shard to reach quorum; failing request", "nodeErrors", totalNodeErrors)
		return h.sendResponseAndClearRequest(ctx, ar, h.constructErrorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to relay nodes", "nodeErrors", totalNodeErrors)
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

	sendErr := ar.SendResponse(payload)

	h.mu.Lock()
	delete(h.activeRequests, ar.req.ID)
	h.mu.Unlock()

	if sendErr != nil {
		h.lggr.Errorw("error sending response to user", "requestID", ar.req.ID, "executionID", ar.labels.ExecutionID, "error", sendErr)
		return sendErr
	}

	h.recordMetrics(ctx, payload.ErrorCode)
	h.lggr.Debugw("response sent to user", "requestID", ar.req.ID, "executionID", ar.labels.ExecutionID, "errorCode", payload.ErrorCode)
	return nil
}

func (h *handler) recordMetrics(ctx context.Context, errorCode api.ErrorCode) {
	//nolint:exhaustive // do not record other errors
	switch errorCode {
	case api.HandlerError:
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donID),
			attribute.String("error", errorCode.String()),
		))
	case api.UnsupportedDONIdError:
		h.metrics.requestUserError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donID),
		))
	case api.NoError:
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donID),
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
