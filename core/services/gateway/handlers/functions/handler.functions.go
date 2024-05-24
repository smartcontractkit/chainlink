package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	fallow "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
	fsub "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions"
)

var (
	ErrNotAllowlisted    = errors.New("sender not allowlisted")
	ErrRateLimited       = errors.New("rate-limited")
	ErrUnsupportedMethod = errors.New("unsupported method")

	promHandlerError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_functions_handler_error",
		Help: "Metric to track functions handler errors",
	}, []string{"don_id", "error"})

	promSecretsSetSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_functions_secrets_set_success",
		Help: "Metric to track successful secrets_set calls",
	}, []string{"don_id"})

	promSecretsSetFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_functions_secrets_set_failure",
		Help: "Metric to track failed secrets_set calls",
	}, []string{"don_id"})

	promSecretsListSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_functions_secrets_list_success",
		Help: "Metric to track successful secrets_list calls",
	}, []string{"don_id"})

	promSecretsListFailure = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_functions_secrets_list_failure",
		Help: "Metric to track failed secrets_list calls",
	}, []string{"don_id"})
)

type FunctionsHandlerConfig struct {
	ChainID string `json:"chainId"`
	// Not specifying OnchainAllowlist config disables allowlist checks
	OnchainAllowlist *fallow.OnchainAllowlistConfig `json:"onchainAllowlist"`
	// Not specifying OnchainSubscriptions config disables minimum balance checks
	OnchainSubscriptions       *fsub.OnchainSubscriptionsConfig `json:"onchainSubscriptions"`
	MinimumSubscriptionBalance *assets.Link                     `json:"minimumSubscriptionBalance"`
	// Not specifying RateLimiter config disables rate limiting
	UserRateLimiter            *hc.RateLimiterConfig `json:"userRateLimiter"`
	NodeRateLimiter            *hc.RateLimiterConfig `json:"nodeRateLimiter"`
	MaxPendingRequests         uint32                `json:"maxPendingRequests"`
	RequestTimeoutMillis       int64                 `json:"requestTimeoutMillis"`
	AllowedHeartbeatInitiators []string              `json:"allowedHeartbeatInitiators"`
}

type functionsHandler struct {
	services.StateMachine

	handlerConfig              FunctionsHandlerConfig
	donConfig                  *config.DONConfig
	don                        handlers.DON
	pendingRequests            hc.RequestCache[PendingRequest]
	allowlist                  fallow.OnchainAllowlist
	subscriptions              fsub.OnchainSubscriptions
	minimumBalance             *assets.Link
	userRateLimiter            *hc.RateLimiter
	nodeRateLimiter            *hc.RateLimiter
	allowedHeartbeatInitiators map[string]struct{}
	chStop                     services.StopChan
	lggr                       logger.Logger
}

