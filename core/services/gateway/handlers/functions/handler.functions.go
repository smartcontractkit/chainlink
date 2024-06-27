package functions

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsHandlerConfig struct {
	OnchainAllowlistChainID string `json:"onchainAllowlistChainId"`
	// Not specifying OnchainAllowlist config disables allowlist checks
	OnchainAllowlist *OnchainAllowlistConfig `json:"onchainAllowlist"`
	// Not specifying RateLimiter config disables rate limiting
	UserRateLimiter      *hc.RateLimiterConfig `json:"userRateLimiter"`
	NodeRateLimiter      *hc.RateLimiterConfig `json:"nodeRateLimiter"`
	MaxPendingRequests   uint32                `json:"maxPendingRequests"`
	RequestTimeoutMillis int64                 `json:"requestTimeoutMillis"`
}

type functionsHandler struct {
	utils.StartStopOnce

	handlerConfig   FunctionsHandlerConfig
	donConfig       *config.DONConfig
	don             handlers.DON
	pendingRequests hc.RequestCache[PendingSecretsRequest]
	allowlist       OnchainAllowlist
	userRateLimiter *hc.RateLimiter
	nodeRateLimiter *hc.RateLimiter
	chStop          utils.StopChan
	lggr            logger.Logger
}

type PendingSecretsRequest struct {
	request    *api.Message
	responses  map[string]*api.Message
	successful []*api.Message
	errors     []*api.Message
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandlerFromConfig(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, chains evm.ChainSet, lggr logger.Logger) (handlers.Handler, error) {
	var cfg FunctionsHandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	var allowlist OnchainAllowlist
	if cfg.OnchainAllowlist != nil {
		chainId, ok := big.NewInt(0).SetString(cfg.OnchainAllowlistChainID, 10)
		if !ok {
			return nil, errors.New("invalid chain ID")
		}
		chain, err := chains.Get(chainId)
		if err != nil {
			return nil, err
		}
		allowlist, err = NewOnchainAllowlist(chain.Client(), *cfg.OnchainAllowlist, lggr)
		if err != nil {
			return nil, err
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
	pendingRequestsCache := hc.NewRequestCache[PendingSecretsRequest](time.Millisecond*time.Duration(cfg.RequestTimeoutMillis), cfg.MaxPendingRequests)
	return NewFunctionsHandler(cfg, donConfig, don, pendingRequestsCache, allowlist, userRateLimiter, nodeRateLimiter, lggr), nil
}

func NewFunctionsHandler(
	cfg FunctionsHandlerConfig,
	donConfig *config.DONConfig,
	don handlers.DON,
	pendingRequestsCache hc.RequestCache[PendingSecretsRequest],
	allowlist OnchainAllowlist,
	userRateLimiter *hc.RateLimiter,
	nodeRateLimiter *hc.RateLimiter,
	lggr logger.Logger) handlers.Handler {
	return &functionsHandler{
		handlerConfig:   cfg,
		donConfig:       donConfig,
		don:             don,
		pendingRequests: pendingRequestsCache,
		allowlist:       allowlist,
		userRateLimiter: userRateLimiter,
		nodeRateLimiter: nodeRateLimiter,
		chStop:          make(utils.StopChan),
		lggr:            lggr,
	}
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	sender := common.HexToAddress(msg.Body.Sender)
	if h.allowlist != nil && !h.allowlist.Allow(sender) {
		h.lggr.Debugw("received a message from a non-allowlisted address", "sender", msg.Body.Sender)
		return errors.New("sender not allowlisted")
	}
	if h.userRateLimiter != nil && !h.userRateLimiter.Allow(msg.Body.Sender) {
		h.lggr.Debug("rate-limited", "sender", msg.Body.Sender)
		return errors.New("rate-limited")
	}
	switch msg.Body.Method {
	case MethodSecretsSet, MethodSecretsList:
		return h.handleSecretsRequest(ctx, msg, callbackCh)
	default:
		h.lggr.Debug("unsupported method", "method", msg.Body.Method)
		return errors.New("unsupported method")
	}
}

func (h *functionsHandler) handleSecretsRequest(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.lggr.Debugw("handleSecretsRequest: processing message", "sender", msg.Body.Sender, "messageId", msg.Body.MessageId)
	err := h.pendingRequests.NewRequest(msg, callbackCh, &PendingSecretsRequest{request: msg, responses: make(map[string]*api.Message)})
	if err != nil {
		h.lggr.Warnw("handleSecretsRequest: error adding new request", "sender", msg.Body.Sender, "err", err)
		return err
	}
	// Send to all nodes.
	for _, member := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, member.Address, msg)
		if err != nil {
			h.lggr.Debugw("handleSecretsRequest: failed to send to a node", "node", member.Address, "err", err)
		}
	}
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("HandleNodeMessage: processing message", "nodeAddr", nodeAddr, "receiver", msg.Body.Receiver, "id", msg.Body.MessageId)
	if h.nodeRateLimiter != nil && !h.nodeRateLimiter.Allow(nodeAddr) {
		h.lggr.Debug("rate-limited", "sender", nodeAddr)
		return errors.New("rate-limited")
	}
	switch msg.Body.Method {
	case MethodSecretsSet, MethodSecretsList:
		return h.pendingRequests.ProcessResponse(msg, h.processSecretsResponse)
	default:
		h.lggr.Debug("unsupported method", "method", msg.Body.Method)
		return errors.New("unsupported method")
	}
}

// Conforms to ResponseProcessor[*PendingSecretsRequest]
func (h *functionsHandler) processSecretsResponse(response *api.Message, responseData *PendingSecretsRequest) (*handlers.UserCallbackPayload, *PendingSecretsRequest, error) {
	if _, exists := responseData.responses[response.Body.Sender]; exists {
		return nil, nil, errors.New("duplicate response")
	}
	responseData.responses[response.Body.Sender] = response
	if response.Body.Method != responseData.request.Body.Method {
		return nil, responseData, errors.New("invalid method")
	}
	var responsePayload SecretsResponseBase
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
	payload := CombinedSecretsResponse{SecretsResponseBase: SecretsResponseBase{Success: success}, NodeResponses: responses}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	userResponse := *request
	userResponse.Body.Receiver = request.Body.Sender
	userResponse.Body.Payload = payloadJson
	return &handlers.UserCallbackPayload{Msg: &userResponse, ErrCode: api.NoError, ErrMsg: ""}, nil
}

func (h *functionsHandler) Start(ctx context.Context) error {
	return h.StartOnce("FunctionsHandler", func() error {
		h.lggr.Info("starting FunctionsHandler")
		if h.allowlist != nil {
			return h.allowlist.Start(ctx)
		}
		return nil
	})
}

func (h *functionsHandler) Close() error {
	return h.StopOnce("FunctionsHandler", func() (err error) {
		close(h.chStop)
		if h.allowlist != nil {
			return h.allowlist.Close()
		}
		return nil
	})
}
