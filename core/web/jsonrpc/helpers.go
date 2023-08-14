package jsonrpc

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	InternalServerErrorMsg string = "Internal server error"
)

// JsonRpcError writes JSON error response back to client over HTTP
func JsonRpcError(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"jsonrpc": "2.0",
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// JsonRpcResponse writes JSON successful response back to client, unless there are deserialization errors
func JsonRpcResponse(c *gin.Context, lggr logger.Logger, response any) {
	c.JSON(http.StatusOK, gin.H{
		"result":  response,
		"jsonrpc": "2.0",
	})
}

type Error struct {
	Code    int
	Message string
}

const (
	ParseError          int = -32700
	InternalError       int = -32603
	InvalidRequestError int = -32600
)