type PendingRequest struct {
	request    *api.Message
	responses  map[string]*api.Message
	successful []*api.Message
	errors     []*api.Message
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandlerFromConfig(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, lggr logger.Logger) (handlers.Handler, error) {
	var cfg FunctionsHandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	lggr = lggr.Named("FunctionsHandler:" + donConfig.DonId)
	var allowlist fallow.OnchainAllowlist
	if cfg.OnchainAllowlist != nil {
		chain, err2 := legacyChains.Get(cfg.ChainID)
		if err2 != nil {
			return nil, err2
		}

		orm, err2 := fallow.NewORM(ds, lggr, cfg.OnchainAllowlist.ContractAddress)
		if err2 != nil {
			return nil, err2
		}
		allowlist, err2 = fallow.NewOnchainAllowlist(chain.Client(), *cfg.OnchainAllowlist, orm, lggr)
		if err2 != nil {
			return nil, err2
		}
	}
	var userRateLimiter, nodeRateLimiter *hc.RateLimiter
	if cfg.UserRateLimiter != nil {
		userRateLimiter, err = hc.NewRateLimiter(*cfg.UserRateLimiter)
		if err != nil {
			return nil, err
		}
	}
	if cfg.NodeRateLimiter != nil {
		nodeRateLimiter, err = hc.NewRateLimiter(*cfg.NodeRateLimiter)
		if err != nil {
			return nil, err
		}
	}
	var subscriptions fsub.OnchainSubscriptions
	if cfg.OnchainSubscriptions != nil {
		chain, err2 := legacyChains.Get(cfg.ChainID)
		if err2 != nil {
			return nil, err2
		}

		orm, err2 := fsub.NewORM(ds, lggr, cfg.OnchainSubscriptions.ContractAddress)
		if err2 != nil {
			return nil, err2
		}

		subscriptions, err2 = fsub.NewOnchainSubscriptions(chain.Client(), *cfg.OnchainSubscriptions, orm, lggr)
		if err2 != nil {
			return nil, err2
		}
	}
	allowedHeartbeatInitiators := make(map[string]struct{})
	for _, initiator := range cfg.AllowedHeartbeatInitiators {
		allowedHeartbeatInitiators[strings.ToLower(initiator)] = struct{}{}
	}
	pendingRequestsCache := hc.NewRequestCache[PendingRequest](time.Millisecond*time.Duration(cfg.RequestTimeoutMillis), cfg.MaxPendingRequests)
	return NewFunctionsHandler(cfg, donConfig, don, pendingRequestsCache, allowlist, subscriptions, cfg.MinimumSubscriptionBalance, userRateLimiter, nodeRateLimiter, allowedHeartbeatInitiators, lggr), nil
}

func NewFunctionsHandler(
	cfg FunctionsHandlerConfig,
	donConfig *config.DONConfig,
	don handlers.DON,
	pendingRequestsCache hc.RequestCache[PendingRequest],
	allowlist fallow.OnchainAllowlist,
	subscriptions fsub.OnchainSubscriptions,
	minimumBalance *assets.Link,
	userRateLimiter *hc.RateLimiter,
	nodeRateLimiter *hc.RateLimiter,
	allowedHeartbeatInitiators map[string]struct{},
	lggr logger.Logger) handlers.Handler {
	return &functionsHandler{
		handlerConfig:              cfg,
		donConfig:                  donConfig,
		don:                        don,
		pendingRequests:            pendingRequestsCache,
		allowlist:                  allowlist,
		subscriptions:              subscriptions,
		minimumBalance:             minimumBalance,
		userRateLimiter:            userRateLimiter,
		nodeRateLimiter:            nodeRateLimiter,
		allowedHeartbeatInitiators: allowedHeartbeatInitiators,
		chStop:                     make(services.StopChan),
		lggr:                       lggr,
	}
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	sender := common.HexToAddress(msg.Body.Sender)
	if h.allowlist != nil && !h.allowlist.Allow(sender) {
		h.lggr.Debugw("received a message from a non-allowlisted address", "sender", msg.Body.Sender)
		promHandlerError.WithLabelValues(h.donConfig.DonId, ErrNotAllowlisted.Error()).Inc()
		return ErrNotAllowlisted
	}
	if h.userRateLimiter != nil && !h.userRateLimiter.Allow(msg.Body.Sender) {
		h.lggr.Debugw("rate-limited", "sender", msg.Body.Sender)
		promHandlerError.WithLabelValues(h.donConfig.DonId, ErrRateLimited.Error()).Inc()
		return ErrRateLimited
	}
	if msg.Body.Method == MethodSecretsSet && h.subscriptions != nil && h.minimumBalance != nil {
		balance, err := h.subscriptions.GetMaxUserBalance(sender)
		if err != nil {
			h.lggr.Debugw("error getting max user balance", "sender", msg.Body.Sender, "err", err)
		}
		if balance == nil {
			balance = big.NewInt(0)
		}
		if err != nil || balance.Cmp(h.minimumBalance.ToInt()) < 0 {
			h.lggr.Debugw("received a message from a user having insufficient balance", "sender", msg.Body.Sender, "balance", balance.String())
			return fmt.Errorf("sender has insufficient balance: %v juels", balance.String())
		}
	}
	switch msg.Body.Method {
	case MethodSecretsSet, MethodSecretsList:
		return h.handleRequest(ctx, msg, callbackCh)
	case MethodHeartbeat:
		if _, ok := h.allowedHeartbeatInitiators[msg.Body.Sender]; !ok {
			h.lggr.Debugw("received heartbeat request from a non-allowed sender", "sender", msg.Body.Sender)
			promHandlerError.WithLabelValues(h.donConfig.DonId, ErrNotAllowlisted.Error()).Inc()
			return ErrUnsupportedMethod
		}
		return h.handleRequest(ctx, msg, callbackCh)
	default:
		h.lggr.Debugw("unsupported method", "method", msg.Body.Method)
		promHandlerError.WithLabelValues(h.donConfig.DonId, ErrUnsupportedMethod.Error()).Inc()
		return ErrUnsupportedMethod
	}
}

func (h *functionsHandler) handleRequest(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handleRequest: processing message", "sender", msg.Body.Sender, "messageId", msg.Body.MessageId)
	err := h.pendingRequests.NewRequest(msg, callbackCh, &PendingRequest{request: msg, responses: make(map[string]*api.Message)})
	if err != nil {
		h.lggr.Warnw("handleRequest: error adding new request", "sender", msg.Body.Sender, "err", err)
		promHandlerError.WithLabelValues(h.donConfig.DonId, err.Error()).Inc()
		return err
	}
	// Send to all nodes.
	for _, member := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, member.Address, msg)
		if err != nil {
			h.lggr.Debugw("handleRequest: failed to send to a node", "node", member.Address, "err", err)
		}
	}
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("HandleNodeMessage: processing message", "nodeAddr", nodeAddr, "receiver", msg.Body.Receiver, "id", msg.Body.MessageId)
	if h.nodeRateLimiter != nil && !h.nodeRateLimiter.Allow(nodeAddr) {
		h.lggr.Debugw("rate-limited", "sender", nodeAddr)
		return errors.New("rate-limited")
	}
	switch msg.Body.Method {
	case MethodSecretsSet, MethodSecretsList:
		return h.pendingRequests.ProcessResponse(msg, h.processSecretsResponse)
	case MethodHeartbeat:
		return h.pendingRequests.ProcessResponse(msg, h.processHeartbeatResponse)
	default:
		h.lggr.Debugw("unsupported method", "method", msg.Body.Method)
		return ErrUnsupportedMethod
	}
}

