package legacygasstation

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/jsonrpc"
)

type (
	handler interface {
		SendTransaction(*gin.Context, types.SendTransactionRequest) (*types.SendTransactionResponse, *jsonrpc.Error)
		CCIPChainSelector() *utils.Big
	}
)

type requestRouter struct {
	handlersMu sync.RWMutex
	handlers   map[uint64]handler // mapping from ccip chain selector to request handler
	lggr       logger.Logger
}

func NewRequestRouter(lggr logger.Logger) *requestRouter {
	handlers := make(map[uint64]handler)
	return &requestRouter{
		handlers:   handlers,
		handlersMu: sync.RWMutex{},
		lggr:       lggr,
	}
}

func (m *requestRouter) RegisterHandler(h handler) error {
	if h.CCIPChainSelector() == nil {
		return errors.New("empty ccip chain selector")
	}
	m.handlersMu.Lock()
	defer m.handlersMu.Unlock()
	m.handlers[h.CCIPChainSelector().ToInt().Uint64()] = h
	return nil
}

func (m *requestRouter) DeregisterHandler(ccipChainSelector *utils.Big) {
	m.handlersMu.Lock()
	defer m.handlersMu.Unlock()
	if _, ok := m.handlers[ccipChainSelector.ToInt().Uint64()]; !ok {
		m.lggr.Warnw("ccip chain selector not found in deregisterHandler", "ccipChainSelector", ccipChainSelector)
	}
	delete(m.handlers, ccipChainSelector.ToInt().Uint64())
}

func (m *requestRouter) SendTransaction(ctx *gin.Context, req types.SendTransactionRequest) (*types.SendTransactionResponse, *jsonrpc.Error) {
	m.handlersMu.Lock()
	defer m.handlersMu.Unlock()
	handler, ok := m.handlers[req.SourceChainID]
	if !ok {
		m.lggr.Warnw("source ccip chain selector not found", "ccipChainSelector", req.SourceChainID)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InvalidRequestError,
			Message: fmt.Sprintf("unsupported source chain ID: %d", req.SourceChainID),
		}
	}
	return handler.SendTransaction(ctx, req)
}
