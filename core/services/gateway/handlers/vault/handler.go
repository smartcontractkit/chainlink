package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gwhandlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlerscommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	defaultCleanUpPeriod                    = 5 * time.Second
	defaultPublicKeyGetCacheDurationSeconds = 300
)

var (
	_                                 gwhandlers.Handler = (*handler)(nil)
	errInsufficientResponsesForQuorum                    = errors.New("insufficient valid responses to reach quorum")
	errQuorumUnobtainable                                = errors.New("quorum unobtainable")
)

type metrics struct {
	requestInternalError metric.Int64Counter
	requestUserError     metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}

	requestUserError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_user_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register user error counter: %w", err)
	}

	requestSuccess, err := beholder.GetMeter().Int64Counter("gateway_vault_request_success")
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

func (ar *activeRequest) copiedResponses() map[string]*jsonrpc.Response[json.RawMessage] {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	copied := make(map[string]*jsonrpc.Response[json.RawMessage], len(ar.responses))
	maps.Copy(copied, ar.responses)
	return copied
}

type capabilitiesRegistry interface {
	DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error)
}

type aggregator interface {
	Aggregate(ctx context.Context, l logger.Logger, resps map[string]*jsonrpc.Response[json.RawMessage], currResp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error)
}

type handler struct {
	services.StateMachine
	methodConfig      Config
	donConfig         *config.DONConfig
	don               gwhandlers.DON
	lggr              logger.Logger
	codec             api.JsonRPCCodec
	mu                sync.RWMutex
	stopCh            services.StopChan
	requestAuthorizer vaultcap.RequestAuthorizer

	nodeRateLimiter *ratelimit.RateLimiter
	requestTimeout  time.Duration

	activeRequests map[string]*activeRequest
	metrics        *metrics

	aggregator aggregator

	cachedPublicKeyGetResponse []byte
	cachedPublicKeyObject      *tdh2easy.PublicKey

	clock clockwork.Clock
}

func (h *handler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *handler) Name() string {
	return h.lggr.Name()
}

type SecretEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type Config struct {
	NodeRateLimiter   ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	RequestTimeoutSec int                         `json:"requestTimeoutSec"`
}

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gwhandlers.DON, capabilitiesRegistry capabilitiesRegistry, requestAuthorizer vaultcap.RequestAuthorizer, lggr logger.Logger, clock clockwork.Clock) (*handler, error) {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal method config: %w", err)
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to create node rate limiter: %w", err)
	}

	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &handler{
		methodConfig:      cfg,
		donConfig:         donConfig,
		don:               don,
		lggr:              logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
		requestTimeout:    time.Duration(cfg.RequestTimeoutSec) * time.Second,
		nodeRateLimiter:   nodeRateLimiter,
		activeRequests:    make(map[string]*activeRequest),
		mu:                sync.RWMutex{},
		requestAuthorizer: requestAuthorizer,
		stopCh:            make(services.StopChan),
		metrics:           metrics,
		aggregator:        &baseAggregator{capabilitiesRegistry: capabilitiesRegistry},
		clock:             clock,
	}, nil
}