// Conforms to ResponseProcessor[*PendingRequest]
func (h *functionsHandler) processSecretsResponse(response *api.Message, responseData *PendingRequest) (*handlers.UserCallbackPayload, *PendingRequest, error) {
	if _, exists := responseData.responses[response.Body.Sender]; exists {
		return nil, nil, errors.New("duplicate response")
	}
	if response.Body.Method != responseData.request.Body.Method {
		return nil, responseData, errors.New("invalid method")
	}
	responseData.responses[response.Body.Sender] = response
	var responsePayload ResponseBase
	err := json.Unmarshal(response.Body.Payload, &responsePayload)
	if err != nil {
		responseData.errors = append(responseData.errors, response)
		return nil, responseData, err
	}
	// user response is ready with either F+1 successes or N-F failures
	if responsePayload.Success {
		responseData.successful = append(responseData.successful, response)
		if len(responseData.successful) >= h.donConfig.F+1 {
			// return success to the user
			callbackPayload, err := newSecretsResponse(responseData.request, true, responseData.successful)
			return callbackPayload, responseData, err
		}
	} else {
		responseData.errors = append(responseData.errors, response)
		if len(responseData.errors) >= len(h.donConfig.Members)-h.donConfig.F {
			// return error to the user
			callbackPayload, err := newSecretsResponse(responseData.request, false, responseData.errors)
			return callbackPayload, responseData, err
		}
	}
	// not ready to be processed yet
	return nil, responseData, nil
}

func newSecretsResponse(request *api.Message, success bool, responses []*api.Message) (*handlers.UserCallbackPayload, error) {
	payload := CombinedResponse{ResponseBase: ResponseBase{Success: success}, NodeResponses: responses}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if request.Body.Method == MethodSecretsSet {
		if success {
			promSecretsSetSuccess.WithLabelValues(request.Body.DonId).Inc()
		} else {
			promSecretsSetFailure.WithLabelValues(request.Body.DonId).Inc()
		}
	} else if request.Body.Method == MethodSecretsList {
		if success {
			promSecretsListSuccess.WithLabelValues(request.Body.DonId).Inc()
		} else {
			promSecretsListFailure.WithLabelValues(request.Body.DonId).Inc()
		}
	}

	userResponse := *request
	userResponse.Body.Receiver = request.Body.Sender
	userResponse.Body.Payload = payloadJson
	return &handlers.UserCallbackPayload{Msg: &userResponse, ErrCode: api.NoError, ErrMsg: ""}, nil
}

// Conforms to ResponseProcessor[*PendingRequest]
func (h *functionsHandler) processHeartbeatResponse(response *api.Message, responseData *PendingRequest) (*handlers.UserCallbackPayload, *PendingRequest, error) {
	if _, exists := responseData.responses[response.Body.Sender]; exists {
		return nil, nil, errors.New("duplicate response")
	}
	if response.Body.Method != responseData.request.Body.Method {
		return nil, responseData, errors.New("invalid method")
	}
	responseData.responses[response.Body.Sender] = response

	// user response is ready with F+1 node responses
	if len(responseData.responses) >= h.donConfig.F+1 {
		var responseList []*api.Message
		for _, response := range responseData.responses {
			responseList = append(responseList, response)
		}
		userResponse := *responseData.request
		userResponse.Body.Receiver = responseData.request.Body.Sender
		// success = true only means that we got F+1 responses
		// it's up to the heartbeat sender to validate computation results
		payload := CombinedResponse{ResponseBase: ResponseBase{Success: true}, NodeResponses: responseList}
		payloadJson, err := json.Marshal(payload)
		if err != nil {
			return &handlers.UserCallbackPayload{Msg: &userResponse, ErrCode: api.NodeReponseEncodingError, ErrMsg: ""}, nil, nil
		}
		userResponse.Body.Payload = payloadJson
		return &handlers.UserCallbackPayload{Msg: &userResponse, ErrCode: api.NoError, ErrMsg: ""}, nil, nil
	}
	// not ready to be processed yet
	return nil, responseData, nil
}

func (h *functionsHandler) Start(ctx context.Context) error {
	return h.StartOnce("FunctionsHandler", func() error {
		h.lggr.Info("starting FunctionsHandler")
		if h.allowlist != nil {
			if err := h.allowlist.Start(ctx); err != nil {
				return err
			}
		}
		if h.subscriptions != nil {
			if err := h.subscriptions.Start(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (h *functionsHandler) Close() error {
	return h.StopOnce("FunctionsHandler", func() (err error) {
		close(h.chStop)
		if h.allowlist != nil {
			err = multierr.Combine(err, h.allowlist.Close())
		}
		if h.subscriptions != nil {
			err = multierr.Combine(err, h.subscriptions.Close())
		}
		return
	})
}
