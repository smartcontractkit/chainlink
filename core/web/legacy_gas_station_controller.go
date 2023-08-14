package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/web/jsonrpc"
)

type LegacyGasStationController struct {
	requestRouter legacygasstation.RequestRouter
	lggr          logger.Logger
}

// SendTransaction is a JSON RPC endpoint that submits meta transaction on-chain
func (c *LegacyGasStationController) SendTransaction(ctx *gin.Context) {
	var req types.SendTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.lggr.Warnw("parse error", "err", err)
		jsonrpc.JsonRpcError(ctx, jsonrpc.ParseError, fmt.Sprintf("Parse error: %s", err.Error()))
		return
	}

	resp, jsonRPCError := c.requestRouter.SendTransaction(ctx, req)
	if jsonRPCError != nil {
		jsonrpc.JsonRpcError(ctx, jsonRPCError.Code, jsonRPCError.Message)
		return
	}
	// should not happen
	if resp == nil {
		jsonrpc.JsonRpcError(ctx, jsonrpc.InternalError, jsonrpc.InternalServerErrorMsg)
		return
	}
	jsonrpc.JsonRpcResponse(ctx, c.lggr, *resp)
}
