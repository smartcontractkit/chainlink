package web

import (
	"fmt"
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

type LogPatchRequest struct {
	Level      string `json:"level"`
	SqlEnabled *bool  `json:"sqlEnabled"`
}

// Get retrieves the current log config settings
func (cc *LogController) Get(c *gin.Context) {
	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		Level:      cc.App.GetStore().Config.LogLevel().String(),
		SqlEnabled: cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}

// Patch sets a log level and enables sql logging for the logger
func (cc *LogController) Patch(c *gin.Context) {
	request := &LogPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if request.Level == "" && request.SqlEnabled == nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("please set either logLevel or logSql as params in order to set the log level"))
		return
	}

	if request.Level != "" {
		var ll zapcore.Level
		err := ll.UnmarshalText([]byte(request.Level))
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		if err = cc.App.GetStore().Config.SetLogLevel(c.Request.Context(), ll.String()); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	if request.SqlEnabled != nil {
		if err := cc.App.GetStore().Config.SetLogSQLStatements(c.Request.Context(), *request.SqlEnabled); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().SetLogging(*request.SqlEnabled)
	}

	// Set default logger with new configurations
	logger.SetLogger(cc.App.GetStore().Config.CreateProductionLogger())

	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		Level:      cc.App.GetStore().Config.LogLevel().String(),
		SqlEnabled: cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}