func (h *handler) Start(_ context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		go func() {
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
			ticker := h.clock.NewTicker(defaultCleanUpPeriod)
			tickerVaultPublicKeyRefresh := h.clock.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			defer tickerVaultPublicKeyRefresh.Stop()
			for {
				select {
				case <-ticker.Chan():
					h.removeExpiredRequests(ctx)
				case <-tickerVaultPublicKeyRefresh.Chan():
					// periodically, fetch vault public key, so we can cache it
					h.fetchVaultPublicKey(ctx)
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce("VaultHandler", func() error {
		h.lggr.Info("closing vault handler")
		close(h.stopCh)
		return nil
	})
}

func (h *handler) fetchVaultPublicKey(ctx context.Context) {
	ctx, cancel := context.WithDeadline(ctx, h.clock.Now().Add(10*time.Second))
	defer cancel()
	param := vaultcommon.GetPublicKeyRequest{}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		h.lggr.Errorw("fetchVaultPublicKey: failed to marshal get public key request", "error", err)
		return
	}
	getPublicKeyRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uuid.New().String(),
		Method:  vaulttypes.MethodPublicKeyGet,
		Params:  (*json.RawMessage)(&paramBytes),
	}
	h.lggr.Debugw("fetchVaultPublicKey: trying to fetch vault public key", "request", getPublicKeyRequest)
	callback := handlerscommon.NewCallback()
	err = h.HandleJSONRPCUserMessage(ctx, getPublicKeyRequest, callback)
	if err != nil {
		h.lggr.Errorw("fetchVaultPublicKey: failed to fetch vault public key", "request", getPublicKeyRequest, "error", err)
		return
	}
	response, err := callback.Wait(ctx)
	if err != nil {
		h.lggr.Errorw("fetchVaultPublicKey: failed to fetch vault public key", "request", getPublicKeyRequest, "error", err)
		return
	}
	httpStatus := api.ToHttpErrorCode(response.ErrorCode)
	jsonCodec := api.JsonRPCCodec{}
	jsonResp, _ := jsonCodec.DecodeRawRequest(response.RawResponse, "")
	if httpStatus != http.StatusOK {
		h.lggr.Errorw("fetchVaultPublicKey: failed to fetch vault public key", "request", getPublicKeyRequest, "httpStatusCode", httpStatus, "rawResponse", jsonResp)
	} else {
		h.lggr.Debugw("fetchVaultPublicKey: successfully fetched vault public key", "request", getPublicKeyRequest, "rawResponse", jsonResp)
	}
}

// removeExpiredRequests removes expired requests from the pending requests map
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
		var nodeResponses string
		for nodeKey, nodeResponse := range er.responses {
			nodeResponses += fmt.Sprintf("%s ---::: %v               ", nodeKey, nodeResponse)
		}
		err := h.sendResponse(ctx, er, h.errorResponse(er.req, api.RequestTimeoutError, errors.New("request expired without getting quorum of responses from nodes. Available responses: "+nodeResponses), []byte(nodeResponses)))
		if err != nil {
			h.lggr.Errorw("error sending response to user", "requestID", er.req.ID, "error", err)
		}
	}
}

func (h *handler) Methods() []string {
	return vaulttypes.GetSupportedMethods(h.lggr)
}

func (h *handler) HandleLegacyUserMessage(_ context.Context, _ *api.Message, _ gwhandlers.Callback) error {
	return errors.New("vault handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callback gwhandlers.Callback) error {
	// Generate a unique ID for the request.
	// We do this ourselves to ensure the ID is unique and can't be tampered with by the user.
	if req.ID == "" {
		return errors.New("request ID cannot be empty")
	}

	h.lggr.Debugw("handling vault request", "method", req.Method, "requestID", req.ID, "request", req)
	// Public key requests don't require authorization,
	// Let's process this request right away.
	// Note we cache this value quite aggressively so don't need to worry about DoS.
	if req.Method == vaulttypes.MethodPublicKeyGet {
		return h.handlePublicKeyGet(ctx, h.newActiveRequest(req, callback))
	} else if req.Method == vaulttypes.MethodSecretsGet {
		// Secrets get is only allowed in non-production builds for testing purposes
		// So no authorization is required
		ar := h.newActiveRequest(req, callback)
		return h.handleSecretsGet(ctx, ar)
	}

	isAuthorized, owner, err := h.requestAuthorizer.AuthorizeRequest(ctx, req)
	if !isAuthorized {
		h.lggr.Errorw("request not authorized", "requestID", req.ID, "owner", owner, "reason:", err)
		return errors.New("request not authorized: " + err.Error())
	}
	// Prefix request id with owner, to ensure uniqueness across different owners
	req.ID = owner + vaulttypes.RequestIDSeparator + req.ID

	h.lggr.Infow("handling authorized vault request", "method", req.Method, "requestID", req.ID, "owner", owner)
	if h.getActiveRequest(req.ID) != nil {
		h.lggr.Errorw("request id already exists", "requestID", req.ID, "owner", owner)
		return errors.New("request ID already exists: " + req.ID)
	}
	ar := h.newActiveRequest(req, callback)

	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		return h.handleSecretsCreate(ctx, ar)
	case vaulttypes.MethodSecretsUpdate:
		return h.handleSecretsUpdate(ctx, ar)
	case vaulttypes.MethodSecretsDelete:
		return h.handleSecretsDelete(ctx, ar)
	case vaulttypes.MethodSecretsList:
		return h.handleSecretsList(ctx, ar)
	default:
		return h.sendResponse(ctx, ar, h.errorResponse(req, api.UnsupportedMethodError, errors.New("this method is unsupported: "+req.Method), nil))
	}
}

