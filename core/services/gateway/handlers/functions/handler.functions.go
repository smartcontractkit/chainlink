package functions

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type FunctionsHandlerConfig struct {
	AllowlistCheckEnabled       bool   `json:"allowlistCheckEnabled"`
	AllowlistChainID            int64  `json:"allowlistChainID"`
	AllowlistContractAddress    string `json:"allowlistContractAddress"`
	AllowlistBlockConfirmations int64  `json:"allowlistBlockConfirmations"`
	AllowlistUpdateFrequencySec int    `json:"allowlistUpdateFrequencySec"`
	AllowlistUpdateTimeoutSec   int    `json:"allowlistUpdateTimeoutSec"`
}

type functionsHandler struct {
	handlerConfig     *FunctionsHandlerConfig
	donConfig         *config.DONConfig
	don               handlers.DON
	allowlist         OnchainAllowlist
	serviceContext    context.Context
	serviceCancel     context.CancelFunc
	shutdownWaitGroup sync.WaitGroup
	lggr              logger.Logger
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, chains evm.ChainSet, lggr logger.Logger) (handlers.Handler, error) {
	cfg, err := ParseConfig(handlerConfig)
	if err != nil {
		return nil, err
	}
	var allowlist OnchainAllowlist
	if cfg.AllowlistCheckEnabled {
		chain, err := chains.Get(big.NewInt(cfg.AllowlistChainID))
		if err != nil {
			return nil, err
		}
		allowlist, err = NewOnchainAllowlist(chain.Client(), common.HexToAddress(cfg.AllowlistContractAddress), cfg.AllowlistBlockConfirmations, lggr)
		if err != nil {
			return nil, err
		}
	}
	serviceContext, serviceCancel := context.WithCancel(context.Background())
	return &functionsHandler{
		handlerConfig:  cfg,
		donConfig:      donConfig,
		don:            don,
		allowlist:      allowlist,
		serviceContext: serviceContext,
		serviceCancel:  serviceCancel,
		lggr:           lggr,
	}, nil
}

func ParseConfig(handlerConfig json.RawMessage) (*FunctionsHandlerConfig, error) {
	var cfg FunctionsHandlerConfig
	if err := json.Unmarshal(handlerConfig, &cfg); err != nil {
		return nil, err
	}
	if cfg.AllowlistCheckEnabled {
		if !common.IsHexAddress(cfg.AllowlistContractAddress) {
			return nil, errors.New("allowlistContractAddress is not a valid hex address")
		}
		if cfg.AllowlistUpdateFrequencySec <= 0 {
			return nil, errors.New("allowlistUpdateFrequencySec must be positive")
		}
		if cfg.AllowlistUpdateTimeoutSec <= 0 {
			return nil, errors.New("allowlistUpdateTimeoutSec must be positive")
		}
	}
	return &cfg, nil
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	if err := msg.Validate(); err != nil {
		h.lggr.Debug("received invalid message", "err", err)
		return err
	}
	sender := common.HexToAddress(msg.Body.Sender)
	if h.allowlist != nil && !h.allowlist.Allow(sender) {
		h.lggr.Debug("received a message from a non-allowlisted address", "sender", msg.Body.Sender)
		return errors.New("sender not allowlisted")
	}
	h.lggr.Debug("received a valid message", "sender", msg.Body.Sender)
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	return nil
}

func (h *functionsHandler) Start(ctx context.Context) error {
	if h.allowlist != nil {
		checkFreq := time.Duration(h.handlerConfig.AllowlistUpdateFrequencySec) * time.Second
		checkTimeout := time.Duration(h.handlerConfig.AllowlistUpdateTimeoutSec) * time.Second
		h.shutdownWaitGroup.Add(1)
		go func() {
			h.allowlist.UpdatePeriodically(h.serviceContext, checkFreq, checkTimeout)
			h.shutdownWaitGroup.Done()
		}()
	}
	return nil
}

func (h *functionsHandler) Close() error {
	h.serviceCancel()
	h.shutdownWaitGroup.Wait()
	return nil
}
