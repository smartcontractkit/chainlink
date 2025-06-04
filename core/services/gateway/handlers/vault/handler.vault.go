package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

var (
	ErrNotAllowlisted    = errors.New("sender not allowlisted")
	ErrRateLimited       = errors.New("rate-limited")
	ErrUnsupportedMethod = errors.New("unsupported method")

	promHandlerError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_handler_error",
		Help: "Metric to track vault handler errors",
	}, []string{"don_id", "error"})

	promSecretsCreateSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_secrets_create_success",
		Help: "Metric to track successful vault secrets_create calls",
	}, []string{"don_id"})

	promSecretsCreateFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_vault_secrets_create_failure",
		Help: "Metric to track failed vault secrets_create calls",
	}, []string{"don_id"})
)

type PendingRequest struct {
	request    *api.Message
	responses  map[string]*api.Message
	successful []*api.Message
	errors     []*api.Message
}

type vaultHandler struct {
	services.StateMachine

	handlerConfig VaultConfig
	donConfig     *config.DONConfig
	don           handlers.DON
	lggr          logger.Logger

	mu                sync.RWMutex
	pendingRequests   map[string]*PendingRequest
	userRateLimiter   *hc.RateLimiter
	nodeRateLimiter   *hc.RateLimiter
	requestTimeoutSec int

	// In-memory secret storage for demo purposes
	// In production, this would be a proper storage backend
	secretsStore map[string]map[string]SecretEntry
	storeMu      sync.RWMutex
}

type SecretEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type VaultConfig struct {
	UserRateLimiterConfig hc.RateLimiterConfig `json:"user_rate_limiter"`
	NodeRateLimiterConfig hc.RateLimiterConfig `json:"node_rate_limiter"`
	RequestTimeoutSec     int                  `json:"request_timeout_sec"`
}

var _ handlers.Handler = (*vaultHandler)(nil)

func NewVaultHandlerFromConfig(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, lggr logger.Logger) (handlers.Handler, error) {
	var cfg VaultConfig
	if err := json.Unmarshal(handlerConfig, &cfg); err != nil {
		return nil, err
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	return &vaultHandler{
		handlerConfig:     cfg,
		donConfig:         donConfig,
		don:               don,
		lggr:              lggr.Named("VaultHandler"),
		pendingRequests:   make(map[string]*PendingRequest),
		requestTimeoutSec: cfg.RequestTimeoutSec,
		secretsStore:      make(map[string]map[string]SecretEntry),
	}, nil
}

func (h *vaultHandler) Start(ctx context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		go h.cleanupExpiredRequests()
		return nil
	})
}

func (h *vaultHandler) Close() error {
	return h.StopOnce("VaultHandler", func() error {
		h.lggr.Info("closing vault handler")
		return nil
	})
}

func (h *vaultHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	switch msg.Body.Method {
	case MethodSecretsCreate:
		return h.handleRequest(ctx, msg, callbackCh)
	default:
		h.lggr.Debugw("unsupported method", "method", msg.Body.Method)
		promHandlerError.WithLabelValues(h.donConfig.DonId, ErrUnsupportedMethod.Error()).Inc()
		return ErrUnsupportedMethod
	}
}

func (h *vaultHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	// Vault handler processes requests directly, doesn't forward to nodes
	h.lggr.Debugw("received unexpected node message", "nodeAddr", nodeAddr, "method", msg.Body.Method)
	return nil
}