func (h *handler) newActiveRequest(req jsonrpc.Request[json.RawMessage], callback gwhandlers.Callback) *activeRequest {
	ar := &activeRequest{
		Callback:  callback,
		req:       req,
		createdAt: h.clock.Now(),
		responses: map[string]*jsonrpc.Response[json.RawMessage]{},
	}

	h.mu.Lock()
	h.activeRequests[req.ID] = ar
	h.mu.Unlock()

	return ar
}

func (h *handler) getActiveRequest(requestID string) *activeRequest {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.activeRequests[requestID]
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	l := logger.With(h.lggr, "method", resp.Method, "requestID", resp.ID, "nodeAddr", nodeAddr)
	l.Debugw("handling node response")

	if !h.nodeRateLimiter.Allow(nodeAddr) {
		l.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}

	ar := h.getActiveRequest(resp.ID)
	if ar == nil {
		// Request is not found, so we don't need to send a response to the user
		// This might happen if the response is stale
		l.Errorw("no pending request found for ID")
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", api.StaleNodeResponseError.String()),
		))
		return nil
	}

	ok := ar.addResponseForNode(nodeAddr, resp)
	if !ok {
		l.Errorw("duplicate response from node, ignoring", "nodeAddr", nodeAddr)
		return nil
	}

	resp, err := h.aggregator.Aggregate(ctx, l, ar.copiedResponses(), resp)
	switch {
	case errors.Is(err, errInsufficientResponsesForQuorum):
		l.Debugw("aggregating responses, waiting for other nodes...", "error", err)
		return nil
	case err != nil:
		l.Error("quorum unobtainable, returning response to user...", "error", err, "responses", maps.Values(ar.responses))
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, err, nil))
	}

	switch resp.Method {
	case vaulttypes.MethodPublicKeyGet:
		h.tryCachePublicKeyResponse(resp, l)
	default:
		// Do nothing for other methods
	}

	return h.sendSuccessResponse(ctx, l, ar, resp)
}

func (h *handler) tryCachePublicKeyResponse(resp *jsonrpc.Response[json.RawMessage], l logger.Logger) {
	if resp.Result == nil {
		l.Infow("no result in public key response, not caching")
		return
	}

	r := &vaultcommon.GetPublicKeyResponse{}
	err := json.Unmarshal(*resp.Result, r)
	if err != nil {
		l.Infow("failed to unmarshal public key response, not caching", "error", err)
		return
	}

	if r.PublicKey == "" {
		l.Infow("no public key in unmarshaled response, not caching", "response", resp, "result", r)
		return
	}
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(r.PublicKey)
	if err != nil {
		l.Infow("failed to decode master public key string", "error", err)
		return
	}
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	if err != nil {
		l.Infow("failed to unmarshal master public key", "error", err)
		return
	}

	h.mu.Lock()
	h.cachedPublicKeyGetResponse = *resp.Result
	h.cachedPublicKeyObject = &masterPublicKey
	h.mu.Unlock()
	l.Infow("successfully cached public key response")
}

