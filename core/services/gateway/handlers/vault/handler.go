package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"

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
)

const (
	defaultCleanUpPeriod = 5 * time.Second
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

	createdAt  time.Time
	callbackCh chan<- gwhandlers.UserCallbackPayload
}

type capabilitiesRegistry interface {
	DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error)
}

type aggregator interface {
	Aggregate(ctx context.Context, l logger.Logger, ar *activeRequest, currResp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error)
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

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gwhandlers.DON, capabilitiesRegistry capabilitiesRegistry, requestAuthorizer vaultcap.RequestAuthorizer, lggr logger.Logger) (*handler, error) {
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
	}, nil
}

func (h *handler) Start(_ context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		go func() {
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
			ticker := time.NewTicker(defaultCleanUpPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
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
	return h.StopOnce("VaultHandler", func() error {
		h.lggr.Info("closing vault handler")
		close(h.stopCh)
		return nil
	})
}

// removeExpiredRequests removes expired requests from the pending requests map
func (h *handler) removeExpiredRequests(ctx context.Context) {
	h.mu.RLock()
	var expiredRequests []*activeRequest
	now := time.Now()
	for _, userRequest := range h.activeRequests {
		if now.Sub(userRequest.createdAt) > h.requestTimeout {
			expiredRequests = append(expiredRequests, userRequest)
		}
	}
	h.mu.RUnlock()

	for _, er := range expiredRequests {
		err := h.sendResponse(ctx, er, h.errorResponse(er.req, api.RequestTimeoutError, errors.New("request expired without getting any response")))
		if err != nil {
			h.lggr.Errorw("error sending response to user", "request_id", er.req.ID, "error", err)
		}
	}
}

func (h *handler) Methods() []string {
	return vaulttypes.Methods
}

func (h *handler) HandleLegacyUserMessage(_ context.Context, _ *api.Message, _ chan<- gwhandlers.UserCallbackPayload) error {
	return errors.New("vault handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callbackCh chan<- gwhandlers.UserCallbackPayload) error {
	// Generate a unique ID for the request.
	// We do this ourselves to ensure the ID is unique and can't be tampered with by the user.
	if req.ID == "" {
		return errors.New("request ID cannot be empty")
	}

	isAuthorized, owner, err := h.requestAuthorizer.AuthorizeRequest(ctx, req)
	if !isAuthorized {
		h.lggr.Errorw("request not authorized", "request_id", req.ID, "owner", owner, "reason:", err)
		return errors.New("request not authorized: " + err.Error())
	}

	// Prefix request id with owner, to ensure uniqueness across different owners
	req.ID = owner + "::" + req.ID

	h.lggr.Debugw("handling authorized vault request", "method", req.Method, "requestID", req.ID, "owner", owner)
	ar := &activeRequest{
		callbackCh: callbackCh,
		req:        req,
		createdAt:  time.Now(),
		responses:  map[string]*jsonrpc.Response[json.RawMessage]{},
	}

	h.mu.Lock()
	h.activeRequests[req.ID] = ar
	h.mu.Unlock()
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		return h.handleSecretsCreate(ctx, ar)
	case vaulttypes.MethodSecretsGet:
		return h.handleSecretsGet(ctx, ar)
	case vaulttypes.MethodSecretsUpdate:
		return h.handleSecretsUpdate(ctx, ar)
	case vaulttypes.MethodSecretsDelete:
		return h.handleSecretsDelete(ctx, ar)
	case vaulttypes.MethodSecretsList:
		return h.handleSecretsList(ctx, ar)
	default:
		return h.sendResponse(ctx, ar, h.errorResponse(req, api.UnsupportedMethodError, errors.New("this method is unsupported: "+req.Method)))
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	l := logger.With(h.lggr, "method", resp.Method, "requestID", resp.ID, "nodeAddr", nodeAddr)
	l.Debugw("handling node response")

	if !h.nodeRateLimiter.Allow(nodeAddr) {
		l.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}

	h.mu.RLock()
	ar, ok := h.activeRequests[resp.ID]
	h.mu.RUnlock()
	if !ok {
		// Request is not found, so we don't need to send a response to the user
		// This might happen if the response is stale
		l.Errorw("no pending request found for ID")
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", api.StaleNodeResponseError.String()),
		))
		return nil
	}

	ar.mu.Lock()
	_, exists := ar.responses[nodeAddr]
	if exists {
		l.Errorw("duplicate response from node, ignoring", "nodeAddr", nodeAddr)
		ar.mu.Unlock()
		return nil
	}
	ar.responses[nodeAddr] = resp
	ar.mu.Unlock()

	resp, err := h.aggregator.Aggregate(ctx, l, ar, resp)
	switch {
	case errors.Is(err, errQuorumUnobtainable):
		l.Error("quorum unobtainable, returning response to user...", "error", err, "responses", maps.Values(ar.responses))
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, err))
	case err != nil:
		l.Debugw("error aggregating responses, waiting for other nodes...", "error", err)
		return nil
	}

	return h.sendSuccessResponse(ctx, l, ar, resp)
}

