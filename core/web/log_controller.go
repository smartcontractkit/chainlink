package web

import (
	"net/http"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

// LogController manages the logger config
type LogController struct {
	App chainlink.Application
}

type loglevelPatchRequest struct {
	EnableDebugLog *bool `json:"isDebugEnabled"`
}

type loglevelPatchResponse struct {
	IsDebugEnabled bool `json:"isDebugEnabled"`
}

// ToggleDebug toggles the debug log mode
func (cc *LogController) ToggleDebug(c *gin.Context) {
	request := &loglevelPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if *request.EnableDebugLog {
		cc.App.GetStore().Config.Set("LOG_LEVEL", zapcore.DebugLevel.String())
	} else {
		cc.App.GetStore().Config.Set("LOG_LEVEL", zapcore.InfoLevel.String())
	}
	logger.SetLogger(cc.App.GetStore().Config.CreateProductionLogger())

	response := &loglevelPatchResponse{
		IsDebugEnabled: *request.EnableDebugLog,
	}
	jsonAPIResponse(c, response, "log")
}