func (h *handler) sendSuccessResponse(ctx context.Context, l logger.Logger, ar *activeRequest, resp *jsonrpc.Response[json.RawMessage]) error {
	// Strip the owner prefix from the response ID before sending it back to the user
	// This ensures compliance with JSONRPC 2.0 spec, which requires response id to match request id
	index := strings.Index(resp.ID, vaulttypes.RequestIDSeparator)
	if index != -1 {
		resp.ID = resp.ID[index+2:]
	}
	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		l.Errorw("failed to encode response", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err), nil))
	}

	var errorCode api.ErrorCode
	if resp.Error != nil {
		errorCode = api.FromJSONRPCErrorCode(resp.Error.Code)
	} else {
		errorCode = api.NoError
	}

	l.Debugw("issued user callback", "errorCode", errorCode)
	successResp := gwhandlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   errorCode,
	}
	return h.sendResponse(ctx, ar, successResp)
}

func (h *handler) handleSecretsCreate(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	createSecretsRequest := &vaultcommon.CreateSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, &createSecretsRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err, nil))
	}
	createSecretsRequest.RequestId = ar.req.ID
	for _, secretItem := range createSecretsRequest.EncryptedSecrets {
		if secretItem.Id != nil && secretItem.Id.Namespace == "" {
			secretItem.Id.Namespace = vaulttypes.DefaultNamespace
		}
	}
	_, cachedPublicKey, _ := h.getCachedPublicKey()
	err := vaultcap.ValidateCreateSecretsRequest(cachedPublicKey, createSecretsRequest)
	if err != nil {
		l.Warnw("failed to validate create secrets request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, fmt.Errorf("failed to validate create secrets request: %w", err), nil))
	}

	reqBytes, err := json.Marshal(createSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err), nil))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	// At this point, we know that the request is valid, and we can send it to the nodes
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsUpdate(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	updateSecretsRequest := &vaultcommon.UpdateSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, updateSecretsRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err, nil))
	}

	updateSecretsRequest.RequestId = ar.req.ID
	for _, secretItem := range updateSecretsRequest.EncryptedSecrets {
		if secretItem.Id != nil && secretItem.Id.Namespace == "" {
			secretItem.Id.Namespace = vaulttypes.DefaultNamespace
		}
	}
	_, cachedPublicKey, _ := h.getCachedPublicKey()
	vaultCapErr := vaultcap.ValidateUpdateSecretsRequest(cachedPublicKey, updateSecretsRequest)
	if vaultCapErr != nil {
		l.Warnw("failed to validate update secrets request", "error", vaultCapErr)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, fmt.Errorf("failed to validate update secrets request: %w", vaultCapErr), nil))
	}

	reqBytes, err := json.Marshal(updateSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err), nil))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsDelete(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	deleteSecretsRequest := &vaultcommon.DeleteSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, deleteSecretsRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err, nil))
	}

	deleteSecretsRequest.RequestId = ar.req.ID
	for _, id := range deleteSecretsRequest.Ids {
		if id != nil && id.Namespace == "" {
			id.Namespace = vaulttypes.DefaultNamespace
		}
	}
	err := vaultcap.ValidateDeleteSecretsRequest(deleteSecretsRequest)
	if err != nil {
		l.Warnw("failed to validate delete secrets request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, fmt.Errorf("failed to validate delete secrets request: %w", err), nil))
	}

	reqBytes, err := json.Marshal(deleteSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err), nil))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsGet(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	secretsGetRequest := &vaultcommon.GetSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, &secretsGetRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err, nil))
	}
	for _, getRequest := range secretsGetRequest.Requests {
		if getRequest.Id != nil && getRequest.Id.Namespace == "" {
			getRequest.Id.Namespace = vaulttypes.DefaultNamespace
		}
	}
	err := vaultcap.ValidateGetSecretsRequest(secretsGetRequest)
	if err != nil {
		l.Warnw("failed to validate get secrets request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, fmt.Errorf("failed to validate get secrets request: %w", err), nil))
	}

	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsList(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	req := &vaultcommon.ListSecretIdentifiersRequest{}
	if err := json.Unmarshal(*ar.req.Params, req); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err, nil))
	}

	req.RequestId = ar.req.ID
	if req.Namespace == "" {
		req.Namespace = vaulttypes.DefaultNamespace
	}
	err := vaultcap.ValidateListSecretIdentifiersRequest(req)
	if err != nil {
		l.Warnw("failed to validate list secret identifiers request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, fmt.Errorf("failed to validate list secret identifiers request: %w", err), nil))
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err), nil))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) getCachedPublicKey() ([]byte, *tdh2easy.PublicKey, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.cachedPublicKeyGetResponse == nil {
		return nil, nil, errors.New("no cached public key response")
	}
	copied := make([]byte, len(h.cachedPublicKeyGetResponse))
	copy(copied, h.cachedPublicKeyGetResponse)
	cachedPublicKeyCopy := *h.cachedPublicKeyObject
	return copied, &cachedPublicKeyCopy, nil
}

