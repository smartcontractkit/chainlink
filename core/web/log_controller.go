package web

import (
	"fmt"
	"net/http"
	"strconv"

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

type LogPatchRequest struct {
	Level      string `json:"level"`
	SqlEnabled *bool  `json:"sqlEnabled"`
}

// SetDebug sets the debug log mode for the logger
func (cc *LogController) SetDebug(c *gin.Context) {
	request := &LogPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if request.Level == "" && request.SqlEnabled == nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("please set either logLevel or logSql as params in order to set the log level"))
		return
	}

	if request.Level != "" {
		var ll zapcore.Level
		err := ll.UnmarshalText([]byte(request.Level))
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().Config.Set("LOG_LEVEL", ll.String())
		err = cc.App.GetStore().SetConfigStrValue("LogLevel", ll.String())
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	if request.SqlEnabled != nil {
		cc.App.GetStore().Config.Set("LOG_SQL", request.SqlEnabled)
		err := cc.App.GetStore().SetConfigStrValue("LogSQLStatements", strconv.FormatBool(*request.SqlEnabled))
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().SetLogging(*request.SqlEnabled)
	}

	// Set default logger with new configurations
	logger.Default = cc.App.GetStore().Config.CreateProductionLogger()

	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		Level:      cc.App.GetStore().Config.LogLevel().String(),
		SqlEnabled: cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}