func (h *vaultHandler) handleRequest(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handling vault request", "method", msg.Body.Method, "sender", msg.Body.Sender)

	// For demonstration, we'll handle vault.secrets.create directly in the gateway
	// In a real implementation, this would be forwarded to nodes
	if msg.Body.Method == MethodSecretsCreate {
		return h.handleSecretsCreateDirect(ctx, msg, callbackCh)
	}

	// For other methods, forward to nodes (placeholder implementation)
	pendingRequest := &PendingRequest{
		request:   msg,
		responses: make(map[string]*api.Message),
	}

	h.mu.Lock()
	h.pendingRequests[msg.Body.MessageId] = pendingRequest
	h.mu.Unlock()

	// Forward to nodes
	for _, member := range h.donConfig.Members {
		nodeMsg := *msg
		nodeMsg.Body.Receiver = member.Address
		if err := h.don.SendToNode(ctx, member.Address, &nodeMsg); err != nil {
			h.lggr.Errorw("failed to send message to node", "node", member.Address, "err", err)
		}
	}

	// Set timeout for request
	go func() {
		timer := time.NewTimer(time.Duration(h.requestTimeoutSec) * time.Second)
		defer timer.Stop()

		select {
		case <-timer.C:
			h.mu.Lock()
			if _, exists := h.pendingRequests[msg.Body.MessageId]; exists {
				delete(h.pendingRequests, msg.Body.MessageId)
				h.mu.Unlock()

				callbackPayload := handlers.UserCallbackPayload{
					Msg:     msg,
					ErrCode: api.RequestTimeoutError,
					ErrMsg:  "request timeout",
				}
				callbackCh <- callbackPayload
			} else {
				h.mu.Unlock()
			}
		case <-ctx.Done():
			return
		}
	}()

	return nil
}

func (h *vaultHandler) handleSecretsCreateDirect(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	var request SecretsCreateRequest
	if err := json.Unmarshal(msg.Body.Payload, &request); err != nil {
		response := SecretsCreateResponse{
			ResponseBase: ResponseBase{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Failed to parse request: %v", err),
			},
		}
		return h.sendResponse(ctx, msg, response, callbackCh)
	}

	// Validate request
	if request.ID == "" {
		response := SecretsCreateResponse{
			ResponseBase: ResponseBase{
				Success:      false,
				ErrorMessage: "Secret ID cannot be empty",
			},
		}
		return h.sendResponse(ctx, msg, response, callbackCh)
	}

	// Store secret
	h.storeMu.Lock()
	defer h.storeMu.Unlock()

	senderAddr := msg.Body.Sender
	if h.secretsStore[senderAddr] == nil {
		h.secretsStore[senderAddr] = make(map[string]SecretEntry)
	}

	// Check if secret already exists
	if _, exists := h.secretsStore[senderAddr][request.ID]; exists {
		response := SecretsCreateResponse{
			ResponseBase: ResponseBase{
				Success:      false,
				ErrorMessage: "Secret with this ID already exists",
			},
		}
		return h.sendResponse(ctx, msg, response, callbackCh)
	}

	// Create new secret
	secret := SecretEntry{
		ID:        request.ID,
		Value:     request.Value,
		CreatedAt: time.Now().Unix(),
	}

	h.secretsStore[senderAddr][request.ID] = secret

	response := SecretsCreateResponse{
		ResponseBase: ResponseBase{
			Success: true,
		},
		ID: request.ID,
	}

	promSecretsCreateSuccess.WithLabelValues(h.donConfig.DonId).Inc()
	return h.sendResponse(ctx, msg, response, callbackCh)
}

func (h *vaultHandler) sendResponse(ctx context.Context, originalMsg *api.Message, response interface{}, callbackCh chan<- handlers.UserCallbackPayload) error {
	responsePayload, err := json.Marshal(response)
	if err != nil {
		promSecretsCreateFailure.WithLabelValues(h.donConfig.DonId).Inc()
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	responseMsg := *originalMsg
	responseMsg.Body.Receiver = originalMsg.Body.Sender
	responseMsg.Body.Payload = responsePayload

	callbackPayload := handlers.UserCallbackPayload{
		Msg:     &responseMsg,
		ErrCode: api.NoError,
		ErrMsg:  "",
	}

	select {
	case callbackCh <- callbackPayload:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *vaultHandler) cleanupExpiredRequests() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	timeout := time.Duration(h.requestTimeoutSec) * time.Second

	for id := range h.pendingRequests {
		// Estimate request time from message ID or use a more sophisticated approach
		if now.Sub(time.Unix(0, 0)) > timeout {
			delete(h.pendingRequests, id)
			h.lggr.Debugw("cleaned up expired request", "id", id)
		}
	}
}
