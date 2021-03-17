package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"go.uber.org/zap/zapcore"
)

// LogController manages the logger config
type LogController struct {
	App chainlink.Application
}

type LoglevelPatchRequest struct {
	EnableDebugLog *bool `json:"debugEnabled"`
}

// SetDebug sets the debug log mode for the logger
func (cc *LogController) SetDebug(c *gin.Context) {
	request := &LoglevelPatchRequest{}
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

	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		DebugEnabled: *request.EnableDebugLog,
	}

	jsonAPIResponse(c, response, "log")
}
