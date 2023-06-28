package functions

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type FunctionsHandlerConfig struct {
	OnchainAllowlistChainID string `json:"onchainAllowlistChainId"`
	// Not specifying OnchainAllowlist config disables allowlist checks
	OnchainAllowlist *OnchainAllowlistConfig `json:"onchainAllowlist"`
}

type functionsHandler struct {
	handlerConfig *FunctionsHandlerConfig
	donConfig     *config.DONConfig
	don           handlers.DON
	allowlist     OnchainAllowlist
	lggr          logger.Logger
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, chains evm.ChainSet, lggr logger.Logger) (handlers.Handler, error) {
	cfg, err := ParseConfig(handlerConfig)
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
	return &functionsHandler{
		handlerConfig: cfg,
		donConfig:     donConfig,
		don:           don,
		allowlist:     allowlist,
		lggr:          lggr,
	}, nil
}

func ParseConfig(handlerConfig json.RawMessage) (*FunctionsHandlerConfig, error) {
	var cfg FunctionsHandlerConfig
	if err := json.Unmarshal(handlerConfig, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	if err := msg.Validate(); err != nil {
		h.lggr.Debugw("received invalid message", "err", err)
		return err
	}
	sender := common.HexToAddress(msg.Body.Sender)
	if h.allowlist != nil && !h.allowlist.Allow(sender) {
		h.lggr.Debugw("received a message from a non-allowlisted address", "sender", msg.Body.Sender)
		return errors.New("sender not allowlisted")
	}
	h.lggr.Debugw("received a valid message", "sender", msg.Body.Sender)
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	return nil
}

func (h *functionsHandler) Start(ctx context.Context) (err error) {
	if h.allowlist != nil {
		err = h.allowlist.Start(ctx)
	}
	return
}

func (h *functionsHandler) Close() (err error) {
	if h.allowlist != nil {
		err = h.allowlist.Close()
	}
	return
}