func (h *handler) handlePublicKeyGet(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	publicKeyResponseBytes, _, err := h.getCachedPublicKey()
	if err == nil {
		l.Debugw("returning cached public key response")
		return h.sendSuccessResponse(ctx, l, ar, &jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      ar.req.ID,
			Method:  ar.req.Method,
			Result:  (*json.RawMessage)(&publicKeyResponseBytes),
		})
	}

	l.Debugw("cache stale: forwarding request to nodes", "now", h.clock.Now(), "err", err)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) fanOutToVaultNodes(ctx context.Context, l logger.Logger, ar *activeRequest) error {
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, &ar.req)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			l.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes"), nil))
	}

	l.Debugw("successfully forwarded request to Vault nodes")
	return nil
}

func (h *handler) errorResponse(
	req jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	err error,
	data []byte,
) gwhandlers.UserCallbackPayload {
	switch errorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
		h.lggr.Errorw(err.Error(), "requestID", req.ID)
		// Intentionally hide the error from the user
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		h.lggr.Errorw("invalid params", "requestID", req.ID, "params", string(*req.Params))
		err = errors.New("invalid params error: " + err.Error())
	case api.UnsupportedMethodError:
		h.lggr.Errorw("unsupported method", "requestID", req.ID, "method", req.Method)
		err = errors.New("unsupported method: " + req.Method)
	case api.UserMessageParseError:
		h.lggr.Errorw("user message parse error", "requestID", req.ID, "error", err.Error())
		err = errors.New("user message parse error: " + err.Error())
	case api.NoError:
	case api.UnsupportedDONIdError:
	case api.HandlerError:
	case api.RequestTimeoutError:
	case api.StaleNodeResponseError:
		// Unused in this handler
	}

	// Strip the owner prefix from the json response ID before sending it back to the user
	// This ensures compliance with JSONRPC 2.0 spec, which requires response id to match request id
	index := strings.Index(req.ID, vaulttypes.RequestIDSeparator)
	if index != -1 {
		req.ID = req.ID[index+2:]
	}

	return gwhandlers.UserCallbackPayload{
		RawResponse: h.codec.EncodeNewErrorResponse(
			req.ID,
			api.ToJSONRPCErrorCode(errorCode),
			err.Error(),
			data,
		),
		ErrorCode: errorCode,
	}
}

func (h *handler) sendResponse(ctx context.Context, userRequest *activeRequest, resp gwhandlers.UserCallbackPayload) error {
	switch resp.ErrorCode {
	case api.StaleNodeResponseError:
	case api.FatalError:
	case api.NodeReponseEncodingError:
	case api.RequestTimeoutError:
	case api.HandlerError:
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", resp.ErrorCode.String()),
		))
	case api.InvalidParamsError:
	case api.UnsupportedMethodError:
	case api.UserMessageParseError:
	case api.UnsupportedDONIdError:
		h.metrics.requestUserError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	case api.NoError:
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	}

	err := userRequest.SendResponse(resp)
	if err != nil {
		h.lggr.Errorw("error sending response to user", "requestID", userRequest.req.ID, "error", err)
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.activeRequests, userRequest.req.ID)
	h.lggr.Debugw("response sent to user", "requestID", userRequest.req.ID, "errorCode", resp.ErrorCode)
	return nil
}