func (h *handler) sendSuccessResponse(ctx context.Context, l logger.Logger, ar *activeRequest, resp *jsonrpc.Response[json.RawMessage]) error {
	// Strip the owner prefix from the response ID before sending it back to the user
	// This ensures compliance with JSONRPC 2.0 spec, which requires response id to match request id
	index := strings.Index(resp.ID, "::")
	if index != -1 {
		resp.ID = resp.ID[index+2:]
	}
	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		l.Errorw("failed to encode response", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err)))
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
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	createSecretsRequest.RequestId = ar.req.ID
	if createSecretsRequest.RequestId == "" {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request_id cannot be empty")))
	}
	if len(createSecretsRequest.EncryptedSecrets) == 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("must have at least 1 request")))
	}
	if len(createSecretsRequest.EncryptedSecrets) >= vaulttypes.MaxBatchSize {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request batch size exceeds maximum of "+strconv.Itoa(vaulttypes.MaxBatchSize))))
	}
	for index, secret := range createSecretsRequest.EncryptedSecrets {
		if secret.Id.Key == "" || secret.EncryptedValue == "" || secret.Id.Owner == "" {
			return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id key, owner and EncryptedValue cannot be empty on index "+strconv.Itoa(index))))
		}
	}

	reqBytes, err := json.Marshal(createSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	// At this point, we know that the request is valid, and we can send it to the nodes
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsUpdate(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	updateSecretsRequest := &vaultcommon.UpdateSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, updateSecretsRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	updateSecretsRequest.RequestId = ar.req.ID

	if updateSecretsRequest.RequestId == "" {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request_id cannot be empty")))
	}
	if len(updateSecretsRequest.EncryptedSecrets) == 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("must have least 1 request")))
	}
	if len(updateSecretsRequest.EncryptedSecrets) >= vaulttypes.MaxBatchSize {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request batch size exceeds maximum of "+strconv.Itoa(vaulttypes.MaxBatchSize))))
	}

	for index, secret := range updateSecretsRequest.EncryptedSecrets {
		if secret.Id.Key == "" || secret.Id.Owner == "" {
			l.Debugw("invalid request parameters: secret id and owner cannot be empty")
			return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id key and owner cannot be empty on index "+strconv.Itoa(index))))
		}
	}

	reqBytes, err := json.Marshal(updateSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsDelete(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestId", ar.req.ID)

	deleteSecretsRequest := &vaultcommon.DeleteSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, deleteSecretsRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	deleteSecretsRequest.RequestId = ar.req.ID

	if deleteSecretsRequest.RequestId == "" {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request_id cannot be empty")))
	}
	if len(deleteSecretsRequest.Ids) == 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("must have least 1 request")))
	}
	if len(deleteSecretsRequest.Ids) >= vaulttypes.MaxBatchSize {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request batch size exceeds maximum of "+strconv.Itoa(vaulttypes.MaxBatchSize))))
	}

	for index, id := range deleteSecretsRequest.Ids {
		if id.Key == "" || id.Owner == "" {
			l.Debugw("invalid request parameters: secret id and owner cannot be empty")
			return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id key and owner cannot be empty on index "+strconv.Itoa(index))))
		}
	}

	reqBytes, err := json.Marshal(deleteSecretsRequest)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsGet(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	secretsGetRequest := &vaultcommon.GetSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, &secretsGetRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	if len(secretsGetRequest.Requests) == 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("must have least 1 request")))
	}
	if len(secretsGetRequest.Requests) >= vaulttypes.MaxBatchSize {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("request batch size exceeds maximum of "+strconv.Itoa(vaulttypes.MaxBatchSize))))
	}
	for index, request := range secretsGetRequest.Requests {
		if request.Id.Key == "" || request.Id.Owner == "" {
			l.Debugw("invalid request parameters: secret id and owner cannot be empty")
			return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id key and owner cannot be empty on index "+strconv.Itoa(index))))
		}
	}

	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsList(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestId", ar.req.ID)

	req := &vaultcommon.ListSecretIdentifiersRequest{}
	if err := json.Unmarshal(*ar.req.Params, req); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	req.RequestId = ar.req.ID

	reqBytes, err := json.Marshal(req)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqBytes)
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
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to Vault nodes")
	return nil
}

func (h *handler) errorResponse(
	req jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	errs ...error,
) gwhandlers.UserCallbackPayload {
	err := errors.New("unknown error")
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}

	switch errorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
		h.lggr.Errorw(err.Error(), "request_id", req.ID)
		// Intentionally hide the error from the user
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		h.lggr.Errorw("invalid params", "request_id", req.ID, "params", string(*req.Params))
		err = errors.New("invalid params error: " + err.Error())
	case api.UnsupportedMethodError:
		h.lggr.Errorw("unsupported method", "request_id", req.ID, "method", req.Method)
		err = errors.New("unsupported method: " + req.Method)
	case api.UserMessageParseError:
		h.lggr.Errorw("user message parse error", "request_id", req.ID, "error", err.Error())
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
	index := strings.Index(req.ID, "::")
	if index != -1 {
		req.ID = req.ID[index+2:]
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

	select {
	case userRequest.callbackCh <- resp:
		h.lggr.Debugw("sent response", "request_id", userRequest.req.ID, "error_code", resp.ErrorCode, "raw_response", string(resp.RawResponse))
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return nil
	case <-ctx.Done():
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return ctx.Err()
	}
}
